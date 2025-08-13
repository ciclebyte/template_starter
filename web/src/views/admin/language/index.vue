<template>
    <div class="languages-manage">
        <div class="page-header">
            <div class="header-left">
                <h1>语言管理</h1>
            </div>
        </div>

        <div class="content-card">
            <div class="toolbar">
                <div class="toolbar-left">
                    <div class="search-section">
                        <span class="search-label">关键词：</span>
                        <n-input v-model:value="searchKeyword" placeholder="输入语言名称或代码进行搜索..." clearable
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

            <n-data-table :columns="columns" :data="languages" :loading="loading" :pagination="pagination"
                :row-key="(row) => row.id" class="languages-table" remote @update:page="handlePageChange"
                @update:page-size="handlePageSizeChange" />
        </div>

        <!-- 添加/编辑语言弹窗 -->
        <n-modal v-model:show="showAddModal" :mask-closable="false">
            <n-card style="width: 500px" :title="editingLanguage ? '编辑语言' : '添加语言'" :bordered="false" size="huge"
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
                    <n-form-item label="语言名称" path="name">
                        <n-input v-model:value="formData.name" placeholder="请输入语言名称（如：JavaScript）" :maxlength="50"
                            show-count />
                    </n-form-item>

                    <n-form-item label="显示名称" path="displayName">
                        <n-input v-model:value="formData.displayName" placeholder="请输入显示名称（如：JavaScript）"
                            :maxlength="50" show-count />
                    </n-form-item>

                    <n-form-item label="语言代码" path="code">
                        <n-input v-model:value="formData.code" placeholder="请输入语言代码（如：js）" :maxlength="20" show-count />
                    </n-form-item>

                    <n-form-item label="语言颜色" path="color">
                        <n-color-picker v-model:value="formData.color" :show-alpha="false" />
                    </n-form-item>

                    <n-form-item label="排序" path="sort">
                        <n-input-number v-model:value="formData.sort" placeholder="排序值，数字越小越靠前" :min="0" :max="9999"
                            style="width: 100%" />
                    </n-form-item>

                    <n-form-item label="热门语言" path="isPopular">
                        <n-switch v-model:value="formData.isPopular" :checked-value="1" :unchecked-value="0">
                            <template #checked>是</template>
                            <template #unchecked>否</template>
                        </n-switch>
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="closeModal">取消</n-button>
                        <n-button type="primary" @click="handleSubmit" :loading="submitting">
                            {{ editingLanguage ? '更新' : '添加' }}
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
                        确定要删除语言 <strong>"{{ deletingLanguage?.name }}"</strong> 吗？
                    </p>
                    <p class="delete-warning">
                        此操作不可撤销，删除后相关模板的语言信息可能会受到影响。
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
    NInputNumber, NSwitch, NColorPicker, NTag, useMessage
} from 'naive-ui'
import {
    AddOutline, SearchOutline, RefreshOutline, CloseOutline,
    TrashOutline, CreateOutline, StarOutline, Star
} from '@vicons/ionicons5'
import { listLanguages, addLanguage, editLanguage, deleteLanguage } from '@/api/languages'

const message = useMessage()

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const searching = ref(false)
const languages = ref([])
const searchKeyword = ref('')

// 弹窗状态
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const editingLanguage = ref(null)
const deletingLanguage = ref(null)

// 表单数据
const formRef = ref(null)
const formData = reactive({
    name: '',
    displayName: '',
    code: '',
    color: '#18a058',
    sort: 0,
    isPopular: 0
})

