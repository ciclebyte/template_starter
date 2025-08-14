<template>
    <div class="simple-var-preset-editor">
        <div class="editor-layout">
            <!-- 左侧树形结构 -->
            <div class="tree-panel" :style="{ width: treePanelWidth + 'px' }" @contextmenu="onTreeAreaContextMenu">
                <div class="tree-header">
                    <h4>数据结构</h4>
                    <div class="tree-actions">
                        <n-button size="small" type="primary" @click="addRootNode">
                            <template #icon>
                                <n-icon><AddOutline /></n-icon>
                            </template>
                            添加根节点
                        </n-button>
                    </div>
                </div>
                <div class="tree-content" ref="treeContainer">
                    <div v-if="treeData.length === 0" class="empty-tree" @contextmenu="onTreeAreaContextMenu">
                        <p>暂无数据，右键或点击"添加根节点"开始创建</p>
                    </div>
                    <n-tree
                        v-if="treeData.length > 0"
                        :data="treeDataComputed"
                        :selected-keys="selectedNode ? [selectedNode.key] : []"
                        :render-label="renderLabel"
                        :node-props="nodeProps"
                        :render-switcher-icon="renderSwitcherIcon"
                        @update:selected-keys="onSelectNode"
                        @update:expanded-keys="updateExpanded"
                        draggable
                        @dragstart="onDragStart"
                        @drag-enter="onDragEnter"
                        @drag-leave="onDragLeave"
                        @drag-over="onDragOver"
                        @drop="onDrop"
                    />
                </div>
                
                <!-- 右键菜单 -->
                <n-dropdown
                    to="body"
                    trigger="manual"
                    :show="showDropdown"
                    :options="dropdownOptions"
                    :x="dropdownX"
                    :y="dropdownY"
                    @select="handleDropdownSelect"
                    @clickoutside="handleDropdownClickoutside"
                    class="tree-dropdown-menu"
                />
            </div>

            <!-- 左侧拖拽分隔线 -->
            <div class="resizer tree-resizer" @mousedown="startResize('tree', $event)"></div>

            <!-- 中间详情编辑区域 -->
            <div class="detail-panel" :style="{ width: detailPanelWidth + 'px' }">
                <div class="detail-header">
                    <h4>{{ selectedNode ? '编辑节点' : '选择节点' }}</h4>
                    <div class="detail-actions" v-if="selectedNode">
                        <n-button size="small" @click="saveCurrentNode" type="primary" :loading="nodeSaving">
                            <template #icon>
                                <n-icon><SaveOutline /></n-icon>
                            </template>
                            保存
                        </n-button>
                    </div>
                </div>
                <div class="detail-content">
                    <div v-if="!selectedNode" class="empty-detail">
                        <n-empty size="medium" description="请从左侧选择一个节点进行编辑">
                            <template #icon>
                                <n-icon size="48" :component="DocumentTextOutline" />
                            </template>
                        </n-empty>
                    </div>
                    <div v-else class="node-form">
                        <n-form ref="nodeFormRef" :model="nodeForm" :rules="nodeFormRules" label-placement="top">
                            <n-form-item label="字段名" path="key">
                                <n-input v-model:value="nodeForm.key" placeholder="请输入字段名" />
                            </n-form-item>

                            <n-form-item label="数据类型" path="type">
                                <n-select v-model:value="nodeForm.type" placeholder="选择数据类型" :options="typeOptions" @update:value="handleTypeChange" />
                            </n-form-item>

                            <n-form-item label="描述" path="description">
                                <n-input v-model:value="nodeForm.description" type="textarea" placeholder="字段描述" :rows="2" />
                            </n-form-item>

                            <n-form-item v-if="nodeForm.type === 'string'" label="默认值" path="defaultValue">
                                <n-input v-model:value="nodeForm.defaultValue" placeholder="默认字符串值" />
                            </n-form-item>

                            <n-form-item v-if="nodeForm.type === 'number'" label="默认值" path="defaultValue">
                                <n-input-number v-model:value="nodeForm.defaultValue" placeholder="默认数值" style="width: 100%" />
                            </n-form-item>

                            <n-form-item v-if="nodeForm.type === 'boolean'" label="默认值" path="defaultValue">
                                <n-switch v-model:value="nodeForm.defaultValue" />
                            </n-form-item>

                            <n-form-item label="是否必填" path="required">
                                <n-switch v-model:value="nodeForm.required" />
                            </n-form-item>
                        </n-form>
                    </div>
                </div>
            </div>

            <!-- 中间拖拽分隔线 -->
            <div class="resizer detail-resizer" @mousedown="startResize('detail', $event)"></div>

            <!-- 右侧 JSON 编辑器 -->
            <div class="editor-panel">
                <div class="editor-header">
                    <h4>JSON 预览</h4>
                    <div class="editor-actions">
                        <n-button size="small" @click="formatJson" quaternary>
                            <template #icon>
                                <n-icon><CodeOutline /></n-icon>
                            </template>
                            格式化
                        </n-button>
                        <n-button size="small" @click="syncFromJson" quaternary>
                            <template #icon>
                                <n-icon><SyncOutline /></n-icon>
                            </template>
                            同步到树
                        </n-button>
                    </div>
                </div>
                <div class="editor-content">
                    <div class="codemirror-container" ref="editorContainer"></div>
                </div>
                <div class="editor-status">
                    <span :class="['status-indicator', jsonValid ? 'valid' : 'invalid']">
                        {{ jsonValid ? '✓ JSON 有效' : '✗ JSON 无效' }}
                    </span>
                </div>
            </div>
        </div>

        <!-- 节点编辑弹窗 -->
        <n-modal v-model:show="showNodeModal" :mask-closable="false">
            <n-card style="width: 500px" :title="nodeModalTitle" :bordered="false" size="huge">
                <template #header-extra>
                    <n-button quaternary circle @click="closeNodeModal">
                        <template #icon>
                            <n-icon><CloseOutline /></n-icon>
                        </template>
                    </n-button>
                </template>

                <n-form ref="nodeFormRef" :model="nodeForm" :rules="nodeFormRules" label-placement="top">
                    <n-form-item label="字段名" path="key">
                        <n-input v-model:value="nodeForm.key" placeholder="请输入字段名" />
                    </n-form-item>

                    <n-form-item label="数据类型" path="type">
                        <n-select v-model:value="nodeForm.type" placeholder="选择数据类型" :options="typeOptions" @update:value="handleTypeChange" />
                    </n-form-item>

                    <n-form-item label="描述" path="description">
                        <n-input v-model:value="nodeForm.description" type="textarea" placeholder="字段描述" :rows="2" />
                    </n-form-item>

                    <n-form-item v-if="nodeForm.type === 'string'" label="默认值" path="defaultValue">
                        <n-input v-model:value="nodeForm.defaultValue" placeholder="默认字符串值" />
                    </n-form-item>

                    <n-form-item v-if="nodeForm.type === 'number'" label="默认值" path="defaultValue">
                        <n-input-number v-model:value="nodeForm.defaultValue" placeholder="默认数值" style="width: 100%" />
                    </n-form-item>

                    <n-form-item v-if="nodeForm.type === 'boolean'" label="默认值" path="defaultValue">
                        <n-switch v-model:value="nodeForm.defaultValue" />
                    </n-form-item>

                    <n-form-item label="是否必填" path="required">
                        <n-switch v-model:value="nodeForm.required" />
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="closeNodeModal">取消</n-button>
                        <n-button type="primary" @click="saveNode" :loading="nodeSaving">
                            {{ editingNode ? '更新' : '添加' }}
                        </n-button>
                    </div>
                </template>
            </n-card>
        </n-modal>
    </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick, h } from 'vue'
