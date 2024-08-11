import React, { useState } from 'react';
import { createPost } from '../services/api';

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
    
            const response = await createPost(message, token)
            setMessage('')
            onPostCreated(response.data)

        } catch (error) {
            if (error.response && error.response.staus === 401) {
                setErrorText('You are not authorized to create a post. Please try again.')
                console.error('Unauthorized: Missing or malformed JWT token')
            } else if (error.response && error.response.staus === 500) {
                console.error('An error occurred while creating the post:', error.response.data.message)
                setErrorText('Internal server error. Please try again later.')
            } else {
                console.error('Error:', error.message)
                setErrorText('Failed to create post. Please try again.')
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
            <button type="submit" className="btn btn-primary">Add Post</button>
        </form>
    );
};

export default CreatePostForm;
