import request from '@/utils/request'

/**
 * 获取模板暴露字段
 * @param {Object} params - 参数对象
 * @param {number} params.templateId - 模板ID
 * @param {string} params.version - 版本号（可选）
 * @returns {Promise}
 */
export const getTemplateExpose = (params) => {
  const { templateId, version } = params
  const url = version 
    ? `/api/v1/templates/${templateId}/expose?version=${version}`
    : `/api/v1/templates/${templateId}/expose`
  
  return request({
    url,
    method: 'GET'
  })
}

/**
 * 设置模板暴露字段
 * @param {Object} params - 参数对象
 * @param {number} params.templateId - 模板ID
 * @param {Object} params.varsSchema - 变量Schema定义
 * @param {string} params.version - 版本号（可选）
 * @returns {Promise}
 */
export const setTemplateExpose = (params) => {
  const { templateId, varsSchema, version, description, ...otherData } = params
  
  return request({
    url: `/api/v1/templates/${templateId}/expose`,
    method: 'PUT',
    data: {
      fieldSchemaJson: JSON.stringify(varsSchema),
      version: version || '1.0',
      description: description || '',
      ...otherData
    }
  })
}

/**
 * 删除模板暴露字段
 * @param {Object} params - 参数对象  
 * @param {number} params.templateId - 模板ID
 * @param {string} params.version - 版本号（可选）
 * @returns {Promise}
 */
export const deleteTemplateExpose = (params) => {
  const { templateId, version } = params
  const url = version 
    ? `/api/v1/templates/${templateId}/expose?version=${version}`
    : `/api/v1/templates/${templateId}/expose`
    
  return request({
    url,
    method: 'DELETE'
  })
}

/**
 * 获取模板暴露字段历史版本列表
 * @param {Object} params - 参数对象
 * @param {number} params.templateId - 模板ID
 * @returns {Promise}
 */
export const getTemplateExposeVersions = (params) => {
  const { templateId } = params
  
  return request({
    url: `/api/v1/templates/${templateId}/expose/versions`,
    method: 'GET'
  })
}