<template>
  <div class="p-4">
    <div class="admin-layout d-flex p-3">
      <!-- Sidebar Menu -->
      <div class="">
        <nav class="sidebar bg-container p-4">
          <div class="d-flex align-items-center mb-5 ms-4">
            <img src="../assets/logo.svg" alt="App Logo" class="logo me-3"/>
            <p class="h2 m-0">Synergo</p>
          </div>
          <ul class="nav flex-column">
            <li class="h5 nav-item">
              <router-link to="/client/dashboard" class="nav-link" active-class="active">
                <i class="bi bi-house"></i>
                <span class="ms-2">Dashboard</span>
              </router-link>
            </li>
            <hr class="hr"/>
            <li class="h5 nav-item">
              <router-link to="/client/appointments" class="nav-link" active-class="active">
                <i class="bi bi-calendar-check me-2"></i>
                Appointments
              </router-link>
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
import Logo from "@/components/Logo.vue";

const router = useRouter();
const userStore = useUserStore();
const businessStore = useBusinessStore();

const logout = () => {
  userStore.token = null;
  userStore.user = null;
  router.push('/auth/login');
};

onMounted(async () => {
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