import {
    NButton, NIcon, NModal, NCard, NForm, NFormItem, NInput, NInputNumber,
    NSelect, NSwitch, NTag, NEmpty, NTree, NDropdown, useMessage
} from 'naive-ui'
import {
    AddOutline, CloseOutline, CreateOutline, TrashOutline, CodeOutline,
    SyncOutline, SaveOutline, DocumentTextOutline, ChevronForwardOutline,
    CopyOutline, CutOutline, ClipboardOutline
} from '@vicons/ionicons5'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { json } from '@codemirror/lang-json'

const props = defineProps({
    modelValue: {
        type: String,
        default: ''
    },
    readonly: {
        type: Boolean,
        default: false
    }
})

const emit = defineEmits(['update:modelValue'])
const message = useMessage()

// 编辑器相关
const editorContainer = ref(null)
let editorView = null

// 数据状态
const treeData = ref([])
const jsonContent = ref('')
const jsonValid = ref(true)

// 面板宽度状态
const treePanelWidth = ref(280)
const detailPanelWidth = ref(320)

// 拖拽状态
const resizing = ref(false)
const resizeType = ref('')
const startX = ref(0)
const startWidth = ref(0)

// 节点选择和编辑相关
const selectedNode = ref(null)
const showNodeModal = ref(false)
const nodeFormRef = ref(null)
const editingNode = ref(null)
const editingNodePath = ref([])
const nodeSaving = ref(false)

// 右键菜单相关
const showDropdown = ref(false)
const dropdownX = ref(0)
const dropdownY = ref(0)
const dropdownNode = ref(null)
const dropdownOptions = ref([])

// 树形结构相关
const treeContainer = ref(null)
const expandedKeys = ref(new Set())
const draggedNode = ref(null)
const dragOverNode = ref(null)
const isDragging = ref(false)

// 编辑状态相关
const isEditing = ref(false)
const editingKey = ref('')
const editingValue = ref('')

const nodeForm = reactive({
    key: '',
    type: 'string',
    description: '',
    defaultValue: '',
    required: false
})

const typeOptions = [
    { label: '字符串', value: 'string' },
    { label: '数字', value: 'number' },
    { label: '布尔值', value: 'boolean' },
    { label: '对象', value: 'object' },
    { label: '数组', value: 'array' }
]

const nodeFormRules = {
    key: {
        required: true,
        message: '请输入字段名',
        trigger: ['input', 'blur']
    },
    type: {
        required: true,
        message: '请选择数据类型',
        trigger: ['change']
    }
}

