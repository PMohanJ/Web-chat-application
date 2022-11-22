import React from 'react'
import { ChatState } from '../context/ChatProvider';
import { Box } from '@chakra-ui/react'
import MyChats from '../components/ChatLayout/MyChats'
import ChatBox from '../components/ChatLayout/ChatBox';
import Header from '../components/Header/Header';

const ChatPage = () => {
  const { user, } = ChatState();

  return (
    <div style={{ width: "100%" }}>
      {user && <Header/>}

      <Box
        display="flex"
        justifyContent="space-between"
        h="92%"
      >
        { user && <MyChats/> }
        { user && <ChatBox/>}
      </Box>
    </div>
  )
}

export default ChatPage;