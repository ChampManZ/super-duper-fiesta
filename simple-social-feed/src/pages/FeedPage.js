import React, { useEffect, useState } from "react";
import axios from "axios";

function FeedPage() {

    const [posts, setPosts] = useState([])
    const [loading, setLoading] = useState(true)
    const [errorText, setErrorText] = useState('')

    // Reminder Note: 
    // The useEffect hook is used to fetch data from the server when the component is mounted.
    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await axios.get('http://localhost:1323/api/v1/posts', {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('token')}`
                    },
                })
                setPosts(response.data)
            } catch (error) {
                setErrorText('Failed to load posts')
            } finally {
                setLoading(false)
            }
        }

        fetchPosts()
    }, [])

    if (loading) {
        return <p>Loading feed...</p>
    }

    if (errorText) {
        return <p>{errorText}</p>
    }

    return (
        <div className="feed-page">
            <h2>Social Feed</h2>
        </div>
    )
}

export default FeedPage;