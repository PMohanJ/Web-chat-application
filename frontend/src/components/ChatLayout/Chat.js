import React, { useEffect, useState } from 'react'
import { Box, Text, Icon, IconButton, Spinner, FormControl, Input, useToast } from '@chakra-ui/react'
import {InfoIcon} from '@chakra-ui/icons'
import { ChatState } from '../../context/ChatProvider';
import { ArrowBackIcon } from '@chakra-ui/icons'
import Profile from '../Header/Profile';
import UpdateGroupChat from '../utils/UpdateGroupChat';
import axios from 'axios'
import MessagesComp from './MessagesComp';

const Chat = ({fetchAgain, setFetchAgain}) => {
  const {selectedChat, setSelectedChat, user} = ChatState();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [newMessage, setNewMessage] = useState("");
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
  const sendMessage = async(e) => {
    if (e.key === "Enter" && newMessage) {
      const content = newMessage.trim()
      try {
        const { data } = await axios.post("http://localhost:8000/api/message/", 
          {
            content: content,
            senderId: user._id,
            chatId: selectedChat._id,
          },
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
        )
        console.log(data);
        setNewMessage("");
        setMessages([...messages, data]);

      } catch (error) {
          toast({
            title: "Failed to send message",
            duration: 5000,
            isClosable: true,
            position: "botton",
          })
      }
    }
  }

  const fetchMessages = async() => {
    try {
      setLoading(true);
      const { data } = await axios.get(`http://localhost:8000/api/message/${selectedChat._id}`, 
        {
          headers: {
            "Content-Type": "application/json"
          }
        }
      )
      console.log(data);
      setLoading(false);
      setMessages(data);

    } catch (error) {
        toast({
          title: "Failed to load messages",
          duration: 5000,
          isClosable: true,
          position: "botton",
        })
        setLoading(false);
    }
  }

  // fetches messages upon changing the selectedChat, 
  // means user switched to other person to chat with 
  useEffect(() => {
    fetchMessages()
  }, [selectedChat])
    
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
            <Profile user={getSenderObj(selectedChat)}>
              <Icon icon={<InfoIcon />}
                display="flex"
              />
            </Profile>
            : 
            <UpdateGroupChat fetchAgain={fetchAgain} setFetchAgain={setFetchAgain}>
              <Icon icon={<InfoIcon />}
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
            : <div style={{overflowY: "auto", padding: "3px"}}>
                <MessagesComp messages={messages} />
              </div>}

        <FormControl onKeyDown={sendMessage} isRequired mt={3} marginTop="auto" marginBottom={1}>
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