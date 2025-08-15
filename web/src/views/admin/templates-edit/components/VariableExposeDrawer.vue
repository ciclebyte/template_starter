<template>
  <n-drawer v-model:show="visible" :width="'90%'" placement="bottom" :height="drawerHeight + 'px'">
    <n-drawer-content title="模板变量定义" :native-scrollbar="false">
      <template #header-extra>
        <n-space>
          <n-button size="small" @click="saveVariables" type="primary" :loading="saving">
            <template #icon>
              <n-icon><SaveOutline /></n-icon>
            </template>
            保存定义
          </n-button>
          <n-button size="small" quaternary @click="visible = false">
            <template #icon>
              <n-icon><CloseOutline /></n-icon>
            </template>
          </n-button>
        </n-space>
      </template>

      <!-- 拖拽手柄 -->
      <div class="resize-handle" @mousedown="startResize">
        <div class="handle-bar"></div>
      </div>

      <div class="expose-layout">
        <!-- 左侧：变量树 -->
        <div class="left-panel" :style="{ width: leftPanelWidth + 'px' }">
          <div class="panel-title">
            变量资源
          </div>
          <div class="variable-tree" @contextmenu="onTreeAreaContextMenu">
            <n-tree
              :data="variableTreeData"
              :selected-keys="selectedKeys"
              :expanded-keys="expandedKeys"
              key-field="key"
              label-field="title"
              children-field="children"
              :node-props="nodeProps"
              block-line
              @update:selected-keys="onSelectVariable"
              @update:expanded-keys="onExpandKeys"
            >
              
              <template #suffix="{ option }">
                <div class="node-actions" @click.stop>
                  <n-dropdown
                    :options="getNodeMenuOptions(option)"
                    @select="(key) => handleNodeAction(key, option)"
                    trigger="click"
                    placement="bottom-end"
                  >
                    <n-button size="tiny" quaternary circle>
                      <template #icon>
                        <n-icon><EllipsisHorizontalOutline /></n-icon>
                      </template>
                    </n-button>
                  </n-dropdown>
                </div>
              </template>
            </n-tree>
            
            <!-- 空状态提示 -->
            <div v-if="!variableTreeData || variableTreeData.length === 0" 
                 class="empty-tree-hint"
                 @contextmenu="onTreeAreaContextMenu">
              暂无变量（右键添加）
            </div>
            
            <!-- 右键上下文菜单 -->
            <n-dropdown
              to="body"
              trigger="manual"
              :x="contextMenuX"
              :y="contextMenuY"
              :options="contextMenuOptions"
              :show="showContextMenuFlag"
              @select="handleContextMenuAction"
              @clickoutside="hideContextMenu"
              placement="bottom-start"
            />
          </div>
        </div>

        <!-- 左右分割线 -->
        <div class="panel-resizer left-resizer" @mousedown="startLeftResize"></div>

        <!-- 中间：变量信息 -->
        <div class="center-panel" :style="{ width: centerPanelWidth + 'px' }">
          <div class="panel-title">变量信息</div>
          
          <!-- 变量编辑表单 -->
          <div v-if="selectedVariableData" class="variable-form">
            <div class="form-header">
              <span class="var-path">{{ selectedVariableData.path }}</span>
              <n-tag :type="getTypeTagType(selectedVariableData.type)" size="small">
                {{ getTypeLabel(selectedVariableData.type) }}
              </n-tag>
            </div>
            
            <n-form ref="formRef" :model="selectedVariableData" size="small" label-placement="top">
              <n-grid :cols="2" :x-gap="16" :y-gap="12">
                <!-- 基础信息 -->
                <n-grid-item :span="2">
                  <div class="form-section-title">基础信息</div>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="显示名称 (title)">
                    <n-input v-model:value="selectedVariableData.title" placeholder="用户友好的显示名称" />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="数据类型 (type)">
                    <n-select 
                      v-model:value="selectedVariableData.type" 
                      :options="typeOptions"
                      @update:value="onTypeChange"
                    />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item :span="2">
                  <n-form-item label="描述信息 (description)">
                    <n-input 
                      v-model:value="selectedVariableData.description" 
                      type="textarea" 
                      :rows="2"
                      placeholder="一句话说明变量的用途"
                    />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="是否必填 (required)">
                    <n-switch v-model:value="selectedVariableData.required">
                      <template #checked>必填</template>
                      <template #unchecked>可选</template>
                    </n-switch>
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="默认值 (default)">
                    <component
                      :is="getInputComponent(selectedVariableData.type)"
                      v-model:value="selectedVariableData.default"
                      :placeholder="`默认${getTypeLabel(selectedVariableData.type)}值`"
                      :disabled="selectedVariableData.type === 'object' || selectedVariableData.type === 'array'"
                    />
                  </n-form-item>
                </n-grid-item>

                <!-- 特殊类型配置 -->
                <template v-if="selectedVariableData.type === 'enum'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">枚举配置</div>
                  </n-grid-item>
                  <n-grid-item :span="2">
                    <n-form-item label="可选值 (enum)">
                      <n-dynamic-tags 
                        v-model:value="selectedVariableData.enum"
                        placeholder="添加可选值"
                      />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <template v-if="selectedVariableData.type === 'array'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">数组配置</div>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="元素类型 (items.type)">
                      <n-select 
                        v-model:value="selectedVariableData.items.type"
                        :options="basicTypeOptions"
                      />
                    </n-form-item>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最小项数 (minItems)">
                      <n-input-number v-model:value="selectedVariableData.minItems" :min="0" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- 数值验证 -->
                <template v-if="selectedVariableData.type === 'integer' || selectedVariableData.type === 'number'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">数值验证</div>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最小值 (min)">
                      <n-input-number v-model:value="selectedVariableData.min" />
                    </n-form-item>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最大值 (max)">
                      <n-input-number v-model:value="selectedVariableData.max" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- 字符串验证 -->
                <template v-if="selectedVariableData.type === 'string'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">字符串验证</div>
                  </n-grid-item>
                  <n-grid-item :span="2">
                    <n-form-item label="正则表达式 (pattern)">
                      <n-input v-model:value="selectedVariableData.pattern" placeholder="^[a-zA-Z0-9]+$" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- UI配置 -->
                <n-grid-item :span="2">
                  <div class="form-section-title">UI配置</div>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="面板显示 (ui.panel)">
                    <n-switch v-model:value="selectedVariableData.ui.panel">
                      <template #checked>显示</template>
                      <template #unchecked>隐藏</template>
                    </n-switch>
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="排序权重 (ui.order)">
                    <n-input-number v-model:value="selectedVariableData.ui.order" :step="10" />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="分组标题 (ui.group)">
                    <n-input v-model:value="selectedVariableData.ui.group" placeholder="基础信息" />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="组件类型 (ui.component)">
                    <n-select 
                      v-model:value="selectedVariableData.ui.component"
                      :options="componentOptions"
                    />
                  </n-form-item>
                </n-grid-item>

                <!-- 命名策略 -->
                <n-grid-item :span="2">
                  <div class="form-section-title">高级配置</div>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="命名策略 (naming_policy)">
                    <n-select 
                      v-model:value="selectedVariableData.naming_policy"
                      :options="namingPolicyOptions"
                    />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </n-form>
          </div>

          <!-- 空状态 -->
          <div v-else class="empty-selection">
            <n-empty description="请从左侧选择一个变量进行编辑">
              <template #icon>
                <n-icon><SettingsOutline /></n-icon>
              </template>
            </n-empty>
          </div>
        </div>

        <!-- 中右分割线 -->
        <div class="panel-resizer right-resizer" @mousedown="startRightResize"></div>

        <!-- 右侧：预览 -->
        <div class="right-panel" :style="{ width: rightPanelWidth + 'px' }">
          <div class="panel-title">
            预览
            <n-button-group size="small">
              <n-button 
                :type="previewFormat === 'json' ? 'primary' : 'default'"
                @click="previewFormat = 'json'"
              >
                JSON
              </n-button>
              <n-button 
                :type="previewFormat === 'yaml' ? 'primary' : 'default'"
                @click="previewFormat = 'yaml'"
              >
                YAML
              </n-button>
            </n-button-group>
          </div>
          <div class="preview-content">
            <div ref="schemaEditorRef" class="code-preview"></div>
          </div>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup>
