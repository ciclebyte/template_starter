import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const routes = [
  // 前台路由
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'home',
        component: () => import('@/views/client/home/index.vue')
      },
      {
        path: '/templates',
        name: 'templates',
        component: () => import('@/views/client/templates-public/index.vue')
      },
    ]
  },
  // 独立的全屏页面（不在布局内）
  {
    path: '/templates/edit/:id',
    name: 'TemplatesEdit',
    component: () => import('@/views/admin/templates-edit/index.vue')
  },
  {
    path: '/template-generator',
    name: 'generator',
    component: () => import('@/views/client/template-generator/index.vue')
  },
  {
    path: '/template-generator/:id',
    name: 'TemplateGenerator',
    component: () => import('@/views/client/template-generator/index.vue')
  },

  // 后台管理路由
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'admin-dashboard',
        component: () => import('@/views/admin/dashboard/index.vue')
      },
      {
        path: 'templates',
        name: 'admin-templates-list',
        component: () => import('@/views/admin/template/index.vue')
      },
      {
        path: 'categories',
        name: 'admin-categories',
        component: () => import('@/views/admin/category/index.vue')
      },
      {
        path: 'languages',
        name: 'admin-languages',
        component: () => import('@/views/admin/language/index.vue')
      },
      {
        path: 'analytics',
        name: 'admin-analytics',
        component: () => import('@/views/admin/Analytics/index.vue')
      },
      {
        path: 'settings',
        name: 'admin-settings',
        component: () => import('@/views/admin/setting/index.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router 