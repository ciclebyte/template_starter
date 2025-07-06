import axios from 'axios';
import qs from 'qs';

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

// 添加请求拦截器
service.interceptors.request.use(
	(config) => {
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// 添加响应拦截器
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
	(error) => {
		console.error('请求错误:', error);
		return Promise.reject(error);
	}
);

// 导出 axios 实例
export default service;
