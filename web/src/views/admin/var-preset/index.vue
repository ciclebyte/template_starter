<template>
    <div class="var-preset-manage">
        <div class="page-header">
            <div class="header-left">
                <h1>变量预设管理</h1>
                <p>管理模板开发时的变量预设，提供快速填充功能</p>
            </div>
        </div>

        <div class="content-card">
            <div class="toolbar">
                <div class="toolbar-left">
                    <div class="search-section">
                        <span class="search-label">预设名称：</span>
                        <n-input v-model:value="searchKeyword" placeholder="输入预设名称进行搜索..." clearable
                            style="width: 300px" @input="handleSearch" @keyup.enter="handleSearch">
                            <template #prefix>
                                <n-icon>
                                    <SearchOutline />
                                </n-icon>
                            </template>
                        </n-input>
                        <n-select v-model:value="searchCategory" placeholder="选择分类" clearable
                            style="width: 120px; margin-left: 12px" @update:value="handleSearch"
                            :options="categoryOptions" />
                        <n-button type="primary" @click="handleSearch" :loading="searching" style="margin-left: 16px">
                            <template #icon>
                                <n-icon>
                                    <SearchOutline />
                                </n-icon>
                            </template>
                            搜索
                        </n-button>
                        <n-button @click="clearSearch" quaternary style="margin-left: 8px" v-if="searchKeyword || searchCategory">
                            清空
                        </n-button>
                        <n-button type="primary" @click="showAddModal = true" style="margin-left: 16px">
                            <template #icon>
                                <n-icon>
                                    <AddOutline />
                                </n-icon>
                            </template>
                            添加预设
                        </n-button>
                    </div>
                </div>
                <div class="toolbar-right">
                    <n-button @click="refreshData" :loading="loading" quaternary>
                        <template #icon>
                            <n-icon>
                                <RefreshOutline />
                            </n-icon>
                        </template>
                        刷新
                    </n-button>
                </div>
            </div>

            <n-data-table :columns="columns" :data="varPresets" :loading="loading" :pagination="pagination"
                :row-key="(row) => row.id" class="var-presets-table" remote @update:page="handlePageChange"
                @update:page-size="handlePageSizeChange" />
        </div>

        <!-- 添加/编辑变量预设弹窗 -->
        <n-modal v-model:show="showAddModal" :mask-closable="false">
            <n-card style="width: 600px" :title="editingPreset ? '编辑变量预设' : '添加变量预设'" :bordered="false" size="huge"
                role="dialog" aria-modal="true">
                <template #header-extra>
                    <n-button quaternary circle @click="closeModal">
                        <template #icon>
                            <n-icon>
                                <CloseOutline />
                            </n-icon>
                        </template>
                    </n-button>
                </template>

                <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
                    <n-grid :cols="2" :x-gap="16">
                        <n-grid-item>
                            <n-form-item label="预设名称" path="name">
                                <n-input v-model:value="formData.name" placeholder="请输入预设名称" :maxlength="50" show-count />
                            </n-form-item>
                        </n-grid-item>
                        <n-grid-item>
                            <n-form-item label="显示名称" path="displayName">
                                <n-input v-model:value="formData.displayName" placeholder="请输入显示名称" :maxlength="100" show-count />
                            </n-form-item>
                        </n-grid-item>
                    </n-grid>

                    <n-form-item label="描述" path="description">
                        <n-input v-model:value="formData.description" type="textarea" placeholder="请输入预设描述"
                            :maxlength="500" show-count :rows="3" />
                    </n-form-item>

                    <n-grid :cols="3" :x-gap="16">
                        <n-grid-item>
                            <n-form-item label="分类" path="category">
                                <n-select v-model:value="formData.category" placeholder="选择分类"
                                    :options="categoryOptions" />
                            </n-form-item>
                        </n-grid-item>
                        <n-grid-item>
                            <n-form-item label="图标" path="icon">
                                <n-input v-model:value="formData.icon" placeholder="图标名称或URL" />
                            </n-form-item>
                        </n-grid-item>
                        <n-grid-item>
                            <n-form-item label="排序权重" path="sort">
                                <n-input-number v-model:value="formData.sort" placeholder="数值越大越靠前" :min="0" :max="999"
                                    style="width: 100%" />
                            </n-form-item>
                        </n-grid-item>
                    </n-grid>

                    <n-grid :cols="2" :x-gap="16">
                        <n-grid-item>
                            <n-form-item label="版本" path="version">
                                <n-input v-model:value="formData.version" placeholder="如: 1.0" />
                            </n-form-item>
                        </n-grid-item>
                        <n-grid-item>
                            <n-form-item label="状态" path="isEnabled">
                                <n-switch v-model:value="formData.isEnabled" />
                                <span style="margin-left: 8px; color: #666;">{{ formData.isEnabled ? '启用' : '禁用' }}</span>
                            </n-form-item>
                        </n-grid-item>
                    </n-grid>

                    <n-form-item label="默认数据" path="defaultDataJson">
                        <n-input v-model:value="formData.defaultDataJson" type="textarea" placeholder='可选，提供默认数据值...'
                            :rows="4" />
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="closeModal">取消</n-button>
                        <n-button type="primary" @click="handleSubmit" :loading="submitting">
                            {{ editingPreset ? '更新' : '添加' }}
                        </n-button>
                    </div>
                </template>
            </n-card>
        </n-modal>

        <!-- 预览弹窗 -->
        <n-modal v-model:show="showPreviewModal" :mask-closable="true">
            <n-card style="width: 80%; max-width: 800px" title="预设详情" :bordered="false" size="huge"
                role="dialog" aria-modal="true">
                <template #header-extra>
                    <n-button quaternary circle @click="showPreviewModal = false">
                        <template #icon>
                            <n-icon>
                                <CloseOutline />
                            </n-icon>
                        </template>
                    </n-button>
                </template>

                <div v-if="previewData" class="preview-content">
                    <div class="preview-info">
                        <n-descriptions :column="2" bordered>
                            <n-descriptions-item label="预设名称">{{ previewData.name }}</n-descriptions-item>
                            <n-descriptions-item label="显示名称">{{ previewData.displayName }}</n-descriptions-item>
                            <n-descriptions-item label="分类">
                                <n-tag :type="previewData.category === 'system' ? 'info' : 'success'">
                                    {{ previewData.category === 'system' ? '系统' : '自定义' }}
                                </n-tag>
                            </n-descriptions-item>
                            <n-descriptions-item label="版本">{{ previewData.version }}</n-descriptions-item>
                            <n-descriptions-item label="状态">
                                <n-tag :type="previewData.isEnabled === 1 ? 'success' : 'error'">
                                    {{ previewData.isEnabled === 1 ? '启用' : '禁用' }}
                                </n-tag>
                            </n-descriptions-item>
                            <n-descriptions-item label="排序权重">{{ previewData.sort }}</n-descriptions-item>
                            <n-descriptions-item label="描述" :span="2">{{ previewData.description || '无' }}</n-descriptions-item>
                        </n-descriptions>
                    </div>

                    <div class="preview-schema" style="margin-top: 20px;">
                        <h4>数据结构模板：</h4>
                        <n-code :code="formatJson(previewData.schemaJson)" language="json" style="max-height: 300px; overflow-y: auto;" />
                    </div>

                    <div v-if="previewData.defaultDataJson" class="preview-default" style="margin-top: 20px;">
                        <h4>默认数据：</h4>
                        <n-code :code="formatJson(previewData.defaultDataJson)" language="json" style="max-height: 200px; overflow-y: auto;" />
                    </div>
                </div>
            </n-card>
        </n-modal>

        <!-- 删除确认弹窗 -->
        <n-modal v-model:show="showDeleteModal" :mask-closable="false">
            <n-card style="width: 400px" title="确认删除" :bordered="false" size="huge" role="dialog" aria-modal="true">
                <div class="delete-content">
                    <div class="delete-icon">
                        <n-icon size="48" color="#d03050">
                            <TrashOutline />
                        </n-icon>
                    </div>
                    <p class="delete-message">
                        确定要删除变量预设 <strong>"{{ deletingPreset?.displayName || deletingPreset?.name }}"</strong> 吗？
                    </p>
                    <p class="delete-warning">
                        删除后将解除与所有模板的关联关系，此操作不可撤销。
                    </p>
                </div>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="showDeleteModal = false">取消</n-button>
                        <n-button type="error" @click="confirmDelete" :loading="deleting">
                            确认删除
                        </n-button>
                    </div>
                </template>
            </n-card>
        </n-modal>
    </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import {
    NButton, NIcon, NInput, NDataTable, NModal, NCard, NForm, NFormItem,
    NInputNumber, NTag, useMessage, NSelect, NSwitch, NGrid, NGridItem,
    NDescriptions, NDescriptionsItem, NCode
} from 'naive-ui'
import {
    AddOutline, SearchOutline, RefreshOutline, CloseOutline,
    TrashOutline, CreateOutline, EyeOutline, ToggleOutline, 
    DocumentTextOutline
} from '@vicons/ionicons5'
import { 
    listVarPresets, addVarPreset, editVarPreset, deleteVarPreset,
    batchDeleteVarPresets, getVarPresetDetail, toggleVarPreset
} from '@/api/varPreset'
import SimpleVarPresetEditor from '@/components/SimpleVarPresetEditor.vue'

