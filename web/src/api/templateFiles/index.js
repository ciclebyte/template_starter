import request from '@/utils/request'

// 新增模板文件
export function addTemplateFile(data) {
  return request({
    url: '/api/v1/templateFiles/add',
    method: 'post',
    data
  })
}

// 编辑模板文件
export function editTemplateFile(data) {
  return request({
    url: '/api/v1/templateFiles/edit',
    method: 'put',
    data
  })
}

// 重命名模板文件
export function renameTemplateFile(data) {
  return request({
    url: '/api/v1/templateFiles/rename',
    method: 'put',
    data
  })
}

// 删除模板文件
export function delTemplateFile(id) {
  return request({
    url: '/api/v1/templateFiles/del',
    method: 'delete',
    params: { id }
  })
}

// 批量删除模板文件
export function batchDelTemplateFile(idList) {
  return request({
    url: '/api/v1/templateFiles/batchdel',
    method: 'delete',
    params: { id: idList }
  })
}

// 获取模板文件详情
export function getTemplateFileDetail(id) {
  return request({
    url: '/api/v1/templateFiles/detail',
    method: 'get',
    params: { id }
  })
}

// 获取模板文件列表
export function listTemplateFiles(params) {
  return request({
    url: '/api/v1/templateFiles/list',
    method: 'get',
    params
  })
}

// 获取模板文件树
export function getTemplateFileTree(templateId) {
  return request({
    url: '/api/v1/templateFiles/fileTree',
    method: 'get',
    params: { templateId }
  })
} 

// 获取模板文件树
export function getTemplateFileContent(id) {
  return request({
    url: '/api/v1/templateFiles/content',
    method: 'get',
    params: { id }
  })
} 

// 上传ZIP包
export function uploadZipFile(templateId, file) {
  const formData = new FormData()
  formData.append('zipFile', file)
  formData.append('templateId', templateId)
  
  return request({
    url: '/api/v1/templateFiles/uploadZip',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传代码文件
export function uploadCodeFile(templateId, file, parentId) {
  const formData = new FormData()
  formData.append('codeFile', file)
  formData.append('templateId', templateId)
  if (parentId) formData.append('parentId', parentId)
  return request({
    url: '/api/v1/templateFiles/uploadCode',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传代码文件
export function uploadCode(data) {
  return request({
    url: '/templateFiles/uploadCode',
    method: 'post',
    data
  })
}

// 渲染模板文件
export function renderTemplate(data) {
  return request({
    url: '/api/v1/templateFiles/render',
    method: 'post',
    data
  })
}