<template>
  <div class="edit-header">
    <div class="header-left">
      <div class="title-area">
        <span class="edit-title">模板编辑</span>
      </div>
      <!-- 文件树切换按钮 -->
      <div class="file-tree-toggle" :class="{ 'active': isFileTreeVisible }" @click="$emit('toggle-file-tree')">
        <n-icon class="toggle-icon">
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path fill="currentColor" d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z" />
          </svg>
        </n-icon>
        <span class="toggle-text">文件树</span>
      </div>
      <!-- 模板变量展开按钮 -->
      <div class="variable-expand-trigger" @click="$emit('toggle-variable-panel')">
        <span class="trigger-text">模板变量</span>
        <n-icon class="trigger-icon" :class="{ 'rotated': isVariablePanelOpen }">
          <ChevronDown />
        </n-icon>
      </div>
      <div v-if="currentFileName" class="file-status">
        <span class="file-name">{{ currentFileName }}</span>
        <span v-if="hasUnsavedChanges" class="unsaved-indicator" title="有未保存的更改">●</span>
      </div>
    </div>
    <div class="header-actions">
      <n-button size="small" @click="$emit('show-variable-expose')">
        <template #icon>
          <n-icon>
            <CodeSlash />
          </n-icon>
        </template>
        变量定义
      </n-button>
      <n-button size="small" @click="$emit('show-variable-manager')">
        <template #icon>
          <n-icon>
            <Settings />
          </n-icon>
        </template>
        变量管理
      </n-button>
      <n-button size="small" @click="$emit('show-settings')" title="编辑器设置">
        <template #icon>
          <n-icon>
            <svg viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor"
                d="M19.14,12.94c0.04-0.3,0.06-0.61,0.06-0.94c0-0.32-0.02-0.64-0.07-0.94l2.03-1.58c0.18-0.14,0.23-0.41,0.12-0.61 l-1.92-3.32c-0.12-0.22-0.37-0.29-0.59-0.22l-2.39,0.96c-0.5-0.38-1.03-0.7-1.62-0.94L14.4,2.81c-0.04-0.24-0.24-0.41-0.48-0.41 h-3.84c-0.24,0-0.43,0.17-0.47,0.41L9.25,5.35C8.66,5.59,8.12,5.92,7.63,6.29L5.24,5.33c-0.22-0.08-0.47,0-0.59,0.22L2.74,8.87 C2.62,9.08,2.66,9.34,2.86,9.48l2.03,1.58C4.84,11.36,4.82,11.69,4.82,12s0.02,0.64,0.07,0.94l-2.03,1.58 c-0.18,0.14-0.23,0.41-0.12,0.61l1.92,3.32c0.12,0.22,0.37,0.29,0.59,0.22l2.39-0.96c0.5,0.38,1.03,0.7,1.62,0.94l0.36,2.54 c0.05,0.24,0.24,0.41,0.48,0.41h3.84c0.24,0,0.44-0.17,0.47-0.41l0.36-2.54c0.59-0.24,1.13-0.56,1.62-0.94l2.39,0.96 c0.22,0.08,0.47,0,0.59-0.22l1.92-3.32c0.12-0.22,0.07-0.47-0.12-0.61L19.14,12.94z M12,15.6c-1.98,0-3.6-1.62-3.6-3.6 s1.62-3.6,3.6-3.6s3.6,1.62,3.6,3.6S13.98,15.6,12,15.6z" />
            </svg>
          </n-icon>
        </template>
        设置
      </n-button>
      <n-button quaternary circle class="edit-close-btn" @click="$emit('close-edit')">
        <template #icon>
          <n-icon><svg viewBox="0 0 24 24" width="20" height="20">
              <path fill="currentColor"
                d="M18.3 5.71a1 1 0 0 0-1.41 0L12 10.59 7.11 5.7A1 1 0 0 0 5.7 7.11L10.59 12l-4.89 4.89a1 1 0 1 0 1.41 1.41L12 13.41l4.89 4.89a1 1 0 0 0 1.41-1.41L13.41 12l4.89-4.89a1 1 0 0 0 0-1.4z" />
            </svg></n-icon>
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { NButton, NIcon } from 'naive-ui'
import { ChevronDown, Settings, CodeSlash } from '@vicons/ionicons5'

defineProps({
  isVariablePanelOpen: {
    type: Boolean,
    required: true
  },
  isFileTreeVisible: {
    type: Boolean,
    required: true
  },
  hasUnsavedChanges: {
    type: Boolean,
    default: false
  },
  currentFileName: {
    type: String,
    default: ''
  }
})

defineEmits(['toggle-variable-panel', 'show-variable-manager', 'show-variable-expose', 'close-edit', 'toggle-file-tree', 'show-settings'])
</script>

<style scoped>
.edit-header {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-area {
  display: flex;
  align-items: center;
  gap: 16px;
}

.edit-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
}

.file-status {
  display: flex;
  align-items: center;
  gap: 6px;
}

.file-name {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.unsaved-indicator {
  color: #ff9500;
  font-size: 20px;
  line-height: 1;
  animation: pulse 2s infinite;
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.5;
  }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 变量展开按钮样式 */
.variable-expand-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.variable-expand-trigger:hover {
  background: #e9ecef;
  border-color: #18a058;
}

.trigger-text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.trigger-icon {
  font-size: 16px;
  color: #666;
  transition: transform 0.2s;
}

.trigger-icon.rotated {
  transform: rotate(180deg);
}

/* 文件树切换按钮样式 */
.file-tree-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.file-tree-toggle:hover {
  background: #e9ecef;
  border-color: #18a058;
}

.file-tree-toggle.active {
  background: #e8f5e8;
  border-color: #18a058;
}

.file-tree-toggle.active .toggle-icon {
  color: #18a058;
}

.file-tree-toggle.active .toggle-text {
  color: #18a058;
}

.toggle-icon {
  font-size: 16px;
  color: #666;
}

.toggle-text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
</style>