import { ref, computed, watch, h, nextTick, onMounted } from 'vue'
import { 
  NDrawer, NDrawerContent, NButton, NIcon, NSpace, NForm, NFormItem, 
  NGrid, NGridItem, NInput, NSelect, NSwitch, NTabs, NTabPane, NInputNumber,
  NTree, NTag, NDropdown, NDynamicTags, NEmpty, NButtonGroup, useMessage 
} from 'naive-ui'
import { 
  SaveOutline, CloseOutline, DocumentOutline, AddOutline, TrashOutline,
  EllipsisHorizontalOutline, SettingsOutline, CodeSlashOutline,
  TextOutline, Calculator, ToggleOutline, ListOutline, ArchiveOutline,
  KeyOutline, LockClosedOutline, Folder, FolderOpenOutline, 
  EllipsisVerticalOutline, CheckboxOutline, CodeOutline
} from '@vicons/ionicons5'
import request from '@/utils/request'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { json } from '@codemirror/lang-json'
import { yaml } from '@codemirror/lang-yaml'
import * as YAML from 'js-yaml'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  templateId: {
    type: [String, Number],
    required: true
  },
  templateVariables: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:show'])

const message = useMessage()

// 抽屉显示状态
const visible = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value)
})

// 数据状态
const saving = ref(false)
const varsSchema = ref({}) // 变量Schema数据
const selectedKeys = ref([])
const expandedKeys = ref([])
const selectedVariableData = ref(null)
const previewFormat = ref('json')
const schemaEditorRef = ref(null)
let schemaEditor = null

