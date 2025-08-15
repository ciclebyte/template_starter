import request from '@/utils/request'

/**
 * 订阅预设变量到模板
 * @param {number} templateId - 模板ID
 * @param {Array} presets - 预设变量订阅列表
 * @returns {Promise}
 */
export function subscribePreset(templateId, presets) {
  return request({
    url: `/templates/${templateId}/preset-variables/subscribe`,
    method: 'POST',
    data: {
      template_id: templateId,
      presets: presets
    }
  })
}

/**
 * 获取模板订阅的预设变量列表
 * @param {number} templateId - 模板ID
 * @returns {Promise}
 */
export function getSubscribedPresets(templateId) {
  return request({
    url: `/templates/${templateId}/preset-variables`,
    method: 'GET'
  })
}

/**
 * 取消订阅预设变量
 * @param {number} templateId - 模板ID
 * @param {number} id - 关联ID
 * @returns {Promise}
 */
export function unsubscribePreset(templateId, id) {
  return request({
    url: `/templates/${templateId}/preset-variables/${id}`,
    method: 'DELETE'
  })
}

/**
 * 更新预设变量映射
 * @param {number} templateId - 模板ID
 * @param {number} id - 关联ID
 * @param {Object} data - 更新数据
 * @returns {Promise}
 */
export function updatePresetMapping(templateId, id, data) {
  return request({
    url: `/templates/${templateId}/preset-variables/${id}`,
    method: 'PUT',
    data: {
      template_id: templateId,
      id: id,
      ...data
    }
  })
}

/**
 * 获取可用预设变量列表
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function getAvailablePresets(params = {}) {
  return request({
    url: '/templates/preset-variables/available',
    method: 'GET',
    params: {
      page: 1,
      size: 20,
      ...params
    }
  })
}