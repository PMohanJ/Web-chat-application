import React from 'react'
import { Box, Text, Icon, IconButton } from '@chakra-ui/react'
import {InfoIcon} from '@chakra-ui/icons'
import { ChatState } from '../../context/ChatProvider';
import { ArrowBackIcon } from '@chakra-ui/icons'
import Profile from '../Header/Profile';
import UpdateGroupChat from '../utils/UpdateGroupChat';

const Chat = ({fetchAgain, setFetchAgain}) => {
  const {selectedChat, setSelectedChat, user} = ChatState();
  
  function getSenderName(chat) {
    if (chat.users[0]._id === user._id) 
      return chat.users[1].name;
    else 
      return chat.users[0].name;
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
            <UpdateGroupChat fetchAgain={fetchAgain} setFetchAgain={setFetchAgain}>
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