// 右键菜单状态
const showContextMenuFlag = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const contextMenuOptions = ref([])
const contextMenuTarget = ref(null)

// 拖拽相关状态
const drawerHeight = ref(Math.floor(window.innerHeight * 0.67)) // 默认三分之二高度
const isResizing = ref(false)
const minHeight = 300
const maxHeight = window.innerHeight - 100

// 面板宽度状态
const leftPanelWidth = ref(250)
const centerPanelWidth = ref(0) // 将通过computed计算
const rightPanelWidth = ref(0) // 将通过computed计算
const isLeftResizing = ref(false)
const isRightResizing = ref(false)
const totalLayoutWidth = ref(1200) // 默认总宽度

// 计算中间和右侧面板宽度
const computedCenterWidth = computed(() => {
  const remaining = totalLayoutWidth.value - leftPanelWidth.value - rightPanelWidth.value - 32 - 16 // 减去padding和间距
  return Math.max(300, remaining) // 最小宽度300px
})

const computedRightWidth = computed(() => {
  return Math.max(300, Math.floor(totalLayoutWidth.value / 3)) // 默认1/3宽度，最小300px
})

// 监听计算值变化
watch(computedCenterWidth, (newWidth) => {
  centerPanelWidth.value = newWidth
})

watch(computedRightWidth, (newWidth) => {
  rightPanelWidth.value = newWidth
}, { immediate: true })

// 节点属性配置（参考TemplateFileTree）
const nodeProps = ({ option }) => {
  return {
    onContextmenu(e) {
      e.preventDefault()
      e.stopPropagation()
      
      // 设置下拉菜单选项
      contextMenuOptions.value = getNodeMenuOptions(option)
      contextMenuTarget.value = option
      contextMenuX.value = e.clientX
      contextMenuY.value = e.clientY
      showContextMenuFlag.value = true
    }
  }
}

// 隐藏右键菜单
const hideContextMenu = () => {
  showContextMenuFlag.value = false
  contextMenuTarget.value = null
}

// 处理右键菜单操作
const handleContextMenuAction = (key) => {
  if (contextMenuTarget.value) {
    handleNodeAction(key, contextMenuTarget.value)
  } else {
    // 空白区域的右键菜单操作
    handleEmptyAreaAction(key)
  }
  hideContextMenu()
}

// 空白区域右键菜单处理
const handleEmptyAreaAction = (key) => {
  switch (key) {
    case 'add-string':
      addRootVariable() // 直接使用现有的添加根变量函数
      break
  }
}


// 树区域右键菜单处理
const onTreeAreaContextMenu = (event) => {
  // 检查是否点击在树节点上
  if (event.target.closest('.n-tree-node')) return
  
  event.preventDefault()
  event.stopPropagation()
  
  // 设置空白区域的右键菜单选项
  contextMenuOptions.value = [
    {
      label: '添加变量',
      key: 'add-string',
      icon: () => h(NIcon, null, { default: () => h(AddOutline) })
    }
  ]
  
  contextMenuTarget.value = null // 标记为空白区域
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  showContextMenuFlag.value = true
  
  console.log('Right-click on empty area at', event.clientX, event.clientY)
}

// 数据类型选项
const typeOptions = [
  { label: '字符串 (string)', value: 'string' },
  { label: '整数 (integer)', value: 'integer' },
  { label: '数字 (number)', value: 'number' },
  { label: '布尔值 (boolean)', value: 'boolean' },
  { label: '数组 (array)', value: 'array' },
  { label: '对象 (object)', value: 'object' },
  { label: '枚举 (enum)', value: 'enum' },
  { label: '密码 (secret)', value: 'secret' }
]

const basicTypeOptions = [
  { label: '字符串', value: 'string' },
  { label: '整数', value: 'integer' },
  { label: '数字', value: 'number' },
  { label: '布尔值', value: 'boolean' }
]

const componentOptions = [
  { label: '输入框 (input)', value: 'input' },
  { label: '下拉选择 (select)', value: 'select' },
  { label: '开关 (switch)', value: 'switch' },
  { label: '文本域 (textarea)', value: 'textarea' }
]

const namingPolicyOptions = [
  { label: 'Go蛇形 (go_snake)', value: 'go_snake' },
  { label: 'TS驼峰 (ts_camel)', value: 'ts_camel' },
  { label: '短横线 (kebab)', value: 'kebab' }
]

// 转换为变量树数据
const variableTreeData = computed(() => {
  const treeData = convertToTreeData(varsSchema.value)
  console.log('Generated tree data:', treeData) // 调试日志
  return treeData
})

