import React from 'react'
import { Box, Text } from '@chakra-ui/react'
import { ChatState } from '../../context/ChatProvider';

const ChatBox = () => {
  const {selectedChat, setSelectedChat, user, chats, setChats} = ChatState();

  function getSenderName(chat) {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].name;
    else 
      return chat.users[0].name;
  }
  

  return (
    <>
      {selectedChat ? 
        <Box
          display={{base: selectedChat? "none" : "flex", md: "flex"}}
          flexDir="column"
          bg="white"
          w={{base: "100%", md: "70%"}}
          borderRadius="10px"
          borderWidth="1px"
          p="5px 10px 5px 10px"
          m="5px"
        >
          {selectedChat &&
            <Box
              display="flex"
              flexDir="row"
              justifyContent="space-between"
              p="10px"
              bg="gray.300"
              borderRadius="xl"
            >
              <Text fontSize="2xl">
                {selectedChat.isGroupChat? selectedChat.chatName: getSenderName(selectedChat)}
              </Text>
            </Box>  
          }
          <Box
            display="flex"
            flexDir="column"
            bg="gray.200"
            h="100%"
          >
            
          </Box>
        </Box>
        : <Box
            display={{base: selectedChat? "none" : "flex", md: "flex"}}
            flexDir="column"
            bg="white"
            w={{base: "100%", md: "70%"}}
            borderRadius="10px"
            borderWidth="1px"
            p="5px"
            m="5px"
          >
            <Text display="flex" fontSize="3xl" justifyContent="center">Please select a user to chat</Text>
          </Box>
      }
    </>
    
  )
}

export default ChatBox