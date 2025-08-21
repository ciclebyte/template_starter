import axios from 'axios';
import qs from 'qs';
import { useAuthStore } from '@/stores/auth'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])

function getApiBaseUrl() {
    // 1. 优先使用显式配置的环境变量
    if (import.meta.env.VITE_API_URL) {
      return import.meta.env.VITE_API_URL
    }
    
    // 生产环境使用当前域名
    const { protocol, hostname, port } = window.location
    let basePort = port ? `:${port}` : ''
    
    // 处理非标准端口
    if ((protocol === 'http:' && port === '80') || 
        (protocol === 'https:' && port === '443')) {
      basePort = ''
    }
    
    return `${protocol}//${hostname}${basePort}/api`
  }



// 配置新建一个 axios 实例
const service = axios.create({
	baseURL: getApiBaseUrl(),  // 修改为根路径，避免重复的/api
	timeout: 50000,
	headers: { 'Content-Type': 'application/json' },
	paramsSerializer: {
		serialize(params) {
			return qs.stringify(params, { allowDots: true, arrayFormat: 'brackets' });
		},
	},
});

// 添加请求拦截器 - 添加认证头
service.interceptors.request.use(
	(config) => {
		const authStore = useAuthStore()
		if (authStore.accessToken) {
			config.headers.Authorization = `Bearer ${authStore.accessToken}`
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// 添加响应拦截器 - 处理认证错误
service.interceptors.response.use(
	(response) => {
		// 如果是下载文件，直接返回response
		if (response.config.responseType === 'blob') {
			return response;
		}
		
		const res = response.data;
		if (res.code !== 0) {
			console.error('API错误:', res.message);
			return Promise.reject(new Error(res.message || '未知错误'));
		}
		
		return response;
	},
	async (error) => {
		const authStore = useAuthStore()
		
		if (error.response?.status === 401) {
			// Token过期，尝试刷新
			if (authStore.refreshTokenValue && !error.config._retry) {
				error.config._retry = true
				const refreshSuccess = await authStore.doRefreshToken()
				
				if (refreshSuccess) {
					// 重新发送原请求
					error.config.headers.Authorization = `Bearer ${authStore.accessToken}`
					return service(error.config)
				}
			}
			
			// 刷新失败或没有刷新令牌，跳转到登录页
			authStore.clearAuth()
			if (window.location.pathname !== '/login') {
				message.error('登录已过期，请重新登录')
				window.location.href = '/login'
			}
		} else if (error.response?.status === 403) {
			message.error('权限不足')
		} else if (error.response?.status >= 500) {
			message.error('服务器错误，请稍后重试')
		}
		
		console.error('请求错误:', error);
		return Promise.reject(error);
	}
);

// 导出 axios 实例
export default service;
