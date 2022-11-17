import React from 'react';
import {Box, Image} from '@chakra-ui/react';


const UserSearchProfile = ({user}) => {
  return (
    <div>
      <Box 
        backgroundColor={"#ede7d8"}
        mb="10px"
        borderRadius={"8px"}
        p={"5px"}
      >
        <Box display={"flex"} justifyContent="space-between">
          {user.name} 
          <Image width={"30px"} borderRadius="full" src={user.pic}/>
        </Box>
        {user.email}
            
        </Box>
    </div>
  )
}

export default UserSearchProfile