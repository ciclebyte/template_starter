<template>
    <div class="categories-manage">
        <div class="page-header">
            <div class="header-left">
                <h1>分类管理</h1>
            </div>
        </div>

        <div class="content-card">
            <div class="toolbar">
                <div class="toolbar-left">
                    <div class="search-section">
                        <span class="search-label">关键词：</span>
                        <n-input v-model:value="searchKeyword" placeholder="输入分类名称或描述进行搜索..." clearable
                            style="width: 300px" @input="handleSearch" @keyup.enter="handleSearch">
                            <template #prefix>
                                <n-icon>
                                    <SearchOutline />
                                </n-icon>
                            </template>
                        </n-input>
                        <n-button type="primary" @click="handleSearch" :loading="searching" style="margin-left: 8px">
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
                            添加
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

            <n-data-table :columns="columns" :data="categories" :loading="loading" :pagination="pagination"
                :row-key="(row) => row.id" class="categories-table" remote @update:page="handlePageChange"
                @update:page-size="handlePageSizeChange" />
        </div>

        <!-- 添加/编辑分类弹窗 -->
        <n-modal v-model:show="showAddModal" :mask-closable="false">
            <n-card style="width: 500px" :title="editingCategory ? '编辑分类' : '添加分类'" :bordered="false" size="huge"
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
                    <n-form-item label="分类名称" path="name">
                        <n-input v-model:value="formData.name" placeholder="请输入分类名称" :maxlength="50" show-count />
                    </n-form-item>

                    <n-form-item label="分类描述" path="description">
                        <n-input v-model:value="formData.description" type="textarea" placeholder="请输入分类描述（可选）"
                            :maxlength="200" show-count :rows="3" />
                    </n-form-item>

                    <n-form-item label="排序" path="sort">
                        <n-input-number v-model:value="formData.sort" placeholder="排序值，数字越小越靠前" :min="0" :max="9999"
                            style="width: 100%" />
                    </n-form-item>

                    <n-form-item label="状态" path="status">
                        <n-switch v-model:value="formData.status" :checked-value="1" :unchecked-value="0">
                            <template #checked>启用</template>
                            <template #unchecked>禁用</template>
                        </n-switch>
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="closeModal">取消</n-button>
                        <n-button type="primary" @click="handleSubmit" :loading="submitting">
                            {{ editingCategory ? '更新' : '添加' }}
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
                        确定要删除分类 <strong>"{{ deletingCategory?.name }}"</strong> 吗？
                    </p>
                    <p class="delete-warning">
                        此操作不可撤销，删除后相关模板的分类信息可能会受到影响。
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
import {
    NButton, NIcon, NInput, NDataTable, NModal, NCard, NForm, NFormItem,
    NInputNumber, NSwitch, useMessage
} from 'naive-ui'
import {
    AddOutline, SearchOutline, RefreshOutline, CloseOutline,
    TrashOutline, CreateOutline, EyeOutline, EyeOffOutline
} from '@vicons/ionicons5'
import { listCategories, addCategory, editCategory, deleteCategory } from '@/api/categories'

const message = useMessage()

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const searching = ref(false)
const categories = ref([])
const searchKeyword = ref('')

// 弹窗状态
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const editingCategory = ref(null)
const deletingCategory = ref(null)

// 表单数据
const formRef = ref(null)
const formData = reactive({
    name: '',
    description: '',
    sort: 0,
    status: 1
})

