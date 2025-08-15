<template>
    <div class="preset-variable-manager" :class="{ 'compact-mode': compact }">
        <!-- 头部操作栏 -->
        <div class="manager-header">
            <div class="header-left">
                <h3 v-if="!compact">预设变量订阅</h3>
                <h4 v-else>预设变量</h4>
                <p v-if="!compact" class="description">从预设变量库中订阅变量到当前模板</p>
            </div>
            <div class="header-right">
                <n-button 
                    :type="compact ? 'default' : 'primary'" 
                    :size="compact ? 'small' : 'medium'"
                    @click="showSubscribeModal = true" 
                    :loading="loading"
                >
                    <template #icon>
                        <n-icon><AddOutline /></n-icon>
                    </template>
                    {{ compact ? '订阅' : '订阅预设变量' }}
                </n-button>
            </div>
        </div>

        <!-- 已订阅的预设变量列表 -->
        <div class="subscribed-list" v-loading="loading">
            <n-empty v-if="subscribedList.length === 0" description="暂无订阅的预设变量">
                <template #extra>
                    <n-button size="small" @click="showSubscribeModal = true">立即订阅</n-button>
                </template>
            </n-empty>
            
            <div v-else>
                <div class="list-item" v-for="item in subscribedList" :key="item.id">
                    <div class="item-header">
                        <div class="item-title">
                            <span class="preset-name">{{ item.preset_name }}</span>
                            <n-tag size="tiny" type="info">{{ item.preset_path }}</n-tag>
                        </div>
                        <div class="item-actions">
                            <n-switch 
                                v-model:value="item.is_active" 
                                :loading="item.updating"
                                @update:value="updateActiveStatus(item)"
                                size="small"
                            />
                            <n-button 
                                size="tiny" 
                                quaternary 
                                @click="editMapping(item)"
                                :loading="item.editing"
                            >
                                编辑
                            </n-button>
                            <n-popconfirm @positive-click="unsubscribe(item)">
                                <template #trigger>
                                    <n-button 
                                        size="tiny" 
                                        quaternary 
                                        type="error"
                                        :loading="item.unsubscribing"
                                    >
                                        删除
                                    </n-button>
                                </template>
                                确定要取消订阅这个预设变量吗？
                            </n-popconfirm>
                        </div>
                    </div>
                    
                    <div class="item-content">
                        <div class="mapping-info">
                            <span class="label">映射为:</span>
                            <span class="value">{{ item.mapped_name }}</span>
                        </div>
                        <div class="variable-info" v-if="item.display_name || item.description">
                            <div class="info-row" v-if="item.display_name">
                                <span class="label">显示名:</span>
                                <span class="value">{{ item.display_name }}</span>
                            </div>
                            <div class="info-row" v-if="item.description">
                                <span class="label">描述:</span>
                                <span class="value">{{ item.description }}</span>
                            </div>
                            <div class="info-row" v-if="item.example">
                                <span class="label">示例:</span>
                                <span class="value">{{ item.example }}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 订阅预设变量弹窗 -->
        <n-modal v-model:show="showSubscribeModal" :mask-closable="false">
            <n-card style="width: 800px" title="订阅预设变量" :bordered="false" size="huge">
                <template #header-extra>
                    <n-button quaternary circle @click="showSubscribeModal = false">
                        <template #icon>
                            <n-icon><CloseOutline /></n-icon>
                        </template>
                    </n-button>
                </template>

                <!-- 搜索栏 -->
                <div class="search-bar">
                    <n-input 
                        v-model:value="searchKeyword" 
                        placeholder="搜索预设变量..." 
                        clearable
                        @update:value="searchPresets"
                    >
                        <template #prefix>
                            <n-icon><SearchOutline /></n-icon>
                        </template>
                    </n-input>
                </div>

                <!-- 可用预设变量列表 -->
                <div class="available-presets" v-loading="presetsLoading">
                    <n-empty v-if="availablePresets.length === 0" description="暂无可用的预设变量" />
                    
                    <div v-else class="presets-list">
                        <div 
                            class="preset-item" 
                            v-for="preset in availablePresets" 
                            :key="preset.id"
                        >
                            <div class="preset-header">
                                <h4>{{ preset.name }}</h4>
                                <p class="preset-description">{{ preset.description }}</p>
                            </div>
                            
                            <div class="variables-list">
                                <div 
                                    class="variable-item"
                                    v-for="variable in preset.variables"
                                    :key="variable.path"
                                >
                                    <n-checkbox 
                                        :checked="isVariableSelected(preset.id, variable.path)"
                                        @update:checked="toggleVariable(preset.id, variable, $event)"
                                    >
                                        <div class="variable-info">
                                            <div class="variable-name">
                                                <span class="path">{{ variable.path }}</span>
                                                <span class="display-name" v-if="variable.display_name">
                                                    ({{ variable.display_name }})
                                                </span>
                                            </div>
                                            <div class="variable-desc" v-if="variable.description">
                                                {{ variable.description }}
                                            </div>
                                        </div>
                                    </n-checkbox>
                                    
                                    <!-- 映射设置 -->
                                    <div 
                                        class="mapping-config" 
                                        v-if="isVariableSelected(preset.id, variable.path)"
                                    >
                                        <n-input 
                                            :value="getVariableMapping(preset.id, variable.path)"
                                            @update:value="updateVariableMapping(preset.id, variable.path, $event)"
                                            placeholder="映射到的模板变量名"
                                            size="small"
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 分页 -->
                <div class="pagination" v-if="totalPresets > pageSize">
                    <n-pagination 
                        v-model:page="currentPage"
                        :page-size="pageSize"
                        :item-count="totalPresets"
                        @update:page="loadAvailablePresets"
                        show-size-picker
                        :page-sizes="[10, 20, 30]"
                        @update:page-size="handlePageSizeChange"
                    />
                </div>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="showSubscribeModal = false">取消</n-button>
                        <n-button 
                            type="primary" 
                            @click="confirmSubscribe" 
                            :loading="subscribing"
                            :disabled="selectedVariables.length === 0"
                        >
                            订阅选中的变量 ({{ selectedVariables.length }})
                        </n-button>
                    </div>
                </template>
            </n-card>
        </n-modal>

        <!-- 编辑映射弹窗 -->
        <n-modal v-model:show="showEditModal" :mask-closable="false">
            <n-card style="width: 500px" title="编辑变量映射" :bordered="false" size="huge">
                <template #header-extra>
                    <n-button quaternary circle @click="showEditModal = false">
                        <template #icon>
                            <n-icon><CloseOutline /></n-icon>
                        </template>
                    </n-button>
                </template>

                <n-form ref="editFormRef" :model="editForm" :rules="editRules" label-placement="top">
                    <n-form-item label="预设变量路径" path="preset_path">
                        <n-input v-model:value="editForm.preset_path" readonly />
                    </n-form-item>
                    
                    <n-form-item label="映射到的模板变量名" path="mapped_name">
                        <n-input v-model:value="editForm.mapped_name" placeholder="请输入映射的变量名" />
                    </n-form-item>
                    
                    <n-form-item label="是否启用" path="is_active">
                        <n-switch v-model:value="editForm.is_active" />
                    </n-form-item>
                </n-form>

                <template #footer>
                    <div class="modal-footer">
                        <n-button @click="showEditModal = false">取消</n-button>
                        <n-button type="primary" @click="confirmEdit" :loading="editSubmitting">
                            保存
                        </n-button>
                    </div>
                </template>
            </n-card>
        </n-modal>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { 
    NButton, NIcon, NModal, NCard, NForm, NFormItem, NInput, NSwitch, 
    NTag, NEmpty, NCheckbox, NPagination, NPopconfirm, useMessage
} from 'naive-ui'
import { 
    AddOutline, CloseOutline, SearchOutline 
} from '@vicons/ionicons5'
import request from '@/utils/request'

