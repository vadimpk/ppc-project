<template>
  <div ref="modalRef" class="modal fade" id="manageServiceModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1"
       aria-labelledby="manageServiceModalLabel" aria-hidden="true">

    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="serviceModalLabel">{{ isEditMode ? 'Edit Service' : 'Add New Service' }}</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitService">
            <div class="mb-3">
              <label for="serviceName" class="form-label">Service Name</label>
              <input type="text" v-model="serviceData.name" id="serviceName" class="form-control" required/>
            </div>
            <div class="mb-3">
              <label for="serviceDescription" class="form-label">Description</label>
              <textarea v-model="serviceData.description" id="serviceDescription" class="form-control"></textarea>
            </div>
            <div class="mb-3">
              <label for="serviceDuration" class="form-label">Duration</label>
              <input
                  type="text"
                  id="durationInput"
                  v-model="durationStr"
                  class="form-control"
                  placeholder="1h 23m or 23m"
                  pattern="^(\d+h\s?)?(\d+m)?$"
                  required
              />
            </div>
            <div class="mb-3">
              <label for="servicePrice" class="form-label">Price</label>
              <input type="number" v-model="price" id="servicePrice" class="form-control" required/>
            </div>
            <button type="submit" class="btn btn-primary w-100">{{ isEditMode ? 'Update' : 'Create' }} Service</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, watch, onMounted} from 'vue';
import {useBusinessStore} from "@/stores/businessStore.js";
import {Modal} from 'bootstrap';
import {formatDuration, parseDurationToInt} from "@/utils/covertors.js";


const props = defineProps({
  isVisible: Boolean,
  isEditMode: Boolean,
  service: {type: Object, default: () => ({})}
});

const durationStr = ref('');
const price = ref(100);

const serviceData = ref({
  id: '',
  name: '',
  description: '',
  duration: 0,
  price: 0
});

const resetForm = () => {
  serviceData.value = {
    id: '',
    name: '',
    description: '',
    duration: 0,
    price: 0
  };
};

watch(() => props.service, (newService) => {
  if (props.isEditMode && newService) {
    serviceData.value = {...newService};
    durationStr.value = formatDuration(newService.duration);
  } else {
    resetForm();
  }
}, {immediate: true});

const businessStore = useBusinessStore();

const modalRef = ref(null);

const submitService = async () => {
  serviceData.value.duration = parseDurationToInt(durationStr.value)
  serviceData.value.price = price.value * 100;
  if (props.isEditMode) {
    await businessStore.updateService(serviceData.value.id, serviceData.value);
  } else {
    await businessStore.createService(serviceData.value);
  }
  Modal.getInstance(modalRef.value)?.hide()
};
</script>

<style scoped>
.modal-body {
  padding: 2rem;
}
</style>
