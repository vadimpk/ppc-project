<template>
  <div ref="modalRef" class="modal fade" id="updateUserProfileModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1" aria-labelledby="updateUserProfileModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title" id="updateUserProfileModalLabel">Update Profile Information</h4>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="updateUserProfile">
            <div class="mb-4">
              <label for="fullName" class="form-label">Full Name</label>
              <input type="text" v-model="fullName" id="fullName" class="form-control" required/>
            </div>
            <div class="mb-4">
              <label for="email" class="form-label">Email</label>
              <input type="email" v-model="email" id="email" class="form-control" required/>
            </div>
            <div class="mb-4">
              <label for="phone" class="form-label">Phone</label>
              <input type="tel" v-model="phone" id="phone" class="form-control"/>
            </div>
            <button type="submit" class="btn btn-primary w-100">Update</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watchEffect } from 'vue';
import { Modal } from 'bootstrap';
import { useUserStore } from '@/stores/userStore';

const userStore = useUserStore();

// Reactive fields for user information
const fullName = ref(userStore.user?.full_name || '');
const email = ref(userStore.user?.email || '');
const phone = ref(userStore.user?.phone || '');

const modalRef = ref(null);

// Method to update user profile information
const updateUserProfile = async () => {
  await userStore.updateUserProfile({
    full_name: fullName.value,
    email: email.value,
    phone: phone.value,
  });
  Modal.getInstance(modalRef.value)?.hide();
};
</script>

<style scoped>
.modal-body {
  padding: 2rem;
}
</style>
