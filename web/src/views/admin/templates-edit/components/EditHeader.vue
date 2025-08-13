<template>
  <div class="edit-header">
    <div class="header-left">
      <span class="edit-title">模板编辑</span>
      <!-- 文件树切换按钮 -->
      <div class="file-tree-toggle" :class="{ 'active': isFileTreeVisible }" @click="$emit('toggle-file-tree')">
        <n-icon class="toggle-icon">
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path fill="currentColor" d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"/>
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
    </div>
    <div class="header-actions">
      <n-button size="small" @click="$emit('show-variable-manager')">
        <template #icon>
          <n-icon><Settings /></n-icon>
        </template>
        变量管理
      </n-button>
      <n-button quaternary circle class="edit-close-btn" @click="$emit('close-edit')">
        <template #icon>
          <n-icon><svg viewBox="0 0 24 24" width="20" height="20"><path fill="currentColor" d="M18.3 5.71a1 1 0 0 0-1.41 0L12 10.59 7.11 5.7A1 1 0 0 0 5.7 7.11L10.59 12l-4.89 4.89a1 1 0 1 0 1.41 1.41L12 13.41l4.89 4.89a1 1 0 0 0 1.41-1.41L13.41 12l4.89-4.89a1 1 0 0 0 0-1.4z"/></svg></n-icon>
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { NButton, NIcon } from 'naive-ui'
import { ChevronDown, Settings } from '@vicons/ionicons5'

defineProps({
  isVariablePanelOpen: {
    type: Boolean,
    required: true
  },
  isFileTreeVisible: {
    type: Boolean,
    required: true
  }
})

defineEmits(['toggle-variable-panel', 'show-variable-manager', 'close-edit', 'toggle-file-tree'])
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
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.edit-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
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