// 表单验证规则
const formRules = {
    name: {
        required: true,
        message: '请输入语言名称',
        trigger: ['input', 'blur']
    },
    displayName: {
        required: true,
        message: '请输入显示名称',
        trigger: ['input', 'blur']
    },
    code: {
        required: true,
        message: '请输入语言代码',
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
        title: '语言名称',
        key: 'name',
        width: 150,
        render: (row) => h('span', { class: 'language-name' }, row.name)
    },
    {
        title: '显示名称',
        key: 'displayName',
        width: 150,
        render: (row) => h('span', { class: 'language-display-name' }, row.displayName || row.display_name || row.name)
    },
    {
        title: '语言代码',
        key: 'code',
        width: 100,
        render: (row) => h(
            NTag,
            {
                type: 'info',
                size: 'small',
                style: { fontFamily: 'monospace' }
            },
            { default: () => row.code }
        )
    },
    {
        title: '颜色',
        key: 'color',
        width: 80,
        render: (row) => h('div', {
            style: {
                display: 'flex',
                alignItems: 'center',
                gap: '8px'
            }
        }, [
            h('div', {
                style: {
                    width: '20px',
                    height: '20px',
                    borderRadius: '4px',
                    backgroundColor: row.color,
                    border: '1px solid #e0e0e0'
                }
            }),
            h('span', { style: { fontSize: '12px', color: '#666' } }, row.color)
        ])
    },
    {
        title: '排序',
        key: 'sort',
        width: 80,
        sorter: (a, b) => (a.sort || 0) - (b.sort || 0)
    },
    {
        title: '热门',
        key: 'isPopular',
        width: 80,
        render: (row) => {
            return h(
                NSwitch,
                {
                    value: (row.isPopular || row.is_popular) === 1,
                    size: 'small',
                    onUpdateValue: (value) => handlePopularChange(row, value)
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
const loadLanguages = async (searchParams = {}) => {
    try {
        loading.value = true
        const params = {
            page: pagination.page,
            pageSize: pagination.pageSize,
            ...searchParams
        }
        const response = await listLanguages(params)
        languages.value = response.data.data.languagesList || []
        pagination.itemCount = response.data.data.total || 0
    } catch (error) {
        console.error('获取语言列表失败:', error)
        message.error('获取语言列表失败')
    } finally {
        loading.value = false
    }
}

const refreshData = () => {
    loadLanguages()
}

const handleSearch = async () => {
    try {
        searching.value = true
        pagination.page = 1 // 重置到第一页

        const searchParams = {}
        if (searchKeyword.value) {
            searchParams.name = searchKeyword.value
        }

        await loadLanguages(searchParams)
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
    editingLanguage.value = null
    resetForm()
    showAddModal.value = true
}

const handleEdit = (language) => {
    editingLanguage.value = language
    formData.name = language.name
    formData.displayName = language.displayName || language.display_name || language.name
    formData.code = language.code || ''
    formData.color = language.color || '#18a058'
    formData.sort = language.sort || 0
    formData.isPopular = language.isPopular || language.is_popular || 0
    showAddModal.value = true
}

const handleDelete = (language) => {
    deletingLanguage.value = language
    showDeleteModal.value = true
}

const closeModal = () => {
    showAddModal.value = false
    editingLanguage.value = null
    resetForm()
}

const resetForm = () => {
    formData.name = ''
    formData.displayName = ''
    formData.code = ''
    formData.color = '#18a058'
    formData.sort = 0
    formData.isPopular = 0
    formRef.value?.restoreValidation()
}

const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        submitting.value = true

        const data = {
            name: formData.name,
            display_name: formData.displayName,
            code: formData.code,
            color: formData.color,
            sort: formData.sort,
            isPopular: formData.isPopular
        }

        if (editingLanguage.value) {
            // 编辑
            await editLanguage({ ...data, id: editingLanguage.value.id })
            message.success('语言更新成功')
        } else {
            // 添加
            await addLanguage(data)
            message.success('语言添加成功')
        }

        closeModal()
        loadLanguages()
    } catch (error) {
        console.error('操作失败:', error)
        message.error(editingLanguage.value ? '更新语言失败' : '添加语言失败')
    } finally {
        submitting.value = false
    }
}

const confirmDelete = async () => {
    try {
        deleting.value = true
        await deleteLanguage({ id: deletingLanguage.value.id })
        message.success('语言删除成功')
        showDeleteModal.value = false
        deletingLanguage.value = null
        loadLanguages()
    } catch (error) {
        console.error('删除语言失败:', error)
        message.error('删除语言失败')
    } finally {
        deleting.value = false
    }
}

const handlePopularChange = async (row, value) => {
    try {
        const newPopular = value ? 1 : 0
        await editLanguage({
            id: row.id,
            name: row.name,
            display_name: row.displayName || row.display_name || row.name,
            code: row.code || '',
            color: row.color || '#18a058',
            sort: row.sort || 0,
            isPopular: newPopular
        })

        // 更新本地数据
        const index = languages.value.findIndex(item => item.id === row.id)
        if (index !== -1) {
            languages.value[index].isPopular = newPopular
            languages.value[index].is_popular = newPopular
        }

        message.success(`语言已${value ? '设为热门' : '取消热门'}`)
    } catch (error) {
        console.error('热门状态切换失败:', error)
        message.error('热门状态切换失败')
        // 刷新数据以恢复原状态
        loadLanguages()
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
    loadLanguages()
})
</script>

<style scoped>
.languages-manage {
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

.languages-table {
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

.language-name {
    font-weight: 500;
    color: #333;
}

.language-display-name {
    color: #666;
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