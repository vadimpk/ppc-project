<template>
  <div class="w-auto ms-4">
    <!-- Header Section -->
    <div class="bg-container p-4 mb-3">
      <div class="d-flex align-items-center mb-3">
        <img :src="business.logo_url" :alt="business.name" class="logo me-3"
             v-if="business.logo_url"/>
        <span class="h3 mb-0">{{ business.name || 'Business Details' }}</span>
      </div>
      <h2 class="h4">{{ service.name || 'Service Name' }}</h2>
      <p class="text-muted">{{ service.description }}</p>

      <p class="badge bg-primary text-dark mb-0 fs-6 me-3">{{ formatDuration(service.duration) }}</p>
      <p class="badge bg-primary text-dark mb-0 fs-6">${{ (service.price / 100).toFixed(2) }}</p>
    </div>

    <div class="row">
      <div class="col-6 pe-3">
        <div class="bg-container p-4">
          <h4>Select an Employee</h4>
          <div v-if="employees.length" class="d-flex flex-column">
            <label
                v-for="employee in employees"
                :key="employee.id"
                class="form-check-label fs-5"
            >
              <input
                  type="radio"
                  :value="employee.id"
                  v-model="selectedEmployee"
                  class="form-check-input me-2"
              />
              {{ employee.user.full_name }}
            </label>
          </div>
          <div v-else>
            <p>No employees available for this service.</p>
          </div>
        </div>
      </div>

      <!-- Date Display and Navigation -->
      <div class="col-6">
      <div class="bg-container p-4" v-if="selectedEmployee">
        <div class="d-flex justify-content-center mb-3">
          <div class="text-end me-3">
            <button @click="prevDay" class="btn btn-sm btn-outline-primary" :disabled="isToday">Previous</button>
          </div>
          <div class="text-center d-flex justify-content-center me-3 h-100">
            <h6 class="mb-0">{{ formattedDate }}</h6>
          </div>
          <div class="">
            <button @click="nextDay" class="btn btn-sm btn-outline-primary">Next</button>
          </div>
        </div>

        <!-- Time Slots Grid -->
        <div class="time-slots-grid">
          <div
              v-for="(slot, index) in timeSlots"
              :key="index"
              class="time-slot"
              :class="{ 'selected': selectedSlot === slot }"
              @click="selectSlot(slot)"
          >
            {{ formatSlot(slot) }}
          </div>
        </div>

        <div class="d-flex justify-content-center align-items-center" v-if="isNotAvailable">
          <p>No available time slots for this date.</p>
        </div>

        <!-- Confirm Booking Button -->
        <div class="mt-4 d-flex justify-content-center">
          <button @click="confirmBooking" class="btn btn-primary" :disabled="!selectedSlot || !selectedEmployee">
            Confirm Booking
          </button>
        </div>
      </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, computed, onMounted, watch} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import {useClientStore} from '@/stores/clientStore.js';
import {formatDuration} from "@/utils/covertors.js";
import {useToast} from "vue-toastification";

const clientStore = useClientStore();
const route = useRoute();
const router = useRouter();

// Extract business_id and service_id from URL parameters
const businessId = parseInt(route.params.business_id);
const serviceId = parseInt(route.params.service_id);

const selectedDate = ref(new Date());
const selectedSlot = ref(null);
const selectedEmployee = ref(null);
const timeSlots = ref([]);
const employees = ref([]);
const business = ref({});
const service = ref({});
const isLoading = ref(true);
const isNotAvailable = ref(false);

// Formatted date for display
const formattedDate = computed(() => {
  return selectedDate.value.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
});

// Check if selected date is today
const isToday = computed(() => {
  const today = new Date();
  return selectedDate.value.toDateString() === today.toDateString();
});

// Fetch initial data on mount
onMounted(async () => {
  await fetchBusinessDetails();
  await fetchServiceDetails();
  await fetchEmployees();
});

// Fetch business details
const fetchBusinessDetails = async () => {
  try {
    business.value = await clientStore.getBusinessDetails(businessId);
  } catch (error) {
    console.error('Failed to fetch business details:', error);
    alert('Error fetching business details.');
  }
};

// Fetch service details
const fetchServiceDetails = async () => {
  service.value = await clientStore.getServiceDetails(businessId, serviceId);
};

// Fetch employees for the selected service
const fetchEmployees = async () => {
  employees.value = (await clientStore.fetchServiceEmployees(businessId, serviceId)).filter(employee => employee.is_active);
};

// Fetch available time slots for the selected date
const updateTimeSlots = async (employeeId) => {
  try {
    timeSlots.value = await clientStore.fetchAvailableTimeSlots(businessId, serviceId, employeeId, selectedDate.value);
    isNotAvailable.value = timeSlots.value.length === 0;
  } catch (error) {
    timeSlots.value = [];
    isNotAvailable.value = true;
    console.error('Failed to fetch time slots:', error);
  }
  isLoading.value = false;
};

watch(() => selectedEmployee?.value, async (newEmployeeId) => {
  console.log("Selected Employee", selectedEmployee.value);
  if (newEmployeeId && selectedDate.value) {
    await updateTimeSlots(newEmployeeId);
  }
}, {immediate: true});


// Navigate to the previous day
const prevDay = () => {
  const newDate = new Date(selectedDate.value);
  newDate.setDate(newDate.getDate() - 1);
  selectedDate.value = newDate;
  selectedSlot.value = null;
  updateTimeSlots(selectedEmployee.value);
};

// Navigate to the next day
const nextDay = () => {
  const newDate = new Date(selectedDate.value);
  newDate.setDate(newDate.getDate() + 1);
  selectedDate.value = newDate;
  selectedSlot.value = null;
  updateTimeSlots(selectedEmployee.value);
};

// Select a time slot
const selectSlot = (slot) => {
  selectedSlot.value = slot;
};

const formatSlot = (slot) => {
  return `${formatToTime(slot.start_time)} - ${formatToTime(slot.end_time)}`;
};

const formatToTime = (datetimeString) => {
  const date = new Date(datetimeString);
  const hours = String(date.getUTCHours()).padStart(2, '0');
  const minutes = String(date.getUTCMinutes()).padStart(2, '0');

  return `${hours}:${minutes}`;
}

// Confirm booking and redirect or perform further actions
const toast = useToast();
const confirmBooking = async () => {
  await clientStore.bookService(businessId, {
    employee_id: selectedEmployee.value,
    service_id: serviceId,
    start_time: combineDateAndTime(selectedDate.value, selectedSlot.value.start_time),
  });
  console.log("Booking confirmed:", combineDateAndTime(selectedDate.value, selectedSlot.value.start_time));
  await router.push('/client/appointments');
};


function combineDateAndTime(dateObj, timeString) {
  // Extract the date components from the dateObj
  const year = dateObj.getFullYear();
  const month = String(dateObj.getMonth() + 1).padStart(2, '0');
  const day = String(dateObj.getDate()).padStart(2, '0');

  // Extract the time components from the time string
  const timePart = timeString.split("T")[1]; // e.g., "07:00:00Z"

  // Combine the date and time in the desired format
  return `${year}-${month}-${day}T${timePart}`;
}
</script>

<style scoped>
.time-slots-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.time-slot {
  padding: 10px;
  background-color: #1e2024;
  border-radius: 10px;
  text-align: center;
  cursor: pointer;
}

.time-slot.selected {
  background-color: #bb86fc;
  color: #1e2024;
}

.logo {
  height: 40px;
}
</style>
