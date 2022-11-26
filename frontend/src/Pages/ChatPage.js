import React, { useState } from 'react'
import { ChatState } from '../context/ChatProvider';
import { Box } from '@chakra-ui/react'
import MyChats from '../components/ChatLayout/MyChats'
import ChatBox from '../components/ChatLayout/ChatBox';
import Header from '../components/Header/Header';

const ChatPage = () => {
  const { user, } = ChatState();
  const [fetchAgain, setFetchAgain] = useState(false)

  return (
    <div style={{ width: "100%" }}>
      {user && <Header/>}

      <Box
        display="flex"
        justifyContent="space-between"
        h="92%"
      >
        { user && <MyChats fetchAgain={fetchAgain}/> }
        { user && <ChatBox fetchAgain={fetchAgain} setFetchAgain={setFetchAgain}/>}
      </Box>
    </div>
  )
}

export default ChatPage;