import request from '@/utils/request'

// ============================================================================
// 管理员 API Key 管理
// ============================================================================

// 获取API Key列表
export function getApiKeys(params) {
    return request({
        url: '/api/v1/api-keys',
        method: 'get',
        params
    })
}

// 创建API Key
export function createApiKey(data) {
    return request({
        url: '/api/v1/api-keys',
        method: 'post',
        data
    })
}

// 更新API Key
export function updateApiKey(id, data) {
    return request({
        url: `/api/v1/api-keys/${id}`,
        method: 'put',
        data
    })
}

// 删除API Key
export function deleteApiKey(id) {
    return request({
        url: `/api/v1/api-keys/${id}`,
        method: 'delete'
    })
}

// 重新生成API Key Secret
export function regenerateApiKey(id) {
    return request({
        url: `/api/v1/api-keys/${id}/regenerate`,
        method: 'post'
    })
}

// 获取API Key使用日志
export function getApiKeyLogs(id, params) {
    return request({
        url: `/api/v1/api-keys/${id}/logs`,
        method: 'get',
        params
    })
}

// ============================================================================
// 个人 API Key 管理
// ============================================================================

// 获取我的API Key列表
export function getMyApiKeys(params) {
    return request({
        url: '/api/v1/profile/api-keys',
        method: 'get',
        params
    })
}

// 创建我的API Key
export function createMyApiKey(data) {
    return request({
        url: '/api/v1/profile/api-keys',
        method: 'post',
        data
    })
}

// 更新我的API Key
export function updateMyApiKey(id, data) {
    return request({
        url: `/api/v1/profile/api-keys/${id}`,
        method: 'put',
        data
    })
}

// 删除我的API Key
export function deleteMyApiKey(id) {
    return request({
        url: `/api/v1/profile/api-keys/${id}`,
        method: 'delete'
    })
}

// 重新生成我的API Key Secret
export function regenerateMyApiKey(id) {
    return request({
        url: `/api/v1/profile/api-keys/${id}/regenerate`,
        method: 'post'
    })
}