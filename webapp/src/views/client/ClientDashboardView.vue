<template>
  <div class="w-auto ms-4">
    <!-- Search Section -->
    <div class="bg-container p-4 d-flex justify-content-center align-items-center">
      <div class="input-group input-group-lg">
        <input
            type="text"
            v-model="searchQuery"
            placeholder="Search for services or businesses..."
            class="form-control search-input form-control-lg"
        />
        <button @click="search" class="btn btn-primary me-2">Search</button>
      </div>
    </div>

    <!-- Results Section -->
    <div class="mt-3">
      <div v-if="isLoading" class="text-center">
        <div class="spinner-border text-primary" role="status">
        </div>
      </div>

      <div v-if="results.services.length > 0" class="bg-container p-4 mb-3">
        <p class="h4 mb-0 text-primary">
          <i class="bi bi-card-checklist me-2"></i>
          Services
        </p>
      </div>
      <div class="">
        <BusinessServiceList :services="results.services" @book="handleBooking"/>
      </div>

      <div v-if="results.businesses.length > 0" class="bg-container p-4 mb-3">
        <h4 class="h4 mb-0 text-primary">
          <i class="bi bi-card-checklist me-2"></i>
          Businesses
        </h4>
      </div>
      <div class="">
        <div v-for="business in results.businesses" :key="business.id"
             class="business-item p-3 bg-container mb-3 d-flex justify-content-between align-items-center">
          <div class="d-flex align-items-center">
            <img :src="business.logo_url" :alt="business.name" class="logo me-3"
                 v-if="business.logo_url"/>
            <span class="h4 mb-0">{{ business.name }}</span>
          </div>
          <button @click="viewBusiness(business.id)" class="btn btn-primary btn-lg">
            <i class="bi bi-eye"></i>
            View
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue';
import {useClientStore} from '@/stores/clientStore';
import BusinessServiceList from "@/components/client/BusinessServiceList.vue";
import {useRouter} from "vue-router";

const clientStore = useClientStore();
const searchQuery = ref('');
const results = ref({
  services: [],
  businesses: [],
});
const isLoading = ref(false);

// Perform search based on query and type
const search = async () => {
  if (!searchQuery.value) {
    return;
  }

  isLoading.value = true;
  results.value = {
    services: [],
    businesses: [],
  }

  try {
    results.value = await clientStore.search(searchQuery.value);
  } catch (error) {
    console.error('Error performing search:', error);
  } finally {
    isLoading.value = false;
  }
};

// View a business
const router = useRouter();
const viewBusiness = (businessId) => {
  router.push({name: 'clientBusiness', params: {id: businessId}});
};


const handleBooking = (serviceId, businessId) => {
  router.push({name: 'newAppointment', params: {business_id: businessId, service_id: serviceId}});
};

const handleBookingConfirmation = ({date, timeSlot}) => {
  alert(`Booking confirmed for ${date.toLocaleDateString()} at ${timeSlot}`);
};
</script>

<style scoped>
.search-panel {
  width: 40%;
  min-width: 200px;
}

.search-panel {
  max-width: 600px;
  width: 100%;
}

.logo {
  height: 40px;
}
</style>
