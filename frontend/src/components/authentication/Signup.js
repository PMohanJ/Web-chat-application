import { FormControl, FormLabel, FormHelperText, VStack,Input, InputGroup, InputRightElement, Button} from '@chakra-ui/react'
import React from 'react'
import axios from "axios"
import { useToast } from '@chakra-ui/react'
import { useNavigate } from "react-router-dom"

const Signup = () => {
  const [name, setName] = React.useState();
  const [email, setEmail] = React.useState();
  const [password, setPassword] = React.useState();
  const [confirmpassword, setConfirmPassowrd] = React.useState();
  const [pic, setPic] = React.useState();
  const [show, setShow] = React.useState();
  const [loading, setLoading] = React.useState(false);
  const toast = useToast();
  const navigate = useNavigate();

  function passwordVisibility(){
    setShow(!show);
  }

  // Performs registration of user with provided details
  const submitHandler = async() => {
    setLoading(true);

    if (!name || !email || !password || !confirmpassword){
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

    if (password !== confirmpassword){
      toast({
        title: "Confirm password does not match",
        status: "warning",
        duration: 4000,
        position: "bottom",
        isClosable: true,
      });
      setLoading(false);
      return;
    }
    console.log(name, email, password, pic)

    try {
      // make a post request to backend with user credentials
      const { data } = await axios.post(
        // As the backend is developed using gin package, the trailing spaces
        // maters, so be cautious with it..
        "http://localhost:8000/api/user/", 
        { 
          name: name, 
          email:email, 
          password: password, 
          pic: pic
        },
        {
          headers: {
            "Content-Type":"application/json",
          }
        }
      );
      
      console.log(data);
      toast({
        title: "Registration Successful",
        status: "success",
        duration: 5000,
        isClosable: true,
        position: "top",
      });

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
  };

  // Check for uploaded image, store it in cloudinary sotrage and get url 
  // of the image/profile 
  function postDetails(pics){
    setLoading(true)

    if(pics === undefined){
      toast({
        title:"Please select a image",
        status: "warning",
        duration: 4000,
        isClosable: true,
        position:"top"

      });
      return; 
    }

    // check for valid formats of the image
    if (pics.type === "image/jpeg" || pics.type === "image/jpg" || pics.type === "image/png"){
      const data = new FormData();
      data.append("file", pics)
      data.append("upload_preset", "webchatapp")

      fetch("https://api.cloudinary.com/v1_1/dkqc4za4f/image/upload",{
        method:"post",
        body: data, 
      }).then((res) => res.json())
        .then((data) => {
          setPic(data.url.toString());
          console.log(data.url.toString());
          setLoading(false);
        })
        .catch((err) => {
          console.log(err);
          setLoading(false)
        })
    } else {
      toast({
        title: "Invalid format",
        description: "Please select jpg/jpeg/png formats",
        duration: 5000,
        isClosable: true,
        position: "bottom"
      });
      setLoading(false)
      return;
    }
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

        <FormControl id="email1" isRequired>
          <FormLabel>Email</FormLabel>
            <Input type="text" 
              placeholder="Enter your email"
              onChange={(e) => setEmail(e.target.value)}
            />
            <FormHelperText mb="10px">We'll never share your email.</FormHelperText>
        </FormControl>

        <FormControl id="password1" isRequired>
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
              // if user selects mul files, get first one
              onChange={(e) => postDetails(e.target.files[0])} 
            />
        </FormControl>

        <Button colorScheme="blue"
          width="30%"
          mt="15px"
          onClick={submitHandler}
          isLoading={loading}
          borderRadius="0 15px 15px 15px"
        >
            Sign Up
        </Button>
    </VStack>
  )
}

export default Signup