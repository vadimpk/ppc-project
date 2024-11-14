<template>
  <div ref="modalRef" class="modal fade" id="updateBusinessModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1"
       aria-labelledby="updateBusinessModalLabel" aria-hidden="true">

    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title" id="updateBusinessModalLabel">Update Business Information</h4>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="updateBusiness">
            <div class="mb-4">
              <label for="businessName" class="form-label">Business Name</label>
              <input type="text" v-model="name" id="businessName" class="form-control" required/>
            </div>
            <div class="mb-4">
              <label for="logoURL" class="form-label">Logo URL</label>
              <input type="url" v-model="logoURL" id="logoURL" class="form-control"/>
            </div>
            <button type="submit" class="btn btn-primary w-100">Update</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, watchEffect} from 'vue';
import {useBusinessStore} from '@/stores/businessStore';
import {Modal} from 'bootstrap';

const props = defineProps(['business'])

const businessStore = useBusinessStore();

const name = ref(props.business?.name || '');
const logoURL = ref(props.business?.logo_url || '');

watchEffect(() => {
  name.value = props.business?.name || '';
  logoURL.value = props.business?.logo_url || '';
})

const modalRef = ref(null);
const updateBusiness = async () => {
  await businessStore.updateBusiness({name: name.value});
  if (logoURL.value) {
    await businessStore.updateBusinessAppearance({logo_url: logoURL.value});
  }
  Modal.getInstance(modalRef.value)?.hide();
};
</script>

<style scoped>
.modal-body {
  padding: 2rem;
}
</style>
