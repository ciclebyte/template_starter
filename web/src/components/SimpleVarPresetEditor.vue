<template>
    <div class="simple-var-preset-editor">
        <div class="editor-layout">
            <!-- 左侧树形结构 -->
            <div class="tree-panel" :style="{ width: treePanelWidth + 'px' }" @contextmenu="onTreeAreaContextMenu">
                <div class="tree-header">
                    <h4>数据结构</h4>
                </div>
                <div class="tree-content" ref="treeContainer">
                    <div v-if="treeData.length === 0" class="empty-tree" @contextmenu="onTreeAreaContextMenu">
                        <p>暂无数据，右键添加根字段开始创建</p>
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
                        <n-button size="small" @click="importJson" quaternary>
                            <template #icon>
                                <n-icon><CloudUploadOutline /></n-icon>
                            </template>
                            导入文件
                        </n-button>
                        <n-button size="small" @click="exportJson" quaternary>
                            <template #icon>
                                <n-icon><CloudDownloadOutline /></n-icon>
                            </template>
                            导出文件
                        </n-button>
                        <n-button size="small" @click="copyJson" quaternary>
                            <template #icon>
                                <n-icon><CopyOutline /></n-icon>
                            </template>
                            复制
                        </n-button>
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

                    <n-form-item label="描述" path="description">
                        <n-input v-model:value="nodeForm.description" type="textarea" placeholder="字段描述（可选）" :rows="2" />
                    </n-form-item>

                    <n-form-item label="包含子字段" path="hasChildren">
                        <n-switch v-model:value="nodeForm.hasChildren" />
                        <span style="margin-left: 8px; color: #666; font-size: 12px;">
                            {{ nodeForm.hasChildren ? '该字段可以包含子字段' : '该字段为普通字段' }}
                        </span>
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
    CopyOutline, CutOutline, ClipboardOutline, CloudUploadOutline, CloudDownloadOutline
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
const originalNodeData = ref(null) // 保存原始节点数据
const isModalCanceling = ref(false) // 标识是否正在取消编辑

// 导入导出相关 (现在使用文件选择器)

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


const nodeForm = reactive({
    key: '',
    description: '',
    hasChildren: false // 是否包含子字段
})

const nodeFormRules = {
    key: {
        required: true,
        message: '请输入字段名',
        trigger: ['input', 'blur']
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
            hasChildren: !!(node.children && node.children.length > 0),
            description: node.description,
            path: node.path,
            prefix: () => getFieldIcon(node),
            suffix: () => getNodeSuffix(node),
            children: node.children ? convertToNTreeFormat(node.children) : [],
            class: getDragClass(nodeKey)
        }
    })
}

// 获取字段图标
const getFieldIcon = (node) => {
    if (node.children !== undefined) {
        if (node.children && node.children.length > 0) {
            // 有子字段的容器节点显示文件夹图标
            return h(NIcon, { size: 16, color: '#1890ff' }, { default: () => h(DocumentTextOutline) })
        } else {
            // 空的容器节点显示空文件夹图标
            return h(NIcon, { size: 16, color: '#52c41a' }, { default: () => h(ClipboardOutline) })
        }
    } else {
        // 不允许子字段的普通字段显示字段图标
        return h(NIcon, { size: 16, color: '#722ed1' }, { default: () => h(CreateOutline) })
    }
}

