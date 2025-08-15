<template>
  <div class="variable-define-panel">
    <!-- 变量列表 -->
    <div class="variable-list">
      <n-tree
        :data="treeData"
        :selected-keys="selectedKeys"
        :expanded-keys="expandedKeys"
        key-field="id"
        label-field="displayLabel"
        children-field="children"
        block-line
        draggable
        @update:selected-keys="onSelectVariable"
        @update:expanded-keys="onExpandChange"
        @drop="onNodeDrop"
      >
        <template #prefix="{ option }">
          <n-icon :class="getVariableIcon(option.type)">
            <component :is="getVariableIcon(option.type)" />
          </n-icon>
        </template>
        
        <template #suffix="{ option }">
          <div class="node-actions" @click.stop>
            <n-dropdown
              :options="getContextMenuOptions(option)"
              @select="(key) => onContextMenuSelect(key, option)"
              trigger="click"
            >
              <n-button size="tiny" quaternary>
                <template #icon>
                  <n-icon><EllipsisHorizontalOutline /></n-icon>
                </template>
              </n-button>
            </n-dropdown>
          </div>
        </template>
      </n-tree>
    </div>

    <!-- 变量详情编辑 -->
    <div class="variable-detail" v-if="selectedVariable">
      <div class="detail-header">
        <h4>{{ selectedVariable.name || '新变量' }}</h4>
        <div class="variable-type-badge" :class="selectedVariable.type">
          {{ getTypeLabel(selectedVariable.type) }}
        </div>
      </div>

      <n-form
        ref="formRef"
        :model="selectedVariable"
        :rules="formRules"
        label-placement="top"
        size="small"
      >
        <n-grid :cols="2" :x-gap="16" :y-gap="16">
          <!-- 基础信息 -->
          <n-grid-item :span="2">
            <div class="form-section">
              <h5 class="section-title">基础信息</h5>
            </div>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="变量名称" path="name">
              <n-input 
                v-model:value="selectedVariable.name"
                placeholder="例如：userName"
                @input="updateDisplayLabel"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="显示名称" path="displayName">
              <n-input 
                v-model:value="selectedVariable.displayName"
                placeholder="例如：用户名称"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="2">
            <n-form-item label="描述" path="description">
              <n-input 
                v-model:value="selectedVariable.description"
                type="textarea"
                placeholder="详细描述变量的用途和含义"
                :rows="3"
              />
            </n-form-item>
          </n-grid-item>

          <!-- 类型配置 -->
          <n-grid-item :span="2">
            <div class="form-section">
              <h5 class="section-title">类型配置</h5>
            </div>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="数据类型" path="type">
              <n-select
                v-model:value="selectedVariable.type"
                :options="typeOptions"
                @update:value="onTypeChange"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="必填项">
              <n-switch
                v-model:value="selectedVariable.required"
                size="small"
              >
                <template #checked>是</template>
                <template #unchecked>否</template>
              </n-switch>
            </n-form-item>
          </n-grid-item>

          <n-grid-item v-if="selectedVariable.type !== 'object' && selectedVariable.type !== 'array'">
            <n-form-item label="默认值" path="defaultValue">
              <component
                :is="getInputComponent(selectedVariable.type)"
                v-model:value="selectedVariable.defaultValue"
                :placeholder="getDefaultPlaceholder(selectedVariable.type)"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="示例值" path="example">
              <n-input 
                v-model:value="selectedVariable.example"
                :placeholder="getExamplePlaceholder(selectedVariable.type)"
              />
            </n-form-item>
          </n-grid-item>

          <!-- 验证规则 -->
          <n-grid-item :span="2" v-if="selectedVariable.type === 'string'">
            <div class="form-section">
              <h5 class="section-title">验证规则</h5>
            </div>
          </n-grid-item>

          <n-grid-item v-if="selectedVariable.type === 'string'">
            <n-form-item label="最小长度">
              <n-input-number
                v-model:value="selectedVariable.validation.minLength"
                :min="0"
                placeholder="最小长度"
                clearable
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item v-if="selectedVariable.type === 'string'">
            <n-form-item label="最大长度">
              <n-input-number
                v-model:value="selectedVariable.validation.maxLength"
                :min="0"
                placeholder="最大长度"
                clearable
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="2" v-if="selectedVariable.type === 'string'">
            <n-form-item label="正则表达式">
              <n-input
                v-model:value="selectedVariable.validation.pattern"
                placeholder="例如：^[a-zA-Z0-9]+$"
              />
            </n-form-item>
          </n-grid-item>

          <!-- 枚举值 -->
          <n-grid-item :span="2" v-if="selectedVariable.type === 'string' || selectedVariable.type === 'number'">
            <n-form-item label="枚举值">
              <n-dynamic-tags
                v-model:value="selectedVariable.validation.enum"
                :max="10"
                placeholder="添加可选值"
              />
            </n-form-item>
          </n-grid-item>

          <!-- 数组配置 -->
          <n-grid-item v-if="selectedVariable.type === 'array'">
            <n-form-item label="元素类型">
              <n-select
                v-model:value="selectedVariable.itemType"
                :options="typeOptions.filter(t => t.value !== 'array')"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item v-if="selectedVariable.type === 'array'">
            <n-form-item label="最小项数">
              <n-input-number
                v-model:value="selectedVariable.validation.minItems"
                :min="0"
                clearable
              />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <!-- 操作按钮 -->
        <div class="form-actions">
          <n-space>
            <n-button type="primary" size="small" @click="saveVariable">
              <template #icon>
                <n-icon><SaveOutline /></n-icon>
              </template>
              保存
            </n-button>
            
            <n-button size="small" @click="addChildVariable" v-if="selectedVariable.type === 'object'">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              添加子字段
            </n-button>
            
            <n-button size="small" type="error" @click="deleteVariable">
              <template #icon>
                <n-icon><TrashOutline /></n-icon>
              </template>
              删除
            </n-button>
          </n-space>
        </div>
      </n-form>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-else>
      <n-empty description="请选择一个变量进行编辑">
        <template #extra>
          <n-button size="small" @click="$emit('add-variable')">
            <template #icon>
              <n-icon><AddOutline /></n-icon>
            </template>
            添加新变量
          </n-button>
        </template>
      </n-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import {
  NTree, NForm, NFormItem, NGrid, NGridItem, NInput, NInputNumber,
  NSelect, NSwitch, NButton, NSpace, NIcon, NEmpty, NDynamicTags,
  NDropdown, useMessage
} from 'naive-ui'
import {
  EllipsisHorizontalOutline, SaveOutline, AddOutline, TrashOutline,
  DocumentTextOutline, CodeOutline, ToggleOutline, CalendarOutline,
  ListOutline, ArchiveOutline
} from '@vicons/ionicons5'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  },
  selectedVariable: {
    type: Object,
    default: null
  },
  templateId: {
    type: [String, Number],
    required: true
  },
  currentFile: {
    type: Object,
    default: null
  }
})

