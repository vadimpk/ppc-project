<template>
  <div class="w-auto ms-4">
    <!-- Header Section -->
    <div class="bg-container p-4 d-flex justify-content-between align-items-center">
      <h2 class="h2 mb-0">Schedule Templates</h2>
      <button v-if="isAdmin" @click="openCreateModal" class="btn btn-success" data-bs-toggle="modal"
              data-bs-target="#scheduleModal">Add New Template
      </button>
    </div>

    <!-- Schedule Templates Table -->
    <div class="p-4 mt-3 bg-container">
      <table class="table table-striped">
        <thead>
        <tr>
          <th>Day of Week</th>
          <th>Start Time</th>
          <th>End Time</th>
          <th>Is Break</th>
          <th v-if="isAdmin">Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="template in scheduleTemplates" :key="template.id">
          <td>{{ daysOfWeek[template.day_of_week - 1] }}</td>
          <td>{{ formatTime(template.start_time) }}</td>
          <td>{{ formatTime(template.end_time) }}</td>
          <td>{{ template.is_break ? 'Yes' : 'No' }}</td>
          <td v-if="isAdmin">
            <button @click="openEditModal(template)" class="btn btn-sm btn-primary me-2" data-bs-toggle="modal"
                    data-bs-target="#scheduleModal">Edit
            </button>
            <button @click="deleteTemplate(template.id)" class="btn btn-sm btn-danger">Delete</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Edit Modal -->
    <EditTemplateModal
        :isEditMode="isEditMode"
        :template="selectedTemplate"
        @save="saveTemplate"
    />
  </div>
</template>

<script setup>
import {ref, onMounted} from 'vue';
import {useBusinessStore} from '@/stores/businessStore';
import EditTemplateModal from "@/components/admin/EditTemplateModal.vue";
import {useRoute} from "vue-router";
import {useUserStore} from "@/stores/userStore.js";
import {USER_ROLE_ADMIN, USER_ROLE_EMPLOYEE} from "@/utils/constants.js";

const route = useRoute();
const employeeId = route.params.id;

const businessStore = useBusinessStore();
const scheduleTemplates = ref([]);
const isEditMode = ref(false);
const selectedTemplate = ref({});

const daysOfWeek = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

onMounted(async () => {
  await fetchScheduleTemplates();
});

const isAdmin = useUserStore().user.role === USER_ROLE_ADMIN

// Fetch schedule templates from the store
const fetchScheduleTemplates = async () => {
  scheduleTemplates.value = await businessStore.fetchScheduleTemplates(employeeId);
};

// Format time to "HH:MM" format
const formatTime = (time) => {
  const date = new Date(time);

  let hours = date.getUTCHours();
  const minutes = date.getUTCMinutes();
  const ampm = hours >= 12 ? 'PM' : 'AM';

  // Convert to 12-hour format
  hours = hours % 12;
  hours = hours ? hours : 12; // The hour '0' should be '12'

  // Format minutes to always have two digits
  const minutesStr = minutes < 10 ? '0' + minutes : minutes;

  return `${hours}:${minutesStr} ${ampm}`;
};


// Open modal to create a new template
const openCreateModal = () => {
  selectedTemplate.value = {day_of_week: null, start_time: '1970-01-01T09:00:00Z', end_time: '1970-01-01T17:00:00Z', is_break: false};
  isEditMode.value = false;
};

// Open modal to edit an existing template
const openEditModal = (template) => {
  selectedTemplate.value = {...template};
  isEditMode.value = true;
};

// Save template (create or update)
const saveTemplate = async (templateData) => {
  if (isEditMode.value) {
    await businessStore.updateScheduleTemplate(employeeId, templateData);
  } else {
    await businessStore.createScheduleTemplate(employeeId, templateData);
  }
  await fetchScheduleTemplates(employeeId);
};

// Delete a schedule template
const deleteTemplate = async (templateId) => {
  if (confirm('Are you sure you want to delete this template?')) {
    await businessStore.deleteScheduleTemplate(employeeId, templateId);
    await fetchScheduleTemplates(employeeId);
  }
};
</script>

<style scoped>

</style>
