<template>
  <div class="p-4">
    <div class="admin-layout d-flex p-3">
      <!-- Sidebar Menu -->
      <div class="">
        <nav class="sidebar bg-container p-4">
          <div class="mb-3 d-flex align-items-center">
            <img :src="business.logo_url" alt="Business Logo" class="logo me-2"
                 v-if="business.logo_url"/>
            <span class="h3">{{ business.name }}</span>
          </div>
          <ul class="nav flex-column">
            <li class="h5 nav-item">
              <router-link to="/employee/appointments" class="nav-link" active-class="active">Appointments</router-link>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <a @click="openSchedule" class="nav-link">Schedule</a>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <a href="#" data-bs-toggle="modal" class="nav-link" data-bs-target="#updateUserProfileModal">
                Profile
              </a>
              <ProfileSettings/>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <button @click="logout" class="nav-link btn btn-link text-danger">Logout</button>
            </li>
          </ul>
        </nav>
      </div>
      <!-- Main Content Area -->
      <main class="content flex-grow-1">
        <router-view/>
      </main>
    </div>
  </div>
</template>

<script setup>
import {useRouter} from 'vue-router';
import {useUserStore} from '@/stores/userStore';
import {computed, onMounted, toRef, watch} from 'vue';
import {useBusinessStore} from '@/stores/businessStore.js';
import BusinessSettings from '@/components/admin/BusinessSettings.vue';
import ProfileSettings from "@/components/ProfileSettings.vue";

const router = useRouter();
const userStore = useUserStore();
const businessStore = useBusinessStore();

// Create a reactive reference to the business data in the store
const business = toRef(businessStore, 'business');

const openSchedule = () => {
  router.push("/employee/schedule/" + userStore.user.employee_id);
};

// Logout function
const logout = () => {
  userStore.token = null;
  userStore.user = null;
  router.push('/auth/login');
};

// Fetch business data on component mount
onMounted(async () => {
  await businessStore.getBusiness(userStore.user.business_id);
  await businessStore.fetchServices();
});

// Automatically watch and react to changes in businessStore.business
watch(business, (newBusiness) => {
});
</script>

<style scoped>
.h5 {
  margin: 0.3em;
}

.admin-layout {
}

.sidebar {
  width: 350px;
}

.logo {
  height: 75px;
}
</style>
