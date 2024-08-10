import React, { useEffect, useState } from "react";
import ActionButton from "../components/ActionButton";
import ButtonLink from "../components/ButtonLink";
import CreatePostForm from "../components/CreatePostForm";
import axios from "axios";

function FeedPage() {

    const [posts, setPosts] = useState([])
    const [loading, setLoading] = useState(true)
    const [errorText, setErrorText] = useState('')
    const [showCreatePostForm, setShowCreatePostForm] = useState(false)

    // Reminder Note: 
    // The useEffect hook is used to fetch data from the server when the component is mounted.
    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const token = localStorage.getItem('token')
                if (!token) {
                    throw new Error('You are not authorized to view this page. Please log in and try again.')
                }

                const response = await axios.get('http://localhost:1323/api/v1/posts')
                setPosts(response.data)
            } catch (error) {
                if (error.response && error.response.status === 401) {
                    setErrorText('Your session has expired. Please log in again.')
                } else {
                    setErrorText('Failed to load feed')
                }
            } finally {
                setLoading(false)
            }
        }

        fetchPosts()
    }, [])

    const postCreator = (newPost) => {
        setPosts([newPost, ...posts])
        setShowCreatePostForm(false)
        console.log('Token:', localStorage.getItem('token'))
    }

    if (loading) {
        return <p>Loading feed...</p>
    }

    if (errorText) {
        return <p>{errorText}</p>
    }

    return (
        <div className="feed-page">
            <h2>Social Feed</h2>
            <ActionButton text={'Create Post'} onClick={() => setShowCreatePostForm(true)} /> <br /> <br />
            {showCreatePostForm && <CreatePostForm onPostCreated={postCreator} />} <br /> <br />
            <ButtonLink href="/" text="Home" />
        </div>
    )
}

export default FeedPage;