const message = useMessage()
const router = useRouter()

// 选项数据
const categoryOptions = [
    { label: '系统', value: 'system' },
    { label: '自定义', value: 'custom' }
]

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const searching = ref(false)
const varPresets = ref([])
const searchKeyword = ref('')
const searchCategory = ref('')

// 弹窗状态
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const showPreviewModal = ref(false)
const editingPreset = ref(null)
const deletingPreset = ref(null)
const previewData = ref(null)

// 表单数据
const formRef = ref(null)
const formData = reactive({
    name: '',
    displayName: '',
    description: '',
    category: 'custom',
    schemaJson: '',
    defaultDataJson: '',
    icon: '',
    sort: 0,
    version: '1.0',
    isEnabled: true
})

// 表单验证规则
const formRules = {
    name: {
        required: true,
        message: '请输入预设名称',
        trigger: ['input', 'blur']
    },
    displayName: {
        required: true,
        message: '请输入显示名称',
        trigger: ['input', 'blur']
    },
    category: {
        required: true,
        message: '请选择分类',
        trigger: ['change']
    }
}

// 分页配置
const pagination = reactive({
    page: 1,
    pageSize: 20,
    showSizePicker: true,
    pageSizes: [20, 50, 100],
    showQuickJumper: true,
    prefix: (info) => `共 ${info.itemCount} 条`
})

