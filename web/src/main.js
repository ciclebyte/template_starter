import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import router from './router'
import naive from 'naive-ui'
import 'vfonts/Lato.css'
import 'vfonts/FiraCode.css'
import { setupRouterGuards } from './router/guards'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(naive)

// 设置路由守卫
setupRouterGuards(router)

// 初始化认证store
const authStore = useAuthStore()
authStore.init()

app.mount('#app')
