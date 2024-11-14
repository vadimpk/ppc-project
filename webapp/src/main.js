import './assets/main.css'

import {createApp} from 'vue'
import {createPinia} from 'pinia'

import piniaPluginPersistence from 'pinia-plugin-persistedstate'

import App from './App.vue'
import router from './router'

import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import './assets/app.css'

const app = createApp(App)


const pinia = createPinia()
pinia.use(piniaPluginPersistence)

app.use(pinia)
app.use(router)
app.use(Toast, {
    transition: "Vue-Toastification__bounce",
    maxToasts: 20,
    newestOnTop: true
});

app.mount('#app')