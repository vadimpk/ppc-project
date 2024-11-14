import axios from 'axios';
import {useUserStore} from '@/stores/userStore.js';
import {useRouter} from 'vue-router';

// Create an instance of axios
const api = axios.create({
    baseURL: "http://localhost:8080/api/v1",
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor to add JWT token
api.interceptors.request.use((config) => {
    const userStore = useUserStore();

    // Check if there's a token and add it to headers
    if (userStore.token) {
        config.headers['Authorization'] = `Bearer ${userStore.token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

// Response interceptor to handle 401 errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        const router = useRouter();
        const userStore = useUserStore();

        if (error.response && error.response.status === 401) {
            // Clear user data and token
            userStore.$reset();

            // Redirect to the login page
            router.push('/auth/login');
        }
        return Promise.reject(error);
    }
);

export default api;
