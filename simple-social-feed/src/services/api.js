import axios from 'axios'
import { ADMIN_USERNAME, ADMIN_PASSWORD } from './constant'

const API_BASE_URL = 'http://localhost:1323/api'
const ADMIN_HEADER = {         
    auth: {
            username: ADMIN_USERNAME,
            password: ADMIN_PASSWORD
    }
}

export const createUser = (userData) => {
    return axios.post(`${API_BASE_URL}/v1/users`, userData)
}

export const loginUser = (loginData) => {
    return axios.post(`${API_BASE_URL}/v1/login`, loginData)
}

export const getUser = (uid) => {
    return axios.get(`${API_BASE_URL}/v1/admin/users?uid=${uid}`, ADMIN_HEADER)
}

export const getUsername = (username) => {
    return axios.get(`${API_BASE_URL}/v1/admin/users?username=${username}`, ADMIN_HEADER)
}

export const createPost = (message, token) => {
    return axios.post('http://localhost:1323/api/v1/restricted/posts', { message }, {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    })
}

export const getPosts = () => {
    return axios.get(`${API_BASE_URL}/v1/posts`)
}

export const logoutUser = () => {
    localStorage.removeItem('token')
    return axios.post(`${API_BASE_URL}/v1/logout`)
}

export const accessProtectedRoute = () => {
    // eslint-disable-next-line no-sequences
    return `${API_BASE_URL}/v1/admin/main`, {
        method: 'GET',
        headers: {
            'Authorization': 'Basic ' + btoa(`${ADMIN_USERNAME}:${ADMIN_PASSWORD}`)
        }
    }
}

export const getMigration = () => {
    return axios.get(`${API_BASE_URL}/v1/admin/get-migrations`, ADMIN_HEADER)
}

export const runMigration = (migrationID) => {
    return axios.post(`${API_BASE_URL}/v1/admin/run-migrations`, { migration_id: migrationID }, ADMIN_HEADER)
}

export const updateUser = (uid, userData, token) => {
    return axios.put(`${API_BASE_URL}/v1/restricted/users/${uid}`, userData, {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
}

export const addComment = (postID, comment, token) => {
    return axios.post(`${API_BASE_URL}/v1/restricted/comments`, {
        post_id: postID,
        comment_msg: comment
    }, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })
}

export const getComment = (postID) => {
    return axios.get(`${API_BASE_URL}/v1/comments/${postID}`)
}