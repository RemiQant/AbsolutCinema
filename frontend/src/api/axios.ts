import axios from 'axios';

const BASE_URL = 'https://absolut-cinema-umwih.ondigitalocean.app/api';
const LOCAL_URL = 'http://localhost:8080/api';

// Use production URL in production, localhost in development
const api = axios.create({
    baseURL: import.meta.env.MODE === 'production' ? BASE_URL : LOCAL_URL,
    withCredentials: true, 
    headers: {
        'Content-Type': 'application/json',
    },
});


export default api;
