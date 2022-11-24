import React, { useEffect, useState } from 'react'
import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, useDisclosure, VStack} from '@chakra-ui/react'
import { FormControl, Input, Button, Box, useToast, Spinner, Stack} from '@chakra-ui/react'
import { ChatState } from '../../context/ChatProvider'
import UserBadge from './UserBadge'
import axios from "axios"
import UserSearchProfile from './UserSearchProfile'

const UpdateGroupChat = ({user, children}) => {
  
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [groupName, setGroupName] = useState("");
  const { selectedChat } = ChatState();
  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState([]);

  const toast = useToast();
  
  const handleAddUserToGroup = async(userToAdd) => {
    if (selectedUsers.includes(userToAdd)) {
        toast({
            title: "User already added",
            status: "warning",
            duration: 4000,
            isClosable: true,
            position: "top",
        });
        return;
    }
    
    try {
      
    } catch (error) {

    }

  }

  const handleRemoveUserFromGroup = (userToBeRemoved) => {
    try {
      
    } catch (error) {

    }
  }

  const handleSearch = async(search) => {
    search = search.trim();
    if (!search){
      return;
    }

    try {
      setLoading(true);

      const { data } = await axios.get(
        `http://localhost:8000/api/user/search?search=${search}`,
        {
          headers: {
            "Content-Type":"application/json",
          }
        }
      );
      console.log(data)
      setLoading(false);
      setSearchResults(data);
    } catch (error) {
        console.log(error.response.data.error)
        toast({
          title: "Error occured",
          decreption: "Failed to load user data",
          duration: 3000,
          status: "warning",
          isClosable: true,
          position: "bottom-left",
        });
        setLoading(false);
    }
  }

  return (
    <>
      <span onClick={onOpen}>{children}</span>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader
            fontSize="25px"
            fontFamily="Work sans"
            display="flex"
            justifyContent="center"
          >
            Update Group Chat
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody
            display="flex" 
            flexDir="column" 
            alignItems="center"
          >
            <Box
              w="100%"
              display="flex"
              flexWrap="wrap"
            >
              {selectedChat.users.map((u) => (
                <UserBadge 
                  key={u._id}
                  handleFunction={() => handleRemoveUserFromGroup(u)}
                  user={u}
                />
              ))}
            </Box>

            <FormControl display="flex" flexDir="row">
              <Input 
                placeholder="Group Name" 
                onChange={(e) => setGroupName(e.target.value)}
               />
              <Button ml="5px">
                Update 
              </Button>
            </FormControl>
            <FormControl>
                <Input placeholder="Add user"
                  onChange={(e) => handleSearch(e.target.value)}
                />
            </FormControl>

           
            {loading ? 
              <Spinner/>
                : <Stack w="100%">
                    { searchResults?.slice(0,4).map((user) => (
                        <UserSearchProfile 
                          key={user._id}
                          handleFunction={() => handleAddUserToGroup(user)}
                          user={user}
                        />
                    )) }
                  </Stack>
            }
           
          </ModalBody>

          <ModalFooter>
            <Button colorScheme={"orange"}>Exit Group</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

export default UpdateGroupChat