import React from 'react'
import { Box, Text, Button, Icon, IconButton } from '@chakra-ui/react'
import {InfoIcon} from '@chakra-ui/icons'
import { ChatState } from '../../context/ChatProvider';
import { ArrowBackIcon } from '@chakra-ui/icons'
import Profile from '../Header/Profile';
import UpdateGroupChat from '../utils/UpdateGroupChat';

const Chat = () => {
  const {selectedChat, setSelectedChat, user, chats, setChats} = ChatState();
  
  function getSenderName(chat) {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].name;
    else 
      return chat.users[0].name;
  }
    
  function getSenderObj(chat) {
    if (chat.users[0]._id === user._id) 
      return chat.users[1];
    else 
      return chat.users[0];
  }
  return (
    <>
      <Box 
        display="flex"
        flexDir="row"
        justifyContent="space-between"
        alignItems={"center"}
        p="10px"
        bg="white"
        w="100%"
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
            <Profile user={selectedChat}>
              <Icon icon={<InfoIcon />}
                display="flex"
              />
            </Profile>
            : 
            <UpdateGroupChat user={selectedChat}>
              <Icon icon={<InfoIcon />}
                display="flex"
              />
            </UpdateGroupChat>
        }
      </Box>  

      <Box>

      </Box>
  </>
  )
}

export default Chat