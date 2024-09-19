import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index.js'

const app = createApp(App)
app.use(router)

import axios from 'axios'
app.config.globalProperties.$axios = axios.create({
    baseURL: 'http://127.0.0.1:51245',
})

import Notifications from '@kyvg/vue3-notification'
app.use(Notifications)

app.mount('#app')