// 获取节点后缀
const getNodeSuffix = (node) => {
    if (node.children !== undefined && (!node.children || node.children.length === 0)) {
        // 空容器节点显示可扩展标识
        return h(NTag, { size: 'tiny', type: 'info', style: { fontSize: '10px' } }, { default: () => '可扩展' })
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
        const { label, children } = node
        
        if (children && children.length > 0) {
            // 有子字段的节点转换为对象
            result[label] = convertTreeToObject(children)
        } else {
            // 普通字段设置为占位符
            result[label] = "{{." + label + "}}"
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
        },
        onDblclick(e) {
            e.preventDefault()
            e.stopPropagation()
            // 双击打开编辑弹窗
            const node = findNodeByKey(treeData.value, option.key)
            if (node) {
                openNodeModal(node)
            }
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
    const canAddChildren = !node.children || node.children.length === 0 || (node.children && node.children.length > 0)
    
    // 检查节点是否允许子字段
    const nodeData = findNodeByKey(treeData.value, node.key)
    const allowChildren = !nodeData || nodeData.children !== undefined
    
    const menuOptions = []
    
    // 只有允许子字段的节点才显示"添加子字段"选项
    if (allowChildren) {
        menuOptions.push({ label: '添加子字段', key: 'addChild', icon: () => h(NIcon, null, { default: () => h(AddOutline) }) })
        menuOptions.push({ type: 'divider', key: 'divider1' })
    }
    
    menuOptions.push(
        { label: '编辑', key: 'edit', icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) },
        { label: '复制', key: 'copy', icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
        { type: 'divider', key: 'divider2' },
        { label: '删除', key: 'delete', icon: () => h(NIcon, null, { default: () => h(TrashOutline) }) }
    )
    
    dropdownOptions.value = menuOptions
}

// 空白区域右键菜单
const onTreeAreaContextMenu = (e) => {
    e.preventDefault()
    e.stopPropagation()
    dropdownNode.value = null
    dropdownOptions.value = [
        { label: '添加根字段', key: 'addRoot', icon: () => h(NIcon, null, { default: () => h(AddOutline) }) }
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
                const nodeData = findNodeByKey(treeData.value, dropdownNode.value.key)
                if (nodeData && nodeData.children !== undefined) {
                    openNodeModal(null, dropdownNode.value.path)
                } else {
                    message.warning('该节点不允许包含子字段')
                }
            }
            break
        case 'addRoot':
            openNodeModal(null, [])
            break
        case 'edit':
            if (dropdownNode.value) {
                openNodeModal(dropdownNode.value)
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

// 打开节点编辑弹窗
const openNodeModal = (node = null, parentPath = []) => {
    resetNodeForm()
    isModalCanceling.value = false
    
    if (node) {
        // 编辑模式 - 保存原始数据
        const nodeData = findNodeByKey(treeData.value, node.key)
        if (nodeData) {
            originalNodeData.value = JSON.parse(JSON.stringify(nodeData))
            editingNode.value = nodeData
            editingNodePath.value = nodeData.path
            nodeForm.key = nodeData.label
            nodeForm.description = nodeData.description || ''
            nodeForm.hasChildren = nodeData.children !== undefined
        }
    } else {
        // 新增模式
        originalNodeData.value = null
        editingNode.value = null
        editingNodePath.value = parentPath
        nodeForm.hasChildren = false
    }
    
    showNodeModal.value = true
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


const addChildNode = (parentNode) => {
    openNodeModal(null, parentNode.path)
}

const editNode = (node) => {
    openNodeModal(node)
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
    nodeForm.description = ''
    nodeForm.hasChildren = false
}


const saveNode = async () => {
    try {
        await nodeFormRef.value?.validate()
        
        // 如果是编辑模式，检查是否要移除子字段支持
        if (editingNode.value && !nodeForm.hasChildren) {
            const currentNode = findNodeByKey(treeData.value, editingNode.value.key)
            if (currentNode && currentNode.children && currentNode.children.length > 0) {
                message.warning('该节点已包含子字段，无法设置为不包含子字段。请先删除所有子字段。')
                return
            }
        }
        
        nodeSaving.value = true
        
        const newNode = {
            key: nodeForm.key,
            label: nodeForm.key,
            path: [...editingNodePath.value, nodeForm.key],
            description: nodeForm.description,
            children: nodeForm.hasChildren ? [] : undefined
        }
        
        if (editingNode.value) {
            // 编辑现有节点，保留现有的子字段
            const currentNode = findNodeByKey(treeData.value, editingNode.value.key)
            if (currentNode && nodeForm.hasChildren && currentNode.children) {
                newNode.children = currentNode.children
            }
            updateNodeInTree(newNode)
        } else {
            // 添加新节点
            addNodeToTree(newNode)
        }
        
        syncTreeToJson()
        
        // 保存成功消息
        const isEdit = !!editingNode.value
        
        // 清空状态
        originalNodeData.value = null
        showNodeModal.value = false
        editingNode.value = null
        resetNodeForm()
        
        message.success(isEdit ? '节点更新成功' : '节点添加成功')
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
    // 如果是编辑模式且有原始数据，需要恢复
    if (editingNode.value && originalNodeData.value) {
        isModalCanceling.value = true
        
        // 恢复原始节点数据
        const restoreNode = (nodes) => {
            for (let i = 0; i < nodes.length; i++) {
                if (nodes[i].key === editingNode.value.key) {
                    nodes[i] = JSON.parse(JSON.stringify(originalNodeData.value))
                    return true
                }
                if (nodes[i].children) {
                    if (restoreNode(nodes[i].children)) {
                        return true
                    }
                }
            }
            return false
        }
        
        restoreNode(treeData.value)
        syncTreeToJson()
    }
    
    showNodeModal.value = false
    editingNode.value = null
    originalNodeData.value = null
    isModalCanceling.value = false
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
watch(() => [nodeForm.key, nodeForm.description, nodeForm.hasChildren], () => {
    // 在取消编辑时不触发更新
    if (isModalCanceling.value) return
    
    if (selectedNode.value && showNodeModal.value && editingNode.value) {
        // 实时更新编辑中的节点信息（仅在弹窗打开时）
        const updatedNode = {
            ...editingNode.value,
            label: nodeForm.key,
            description: nodeForm.description,
            children: nodeForm.hasChildren ? (editingNode.value.children || []) : undefined
        }
        
        // 更新节点路径中的key
        const newPath = [...editingNode.value.path]
        if (newPath.length > 0) {
            newPath[newPath.length - 1] = nodeForm.key
            updatedNode.path = newPath
            updatedNode.key = newPath.join('.')
        }
        
        // 更新树中的节点
        updateNodeInTree(updatedNode, editingNode.value.path)
        editingNode.value = updatedNode
        
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
    if (info.node) {
        const nodeData = findNodeByKey(treeData.value, info.node.key)
        // 只有允许子字段的节点才能接受拖拽
        if (nodeData && nodeData.children !== undefined) {
            dragOverNode.value = info.node
        }
    }
}

const onDragLeave = () => {
    dragOverNode.value = null
}

const onDragOver = (info) => {
    const { event, node } = info
    
    if (node) {
        const nodeData = findNodeByKey(treeData.value, node.key)
        // 检查目标节点是否允许子字段
        if (nodeData && nodeData.children !== undefined) {
            event.preventDefault()
            event.dataTransfer.dropEffect = 'move'
            dragOverNode.value = node
        } else {
            event.dataTransfer.dropEffect = 'none'
        }
    }
}

const onDrop = (info) => {
    const { dragNode, node } = info
    
    if (!dragNode || !node || dragNode.key === node.key) {
        clearDragState()
        return
    }
    
    // 检查目标节点是否允许子字段
    const targetNodeData = findNodeByKey(treeData.value, node.key)
    if (!targetNodeData || targetNodeData.children === undefined) {
        message.warning('该节点不允许包含子字段')
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

// JSON文件导入导出功能

const importJson = () => {
    // 创建隐藏的文件输入元素
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.json'
    input.style.display = 'none'
    
    input.onchange = (event) => {
        const file = event.target.files[0]
        if (file) {
            handleJsonFileImport(file)
        }
        document.body.removeChild(input)
    }
    
    document.body.appendChild(input)
    input.click()
}

const exportJson = () => {
    if (!editorView) return
    
    try {
        const content = editorView.state.doc.toString()
        
        // 验证JSON格式
        try {
            JSON.parse(content)
        } catch (error) {
            message.error('当前JSON格式无效，无法导出')
            return
        }
        
        // 创建并下载文件
        const blob = new Blob([content], { type: 'application/json' })
        const url = URL.createObjectURL(blob)
        
        const link = document.createElement('a')
        link.href = url
        link.download = 'var-preset-schema.json'
        link.style.display = 'none'
        
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        
        // 清理URL对象
        URL.revokeObjectURL(url)
        
        message.success('JSON文件导出成功')
    } catch (error) {
        console.error('导出失败:', error)
        message.error('导出失败')
    }
}

const handleJsonFileImport = (file) => {
    if (!file.name.toLowerCase().endsWith('.json')) {
        message.error('请选择.json格式的文件')
        return
    }
    
    const reader = new FileReader()
    
    reader.onload = (event) => {
        try {
            const content = event.target.result
            
            // 解析JSON内容
            let jsonObject
            try {
                jsonObject = JSON.parse(content)
            } catch (error) {
                message.error('文件内容不是有效的JSON格式: ' + error.message)
                return
            }
            
            // 转换为树结构
            const newTreeData = convertJsonToTree(jsonObject)
            
            // 更新树数据
            treeData.value = newTreeData
            
            // 同步到编辑器
            syncTreeToJson()
            
            message.success(`JSON文件 "${file.name}" 导入成功`)
        } catch (error) {
            console.error('导入失败:', error)
            message.error('导入文件失败: ' + error.message)
        }
    }
    
    reader.onerror = () => {
        message.error('读取文件失败')
    }
    
    reader.readAsText(file, 'UTF-8')
}

const copyJson = () => {
    if (!editorView) return
    
    try {
        const content = editorView.state.doc.toString()
        
        // 验证JSON格式
        try {
            JSON.parse(content)
        } catch (error) {
            message.error('当前JSON格式无效，无法复制')
            return
        }
        
        // 复制到剪贴板
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(content).then(() => {
                message.success('JSON内容已复制到剪贴板')
            }).catch(() => {
                // 备用方法
                fallbackCopyToClipboard(content)
            })
        } else {
            // 备用方法
            fallbackCopyToClipboard(content)
        }
    } catch (error) {
        console.error('复制失败:', error)
        message.error('复制失败')
    }
}

const fallbackCopyToClipboard = (text) => {
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    
    try {
        const successful = document.execCommand('copy')
        if (successful) {
            message.success('JSON内容已复制到剪贴板')
        } else {
            message.error('复制失败，请手动复制')
        }
    } catch (err) {
        message.error('复制失败，请手动复制')
        console.error('复制失败:', err)
    }
    
    document.body.removeChild(textArea)
}

const convertJsonToTree = (obj, parentPath = []) => {
    if (!obj || typeof obj !== 'object') {
        return []
    }

    const result = []
    
    for (const [key, value] of Object.entries(obj)) {
        const currentPath = [...parentPath, key]
        const node = {
            key: key,
            label: key,
            path: currentPath,
            description: ''
        }

        if (value !== null && typeof value === 'object' && !Array.isArray(value)) {
            // 对象类型，递归处理子节点
            node.children = convertJsonToTree(value, currentPath)
        } else if (Array.isArray(value) && value.length > 0 && typeof value[0] === 'object') {
            // 对象数组，转换为可扩展的容器节点
            node.children = []
            const sampleObject = value[0]
            if (sampleObject && typeof sampleObject === 'object') {
                node.children = convertJsonToTree(sampleObject, currentPath)
            }
        } else {
            // 原始值，设为不包含子字段的节点
            node.children = undefined
        }

        result.push(node)
    }

    return result
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