// API函数
const subscribePreset = (templateId, presets) => {
  return request({
    url: `/api/v1/templates/${templateId}/preset-variables/subscribe`,
    method: 'POST',
    data: { template_id: templateId, presets: presets }
  })
}

const getSubscribedPresets = (templateId) => {
  return request({
    url: `/api/v1/templates/${templateId}/preset-variables`,
    method: 'GET'
  })
}

const unsubscribePreset = (templateId, id) => {
  return request({
    url: `/api/v1/templates/${templateId}/preset-variables/${id}`,
    method: 'DELETE'
  })
}

const updatePresetMapping = (templateId, id, data) => {
  return request({
    url: `/api/v1/templates/${templateId}/preset-variables/${id}`,
    method: 'PUT',
    data: { template_id: templateId, id: id, ...data }
  })
}

const getAvailablePresets = (params = {}) => {
  return request({
    url: '/api/v1/templates/preset-variables/available',
    method: 'GET',
    params: { pageNum: 1, pageSize: 20, ...params }
  })
}

const props = defineProps({
    templateId: {
        type: [Number, String],
        required: true
    },
    compact: {
        type: Boolean,
        default: false
    }
})

const message = useMessage()

// 数据状态
const loading = ref(false)
const subscribedList = ref([])
const showSubscribeModal = ref(false)
const showEditModal = ref(false)

