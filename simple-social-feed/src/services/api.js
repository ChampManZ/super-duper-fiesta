import axios from 'axios'

const API_BASE_URL = 'http://localhost:1323/api'

export const createUser = (userData) => {
    return axios.post(`${API_BASE_URL}/v1/users`, userData)
}

export const loginUser = (loginData) => {
    return axios.post(`${API_BASE_URL}/v1/login`, loginData)
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