const emit = defineEmits([
  'update:variables',
  'update:selectedVariable',
  'variable-change',
  'add-variable'
])

const message = useMessage()
const formRef = ref(null)

// 选中状态
const selectedKeys = ref([])
const expandedKeys = ref([])

// 数据类型选项
const typeOptions = [
  { label: '字符串', value: 'string' },
  { label: '数字', value: 'number' },
  { label: '布尔值', value: 'boolean' },
  { label: '对象', value: 'object' },
  { label: '数组', value: 'array' },
  { label: '日期', value: 'date' },
  { label: '文件', value: 'file' }
]

// 表单验证规则
const formRules = {
  name: {
    required: true,
    message: '请输入变量名称',
    trigger: ['input', 'blur']
  },
  displayName: {
    required: true,
    message: '请输入显示名称',
    trigger: ['input', 'blur']
  }
}

// 转换为树形数据
const treeData = computed(() => {
  return props.variables.map(variable => ({
    ...variable,
    displayLabel: variable.displayName || variable.name || '未命名变量',
    children: variable.children || []
  }))
})

// 获取变量图标
const getVariableIcon = (type) => {
  const iconMap = {
    string: DocumentTextOutline,
    number: CodeOutline,
    boolean: ToggleOutline,
    date: CalendarOutline,
    array: ListOutline,
    object: ArchiveOutline,
    file: DocumentTextOutline
  }
  return iconMap[type] || DocumentTextOutline
}

// 获取类型标签
const getTypeLabel = (type) => {
  const typeMap = {
    string: '字符串',
    number: '数字',
    boolean: '布尔',
    date: '日期',
    array: '数组',
    object: '对象',
    file: '文件'
  }
  return typeMap[type] || type
}

// 获取输入组件
const getInputComponent = (type) => {
  const componentMap = {
    string: 'n-input',
    number: 'n-input-number',
    boolean: 'n-switch',
    date: 'n-date-picker'
  }
  return componentMap[type] || 'n-input'
}

// 获取默认值占位符
const getDefaultPlaceholder = (type) => {
  const placeholderMap = {
    string: '默认字符串值',
    number: '默认数字值',
    boolean: '默认布尔值',
    date: '选择默认日期'
  }
  return placeholderMap[type] || '默认值'
}

// 获取示例值占位符
const getExamplePlaceholder = (type) => {
  const placeholderMap = {
    string: 'example string',
    number: '123',
    boolean: 'true',
    date: '2024-01-01',
    array: '[item1, item2]',
    object: '{ key: value }'
  }
  return placeholderMap[type] || '示例值'
}

// 选择变量
const onSelectVariable = (keys) => {
  selectedKeys.value = keys
  if (keys.length > 0) {
    const variable = findVariableById(keys[0], props.variables)
    emit('update:selectedVariable', variable)
  }
}

