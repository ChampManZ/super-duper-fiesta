import React from "react";
import { useNavigate } from "react-router-dom";
import ActionButton from "./ActionButton";
import { logoutUser } from "../services/api";

const LogoutButton = ( style ) => {
    const navigate = useNavigate()
    const logoutButtonStyle = {
        backgroundColor: "#ff4d4d",
        ...style
    }

    const handleLogout = async () => {
        try {
            await logoutUser()
            console.log('User logged out')
            console.log('Token:', localStorage.getItem('token'))
            navigate('/')
        } catch (error) {
            console.error('Error logging out:', error.message)
        }
    }

    return (
        <ActionButton text="Logout" onClick={handleLogout} style={logoutButtonStyle} />
    )
}

export default LogoutButton;