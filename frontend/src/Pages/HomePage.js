import React, { useEffect } from 'react'
import {Container, Box, Text} from "@chakra-ui/react";
import {Tabs, TabList, Tab, TabPanels, TabPanel} from "@chakra-ui/react";
import Login from "../components/authentication/Login"
import Signup from "../components/authentication/Signup"
import { ChatState } from '../context/ChatProvider';
import { useNavigate } from 'react-router-dom';

const HomePage = () => {
  const {user, } = ChatState();
  const navigate = useNavigate();

  useEffect(() => {
    if(user){
      navigate("/chat")
    }
  }, [])
  

  return (
    <Container maxW="xl"  centerContent>
      <Box 
        display="flex"
        justifyContent="center"
        fontFamily="Work Sans"
        width="200%"
        fontSize="4xl"
        background="#cacfcc"
        margin="0 40px 20px 40px"
        borderRadius="40px 40px 10px 10px"
        boxShadow="0 5px 10px gray "
      >
        <Text color="white" fontWeight="bold">Chat App</Text>
      </Box>

      <Box bg="white" 
        width="100%" p={4} 
        borderWidth="1px"
        borderRadius="0 0 10px 10px"
        boxShadow="0 5px 10px gray "
      >
        <Tabs variant='solid-rounded'>
          <TabList >
            <Tab _selected={{ color: 'white', bg: 'blue.300' }} width="50%">Login</Tab>
            <Tab _selected={{ color: 'white', bg: 'blue.300' }} width="50%">Sign Up</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <Login/>
            </TabPanel>
            <TabPanel>
              <Signup/>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
    </Container>
  )
}

export default HomePage;