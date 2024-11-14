<template>
  <div class="w-auto ms-4">
    <!-- Business Header -->
    <div class="bg-container p-4 mb-3">
      <h2 class="h4">{{ business.name || 'Business Details' }}</h2>
      <p v-if="business.description">{{ business.description }}</p>
    </div>

    <div class="services-section">
      <BusinessServiceList :services="services" @book="handleBooking" />

      <div v-if="!services.length && !isLoading" class="text-center mt-4">
        <p>No services found for this business.</p>
      </div>

      <div v-if="isLoading" class="text-center mt-4">
        <p>Loading services...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import {useRoute, useRouter} from 'vue-router';
import { useClientStore } from '@/stores/clientStore';
import BusinessServiceList from "@/components/client/BusinessServiceList.vue";

const clientStore = useClientStore();
const route = useRoute();
const businessId = route.params.id;

const business = ref({});
const services = ref([]);
const isLoading = ref(false);

// Fetch business details and services on mount
onMounted(async () => {
  await fetchBusinessDetails();
  await fetchBusinessServices();
});

// Fetch business details by ID
const fetchBusinessDetails = async () => {
  isLoading.value = true;
  try {
    business.value = await clientStore.getBusinessDetails(businessId);
  } catch (error) {
    console.error('Error fetching business details:', error);
    alert('An error occurred while fetching business details.');
  } finally {
    isLoading.value = false;
  }
};

// Fetch services offered by this business
const fetchBusinessServices = async () => {
  isLoading.value = true;
  try {
    services.value = await clientStore.getBusinessServices(businessId);
  } catch (error) {
    console.error('Error fetching business services:', error);
    alert('An error occurred while fetching business services.');
  } finally {
    isLoading.value = false;
  }
};

// Handle booking action for a service
const router = useRouter();
const handleBooking = (serviceId, businessId) => {
  router.push({name: 'newAppointment', params: {business_id: businessId, service_id: serviceId}});
};
</script>

<style scoped>

</style>
