import React, { useState, useEffect } from "react";
import { jwtDecode } from "jwt-decode";
import EditUserForm from "../components/EditUserForm";
import ButtonLink from "../components/ButtonLink";
import LogoutButton from "../components/LogoutButton";
import { useParams } from "react-router-dom";

function UserProfilePage() {

    let { username } = useParams()
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [uid, setUid] = useState(null)
    const [currentUsername, setCurrentUsername] = useState('')

    useEffect(() => {
        const token = localStorage.getItem('token')
        if (token) {
            try {
                const decodedToken = jwtDecode(token)
                console.log("Decoded token:", decodedToken)
                setUid(decodedToken.uid)
                setIsLoggedIn(true)
                setCurrentUsername(decodedToken.username)
            } catch (error) {
                console.error("Error decoding token:", error)
                setIsLoggedIn(false)
            }
        } else {
            setIsLoggedIn(false)
        }
    }, [])

    const isCurrentUserProfile = username === currentUsername

    return (
        <div>
            { isLoggedIn && isCurrentUserProfile ? (
                <div>
                    <h1>{username}</h1>
                    { uid && <EditUserForm uid={uid} />}
                    <LogoutButton />
                </div>
                ) : (
                <div>
                    <h4>{username}, you're not logging in. Join the fun now!</h4>
                    <ButtonLink text="Login" href="/login" />
                    <ButtonLink text="Register" href="/register" />
                </div> 
        )}
        </div>
    )
}

export default UserProfilePage;