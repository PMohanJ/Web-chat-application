import { useToast } from '@chakra-ui/react';
import React from 'react'
import { ChatState } from '../../context/ChatProvider';
import { useState } from 'react';
import axios from "axios"

const MyChats = () => {
  const [loggedUser, setLoggedUser] = useState();
  const {selectedChat, setSelectedChat, user, chats, setChats} = ChatState();
  const toast = useToast();

  const fetchChats = async() => {
    try {
      const { data } = await axios.get("http://localhost:8000/api/chat/", 
      {
        headers: {
          "Content-Type":"application/json"
        }
      })
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
  
  return (
    <div>
      MyChats
    </div>
  )
}

export default MyChats