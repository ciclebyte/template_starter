import request from '@/utils/request'

// 用户登录
export function login(data) {
    return request({
        url: '/v1/auth/login',
        method: 'post',
        data
    })
}

// 用户注册
export function register(data) {
    return request({
        url: '/v1/auth/register',
        method: 'post',
        data
    })
}

// 用户登出
export function logout() {
    return request({
        url: '/v1/auth/logout',
        method: 'post'
    })
}

// 获取用户信息
export function getUserInfo() {
    return request({
        url: '/v1/auth/me',
        method: 'get'
    })
}

// 刷新Token
export function refreshToken(data) {
    return request({
        url: '/v1/auth/refresh',
        method: 'post',
        data
    })
}

// 检查权限
export function checkPermission(data) {
    return request({
        url: '/v1/auth/check-permission',
        method: 'post',
        data
    })
}

// 检查角色
export function checkRole(data) {
    return request({
        url: '/v1/auth/check-role',
        method: 'post',
        data
    })
}