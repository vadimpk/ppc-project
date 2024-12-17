<template>
  <div class="p-4">
    <div class="admin-layout d-flex p-3">
      <!-- Sidebar Menu -->
      <div class="">
        <nav class="sidebar bg-container p-4">
          <div class="mb-4 ms-3 d-flex align-items-center">
            <img :src="business.logo_url" alt="Business Logo" class="logo me-3"
                 v-if="business.logo_url"/>
            <span class="h3">{{ business.name }}</span>
          </div>
          <ul class="nav flex-column">
            <li class="h5 nav-item">
              <router-link to="/admin/appointments" class="nav-link" active-class="active">
                <i class="bi bi-calendar-check me-2"></i>
                Appointments
              </router-link>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <router-link to="/admin/services" class="nav-link" active-class="active">
                <i class="bi bi-card-checklist me-2"></i>
                Services
              </router-link>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <router-link to="/admin/employees" class="nav-link" active-class="active">
                <i class="bi bi-people me-2"></i>
                Employees
              </router-link>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <a href="#" data-bs-toggle="modal" class="nav-link" data-bs-target="#updateBusinessModal">
                <i class="bi bi-gear me-2"></i>
                Settings
              </a>
              <BusinessSettings :business="business"/>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <a href="#" data-bs-toggle="modal" class="nav-link" data-bs-target="#updateUserProfileModal">
                <i class="bi bi-person-lines-fill me-2"></i>
                Profile
              </a>
              <ProfileSettings/>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <button @click="logout" class="nav-link btn btn-link text-danger">
                <i class="bi bi-box-arrow-left me-2"></i>
                Logout
              </button>
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
import {onMounted, toRef, watch} from 'vue';
import {useBusinessStore} from '@/stores/businessStore.js';
import BusinessSettings from '@/components/admin/BusinessSettings.vue';
import ProfileSettings from "@/components/ProfileSettings.vue";

const router = useRouter();
const userStore = useUserStore();
const businessStore = useBusinessStore();

// Create a reactive reference to the business data in the store
const business = toRef(businessStore, 'business');

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
  height: 50px;
}
</style>
