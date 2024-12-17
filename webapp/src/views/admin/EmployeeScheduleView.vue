<template>
  <div class="w-auto ms-4">
    <!-- Header Section -->
    <div class="bg-container p-4 d-flex justify-content-between align-items-center">
      <h2 class="h2 mb-0">
        <i class="bi bi-calendar2-range me-2"></i>
        Schedule Templates
      </h2>
      <button v-if="isAdmin" @click="openCreateModal" class="btn btn-primary btn-lg" data-bs-toggle="modal"
              data-bs-target="#scheduleModal">
        <i class="bi bi-plus-lg me-2"></i>
        Add New Template
      </button>
    </div>

    <!-- Schedule Templates Table -->
    <div class="p-4 mt-3 bg-container">
      <table class="table">
        <thead>
        <tr>
          <th class="table-bg text-primary fs-4">Day of Week</th>
          <th class="table-bg text-primary fs-4">Start Time</th>
          <th class="table-bg text-primary fs-4">End Time</th>
          <th class="table-bg text-primary fs-4">Is Break</th>
          <th class="table-bg text-primary fs-4" v-if="isAdmin">Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="template in scheduleTemplates" :key="template.id">
          <td class="table-bg fs-6">{{ daysOfWeek[template.day_of_week - 1] }}</td>
          <td class="table-bg fs-6">{{ formatTime(template.start_time) }}</td>
          <td class="table-bg fs-6">{{ formatTime(template.end_time) }}</td>
          <td class="table-bg fs-6">{{ template.is_break ? 'Yes' : 'No' }}</td>
          <td class="table-bg fs-6" v-if="isAdmin">
            <button @click="openEditModal(template)" class="btn btn-sm btn-outline-primary me-2" data-bs-toggle="modal"
                    data-bs-target="#scheduleModal">Edit
            </button>
            <button @click="deleteTemplate(template.id)" class="btn btn-sm btn-outline-danger">Delete</button>
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
import {USER_ROLE_ADMIN} from "@/utils/constants.js";
import {formatTime} from "@/utils/covertors.js";

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


// Open modal to create a new template
const openCreateModal = () => {
  selectedTemplate.value = {
    day_of_week: null,
    start_time: '1970-01-01T09:00:00Z',
    end_time: '1970-01-01T17:00:00Z',
    is_break: false
  };
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
.table-bg {
  background: #1e2024;
  padding: 15px;
}
</style>
