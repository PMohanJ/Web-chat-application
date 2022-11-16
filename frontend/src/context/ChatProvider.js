import { createContext, useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const ChatContext = createContext();

const ChatProvider = ({childern}) => {
    const [user, setUser] = useState();
    const navigate = useNavigate();

    // get the user state of login
    useEffect(() => {
        // we still yet to provide jwt functionailty...
        const userInfo = JSON.parse(localStorage.getItem("userInfo"));
        setUser(userInfo);

        if (!userInfo){
            navigate("/");
        }
    }, [history]);
    return(
        <ChatContext.Provider value={{user,setUser}}>
          {childern}
        </ChatContext.Provider>
    ) 
};

export const ChatState = () => {
    return useContext(ChatContext);
}

export default ChatProvider;
