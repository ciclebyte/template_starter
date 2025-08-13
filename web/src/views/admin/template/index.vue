<template>
    <div class="templates-manage">
        <div class="page-header">
            <div class="header-left">
                <h1>模板管理</h1>
            </div>
        </div>

        <div class="content-card">
            <div class="toolbar">
                <div class="toolbar-left">
                    <div class="search-section">
                        <span class="search-label">关键词：</span>
                        <n-input v-model:value="searchKeyword" placeholder="输入模板名称或描述进行搜索..." clearable
                            style="width: 300px" @input="handleSearch" @keyup.enter="handleSearch">
                            <template #prefix>
                                <n-icon>
                                    <SearchOutline />
                                </n-icon>
                            </template>
                        </n-input>
                        <span class="search-label" style="margin-left: 16px">分类：</span>
                        <n-select v-model:value="selectedCategory" placeholder="选择分类" style="width: 160px"
                            :options="categoryOptions" clearable @update:value="handleSearch" />
                        <span class="search-label" style="margin-left: 16px">语言：</span>
                        <n-select v-model:value="selectedLanguage" placeholder="选择语言" style="width: 160px"
                            :options="languageOptions" clearable @update:value="handleSearch" />
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

            <n-data-table :columns="columns" :data="templates" :loading="loading" :pagination="pagination"
                :row-key="(row) => row.id" class="templates-table" remote @update:page="handlePageChange"
                @update:page-size="handlePageSizeChange" />
        </div>

        <!-- 添加/编辑模板弹窗 -->
        <n-modal v-model:show="showAddModal" :mask-closable="false">
            <n-card style="width: 600px" :title="editingTemplate ? '编辑模板' : '添加模板'" :bordered="false" size="huge"
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
                    <n-form-item label="模板名称" path="name">
                        <n-input v-model:value="formData.name" placeholder="请输入模板名称" :maxlength="100" show-count />
                    </n-form-item>

                    <n-form-item label="模板描述" path="description">
                        <n-input v-model:value="formData.description" type="textarea" placeholder="请输入模板描述"
                            :maxlength="500" show-count :rows="3" />
                    </n-form-item>

                    <n-form-item label="详细介绍" path="introduction">
                        <n-input v-model:value="formData.introduction" type="textarea"
                            placeholder="请输入详细介绍（支持Markdown格式）" :maxlength="2000" show-count :rows="4" />
                    </n-form-item>

                    <n-form-item label="所属分类" path="categoryId">
                        <n-select v-model:value="formData.categoryId" placeholder="请选择分类"
                            :options="categorySelectOptions" style="width: 100%" />
                    </n-form-item>

                    <n-form-item label="支持语言" path="languages">
                        <n-select v-model:value="formData.languages" placeholder="请选择支持的语言"
                            :options="languageSelectOptions" multiple style="width: 100%"
                            @update:value="onLanguagesChange" />
                    </n-form-item>

                    <n-form-item label="主语言" path="primaryLanguage">
                        <n-select v-model:value="formData.primaryLanguage" placeholder="请选择主语言"
                            :options="primaryLanguageOptions" style="width: 100%" />
                    </n-form-item>

                    <n-form-item label="模板图标" path="icon">
                        <n-input v-model:value="formData.icon" placeholder="请输入图标名称（可选）" />
                    </n-form-item>

                    <n-form-item label="推荐模板" path="isFeatured">
                        <n-switch v-model:value="formData.isFeatured" :checked-value="1" :unchecked-value="0">
                            <template #checked>是</template>
                            <template #unchecked>否</template>
                        </n-switch>
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="closeModal">取消</n-button>
                        <n-button type="primary" @click="handleSubmit" :loading="submitting">
                            {{ editingTemplate ? '更新' : '添加' }}
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
                        确定要删除模板 <strong>"{{ deletingTemplate?.name }}"</strong> 吗？
                    </p>
                    <p class="delete-warning">
                        此操作不可撤销，删除后相关模板文件和配置将无法恢复。
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
import { ref, reactive, computed, onMounted, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
    NButton, NIcon, NInput, NDataTable, NModal, NCard, NForm, NFormItem,
    NSelect, NSwitch, NTag, useMessage
} from 'naive-ui'
import {
    AddOutline, SearchOutline, RefreshOutline, CloseOutline,
    TrashOutline, CreateOutline, EyeOutline, Star, CodeOutline
} from '@vicons/ionicons5'
import { listTemplates, addTemplate, editTemplate, deleteTemplate } from '@/api/templates'
import { useCategoryStore } from '@/stores/categoryStore'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'

