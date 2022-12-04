import React, { useEffect, useState } from 'react'
import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, useDisclosure } from '@chakra-ui/react'
import { FormControl, Input, Button, Box, useToast, Spinner, Stack} from '@chakra-ui/react'
import { ChatState } from '../../context/ChatProvider'
import UserBadge from './UserBadge'
import axios from "axios"
import UserSearchProfile from './UserSearchProfile'

const UpdateGroupChat = ({fetchAgain, setFetchAgain, children}) => {
  
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [groupName, setGroupName] = useState("");
  const { selectedChat, setSelectedChat, user } = ChatState();
  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState([]);

  const toast = useToast();
 
  const handleGroupRename = async() => {

    if (groupName.trim() === "") {
      toast({
        title: "Please trype a name",
        isClosable: true,
        duration: 4000,
        status: "warning",
        position: "bottom"
      });
      return;
    }

    try {
      setLoading(true);
      
      const { data } = await axios.put("http://localhost:8000/api/chat/grouprename",
        {
          chatId: selectedChat._id,
          groupName: groupName,
        },
        {
          headers:{
            "Content-Type": "application/json",
          }
        }
      )
      console.log(data)
      setLoading(false);
      const updatedSelectedChat = JSON.parse(JSON.stringify(selectedChat));
      updatedSelectedChat.chatName = data.updatedGroupName;
      setSelectedChat(updatedSelectedChat)
    } catch (error) {
        toast({
          title: "Error occured",
          decreption: "Failed to rename group",
          duration: 4000,
          status: "warning",
          isClosable: true,
          position: "bottom",
        });
        setLoading(false);
    }
    setGroupName("")
  }
  const handleAddUserToGroup = async(userToBeAdd) => {

    // check if the user already exists in group
    if (selectedUsers.includes(userToBeAdd._id)) {
      toast({
          title: "User already added",
          status: "warning",
          duration: 4000,
          isClosable: true,
          position: "bottom",
      });
      return;
    }
    
    // only an admin can add a user to group
    if (user._id !== selectedChat.groupAdmin) {
      toast({
        title: "Sorry, only admin can add users",
        status: "warning",
        duration: 4000,
        isClosable: true,
        position: "bottom",
      });
      return;
    }

    try {
      const { data } = await axios.put("http://localhost:8000/api/chat/groupadd",
        {
          chatId: selectedChat._id,
          userId: userToBeAdd._id,
        },
        {
          headers: {
            "Content-Type": "application/json",
          }
        }
      )
      console.log(data);
      setSelectedChat(data);
    } catch (error) {
        toast({
          title: "Error occured",
          decription: "Failed to add user",
          duration: 3000,
          status: "warning",
          isClosable: true,
          position: "bottom",
        });
    }

  }

  const handleRemoveUserFromGroup = async(userToBeRemoved) => {

    // only an admin can remove users from group
    if (user._id !== selectedChat.groupAdmin) {
      toast({
        title: "Sorry, only admin can remove user",
        status: "warning",
        duration: 4000,
        isClosable: true,
        position: "bottom",
      });
      return;
    }

    try {
      const { data } = await axios.put("http://localhost:8000/api/chat/groupremove",
        {
          chatId: selectedChat._id,
          userId: userToBeRemoved._id,
        },
        {
          headers: {
            "Content-Type": "application/json",
          }
        }
      )
      console.log(data);
      setSelectedChat(data);
    } catch (error) {
        toast({
          title: "Error occured",
          decription: "Failed to remove user",
          duration: 3000,
          status: "warning",
          isClosable: true,
          position: "bottom",
        });
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
          decription: "Failed to load users data",
          duration: 3000,
          status: "warning",
          isClosable: true,
          position: "bottom",
        });
        setLoading(false);
    }
  }

  // hanldeExitGroup lets the user to exit from group
  const handleExitGroup = async(userToExit) => {
    try {
      const { data } = await axios.put("http://localhost:8000/api/chat/groupexit",
        {
          chatId: selectedChat._id,
          userId: userToExit._id,
        },
        {
          headers:{
            "Content-Type": "application/json",
          }
        }
      )

      console.log(data);
      toast({
        title: data.message,
        duration: 3000,
        status: "success",
        isClosable: true,
        position: "bottom",
      });
      setFetchAgain(!fetchAgain);
      setSelectedChat("")
    } catch (error) {
        toast({
          title: "Error occured",
          decription: "Failed to exit group",
          duration: 3000,
          status: "warning",
          isClosable: true,
          position: "bottom",
        });
    }
  }

  // Store the group users, so that when adding users we can perform checking
  useEffect(() => {
    setSelectedUsers(selectedChat.users.map((u) => (u._id)));
  },[selectedChat.users])

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
              <Button ml="5px" onClick={() => handleGroupRename()} isLoading={loading}>
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
            <Button colorScheme={"orange"} onClick={() => handleExitGroup(user)}>Exit Group</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

export default UpdateGroupChat