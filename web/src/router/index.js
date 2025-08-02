import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'
import Home from '@/views/Home.vue'
import Templates from '@/views/Templates.vue'
import TemplatesEdit from '@/views/TemplatesEdit.vue'

const routes = [
  // 前台路由
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'home',
        component: Home
      },
      {
        path: '/templates',
        name: 'templates',
        component: () => import('@/views/TemplatesPublic.vue')
      },
      {
        path: '/template-generator',
        name: 'generator',
        component: () => import('@/views/TemplateGenerator.vue')
      },
      {
        path: '/template-generator/:id',
        name: 'TemplateGenerator',
        component: () => import('@/views/TemplateGenerator.vue')
      }
    ]
  },
  // 独立的模板编辑页面（不在布局内）
  {
    path: '/templates/edit/:id',
    name: 'TemplatesEdit',
    component: TemplatesEdit
  },

  // 后台管理路由
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'admin-dashboard',
        component: () => import('@/views/admin/Dashboard.vue')
      },
      {
        path: 'templates',
        name: 'admin-templates-list',
        component: () => import('@/views/admin/TemplatesList.vue')
      },
      {
        path: 'categories',
        name: 'admin-categories',
        component: () => import('@/views/admin/CategoriesManage.vue')
      },
      {
        path: 'languages',
        name: 'admin-languages',
        component: () => import('@/views/admin/LanguagesManage.vue')
      },
      {
        path: 'files',
        name: 'admin-files',
        component: () => import('@/views/Templates.vue') // 临时使用现有组件
      },
      {
        path: 'analytics',
        name: 'admin-analytics',
        component: () => import('@/views/Templates.vue') // 临时使用现有组件
      },
      {
        path: 'settings',
        name: 'admin-settings',
        component: () => import('@/views/Templates.vue') // 临时使用现有组件
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router 