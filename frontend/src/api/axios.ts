import axios from 'axios';

const api = axios.create({
    baseURL: 'https://absolut-cinema-umwih.ondigitalocean.app/api', 
    withCredentials: true, 
    headers: {
        'Content-Type': 'application/json',
    },
});


export default api;
