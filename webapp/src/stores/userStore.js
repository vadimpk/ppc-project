import {defineStore} from 'pinia';
import {useToast} from 'vue-toastification';
import api from "@/axios/index.js";

const toast = useToast();

export const useUserStore = defineStore('userStore', {
    state: () => ({
        token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJidXNpbmVzc19pZCI6NCwicm9sZSI6ImVtcGxveWVlIiwiZXhwIjoxNzMxNTk1MTIyLCJpYXQiOjE3MzE1MDg3MjJ9.AW3mIqU_jMHLEvttukMlvFMUIbOVgq8yJV1PlTgxbic",
        user: {
            id: 2,
            business_id: 4,
            employee_id: null,
            role: 'employee',
        }
    }),
    actions: {
        async registerUser(payload) {
            try {
                const response = await api.post('auth/register', payload);
                const {success, data, error} = response.data;
                console.log(data)

                if (!success) {
                    toast.error(error?.message || 'An error occurred during registration', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                this.token = data.token;
                this.user = data.user;
                return data.user;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during registration', {
                    position: 'top-right',
                    timeout: 5000,
                });
            }
        },
        async loginUser(payload) {
            try {
                const response = await api.post('auth/login', payload);
                const {success, data, error} = response.data;

                if (!success) {
                    throw new Error(error?.message || 'Invalid credentials');
                }

                this.token = data.token;
                this.user = data.user;
            } catch (error) {
                toast.error('Invalid credentials', {
                    position: 'top-right',
                    timeout: 5000,
                });
                throw new Error(error || 'Invalid credentials');
            }
        },
        async updateUserProfile(payload) {
            try {
                const response = await api.put(`users/${this.user.id}/`, payload);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'This email or phone number is already in use', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                    throw new Error(error?.message || 'Profile update failed');
                } else {
                    toast.success('Profile updated successfully', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }
                this.user = data.user;
                this.token = data.token;

                console.log(data)
                console.log(this.user)
            } catch (error) {
                throw new Error(error.response?.data?.error?.message || 'Profile update failed');
            }
        }
    },
    persist: {
        enabled: true,
        strategies: [
            {storage: localStorage, paths: ['token', 'user']},
        ],
    },
});
