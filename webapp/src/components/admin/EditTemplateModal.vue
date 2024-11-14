<template>
  <div ref="modalRef" class="modal fade" id="scheduleModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1"
       aria-labelledby="scheduleModalLabel" aria-hidden="true">

    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="scheduleModalLabel">
            {{ isEditMode ? 'Edit Schedule Template' : 'Add New Schedule Template' }}
          </h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="save">
            <!-- Day of Week Selection (Only for Create Mode) -->
            <div class="mb-3" v-if="!isEditMode">
              <label for="day_of_week" class="form-label">Day of the Week</label>
              <select v-model="templateData.day_of_week" class="form-select" id="day_of_week" required>
                <option value="" disabled>Select a Day</option>
                <option v-for="(day, index) in daysOfWeek" :key="index" :value="index + 1">{{ day }}</option>
              </select>
            </div>

            <!-- Start Time Input -->
            <div class="mb-3">
              <label for="start_time" class="form-label">Start Time</label>
              <input
                  type="time"
                  v-model="templateData.start_time_fancy"
                  :disabled="templateData.is_break"
                  class="form-control"
                  id="start_time"
                  required
              />
            </div>

            <!-- End Time Input -->
            <div class="mb-3">
              <label for="end_time" class="form-label">End Time</label>
              <input
                  type="time"
                  v-model="templateData.end_time_fancy"
                  :disabled="templateData.is_break"
                  class="form-control"
                  id="end_time"
                  required
              />
            </div>

            <!-- Break Toggle -->
            <div class="form-check mb-3">
              <input type="checkbox" v-model="templateData.is_break" class="form-check-input" id="is_break"/>
              <label class="form-check-label" for="is_break">Is Break</label>
            </div>

            <button type="submit" class="btn btn-primary w-100">{{
                isEditMode ? 'Save Changes' : 'Add Template'
              }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, watch} from 'vue';
import {Modal} from "bootstrap";

// Props and events for the modal
const props = defineProps({
  isEditMode: Boolean,
  template: {
    type: Object,
    default: () => ({
      day_of_week: null,
      start_time: '',
      start_time_fancy: '',
      end_time: '',
      end_time_fancy: '',
      is_break: false
    }),
  },
});
const emit = defineEmits(['save']);

// Days of the week options

const daysOfWeek = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]

// Reactive local copy of the template data
const templateData = ref({...props.template});

function formatToHHMM(isoTimeString) {
  const date = new Date(isoTimeString);

  let hours = date.getUTCHours();
  const minutes = date.getUTCMinutes();

  console.log("H", hours, minutes);

  // Format hours and minutes to always have two digits
  const hoursStr = hours < 10 ? '0' + hours : hours;
  const minutesStr = minutes < 10 ? '0' + minutes : minutes;

  return `${hoursStr}:${minutesStr}`;
}


// Watch for changes to the `template` prop and update local data accordingly
watch(
    () => props.template,
    (newTemplate) => {
      templateData.value = {...newTemplate};
      templateData.value.start_time_fancy = formatToHHMM(newTemplate.start_time);
      templateData.value.end_time_fancy = formatToHHMM(newTemplate.end_time)
    },
    {immediate: true}
);

const modalRef = ref(null);
const save = () => {
  // Validate day_of_week if in create mode
  if (!props.isEditMode && !templateData.value.day_of_week) {
    alert('Please select a day of the week.');
    return;
  }

  templateData.value.start_time = `2021-01-01T${templateData.value.start_time_fancy}:00.000Z`
  templateData.value.end_time = `2021-01-01T${templateData.value.end_time_fancy}:00.000Z`

  console.log("T", templateData.value);

  emit('save', {...templateData.value});
  Modal.getInstance(modalRef.value)?.hide()
};

</script>

<style scoped>
</style>
