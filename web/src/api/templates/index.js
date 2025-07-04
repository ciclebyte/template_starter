import request from '@/utils/request'

// 模板-新增
export function addTemplate(data) {
  return request({
    url: '/api/v1/templates/add',
    method: 'post',
    data
  })
}

// 模板-批量删除
export function batchDeleteTemplate(params) {
  return request({
    url: '/api/v1/templates/batchdel',
    method: 'delete',
    params
  })
}

// 模板-删除
export function deleteTemplate(params) {
  return request({
    url: '/api/v1/templates/del',
    method: 'delete',
    params
  })
}

// 模板-详情
export function getTemplateDetail(params) {
  return request({
    url: '/api/v1/templates/detail',
    method: 'get',
    params
  })
}

// 模板-修改
export function editTemplate(data) {
  return request({
    url: '/api/v1/templates/edit',
    method: 'put',
    data
  })
}

// 模板-列表
export function listTemplates(params) {
  return request({
    url: '/api/v1/templates/list',
    method: 'get',
    params
  })
}