// 转换变量Schema为树形数据
const convertToTreeData = (schema, parentPath = '') => {
  const treeData = []
  
  if (!schema || typeof schema !== 'object') {
    return treeData
  }
  
  Object.entries(schema).forEach(([key, value]) => {
    const currentPath = parentPath ? `${parentPath}.${key}` : key
    const varType = value.type || 'string'
    
    // 检查对象类型是否有子变量
    const hasChildren = varType === 'object' && value.properties && Object.keys(value.properties).length > 0
    
    const node = {
      key: currentPath,
      title: value.title || key,
      type: varType,
      path: currentPath,
      data: value,
      isLeaf: !hasChildren,
      // 添加prefix函数，参考模板资源树的实现
      prefix: () => h(NIcon, { class: `var-icon var-${varType}` }, {
        default: () => h(getVariableIconComponent(varType, hasChildren))
      })
    }
    
    // 如果是对象类型且有properties，递归生成子节点
    if (hasChildren) {
      node.children = convertToTreeData(value.properties, currentPath)
    }
    
    treeData.push(node)
  })
  
  return treeData
}

// 获取变量图标组件（参考模板资源树的实现）
const getVariableIconComponent = (type, hasChildren = false) => {
  // 对象类型根据是否有子变量显示不同图标
  if (type === 'object') {
    return hasChildren ? FolderOpenOutline : Folder
  }
  
  // 其他类型使用标准图标
  const iconMap = {
    string: TextOutline,        // 📝 文本图标
    integer: Calculator,        // 🔢 计算器图标
    number: Calculator,         // 🔢 计算器图标 
    boolean: CheckboxOutline,   // ☑️ 复选框图标
    array: ListOutline,         // 📋 列表图标
    enum: EllipsisVerticalOutline, // ⋮ 选择图标
    secret: LockClosedOutline   // 🔒 锁图标
  }
  return iconMap[type] || TextOutline
}


// 获取类型标签样式
const getTypeTagType = (type) => {
  const tagMap = {
    string: 'info',
    integer: 'success',
    number: 'success', 
    boolean: 'warning',
    array: 'error',
    object: 'default',
    enum: 'info',
    secret: 'error'
  }
  return tagMap[type] || 'default'
}

// 获取类型标签
const getTypeLabel = (type) => {
  const typeMap = {
    string: '字符串',
    integer: '整数',
    number: '数字',
    boolean: '布尔值',
    array: '数组',
    object: '对象',
    enum: '枚举',
    secret: '密码'
  }
  return typeMap[type] || type
}

// 获取输入组件
const getInputComponent = (type) => {
  const componentMap = {
    string: 'n-input',
    integer: 'n-input-number',
    number: 'n-input-number',
    boolean: 'n-switch',
    secret: 'n-input'
  }
  return componentMap[type] || 'n-input'
}

// 获取变量模板字符串
const getVariableTemplate = (varName) => {
  return `{{.${varName || 'varName'}}}`
}

// 创建默认变量配置
const createDefaultVariable = (name = '', type = 'string') => {
  return {
    type,
    title: name,
    description: '',
    required: false,
    default: getDefaultValue(type),
    enum: type === 'enum' ? [] : undefined,
    items: type === 'array' ? { type: 'string' } : undefined,
    properties: type === 'object' ? {} : undefined,
    min: ['integer', 'number'].includes(type) ? undefined : undefined,
    max: ['integer', 'number'].includes(type) ? undefined : undefined,
    minItems: type === 'array' ? undefined : undefined,
    pattern: type === 'string' ? undefined : undefined,
    ui: {
      panel: true,
      order: 10,
      group: '基础信息',
      component: getDefaultComponent(type)
    },
    naming_policy: 'go_snake',
    conditional: undefined
  }
}

// 获取类型默认值
const getDefaultValue = (type) => {
  const defaults = {
    string: '',
    integer: 0,
    number: 0.0,
    boolean: false,
    array: [],
    object: {},
    enum: '',
    secret: ''
  }
  return defaults[type]
}

// 获取类型默认组件
const getDefaultComponent = (type) => {
  const components = {
    string: 'input',
    integer: 'input',
    number: 'input',
    boolean: 'switch',
    array: 'input',
    object: 'input',
    enum: 'select',
    secret: 'input'
  }
  return components[type] || 'input'
}

// 树选择处理
const onSelectVariable = (selectedKeys) => {
  if (selectedKeys.length > 0) {
    const selectedKey = selectedKeys[0]
    const pathParts = selectedKey.split('.')
    let current = varsSchema.value
    
    for (const part of pathParts) {
      if (current[part]) {
        current = current[part]
      }
    }
    
    // 创建编辑数据的响应式副本
    selectedVariableData.value = {
      path: selectedKey,
      ...JSON.parse(JSON.stringify(current)),
      // 确保必要的嵌套对象存在
      ui: current.ui || {
        panel: true,
        order: 10,
        group: '基础信息',
        component: 'input'
      },
      items: current.items || (current.type === 'array' ? { type: 'string' } : undefined),
      enum: current.enum || (current.type === 'enum' ? [] : undefined)
    }
  } else {
    selectedVariableData.value = null
  }
}

