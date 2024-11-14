<template>
  <div class="w-auto ms-4">
    <!-- Header Section -->
    <div class="bg-container p-4 d-flex justify-content-between align-items-center">
      <h2 class="h2 mb-0">Employees</h2>
      <button @click="generateRegistrationLink" class="btn btn-success">Generate Registration Link</button>
    </div>

    <!-- Employees List -->
    <div class="py-3">
      <div class="p-4 mb-3 bg-container d-flex justify-content-between align-items-center" v-for="employee in employees"
           :key="employee.id">
        <div>
          <h5>{{ employee.user?.full_name }}</h5>
          <p>Specialization: {{ employee.specialization || 'Not Specified' }}</p>
          <small>Joined: {{ new Date(employee.created_at).toLocaleDateString() }}</small>
        </div>
        <div class="d-flex justify-content-center align-items-center">
          <button @click="navigateToSchedule(employee.id)" class="btn btn-primary me-2" data-bs-toggle="modal"
                  data-bs-target="#scheduleModal">Schedule
          </button>
          <button @click="viewEmployee(employee)" class="btn btn-primary me-2" data-bs-toggle="modal"
                  data-bs-target="#employeeModal">View
          </button>
          <button @click="deleteEmployee(employee.id)" class="btn btn-danger">Delete</button>
        </div>
      </div>
    </div>

    <EmployeeModal :employee="selectedEmployee" :services="employeeServices"/>
  </div>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue';
import {useBusinessStore} from '@/stores/businessStore';
import EmployeeModal from "@/components/admin/EmployeeModal.vue";
import {useRoute, useRouter} from "vue-router";

const businessStore = useBusinessStore();

const employees = computed(() => businessStore.employees);
const employeeServices = ref([]);
const selectedEmployee = ref(null);

onMounted(async () => {
  await businessStore.fetchServices();
  await businessStore.fetchEmployees();
});

// Generate a registration link for a new employee
const generateRegistrationLink = async () => {
  const link = businessStore.generateRegistrationLink();
  alert(`Registration link: ${link}`);
};

// Open the modal to view an employee's details and their services
const viewEmployee = async (employee) => {
  selectedEmployee.value = employee;
  console.log(employee);
  employeeServices.value = await businessStore.fetchEmployeeServices(employee.id);
};

const router = useRouter();
const navigateToSchedule = (employeeId) => {
  router.push({name: 'schedule', params: {id: employeeId}});
};

// Delete an employee
const deleteEmployee = async (employeeId) => {
  if (confirm('Are you sure you want to delete this employee?')) {
    await businessStore.deleteEmployee(employeeId);
  }
};
</script>

<style scoped>
</style>