// 可用预设变量相关
const presetsLoading = ref(false)
const availablePresets = ref([])
const totalPresets = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')

// 选中的变量
const selectedVariables = ref([])

// 订阅状态
const subscribing = ref(false)

// 编辑表单
const editFormRef = ref(null)
const editForm = reactive({
    id: null,
    preset_path: '',
    mapped_name: '',
    is_active: true
})

const editRules = {
    mapped_name: {
        required: true,
        message: '请输入映射的变量名',
        trigger: ['input', 'blur']
    }
}

const editSubmitting = ref(false)

// 加载已订阅的预设变量
const loadSubscribedPresets = async () => {
    loading.value = true
    try {
        const response = await getSubscribedPresets(props.templateId)
        subscribedList.value = response.data || []
    } catch (error) {
        console.error('加载订阅列表失败:', error)
        message.error('加载订阅列表失败')
    } finally {
        loading.value = false
    }
}

// 加载可用预设变量
const loadAvailablePresets = async (page = 1) => {
    presetsLoading.value = true
    try {
        const response = await getAvailablePresets({
            pageNum: page,
            pageSize: pageSize.value,
            keyword: searchKeyword.value
        })
        
        const data = response.data.data || {}
        console.log('API返回数据:', response.data)
        console.log('解析后数据:', data)
        console.log('预设列表:', data.list)
        availablePresets.value = data.list || []
        totalPresets.value = data.total || 0
        currentPage.value = data.pageNum || 1
        console.log('设置后的availablePresets:', availablePresets.value)
    } catch (error) {
        console.error('加载可用预设变量失败:', error)
        message.error('加载可用预设变量失败')
    } finally {
        presetsLoading.value = false
    }
}

// 搜索预设变量
let searchTimeout = null
const searchPresets = () => {
    if (searchTimeout) {
        clearTimeout(searchTimeout)
    }
    searchTimeout = setTimeout(() => {
        currentPage.value = 1
        loadAvailablePresets()
    }, 300)
}

// 分页处理
const handlePageSizeChange = (newSize) => {
    pageSize.value = newSize
    currentPage.value = 1
    loadAvailablePresets()
}

// 变量选择相关
const isVariableSelected = (presetId, variablePath) => {
    return selectedVariables.value.some(v => 
        v.var_preset_id === presetId && v.preset_path === variablePath
    )
}

const toggleVariable = (presetId, variable, checked) => {
    if (checked) {
        // 添加选中的变量
        selectedVariables.value.push({
            var_preset_id: presetId,
            preset_path: variable.path,
            mapped_name: variable.path.split('.').pop(), // 默认使用路径的最后一部分作为映射名
            is_active: 1,
            sort: selectedVariables.value.length
        })
    } else {
        // 移除选中的变量
        const index = selectedVariables.value.findIndex(v => 
            v.var_preset_id === presetId && v.preset_path === variable.path
        )
        if (index > -1) {
            selectedVariables.value.splice(index, 1)
        }
    }
}

const getVariableMapping = (presetId, variablePath) => {
    const variable = selectedVariables.value.find(v => 
        v.var_preset_id === presetId && v.preset_path === variablePath
    )
    return variable?.mapped_name || ''
}

const updateVariableMapping = (presetId, variablePath, mappedName) => {
    const variable = selectedVariables.value.find(v => 
        v.var_preset_id === presetId && v.preset_path === variablePath
    )
    if (variable) {
        variable.mapped_name = mappedName
    }
}

// 确认订阅
const confirmSubscribe = async () => {
    if (selectedVariables.value.length === 0) {
        message.warning('请选择要订阅的变量')
        return
    }

    subscribing.value = true
    try {
        await subscribePreset(props.templateId, selectedVariables.value)
        message.success('订阅成功')
        showSubscribeModal.value = false
        selectedVariables.value = []
        await loadSubscribedPresets()
    } catch (error) {
        console.error('订阅失败:', error)
        message.error('订阅失败')
    } finally {
        subscribing.value = false
    }
}

// 更新激活状态
const updateActiveStatus = async (item) => {
    item.updating = true
    try {
        await updatePresetMapping(props.templateId, item.id, {
            mapped_name: item.mapped_name,
            is_active: item.is_active ? 1 : 0,
            sort: item.sort
        })
        message.success('更新成功')
    } catch (error) {
        console.error('更新失败:', error)
        message.error('更新失败')
        // 恢复原状态
        item.is_active = !item.is_active
    } finally {
        item.updating = false
    }
}

// 编辑映射
const editMapping = (item) => {
    editForm.id = item.id
    editForm.preset_path = item.preset_path
    editForm.mapped_name = item.mapped_name
    editForm.is_active = item.is_active === 1
    showEditModal.value = true
}