// 展开键处理
const onExpandKeys = (keys) => {
  expandedKeys.value = keys
}

// 添加根变量
const addRootVariable = () => {
  const varName = `variable_${Date.now()}`
  varsSchema.value[varName] = createDefaultVariable(varName, 'string')
  
  // 自动选中新创建的变量
  selectedKeys.value = [varName]
  onSelectVariable([varName])
}


// 获取节点菜单选项
const getNodeMenuOptions = (option) => {
  const menuOptions = []
  
  // 对象类型可以添加子属性
  if (option.type === 'object') {
    menuOptions.push({
      label: '新增子变量',
      key: 'add-child',
      icon: () => h(NIcon, null, { default: () => h(AddOutline) })
    })
  }
  
  // 复制变量
  menuOptions.push({
    label: '复制变量',
    key: 'copy',
    icon: () => h(NIcon, null, { default: () => h(DocumentOutline) })
  })
  
  // 删除变量
  menuOptions.push({
    label: '删除变量',
    key: 'delete',
    icon: () => h(NIcon, null, { default: () => h(TrashOutline) })
  })
  
  return menuOptions
}

// 处理节点操作
const handleNodeAction = (key, option) => {
  switch (key) {
    case 'add-child':
      addChildVariable(option.path)
      break
    case 'copy':
      copyVariable(option.path)
      break
    case 'delete':
      deleteVariable(option.path)
      break
  }
}

// 添加子变量
const addChildVariable = (parentPath) => {
  console.log('Adding child to:', parentPath) // 调试日志
  
  const pathParts = parentPath.split('.')
  let current = varsSchema.value
  
  // 导航到父级变量
  for (const part of pathParts) {
    if (current[part]) {
      current = current[part]
    } else {
      console.error('Path not found:', part, 'in', current)
      return
    }
  }
  
  // 确保父级是对象类型并有properties
  if (current.type !== 'object') {
    current.type = 'object'
  }
  
  if (!current.properties) {
    current.properties = {}
  }
  
  const childName = `child_${Date.now()}`
  current.properties[childName] = createDefaultVariable(childName, 'string')
  
  console.log('Added child:', childName, 'to', parentPath) // 调试日志
  
  // 展开父级并选中新创建的子变量
  const childPath = `${parentPath}.${childName}`
  if (!expandedKeys.value.includes(parentPath)) {
    expandedKeys.value = [...expandedKeys.value, parentPath]
  }
  selectedKeys.value = [childPath]
  onSelectVariable([childPath])
  
  // 强制更新树数据
  varsSchema.value = { ...varsSchema.value }
}


// 复制变量
const copyVariable = (path) => {
  console.log('Copying variable:', path)
  
  const pathParts = path.split('.')
  let current = varsSchema.value
  
  // 导航到目标变量
  for (let i = 0; i < pathParts.length - 1; i++) {
    if (current[pathParts[i]]) {
      if (current[pathParts[i]].properties) {
        current = current[pathParts[i]].properties
      } else {
        current = current[pathParts[i]]
      }
    }
  }
  
  const originalVarName = pathParts[pathParts.length - 1]
  if (current[originalVarName]) {
    const originalVariable = current[originalVarName]
    const copyVarName = `${originalVarName}_copy_${Date.now()}`
    
    // 深度复制变量定义
    current[copyVarName] = JSON.parse(JSON.stringify(originalVariable))
    current[copyVarName].title = (current[copyVarName].title || originalVarName) + ' 复制'
    
    console.log('Variable copied successfully:', copyVarName)
    
    // 强制更新树数据
    varsSchema.value = { ...varsSchema.value }
    
    // 选中新复制的变量
    const newPath = pathParts.slice(0, -1).concat([copyVarName]).join('.')
    selectedKeys.value = [newPath]
    onSelectVariable([newPath])
    
    message.success('变量复制成功')
  }
}


// 删除变量
const deleteVariable = (path) => {
  console.log('Deleting variable:', path) // 调试日志
  
  const pathParts = path.split('.')
  const varName = pathParts[pathParts.length - 1]
  
  if (pathParts.length === 1) {
    // 删除根级变量
    delete varsSchema.value[varName]
    console.log('Deleted root variable:', varName)
  } else {
    // 删除嵌套变量
    let current = varsSchema.value
    for (let i = 0; i < pathParts.length - 2; i++) {
      if (current[pathParts[i]]) {
        current = current[pathParts[i]]
      } else {
        console.error('Parent path not found:', pathParts[i])
        return
      }
    }
    
    const parentVarName = pathParts[pathParts.length - 2]
    if (current[parentVarName] && current[parentVarName].properties) {
      delete current[parentVarName].properties[varName]
      console.log('Deleted nested variable:', varName, 'from', parentVarName)
    } else {
      console.error('Parent variable properties not found:', parentVarName)
    }
  }
  
  // 清除选择
  if (selectedKeys.value.includes(path)) {
    selectedKeys.value = []
    selectedVariableData.value = null
  }
  
  // 强制更新树数据
  varsSchema.value = { ...varsSchema.value }
}

