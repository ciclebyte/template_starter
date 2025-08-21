<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h1>注册 Template Starter</h1>
        <p>创建您的新账户</p>
      </div>
      
      <n-form
        ref="formRef"
        :model="registerForm"
        :rules="rules"
        size="large"
        label-placement="top"
      >
        <n-form-item label="用户名" path="username">
          <n-input
            v-model:value="registerForm.username"
            placeholder="请输入用户名 (3-20位)"
            :disabled="loading"
          />
        </n-form-item>
        
        <n-form-item label="邮箱" path="email">
          <n-input
            v-model:value="registerForm.email"
            placeholder="请输入邮箱地址"
            :disabled="loading"
          />
        </n-form-item>
        
        <n-form-item label="昵称" path="nickname">
          <n-input
            v-model:value="registerForm.nickname"
            placeholder="请输入昵称 (可选)"
            :disabled="loading"
          />
        </n-form-item>
        
        <n-form-item label="密码" path="password">
          <n-input
            v-model:value="registerForm.password"
            type="password"
            placeholder="请输入密码 (6-20位)"
            :disabled="loading"
            show-password-on="click"
          />
        </n-form-item>
        
        <n-form-item label="确认密码" path="confirmPassword">
          <n-input
            v-model:value="registerForm.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            :disabled="loading"
            show-password-on="click"
            @keydown.enter="handleRegister"
          />
        </n-form-item>
        
        <n-form-item>
          <n-button
            type="primary"
            size="large"
            block
            :loading="loading"
            @click="handleRegister"
          >
            注册
          </n-button>
        </n-form-item>
      </n-form>
      
      <div class="register-footer">
        <p>
          已有账户？
          <router-link to="/login" class="link">立即登录</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMessage } from 'naive-ui'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const loading = ref(false)
const formRef = ref(null)

const registerForm = reactive({
  username: '',
  email: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

const validatePasswordSame = (rule, value) => {
  return value === registerForm.password || new Error('两次输入的密码不一致')
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20位', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  nickname: [
    { min: 2, max: 20, message: '昵称长度为2-20位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validatePasswordSame, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true
    
    const { confirmPassword, ...registerData } = registerForm
    const result = await authStore.doRegister(registerData)
    
    if (result.success) {
      message.success('注册成功，已自动登录')
      router.push('/')
    } else {
      message.error(result.message || '注册失败')
    }
  } catch (error) {
    console.error('注册错误:', error)
    message.error('注册失败，请检查输入信息')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // 如果已经登录，直接跳转到首页
  if (authStore.isLoggedIn) {
    router.push('/')
  }
})
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-header h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.register-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.register-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #666;
}

.link {
  color: #18a058;
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .register-card {
    padding: 24px;
  }
  
  .register-header h1 {
    font-size: 20px;
  }
}
</style>