<template>
  <div class="vh-100 d-flex align-items-center justify-content-center flex-column">
    <div class="mb-4"><Logo/></div>
    <div class="container p-5 bg-container">
      <h2 class="mb-4 h2 text-center" v-if="!business_id">Start your journey of online booking system!</h2>
      <h2 class="mb-4 h2 text-center" v-else>Join to the business!</h2>

      <div class="d-flex justify-content-center mb-4" v-if="!business_id">
        <button
            type="button"
            class="btn me-2"
            :class="accountType === 'user' ? 'btn-primary' : 'btn-outline-primary'"
            @click="setAccountType('user')">
          Client
        </button>
        <button
            type="button"
            class="btn"
            :class="accountType === 'business' ? 'btn-primary' : ' btn-outline-primary'"
            @click="setAccountType('business')">
          Business
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="mx-auto" style="max-width: 400px;">
        <div class="mb-4" v-if="accountType === 'business'">
          <input type="text" v-model="formData.business_name" id="businessName" class="form-control form-control-lg"
                 :required="accountType === 'business'" placeholder="Business Name"/>
        </div>

        <div class="mb-4">
          <input type="text" v-model="formData.full_name" id="fullName" class="form-control form-control-lg" required placeholder="Full Name"/>
        </div>

        <div class="mb-4">
          <input type="email" v-model="formData.email" id="email" class="form-control form-control-lg" required placeholder="Email"/>
        </div>

        <div class="mb-4">
          <input type="tel" v-model="formData.phone" id="phone" class="form-control form-control-lg" placeholder="Phone"/>
        </div>

        <div class="mb-4">
          <input type="password" v-model="formData.password" id="password" class="form-control form-control-lg" required placeholder="Password"/>
        </div>

        <button type="submit" class="btn btn-lg btn-primary w-100">Register</button>
      </form>

      <p class="text-center mt-4" v-if="!business_id">
        Already have an account?
        <router-link to="/auth/login" class="text-decoration-none">Log in</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import {computed, ref} from 'vue';
import {useToast} from 'vue-toastification';
import {useUserStore} from '@/stores/userStore';
import {USER_ROLE_ADMIN, USER_ROLE_CLIENT, USER_ROLE_EMPLOYEE} from "@/utils/constants.js";
import {useRoute, useRouter} from "vue-router";
import Logo from "@/components/Logo.vue";

const userStore = useUserStore();
const router = useRouter();

const formData = ref({
  full_name: '',
  email: '',
  phone: '',
  password: '',
  business_name: '',
  business_id: 0,
});

const route = useRoute()

const business_id = computed(() => route.query.business_id);

const accountType = ref('user');
const setAccountType = (type) => {
  accountType.value = type;
  if (type === 'user') formData.value.business_name = '';
};

const handleSubmit = async () => {
  formData.value.business_id = parseInt(route.query.business_id);
  const payload = {...formData.value};
  if (accountType.value === 'user') delete payload.business_name;

  userStore.registerUser(payload).then(() => {
    console.log(userStore.user);
    if (userStore.user.role === USER_ROLE_ADMIN) {
      router.push('/admin/services');
    } else if (userStore.user.role === USER_ROLE_EMPLOYEE) {
      router.push(`/employee/schedule/${userStore.user.employee_id}`);
    } else if (userStore.user.role === USER_ROLE_CLIENT)
      router.push('/client/dashboard');
  });
};
</script>

<style scoped>
.container {
  max-width: 600px;
  width: 550px;
}
</style>
