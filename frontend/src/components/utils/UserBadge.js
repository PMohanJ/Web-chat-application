import React from 'react'
import { Badge } from '@chakra-ui/react'
import {CloseIcon} from '@chakra-ui/icons'

const UserBadge = ({user, handleFunction}) => {
  return (
    <>
      <Badge 
        variant='solid' 
        colorScheme='green' 
        onClick={handleFunction} 
        mr="10px"
        mb="5px"
        p="3px"
        fontSize="12px"
        cursor="default"
        borderRadius="5px"
      >
        {user.name}
        <CloseIcon pl="2px" h="10px" w="10px"/>
      </Badge>
    </>
  )
}

export default UserBadge