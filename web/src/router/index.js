import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const routes = [
  // 认证路由 (独立页面，不在布局内)
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: '注册' }
  },
  
  // 错误页面
  {
    path: '/401',
    name: 'error-401',
    component: () => import('@/views/errors/401.vue'),
    meta: { title: '未授权' }
  },
  {
    path: '/403',
    name: 'error-403',
    component: () => import('@/views/errors/403.vue'),
    meta: { title: '权限不足' }
  },
  
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
      {
        path: '/permission-demo',
        name: 'permission-demo',
        component: () => import('@/views/demo/PermissionDemo.vue'),
        meta: { title: '权限演示' }
      },
      {
        path: '/profile',
        name: 'profile',
        component: () => import('@/views/client/Profile.vue'),
        meta: { title: '个人资料', requiresAuth: true }
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
    path: '/admin/var-presets/:id/design',
    name: 'var-preset-design',
    component: () => import('@/views/admin/var-preset/design.vue')
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
    meta: { 
      requiresAuth: true, 
      roles: ['super_admin', 'system_admin', 'org_admin'],
      title: '管理后台'
    },
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
        path: 'tags',
        name: 'admin-tags',
        component: () => import('@/views/admin/tag/index.vue')
      },
      {
        path: 'var-presets',
        name: 'admin-var-presets',
        component: () => import('@/views/admin/var-preset/index.vue')
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
      },
      {
        path: 'permissions',
        name: 'admin-permissions',
        component: () => import('@/views/admin/permission/index.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router 