// 类型改变处理
const onTypeChange = (newType) => {
  if (!selectedVariableData.value) return
  
  // 重置类型相关的字段
  selectedVariableData.value.default = getDefaultValue(newType)
  selectedVariableData.value.ui.component = getDefaultComponent(newType)
  
  // 根据类型设置特殊字段
  if (newType === 'enum') {
    selectedVariableData.value.enum = []
  } else {
    delete selectedVariableData.value.enum
  }
  
  if (newType === 'array') {
    selectedVariableData.value.items = { type: 'string' }
    selectedVariableData.value.minItems = undefined
  } else {
    delete selectedVariableData.value.items
    delete selectedVariableData.value.minItems
  }
  
  if (newType === 'object') {
    selectedVariableData.value.properties = {}
  } else {
    delete selectedVariableData.value.properties
  }
  
  if (['integer', 'number'].includes(newType)) {
    selectedVariableData.value.min = undefined
    selectedVariableData.value.max = undefined
  } else {
    delete selectedVariableData.value.min
    delete selectedVariableData.value.max
  }
  
  if (newType === 'string') {
    selectedVariableData.value.pattern = undefined
  } else {
    delete selectedVariableData.value.pattern
  }
}

// Schema 预览内容
const schemaContent = computed(() => {
  const schema = { vars_schema: varsSchema.value }
  if (previewFormat.value === 'yaml') {
    return YAML.dump(schema, {
      indent: 2,
      lineWidth: -1,
      noRefs: true,
      sortKeys: false
    })
  } else {
    return JSON.stringify(schema, null, 2)
  }
})

// 初始化CodeMirror编辑器
const initSchemaEditor = () => {
  if (!schemaEditorRef.value) return
  
  const extensions = [
    basicSetup,
    EditorView.theme({
      '&': { 
        height: '100%',
        minHeight: '100%'
      },
      '.cm-scroller': { 
        fontFamily: 'Monaco, Menlo, monospace',
        minHeight: '100%'
      },
      '.cm-content': {
        minHeight: '100%',
        padding: '12px'
      },
      '.cm-editor': {
        height: '100%',
        minHeight: '100%'
      },
      '.cm-focused': { outline: 'none' }
    }),
    EditorState.readOnly.of(true),
    EditorView.lineWrapping
  ]
  
  if (previewFormat.value === 'yaml') {
    extensions.push(yaml())
  } else {
    extensions.push(json())
  }
  
  if (schemaEditor) {
    schemaEditor.destroy()
  }
  
  schemaEditor = new EditorView({
    state: EditorState.create({
      doc: schemaContent.value || '{}', // 确保有默认内容
      extensions
    }),
    parent: schemaEditorRef.value
  })
}

// 更新编辑器内容
const updateSchemaEditor = () => {
  if (!schemaEditor) return
  
  const transaction = schemaEditor.state.update({
    changes: {
      from: 0,
      to: schemaEditor.state.doc.length,
      insert: schemaContent.value
    }
  })
  
  schemaEditor.dispatch(transaction)
}

// 重新初始化编辑器（切换格式时）
const reinitSchemaEditor = async () => {
  await nextTick()
  initSchemaEditor()
}

// 加载现有变量定义
const loadVariableDefinitions = async () => {
  try {
    const response = await request({
      url: `/api/v1/templates/${props.templateId}/vars-schema`,
      method: 'GET'
    })
    if (response.data.data && response.data.data.vars_schema) {
      varsSchema.value = response.data.data.vars_schema
    } else {
      // 如果没有数据，创建一些测试数据用于调试
      varsSchema.value = {
        'app_name': {
          type: 'string',
          title: '应用名称',
          description: '应用的名称',
          required: true,
          default: 'my-app',
          ui: { panel: true, order: 10, group: '基础信息', component: 'input' },
          naming_policy: 'go_snake'
        },
        'database': {
          type: 'object',
          title: '数据库配置',
          description: '数据库相关配置',
          required: true,
          properties: {
            'host': {
              type: 'string',
              title: '主机地址',
              description: '数据库主机地址',
              required: true,
              default: 'localhost',
              ui: { panel: true, order: 20, group: '数据库', component: 'input' }
            },
            'port': {
              type: 'integer',
              title: '端口号',
              description: '数据库端口号',
              required: true,
              default: 3306,
              ui: { panel: true, order: 21, group: '数据库', component: 'input' }
            }
          },
          ui: { panel: true, order: 30, group: '基础信息', component: 'input' }
        }
      }
      console.log('Created test data:', varsSchema.value)
    }
  } catch (error) {
    console.error('加载变量定义失败:', error)
    // 出错时也创建测试数据
    varsSchema.value = {
      'test_var': {
        type: 'string',
        title: '测试变量',
        description: '用于测试的变量',
        required: false,
        default: 'test',
        ui: { panel: true, order: 10, group: '测试', component: 'input' }
      }
    }
  }
}