// 递归查找变量
const findVariableById = (id, variables) => {
  for (const variable of variables) {
    if (variable.id === id) {
      return variable
    }
    if (variable.children) {
      const found = findVariableById(id, variable.children)
      if (found) return found
    }
  }
  return null
}

// 展开状态变化
const onExpandChange = (keys) => {
  expandedKeys.value = keys
}

// 更新显示标签
const updateDisplayLabel = () => {
  if (props.selectedVariable && !props.selectedVariable.displayName) {
    props.selectedVariable.displayLabel = props.selectedVariable.name
  }
}

// 类型变化处理
const onTypeChange = (type) => {
  if (props.selectedVariable) {
    // 初始化验证规则
    props.selectedVariable.validation = {
      minLength: null,
      maxLength: null,
      pattern: '',
      enum: [],
      minItems: null,
      maxItems: null
    }
    
    // 如果改为对象类型，初始化children数组
    if (type === 'object' && !props.selectedVariable.children) {
      props.selectedVariable.children = []
    }
    
    // 如果改为数组类型，设置默认元素类型
    if (type === 'array' && !props.selectedVariable.itemType) {
      props.selectedVariable.itemType = 'string'
    }
  }
}

// 保存变量
const saveVariable = async () => {
  try {
    await formRef.value?.validate()
    emit('variable-change', props.selectedVariable)
    message.success('变量保存成功')
  } catch (error) {
    message.error('请填写必填字段')
  }
}

// 添加子变量
const addChildVariable = () => {
  if (props.selectedVariable && props.selectedVariable.type === 'object') {
    if (!props.selectedVariable.children) {
      props.selectedVariable.children = []
    }
    
    const childVariable = {
      id: Date.now().toString(),
      name: '',
      displayName: '',
      description: '',
      type: 'string',
      required: false,
      defaultValue: '',
      example: '',
      validation: {
        minLength: null,
        maxLength: null,
        pattern: '',
        enum: []
      }
    }
    
    props.selectedVariable.children.push(childVariable)
    
    // 展开当前节点
    if (!expandedKeys.value.includes(props.selectedVariable.id)) {
      expandedKeys.value.push(props.selectedVariable.id)
    }
    
    // 选中新创建的子变量
    nextTick(() => {
      emit('update:selectedVariable', childVariable)
      selectedKeys.value = [childVariable.id]
    })
  }
}

// 删除变量
const deleteVariable = () => {
  if (props.selectedVariable) {
    // 从父级中移除
    const removeFromParent = (variables, targetId) => {
      for (let i = 0; i < variables.length; i++) {
        if (variables[i].id === targetId) {
          variables.splice(i, 1)
          return true
        }
        if (variables[i].children) {
          if (removeFromParent(variables[i].children, targetId)) {
            return true
          }
        }
      }
      return false
    }
    
    removeFromParent(props.variables, props.selectedVariable.id)
    emit('update:selectedVariable', null)
    selectedKeys.value = []
    message.success('变量删除成功')
  }
}

// 获取上下文菜单选项
const getContextMenuOptions = (option) => {
  const options = [
    {
      label: '编辑',
      key: 'edit'
    }
  ]
  
  if (option.type === 'object') {
    options.push({
      label: '添加子字段',
      key: 'addChild'
    })
  }
  
  options.push(
    {
      type: 'divider'
    },
    {
      label: '删除',
      key: 'delete'
    }
  )
  
  return options
}

// 上下文菜单选择
const onContextMenuSelect = (key, option) => {
  switch (key) {
    case 'edit':
      emit('update:selectedVariable', option)
      selectedKeys.value = [option.id]
      break
    case 'addChild':
      emit('update:selectedVariable', option)
      selectedKeys.value = [option.id]
      nextTick(() => {
        addChildVariable()
      })
      break
    case 'delete':
      emit('update:selectedVariable', option)
      selectedKeys.value = [option.id]
      nextTick(() => {
        deleteVariable()
      })
      break
  }
}

// 节点拖拽
const onNodeDrop = (info) => {
  // 实现拖拽重排序逻辑
  console.log('Node drop:', info)
}

// 监听选中变量变化
watch(() => props.selectedVariable, (newVar) => {
  if (newVar) {
    selectedKeys.value = [newVar.id]
  }
})
</script>

<style scoped>
.variable-define-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.variable-list {
  flex: 1;
  min-height: 300px;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  overflow-y: auto;
}

.variable-detail {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.detail-header h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.variable-type-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: white;
}

.variable-type-badge.string {
  background: #52c41a;
}

.variable-type-badge.number {
  background: #1890ff;
}

.variable-type-badge.boolean {
  background: #722ed1;
}

.variable-type-badge.object {
  background: #fa8c16;
}

.variable-type-badge.array {
  background: #eb2f96;
}

.variable-type-badge.date {
  background: #13c2c2;
}

.variable-type-badge.file {
  background: #666;
}

.form-section {
  margin: 16px 0 8px 0;
}

.section-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 4px;
}

.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.node-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

:deep(.n-tree-node):hover .node-actions {
  opacity: 1;
}
</style>