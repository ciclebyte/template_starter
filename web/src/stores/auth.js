import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { login, register, logout, getUserInfo, refreshToken } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
    // 状态
    const user = ref(null)
    const accessToken = ref(localStorage.getItem('access_token') || '')
    const refreshTokenValue = ref(localStorage.getItem('refresh_token') || '')
    const permissions = ref([])
    const roles = ref([])
    
    // 计算属性
    const isLoggedIn = computed(() => !!user.value && !!accessToken.value)
    const isAdmin = computed(() => roles.value.includes('super_admin') || roles.value.includes('system_admin'))
    const userDisplayName = computed(() => user.value?.nickname || user.value?.username || '未知用户')
    
    // 方法
    const setTokens = (tokens) => {
        accessToken.value = tokens.access_token
        refreshTokenValue.value = tokens.refresh_token
        localStorage.setItem('access_token', tokens.access_token)
        localStorage.setItem('refresh_token', tokens.refresh_token)
    }
    
    const clearAuth = () => {
        user.value = null
        accessToken.value = ''
        refreshTokenValue.value = ''
        permissions.value = []
        roles.value = []
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
    }
    
    const setUserInfo = (userInfo) => {
        user.value = userInfo
        permissions.value = userInfo.permissions || []
        roles.value = userInfo.roles || []
    }
    
    // 登录
    const doLogin = async (credentials) => {
        try {
            const response = await login(credentials)
            if (response.data.code === 0) {
                const { user: userInfo, ...tokens } = response.data.data
                setTokens(tokens)
                setUserInfo(userInfo)
                return { success: true }
            } else {
                return { success: false, message: response.data.message }
            }
        } catch (error) {
            return { success: false, message: error.message || '登录失败' }
        }
    }
    
    // 注册
    const doRegister = async (registerData) => {
        try {
            const response = await register(registerData)
            if (response.data.code === 0) {
                const { user: userInfo, ...tokens } = response.data.data
                setTokens(tokens)
                setUserInfo(userInfo)
                return { success: true }
            } else {
                return { success: false, message: response.data.message }
            }
        } catch (error) {
            return { success: false, message: error.message || '注册失败' }
        }
    }
    
    // 登出
    const doLogout = async () => {
        try {
            await logout()
        } catch (error) {
            console.error('登出失败:', error)
        } finally {
            clearAuth()
        }
    }
    
    // 获取用户信息
    const fetchUserInfo = async () => {
        if (!accessToken.value) {
            return false
        }
        
        try {
            const response = await getUserInfo()
            if (response.data.code === 0) {
                setUserInfo(response.data.data.user)
                return true
            } else {
                console.warn('获取用户信息失败:', response.data.message)
                clearAuth()
                return false
            }
        } catch (error) {
            console.error('获取用户信息失败:', error)
            // 只有在401错误时才清除认证信息，其他错误保留token以便重试
            if (error.response?.status === 401) {
                clearAuth()
            }
            return false
        }
    }
    
    // 刷新令牌
    const doRefreshToken = async () => {
        try {
            const response = await refreshToken({ refresh_token: refreshTokenValue.value })
            if (response.data.code === 0) {
                setTokens(response.data.data)
                return true
            } else {
                clearAuth()
                return false
            }
        } catch (error) {
            console.error('刷新令牌失败:', error)
            clearAuth()
            return false
        }
    }
    
    // 权限检查
    const hasPermission = (permission) => {
        if (!permission) return true
        if (isAdmin.value) return true // 管理员拥有所有权限
        return permissions.value.includes(permission)
    }
    
    const hasRole = (role) => {
        if (!role) return true
        return roles.value.includes(role)
    }
    
    const hasAnyRole = (roleList) => {
        if (!roleList || roleList.length === 0) return true
        return roleList.some(role => roles.value.includes(role))
    }
    
    const hasAllPermissions = (permissionList) => {
        if (!permissionList || permissionList.length === 0) return true
        if (isAdmin.value) return true
        return permissionList.every(permission => permissions.value.includes(permission))
    }
    
    // 初始化
    const init = async () => {
        if (accessToken.value && !user.value) {
            const success = await fetchUserInfo()
            if (!success) {
                // 如果获取用户信息失败，尝试刷新token
                if (refreshTokenValue.value) {
                    const refreshSuccess = await doRefreshToken()
                    if (refreshSuccess) {
                        await fetchUserInfo()
                    }
                }
            }
        }
    }
    
    return {
        // 状态
        user: readonly(user),
        accessToken: readonly(accessToken),
        refreshTokenValue: readonly(refreshTokenValue),
        permissions: readonly(permissions),
        roles: readonly(roles),
        
        // 计算属性
        isLoggedIn,
        isAdmin,
        userDisplayName,
        
        // 方法
        doLogin,
        doRegister,
        doLogout,
        fetchUserInfo,
        doRefreshToken,
        hasPermission,
        hasRole,
        hasAnyRole,
        hasAllPermissions,
        init,
        clearAuth
    }
})