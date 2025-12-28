import axios from 'axios';

// In production, use relative path since frontend and backend are on same domain
// In development, use full localhost URL
const api = axios.create({
    baseURL: import.meta.env.MODE === 'production' ? '/api' : 'http://localhost:8080/api',
    withCredentials: true, 
    headers: {
        'Content-Type': 'application/json',
    },
});


export default api;
