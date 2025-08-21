import { useAuthStore } from '@/stores/auth'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])

// 路由元信息接口定义
// meta: {
//   requiresAuth: boolean,      // 是否需要认证
//   roles: string[],           // 需要的角色
//   permissions: string[],     // 需要的权限
//   title: string             // 页面标题
// }

export function setupRouterGuards(router) {
    // 全局前置守卫
    router.beforeEach(async (to, from, next) => {
        const authStore = useAuthStore()
        
        // 设置页面标题
        if (to.meta?.title) {
            document.title = `${to.meta.title} - Template Starter`
        }
        
        // 如果用户已登录但用户信息为空，尝试获取用户信息
        if (authStore.accessToken && !authStore.user) {
            await authStore.fetchUserInfo()
        }
        
        // 检查是否需要认证
        if (to.meta?.requiresAuth && !authStore.isLoggedIn) {
            message.warning('请先登录')
            next({
                path: '/login',
                query: { redirect: to.fullPath }
            })
            return
        }
        
        // 检查角色权限
        if (to.meta?.roles && to.meta.roles.length > 0) {
            if (!authStore.hasAnyRole(to.meta.roles)) {
                message.error('权限不足')
                next('/403')
                return
            }
        }
        
        // 检查功能权限
        if (to.meta?.permissions && to.meta.permissions.length > 0) {
            if (!authStore.hasAllPermissions(to.meta.permissions)) {
                message.error('权限不足')
                next('/403')
                return
            }
        }
        
        // 已登录用户访问登录页，重定向到首页
        if (to.path === '/login' && authStore.isLoggedIn) {
            next('/')
            return
        }
        
        next()
    })
    
    // 全局后置钩子
    router.afterEach((to, from) => {
        // 可以在这里记录页面访问日志
        console.log(`导航到: ${to.path}`)
    })
}