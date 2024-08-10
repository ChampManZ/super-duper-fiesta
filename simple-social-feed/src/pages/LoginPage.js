import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function LoginPage() {

    const [loginData, setLoginData] = useState({
        identifier: '',
        password: ''
    })
    const [errorText, setErrorText] = useState('')
    const navigate = useNavigate()

    const handleInputEventChange = (e) => {
        const { name, value } = e.target
        setLoginData({
            ...loginData,
            [name]: value
        })
    }

    const handleLogin = async (e) => {
        e.preventDefault()
        setErrorText('')

        try {
            const response = await axios.post('http://localhost:1323/api/v1/login', loginData)
            console.log('User logged in:', response.data)
            console.log('Token:', response.data.token)
            localStorage.setItem('token', response.data.token)
            navigate('/feed')
        } catch (error) {
            if (error.response && error.response.data) {
                setErrorText(error.response.data.message)
            } else {
                console.error('Error logging user:', error.message)
            }
        }
    }
    
    return (
        <div className="login-page">
            <h2>Login</h2>
            <form onSubmit={handleLogin}>
                <div className="mb-3">
                    <label htmlFor="identifier" className="form-label">Username or Email</label>
                    <input 
                        type="text" 
                        className="form-control" 
                        id="identifier" 
                        name="identifier" 
                        value={loginData.identifier} 
                        onChange={handleInputEventChange} 
                        required
                    />
                </div>
                <div className="mb-3">
                    <label htmlFor="password" className="form-label">Password</label>
                    <input 
                        type="password" 
                        className="form-control" 
                        id="password" 
                        name="password" 
                        value={loginData.password} 
                        onChange={handleInputEventChange} 
                        required
                    />
                </div>
                <button type="submit" className="btn btn-primary">Login</button>
            </form>
            { errorText && <p className='text-danger'>{errorText}</p> }
        </div>
    )
}

export default LoginPage