const nodeModalTitle = computed(() => {
    return editingNode.value ? '编辑节点' : '添加节点'
})

// 转换为NTree格式的计算属性
const treeDataComputed = computed(() => {
    return convertToNTreeFormat(treeData.value)
})

// 转换树数据为NTree格式
const convertToNTreeFormat = (nodes) => {
    if (!Array.isArray(nodes)) return []
    
    return nodes.map(node => {
        const nodeKey = node.key || node.path?.join('.') || ''
        const isExpanded = expandedKeys.value.has(nodeKey)
        
        return {
            label: node.label,
            key: nodeKey,
            isLeaf: !node.children || node.children.length === 0,
            nodeType: node.nodeType,
            value: node.value,
            description: node.description,
            required: node.required,
            path: node.path,
            isEditing: node.isEditing || false,
            prefix: () => getTypeIcon(node.nodeType),
            suffix: () => getNodeSuffix(node),
            children: node.children ? convertToNTreeFormat(node.children) : [],
            class: getDragClass(nodeKey)
        }
    })
}

// 获取类型图标
const getTypeIcon = (type) => {
    const iconMap = {
        object: () => h(NIcon, { size: 16, color: '#1890ff' }, { default: () => h(DocumentTextOutline) }),
        array: () => h(NIcon, { size: 16, color: '#52c41a' }, { default: () => h(ClipboardOutline) }),
        string: () => h(NIcon, { size: 16, color: '#722ed1' }, { default: () => h(CreateOutline) }),
        number: () => h(NIcon, { size: 16, color: '#eb2f96' }, { default: () => h(CreateOutline) }),
        boolean: () => h(NIcon, { size: 16, color: '#fa8c16' }, { default: () => h(CreateOutline) })
    }
    return iconMap[type] || iconMap.string
}

// 获取节点后缀
const getNodeSuffix = (node) => {
    if (node.required) {
        return h(NTag, { size: 'tiny', type: 'error' }, { default: () => '*' })
    }
    return null
}

// 获取拖拽类名
const getDragClass = (nodeKey) => {
    const classes = []
    if (draggedNode.value && draggedNode.value.key === nodeKey) {
        classes.push('dragging')
    }
    if (dragOverNode.value && dragOverNode.value.key === nodeKey) {
        classes.push('drag-over')
    }
    return classes.join(' ')
}

// 初始化 CodeMirror 编辑器
const initEditor = () => {
    if (!editorContainer.value) return

    const state = EditorState.create({
        doc: props.modelValue || '{}',
        extensions: [
            basicSetup,
            json(),
            EditorView.updateListener.of((update) => {
                if (update.docChanged) {
                    const content = update.state.doc.toString()
                    jsonContent.value = content
                    validateJsonContent(content)
                    emit('update:modelValue', content)
                }
            })
        ]
    })

    editorView = new EditorView({
        state,
        parent: editorContainer.value
    })

    // 初始化时同步到树
    if (props.modelValue) {
        syncJsonToTree(props.modelValue)
    }
}

// JSON 处理
const handleJsonChange = (value) => {
    jsonContent.value = value
    validateJsonContent(value)
    emit('update:modelValue', value)
}

const validateJsonContent = (content) => {
    try {
        if (content.trim()) {
            JSON.parse(content)
        }
        jsonValid.value = true
    } catch (e) {
        jsonValid.value = false
    }
}

// 将 JSON 同步到树形结构
const syncJsonToTree = (jsonStr) => {
    try {
        const data = JSON.parse(jsonStr || '{}')
        treeData.value = convertObjectToTree(data)
    } catch (e) {
        console.error('JSON 解析失败:', e)
        treeData.value = []
    }
}

// 将对象转换为树形结构
const convertObjectToTree = (obj, path = []) => {
    if (typeof obj !== 'object' || obj === null) {
        return []
    }

    return Object.entries(obj).map(([key, value]) => {
        const currentPath = [...path, key]
        const nodeType = Array.isArray(value) ? 'array' : typeof value
        
        const node = {
            key: currentPath.join('.'),
            label: key,
            nodeType,
            value,
            path: currentPath,
            description: '',
            required: false,
            children: []
        }

        if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
            node.children = convertObjectToTree(value, currentPath)
        } else if (Array.isArray(value) && value.length > 0 && typeof value[0] === 'object') {
            node.children = convertObjectToTree(value[0], [...currentPath, '0'])
        }

        return node
    })
}

// 将树形结构转换为对象
const convertTreeToObject = (nodes) => {
    const result = {}
    
    nodes.forEach(node => {
        const { label, nodeType, value, children } = node
        
        if (nodeType === 'object' && children && children.length > 0) {
            result[label] = convertTreeToObject(children)
        } else if (nodeType === 'array') {
            if (children && children.length > 0) {
                result[label] = [convertTreeToObject(children)]
            } else {
                result[label] = []
            }
        } else {
            // 根据类型设置默认值
            switch (nodeType) {
                case 'string':
                    result[label] = value || ''
                    break
                case 'number':
                    result[label] = typeof value === 'number' ? value : 0
                    break
                case 'boolean':
                    result[label] = typeof value === 'boolean' ? value : false
                    break
                default:
                    result[label] = value
            }
        }
    })
    
    return result
}

