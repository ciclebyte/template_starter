import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import Home from '@/views/Home.vue'
import Templates from '@/views/Templates.vue'
import TemplatesEdit from '@/views/TemplatesEdit.vue'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: Home
      },
      {
        path: '/templates',
        name: 'Templates',
        component: Templates
      }
      // 其他普通页面
    ]
  },
  {
    path: '/templates/edit/:id',
    name: 'TemplatesEdit',
    component: TemplatesEdit
  },
  {
    path: '/templates/generate/:id',
    name: 'TemplateGenerator',
    component: () => import('@/views/TemplateGenerator.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router 