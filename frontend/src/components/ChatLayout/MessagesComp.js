import { Text, Box, useToast, Menu, MenuButton, MenuList, MenuItem } from '@chakra-ui/react'
import React from 'react'
import { ChatState } from '../../context/ChatProvider'
import { isSenderTheLoggedInUser } from '../utils/messagesRendering'

const MessagesComp = ({ messages, deleteMessage, handleEditMessage }) => {
    const {user} = ChatState();

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
                  <MenuItem onClick={() => handleEditMessage(m._id)}>Edit Message</MenuItem>
                </MenuList>
              </Menu>
              /*<IconButton colorScheme="none" size="20px" 
                icon={<ChevronDownIcon width="15px" color="blackAlpha.500"/>}
                onClick={() => deleteMessage(m._id)}
              /> */
              : <Text> {m.content} </Text>}
            </Box>
        </div>))}
    </>
  )
}

export default MessagesComp