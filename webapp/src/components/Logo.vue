<template>
  <div class="d-flex justify-content-center align-items-center">
    <img src="../assets/logo.svg" alt="App Logo" class="logo me-3"/>
    <p class="m-0" :class="size">Synergo</p>
  </div>
</template>

<script setup>
import {computed, defineProps} from 'vue';

const props = defineProps({
  size: {
    type: String,
    default: 'h1'
  }
});

const size = computed(() => {
  return props.size;
});
</script>

<style scoped>
.logo {
  width: 60px;
}
</style>