import axios from 'axios'

const API_BASE_URL = 'http://localhost:1323/api'

export const createUser = (userData) => {
    return axios.post(`${API_BASE_URL}/v1/users`, userData);
}