// 拖拽调整面板宽度
const startResize = (type, event) => {
    resizing.value = true
    resizeType.value = type
    startX.value = event.clientX
    
    if (type === 'tree') {
        startWidth.value = treePanelWidth.value
    } else if (type === 'detail') {
        startWidth.value = detailPanelWidth.value
    }
    
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
    
    event.preventDefault()
}

const handleResize = (event) => {
    if (!resizing.value) return
    
    const deltaX = event.clientX - startX.value
    
    if (resizeType.value === 'tree') {
        const newWidth = Math.max(200, Math.min(500, startWidth.value + deltaX))
        treePanelWidth.value = newWidth
    } else if (resizeType.value === 'detail') {
        const newWidth = Math.max(250, Math.min(600, startWidth.value + deltaX))
        detailPanelWidth.value = newWidth
    }
}

const stopResize = () => {
    resizing.value = false
    resizeType.value = ''
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
}

// 树形交互方法
const onSelectNode = (keys) => {
    if (keys && keys.length > 0) {
        const nodeKey = keys[0]
        const node = findNodeByKey(treeData.value, nodeKey)
        if (node) {
            handleNodeSelect(node)
        }
    }
}

const updateExpanded = (keys) => {
    expandedKeys.value = new Set(keys)
}

const renderLabel = ({ option }) => {
    if (option.isEditing) {
        return h('input', {
            class: 'tree-input',
            value: editingValue.value,
            autofocus: true,
            placeholder: '输入字段名',
            onInput: e => (editingValue.value = e.target.value),
            onKeydown: e => {
                if (e.key === 'Enter') {
                    e.preventDefault()
                    confirmEdit()
                }
                if (e.key === 'Escape') {
                    e.preventDefault()
                    cancelEdit()
                }
            },
            onBlur: () => {
                setTimeout(() => {
                    if (isEditing.value) {
                        confirmEdit()
                    }
                }, 100)
            }
        })
    }
    return option.label
}

const renderSwitcherIcon = () => {
    return h(NIcon, null, { default: () => h(ChevronForwardOutline) })
}

const nodeProps = ({ option }) => {
    return {
        onContextmenu(e) {
            e.preventDefault()
            e.stopPropagation()
            dropdownNode.value = option
            setupContextMenu(option)
            showDropdown.value = true
            dropdownX.value = e.clientX
            dropdownY.value = e.clientY
        }
    }
}

// 查找节点
const findNodeByKey = (nodes, key) => {
    if (!nodes || !Array.isArray(nodes)) return null
    
    for (const node of nodes) {
        if (node.key === key || node.path?.join('.') === key) {
            return node
        }
        if (node.children) {
            const found = findNodeByKey(node.children, key)
            if (found) return found
        }
    }
    return null
}

// 设置右键菜单
const setupContextMenu = (node) => {
    const canHaveChildren = node.nodeType === 'object' || node.nodeType === 'array'
    
    if (canHaveChildren) {
        dropdownOptions.value = [
            { label: '添加子字段', key: 'addChild', icon: () => h(NIcon, null, { default: () => h(AddOutline) }) },
            { type: 'divider', key: 'divider1' },
            { label: '重命名', key: 'rename', icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) },
            { label: '复制', key: 'copy', icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
            { type: 'divider', key: 'divider2' },
            { label: '删除', key: 'delete', icon: () => h(NIcon, null, { default: () => h(TrashOutline) }) }
        ]
    } else {
        dropdownOptions.value = [
            { label: '重命名', key: 'rename', icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) },
            { label: '复制', key: 'copy', icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
            { type: 'divider', key: 'divider1' },
            { label: '删除', key: 'delete', icon: () => h(NIcon, null, { default: () => h(TrashOutline) }) }
        ]
    }
}

// 空白区域右键菜单
const onTreeAreaContextMenu = (e) => {
    e.preventDefault()
    e.stopPropagation()
    dropdownNode.value = null
    dropdownOptions.value = [
        { label: '添加根字段 (对象)', key: 'addRootObject', icon: () => h(NIcon, null, { default: () => h(DocumentTextOutline) }) },
        { label: '添加根字段 (数组)', key: 'addRootArray', icon: () => h(NIcon, null, { default: () => h(ClipboardOutline) }) },
        { label: '添加根字段 (字符串)', key: 'addRootString', icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) },
        { label: '添加根字段 (数字)', key: 'addRootNumber', icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) },
        { label: '添加根字段 (布尔)', key: 'addRootBoolean', icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) }
    ]
    showDropdown.value = true
    dropdownX.value = e.clientX
    dropdownY.value = e.clientY
}

