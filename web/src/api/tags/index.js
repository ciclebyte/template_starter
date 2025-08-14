import request from '@/utils/request'

// 标签-新增
export function addTag(data) {
  return request({
    url: '/api/v1/tags/add',
    method: 'post',
    data
  })
}

// 标签-批量删除
export function batchDeleteTags(data) {
  return request({
    url: '/api/v1/tags/batchdel',
    method: 'delete',
    data
  })
}

// 标签-删除
export function deleteTag(params) {
  return request({
    url: '/api/v1/tags/del',
    method: 'delete',
    params
  })
}

// 标签-详情
export function getTagDetail(params) {
  return request({
    url: '/api/v1/tags/detail',
    method: 'get',
    params
  })
}

// 标签-修改
export function editTag(data) {
  return request({
    url: '/api/v1/tags/edit',
    method: 'put',
    data
  })
}

// 标签-列表（分页）
export function listTags(params) {
  return request({
    url: '/api/v1/tags/list',
    method: 'get',
    params
  })
}

// 标签-全部（不分页）
export function getAllTags() {
  return request({
    url: '/api/v1/tags/all',
    method: 'get'
  })
}

// 标签-带模板数量统计
export function getTagsWithCount() {
  return request({
    url: '/api/v1/tags/with-count',
    method: 'get'
  })
}

// 模板标签-添加
export function addTemplateTags(templateId, data) {
  return request({
    url: `/api/v1/templates/${templateId}/tags/add`,
    method: 'post',
    data
  })
}

// 模板标签-移除
export function removeTemplateTags(templateId, data) {
  return request({
    url: `/api/v1/templates/${templateId}/tags/remove`,
    method: 'delete',
    data
  })
}

// 模板标签-设置（覆盖）
export function setTemplateTags(templateId, data) {
  return request({
    url: `/api/v1/templates/${templateId}/tags/set`,
    method: 'put',
    data
  })
}

// 模板标签-获取列表
export function getTemplateTags(templateId) {
  return request({
    url: `/api/v1/templates/${templateId}/tags`,
    method: 'get'
  })
}

// 标签模板-根据标签获取模板列表
export function getTemplatesByTag(tagId, params) {
  return request({
    url: `/api/v1/tags/${tagId}/templates`,
    method: 'get',
    params
  })
}