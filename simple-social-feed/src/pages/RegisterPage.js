import React, { useState } from 'react'
import { createUser } from '../services/api'
import { useNavigate } from "react-router-dom";

function RegisterPage() {
    const [registerData, setRegisterData] = useState({
        username: '',
        firstname: '',
        surname: '',
        email: '',
        password: ''
    })

    const [errorText, setErrorText] = useState("")
    const [validation, setValidation] = useState({
        usernameValid: false,
        firstnameValid: false,
        surnameValid: false,
        emailValid: false,
        passwordValid: false
    })
    const navigate = useNavigate()

    const handleInputEventChange = (e) => {
        const { name, value } = e.target
        setRegisterData({
            ...registerData,
            [name]: value
        })

        validateField(name, value)
    }

    const validateField = (name, value) => {
        let isValid;

        switch (name) {
            case 'username':
                isValid = value.length >= 3 && value.length <= 32
                setValidation((prev) => ({ ...prev, usernameValid: isValid }))
                break
            case 'firstname':
                isValid = value.trim() !== ''
                setValidation((prev) => ({ ...prev, firstnameValid: isValid }))
                break
            case 'surname':
                isValid = value.trim() !== ''
                setValidation((prev) => ({ ...prev, surnameValid: isValid }))
                break
            case 'email':
                const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;  // Regex pattern for email validation
                isValid = emailPattern.test(value)
                setValidation((prev) => ({ ...prev, emailValid: isValid }))
                break
            case 'password':
                isValid = value.length >= 8
                setValidation((prev) => ({ ...prev, passwordValid: isValid }))
                break
            default:
                break
        }
    }

    const bannerStyle = {
        border: "2px solid",
        padding: "10px",
        marginBottom: "20px",
        borderRadius: "5px",
        fontWeight: "bold"
      }
      
    const bannerSuccessStyle = {
        backgroundColor: "#d4edda",
        borderColor: "#c3e6cb",
        color: "#155724"
      }
      
    const bannerDangerStyle = {
        backgroundColor: "#f8d7da",
        borderColor: "#f5c6cb",
        color: "#721c24"
      }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setErrorText("")

        try {
            const response = await createUser(registerData)
            // In bigger test environment and production, I will avoid unnecessary logging all error to 
            // keep the console log important information only. This log is for debugging purpose.
            console.log('User registered successfully:', response.data)
            navigate('/success')
        } catch (error) {
            // Error log from the server side
            if (error.response && error.response.data) {
                setErrorText(error.response.data.message)
            } else {
                console.error('Error registering user:', error.message)
            }
        }
    }

    const requirementBanner = () => {
        const validator = Object.values(validation).every(Boolean)
        const bannerStyles = validator ? {...bannerStyle, ...bannerSuccessStyle} : {...bannerStyle, ...bannerDangerStyle}
        const symbol = (valid) => valid ? '✅' : '❌'

        return (
            <div style={bannerStyles} role="alert">
                <p>A username must be minimum 3 to 32 characters maxmimum {symbol(validation.usernameValid)} </p>
                <p>Firstname is required {symbol(validation.firstnameValid)} </p>
                <p>Surname is required {symbol(validation.surnameValid)} </p>
                <p>E-Mail is required {symbol(validation.emailValid)} </p>
                <p>Password must be minimum 8 characters {symbol(validation.passwordValid)} </p>
            </div>
        )
    }

    return (
        <div>
            <h2>Register</h2>
            {requirementBanner()}
            <form onSubmit={handleSubmit}>
            <div className="mb-3">
                <label htmlFor="username" className="form-label">Username</label>
                <input type="text" name="username" id="username" className="form-control" value={registerData.username} onChange={handleInputEventChange} required /> 
                <br />
            </div>
            <div className="mb-3">
                <label htmlFor="firstname" className="form-label">Firstname</label>
                <input type="text" name="firstname" id="firstname" className="form-control" value={registerData.firstname} onChange={handleInputEventChange} required /> 
                <br />
            </div>
            <div className="mb-3">
                <label htmlFor="surname" className="form-label">Surname</label>
                <input type="text" name="surname" id="surname" className="form-control" value={registerData.surname} onChange={handleInputEventChange} required />
                <br />
            </div>
            <div className="mb-3">
                <label htmlFor="email" className="form-label">E-Mail</label>
                <input type="email" name="email" id="email" className="form-control" value={registerData.email} onChange={handleInputEventChange} required /> 
                <br />
            </div>
            <div className="mb-3">
                <label htmlFor="password" className="form-label">Password</label>
                <input type="password" name="password" id="password" className="form-control" value={registerData.password} onChange={handleInputEventChange} required />
            </div>
            <button type="submit" className="btn btn-success">Register</button>
            { errorText && <p className='text-danger'>{errorText}</p> }
            </form>
        </div>
    )
}

export default RegisterPage