// 保存变量定义
const saveVariables = async () => {
  saving.value = true
  try {
    // 同步当前编辑的变量数据到schema
    if (selectedVariableData.value && selectedKeys.value.length > 0) {
      const path = selectedKeys.value[0]
      updateVariableInSchema(path, selectedVariableData.value)
    }
    
    await request({
      url: `/api/v1/templates/${props.templateId}/vars-schema`,
      method: 'POST',
      data: {
        vars_schema: varsSchema.value
      }
    })
    message.success('变量定义保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 更新变量到schema
const updateVariableInSchema = (path, variableData) => {
  const pathParts = path.split('.')
  let current = varsSchema.value
  
  // 导航到目标位置
  for (let i = 0; i < pathParts.length - 1; i++) {
    if (current[pathParts[i]]) {
      if (current[pathParts[i]].properties) {
        current = current[pathParts[i]].properties
      } else {
        current = current[pathParts[i]]
      }
    }
  }
  
  const varName = pathParts[pathParts.length - 1]
  if (current[varName]) {
    // 只更新变量的配置数据，不包括path
    const { path: _, ...configData } = variableData
    Object.assign(current[varName], configData)
  }
}

// 监听变量数据变化，实时同步到schema
watch(selectedVariableData, (newData) => {
  if (newData && selectedKeys.value.length > 0) {
    updateVariableInSchema(selectedKeys.value[0], newData)
  }
}, { deep: true })

// 拖拽调整高度
const startResize = (e) => {
  isResizing.value = true
  const startY = e.clientY
  const startHeight = drawerHeight.value

  const onMouseMove = (e) => {
    if (!isResizing.value) return
    
    const deltaY = startY - e.clientY // 向上为正，向下为负
    const newHeight = Math.min(Math.max(startHeight + deltaY, minHeight), maxHeight)
    drawerHeight.value = newHeight
  }

  const onMouseUp = () => {
    isResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }

  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ns-resize'
}

// 初始化面板宽度
const initPanelWidths = () => {
  const layoutElement = document.querySelector('.expose-layout')
  if (layoutElement) {
    totalLayoutWidth.value = layoutElement.offsetWidth
    rightPanelWidth.value = Math.floor(totalLayoutWidth.value / 3)
    leftPanelWidth.value = 250
  }
}

// 左侧面板拖拽调整
const startLeftResize = (e) => {
  isLeftResizing.value = true
  const startX = e.clientX
  const startWidth = leftPanelWidth.value
  
  const onMouseMove = (e) => {
    if (!isLeftResizing.value) return
    
    const deltaX = e.clientX - startX
    const newWidth = Math.min(Math.max(startWidth + deltaX, 100), totalLayoutWidth.value - 400) // 最小100px，保证其他面板至少400px
    leftPanelWidth.value = newWidth
  }
  
  const onMouseUp = () => {
    isLeftResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }
  
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ew-resize'
}

// 右侧面板拖拽调整
const startRightResize = (e) => {
  isRightResizing.value = true
  const startX = e.clientX
  const startWidth = rightPanelWidth.value
  
  const onMouseMove = (e) => {
    if (!isRightResizing.value) return
    
    const deltaX = startX - e.clientX // 右侧面板向左拖拽为正
    const newWidth = Math.min(Math.max(startWidth + deltaX, 150), totalLayoutWidth.value - leftPanelWidth.value - 200) // 最小150px，保证左侧和中间至少200px
    rightPanelWidth.value = newWidth
  }
  
  const onMouseUp = () => {
    isRightResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }
  
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ew-resize'
}

// 监听抽屉显示状态
watch(visible, async (show) => {
  if (show) {
    loadVariableDefinitions()
    await nextTick()
    initPanelWidths()
    initSchemaEditor()
  }
})

// 监听预览格式变化
watch(previewFormat, () => {
  reinitSchemaEditor()
})

// 监听varsSchema变化
watch(varsSchema, () => {
  updateSchemaEditor()
}, { deep: true })

// 组件销毁时清理编辑器
onMounted(() => {
  // 组件被销毁时清理
  return () => {
    if (schemaEditor) {
      schemaEditor.destroy()
      schemaEditor = null
    }
  }
})
</script>

<style scoped>
.resize-handle {
  position: relative;
  height: 12px;
  cursor: ns-resize;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 8px;
}

.resize-handle:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.handle-bar {
  width: 60px;
  height: 4px;
  background: #d9d9d9;
  border-radius: 2px;
  transition: background-color 0.2s;
}

.resize-handle:hover .handle-bar {
  background: #1890ff;
}

.expose-layout {
  display: flex;
  height: calc(100% - 20px);
  gap: 0;
  position: relative;
}

.left-panel,
.center-panel,
.right-panel {
  background: #f8f8f8;
  border-radius: 6px;
  padding: 16px;
  overflow-y: auto;
  flex-shrink: 0;
}

.left-panel {
  min-width: 100px;
  margin-right: 4px;
}

.center-panel {
  min-width: 200px;
  margin: 0 4px;
  flex: 1;
}

.right-panel {
  min-width: 150px;
  margin-left: 4px;
}

.panel-resizer {
  width: 8px;
  background: transparent;
  cursor: ew-resize;
  position: relative;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.panel-resizer:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.panel-resizer::before {
  content: '';
  width: 2px;
  height: 20px;
  background: #d9d9d9;
  border-radius: 1px;
  transition: background-color 0.2s;
}

.panel-resizer:hover::before {
  background: #1890ff;
}

.left-resizer {
  order: 1;
}

.center-panel {
  order: 2;
}

.right-resizer {
  order: 3;
}

.right-panel {
  order: 4;
}

.left-panel {
  order: 0;
}

.panel-title {
  font-weight: 600;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.file-list {
  margin-bottom: 16px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 8px;
  margin-bottom: 4px;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.file-item:hover {
  background: #e6f7ff;
}

.file-item.active {
  background: #1890ff;
  color: white;
}

.file-icon {
  margin-right: 8px;
}

.usage-section {
  border-top: 1px solid #e0e0e0;
  padding-top: 12px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #666;
}

.usage-list {
  max-height: 150px;
  overflow-y: auto;
}

.usage-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 8px;
  margin-bottom: 2px;
  background: white;
  border-radius: 3px;
  font-size: 12px;
}

.var-name {
  color: #1890ff;
  font-family: monospace;
}

.var-count {
  color: #666;
}

.variables-list {
  margin-bottom: 16px;
  max-height: 200px;
  overflow-y: auto;
}

.variable-card {
  background: white;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  border: 1px solid #e0e0e0;
  transition: all 0.2s;
}

.variable-card:hover {
  border-color: #1890ff;
}

.variable-card.active {
  border-color: #1890ff;
  background: #f0f9ff;
}

.variable-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.variable-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.variable-name {
  font-weight: 600;
}

.variable-type {
  font-size: 12px;
  padding: 2px 6px;
  background: #f0f0f0;
  border-radius: 3px;
  color: #666;
}

.variable-desc {
  font-size: 12px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.variable-form {
  border-top: 1px solid #e0e0e0;
  padding-top: 16px;
}

.preview-content {
  height: calc(100% - 50px);
  overflow: hidden;
}

.code-preview {
  height: 100%;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.code-preview :deep(.cm-editor) {
  height: 100% !important;
  min-height: 100%;
}

.code-preview :deep(.cm-scroller) {
  font-size: 12px;
  min-height: 100%;
}

.code-preview :deep(.cm-content) {
  min-height: 100%;
  padding: 12px;
}

.code-preview :deep(.cm-focused) {
  outline: none;
}

.json-preview {
  font-family: monospace;
  font-size: 12px;
  background: white;
  padding: 12px;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
  margin: 0;
  white-space: pre-wrap;
}


/* 变量树样式 */
.variable-tree {
  height: calc(100% - 50px);
  overflow-y: auto;
}

.var-icon {
  margin-right: 4px;
}

.var-icon.var-string {
  color: #1890ff; /* 蓝色 - 文本 */
}

.var-icon.var-integer,
.var-icon.var-number {
  color: #52c41a; /* 绿色 - 数字 */
}

.var-icon.var-boolean {
  color: #fa8c16; /* 橙色 - 布尔值 */
}

.var-icon.var-array {
  color: #f5222d; /* 红色 - 数组 */
}

.var-icon.var-object {
  color: #722ed1; /* 紫色 - 对象 */
}

.var-icon.var-enum {
  color: #13c2c2; /* 青色 - 枚举 */
}

.var-icon.var-secret {
  color: #eb2f96; /* 粉色 - 密码 */
}

.node-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.n-tree-node:hover .node-actions {
  opacity: 1;
}

.tree-node-content {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
}

.node-title {
  flex: 1;
}

/* 表单样式 */
.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px;
  background: #fafafa;
  border-radius: 6px;
}

.var-path {
  font-family: monospace;
  color: #666;
  font-size: 12px;
}

.form-section-title {
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  padding-bottom: 4px;
  border-bottom: 1px solid #e8e8e8;
}

.empty-selection {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.empty-tree-hint {
  padding: 32px;
  color: #888;
  text-align: center;
  user-select: none;
  cursor: context-menu;
  font-size: 14px;
  border: 2px dashed #e0e0e0;
  border-radius: 6px;
  margin: 16px;
  transition: all 0.2s;
}

.empty-tree-hint:hover {
  border-color: #1890ff;
  color: #1890ff;
  background: rgba(24, 144, 255, 0.05);
}
</style>