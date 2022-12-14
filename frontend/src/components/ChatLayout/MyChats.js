import { useToast, Box, Text, Button, Stack} from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons'
import React, { useEffect } from 'react'
import { ChatState } from '../../context/ChatProvider';
import axios from "axios"
import GroupChatModel from '../utils/GroupChatModel';

const MyChats = ({ fetchAgain }) => {
  const {selectedChat, setSelectedChat, user, chats, setChats} = ChatState();
  const toast = useToast();

  const fetchChats = async() => {
    try {
      const { data } = await axios.get(`/api/chat/${user._id}`, 
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

  useEffect(() => {
    fetchChats();
  },[fetchAgain])

  function getSenderName(chat) {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].name;
    else 
      return chat.users[0].name;
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
              >
                <Text fontSize="large">
                  {chat.isGroupChat? chat.chatName: getSenderName(chat)}
                </Text>
                <Text fontSize="small">
                  {chat.latestMessage.length > 0 ? chat.latestMessage[0].content: ""}
                </Text>
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