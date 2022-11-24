import React from 'react'
import { Box, Text } from '@chakra-ui/react'
import { ChatState } from '../../context/ChatProvider';
import SingleChat from './SingleChat';

const ChatBox = () => {
  const {selectedChat} = ChatState();

  return (
    <> 
      <Box
        display={{base: selectedChat? "flex" : "none", md: "flex"}}
        flexDir="column"
        alignItems="center"
        bg="white"
        w={{base: "100%", md: "70%"}}
        borderRadius="10px"
        borderWidth="1px"
      >
        {selectedChat? <SingleChat />
          : <Box display="flex" justifyContent="center" alignItems="center" h="100%">
              <Text fontSize="3xl" fontWeight="light">Please select a user to chat</Text>
            </Box>
        }
        
      </Box>
    </>
  )
}
export default ChatBox