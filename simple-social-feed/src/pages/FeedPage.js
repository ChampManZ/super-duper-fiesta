import React, { useEffect, useState } from "react";
import ActionButton from "../components/ActionButton";
import ButtonLink from "../components/ButtonLink";
import CreatePostForm from "../components/CreatePostForm";
import MessageFeed from "../components/MessageFeed";
import LogoutButton from "../components/LogoutButton";
import Comment from "../components/Comment";
import CommentForm from "../components/CommentForm";
import { getPosts, getComment } from "../services/api";
import { jwtDecode } from "jwt-decode";

function FeedPage() {
    const [posts, setPosts] = useState([])
    const [loading, setLoading] = useState(true)
    const [errorText, setErrorText] = useState('')
    const [showCreatePostForm, setShowCreatePostForm] = useState(false)
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [username, setUsername] = useState('')

    const padding = {
        padding: '12px'
    }

    useEffect(() => {

        const token = localStorage.getItem('token')
        if (token) {
            const decodedToken = jwtDecode(token)
            setUsername(decodedToken.username)
            setIsLoggedIn(true)
        } else {
            setIsLoggedIn(false)
        }

        const fetchPosts = async () => {
            try {
                const response = await getPosts()
                if (response.data && Array.isArray(response.data)) {
                    setPosts(response.data)
                } else {
                    setPosts([])
                }
            } catch (error) {
                if (error.response && error.response.status === 401) {
                    setErrorText('Your session has expired. Please log in again.')
                } else {
                    setErrorText('Failed to load feed.')
                }
            } finally {
                setLoading(false)
            }
        }

        fetchPosts()
    }, [])

    const handlePostCreated = (newPost) => {
        setPosts(prevPosts => [newPost, ...prevPosts])
        setShowCreatePostForm(false)
        window.location.reload()
    }

    const handleCommentCreated = async (postId) => {
        try {
            const response = await getComment(postId)
            return response.data || []
        } catch (error) {
            console.error('Failed to load comments', error)
            return []
        }
    }

    if (loading) {
        return <p>Loading feed...</p>
    }

    if (errorText) {
        return <p>{errorText}</p>
    }

    const isPostAvailable = posts.length > 0

    return (
        <div className="feed-page">
            <h2>Social Feed</h2>
            {isPostAvailable ? (
                posts.map(post => (
                    <div key={post.post_id}>
                        <MessageFeed
                            username={post.username}
                            postTime={post.post_created_at ? new Date(post.post_created_at).toLocaleString() : 'Unknown'}
                            message={post.post_message}
                        />
                        <Comment postID={post.post_id} />
                        <CommentForm postID={post.post_id} onCommentAdded={(newComment) => handleCommentCreated(post.post_id, newComment)} />
                    </div>
                ))
            ) : (
                <p>No post available</p>
            )}
            <div className="create-post">
                <ActionButton text="Create Post" onClick={() => setShowCreatePostForm(true)} />
                {showCreatePostForm && <CreatePostForm onPostCreated={handlePostCreated} />}
            </div> <br />
            { isLoggedIn ? (
                <div>
                    <ButtonLink text="My Profile" href={`/${username}`} style={padding} />
                    <LogoutButton />
                </div>
                ) : <div>
                <h4>Join the fun now!</h4>
                <ButtonLink text="Login" href="/login" />
                <ButtonLink text="Register" href="/register" />
            </div> 
            } 
        </div>
    )
}

export default FeedPage;
