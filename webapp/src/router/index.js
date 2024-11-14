import {createRouter, createWebHistory} from 'vue-router'
import RegisterView from "@/views/RegisterView.vue";
import LoginView from "@/views/LoginView.vue";
import {useUserStore} from "@/stores/userStore.js";
import {USER_ROLE_CLIENT, USER_ROLE_EMPLOYEE} from "@/utils/constants.js";
import AdminLayout from "@/layouts/AdminLayout.vue";
import BusinessServicesView from "@/views/admin/BusinessServicesView.vue";
import EmployeesView from "@/views/admin/EmployeesView.vue";
import EmployeeScheduleView from "@/views/admin/EmployeeScheduleView.vue";
import AppointmentView from "@/views/admin/AppointmentView.vue";
import EmployeesLayout from "@/layouts/EmployeesLayout.vue";
import EmployeeAppointmentView from "@/views/employee/EmployeeAppointmentView.vue";
import ClientLayout from "@/layouts/ClientLayout.vue";
import ClientAppointmentView from "@/views/client/ClientAppointmentView.vue";
import ClientDashboard from "@/views/client/ClientDashboardView.vue";
import ClientBusinessView from "@/views/client/ClientBusinessView.vue";
import BookServiceView from "@/views/client/BookServiceView.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: "/auth",
            children: [
                {
                    path: "register",
                    name: "register",
                    component: RegisterView,
                },
                {
                    path: "login",
                    name: "login",
                    component: LoginView,
                }
            ]
        },

        {
            path: '/admin',
            component: AdminLayout,
            meta: {requiresAuth: true, role: 'admin'},
            children: [
                {
                    path: 'services',
                    component: BusinessServicesView,
                },
                {
                    path: 'employees',
                    component: EmployeesView,
                },
                {
                    path: 'schedule/:id',
                    name: 'schedule',
                    component: EmployeeScheduleView,
                },
                {
                    path: 'appointments',
                    component: AppointmentView,
                }
            ],
        },
        {
            path: '/employee',
            meta: {requiresAuth: true, role: 'employee'},
            component: EmployeesLayout,
            children: [
                {
                    path: 'schedule/:id',
                    component: EmployeeScheduleView,
                },
                {
                    path: 'appointments',
                    component: EmployeeAppointmentView,
                }
            ],
        },
        {
            path: '/client',
            meta: {requiresAuth: true, role: 'client'},
            component: ClientLayout,
            children: [
                {
                    path: 'dashboard',
                    component: ClientDashboard,
                },
                {
                    path: 'appointments',
                    component: ClientAppointmentView,
                },
                {
                    path: 'business/:id',
                    name: 'clientBusiness',
                    component: ClientBusinessView,
                },
                {
                    path: "appointments/new/:business_id/:service_id",
                    name: "newAppointment",
                    component: BookServiceView,
                }
            ],
        },

        {
            path: '/:catchAll(.*)',
            redirect: '/auth/login',
        },
    ],
})

// Add a navigation guard
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()

    // Redirect to login if user is not authenticated
    if (to.meta.requiresAuth && (!userStore.token || !userStore.user)) {
        return next({path: '/auth/login'})
    }

    // Restrict access based on user role
    if (to.meta.role) {
        const userRole = userStore.user?.role

        if (to.meta.role === USER_ROLE_CLIENT && userRole !== USER_ROLE_CLIENT && userRole !== 'admin') {
            console.error("Unknown role: ", userRole)
            return next({path: '/auth/login'})
        }

        if (to.meta.role !== userRole) {
            return next({path: userRole + '/dashboard'})
        }

        console.log("AB", to.path, to.params.id)
        if (to.meta.role === USER_ROLE_EMPLOYEE && to.path.startsWith('/employee/schedule') && to.params.id !== userStore.user?.employee_id.toString()) {
            return next({path: userRole + '/schedule/' + userStore.user?.employee_id})
        }
    }

    // Proceed to the requested route
    next()
})

export default router
