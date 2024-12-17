<template>
  <div>
    <div
        v-for="service in services"
        :key="service.id"
        class="bg-container p-4 mb-3 d-flex justify-content-between align-items-center"
    >
      <div class="">
        <h5>{{ service.name }}</h5>
        <p v-if="service.description">{{ service.description }}</p>
        <p class="badge bg-primary text-dark mb-0 fs-6 me-3">{{ formatDuration(service.duration) }}</p>
        <p class="badge bg-primary text-dark mb-0 fs-6">${{ (service.price / 100).toFixed(2) }}</p>
      </div>
      <button
          @click="bookService(service.id, service.business_id)"
          class="btn btn-primary btn-lg"
          :disabled="!service.is_active"
      >
        <i class="bi bi-journal-check"></i>
        {{ service.is_active ? 'Book' : 'Unavailable' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import {defineProps, defineEmits} from 'vue';

const props = defineProps({
  services: {
    type: Array,
    required: true,
  },
});
const emit = defineEmits(['book']);

// Function to format duration in minutes to "HHh MMm"
const formatDuration = (minutes) => {
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  if (remainingMinutes === 0) {
    return `${hours}h`;
  }
  return `${hours > 0 ? `${hours}h ` : ''}${remainingMinutes}m`;
};

// Function to emit booking event
const bookService = (serviceId, businessId) => {
  emit('book', serviceId, businessId);
};
</script>

<style scoped>
.service-list {
  display: flex;
  flex-direction: column;
}

.service-item {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}
</style>
