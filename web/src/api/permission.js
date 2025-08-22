import request from '@/utils/request'

// ============================================================================
// 权限管理
// ============================================================================

// 获取权限列表
export function getPermissions(params) {
  return request({
    url: '/api/v1/permissions',
    method: 'get',
    params
  })
}

// 创建权限
export function createPermission(data) {
  return request({
    url: '/api/v1/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(id, data) {
  return request({
    url: `/api/v1/permissions/${id}`,
    method: 'put',
    data
  })
}

// 删除权限
export function deletePermission(id) {
  return request({
    url: `/api/v1/permissions/${id}`,
    method: 'delete'
  })
}

// 获取权限详情
export function getPermission(id) {
  return request({
    url: `/api/v1/permissions/${id}`,
    method: 'get'
  })
}

// ============================================================================
// 角色管理
// ============================================================================

// 获取角色列表
export function getRoles(params) {
  return request({
    url: '/api/v1/roles',
    method: 'get',
    params
  })
}

// 创建角色
export function createRole(data) {
  return request({
    url: '/api/v1/roles',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(id, data) {
  return request({
    url: `/api/v1/roles/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(id) {
  return request({
    url: `/api/v1/roles/${id}`,
    method: 'delete'
  })
}

// 获取角色详情
export function getRole(id) {
  return request({
    url: `/api/v1/roles/${id}`,
    method: 'get'
  })
}

// 分配角色权限
export function assignRolePermissions(id, data) {
  return request({
    url: `/api/v1/roles/${id}/permissions`,
    method: 'post',
    data
  })
}

// ============================================================================
// 用户角色管理
// ============================================================================

// 获取用户角色列表
export function getUserRoles(userId) {
  return request({
    url: `/api/v1/users/${userId}/roles`,
    method: 'get'
  })
}

// 分配用户角色
export function assignUserRoles(userId, data) {
  return request({
    url: `/api/v1/users/${userId}/roles`,
    method: 'post',
    data
  })
}

// 移除用户角色
export function removeUserRole(userId, roleId) {
  return request({
    url: `/api/v1/users/${userId}/roles/${roleId}`,
    method: 'delete'
  })
}