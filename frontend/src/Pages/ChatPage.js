import React, { useState } from 'react'
import { ChatState } from '../context/ChatProvider';
import { Box } from '@chakra-ui/react'
import MyChats from '../components/ChatLayout/MyChats'
import ChatBox from '../components/ChatLayout/ChatBox';
import Header from '../components/Header/Header';

const ChatPage = () => {
  const { user, } = ChatState();
  const [fetchChat, setFetchState] = useState(false)

  return (
    <div style={{ width: "100%" }}>
      {user && <Header/>}

      <Box
        display="flex"
        justifyContent="space-between"
        h="92%"
      >
        { user && <MyChats fetchChat={fetchChat}/> }
        { user && <ChatBox />}
      </Box>
    </div>
  )
}

export default ChatPage;