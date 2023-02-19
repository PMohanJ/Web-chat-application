import { useToast, Box, Text, Button, Stack} from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons'
import React, { useEffect } from 'react'
import { ChatState } from '../../context/ChatProvider';
import axios from "axios"
import GroupChatModel from '../utils/GroupChatModel';

const MyChats = ({ fetchAgain }) => {
  const {selectedChat, setSelectedChat, user, chats, setChats, latestMessages, setLatestMessages} = ChatState();
  const toast = useToast();
  
  const fetchChats = async() => {
    try {
      const { data } = await axios.get(`/api/chat/`, 
        {
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${user.token}`
          }
        }
      )
      console.log(data)
      setChats(data);
    } catch (error) {
        toast({
          title:"Error Occured",
          duration: 4000,
          isClosable: true,
          status: "error",
        });
    }
  }

  const getAllLatestMessages = () => {
    if (chats.length > 0) {
      //console.log("getAllLatestMsg called in func")
      setLatestMessages(chats.map((chat) => ({
        chatId: chat._id,
        message: chat.latestMessage.length > 0 ? chat.latestMessage[0].content: "",
      }))) 
    }
  }

  useEffect(() => {
    fetchChats();
    //console.log("fetchChats is called");
  },[fetchAgain]);

  useEffect(() => {
    //console.log("GetallMsg is called")
    getAllLatestMessages();
  }, [chats])

  const getSenderName = (chat) => {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].name;
    else 
      return chat.users[0].name;
  }
  
  const getSenderPic = (chat) => {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].pic;
    else 
      return chat.users[0].pic;
  }

  const getLatestMessage = (chat) => {
    if (latestMessages.length > 0) {
      const obj = latestMessages.find((obj) => obj.chatId === chat._id)
      if (obj && obj.message) {
        return obj.message;
      }
      return ""
    }
  }

  return (
    <Box
      display={{base: selectedChat? "none" : "flex", md: "flex"}}
      flexDir="column"
      bg="white"
      w={{base: "100%", md: "30%"}}
      borderRadius="10px"
      borderWidth="1px"
      p="5px"
      m="5px"
      minWidth="250px"
    >
      <Box
        display="flex"
        justifyContent="space-between"
        fontSize={"3xl"}
        fontFamily="-moz-initial"
      >
        <Text>Chats</Text>

        <GroupChatModel>
          <Button
            display="flex"
            rightIcon={<AddIcon/>}
            fontFamily="-moz-initial"
            fontWeight="semibold"
            >
            Create a Group
          </Button>
        </GroupChatModel>
      </Box>

      {chats? (
        <Stack >
            {chats.map((chat) => (
              <Box
                bg={selectedChat === chat ? "#62c8f0": "gray.400"}
                color={selectedChat === chat ? "black": "black"}
                key={chat._id}
                borderRadius="5px"
                p="7px"
                onClick={() => setSelectedChat(chat)}
                maxHeight="60px"
                minWidth="180px"
              >
                <Box display="flex" flexDir="row">
                  <img style={{minWidth:"50px", height:"50px", borderRadius: "50%", border: "3px solid #7479e3"}}
                    alt="Profile Pic"
                    src={getSenderPic(chat)}
                  />
                  <Box marginLeft="5px">
                    <Text fontSize="large">
                      {chat.isGroupChat? chat.chatName: getSenderName(chat)}
                    </Text>
                    <p style={{whiteSpace: "nowrap", overflow: "hidden", width:"250px", textOverflow: "ellipsis", fontSize:"14px"}}>
                      {getLatestMessage(chat)}
                    </p>
                  </Box>
                </Box>
              </Box>
            ))}
        </Stack>
      ): <Box >
            <Text>Please add some friends to chat with</Text>
         </Box>
        }
    </Box>
  )
}

export default MyChats