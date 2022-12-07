 import React from 'react'
 import { useState } from 'react'
 import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  useDisclosure,
  Button,
  useToast,
  FormControl,
  Input,
  Box,
  Spinner,
  Stack
  } from '@chakra-ui/react'
import { ChatState } from '../../context/ChatProvider';
import axios from 'axios'
import UserSearchProfile from './UserSearchProfile';
import UserBadge from './UserBadge';

const GroupChatModel = ({children}) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [groupName, setGroupName] = useState("");
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [search, setSearch] = useState("");
    const [searchResults, setSearchResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const toast = useToast();

    const {user, chats, setChats} = ChatState();

    // handleSearch retrieves the users that match the search query
    const handleSearch = async(name) => {
      let query = name.trim() 
      if (!query) {
          return;
      }
      setSearch(query)
      try {
        setLoading(true)
        const url = "http://localhost:8000/api/user/search?search="+query.toString()
        const { data } = await axios.get(url,
          {
            headers:{
              "Content-Type": "application/json",
              "Authorization": `Bearer ${user.token}`
            }
          }
        )
        
        console.log(data)
        setSearchResults(data.filter((u) => u._id !== user._id))
        setLoading(false)

      } catch (error) {
        toast({
          title: "Error Occured!",
          description: "Failed to Load the Search Results",
          status: "error",
          duration: 5000,
          isClosable: true,
          position: "bottom-left",
        });
        setLoading(false);
        }
    };

    // handleGroup adds the users to selectedUsers list while creating the groupchat
    const  handleGroup = (userToAdd) => {
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

        setSelectedUsers([...selectedUsers, userToAdd])
    };

    const handleRemoveUser = (user) => {
        setSelectedUsers(selectedUsers.filter((sel) => sel._id !== user._id))
    }

    // handleSubmit creates the groupchat with given details
    const handleSubmit = async() => {
      if (!groupName || !selectedUsers) {
        toast({
          title: "Please fill all details",
          duration: 4000,
          status: "warning",
          isClosable: true,
          position: "top",
        });
        return;
      }

      if (selectedUsers.length < 2){
        toast({
          title: "Please select atleast 2 members",
          duration: 4000,
          status: "warning",
          isClosable: true,
          position: "top",
        });
        return;
      }

      try {
        const { data } = await axios.post("http://localhost:8000/api/chat/group",
          {
            groupName: groupName,
            users: selectedUsers.map((u) => u._id)
          },
          {
            headers:{
              "Content-Type": "application/json",
              "Authorization": `Bearer ${user.token}`
            }
          }
        )

        console.log(data);
        setChats([data, ...chats]);
        onClose();
      } catch (error) {
          toast({
            title: "Failed to Create the Chat!",
            description: error.response.data,
            status: "error",
            duration: 5000,
            isClosable: true,
            position: "bottom",
          });
      }
    };

   return (
    <>
    <span onClick={onOpen}>{children}</span>

    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader
          fontSize="25px"
          fontFamily="Work sans"
          display="flex"
          justifyContent="center"
        >
            Create Group Chat
        </ModalHeader>

        <ModalCloseButton />
        <ModalBody
          display="flex" 
          flexDir="column" 
          alignItems="center"
        >
          <FormControl>
            <Input 
              placeholder="Chat Name"
              mb={2}
              onChange={(e) => setGroupName(e.target.value)}
            />
          </FormControl>

          <FormControl>
            <Input
              placeholder="Add Users eg: Rohan" 
              mb={1}
              onChange={(e) => handleSearch(e.target.value)}
            />
          </FormControl>
          
          <Box
            w="100%"
            display="flex"
            flexWrap="wrap"
          >
            {selectedUsers.map((u) => (
                <UserBadge 
                  key={u._id}
                  handleFunction={() => handleRemoveUser(u)}
                  user={u}
                />
            ))}
          </Box>

          {loading ? 
            <Spinner/>
            : <Stack w="100%">
                { searchResults?.slice(0, 4).map((u) => (
                    <UserSearchProfile 
                        key={u._id}
                        handleFunction={() => handleGroup(u)}
                        user={u}
                    /> 
                  )) }
              </Stack>
          }
        </ModalBody>

        <ModalFooter>
          <Button onClick={handleSubmit} colorScheme="blue" isLoading={loading}>
              Create Chat
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  </>
   )
 }
 
 export default GroupChatModel