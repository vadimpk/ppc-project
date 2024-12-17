<template>
  <div ref="modalRef" class="modal fade" id="updateUserProfileModal" data-bs-backdrop="static" data-bs-keyboard="false"
       tabindex="-1" aria-labelledby="updateUserProfileModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title" id="updateUserProfileModalLabel">
            <i class="bi bi-person-lines-fill me-2"></i>
            Update Profile Information
          </h4>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="updateUserProfile">
            <div class="mb-2">
              <input type="text" v-model="fullName" id="fullName" class="form-control-plaintext" required placeholder="Full Name"/>
            </div>
            <div class="mb-2">
              <input type="email" v-model="email" id="email" class="form-control-plaintext" required placeholder="Email"/>
            </div>
            <div class="mb-3">
              <input type="tel" v-model="phone" id="phone" class="form-control-plaintext" placeholder="Phone"/>
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