// 右键菜单选择处理
const handleDropdownSelect = (key) => {
    showDropdown.value = false
    
    switch (key) {
        case 'addChild':
            if (dropdownNode.value) {
                addChildNode(dropdownNode.value)
            }
            break
        case 'addRootObject':
            addRootNodeWithType('object')
            break
        case 'addRootArray':
            addRootNodeWithType('array')
            break
        case 'addRootString':
            addRootNodeWithType('string')
            break
        case 'addRootNumber':
            addRootNodeWithType('number')
            break
        case 'addRootBoolean':
            addRootNodeWithType('boolean')
            break
        case 'rename':
            if (dropdownNode.value) {
                startRename(dropdownNode.value)
            }
            break
        case 'copy':
            if (dropdownNode.value) {
                copyNode(dropdownNode.value)
            }
            break
        case 'delete':
            if (dropdownNode.value) {
                deleteNode(dropdownNode.value)
            }
            break
    }
}

const handleDropdownClickoutside = () => {
    showDropdown.value = false
}

// 节点选择
const handleNodeSelect = (node) => {
    selectedNode.value = node
    nodeForm.key = node.label
    nodeForm.type = node.nodeType
    nodeForm.description = node.description || ''
    nodeForm.required = node.required || false
    
    // 设置默认值
    if (node.nodeType === 'string') {
        nodeForm.defaultValue = node.value || ''
    } else if (node.nodeType === 'number') {
        nodeForm.defaultValue = node.value || 0
    } else if (node.nodeType === 'boolean') {
        nodeForm.defaultValue = node.value || false
    } else {
        nodeForm.defaultValue = ''
    }
}

// 保存当前编辑的节点
const saveCurrentNode = async () => {
    if (!selectedNode.value) return
    
    try {
        await nodeFormRef.value?.validate()
        nodeSaving.value = true
        
        // 更新选中节点的信息
        const updatedNode = {
            ...selectedNode.value,
            label: nodeForm.key,
            nodeType: nodeForm.type,
            value: nodeForm.defaultValue,
            description: nodeForm.description,
            required: nodeForm.required
        }
        
        // 更新节点路径中的key
        const newPath = [...selectedNode.value.path]
        newPath[newPath.length - 1] = nodeForm.key
        updatedNode.path = newPath
        updatedNode.key = newPath.join('.')
        
        // 更新树中的节点
        updateNodeInTree(updatedNode, selectedNode.value.path)
        selectedNode.value = updatedNode
        
        syncTreeToJson()
        message.success('节点保存成功')
    } catch (error) {
        console.error('保存节点失败:', error)
    } finally {
        nodeSaving.value = false
    }
}

// 节点操作
const addRootNode = () => {
    resetNodeForm()
    editingNode.value = null
    editingNodePath.value = []
    showNodeModal.value = true
}

const addRootNodeWithType = (type) => {
    const newKey = `field_${Date.now()}`
    const newNode = {
        key: newKey,
        label: newKey,
        nodeType: type,
        value: getDefaultValue(type),
        path: [newKey],
        description: '',
        required: false,
        children: (type === 'object' || type === 'array') ? [] : undefined,
        isEditing: true
    }
    
    treeData.value.push(newNode)
    selectedNode.value = newNode
    startEdit(newNode, newKey)
    syncTreeToJson()
}

const startEdit = (node, initialValue) => {
    isEditing.value = true
    editingKey.value = node.key
    editingValue.value = initialValue || node.label
    
    // 将节点设置为编辑状态
    const updateNodeEditState = (nodes) => {
        for (const n of nodes) {
            if (n.key === node.key) {
                n.isEditing = true
                break
            }
            if (n.children) {
                updateNodeEditState(n.children)
            }
        }
    }
    updateNodeEditState(treeData.value)
}

const confirmEdit = () => {
    if (!isEditing.value || !editingValue.value.trim()) return
    
    const newLabel = editingValue.value.trim()
    
    // 更新节点
    const updateNode = (nodes) => {
        for (const node of nodes) {
            if (node.key === editingKey.value) {
                // 更新标签和路径
                const oldLabel = node.label
                node.label = newLabel
                node.isEditing = false
                
                // 更新路径
                if (node.path && node.path.length > 0) {
                    node.path[node.path.length - 1] = newLabel
                    node.key = node.path.join('.')
                }
                
                // 如果这是当前选中的节点，更新表单
                if (selectedNode.value && selectedNode.value.key === editingKey.value) {
                    selectedNode.value = node
                    nodeForm.key = newLabel
                }
                
                break
            }
            if (node.children) {
                updateNode(node.children)
            }
        }
    }
    
    updateNode(treeData.value)
    
    isEditing.value = false
    editingKey.value = ''
    editingValue.value = ''
    
    syncTreeToJson()
    message.success('节点重命名成功')
}

const cancelEdit = () => {
    // 如果是新节点，删除它
    if (editingKey.value.includes('field_' + Date.now())) {
        const removeNewNode = (nodes) => {
            for (let i = nodes.length - 1; i >= 0; i--) {
                if (nodes[i].key === editingKey.value) {
                    nodes.splice(i, 1)
                    break
                }
                if (nodes[i].children) {
                    removeNewNode(nodes[i].children)
                }
            }
        }
        removeNewNode(treeData.value)
    } else {
        // 取消编辑状态
        const cancelNodeEdit = (nodes) => {
            for (const node of nodes) {
                if (node.key === editingKey.value) {
                    node.isEditing = false
                    break
                }
                if (node.children) {
                    cancelNodeEdit(node.children)
                }
            }
        }
        cancelNodeEdit(treeData.value)
    }
    
    isEditing.value = false
    editingKey.value = ''
    editingValue.value = ''
}

