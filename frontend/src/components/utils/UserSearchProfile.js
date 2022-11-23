import React from 'react';
import {Avatar, Box, Text} from '@chakra-ui/react';


const UserSearchProfile = ({user, handleFunction}) => {
  return (
    <div>
      <Box 
        display={"flex"} 
        flexDir="row"
        justifyContent="space-between"
        bg={"#ede7d8"}
        mb="10px"
        borderRadius={"8px"}
        p={"5px"}
        onClick={handleFunction}
        _hover={{background:"blue.200", color:"white"}}
        color="black"
      >
        <Box >
          <Text fontSize="xl" color={"black"}>
          {user.name} 
          </Text>
          
          <Text color={"blackAlpha.700"}>
            Email: {user.email}
          </Text>
        </Box>

        <Avatar size={"md"} src={user.pic}/>
            
        </Box>
    </div>
  )
}

export default UserSearchProfile