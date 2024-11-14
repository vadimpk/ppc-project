<template>
  <div>
    <div
        v-for="service in services"
        :key="service.id"
        class="bg-container p-3 mb-3 d-flex justify-content-between align-items-center"
    >
      <div class="">
        <h5>{{ service.name }}</h5>
        <p v-if="service.description">{{ service.description }}</p>
        <p class="mb-0"><strong>Duration:</strong> {{ formatDuration(service.duration) }} | <strong>Price:</strong>
          ${{ (service.price / 100).toFixed(2) }}</p>
      </div>
      <button
          @click="bookService(service.id, service.business_id)"
          class="btn btn-primary"
          :disabled="!service.is_active"
      >
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
