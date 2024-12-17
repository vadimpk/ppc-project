<template>
  <div ref="modalRef" class="modal fade" id="appointmentModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1"
       aria-labelledby="appointmentModalLabel" aria-hidden="true">

    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="text-center h4">Appointment to <strong>{{ appointment.service?.name || 'N/A' }}</strong></h5>
          <button type="button" class="btn-close" @click="closeModal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <div class="row mb-2">
            <div class="col-6">
              <h5 class="h4 mb-3 text-center">Client</h5>
              <p>
                <span class="text-muted me-2">Name: </span>
                <strong>{{ appointment.client?.full_name || 'N/A' }}</strong>
              </p>
              <p>
                <span class="text-muted me-2">Email: </span>
                <strong>{{ appointment.client?.email || 'N/A' }}</strong>
              </p>
              <p>
                <span class="text-muted me-2">Phone: </span>
                <strong>{{ appointment.client?.phone || 'N/A' }}</strong>
              </p>
            </div>
            <div class="col-6">
              <h5 class="h4 mb-3 text-center">Employee</h5>
              <p>
                <span class="text-muted me-2">Name: </span>
                <strong>{{ appointment.employee?.full_name || 'N/A' }}</strong>
              </p>
              <p>
                <span class="text-muted me-2">Email: </span>
                <strong>{{ appointment.employee?.email || 'N/A' }}</strong>
              </p>
              <p>
                <span class="text-muted me-2">Phone: </span>
                <strong>{{ appointment.employee?.phone || 'N/A' }}</strong>
              </p>
            </div>
          </div>
          <div class="d-flex justify-content-center">
            <p class="badge bg-primary text-dark mb-0 fs-6 me-3">{{ formatTime(appointment.start_time) }} - {{ formatTime(appointment.end_time) }}</p>
            <p class="badge bg-primary text-dark mb-0 fs-6">{{ STATUS_MAP[appointment.status] }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue';
import {formatTime} from "@/utils/covertors.js";
import {Modal} from "bootstrap";
import {STATUS_MAP} from "@/utils/constants.js";

const props = defineProps({
  appointment: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(['close']);

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