const router = useRouter()
const message = useMessage()

const categoryStore = useCategoryStore()
const { categoriesList } = storeToRefs(categoryStore)

const languageStore = useLanguageStore()
const { languagesList } = storeToRefs(languageStore)

// 数据状态
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const searching = ref(false)
const templates = ref([])
const searchKeyword = ref('')
const selectedCategory = ref(null)
const selectedLanguage = ref(null)

// 弹窗状态
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const editingTemplate = ref(null)
const deletingTemplate = ref(null)

// 表单数据
const formRef = ref(null)
const formData = reactive({
    name: '',
    description: '',
    introduction: '',
    categoryId: null,
    languages: [],
    primaryLanguage: null,
    icon: '',
    isFeatured: 0
})

// 表单验证规则
const formRules = {
    name: {
        required: true,
        message: '请输入模板名称',
        trigger: ['input', 'blur']
    },
    description: {
        required: true,
        message: '请输入模板描述',
        trigger: ['input', 'blur']
    },
    categoryId: {
        required: true,
        type: 'number',
        message: '请选择分类',
        trigger: 'change'
    },
    languages: {
        required: true,
        type: 'array',
        min: 1,
        message: '请选择支持的语言',
        trigger: 'change'
    },
    primaryLanguage: {
        required: true,
        type: 'number',
        message: '请选择主语言',
        trigger: 'change'
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

// 选项数据
const categoryOptions = computed(() => [
    { label: '全部分类', value: null },
    ...categoriesList.value.map(c => ({ label: c.name, value: Number(c.id) }))
])

const languageOptions = computed(() => [
    { label: '全部语言', value: null },
    ...languagesList.value.map(lang => ({ label: lang.name, value: Number(lang.id) }))
])

const categorySelectOptions = computed(() =>
    categoriesList.value.map(c => ({ label: c.name, value: Number(c.id) }))
)

const languageSelectOptions = computed(() =>
    languagesList.value.map(lang => ({ label: lang.name, value: Number(lang.id) }))
)

const primaryLanguageOptions = computed(() =>
    languagesList.value
        .filter(lang => formData.languages.includes(Number(lang.id)))
        .map(lang => ({
            label: lang.name,
            value: Number(lang.id)
        }))
)

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
        title: '模板名称',
        key: 'name',
        width: 200,
        render: (row) => h('div', { style: 'display: flex; align-items: center; gap: 8px' }, [
            row.isFeatured ? h(NIcon, { color: '#f0a020', size: 16 }, { default: () => h(Star) }) : null,
            h('span', { class: 'template-name' }, row.name)
        ])
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
        title: '分类',
        key: 'categoryId',
        width: 120,
        render: (row) => {
            const categoryName = getCategoryName(row.categoryId)
            return categoryName ? h(NTag, { type: 'info', size: 'small' }, { default: () => categoryName }) : '-'
        }
    },
    {
        title: '语言',
        key: 'languages',
        width: 200,
        render: (row) => {
            if (!row.languages || row.languages.length === 0) return '-'
            return h('div', { style: 'display: flex; flex-wrap: wrap; gap: 4px' },
                row.languages.slice(0, 3).map(lang => {
                    const isPrimary = lang.isPrimary === 1 || lang.is_primary === 1
                    if (isPrimary) {
                        // 主语言使用类似分类的样式
                        return h(NTag, {
                            type: 'info',
                            size: 'small'
                        }, {
                            default: () => getLanguageName(lang.languageId)
                        })
                    } else {
                        // 次语言保持彩色样式
                        return h(NTag, {
                            size: 'small',
                            color: {
                                color: '#fff',
                                backgroundColor: getLanguageColor(lang.languageId)
                            }
                        }, {
                            default: () => getLanguageName(lang.languageId)
                        })
                    }
                }).concat(
                    row.languages.length > 3 ? [h('span', { style: 'color: #999; font-size: 12px' }, `+${row.languages.length - 3}`)] : []
                )
            )
        }
    },
    {
        title: '推荐',
        key: 'isFeatured',
        width: 80,
        render: (row) => {
            return h(
                NSwitch,
                {
                    value: row.isFeatured === 1,
                    size: 'small',
                    onUpdateValue: (value) => handleFeaturedChange(row, value)
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
        width: 400,
        render: (row) => {
            return h('div', { class: 'action-buttons' }, [
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'info',
                        secondary: true,
                        onClick: () => handleView(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(EyeOutline) }),
                        default: () => '查看'
                    }
                ),
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'warning',
                        secondary: true,
                        style: { marginLeft: '8px' },
                        onClick: () => handleContentEdit(row)
                    },
                    {
                        icon: () => h(NIcon, null, { default: () => h(CodeOutline) }),
                        default: () => '内容编辑'
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
const loadTemplates = async (searchParams = {}) => {
    try {
        loading.value = true
        const params = {
            pageNum: pagination.page,
            pageSize: pagination.pageSize,
            ...searchParams
        }
        const response = await listTemplates(params)
        templates.value = response.data.data.templatesList || []
        pagination.itemCount = response.data.data.total || 0
    } catch (error) {
        console.error('获取模板列表失败:', error)
        message.error('获取模板列表失败')
    } finally {
        loading.value = false
    }
}

const refreshData = () => {
    loadTemplates()
}

const handleSearch = async () => {
    try {
        searching.value = true
        pagination.page = 1 // 重置到第一页

        const searchParams = {}
        if (searchKeyword.value) {
            searchParams.name = searchKeyword.value
        }
        if (selectedCategory.value) {
            searchParams.categoryId = selectedCategory.value
        }
        if (selectedLanguage.value) {
            searchParams.languageId = selectedLanguage.value
        }

        await loadTemplates(searchParams)
    } catch (error) {
        console.error('搜索失败:', error)
        message.error('搜索失败')
    } finally {
        searching.value = false
    }
}

const clearSearch = async () => {
    searchKeyword.value = ''
    selectedCategory.value = null
    selectedLanguage.value = null
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

const handleView = (template) => {
    router.push(`/templates/edit/${template.id}`)
}

const handleContentEdit = (template) => {
    // 在新标签页中打开模板编辑页面
    window.open(`/templates/edit/${template.id}`, '_blank')
}

const handleEdit = (template) => {
    editingTemplate.value = template
    formData.name = template.name
    formData.description = template.description
    formData.introduction = template.introduction || ''
    formData.categoryId = template.categoryId || template.category_id
    formData.isFeatured = template.isFeatured || 0
    formData.icon = template.icon || ''

    // 语言回显
    if (template.languages && template.languages.length > 0) {
        formData.languages = template.languages.map(l => Number(l.languageId || l.id))
        const primary = template.languages.find(l => l.isPrimary === 1 || l.is_primary === 1)
        formData.primaryLanguage = primary ? Number(primary.languageId || primary.id) : null
    } else {
        formData.languages = []
        formData.primaryLanguage = null
    }

    showAddModal.value = true
}

const handleDelete = (template) => {
    deletingTemplate.value = template
    showDeleteModal.value = true
}

const closeModal = () => {
    showAddModal.value = false
    editingTemplate.value = null
    resetForm()
}

const resetForm = () => {
    formData.name = ''
    formData.description = ''
    formData.introduction = ''
    formData.categoryId = null
    formData.languages = []
    formData.primaryLanguage = null
    formData.icon = ''
    formData.isFeatured = 0
    formRef.value?.restoreValidation()
}

const onLanguagesChange = (val) => {
    if (!val.includes(formData.primaryLanguage)) {
        formData.primaryLanguage = null
    }
}

const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        submitting.value = true

        const languagesArr = formData.languages.map(langId => ({
            languageId: langId,
            isPrimary: langId === formData.primaryLanguage ? 1 : 0
        }))

        const data = {
            name: formData.name,
            description: formData.description,
            introduction: formData.introduction,
            categoryId: formData.categoryId,
            isFeatured: formData.isFeatured,
            icon: formData.icon,
            languages: languagesArr
        }

        if (editingTemplate.value) {
            await editTemplate({ ...data, id: editingTemplate.value.id })
            message.success('模板更新成功')
        } else {
            await addTemplate(data)
            message.success('模板添加成功')
        }

        closeModal()
        loadTemplates()
    } catch (error) {
        console.error('操作失败:', error)
        message.error(editingTemplate.value ? '更新模板失败' : '添加模板失败')
    } finally {
        submitting.value = false
    }
}

const confirmDelete = async () => {
    try {
        deleting.value = true
        await deleteTemplate({ id: deletingTemplate.value.id })
        message.success('模板删除成功')
        showDeleteModal.value = false
        deletingTemplate.value = null
        loadTemplates()
    } catch (error) {
        console.error('删除模板失败:', error)
        message.error('删除模板失败')
    } finally {
        deleting.value = false
    }
}

const handleFeaturedChange = async (row, value) => {
    try {
        const newFeatured = value ? 1 : 0
        await editTemplate({
            id: row.id,
            name: row.name,
            description: row.description,
            introduction: row.introduction || '',
            categoryId: row.categoryId,
            isFeatured: newFeatured,
            icon: row.icon || '',
            languages: row.languages || []
        })

        const index = templates.value.findIndex(item => item.id === row.id)
        if (index !== -1) {
            templates.value[index].isFeatured = newFeatured
        }

        message.success(`模板已${value ? '设为推荐' : '取消推荐'}`)
    } catch (error) {
        console.error('推荐状态切换失败:', error)
        message.error('推荐状态切换失败')
        loadTemplates()
    }
}

const getCategoryName = (categoryId) => {
    if (!categoryId) return null
    const category = categoriesList.value.find(cat => cat.id === Number(categoryId))
    return category ? category.name : null
}

const getLanguageName = (languageId) => {
    if (!languageId) return ''
    const language = languagesList.value.find(lang => lang.id === Number(languageId))
    return language ? language.name : `未知语言(${languageId})`
}

const getLanguageColor = (languageId) => {
    if (!languageId) return '#666'
    const language = languagesList.value.find(lang => lang.id === Number(languageId))
    return language ? language.color : '#999'
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

watch(
    () => formData.languages,
    (langs) => {
        if (!langs.includes(formData.primaryLanguage)) {
            formData.primaryLanguage = null
        }
    }
)

// 生命周期
onMounted(async () => {
    try {
        // 确保分类和语言数据先加载完成
        await Promise.all([
            categoryStore.getCategories(),
            languageStore.getLanguages()
        ])
        // 然后加载模板数据
        await loadTemplates()
    } catch (error) {
        console.error('初始化数据加载失败:', error)
        message.error('初始化数据加载失败')
    }
})
</script>

<style scoped>
.templates-manage {
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
    flex-wrap: wrap;
}

.search-label {
    font-size: 14px;
    color: #333;
    font-weight: 500;
    white-space: nowrap;
}

.templates-table {
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

.template-name {
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