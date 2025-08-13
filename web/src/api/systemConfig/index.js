import request from '@/utils/request'

// 获取配置列表
export function getConfigList(params) {
  return request({
    url: '/api/v1/systemConfig/list',
    method: 'get',
    params
  })
}

// 获取配置详情
export function getConfigDetail(id) {
  return request({
    url: '/api/v1/systemConfig/detail',
    method: 'get',
    params: { id }
  })
}

// 新增配置
export function addConfig(data) {
  return request({
    url: '/api/v1/systemConfig/add',
    method: 'post',
    data
  })
}

// 编辑配置
export function editConfig(data) {
  return request({
    url: '/api/v1/systemConfig/edit',
    method: 'put',
    data
  })
}

// 删除配置
export function deleteConfig(id) {
  return request({
    url: '/api/v1/systemConfig/del',
    method: 'delete',
    data: { id }
  })
}

// 批量删除配置
export function batchDeleteConfig(ids) {
  return request({
    url: '/api/v1/systemConfig/batchdel',
    method: 'delete',
    data: { ids }
  })
}

// 根据键名获取配置
export function getConfigByKey(configKey) {
  return request({
    url: '/api/v1/systemConfig/getByKey',
    method: 'get',
    params: { configKey }
  })
}

// 根据分组获取配置
export function getConfigByGroup(configGroup, isPublic) {
  return request({
    url: '/api/v1/systemConfig/getByGroup',
    method: 'get',
    params: { configGroup, isPublic }
  })
}

// 批量更新配置
export function batchUpdateConfig(configs) {
  return request({
    url: '/api/v1/systemConfig/batchUpdate',
    method: 'put',
    data: { configs }
  })
}

// 获取公开配置
export function getPublicConfig(configGroup) {
  return request({
    url: '/api/v1/systemConfig/public',
    method: 'get',
    params: { configGroup }
  })
}

// 重置配置
export function resetConfig(configKey) {
  return request({
    url: '/api/v1/systemConfig/reset',
    method: 'put',
    data: { configKey }
  })
}

// 验证配置值
export function validateConfig(configKey, configValue) {
  return request({
    url: '/api/v1/systemConfig/validate',
    method: 'post',
    data: { configKey, configValue }
  })
}

// 获取AI配置
export function getAIConfig() {
  return getConfigByGroup('ai', false)
}

// 获取系统配置
export function getSystemConfig() {
  return getConfigByGroup('system', true)
}

// 获取模板配置
export function getTemplateConfig() {
  return getConfigByGroup('template', false)
}

// 获取UI配置
export function getUIConfig() {
  return getConfigByGroup('ui', true)
}

// 保存AI配置
export function saveAIConfig(configs) {
  const configArray = Object.entries(configs).map(([key, value]) => ({
    configKey: key,
    configValue: typeof value === 'object' ? JSON.stringify(value) : String(value)
  }))
  return batchUpdateConfig(configArray)
}