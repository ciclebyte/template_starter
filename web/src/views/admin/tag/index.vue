<template>
    <div class="tags-manage">
        <div class="page-header">
            <div class="header-left">
                <h1>标签管理</h1>
                <p>管理模板标签，用于分类和筛选模板</p>
            </div>
        </div>

        <div class="content-card">
            <div class="toolbar">
                <div class="toolbar-left">
                    <div class="search-section">
                        <span class="search-label">标签名称：</span>
                        <n-input v-model:value="searchKeyword" placeholder="输入标签名称进行搜索..." clearable
                            style="width: 300px" @input="handleSearch" @keyup.enter="handleSearch">
                            <template #prefix>
                                <n-icon>
                                    <SearchOutline />
                                </n-icon>
                            </template>
                        </n-input>
                        <n-button type="primary" @click="handleSearch" :loading="searching" style="margin-left: 16px">
                            <template #icon>
                                <n-icon>
                                    <SearchOutline />
                                </n-icon>
                            </template>
                            搜索
                        </n-button>
                        <n-button @click="clearSearch" quaternary style="margin-left: 8px" v-if="searchKeyword">
                            清空
                        </n-button>
                        <n-button type="primary" @click="showAddModal = true" style="margin-left: 16px">
                            <template #icon>
                                <n-icon>
                                    <AddOutline />
                                </n-icon>
                            </template>
                            添加标签
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

            <n-data-table :columns="columns" :data="tags" :loading="loading" :pagination="pagination"
                :row-key="(row) => row.id" class="tags-table" remote @update:page="handlePageChange"
                @update:page-size="handlePageSizeChange" />
        </div>

        <!-- 添加/编辑标签弹窗 -->
        <n-modal v-model:show="showAddModal" :mask-closable="false">
            <n-card style="width: 500px" :title="editingTag ? '编辑标签' : '添加标签'" :bordered="false" size="huge"
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
                    <n-form-item label="标签名称" path="name">
                        <n-input v-model:value="formData.name" placeholder="请输入标签名称" :maxlength="50" show-count />
                    </n-form-item>

                    <n-form-item label="标签描述" path="description">
                        <n-input v-model:value="formData.description" type="textarea" placeholder="请输入标签描述"
                            :maxlength="500" show-count :rows="3" />
                    </n-form-item>

                    <n-form-item label="排序权重" path="sort">
                        <n-input-number v-model:value="formData.sort" placeholder="数值越大越靠前" :min="0" :max="999"
                            style="width: 100%" />
                        <template #suffix>
                            <span class="input-tip">数值越大越靠前</span>
                        </template>
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="closeModal">取消</n-button>
                        <n-button type="primary" @click="handleSubmit" :loading="submitting">
                            {{ editingTag ? '更新' : '添加' }}
                        </n-button>
                    </div>
                </template>
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
                        确定要删除标签 <strong>"{{ deletingTag?.name }}"</strong> 吗？
                    </p>
                    <p class="delete-warning">
                        删除后将自动解除与所有模板的关联关系，此操作不可撤销。
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

        <!-- 批量删除确认弹窗 -->
        <n-modal v-model:show="showBatchDeleteModal" :mask-closable="false">
            <n-card style="width: 400px" title="批量删除确认" :bordered="false" size="huge" role="dialog" aria-modal="true">
                <div class="delete-content">
                    <div class="delete-icon">
                        <n-icon size="48" color="#d03050">
                            <TrashOutline />
                        </n-icon>
                    </div>
                    <p class="delete-message">
                        确定要删除选中的 <strong>{{ checkedRowKeys.length }}</strong> 个标签吗？
                    </p>
                    <p class="delete-warning">
                        删除后将自动解除与所有模板的关联关系，此操作不可撤销。
                    </p>
                </div>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="showBatchDeleteModal = false">取消</n-button>
                        <n-button type="error" @click="confirmBatchDelete" :loading="batchDeleting">
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
import {
    NButton, NIcon, NInput, NDataTable, NModal, NCard, NForm, NFormItem,
    NInputNumber, NTag, useMessage, NCheckbox
} from 'naive-ui'
import {
    AddOutline, SearchOutline, RefreshOutline, CloseOutline,
    TrashOutline, CreateOutline, PricetagOutline
} from '@vicons/ionicons5'
import { 
    listTags, addTag, editTag, deleteTag, batchDeleteTags, getTagsWithCount
} from '@/api/tags'

const message = useMessage()

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const batchDeleting = ref(false)
const searching = ref(false)
const tags = ref([])
const searchKeyword = ref('')
const checkedRowKeys = ref([])

// 弹窗状态
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const showBatchDeleteModal = ref(false)
const editingTag = ref(null)
const deletingTag = ref(null)

// 表单数据
const formRef = ref(null)
const formData = reactive({
    name: '',
    description: '',
    sort: 0
})

