import React, { useState } from 'react';
import axios from 'axios';

const CreatePostForm = ({ onPostCreated }) => {
    const [message, setMessage] = useState('')
    const [errorText, setErrorText] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault()
        setErrorText('')

        try {
            const token = localStorage.getItem('token')
            if (!token) {
                throw new Error('You are not authorized to post. Please re-log in and try again.')
            }

            const response = await axios.post('http://localhost:1323/api/v1/posts', { message }, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            })

            setMessage('')
            onPostCreated(response.data)
        } catch (error) {
            if (error.response && error.response.staus === 401) {
                setErrorText('You are not authorized to create a post. Please try again.')
                console.error('Unauthorized: Missing or malformed JWT token')
            } else if (error.response) {
                console.error('An error occurred while creating the post:', error.response.data.message)
                setErrorText(error.response.data.message || 'Failed to create post. Please try again.')
            } else {
                console.error('Error:', error.message)
                setErrorText('Failed to create post. Please check your network connection and try again.')
            }
        }
    }

    return (
        <form onSubmit={handleSubmit}>
            <div className="mb-3">
                <label htmlFor="message" className="form-label">Message</label>
                <textarea
                    id="message"
                    className="form-control"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    required
                />
            </div>
            {errorText && <p className="text-danger">{errorText}</p>}
            <button type="submit" className="btn btn-primary">Create Post</button>
        </form>
    );
};

export default CreatePostForm;
