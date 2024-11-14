import {defineStore} from 'pinia';
import {useToast} from 'vue-toastification';
import api from "@/axios/index.js";

const toast = useToast();

export const useBusinessStore = defineStore('businessStore', {
    state: () => {
        return {
            business: {
                id: null,
                name: null,
                logo_url: null,
                color_scheme: null,
                created_at: null,
            },
            services: [],
            employees: [],
            employeeServices: [],
            scheduleTemplate: [],
            appointments: [],
        }
    },
    actions: {
        async getBusiness(business_id) {
            try {
                const response = await api.get(`/businesses/${business_id}`);
                if (response.data.success) {
                    this.business = response.data.data;
                    return response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to get business');
                }
            } catch (error) {
                toast.error(error.message || 'Error getting business');
                throw error;
            }
        },
        async updateBusiness(payload) {
            try {
                const response = await api.put(`/businesses/${this.business.id}/`, payload);
                if (response.data.success) {
                    this.business.name = payload.name;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to update business name');
                }
            } catch (error) {
                toast.error(error.message || 'Error updating business name');
                throw error;
            }
        },
        async updateBusinessAppearance(payload) {
            try {
                const response = await api.patch(`/businesses/${this.business.id}/appearance`, payload);
                if (response.data.success) {
                    this.business.logoURL = payload.logoURL;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to update business logo');
                }
            } catch (error) {
                toast.error(error.message || 'Error updating business logo');
                throw error;
            }
        },


        async fetchServices() {
            try {
                const response = await api.get(`/businesses/${this.business.id}/services`);
                if (response.data.success) {
                    console.log("new services", response.data.data);
                    this.services = response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to fetch services');
                }
            } catch (error) {
                toast.error(error.message || 'Error fetching services');
                console.error("Failed to fetch services:", error);
            }
            const toast = useToast();

        },
        async createService(serviceData) {
            try {
                const response = await api.post(`/businesses/${this.business.id}/services`, serviceData);
                if (response.data.success) {
                    console.log("new service", response.data.data);
                    this.services.push(response.data.data);
                } else {
                    throw new Error(response.data.error?.message || 'Failed to create service');
                }
            } catch (error) {
                toast.error(error.message || 'Error creating service');
                console.error("Failed to create service:", error);
            }
        },
        async updateService(serviceId, serviceData) {
            try {
                const response = await api.put(`/businesses/${this.business.id}/services/${serviceId}`, serviceData);
                if (response.data.success) {
                    const index = this.services.findIndex(service => service.id === serviceId);
                    if (index !== -1) {
                        this.services[index] = response.data.data;
                    }
                } else {
                    throw new Error(response.data.error?.message || 'Failed to update service');
                }
            } catch (error) {
                toast.error(error.message || 'Error updating service');
                console.error("Failed to update service:", error);
            }
        },
        async deleteService(serviceId) {
            try {
                const response = await api.delete(`/businesses/${this.business.id}/services/${serviceId}`);
                if (response.data.success) {
                    this.services = this.services.filter(service => service.id !== serviceId);
                } else {
                    throw new Error(response.data.error?.message || 'Failed to delete service');
                }
            } catch (error) {
                toast.error(error.message || 'Error deleting service');
                console.error("Failed to delete service:", error);
            }
        },


        // Employees
        async fetchEmployeeAppointments(employeeId, { startDate, endDate }) {
            try {
                const params = {};
                if (startDate) params.start_date = startDate;
                if (endDate) params.end_date = endDate;

                const response = await api.get(`/businesses/${this.business.id}/appointments/employee/${employeeId}`, {params});
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

        async fetchEmployees() {
            try {
                const response = await api.get(`/businesses/${this.business.id}/employees`);
                if (response.data.success) {
                    this.employees = response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to fetch employees');
                }
            } catch (error) {
                toast.error(error.message || 'Error fetching employees');
                console.error("Failed to fetch employees:", error);
            }
        },
        generateRegistrationLink() {
            return `${window.location.origin}/auth/register?business_id=${this.business.id}`;
        },
        async fetchEmployeeServices(employeeId) {
            try {
                const response = await api.get(`/businesses/${this.business.id}/employees/${employeeId}/services`);
                if (response.data.success) {
                    this.employeeServices = response.data.data;
                    return response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to fetch employee services');
                }
            } catch (error) {
                toast.error(error.message || 'Error fetching employee services');
                console.error("Failed to fetch employee services:", error);
            }
        },
        async assignService(employeeId, serviceId) {
            try {
                const response = await api.post(`/businesses/${this.business.id}/employees/${employeeId}/services`, {service_ids: [serviceId]});
                if (response.data.success) {
                    this.employeeServices.push(this.services.find(service => service.id === serviceId));
                    return response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to assign service');
                }
            } catch (error) {
                toast.error(error.message || 'Error assigning service');
                console.error("Failed to assign service:", error);
            }
        },
        async removeService(employeeId, serviceId) {
            try {
                const response = await api.delete(`/businesses/${this.business.id}/employees/${employeeId}/services`, {data: {service_ids: [serviceId]}});
                if (response.data.success) {
                    this.employeeServices = this.employeeServices.filter(service => service.id !== serviceId);
                    return response.data.data;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to remove service');
                }
            } catch (error) {
                toast.error(error.message || 'Error removing service');
                console.error("Failed to remove service:", error);
            }
        },
        async deleteEmployee(employeeId) {
            try {
                const response = await api.delete(`/businesses/${this.business.id}/employees/${employeeId}`);
                if (response.data.success) {
                    this.employees = this.employees.filter(employee => employee.id !== employeeId);
                } else {
                    throw new Error(response.data.error?.message || 'Failed to delete employee');
                }
            } catch (error) {
                toast.error(error.message || 'Error deleting employee');
                console.error("Failed to delete employee:", error);
            }
        },


        async fetchScheduleTemplates(employeeId) {
            try {
                const response = await api.get(`/businesses/${this.business.id}/employees/${employeeId}/schedule/templates`);
                if (response.data.success) {
                    this.scheduleTemplates = response.data.data;
                    return this.scheduleTemplates;
                } else {
                    throw new Error(response.data.error?.message || 'Failed to fetch schedule templates');
                }
            } catch (error) {
                toast.error(error.message || 'Error fetching schedule templates');
                console.error("Failed to fetch schedule templates:", error);
            }
        },
        async deleteScheduleTemplate(employeeId, templateId) {
            try {
                const response = await api.delete(`/businesses/${this.business.id}/employees/${employeeId}/schedule/templates/${templateId}`);
                if (response.data.success) {
                    this.scheduleTemplates = this.scheduleTemplates.filter(t => t.id !== templateId);
                } else {
                    throw new Error(response.data.error?.message || 'Failed to delete schedule template');
                }
            } catch (error) {
                toast.error(error.message || 'Error deleting schedule template');
                console.error("Failed to delete schedule template:", error);
            }
        },
        async updateScheduleTemplate(employeeId, updatedTemplate) {
            try {
                const response = await api.put(`/businesses/${this.business.id}/employees/${employeeId}/schedule/templates/${updatedTemplate.id}`, updatedTemplate);
                if (response.data.success) {
                    const index = this.scheduleTemplates.findIndex(template => template.id === updatedTemplate.id);
                    if (index !== -1) {
                        this.scheduleTemplates[index] = updatedTemplate;
                    }
                } else {
                    throw new Error(response.data.error?.message || 'Failed to update schedule template');
                }
            } catch (error) {
                toast.error(error.message || 'Error updating schedule template');
                console.error("Failed to update schedule template:", error);
            }
        },
        async createScheduleTemplate(employeeId, createTemplate) {
            try {
                const response = await api.post(`/businesses/${this.business.id}/employees/${employeeId}/schedule/templates`, createTemplate);
                if (response.data.success) {
                    this.scheduleTemplates.push(response.data.data);
                } else {
                    throw new Error(response.data.error?.message || 'Failed to update schedule template');
                }
            } catch (error) {
                toast.error(error.message || 'Error updating schedule template');
                console.error("Failed to update schedule template:", error);
            }
        },

        async fetchAppointments({ startDate, endDate }) {
            try {
                const params = {};
                if (startDate) params.start_date = startDate;
                if (endDate) params.end_date = endDate;

                const response = await api.get(`/businesses/${this.business.id}/appointments/`, { params });
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
        async cancelAppointment(appointmentId) {
            try {
                const response = await api.delete(`/businesses/${this.business.id}/appointments/${appointmentId}`);
                if (response.data.success) {
                    this.appointments = this.appointments.filter(appointment => appointment.id !== appointmentId);
                } else {
                    throw new Error(response.data.error?.message || 'Failed to cancel appointment');
                }
            } catch (error) {
                toast.error(error.message || 'Error cancelling appointment');
                console.error("Failed to cancel appointment:", error);
            }
        },
    },
    persist: {
        enabled: true,
        strategies: [{storage: localStorage, paths: ['business', 'services', "employees"]}],
    },
});