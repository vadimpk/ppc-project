<template>
  <div ref="modalRef" class="modal fade" id="employeeModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1"
       aria-labelledby="employeeModalLabel" aria-hidden="true">

    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="employeeModalLabel">{{ props.employee?.user?.full_name }} - Services</h5>
          <button type="button" class="btn-close" @click="closeModal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <h6>Assigned Services</h6>
          <ul class="list-group mb-4">
            <li class="list-group-item d-flex justify-content-between align-items-center"
                v-for="service in employeeServices" :key="service.id">
              <span>{{ service.name }}</span>
              <button @click="removeService(service.id)" class="btn btn-sm btn-outline-danger">
                <i class="bi bi-trash"></i>
                Remove
              </button>
            </li>
          </ul>

          <h6>Assign New Service</h6>
          <select v-model="newServiceId" class="form-select mb-3">
            <option v-for="service in services" :key="service.id" :value="service.id">{{ service.name }}</option>
          </select>
          <button @click="assignService" class="btn btn-primary" :disabled="services.length === 0">Assign Service</button>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup>
import {Modal} from 'bootstrap';
import {computed, ref, watch} from 'vue';
import {useBusinessStore} from '@/stores/businessStore.js';

const props = defineProps(['employee']);
const newServiceId = ref(null);

const businessStore = useBusinessStore();

// Computed property for retrieving the services of the current employee
const employeeServices = computed(() => {
  console.log(businessStore.employeeServices)
  return businessStore.employeeServices || [];
});

// List of available services that are not already assigned to the employee
const services = computed(() => {
      return businessStore.services.filter(service => !businessStore.employeeServices.some(s => s.id === service.id))
    }
);

// Watch the employee prop and fetch services when it changes
watch(
    () => props.employee,
    async (newEmployee) => {
      if (newEmployee) {
        await businessStore.fetchEmployeeServices(newEmployee.id);
      }
    },
    {immediate: true}
);

// Assign a new service to the selected employee
const assignService = async () => {
  if (newServiceId.value && props.employee) {
    await businessStore.assignService(props.employee.id, newServiceId.value);
    newServiceId.value = null;
  }
};

// Remove a service from the selected employee
const removeService = async (serviceId) => {
  if (props.employee) {
    await businessStore.removeService(props.employee.id, serviceId);
  }
};

// Close the modal
const modalRef = ref(null);
const closeModal = () => {
  Modal.getInstance(modalRef.value)?.hide();
};
</script>

<style scoped>

</style>