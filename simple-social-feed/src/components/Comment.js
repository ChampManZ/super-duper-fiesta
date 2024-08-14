import React, { useEffect, useState } from "react";
import { getComment } from "../services/api";

function Comment({ postID }) {

    const [comments, setComments] = useState([])
    const [loading, setLoading] = useState(true)
    const [errorText, setErrorText] = useState('')

    useEffect(() => {
        const fetchComments = async () => {
            try {
                const response = await getComment(postID)
                if (response.status === 200) {
                    setComments(response.data)
                } else {
                    setErrorText('Failed to load comments.')
                }
            } catch (error) {
                setErrorText('An error occured while loading comments.')
            } finally {
                setLoading(false)
            }
        }

        fetchComments()
    }, [postID])

    if (loading) {
        return <p>Loading comments...</p>
    }

    if (errorText) {
        return <p>{errorText}</p>
    }

    return (
        <div>
            {comments && comments.length > 0 ? (
                comments.map((comment, index) => (
                    <div key={index}>
                        <p><strong>{comment.username}</strong>: {comment.comment_msg}</p>
                    </div>
                ))
            ) : (
                <p>No comments available.</p>
            )}
        </div>
    );
}

export default Comment;