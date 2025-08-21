<template>
  <div class="error-page">
    <div class="error-content">
      <div class="error-code">401</div>
      <div class="error-title">未授权</div>
      <div class="error-description">
        抱歉，您需要登录才能访问此页面。
      </div>
      <div class="error-actions">
        <n-button type="primary" @click="goToLogin">
          去登录
        </n-button>
        <n-button @click="goHome">
          返回首页
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goToLogin = () => {
  router.push('/login')
}

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
}

.error-content {
  text-align: center;
  padding: 40px;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  color: #ff6b6b;
  margin-bottom: 20px;
  line-height: 1;
}

.error-title {
  font-size: 32px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.error-description {
  font-size: 16px;
  color: #666;
  margin-bottom: 32px;
  max-width: 400px;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

@media (max-width: 768px) {
  .error-code {
    font-size: 80px;
  }
  
  .error-title {
    font-size: 24px;
  }
  
  .error-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .error-actions .n-button {
    width: 200px;
  }
}
</style>