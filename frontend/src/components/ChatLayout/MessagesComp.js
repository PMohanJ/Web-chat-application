import { IconButton, Button, Text, Box, useToast, Menu, MenuButton, MenuList, MenuItem } from '@chakra-ui/react'
import React from 'react'
import { ChatState } from '../../context/ChatProvider'
import { isSenderTheLoggedInUser } from '../utils/messagesRendering'
import axios from "axios"

const MessagesComp = ({ messages, setMessages }) => {
    const {user} = ChatState();
    const toast = useToast();
  
    const deleteMessage = async(messageId) => {
      try {
        const url = `http://localhost:8000/api/message/${messageId}`
        const { data } = await axios.delete(url, 
          {
            headers: {
              "Content-Type": "application/json",
              "Authorization": `Bearer ${user.token}`
            }
          }
        )
        setMessages(messages.filter((obj)=> obj._id != messageId))
        toast({
          title: "Message deleted",
          status: "success",
          duration: 4000,
          isClosable: true,
          position: "botton",
        });
      } catch (error) {
        toast({
          title: "Message not deleted",
          status: "error",
          duration: 4000,
          isClosable: true,
          position: "botton",
        });
      }
    }

    // reversing the order of msg, so that bottom location
    // is prioritized as needed for chatting
    const reveMessages = messages.slice().reverse();
    
  return (
    <>
        {reveMessages && reveMessages.map((m) => (
        <div style={{display: "flex"}} key={m._id}>
            <Box
              backgroundColor={isSenderTheLoggedInUser(m, user._id)? "#BEE3F8": "#B9F5D0"}
              maxWidth="75%"
              p="5px 10px"
              mt="5px"
              borderRadius="5px"
              ml={isSenderTheLoggedInUser(m, user._id)? "auto": 0}
              display="flex"
            > 
           
            {isSenderTheLoggedInUser(m, user._id)? 
              <Menu>
                <MenuButton>
                  <Text> {m.content} </Text>
                </MenuButton>
                <MenuList minWidth="150px">
                  <MenuItem onClick={() => deleteMessage(m._id)}>Delete Message</MenuItem>
                </MenuList>
              </Menu>
              /*<IconButton colorScheme="none" size="20px" 
                icon={<ChevronDownIcon width="15px" color="blackAlpha.500"/>}
                onClick={() => deleteMessage(m._id)}
              /> */
              : null}
            </Box>
        </div>))}
    </>
  )
}

export default MessagesComp