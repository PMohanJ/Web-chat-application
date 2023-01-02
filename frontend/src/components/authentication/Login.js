import React, { useState } from 'react'
import { FormControl, FormLabel, VStack,Input, InputGroup, InputRightElement, Button } from '@chakra-ui/react'
import { useToast } from '@chakra-ui/react'
import axios from "axios"
import { useNavigate } from "react-router-dom"

const Login = () => {
  const [show, setShow] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [loading, setLoading] = useState(false);
  const handleClick = () => setShow(!show);
  const toast = useToast();
  const navigate = useNavigate();

  const submitHandler = async() => {
    setLoading(true);

    if (!email || !password){
      toast({
        title: "Please fill all details",
        status: "warning",
        isClosable: true,
        position: "bottom",
        duration: 4000,
      });
      setLoading(false);
      return;
    } 

    try {
      // make a post request to backend with user credentials
      const { data } = await axios.post(
        "/api/user/login", 
        {  
          email:email, 
          password: password
        },
        {
          headers: {
            "Content-Type":"application/json",
          }
        }
      );
      
      toast({
        title: "Login Successful",
        status: "success",
        duration: 4000,
        isClosable: true,
        position: "top",
      });

      if (localStorage.getItem("userInfo") !== null) {
        localStorage.removeItem("userInfo");
      }
      // storing user details in localstorage
      localStorage.setItem("userInfo", JSON.stringify(data));
      setLoading(false);
      navigate("/chat");

    } catch (error) {
      toast({
        title: "Error Occured!",
        description: error.response.data.error,
        status: "error",
        duration: 5000,
        isClosable: true,
        position: "bottom",
      });
      setLoading(false);
    }
  }

  return (
    <VStack spacing="10px">
      <FormControl id="email2" isRequired>
        <FormLabel>Email Address</FormLabel>
        <Input
          value={email}
          type="email"
          placeholder="Enter Your Email Address"
          onChange={(e) => setEmail(e.target.value)}
        />
      </FormControl>

      <FormControl id="password2" isRequired>
        <FormLabel>Password</FormLabel>
        <InputGroup size="md">
          <Input
            value={password}
            type={show ? "text" : "password"}
            placeholder="Enter password"
            onChange={(e) => setPassword(e.target.value)}
          />
          <InputRightElement width="4.5rem">
            <Button h="1.75rem" size="sm" onClick={handleClick}>
              {show ? "Hide" : "Show"}
            </Button>
          </InputRightElement>
        </InputGroup>
      </FormControl>

      <Button
        colorScheme="blue"
        width="50%"
        mt={15}
        onClick={submitHandler}
        isLoading={loading}
      >
        Login
      </Button>

      <Button
        variant="solid"
        colorScheme="gray"
        width="50%"
        onClick={() => {
            setEmail("guest@example.com");
            setPassword("123456");
        }}
        >
        Get Guest User Credentials
      </Button>
    </VStack>
    );
}

export default Login