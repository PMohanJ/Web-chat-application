import React, { useEffect, useState } from 'react'
import { Box, Text, Icon, IconButton, Spinner, FormControl, Input, useToast, MenuButton, MenuList, MenuItem, Menu, InputGroup } from '@chakra-ui/react'
import { SettingsIcon } from '@chakra-ui/icons'
import { ChatState } from '../../context/ChatProvider';
import { ArrowBackIcon } from '@chakra-ui/icons'
import Profile from '../Header/Profile';
import UpdateGroupChat from '../utils/UpdateGroupChat';
import axios from 'axios'
import MessagesComp from './MessagesComp';

var socket = new WebSocket(`ws://${process.env.WDS_SOCKET_HOST}:${process.env.WDS_SOCKET_PORT}${process.env.WDS_SOCKET_PATH}`);

const Chat = ({fetchAgain, setFetchAgain}) => {
  const {selectedChat, setSelectedChat, user, latestMessages, setLatestMessages} = ChatState();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [newMessage, setNewMessage] = useState("");
  const [editingMessageId, setEditingMessageId] = useState("");
  const toast = useToast();
  const [isTyping, setIsTyping] = useState(false);
  const [senderTyped, setSenderTyped] = useState(false);

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
      await axios.delete(`/api/message/${messageId}`, 
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
      await axios.delete(`/api/chat/${chatId}`, 
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
  
  let clearTypingId;

  // appends message received from websocket server
  const addMessage = (msg) => {
    if (!msg.typing)  updateLatestMessage(msg)

    // typing obj to indicate user typing
    if (msg.typing) {
      senderTyping(msg)
      setIsTyping(true);
      
      if(clearTypingId) {
      clearTimeout(clearTypingId);
    }
      clearTypingId = setTimeout(() => {
        setIsTyping(false);
      }, 3000); 
    }
    else if (msg.delete) {
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

  const isLastMessage = (msg) => {
    let index = 0;
    for(let i=0; i<messages.length; i++) {
      if (messages[i]._id === msg._id) {
        index = i;
        break;
      }
    }
    
    if (index === messages.length - 1) {
      return true
    }
    return false
  }

  // update latest message upon receiving new msg,
  // edited msg or deleted msg accordingly
  const updateLatestMessage = (msg) => {
    let indexChatId = 0;
    for(let i=0; i<latestMessages.length; i++) {
      if(selectedChat._id === latestMessages[i].chatId){
        indexChatId = i;
        break;
      }
    }

    let islastmessage = isLastMessage(msg)

    let indexDelete = messages.length - 1
    // if previous message exists
    if (messages.length > 1 && msg.delete && islastmessage){
      setLatestMessages([...latestMessages.slice(0, indexChatId), {chatId: selectedChat._id, message: messages[indexDelete - 1].content},
        ...latestMessages.slice(indexChatId+1, latestMessages.length)]);
    } 

    else if(msg.isedited && islastmessage) {
      setLatestMessages([...latestMessages.slice(0, indexChatId), {chatId: selectedChat._id, message: msg.content},
      ...latestMessages.slice(indexChatId+1, latestMessages.length)]);
    }

    else if(!msg.delete && !msg.isedited){
      setLatestMessages([...latestMessages.slice(0, indexChatId), {chatId: selectedChat._id, message: msg.content},
        ...latestMessages.slice(indexChatId+1, latestMessages.length)]);
    }
  }

  // sends websocket msg on user typing
  const handleTyping = (event) => {
    sendmessage(JSON.stringify({"_id":user._id, "typing": true, "chat": selectedChat._id}));
  }

  // to know whether typing object is by sender, coz
  // we don't need to render typing... to sender
  const senderTyping = (msg) => {
    if(msg._id === user._id) {
      setSenderTyped(true)
    } else {
      setSenderTyped(false);
    }
  }

  const sendmessage = (msg) => {
    socket.send(msg);
  }

  // fetch messages upon changing the selectedChat, which
  // means user switched to other person to chat with 
  useEffect(() => {
    fetchMessages();
  }, [selectedChat])

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

        <Box mt={3} marginTop="auto" marginBottom={1}>
          <Text 
            display={isTyping && !senderTyped? "block": "none"}
            color="#14ed05"
            ml={"5px"}
            fontWeight={"extrabold"}
          >
            Typing...
          </Text>
          <FormControl onKeyDown={handleSendMessage} isRequired >
            <InputGroup>
            </InputGroup>
            <Input 
              placeholder="Enter your text" 
              bg="#E0E0E0"
              value={newMessage} 
              onChange={(e) => setNewMessage(e.target.value)}
              onKeyDown={handleTyping}/>
          </FormControl>
        </Box>
        
      </Box>
  </>

  )
}
export default Chat