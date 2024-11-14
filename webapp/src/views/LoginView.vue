<template>
  <div class="vh-100 d-flex align-items-center justify-content-center">
    <div class="container p-5 bg-container">
      <h2 class="mb-4 h2 text-center">Start your journey of project management!</h2>

      <!--      <div class="d-flex justify-content-center mb-4">-->
      <!--        <button-->
      <!--            type="button"-->
      <!--            class="btn me-2"-->
      <!--            :class="accountType === 'user' ? 'btn-primary' : 'btn-outline-primary'"-->
      <!--            @click="setAccountType('user')">-->
      <!--          Client-->
      <!--        </button>-->
      <!--        <button-->
      <!--            type="button"-->
      <!--            class="btn"-->
      <!--            :class="accountType === 'business' ? 'btn-primary' : 'btn-outline-primary'"-->
      <!--            @click="setAccountType('business')">-->
      <!--          Business-->
      <!--        </button>-->
      <!--      </div>-->

      <form @submit.prevent="handleSubmit" class="mx-auto" style="max-width: 400px;">
        <!--        <div class="mb-3" v-if="accountType === 'business'">-->
        <!--          <label for="businessID" class="form-label">Business ID</label>-->
        <!--          <input type="number" v-model="formData.business_id" id="businessID" class="form-control" :required="accountType === 'business'" />-->
        <!--        </div>-->

        <div class="mb-3">
          <label for="emailOrPhone" class="form-label">Email or Phone</label>
          <input type="text" v-model="formData.email_or_phone" id="emailOrPhone" class="form-control" required/>
        </div>

        <div class="mb-5">
          <label for="password" class="form-label">Password</label>
          <input type="password" v-model="formData.password" id="password" class="form-control" required/>
        </div>

        <button type="submit" class="btn btn-lg btn-success w-100">Login</button>
      </form>

      <p class="text-center mt-3">
        Don't have an account?
        <router-link to="/auth/register" class="text-decoration-none">Register here</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue';
import {useUserStore} from '@/stores/userStore';
import {useRouter} from "vue-router";
import {USER_ROLE_ADMIN, USER_ROLE_CLIENT, USER_ROLE_EMPLOYEE} from "@/utils/constants.js";

const userStore = useUserStore();

const formData = ref({
  email_or_phone: '',
  password: '',
  business_id: 4
});

const accountType = ref('user');

const router = useRouter();

const handleSubmit = async () => {
      const payload = {password: formData.value.password, business_id: 1};
      payload.business_id = formData.value.business_id;
      if (formData.value.email_or_phone.includes('@')) {
        payload.email = formData.value.email_or_phone;
      } else {
        payload.phone = formData.value.email_or_phone;
      }

      userStore.loginUser(payload).then(() => {
        if (userStore.user.role === USER_ROLE_ADMIN) {
          router.push('/admin/services');
        } else if (userStore.user.role === USER_ROLE_EMPLOYEE) {
          router.push(`/employee/schedule/${userStore.user.employee_id}`);
        } else if (userStore.user.role === USER_ROLE_CLIENT)
          router.push('/client/dashboard');
      });
    }
;
</script>

<style scoped>
.container {
  max-width: 600px;
  width: 500px;
  margin: auto;
}
</style>
