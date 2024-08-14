import React, { useState } from "react";
import { addComment } from "../services/api";

function CommentForm({ postID, onCommentAdded }) {
    const [comment, setComment] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault()
        try {
            const token = localStorage.getItem('token')
            if (!token) {
                throw new Error('You are not authorized to post. Please re-log in and try again.')
            }

            const response = await addComment(postID, comment, token)
            if (response.status === 201) {
                onCommentAdded(response.data)
                setComment('')
                window.location.reload()
            } else {
                alert('Failed to add comment. Please try again.')
            }
        } catch (error) {
            if (error.response && error.response.status === 401) {
                alert('You are not authorized to post. Please re-log in and try again.')
            } else {
                alert('An error occured while adding the comment. Please try again.')
            }
        }
    }

    return (
        <form onSubmit={handleSubmit}>
            <input type="text" value={comment} onChange={(e) => setComment(e.target.value)} />
            <button type="submit">Comment</button>
        </form>
    )
}

export default CommentForm;