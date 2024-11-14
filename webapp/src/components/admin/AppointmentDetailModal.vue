<template>
  <div ref="modalRef" class="modal fade" id="appointmentModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1"
       aria-labelledby="appointmentModalLabel" aria-hidden="true">

    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title" id="appointmentModalLabel">Appointment Details</h4>
          <button type="button" class="btn-close" @click="closeModal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <!-- Appointment Information -->
          <h5>Appointment to <strong>{{ appointment.service?.name || 'N/A' }}</strong></h5>
          <div class="row">
            <div class="col-6">
              <p><strong>Start Time:</strong> {{ formatTime(appointment.start_time) }}</p>
            </div>
            <div class="col-6">
              <p><strong>End Time:</strong> {{ formatTime(appointment.end_time) }}</p>
            </div>
          </div>
          <div class="row">
            <div class="col-6">
              <p><strong>Status:</strong> {{ appointment.status }}</p>
            </div>
            <div class="col-6">
              <p><strong>Reminder Time:</strong>
                {{ appointment.reminder_time ? appointment.reminder_time + ' minutes' : 'None' }}</p>
            </div>
          </div>

          <div class="row">
            <div class="col-6">
              <h5 class="h5 mt-4">Client Information</h5>
              <p><strong>Name:</strong> {{ appointment.client?.full_name || 'N/A' }}</p>
              <p><strong>Email:</strong> {{ appointment.client?.email || 'N/A' }}</p>
              <p><strong>Phone:</strong> {{ appointment.client?.phone || 'N/A' }}</p>
              <p><strong>Registered At:</strong> {{ formatDate(appointment.client?.created_at) }}</p>
            </div>
            <div class="col-6">
              <h5 class="h5 mt-4">Employee Information</h5>
              <p><strong>Name:</strong> {{ appointment.employee?.full_name || 'N/A' }}</p>
              <p><strong>Email:</strong> {{ appointment.employee?.email || 'N/A' }}</p>
              <p><strong>Phone:</strong> {{ appointment.employee?.phone || 'N/A' }}</p>
              <p><strong>Joined At:</strong> {{ formatDate(appointment.employee?.created_at) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue';
import {parseDurationToInt} from "@/utils/covertors.js";
import {Modal} from "bootstrap";

const props = defineProps({
  appointment: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(['close']);

// Format time for display
const formatTime = (time) => {
  const date = new Date(time);
  return date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'});
};

// Format date for display
const formatDate = (date) => {
  return date ? new Date(date).toLocaleDateString() : 'N/A';
};

const modalRef = ref(null);
const closeModal = async () => {
  Modal.getInstance(modalRef.value)?.hide()
  emit('close');
};
</script>

<style scoped>

</style>
