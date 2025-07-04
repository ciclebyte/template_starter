import request from '@/utils/request'

// 模板语言-新增
export function addTemplateLanguage(data) {
  return request({
    url: '/api/v1/templateLanguages/add',
    method: 'post',
    data
  })
}

// 模板语言-批量删除
export function batchDeleteTemplateLanguage(params) {
  return request({
    url: '/api/v1/templateLanguages/batchdel',
    method: 'delete',
    params
  })
}

// 模板语言-删除
export function deleteTemplateLanguage(params) {
  return request({
    url: '/api/v1/templateLanguages/del',
    method: 'delete',
    params
  })
}

// 模板语言-详情
export function getTemplateLanguageDetail(params) {
  return request({
    url: '/api/v1/templateLanguages/detail',
    method: 'get',
    params
  })
}

// 模板语言-修改
export function editTemplateLanguage(data) {
  return request({
    url: '/api/v1/templateLanguages/edit',
    method: 'put',
    data
  })
}

// 模板语言-列表
export function listTemplateLanguages(params) {
  return request({
    url: '/api/v1/templateLanguages/list',
    method: 'get',
    params
  })
}
