import request from '@/utils/request'

// 模板变量-新增
export function addTemplateVariable(data) {
  return request({
    url: '/api/v1/templateVariables/add',
    method: 'post',
    data
  })
}

// 模板变量-批量删除
export function batchDeleteTemplateVariable(params) {
  return request({
    url: '/api/v1/templateVariables/batchdel',
    method: 'delete',
    params
  })
}

// 模板变量-删除
export function deleteTemplateVariable(params) {
  return request({
    url: '/api/v1/templateVariables/del',
    method: 'delete',
    params
  })
}

// 模板变量-详情
export function getTemplateVariableDetail(params) {
  return request({
    url: '/api/v1/templateVariables/detail',
    method: 'get',
    params
  })
}

// 模板变量-修改
export function editTemplateVariable(data) {
  return request({
    url: '/api/v1/templateVariables/edit',
    method: 'put',
    data
  })
}

// 模板变量-列表
export function listTemplateVariables(params) {
  return request({
    url: '/api/v1/templateVariables/list',
    method: 'get',
    params
  })
}
