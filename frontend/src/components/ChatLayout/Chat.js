import React, { useEffect, useState } from 'react'
import { Box, Text, Icon, IconButton, Spinner, FormControl, Input, useToast, MenuButton, MenuList, MenuItem, Menu } from '@chakra-ui/react'
import { SettingsIcon } from '@chakra-ui/icons'
import { ChatState } from '../../context/ChatProvider';
import { ArrowBackIcon } from '@chakra-ui/icons'
import Profile from '../Header/Profile';
import UpdateGroupChat from '../utils/UpdateGroupChat';
import axios from 'axios'
import MessagesComp from './MessagesComp';

var socket = new WebSocket(`ws://${process.env.WDS_SOCKET_HOST}:${process.env.WDS_SOCKET_PORT}${process.env.WDS_SOCKET_PATH}`);

const Chat = ({fetchAgain, setFetchAgain}) => {
  const {selectedChat, setSelectedChat, user} = ChatState();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [newMessage, setNewMessage] = useState("");
  const [editingMessageId, setEditingMessageId] = useState("");
  const toast = useToast();

  function getSenderName(chat) {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].name;
    else 
      return chat.users[0].name;
  }

  function getSenderObj(chat){
    if (chat.users[0]._id === user._id) 
      return chat.users[1];
    else 
      return chat.users[0];
  }

  const handleEditMessage = (messageId) => {
    setEditingMessageId(messageId);

    // place the message in input element for user to edit 
    setNewMessage(messages.find((objId) => objId._id===messageId).content);
  }

  const handleSendMessage = (event) => {
    if (event.key === "Enter" && newMessage) {
      const content = newMessage.trim();
      if (editingMessageId !== "") {
        editMessage(content);
        setEditingMessageId("");
      } else {
        sendMessage(content);
      }
    }
  }

  const sendMessage = async(content) => {
      try {
        const { data } = await axios.post("/api/message/", 
          {
            content: content,
            chatId: selectedChat._id,
          },
          {
            headers: {
              "Content-Type": "application/json",
              "Authorization": `Bearer ${user.token}`
            }
          }
        )
        // send the message to server, so that server will broadcast it
        sendmessage(JSON.stringify(data));
        setNewMessage("");
      } catch (error) {
          toast({
            title: "Failed to send message",
            status: "error",
            duration: 4000,
            isClosable: true,
            position: "botton",
          });
      }
  }

  const editMessage = async(content) => {
      try {
        const { data } = await axios.put("/api/message/", 
          {
            content: content,
            messageId: editingMessageId,
          },
          {
            headers: {
              "Content-Type": "application/json",
              "Authorization": `Bearer ${user.token}`
            }
          }
        )
        // send the message to server, so that server will broadcast it
        sendmessage(JSON.stringify(data));
        toast({
          title: "Message edited",
          status: "success",
          duration: 3000,
          isClosable: true,
          position: "top",
        });
        setNewMessage("");
      } catch (error) {
          toast({
            title: "Failed to edit message",
            status: "error",
            duration: 4000,
            isClosable: true,
            position: "botton",
          });
      }
  }

  const deleteMessage = async(messageId) => {
    try {
      const { data } = await axios.delete(`/api/message/${messageId}`, 
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${user.token}`
          }
        }
      )

      // send websocket server to delete this message on all active users chats
      sendmessage(JSON.stringify({"_id":messageId, "delete": true, "chat": selectedChat._id}));
      toast({
        title: "Message deleted",
        status: "success",
        duration: 3000,
        isClosable: true,
        position: "top",
      });
    } catch (error) {
      toast({
        title: "Message not deleted",
        status: "error",
        duration: 4000,
        isClosable: true,
        position: "botton",
      });
    }
  }

  const deleteConversation = async(chatId) => {
    try {
      const _ = await axios.delete(`/api/chat/${chatId}`, 
        {
          headers: {
            "Content-Type" : "application/json",
            "Authorization": `Bearer ${user.token}`
          }
        }
      )
      setSelectedChat("");
      setFetchAgain(!fetchAgain);
      toast({
        title: "Chat deleted successfully",
        status: "success",
        duration: 4000,
        isClosable: true,
        position: "botton",
      });
    } catch (error) {
        toast({
          title: "Failed to delete chat",
          status: "error",
          duration: 4000,
          isClosable: true,
          position: "botton",
        });
    }
  }

  const fetchMessages = async() => {
    try {
      setLoading(true);
      const { data } = await axios.get(`/api/message/${selectedChat._id}`, 
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${user.token}`
          }
        }
      )
      console.log(data);
      setLoading(false);
      if (data == null){
        setMessages([]);
      } else {
        setMessages(data);
      }

      sendmessage(JSON.stringify({"messageType":"setup", "chat": selectedChat._id}));
    } catch (error) {
        toast({
          title: "Failed to load messages",
          status: "error",
          duration: 4000,
          isClosable: true,
          position: "botton",
        });
        setLoading(false);
    }
  }

  // fetch messages upon changing the selectedChat, 
  // means user switched to other person to chat with 
  useEffect(() => {
    fetchMessages();
  }, [selectedChat])
  
  const addMessage = (msg) => {
    if (msg.delete) {
      setMessages(messages.filter((obj)=> obj._id !== msg._id));
    }
    else if (msg.isedited) {
      let index = 0;
        for(let i=0; i<messages.length; i++) {
          if (messages[i]._id === msg._id){
            index = i;
            break;
          }
        }
        setMessages([...messages.slice(0, index), msg, ...messages.slice(index+1, messages.length)]);
    } else {
      setMessages([...messages, msg]);
    }
  }

  const sendmessage = (msg) => {
    socket.send(msg);
  }

  // need to handle the reconnection to socket if any interruption 
  // occured to websocket connection
  useEffect(() => {
    window.client = socket;

    socket.onopen = () => {
      console.log("Successfully connected");
    }

    socket.onmessage = (msg) => {
        addMessage(JSON.parse(msg.data));
    }

    socket.onclose = (event) => {
        console.log("Socket closed connection: ", event);
    }

    socket.onerror = (error) => {
        console.log("Socket error: ", error);
    }
  })

  return (
    <>
      <Box 
        display="flex"
        flexDir="row"
        justifyContent="space-between"
        alignItems="center"
        p="10px"
        bg="white"
        w="100%"
        h="10%"
        borderRadius="xl"
        boxShadow="0 5px 10px gray"
        
      >
        <IconButton 
          display={{base: "flex", md: "none"}}
          icon={<ArrowBackIcon/>}
          onClick={() => setSelectedChat("")}
        />
        
        <Text fontSize="2xl">
          {selectedChat.isGroupChat? selectedChat.chatName: getSenderName(selectedChat)}
        </Text>

        {!selectedChat.isGroupChat? 
            <div>
              <Menu>
                <MenuButton as={IconButton} icon={<SettingsIcon/>} />
                <MenuList minWidth="150px" >
                  <Profile user={getSenderObj(selectedChat)}>
                    <MenuItem>User Profile</MenuItem>
                  </Profile>
                  <MenuItem onClick={() => {window.confirm("Sure you want to delete user, all conversation will be lost?") 
                                                    && deleteConversation(selectedChat._id)} }>Delete</MenuItem>
                </MenuList>
              </Menu>
            </div>
            : 
            <UpdateGroupChat fetchAgain={fetchAgain} setFetchAgain={setFetchAgain}>
               <Icon icon={<SettingsIcon/>}
                display="flex"
              />
            </UpdateGroupChat>
        }
      </Box>  

      <Box 
        display="flex" 
        flexDir="column" 
        h="100%" 
        w="100%" 
        flex={1} p={1}
        bg="#E8E8E8" 
        borderRadius="xl"
        mt="2px"
        boxShadow="0 5px 10px black"
        height="90%"
      >
        {loading? 
          (<Spinner size="lg" alignSelf="center" margin="auto"/>)
            : <div style={{ display:"flex", flexDirection:"column-reverse", overflowY:"auto", padding: "3px"}}>
                <MessagesComp messages={messages} deleteMessage={deleteMessage} handleEditMessage={handleEditMessage}/>
              </div>}

        <FormControl onKeyDown={handleSendMessage} isRequired mt={3} marginTop="auto" marginBottom={1}>
          <Input 
            placeholder="Enter your text" 
            bg="#E0E0E0"
            value={newMessage} 
            onChange={(e) => setNewMessage(e.target.value)}/>
        </FormControl>
      </Box>
  </>

  )
}
export default Chat