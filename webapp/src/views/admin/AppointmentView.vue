<template>
  <div class="w-auto ms-4">
    <!-- Header Section with Date Range Filters -->
    <div class="bg-container p-4 d-flex justify-content-between align-items-center">
      <h2 class="h2 mb-0">
        <i class="bi bi-calendar-check me-2"></i>
        Appointments
      </h2>
      <div class="d-flex align-items-center">
        <input type="date" v-model="startDate" class="form-control me-2" placeholder="Start Date"/>
        <input type="date" v-model="endDate" class="form-control me-2" placeholder="End Date"/>
        <button @click="fetchAppointments" class="btn btn-primary">Filter</button>
      </div>
    </div>

    <!-- Appointments List Table -->
    <div class="p-4 mt-3 bg-container">
      <table class="table">
        <thead>
        <tr>
          <th class="table-bg text-primary fs-4">Client</th>
          <th class="table-bg text-primary fs-4">Employee</th>
          <th class="table-bg text-primary fs-4">Service</th>
          <th class="table-bg text-primary fs-4">Date</th>
          <th class="table-bg text-primary fs-4">Time</th>
          <th class="table-bg text-primary fs-4">Status</th>
          <th class="table-bg text-primary fs-4">Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="appointment in appointments" :key="appointment.id">
          <td class="table-bg fs-6">{{ appointment.client?.full_name || 'N/A' }}</td>
          <td class="table-bg fs-6">{{ appointment.employee?.full_name || 'N/A' }}</td>
          <td class="table-bg fs-6">{{ appointment.service?.name || 'N/A' }}</td>
          <td class="table-bg fs-6">{{formatToFancyDateString(appointment.start_time)}}</td>
          <td class="table-bg fs-6">{{ formatTime(appointment.start_time) }} - {{ formatTime(appointment.end_time) }}</td>
          <td class="table-bg fs-6">{{ STATUS_MAP[appointment.status] }}</td>
          <td class="table-bg fs-6">
            <button @click="viewAppointment(appointment)" class="btn btn-sm btn-outline-primary me-2" data-bs-toggle="modal"
                    data-bs-target="#appointmentModal">
              <i class="bi bi-eye"></i>
              View
            </button>
            <button @click="cancelAppointment(appointment.id)" class="btn btn-sm btn-outline-danger"
                    :disabled="appointment.status === 'cancelled' || appointment.status === 'completed'">Cancel
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Appointment Detail Modal -->
    <AppointmentDetailModal v-if="selectedAppointment" :appointment="selectedAppointment"
                            @close="selectedAppointment = null"/>
  </div>
</template>

<script setup>
import {ref, onMounted} from 'vue';
import {useBusinessStore} from '@/stores/businessStore';
import AppointmentDetailModal from "@/components/admin/AppointmentDetailModal.vue";
import {STATUS_MAP} from "@/utils/constants.js";
import {formatTime, formatToFancyDateString} from "@/utils/covertors.js";

// Helper function to format dates to "YYYY-MM-DD" format
const formatDate = (date) => date.toISOString().split('T')[0];

// Default values: `startDate` set to 7 days before today, `endDate` set to today
const endDate = ref(formatDate(new Date(Date.now() + 14 * 24 * 60 * 60 * 1000)));
const startDate = ref(formatDate(new Date(Date.now() - 24 * 60 * 60 * 1000)));

const businessStore = useBusinessStore();
const appointments = ref([]);
const selectedAppointment = ref(null);

// Fetch appointments based on date range
const fetchAppointments = async () => {
  appointments.value = await businessStore.fetchAppointments({startDate: startDate.value, endDate: endDate.value});
};

// Open modal to view detailed appointment information
const viewAppointment = (appointment) => {
  selectedAppointment.value = appointment;
};

// Cancel an appointment
const cancelAppointment = async (appointmentId) => {
  if (confirm('Are you sure you want to cancel this appointment?')) {
    await businessStore.cancelAppointment(appointmentId);
    await fetchAppointments(); // Refresh appointments after cancellation
  }
};

// Fetch initial appointments on mount
onMounted(fetchAppointments);
</script>

<style scoped>
.table-bg {
  background: #1e2024;
  padding: 15px;
}
</style>
