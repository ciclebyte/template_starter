import request from '@/utils/request'

// 获取个人资料
export function getProfile() {
    return request({
        url: '/api/v1/profile',
        method: 'get'
    })
}

// 更新个人资料
export function updateProfile(data) {
    return request({
        url: '/api/v1/profile',
        method: 'put',
        data
    })
}

// 修改密码
export function changePassword(data) {
    return request({
        url: '/api/v1/profile/password',
        method: 'put',
        data
    })
}

// 更新邮箱
export function updateEmail(data) {
    return request({
        url: '/api/v1/profile/email',
        method: 'put',
        data
    })
}

// 上传头像
export function uploadAvatar(data) {
    return request({
        url: '/api/v1/profile/avatar',
        method: 'post',
        data
    })
}

// 获取账户安全信息
export function getSecurityInfo() {
    return request({
        url: '/api/v1/profile/security',
        method: 'get'
    })
}

// 获取登录历史
export function getLoginHistory(params) {
    return request({
        url: '/api/v1/profile/login-history',
        method: 'get',
        params
    })
}