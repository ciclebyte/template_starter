import request from '@/utils/request'

// 获取Sprig函数列表
export function getSprigFunctions() {
  return request({
    url: '/api/v1/sprig-functions',
    method: 'get'
  })
}