const startRename = (node) => {
    startEdit(node, node.label)
}

const copyNode = (node) => {
    const copiedNode = JSON.parse(JSON.stringify(node))
    const newKey = `${node.label}_copy_${Date.now()}`
    
    // 更新复制节点的key和路径
    const updateCopiedNodeKeys = (n, parentPath = []) => {
        n.key = [...parentPath, n.label].join('.')
        n.path = [...parentPath, n.label]
        if (n.children) {
            n.children.forEach(child => updateCopiedNodeKeys(child, n.path))
        }
    }
    
    copiedNode.label = newKey
    updateCopiedNodeKeys(copiedNode, node.path.slice(0, -1))
    
    // 找到父节点并添加复制的节点
    if (node.path.length === 1) {
        // 根节点，直接添加到树中
        treeData.value.push(copiedNode)
    } else {
        // 子节点，添加到父节点的children中
        const parentPath = node.path.slice(0, -1)
        const parent = findNodeByPath(treeData.value, parentPath)
        if (parent && parent.children) {
            parent.children.push(copiedNode)
        }
    }
    
    syncTreeToJson()
    message.success('节点复制成功')
}

const findNodeByPath = (nodes, path) => {
    if (!path || path.length === 0) return null
    
    for (const node of nodes) {
        if (node.path && node.path.join('.') === path.join('.')) {
            return node
        }
        if (node.children) {
            const found = findNodeByPath(node.children, path)
            if (found) return found
        }
    }
    return null
}

const getDefaultValue = (type) => {
    switch (type) {
        case 'string': return ''
        case 'number': return 0
        case 'boolean': return false
        case 'object': return {}
        case 'array': return []
        default: return ''
    }
}

const addChildNode = (parentNode) => {
    resetNodeForm()
    editingNode.value = null
    editingNodePath.value = parentNode.path
    showNodeModal.value = true
}

const editNode = (node) => {
    editingNode.value = node
    editingNodePath.value = node.path
    nodeForm.key = node.label
    nodeForm.type = node.nodeType
    nodeForm.description = node.description || ''
    nodeForm.required = node.required || false
    
    // 设置默认值
    if (node.nodeType === 'string') {
        nodeForm.defaultValue = node.value || ''
    } else if (node.nodeType === 'number') {
        nodeForm.defaultValue = node.value || 0
    } else if (node.nodeType === 'boolean') {
        nodeForm.defaultValue = node.value || false
    }
    
    showNodeModal.value = true
}

const deleteNode = (node) => {
    // 从树中删除节点
    const deleteFromArray = (nodes, targetPath) => {
        for (let i = 0; i < nodes.length; i++) {
            if (nodes[i].path.join('.') === targetPath.join('.')) {
                nodes.splice(i, 1)
                return true
            }
            if (nodes[i].children) {
                if (deleteFromArray(nodes[i].children, targetPath)) {
                    return true
                }
            }
        }
        return false
    }
    
    deleteFromArray(treeData.value, node.path)
    syncTreeToJson()
    message.success('节点删除成功')
}

const resetNodeForm = () => {
    nodeForm.key = ''
    nodeForm.type = 'string'
    nodeForm.description = ''
    nodeForm.defaultValue = ''
    nodeForm.required = false
}

const handleTypeChange = (type) => {
    // 根据类型重置默认值
    switch (type) {
        case 'string':
            nodeForm.defaultValue = ''
            break
        case 'number':
            nodeForm.defaultValue = 0
            break
        case 'boolean':
            nodeForm.defaultValue = false
            break
        case 'object':
        case 'array':
            nodeForm.defaultValue = ''
            break
    }
}

const saveNode = async () => {
    try {
        await nodeFormRef.value?.validate()
        nodeSaving.value = true
        
        const newNode = {
            key: nodeForm.key,
            label: nodeForm.key,
            nodeType: nodeForm.type,
            value: nodeForm.defaultValue,
            path: [...editingNodePath.value, nodeForm.key],
            description: nodeForm.description,
            required: nodeForm.required,
            children: nodeForm.type === 'object' || nodeForm.type === 'array' ? [] : undefined
        }
        
        if (editingNode.value) {
            // 编辑现有节点
            updateNodeInTree(newNode)
        } else {
            // 添加新节点
            addNodeToTree(newNode)
        }
        
        syncTreeToJson()
        closeNodeModal()
        message.success(editingNode.value ? '节点更新成功' : '节点添加成功')
    } catch (error) {
        console.error('保存节点失败:', error)
    } finally {
        nodeSaving.value = false
    }
}

const addNodeToTree = (newNode) => {
    if (editingNodePath.value.length === 0) {
        // 添加根节点
        treeData.value.push(newNode)
    } else {
        // 添加子节点
        const findParent = (nodes, path) => {
            for (const node of nodes) {
                if (node.path.join('.') === path.join('.')) {
                    return node
                }
                if (node.children) {
                    const found = findParent(node.children, path)
                    if (found) return found
                }
            }
            return null
        }
        
        const parent = findParent(treeData.value, editingNodePath.value)
        if (parent) {
            if (!parent.children) parent.children = []
            parent.children.push(newNode)
        }
    }
}