// 确认编辑
const confirmEdit = async () => {
    try {
        await editFormRef.value?.validate()
        editSubmitting.value = true
        
        await updatePresetMapping(props.templateId, editForm.id, {
            mapped_name: editForm.mapped_name,
            is_active: editForm.is_active ? 1 : 0,
            sort: 0
        })
        
        message.success('更新成功')
        showEditModal.value = false
        await loadSubscribedPresets()
    } catch (error) {
        console.error('更新失败:', error)
        if (error.message) {
            message.error(error.message)
        } else {
            message.error('更新失败')
        }
    } finally {
        editSubmitting.value = false
    }
}

// 取消订阅
const unsubscribe = async (item) => {
    item.unsubscribing = true
    try {
        await unsubscribePreset(props.templateId, item.id)
        message.success('取消订阅成功')
        await loadSubscribedPresets()
    } catch (error) {
        console.error('取消订阅失败:', error)
        message.error('取消订阅失败')
    } finally {
        item.unsubscribing = false
    }
}

// 初始化
onMounted(() => {
    loadSubscribedPresets()
})

// 监听弹窗打开，加载可用预设变量
const watchSubscribeModal = () => {
    if (showSubscribeModal.value) {
        selectedVariables.value = []
        currentPage.value = 1
        searchKeyword.value = ''
        loadAvailablePresets()
    }
}

// 监听弹窗状态
watch(() => showSubscribeModal.value, watchSubscribeModal)
</script>

<style scoped>
.preset-variable-manager {
    padding: 16px;
}

.manager-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
}

.header-left h3 {
    margin: 0 0 4px 0;
    color: #333;
}

.header-left .description {
    margin: 0;
    color: #666;
    font-size: 14px;
}

.subscribed-list .list-item {
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 12px;
    background: #fff;
}

.item-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.item-title {
    display: flex;
    align-items: center;
    gap: 8px;
}

.preset-name {
    font-weight: 500;
    color: #333;
}

.item-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.item-content {
    padding-left: 0;
}

.mapping-info, .info-row {
    display: flex;
    margin-bottom: 4px;
}

.mapping-info .label, .info-row .label {
    width: 60px;
    color: #666;
    font-size: 12px;
}

.mapping-info .value, .info-row .value {
    color: #333;
    font-size: 12px;
}

.search-bar {
    margin-bottom: 16px;
}

.available-presets {
    max-height: 400px;
    overflow-y: auto;
}

.preset-item {
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
}

.preset-header h4 {
    margin: 0 0 8px 0;
    color: #333;
}

.preset-description {
    margin: 0 0 16px 0;
    color: #666;
    font-size: 14px;
}

.variable-item {
    padding: 8px 0;
    border-bottom: 1px solid #f0f0f0;
}

.variable-item:last-child {
    border-bottom: none;
}

.variable-info {
    margin-left: 24px;
}

.variable-name {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
}

.path {
    font-family: monospace;
    color: #333;
    background: #f5f5f5;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 12px;
}

.display-name {
    color: #666;
    font-size: 12px;
}

.variable-desc {
    color: #999;
    font-size: 12px;
    line-height: 1.4;
}

.mapping-config {
    margin: 8px 0 0 24px;
    max-width: 200px;
}

.pagination {
    margin-top: 16px;
    text-align: center;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}

/* 紧凑模式样式 - 针对在变量面板中的使用 */
.compact-mode .preset-variable-manager {
    padding: 8px;
}

.compact-mode .manager-header {
    margin-bottom: 12px;
}

.compact-mode .manager-header h4 {
    font-size: 14px;
    margin-bottom: 2px;
}

.compact-mode .subscribed-list .list-item {
    padding: 8px;
    margin-bottom: 8px;
    border-radius: 6px;
}

.compact-mode .item-header {
    margin-bottom: 6px;
}

.compact-mode .preset-name {
    font-size: 13px;
}

.compact-mode .item-actions .n-button {
    padding: 2px 6px;
    font-size: 11px;
    height: 24px;
}

.compact-mode .item-actions .n-switch {
    transform: scale(0.85);
}

.compact-mode .mapping-info .label,
.compact-mode .info-row .label {
    width: 50px;
    font-size: 11px;
}

.compact-mode .mapping-info .value,
.compact-mode .info-row .value {
    font-size: 11px;
}

.compact-mode .n-tag {
    font-size: 10px;
    padding: 2px 6px;
}

.compact-mode .header-right .n-button {
    padding: 4px 8px;
    font-size: 12px;
    height: 28px;
}
</style>