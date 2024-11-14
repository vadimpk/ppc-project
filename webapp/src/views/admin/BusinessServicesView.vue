<template>
  <div class="w-auto ms-4">
    <div class="bg-container p-4 d-flex justify-content-between align-items-center">
      <h2 class="h2 mb-0">Business Services</h2>
      <button @click="openCreateModal" data-bs-toggle="modal" data-bs-target="#manageServiceModal"
              class="btn btn-success">Add New Service
      </button>
      <ServiceModal :service="selectedService" :is-edit-mode="isEditMode"/>
    </div>

    <div class="py-3">
      <div class="p-4 mb-3 bg-container d-flex justify-content-between align-items-center" v-for="service in services"
           :key="service.id">
        <div>
          <h5>{{ service.name }}</h5>
          <p>{{ service.description }}</p>
          <small>Duration: {{ formatDuration(service.duration) }}  |  Price: ${{ service.price / 100 }}</small>
        </div>
        <div>
          <button @click="openEditModal(service)" class="btn btn-sm btn-primary me-2" data-bs-toggle="modal"
                  data-bs-target="#manageServiceModal">Edit
          </button>
          <button @click="deleteService(service.id)" class="btn btn-sm btn-danger">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useBusinessStore } from '@/stores/businessStore';
import ServiceModal from "@/components/admin/ServiceModal.vue";
import {formatDuration} from "@/utils/covertors.js";

const businessStore = useBusinessStore();

// Use computed to make services reactive to businessStore.services changes
const services = computed(() => businessStore.services);

// Watch the business ID and fetch services when it changes
watch(() => businessStore.business.id, async (newId) => {
  if (newId) {
    await businessStore.fetchServices();
  }
});

const isModalVisible = ref(false);
const isEditMode = ref(false);
const selectedService = ref(null);

// Function to open modal for creating a new service
const openCreateModal = () => {
  selectedService.value = null;
  isEditMode.value = false;
  isModalVisible.value = true;
};

// Function to open modal for editing an existing service
const openEditModal = (service) => {
  selectedService.value = service;
  isEditMode.value = true;
  isModalVisible.value = true;
};

// Function to delete a service
const deleteService = async (serviceId) => {
  if (confirm('Are you sure you want to delete this service?')) {
    await businessStore.deleteService(serviceId);
    // No need to manually set services; computed property keeps it in sync
  }
};
</script>


<style scoped>
</style>
