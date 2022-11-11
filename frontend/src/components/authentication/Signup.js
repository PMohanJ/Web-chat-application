import { FormControl, FormLabel, FormHelperText, VStack,Input, InputGroup, InputRightElement, Button } from '@chakra-ui/react'
import React from 'react'

const Signup = () => {
    const [name, setName] = React.useState();
    const [email, setEmail] = React.useState();
    const [password, setPassword] = React.useState();
    const [confirmpassword, setConfirmPassowrd] = React.useState();
    const [pic, setPic] = React.useState();
    const [show, setShow] = React.useState();
    const [picLoading, setPicLoading] = React.useState(false);

    function passwordVisibility(){
        setShow(!show);
    }

    function postDetails(){

    }

    function submitHandler(){

    }
  return (
    <VStack spacing="5px">
        <FormControl id="first-name" isRequired>
          <FormLabel>Name</FormLabel>
            <Input mb="10px"
              type="text" 
              placeholder="Enter your name"
              onChange={(e) => setName(e.target.value)}
            />
        </FormControl>

        <FormControl id="email" isRequired>
          <FormLabel>Email</FormLabel>
            <Input type="text" 
              placeholder="Enter your email"
              onChange={(e) => setEmail(e.target.value)}
            />
            <FormHelperText mb="10px">We'll never share your email.</FormHelperText>
        </FormControl>

        <FormControl id="password" isRequired>
          <FormLabel>Password</FormLabel>
            <InputGroup>
              <Input mb="10px" 
                type={show ? "text" : "password"}
                placeholder="Enter your password"
                onChange={(e) => setPassword(e.target.value)}
              />
              <InputRightElement width="4.5rem">
                <Button h="1.75rem" size="sm" onClick={passwordVisibility}>
                  {show ? "Hide" : "Show"}
                </Button>
              </InputRightElement>
            </InputGroup>
        </FormControl>
        
        <FormControl id="password" isRequired>
          <FormLabel>Confirm Password</FormLabel>
            <InputGroup>
              <Input mb="10px" 
                type={show ? "text" : "password"}
                placeholder="Confirm your password"
                onChange={(e) => setConfirmPassowrd(e.target.value)}
              />
              <InputRightElement width="4.5rem">
                <Button h="1.75rem" size="sm" onClick={passwordVisibility}>
                  {show ? "Hide" : "Show"}
                </Button>
              </InputRightElement>
            </InputGroup>
        </FormControl>

        <FormControl id="pic">
            <FormLabel>Upload your pic</FormLabel>
            <Input type="file"
              p={1.5}
              accept="image/*"
              onChange={(e) => postDetails(e.target.files[0])}
            />
        </FormControl>

        <Button colorScheme="blue"
          width="30%"
          mt="15px"
          onClick={submitHandler}
          isLoading={picLoading}
          borderRadius="0 15px 15px 15px"
        >
            Sign Up
        </Button>
    </VStack>
  )
}

export default Signup