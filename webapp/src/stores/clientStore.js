import {defineStore} from 'pinia';
import {useToast} from 'vue-toastification';
import api from "@/axios/index.js";

const toast = useToast();

export const useClientStore = defineStore('ClientStore', {
    state: () => ({}),
    actions: {
        async search(query) {
            try {
                const response = await api.get('/businesses/search?search=' + query);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'An error occurred during loading', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                return data;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during registration', {
                    position: 'top-right',
                    timeout: 5000,
                });
            }
        },

        async getBusinessDetails(businessId) {
            try {
                const response = await api.get(`/businesses/${businessId}`);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'An error occurred during loading', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                return data;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during registration', {
                    position: 'top-right',
                    timeout: 5000,
                });
            }
        },
        async getBusinessServices(businessId) {
            try {
                const response = await api.get(`/businesses/${businessId}/services`);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'An error occurred during loading', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                return data;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during registration', {
                    position: 'top-right',
                    timeout: 5000,
                });
            }
        },
        async fetchServiceEmployees(businessId, serviceId) {
            try {
                const response = await api.get(`/businesses/${businessId}/services/${serviceId}/employees`);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'An error occurred during loading', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                return data;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during registration', {
                    position: 'top-right',
                    timeout: 5000,
                });
            }
        },
        async getServiceDetails(businessId, serviceId) {
            try {
                const response = await api.get(`/businesses/${businessId}/services/${serviceId}`);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'An error occurred during loading', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                return data;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during registration', {
                    position: 'top-right',
                    timeout: 5000,
                });
            }
        },
        async fetchAvailableTimeSlots(businessId, serviceId, employeeId, date) {
            const response = await api.get(`/businesses/${businessId}/appointments/slots?service_id=${serviceId}&employee_id=${employeeId}&date=${formatDate(date)}`);
            const {success, data, error} = response.data;

            if (!success) {
                throw new Error(error?.message || 'An error occurred during loading');
            }

            return data;
        },
        async bookService(businessId, bookingData) {
            try {
                const response = await api.post(`/businesses/${businessId}/appointments`, bookingData);
                const {success, data, error} = response.data;

                if (!success) {
                    toast.error(error?.message || 'An error occurred during booking', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                } else {
                    toast.success('Booking successful!', {
                        position: 'top-right',
                        timeout: 5000,
                    });
                }

                return data;
            } catch (error) {
                toast.error(error?.message || 'An error occurred during booking', {
                    position: 'top-right',
                    timeout: 5000,
                });
                throw error;
            }
        },
        async fetchAppointments(userId, { startDate, endDate }) {
            try {
                const params = {};
                if (startDate) params.start_date = startDate;
                if (endDate) params.end_date = endDate;

                const response = await api.get(`/users/${userId}/appointments`, {params});
                if (response.data.success) {
                    this.appointments = response.data.data;
                    return response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to fetch appointments');
                }
            } catch (error) {
                toast.error(error.message || 'Error fetching appointments');
                console.error("Failed to fetch appointments:", error);
            }
        },
    },
    persist: {
        enabled: true,
        strategies: [
            {storage: localStorage, paths: []},
        ],
    },
});

function formatDate(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
}