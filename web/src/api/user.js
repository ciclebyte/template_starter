import request from '@/utils/request'

// 获取用户列表
export function getUsers(params) {
    return request({
        url: '/api/v1/users',
        method: 'get',
        params
    })
}

// 创建用户
export function createUser(data) {
    return request({
        url: '/api/v1/users',
        method: 'post',
        data
    })
}

// 更新用户
export function updateUser(id, data) {
    return request({
        url: `/api/v1/users/${id}`,
        method: 'put',
        data: { id, ...data }
    })
}

// 删除用户
export function deleteUser(id) {
    return request({
        url: `/api/v1/users/${id}`,
        method: 'delete'
    })
}

// 获取用户详情
export function getUser(id) {
    return request({
        url: `/api/v1/users/${id}`,
        method: 'get'
    })
}

// 重置用户密码
export function resetUserPassword(id, data) {
    return request({
        url: `/api/v1/users/${id}/reset-password`,
        method: 'post',
        data: { id, ...data }
    })
}

// 更新用户状态
export function updateUserStatus(id, status) {
    return request({
        url: `/api/v1/users/${id}/status`,
        method: 'put',
        data: { id, status }
    })
}

// 分配用户角色
export function assignUserRoles(id, data) {
    return request({
        url: `/api/v1/users/${id}/roles`,
        method: 'put',
        data: { id, ...data }
    })
}