// 表单验证规则
const formRules = {
    name: {
        required: true,
        message: '请输入标签名称',
        trigger: ['input', 'blur']
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
        type: 'selection',
        multiple: true
    },
    {
        title: '序号',
        key: 'index',
        width: 70,
        render: (row, index) => {
            return (pagination.page - 1) * pagination.pageSize + index + 1
        }
    },
    {
        title: '标签名称',
        key: 'name',
        width: 150,
        render: (row) => h(NTag, { 
            type: 'info', 
            size: 'medium',
            style: { cursor: 'pointer' }
        }, { 
            default: () => row.name,
            icon: () => h(NIcon, { size: 14 }, { default: () => h(PricetagOutline) })
        })
    },
    {
        title: '描述',
        key: 'description',
        minWidth: 200,
        ellipsis: {
            tooltip: true
        },
        render: (row) => row.description || h('span', { class: 'text-placeholder' }, '暂无描述')
    },
    {
        title: '排序权重',
        key: 'sort',
        width: 100,
        render: (row) => row.sort || 0
    },
    {
        title: '关联模板数',
        key: 'templateCount',
        width: 120,
        render: (row) => {
            const count = row.templateCount || 0
            return h(NTag, { 
                type: count > 0 ? 'success' : 'default', 
                size: 'small' 
            }, { 
                default: () => `${count} 个`
            })
        }
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
        width: 200,
        render: (row) => {
            return h('div', { class: 'action-buttons' }, [
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'primary',
                        secondary: true,
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
                        type: 'error',
                        secondary: true,
                        style: { marginLeft: '8px' },
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
const loadTags = async (searchParams = {}) => {
    try {
        loading.value = true
        const params = {
            pageNum: pagination.page,
            pageSize: pagination.pageSize,
            ...searchParams
        }
        
        // 使用带统计的接口
        const response = await getTagsWithCount()
        tags.value = response.data.data.tagsList || []
        
        // 如果有搜索条件，需要过滤数据
        let filteredTags = tags.value
        if (searchParams.name) {
            filteredTags = tags.value.filter(tag => 
                tag.name.toLowerCase().includes(searchParams.name.toLowerCase())
            )
        }
        
        // 手动分页处理
        const total = filteredTags.length
        const startIndex = (pagination.page - 1) * pagination.pageSize
        const endIndex = startIndex + pagination.pageSize
        tags.value = filteredTags.slice(startIndex, endIndex)
        pagination.itemCount = total
    } catch (error) {
        console.error('获取标签列表失败:', error)
        message.error('获取标签列表失败')
    } finally {
        loading.value = false
    }
}

const refreshData = () => {
    loadTags()
}

const handleSearch = async () => {
    try {
        searching.value = true
        pagination.page = 1 // 重置到第一页

        const searchParams = {}
        if (searchKeyword.value) {
            searchParams.name = searchKeyword.value
        }

        await loadTags(searchParams)
    } catch (error) {
        console.error('搜索失败:', error)
        message.error('搜索失败')
    } finally {
        searching.value = false
    }
}

const clearSearch = async () => {
    searchKeyword.value = ''
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

const handleEdit = (tag) => {
    editingTag.value = tag
    formData.name = tag.name
    formData.description = tag.description || ''
    formData.sort = tag.sort || 0
    showAddModal.value = true
}

const handleDelete = (tag) => {
    deletingTag.value = tag
    showDeleteModal.value = true
}

const handleBatchDelete = () => {
    if (checkedRowKeys.value.length === 0) {
        message.warning('请选择要删除的标签')
        return
    }
    showBatchDeleteModal.value = true
}

const closeModal = () => {
    showAddModal.value = false
    editingTag.value = null
    resetForm()
}

const resetForm = () => {
    formData.name = ''
    formData.description = ''
    formData.sort = 0
    formRef.value?.restoreValidation()
}

const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        submitting.value = true

        const data = {
            name: formData.name,
            description: formData.description,
            sort: formData.sort
        }

        if (editingTag.value) {
            await editTag({ ...data, id: editingTag.value.id })
            message.success('标签更新成功')
        } else {
            await addTag(data)
            message.success('标签添加成功')
        }

        closeModal()
        loadTags()
    } catch (error) {
        console.error('操作失败:', error)
        message.error(editingTag.value ? '更新标签失败' : '添加标签失败')
    } finally {
        submitting.value = false
    }
}

const confirmDelete = async () => {
    try {
        deleting.value = true
        await deleteTag({ id: deletingTag.value.id })
        message.success('标签删除成功')
        showDeleteModal.value = false
        deletingTag.value = null
        loadTags()
    } catch (error) {
        console.error('删除标签失败:', error)
        message.error('删除标签失败')
    } finally {
        deleting.value = false
    }
}

const confirmBatchDelete = async () => {
    try {
        batchDeleting.value = true
        await batchDeleteTags({ ids: checkedRowKeys.value })
        message.success(`成功删除 ${checkedRowKeys.value.length} 个标签`)
        showBatchDeleteModal.value = false
        checkedRowKeys.value = []
        loadTags()
    } catch (error) {
        console.error('批量删除标签失败:', error)
        message.error('批量删除标签失败')
    } finally {
        batchDeleting.value = false
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

// 生命周期
onMounted(async () => {
    await loadTags()
})
</script>

<style scoped>
.tags-manage {
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

.tags-table {
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

.input-tip {
    color: #999;
    font-size: 12px;
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
</style>