const updateNodeInTree = (updatedNode, originalPath = null) => {
    const targetPath = originalPath || editingNode.value?.path
    if (!targetPath) return
    
    const updateInArray = (nodes) => {
        for (let i = 0; i < nodes.length; i++) {
            if (nodes[i].path.join('.') === targetPath.join('.')) {
                // 保留子节点
                updatedNode.children = nodes[i].children
                nodes[i] = updatedNode
                return true
            }
            if (nodes[i].children) {
                if (updateInArray(nodes[i].children)) {
                    return true
                }
            }
        }
        return false
    }
    
    updateInArray(treeData.value)
}

const closeNodeModal = () => {
    showNodeModal.value = false
    editingNode.value = null
    resetNodeForm()
}

// 同步树到 JSON
const syncTreeToJson = () => {
    const obj = convertTreeToObject(treeData.value)
    const jsonStr = JSON.stringify(obj, null, 2)
    
    if (editorView) {
        editorView.dispatch({
            changes: {
                from: 0,
                to: editorView.state.doc.length,
                insert: jsonStr
            }
        })
    }
    
    jsonContent.value = jsonStr
    emit('update:modelValue', jsonStr)
}

// 其他操作
const formatJson = () => {
    if (!editorView) return
    
    try {
        const content = editorView.state.doc.toString()
        const parsed = JSON.parse(content)
        const formatted = JSON.stringify(parsed, null, 2)
        
        editorView.dispatch({
            changes: {
                from: 0,
                to: editorView.state.doc.length,
                insert: formatted
            }
        })
        
        message.success('JSON 格式化成功')
    } catch (e) {
        message.error('JSON 格式错误，无法格式化')
    }
}

const syncFromJson = () => {
    if (!editorView) return
    
    try {
        const content = editorView.state.doc.toString()
        syncJsonToTree(content)
        message.success('已从 JSON 同步到树形结构')
    } catch (e) {
        message.error('JSON 格式错误，无法同步')
    }
}

// 监听 modelValue 变化
watch(() => props.modelValue, (newValue) => {
    if (editorView && newValue !== editorView.state.doc.toString()) {
        editorView.dispatch({
            changes: {
                from: 0,
                to: editorView.state.doc.length,
                insert: newValue || '{}'
            }
        })
        syncJsonToTree(newValue)
    }
})

// 监听nodeForm变化，实时更新JSON
watch(() => [nodeForm.key, nodeForm.type, nodeForm.defaultValue, nodeForm.description, nodeForm.required], () => {
    if (selectedNode.value) {
        // 实时更新选中节点的信息
        const updatedNode = {
            ...selectedNode.value,
            label: nodeForm.key,
            nodeType: nodeForm.type,
            value: nodeForm.defaultValue,
            description: nodeForm.description,
            required: nodeForm.required
        }
        
        // 更新节点路径中的key
        const newPath = [...selectedNode.value.path]
        newPath[newPath.length - 1] = nodeForm.key
        updatedNode.path = newPath
        updatedNode.key = newPath.join('.')
        
        // 更新树中的节点
        updateNodeInTree(updatedNode, selectedNode.value.path)
        selectedNode.value = updatedNode
        
        syncTreeToJson()
    }
}, { deep: true })

// 拖拽功能
const onDragStart = (info) => {
    if (info.dragNode) {
        draggedNode.value = info.dragNode
        isDragging.value = true
    }
}

const onDragEnter = (info) => {
    if (info.node && info.node.nodeType === 'object') {
        dragOverNode.value = info.node
    }
}

const onDragLeave = () => {
    dragOverNode.value = null
}

const onDragOver = (info) => {
    const { event, node } = info
    
    if (node && node.nodeType === 'object') {
        event.preventDefault()
        event.dataTransfer.dropEffect = 'move'
        dragOverNode.value = node
    }
}

const onDrop = (info) => {
    const { dragNode, node } = info
    
    if (!dragNode || !node || dragNode.key === node.key) {
        clearDragState()
        return
    }
    
    // 只能拖拽到对象类型的节点
    if (node.nodeType !== 'object') {
        clearDragState()
        return
    }
    
    // 不能拖拽到自己的子节点
    if (isDescendant(dragNode, node)) {
        clearDragState()
        return
    }
    
    // 执行移动操作
    moveNode(dragNode, node)
    clearDragState()
}

const clearDragState = () => {
    draggedNode.value = null
    dragOverNode.value = null
    isDragging.value = false
}

const isDescendant = (parentNode, childNode) => {
    if (!parentNode.children) return false
    
    for (const child of parentNode.children) {
        if (child.key === childNode.key) return true
        if (isDescendant(child, childNode)) return true
    }
    return false
}