// 表格列配置
const columns = computed(() => [
    {
        title: '序号',
        key: 'index',
        width: 70,
        render: (row, index) => {
            return (pagination.page - 1) * pagination.pageSize + index + 1
        }
    },
    {
        title: '预设名称',
        key: 'name',
        width: 120,
        render: (row) => h(NTag, { 
            type: 'info', 
            size: 'medium'
        }, { 
            default: () => row.name
        })
    },
    {
        title: '显示名称',
        key: 'displayName',
        width: 150,
        render: (row) => row.displayName || row.name
    },
    {
        title: '分类',
        key: 'category',
        width: 80,
        render: (row) => h(NTag, { 
            type: row.category === 'system' ? 'info' : 'success',
            size: 'small'
        }, { 
            default: () => row.category === 'system' ? '系统' : '自定义'
        })
    },
    {
        title: '描述',
        key: 'description',
        width: 180,
        ellipsis: {
            tooltip: true
        },
        render: (row) => row.description || h('span', { class: 'text-placeholder' }, '暂无描述')
    },
    {
        title: '版本',
        key: 'version',
        width: 80,
        render: (row) => row.version || '1.0'
    },
    {
        title: '排序',
        key: 'sort',
        width: 80,
        render: (row) => row.sort || 0
    },
    {
        title: '状态',
        key: 'isEnabled',
        width: 80,
        render: (row) => h(NTag, { 
            type: row.isEnabled === 1 ? 'success' : 'error',
            size: 'small'
        }, { 
            default: () => row.isEnabled === 1 ? '启用' : '禁用'
        })
    },
    {
        title: '创建时间',
        key: 'createdAt',
        width: 150,
        render: (row) => formatDate(row.createdAt)
    },
    {
        title: '操作',
        key: 'actions',
        width: 380,
        render: (row) => {
            return h('div', { class: 'action-buttons' }, [
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'info',
                        secondary: true,
                        onClick: () => handlePreview(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(EyeOutline) }),
                        default: () => '预览'
                    }
                ),
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'primary',
                        secondary: true,
                        style: { marginLeft: '8px' },
                        onClick: () => handleEdit(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(CreateOutline) }),
                        default: () => '编辑'
                    }
                ),
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'info',
                        secondary: true,
                        style: { marginLeft: '8px' },
                        onClick: () => handleDesignSchema(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(DocumentTextOutline) }),
                        default: () => '设计结构'
                    }
                ),
                h(
                    NButton,
                    {
                        size: 'small',
                        type: row.isEnabled === 1 ? 'warning' : 'success',
                        secondary: true,
                        style: { marginLeft: '8px' },
                        onClick: () => handleToggle(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(ToggleOutline) }),
                        default: () => row.isEnabled === 1 ? '禁用' : '启用'
                    }
                ),
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'error',
                        secondary: true,
                        style: { marginLeft: '8px' },
                        disabled: row.category === 'system',
                        onClick: () => handleDelete(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(TrashOutline) }),
                        default: () => '删除'
                    }
                )
            ])
        }
    }
])