// 表单验证规则
const formRules = {
    name: {
        required: true,
        message: '请输入分类名称',
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
const columns = [
    {
        title: '序号',
        key: 'index',
        width: 70,
        render: (row, index) => {
            return (pagination.page - 1) * pagination.pageSize + index + 1
        }
    },
    {
        title: '分类名称',
        key: 'name',
        width: 180,
        render: (row) => h('span', { class: 'category-name' }, row.name)
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
        title: '排序',
        key: 'sort',
        width: 80,
        sorter: (a, b) => a.sort - b.sort
    },
    {
        title: '状态',
        key: 'status',
        width: 80,
        render: (row) => {
            return h(
                NSwitch,
                {
                    value: row.status === 1,
                    size: 'small',
                    onUpdateValue: (value) => handleStatusChange(row, value)
                }
            )
        }
    },
    {
        title: '创建时间',
        key: 'created_at',
        width: 150,
        render: (row) => formatDate(row.created_at)
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
]

// 计算属性已移除，现在使用后端搜索

// 方法
const loadCategories = async (searchParams = {}) => {
    try {
        loading.value = true
        const params = {
            page: pagination.page,
            pageSize: pagination.pageSize,
            ...searchParams
        }
        const response = await listCategories(params)
        categories.value = response.data.data.categoriesList || []
        pagination.itemCount = response.data.data.total || 0
    } catch (error) {
        console.error('获取分类列表失败:', error)
        message.error('获取分类列表失败')
    } finally {
        loading.value = false
    }
}

const refreshData = () => {
    loadCategories()
}

const handleSearch = async () => {
    try {
        searching.value = true
        pagination.page = 1 // 重置到第一页

        const searchParams = {}
        if (searchKeyword.value) {
            searchParams.name = searchKeyword.value
        }

        await loadCategories(searchParams)
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

const handleAdd = () => {
    editingCategory.value = null
    resetForm()
    showAddModal.value = true
}

const handleEdit = (category) => {
    editingCategory.value = category
    formData.name = category.name
    formData.description = category.description || ''
    formData.sort = category.sort || 0
    formData.status = category.status || 1
    showAddModal.value = true
}

const handleDelete = (category) => {
    deletingCategory.value = category
    showDeleteModal.value = true
}

const closeModal = () => {
    showAddModal.value = false
    editingCategory.value = null
    resetForm()
}

const resetForm = () => {
    formData.name = ''
    formData.description = ''
    formData.sort = 0
    formData.status = 1
    formRef.value?.restoreValidation()
}

const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        submitting.value = true

        const data = {
            name: formData.name,
            description: formData.description,
            sort: formData.sort,
            status: formData.status
        }

        if (editingCategory.value) {
            // 编辑
            await editCategory({ ...data, id: editingCategory.value.id })
            message.success('分类更新成功')
        } else {
            // 添加
            await addCategory(data)
            message.success('分类添加成功')
        }

        closeModal()
        loadCategories()
    } catch (error) {
        console.error('操作失败:', error)
        message.error(editingCategory.value ? '更新分类失败' : '添加分类失败')
    } finally {
        submitting.value = false
    }
}

const confirmDelete = async () => {
    try {
        deleting.value = true
        await deleteCategory({ id: deletingCategory.value.id })
        message.success('分类删除成功')
        showDeleteModal.value = false
        deletingCategory.value = null
        loadCategories()
    } catch (error) {
        console.error('删除分类失败:', error)
        message.error('删除分类失败')
    } finally {
        deleting.value = false
    }
}

const handleStatusChange = async (row, value) => {
    try {
        const newStatus = value ? 1 : 0
        await editCategory({
            id: row.id,
            name: row.name,
            description: row.description || '',
            sort: row.sort || 0,
            status: newStatus
        })

        // 更新本地数据
        const index = categories.value.findIndex(item => item.id === row.id)
        if (index !== -1) {
            categories.value[index].status = newStatus
        }

        message.success(`分类已${value ? '启用' : '禁用'}`)
    } catch (error) {
        console.error('状态切换失败:', error)
        message.error('状态切换失败')
        // 刷新数据以恢复原状态
        loadCategories()
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
onMounted(() => {
    loadCategories()
})
</script>

<style scoped>
.categories-manage {
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

.page-description {
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
}

.search-label {
    font-size: 14px;
    color: #333;
    font-weight: 500;
    white-space: nowrap;
}

.categories-table {
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

.category-name {
    font-weight: 500;
    color: #333;
}

.text-placeholder {
    color: #999;
    font-style: italic;
}

.action-buttons {
    display: flex;
    gap: 8px;
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