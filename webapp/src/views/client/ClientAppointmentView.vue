<template>
  <div class="w-auto ms-4">
    <!-- Header Section with Date Range Filters -->
    <div class="bg-container p-4 d-flex justify-content-between align-items-center">
      <h2 class="h2 mb-0">Appointments</h2>
      <div class="d-flex align-items-center">
        <input type="date" v-model="startDate" class="form-control me-2" placeholder="Start Date"/>
        <input type="date" v-model="endDate" class="form-control me-2" placeholder="End Date"/>
        <button @click="fetchAppointments" class="btn btn-primary">Filter</button>
      </div>
    </div>

    <!-- Appointments List Table -->
    <div class="p-4 mt-3 bg-container">
      <table class="table table-striped">
        <thead>
        <tr>
          <th>Employee</th>
          <th>Service</th>
          <th>Date</th>
          <th>Start Time</th>
          <th>End Time</th>
          <th>Status</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="appointment in appointments" :key="appointment.id">
          <td>{{ appointment.employee?.full_name || 'N/A' }}</td>
          <td>{{ appointment.service?.name || 'N/A' }}</td>
          <td>{{formatToFancyDateString(appointment.start_time)}}</td>
          <td>{{ formatTime(appointment.start_time) }}</td>
          <td>{{ formatTime(appointment.end_time) }}</td>
          <td>{{ STATUS_MAP[appointment.status] }}</td>
          <td>
            <button @click="viewAppointment(appointment)" class="btn btn-sm btn-primary me-2" data-bs-toggle="modal"
                    data-bs-target="#appointmentModal">
              View
            </button>
            <button @click="cancelAppointment(appointment.id)" class="btn btn-sm btn-danger"
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
import AppointmentDetailModal from "@/components/admin/AppointmentDetailModal.vue";
import {useUserStore} from "@/stores/userStore.js";
import {STATUS_MAP} from "@/utils/constants.js";
import {useClientStore} from "@/stores/clientStore.js";
import {useBusinessStore} from "@/stores/businessStore.js";
import {formatToFancyDateString} from "@/utils/covertors.js";

// Helper function to format dates to "YYYY-MM-DD" format
const formatDate = (date) => date.toISOString().split('T')[0];

// Default values: `startDate` set to 7 days before today, `endDate` set to today
const endDate = ref(formatDate(new Date(Date.now() + 14 * 24 * 60 * 60 * 1000)));
const startDate = ref(formatDate(new Date(Date.now() - 24 * 60 * 60 * 1000)));

const clientStore = useClientStore();
const appointments = ref([]);
const selectedAppointment = ref(null);

const userStore = useUserStore();
const businessStore = useBusinessStore();

// Fetch appointments based on date range
const fetchAppointments = async () => {
  appointments.value = await clientStore.fetchAppointments(userStore.user.id, {
    startDate: startDate.value,
    endDate: endDate.value
  });

  appointments.value = appointments.value || [];
  appointments.value.sort((a, b) => a.start_time - b.start_time);
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

// Format time to "HH:MM" format for display
const formatTime = (time) => {
  const date = new Date(time);
  return date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'});
};

// Fetch initial appointments on mount
onMounted(fetchAppointments);
</script>

<style scoped>

</style>