const moveNode = (sourceNode, targetNode) => {
    // 从原位置移除
    const removeNode = (nodes, nodeKey) => {
        for (let i = 0; i < nodes.length; i++) {
            if (nodes[i].key === nodeKey) {
                return nodes.splice(i, 1)[0]
            }
            if (nodes[i].children) {
                const removed = removeNode(nodes[i].children, nodeKey)
                if (removed) return removed
            }
        }
        return null
    }
    
    const movedNode = removeNode(treeData.value, sourceNode.key)
    if (!movedNode) return
    
    // 更新路径
    const updateNodePath = (node, newParentPath) => {
        node.path = [...newParentPath, node.label]
        node.key = node.path.join('.')
        if (node.children) {
            node.children.forEach(child => updateNodePath(child, node.path))
        }
    }
    
    updateNodePath(movedNode, targetNode.path)
    
    // 添加到新位置
    if (!targetNode.children) {
        targetNode.children = []
    }
    targetNode.children.push(movedNode)
    
    syncTreeToJson()
    message.success('节点移动成功')
}

// 生命周期
onMounted(async () => {
    await nextTick()
    initEditor()
})

onUnmounted(() => {
    if (editorView) {
        editorView.destroy()
    }
    // 清理拖拽事件监听器
    stopResize()
})
</script>

<style scoped>
.simple-var-preset-editor {
    width: 100%;
    height: 100%;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    overflow: hidden;
}

.editor-layout {
    display: flex;
    height: 100%;
}

.tree-panel {
    border-right: 1px solid #e0e0e0;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
}

.detail-panel {
    border-right: 1px solid #e0e0e0;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
}

.editor-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.resizer {
    width: 4px;
    background: #f0f0f0;
    cursor: col-resize;
    flex-shrink: 0;
    position: relative;
    transition: background-color 0.2s ease;
}

.resizer:hover {
    background: #d0d0d0;
}

.resizer:active {
    background: #1890ff;
}

.resizer::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 2px;
    height: 20px;
    background: #bbb;
    border-radius: 1px;
}

.resizer:hover::before {
    background: #999;
}

.resizer:active::before {
    background: #fff;
}

.tree-header,
.detail-header,
.editor-header {
    padding: 12px 16px;
    border-bottom: 1px solid #e0e0e0;
    background: #fafafa;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.tree-header h4,
.detail-header h4,
.editor-header h4 {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
    color: #333;
}

.tree-actions,
.detail-actions,
.editor-actions {
    display: flex;
    gap: 8px;
}

.tree-content {
    flex: 1;
    padding: 8px;
    overflow-y: auto;
}

.detail-content {
    flex: 1;
    overflow-y: auto;
}

.empty-detail {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
}

.node-form {
    padding: 16px;
}

.empty-tree {
    text-align: center;
    padding: 40px 20px;
    color: #999;
}

.editor-content {
    flex: 1;
    overflow: hidden;
}

.codemirror-container {
    height: 100%;
    width: 100%;
}

.editor-status {
    padding: 8px 16px;
    border-top: 1px solid #e0e0e0;
    background: #fafafa;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
}

.status-indicator {
    font-weight: 500;
}

.status-indicator.valid {
    color: #52c41a;
}

.status-indicator.invalid {
    color: #ff4d4f;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}

/* CodeMirror 样式 */
:deep(.cm-editor) {
    height: 100%;
    font-size: 13px;
    border: none;
}

:deep(.cm-focused) {
    outline: none;
}

:deep(.cm-content) {
    padding: 8px;
}

:deep(.cm-scroller) {
    font-family: 'Courier New', 'Consolas', 'Monaco', monospace;
}

/* 树形结构样式 */
.tree-dropdown-menu {
    z-index: 2147483647 !important;
}

.tree-input {
    width: 100%;
    height: 22px;
    padding: 1px 4px;
    font-size: 13px;
    font-family: 'Segoe UI', 'Consolas', 'Monaco', monospace;
    background: #ffffff;
    border: 1px solid #007acc;
    border-radius: 2px;
    outline: none;
    color: #333333;
    line-height: 18px;
    box-shadow: 0 0 0 1px #007acc;
    margin: 0;
    display: block;
}

:deep(.tree-input:focus) {
    border-color: #1890ff;
    box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

/* 拖拽样式 */
:deep(.n-tree-node.dragging) {
    opacity: 0.5;
    background: rgba(24, 144, 255, 0.1);
}

:deep(.n-tree-node.drag-over) {
    background: rgba(24, 144, 255, 0.2);
    border: 2px dashed #1890ff;
    border-radius: 4px;
}

:deep(.n-tree-node[draggable="true"]) {
    cursor: grab;
}

:deep(.n-tree-node[draggable="true"]:active) {
    cursor: grabbing;
}

/* 树节点样式优化 */
:deep(.n-tree-node-content) {
    padding: 4px 8px;
    border-radius: 4px;
    transition: all 0.2s ease;
}

:deep(.n-tree-node-content:hover) {
    background: #f5f5f5;
}

:deep(.n-tree-node--selected .n-tree-node-content) {
    background: #e6f7ff;
    border: 1px solid #91d5ff;
}

:deep(.n-tree-node--selected .n-tree-node-content:hover) {
    background: #e6f7ff;
}

/* 类型图标样式 */
:deep(.n-tree-node-content__prefix) {
    margin-right: 6px;
}

:deep(.n-tree-node-content__suffix) {
    margin-left: 6px;
}
</style>