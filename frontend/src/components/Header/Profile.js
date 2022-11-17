import React from 'react'
import { Text, useDisclosure, Drawer, DrawerBody, DrawerOverlay, DrawerContent, DrawerCloseButton } from '@chakra-ui/react'

const Profile = ({user, children}) => {

  const { isOpen, onOpen, onClose } = useDisclosure()

  return (
    <>  <span onClick={onOpen}>{children}</span>
        <Drawer
          isOpen={isOpen}
          placement="right"
          onClose={onClose}
        >
          <DrawerOverlay />
           <DrawerContent>
              <DrawerCloseButton />

              <DrawerBody 
                display={"flex"}
                flexDir="column"
                alignItems={"center"}
                
              >
                <Text 
                  fontSize={"2xl"}
                  fontFamily={"serif"}
                  padding={"5px"}
                >
                    {user.name}
                </Text>

                <img 
                  style={{width:"120px", height:"120px", borderRadius: "50%", border: "3px solid #7479e3"}}
                  alt="Profile Pic"
                  src={user.pic}
                />

                <Text 
                  mt={"20px"}
                  fontFamily="serif"
                  fontSize={"2xl"}
                >
                    {user.email}
                </Text>
              </DrawerBody>
          </DrawerContent>
        </Drawer>
    </>
  )
}

export default Profile 