import React, { useState, useEffect } from "react";
import { getUser, updateUser } from "../services/api";
import { useNavigate } from "react-router-dom";

const EditUserForm = ({uid}) => {

    const [userData, setUserData] = useState({
        username: "",
        firstname: "",
        lastname: ""
    })
    const [loading, setLoading] = useState(true)
    const [text, setText] = useState("")
    const navigate = useNavigate()

    useEffect(() => {
        const fetchUserData = async () => {
            try {
                const res = await getUser(uid)
                setUserData({
                    username: res.data.username,
                    firstname: res.data.firstname,
                    surname: res.data.surname
                })
            } catch (error) {
                setText("Failed to fetch user data")
            } finally {
                setLoading(false)
            }
        }

        fetchUserData()
    }, [uid])

    const handleChange = (e) => {
        const { name, value } = e.target
        setUserData({
            ...userData,
            [name]: value
        })
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        try {
            const token = localStorage.getItem('token')
            const response = await updateUser(uid, userData, token)
            setText('User updated successfully. Please re-login. Logging out in 3 seconds...')

            setTimeout(() => {
                // Force user to log in again after updating their profile to reset the token
                localStorage.removeItem('token')
                navigate('/')
            }, 3000)
            console.log('User updated:', response.data)
        } catch (error) {
            setText('Failed to update user')
        }
    }

    if (loading) {
        return <p>Loading...</p>
    }

    console.log('User data:', userData)

    return (
        <div>
            <h3>Edit User</h3>
            { text && <p>{text}</p> }
            <form onSubmit={handleSubmit}>
                <div className="mb-3">
                    <label htmlFor="username" className="form-label">Username</label>
                    <input 
                        type="text" 
                        name="username" 
                        id="username" 
                        className="form-control" 
                        value={userData.username} 
                        onChange={handleChange} 
                        required 
                    />
                </div>
                <div className="mb-3">
                    <label htmlFor="firstname" className="form-label">Firstname</label>
                    <input 
                        type="text" 
                        name="firstname" 
                        id="firstname" 
                        className="form-control" 
                        value={userData.firstname} 
                        onChange={handleChange} 
                        required 
                    />
                </div>
                <div className="mb-3">
                    <label htmlFor="surname" className="form-label">Surname</label>
                    <input 
                        type="text" 
                        name="surname" 
                        id="surname" 
                        className="form-control" 
                        value={userData.surname} 
                        onChange={handleChange} 
                        required 
                    />
                </div>
                <button type="submit" className="btn btn-primary">Update Profile</button>
            </form>
        </div>
    )
}

export default EditUserForm;