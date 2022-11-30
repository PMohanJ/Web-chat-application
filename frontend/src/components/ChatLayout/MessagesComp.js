import { Box, Text} from '@chakra-ui/react'
import React from 'react'
import { ChatState } from '../../context/ChatProvider'
import { isSenderTheLoggedInUser } from '../utils/messagesRendering'

const MessagesComp = ({ messages }) => {
    const {user} = ChatState();

  return (
    <>
        {messages && messages.map((m) => (
        <div style={{display: "flex"}} key={m._id}>
            <Text
              backgroundColor={isSenderTheLoggedInUser(m, user._id)? "#BEE3F8": "#B9F5D0"}
              maxWidth="75%"
              p="5px 10px"
              mt="5px"
              borderRadius="5px"
              ml={isSenderTheLoggedInUser(m, user._id)? "auto": 0}
            >
                {m.content}
            </Text>
        </div>))}
    </>
  )
}

export default MessagesComp