// 方法
const loadVarPresets = async (searchParams = {}) => {
    try {
        loading.value = true
        const params = {
            pageNum: pagination.page,
            pageSize: pagination.pageSize,
            ...searchParams
        }
        
        const response = await listVarPresets(params)
        varPresets.value = response.data.data.varPresetsList || []
        pagination.itemCount = response.data.data.total || 0
    } catch (error) {
        console.error('获取变量预设列表失败:', error)
        message.error('获取变量预设列表失败')
    } finally {
        loading.value = false
    }
}

const refreshData = () => {
    loadVarPresets()
}

const handleSearch = async () => {
    try {
        searching.value = true
        pagination.page = 1 // 重置到第一页

        const searchParams = {}
        if (searchKeyword.value) {
            searchParams.name = searchKeyword.value
        }
        if (searchCategory.value) {
            searchParams.category = searchCategory.value
        }

        await loadVarPresets(searchParams)
    } catch (error) {
        console.error('搜索失败:', error)
        message.error('搜索失败')
    } finally {
        searching.value = false
    }
}

const clearSearch = async () => {
    searchKeyword.value = ''
    searchCategory.value = ''
    await handleSearch()
}

const handlePageChange = async (page) => {
    pagination.page = page
    await handleSearch()
}

const handlePageSizeChange = async (pageSize) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    await handleSearch()
}

const handlePreview = async (preset) => {
    try {
        const response = await getVarPresetDetail({ id: preset.id })
        previewData.value = response.data.data.varPreset || response.data.data
        showPreviewModal.value = true
    } catch (error) {
        console.error('获取预设详情失败:', error)
        message.error('获取预设详情失败')
    }
}

const handleEdit = (preset) => {
    editingPreset.value = preset
    formData.name = preset.name
    formData.displayName = preset.displayName
    formData.description = preset.description || ''
    formData.category = preset.category
    formData.schemaJson = preset.schemaJson || ''
    formData.defaultDataJson = preset.defaultDataJson || ''
    formData.icon = preset.icon || ''
    formData.sort = preset.sort || 0
    formData.version = preset.version || '1.0'
    formData.isEnabled = preset.isEnabled === 1
    showAddModal.value = true
}

const handleDesignSchema = (preset) => {
    // 跳转到数据结构设计页面
    const route = router.resolve({
        name: 'var-preset-design',
        params: { id: preset.id }
    })
    window.open(route.href, '_blank')
}

const handleDelete = (preset) => {
    if (preset.category === 'system') {
        message.warning('系统预设不能删除')
        return
    }
    deletingPreset.value = preset
    showDeleteModal.value = true
}

const handleToggle = async (preset) => {
    try {
        const newStatus = preset.isEnabled === 1 ? 0 : 1
        await toggleVarPreset({
            id: preset.id,
            isEnabled: newStatus
        })
        message.success(`预设已${newStatus === 1 ? '启用' : '禁用'}`)
        await loadVarPresets()
    } catch (error) {
        console.error('切换预设状态失败:', error)
        message.error('切换预设状态失败')
    }
}

