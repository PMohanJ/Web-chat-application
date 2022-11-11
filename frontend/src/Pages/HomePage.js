import React from 'react'
import {Container, Box, Text} from "@chakra-ui/react";
import {Tabs, TabList, Tab, TabPanels, TabPanel} from "@chakra-ui/react";

const HomePage = () => {
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
        borderRadius="0 0 10px 10px"
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
              
            </TabPanel>
            <TabPanel>
              
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
    </Container>
  )
}

export default HomePage;