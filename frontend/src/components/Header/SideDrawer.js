import React from 'react'
import { Search2Icon } from '@chakra-ui/icons'
import { Tooltip, Box, Button, Text, DrawerHeader, useToast } from '@chakra-ui/react';
import { Menu, MenuButton, MenuList, MenuItem, Avatar} from '@chakra-ui/react'
import { ChevronDownIcon, BellIcon } from '@chakra-ui/icons'
import { useState } from 'react';
import { ChatState } from '../../context/ChatProvider';
import { useNavigate } from 'react-router-dom';
import { Input, useDisclosure, Drawer, DrawerBody, DrawerOverlay, DrawerContent, Spinner } from '@chakra-ui/react'
import Profile from './Profile';
import axios from "axios"
import UserSearchProfile from '../utils/UserSearchProfile'

const SideDrawer = () => {
  const [search, setSearch] = useState("")
  const [searchResults, setSearchResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure()
  const navigate = useNavigate();
  const toast = useToast();

  const {user, setSelectedChat, chats, setChats} = ChatState();

  const handleLogout = () => {
    localStorage.removeItem("userInfo");
    navigate('/')
  }

  // get the selected user chats 
  const getSelectedUserChat = async(userId)=> {
    try {
      setLoading(true);

      const {data} = await axios.post("http://localhost:8000/api/chat/",
        {
          userToBeAdded: userId
        },
        {
          headers:{
            "Content-Type": "application/json",
            "Authorization": `Bearer ${user.token}`
          }
        }
      )
      console.log(data);
      if (!chats) {
          setChats([data]);
      }
      else if (!chats.find((c) => c._id === data._id)) {
        setChats([data, ...chats]);
      }

      setLoading(false);
      setSelectedChat(data)
      
      onClose();
    } catch (error){
        console.log(error)
        toast({
          title: "Error while adding users",
          duration: 5000,
          isClosable: true,
          position: "bottom-left",
          status: "error"
        });
        setLoading(false);
    }
  }

  // Search for the user by given name or email
  const handleSearch = async() => {
    if (!search){
      toast({
        title: "Please enter a name or email",
        duration: 4000,
        status: "warning",
        position: "top-left",
        isClosable: true,
      });
      return;
    }

    try {
      setLoading(true);

      const url = `http://localhost:8000/api/user/search?search=${search}`
      const { data } = await axios.get(url,
        {
          headers: {
            "Content-Type":"application/json",
            "Authorization": `Bearer ${user.token}`
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
    <div>
      <Box 
        display="flex"
        justifyContent={"space-between"}
        backgroundColor="white"
        p={"5px"}
      >
        <Tooltip label="Search Users" hasArrow placement='right'>
          <Button onClick={onOpen} variant="ghost">
            <Search2Icon/>
            <Text ml="5px" display={{base:"none", md: "flex"}} p="5px">Search</Text>
            </Button>
        </Tooltip>
        
        <Text fontSize={"2xl"} fontFamily={"sans-serif"}>
          Chat App
        </Text>

        <div>
          <Menu>
            <MenuButton as={Button}>
              <BellIcon/>
            </MenuButton>
            <MenuList>
              {/*here comes notifications list */}
            </MenuList>
          </Menu>

          <Menu>
            <MenuButton as={Button} rightIcon={<ChevronDownIcon />}>
            <Avatar
              size="sm" 
              cursor="pointer" 
              name={user.name}
              src={user.pic}
            />
            </MenuButton>
            <MenuList>
              <Profile user={user}>
                <MenuItem >My Profile</MenuItem>
              </Profile>
              <MenuItem onClick={handleLogout}>Logout</MenuItem>
            </MenuList>
          </Menu>

        </div>
       
      </Box>

      <Drawer
        isOpen={isOpen}
        placement='left'
        onClose={onClose}
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerHeader borderBottomWidth={"1px"}> Search Users</DrawerHeader>
          <DrawerBody>
            <Box display={"flex"}>
              <Input 
                placeholder='Search by name or email...' 
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
              <Button 
                backgroundColor={"#e1e1e8"} 
                ml="5px" 
                onClick={handleSearch}
              >
                Go
              </Button>
            </Box>

            <Box textAlign={"center"} mt="10px">
              {loading && <Spinner color='red.500'/>} 
            </Box>
            
            {searchResults? searchResults.map((item) =>(
                <UserSearchProfile key={item._id} user={item} handleFunction={() => getSelectedUserChat(item._id)}/>
            )): null}
           
          </DrawerBody>
        </DrawerContent>
      </Drawer>
    </div>
  )
}

export default SideDrawer