const closeModal = () => {
    showAddModal.value = false
    editingPreset.value = null
    resetForm()
}

const resetForm = () => {
    formData.name = ''
    formData.displayName = ''
    formData.description = ''
    formData.category = 'custom'
    formData.schemaJson = '{}'
    formData.defaultDataJson = ''
    formData.icon = ''
    formData.sort = 0
    formData.version = '1.0'
    formData.isEnabled = true
    formRef.value?.restoreValidation()
}

const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        
        // 验证默认数据JSON格式
        if (formData.defaultDataJson) {
            try {
                JSON.parse(formData.defaultDataJson)
            } catch (e) {
                message.error('默认数据格式不正确')
                return
            }
        }

        submitting.value = true

        const data = {
            name: formData.name,
            displayName: formData.displayName,
            description: formData.description,
            category: formData.category,
            schemaJson: formData.schemaJson || '{}', // 如果为空则设置为空对象
            defaultDataJson: formData.defaultDataJson,
            icon: formData.icon,
            sort: formData.sort,
            version: formData.version,
            isEnabled: formData.isEnabled ? 1 : 0
        }

        if (editingPreset.value) {
            await editVarPreset({ ...data, id: editingPreset.value.id })
            message.success('变量预设更新成功')
        } else {
            await addVarPreset(data)
            message.success('变量预设添加成功')
        }

        closeModal()
        loadVarPresets()
    } catch (error) {
        console.error('操作失败:', error)
        message.error(editingPreset.value ? '更新变量预设失败' : '添加变量预设失败')
    } finally {
        submitting.value = false
    }
}

const confirmDelete = async () => {
    try {
        deleting.value = true
        await deleteVarPreset({ id: deletingPreset.value.id })
        message.success('变量预设删除成功')
        showDeleteModal.value = false
        deletingPreset.value = null
        loadVarPresets()
    } catch (error) {
        console.error('删除变量预设失败:', error)
        message.error('删除变量预设失败')
    } finally {
        deleting.value = false
    }
}

const formatDate = (dateString) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    })
}

const formatJson = (jsonString) => {
    try {
        const obj = JSON.parse(jsonString)
        return JSON.stringify(obj, null, 2)
    } catch (e) {
        return jsonString
    }
}

// 生命周期
onMounted(async () => {
    await loadVarPresets()
})
</script>

<style scoped>
.var-preset-manage {
    padding: 0;
    background: transparent;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;
}

.header-left h1 {
    margin: 0 0 8px 0;
    font-size: 24px;
    font-weight: 600;
    color: #333;
}

.header-left p {
    margin: 0;
    color: #666;
    font-size: 14px;
}

.content-card {
    background: white;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.toolbar-left {
    flex: 1;
}

.toolbar-right {
    display: flex;
    gap: 12px;
}

.search-section {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.search-label {
    font-size: 14px;
    color: #333;
    font-weight: 500;
    white-space: nowrap;
}

.var-presets-table {
    margin-top: 16px;
}


.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}

.delete-content {
    text-align: center;
    padding: 20px 0;
}

.delete-icon {
    margin-bottom: 16px;
}

.delete-message {
    font-size: 16px;
    color: #333;
    margin-bottom: 12px;
}

.delete-warning {
    font-size: 14px;
    color: #d03050;
    background: #fff2f0;
    padding: 12px;
    border-radius: 6px;
    border: 1px solid #ffccc7;
}

.text-placeholder {
    color: #999;
    font-style: italic;
}

.action-buttons {
    display: flex;
    gap: 8px;
}

.preview-content {
    max-height: 70vh;
    overflow-y: auto;
}

.preview-info {
    margin-bottom: 20px;
}

:deep(.n-data-table th) {
    background: #fafbfc;
    font-weight: 600;
}

:deep(.n-data-table td) {
    border-bottom: 1px solid #f0f0f0;
}

:deep(.n-data-table tr:hover td) {
    background: #f9f9f9;
}

:deep(.n-input__textarea-el) {
    font-family: 'Courier New', monospace;
    font-size: 13px;
}
</style>