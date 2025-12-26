import axios from 'axios';

const BASE_URL = 'https://absolut-cinema-umwih.ondigitalocean.app/api';
const api = axios.create({
    baseURL: BASE_URL,
    withCredentials: true, 
    headers: {
        'Content-Type': 'application/json',
    },
});

export default api;