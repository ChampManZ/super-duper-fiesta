import React, { useState, useEffect } from "react";
import { jwtDecode } from "jwt-decode";
import EditUserForm from "../components/EditUserForm";
import ButtonLink from "../components/ButtonLink";
import LogoutButton from "../components/LogoutButton";
import { useParams } from "react-router-dom";
import { getUsername } from "../services/api";

function UserProfilePage() {

    const { username } = useParams()
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [uid, setUid] = useState(null)
    const [currentUsername, setCurrentUsername] = useState('')
    const [userExists, setUserExists] = useState(null)

    useEffect(() => {
        const checkUser = async () => {
            const token = localStorage.getItem('token')
            let decodedToken = null

            if (token) {
                try {
                    decodedToken = jwtDecode(token)
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

            try {
                const response = await getUsername(username)
                console.log(response.data)
                if (response.data && response.data.username) {
                    setUserExists(true)
                } else {
                    setUserExists(false)
                }
            } catch (error) {
                console.error("Error checking user:", error)
                setUserExists(false)
            }
        }

        checkUser()
    }, [username])

    const isCurrentUserProfile = username === currentUsername
    // console.log(username, currentUsername)
    // console.log(typeof username, typeof currentUsername)
    // console.log(isCurrentUserProfile)

    if (userExists === null) {
        return <h1>Loading...</h1>
    }

    if (!userExists) {
        return (
            <div>
                <h4>{username} does not exist</h4>
                <ButtonLink text="Home" href="/" />
                <ButtonLink text="Register" href="/register" />
            </div>
        )
    }

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