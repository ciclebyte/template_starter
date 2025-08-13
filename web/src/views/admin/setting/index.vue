<template>
    <div class="settings">
        <n-spin :show="loading">
            <!-- 页面标题 -->
            <div class="page-header">
                <h1>系统设置</h1>
                <p>管理系统配置、AI功能、界面偏好等设置</p>
            </div>

            <!-- 设置分类导航 -->
            <n-grid :cols="5" :x-gap="20">
                <!-- 左侧分类菜单 -->
                <n-gi :span="1">
                    <n-card>
                        <n-menu :value="activeTab" :options="tabOptions" @update:value="handleTabChange"
                            class="settings-menu" />
                    </n-card>
                </n-gi>

                <!-- 右侧配置内容 -->
                <n-gi :span="4">
                    <n-card>
                        <!-- AI配置 -->
                        <div v-if="activeTab === 'ai'" class="config-section">
                            <div class="section-header">
                                <h2>AI功能配置</h2>
                                <p>配置人工智能相关功能和服务商设置</p>
                                <n-divider />
                            </div>

                            <n-form ref="aiFormRef" :model="aiConfig" :rules="aiRules" label-placement="left"
                                label-width="140" class="config-form">
                                <!-- 基础开关 -->
                                <n-form-item label="启用AI功能" path="enabled">
                                    <n-switch v-model:value="aiConfig.enabled" />
                                    <span class="form-tip">开启后可使用AI辅助模板生成等功能</span>
                                </n-form-item>

                                <n-form-item label="AI服务商" path="provider">
                                    <n-select v-model:value="aiConfig.provider" :options="providerOptions"
                                        :disabled="!aiConfig.enabled" />
                                </n-form-item>

                                <!-- OpenAI配置 -->
                                <template v-if="aiConfig.provider === 'openai'">
                                    <n-divider title-placement="left">OpenAI 配置</n-divider>

                                    <n-form-item label="API Key" path="openai.apiKey">
                                        <n-input v-model:value="aiConfig.openai.apiKey" type="password"
                                            placeholder="sk-your-openai-api-key" show-password-on="click"
                                            :disabled="!aiConfig.enabled" />
                                    </n-form-item>

                                    <n-form-item label="API地址" path="openai.baseUrl">
                                        <n-input v-model:value="aiConfig.openai.baseUrl"
                                            placeholder="https://api.openai.com" :disabled="!aiConfig.enabled" />
                                    </n-form-item>

                                    <n-form-item label="模型" path="openai.model">
                                        <n-select v-model:value="aiConfig.openai.model" :options="openaiModelOptions"
                                            :disabled="!aiConfig.enabled" filterable tag placeholder="选择或输入模型名称" />
                                    </n-form-item>

                                    <n-form-item label="最大Token数" path="openai.maxTokens">
                                        <n-input-number v-model:value="aiConfig.openai.maxTokens" :min="100" :max="8000"
                                            :disabled="!aiConfig.enabled" />
                                    </n-form-item>

                                    <n-form-item label="创造性参数" path="openai.temperature">
                                        <n-slider v-model:value="aiConfig.openai.temperature" :min="0" :max="1"
                                            :step="0.1" :disabled="!aiConfig.enabled" />
                                        <span class="form-tip">控制AI回复的创造性，0为保守，1为创新</span>
                                    </n-form-item>
                                </template>

                                <!-- Claude配置 -->
                                <template v-if="aiConfig.provider === 'claude'">
                                    <n-divider title-placement="left">Claude 配置</n-divider>

                                    <n-form-item label="API Key" path="claude.apiKey">
                                        <n-input v-model:value="aiConfig.claude.apiKey" type="password"
                                            placeholder="your-claude-api-key" show-password-on="click"
                                            :disabled="!aiConfig.enabled" />
                                    </n-form-item>

                                    <n-form-item label="API地址" path="claude.baseUrl">
                                        <n-input v-model:value="aiConfig.claude.baseUrl"
                                            placeholder="https://api.anthropic.com" :disabled="!aiConfig.enabled" />
                                    </n-form-item>

                                    <n-form-item label="模型" path="claude.model">
                                        <n-select v-model:value="aiConfig.claude.model" :options="claudeModelOptions"
                                            :disabled="!aiConfig.enabled" filterable tag placeholder="选择或输入模型名称" />
                                    </n-form-item>
                                </template>

                                <!-- 功能开关 -->
                                <n-divider title-placement="left">功能模块</n-divider>

                                <n-form-item label="模板生成">
                                    <n-switch v-model:value="aiConfig.features.templateGeneration"
                                        :disabled="!aiConfig.enabled" />
                                    <span class="form-tip">AI辅助生成项目模板</span>
                                </n-form-item>

                                <n-form-item label="代码审查">
                                    <n-switch v-model:value="aiConfig.features.codeReview"
                                        :disabled="!aiConfig.enabled" />
                                    <span class="form-tip">AI代码质量检查和建议</span>
                                </n-form-item>

                                <n-form-item label="变量建议">
                                    <n-switch v-model:value="aiConfig.features.variableSuggestion"
                                        :disabled="!aiConfig.enabled" />
                                    <span class="form-tip">智能推荐模板变量</span>
                                </n-form-item>

                                <n-form-item label="自动文档">
                                    <n-switch v-model:value="aiConfig.features.autoDocumentation"
                                        :disabled="!aiConfig.enabled" />
                                    <span class="form-tip">AI自动生成项目文档</span>
                                </n-form-item>

                                <!-- 使用限制 -->
                                <n-divider title-placement="left">使用限制</n-divider>

                                <n-form-item label="每小时请求数">
                                    <n-input-number v-model:value="aiConfig.rateLimit.requestsPerHour" :min="1"
                                        :max="1000" :disabled="!aiConfig.enabled" />
                                </n-form-item>

                                <n-form-item label="每日请求数">
                                    <n-input-number v-model:value="aiConfig.rateLimit.requestsPerDay" :min="10"
                                        :max="10000" :disabled="!aiConfig.enabled" />
                                </n-form-item>
                            </n-form>
                        </div>

                        <!-- 系统配置 -->
                        <div v-if="activeTab === 'system'" class="config-section">
                            <div class="section-header">
                                <h2>系统基础配置</h2>
                                <p>配置系统基本信息和显示设置</p>
                                <n-divider />
                            </div>

                            <n-form ref="systemFormRef" :model="systemConfig" :rules="systemRules"
                                label-placement="left" label-width="120" class="config-form">
                                <n-form-item label="系统名称" path="name">
                                    <n-input v-model:value="systemConfig.name" placeholder="Template Starter" />
                                </n-form-item>

                                <n-form-item label="系统描述" path="description">
                                    <n-input v-model:value="systemConfig.description" type="textarea"
                                        placeholder="智能代码模板生成平台" :rows="3" />
                                </n-form-item>

                                <n-form-item label="系统版本" path="version">
                                    <n-input v-model:value="systemConfig.version" placeholder="1.0.0" />
                                </n-form-item>

                                <n-form-item label="Logo地址" path="logoUrl">
                                    <n-input v-model:value="systemConfig.logoUrl" placeholder="/logo.png" />
                                </n-form-item>
                            </n-form>
                        </div>

                        <!-- 模板配置 -->
                        <div v-if="activeTab === 'template'" class="config-section">
                            <div class="section-header">
                                <h2>模板相关配置</h2>
                                <p>设置模板生成的默认参数和限制</p>
                                <n-divider />
                            </div>

                            <n-form ref="templateFormRef" :model="templateConfig" :rules="templateRules"
                                label-placement="left" label-width="140" class="config-form">
                                <n-form-item label="默认作者" path="defaultAuthor">
                                    <n-input v-model:value="templateConfig.defaultAuthor" placeholder="系统管理员" />
                                </n-form-item>

                                <n-form-item label="默认公司" path="defaultCompany">
                                    <n-input v-model:value="templateConfig.defaultCompany" placeholder="我的公司" />
                                </n-form-item>

                                <n-form-item label="文件大小限制" path="maxFileSize">
                                    <n-input-number v-model:value="templateConfig.maxFileSize" :min="1" :max="100"
                                        suffix="MB" />
                                </n-form-item>

                                <n-form-item label="文件数量限制" path="maxFilesPerTemplate">
                                    <n-input-number v-model:value="templateConfig.maxFilesPerTemplate" :min="10"
                                        :max="1000" suffix="个" />
                                </n-form-item>
                            </n-form>
                        </div>

                        <!-- UI配置 -->
                        <div v-if="activeTab === 'ui'" class="config-section">
                            <div class="section-header">
                                <h2>界面配置</h2>
                                <p>自定义用户界面和体验设置</p>
                                <n-divider />
                            </div>

                            <n-form ref="uiFormRef" :model="uiConfig" :rules="uiRules" label-placement="left"
                                label-width="120" class="config-form">
                                <n-form-item label="默认主题" path="theme">
                                    <n-select v-model:value="uiConfig.theme" :options="themeOptions" />
                                </n-form-item>

                                <n-form-item label="默认语言" path="language">
                                    <n-select v-model:value="uiConfig.language" :options="languageOptions" />
                                </n-form-item>

                                <n-form-item label="分页大小" path="pageSize">
                                    <n-input-number v-model:value="uiConfig.pageSize" :min="10" :max="100" />
                                </n-form-item>

                                <n-form-item label="编辑器主题" path="editorTheme">
                                    <n-select v-model:value="uiConfig.editorTheme" :options="editorThemeOptions" />
                                </n-form-item>

                                <n-form-item label="编辑器字体大小" path="editorFontSize">
                                    <n-input-number v-model:value="uiConfig.editorFontSize" :min="12" :max="24"
                                        suffix="px" />
                                </n-form-item>
                            </n-form>
                        </div>

                        <!-- 操作按钮 -->
                        <div class="config-actions">
                            <n-space>
                                <n-button type="primary" @click="saveConfig" :loading="saving">
                                    <template #icon>
                                        <n-icon>
                                            <SaveOutline />
                                        </n-icon>
                                    </template>
                                    保存配置
                                </n-button>
                                <n-button @click="resetConfig">
                                    <template #icon>
                                        <n-icon>
                                            <RefreshOutline />
                                        </n-icon>
                                    </template>
                                    重置
                                </n-button>
                                <n-button @click="testConnection" v-if="activeTab === 'ai'" :loading="testing">
                                    <template #icon>
                                        <n-icon>
                                            <CheckmarkOutline />
                                        </n-icon>
                                    </template>
                                    测试连接
                                </n-button>
                            </n-space>
                        </div>
                    </n-card>
                </n-gi>
            </n-grid>
        </n-spin>
    </div>
</template>

<script setup>
import { ref, onMounted, watch, h } from 'vue'
import {
    NGrid,
    NGi,
    NCard,
    NMenu,
    NForm,
    NFormItem,
    NInput,
    NInputNumber,
    NSelect,
    NSwitch,
    NSlider,
    NButton,
    NSpace,
    NIcon,
    NDivider,
    NSpin,
    useMessage
} from 'naive-ui'
import {
    SaveOutline,
    RefreshOutline,
    CheckmarkOutline,
    SettingsOutline,
    ColorPaletteOutline,
    DocumentTextOutline,
    CodeOutline,
    CloudOutline
} from '@vicons/ionicons5'
import {
    getAIConfig,
    getSystemConfig,
    getTemplateConfig,
    getUIConfig,
    batchUpdateConfig
} from '@/api/systemConfig'

const message = useMessage()

// 数据状态
const loading = ref(true)
const saving = ref(false)
const testing = ref(false)
const activeTab = ref('system')

// 表单引用
const aiFormRef = ref(null)
const systemFormRef = ref(null)
const templateFormRef = ref(null)
const uiFormRef = ref(null)

// 配置数据
const aiConfig = ref({
    enabled: false,
    provider: 'openai',
    openai: {
        apiKey: '',
        baseUrl: 'https://api.openai.com',
        model: 'gpt-4',
        maxTokens: 4000,
        temperature: 0.7
    },
    claude: {
        apiKey: '',
        baseUrl: 'https://api.anthropic.com',
        model: 'claude-3-sonnet-20240229'
    },
    features: {
        templateGeneration: true,
        codeReview: false,
        variableSuggestion: true,
        autoDocumentation: false
    },
    rateLimit: {
        requestsPerHour: 100,
        requestsPerDay: 1000
    }
})

const systemConfig = ref({
    name: 'Template Starter',
    description: '智能代码模板生成平台',
    version: '1.0.0',
    logoUrl: '/logo.png'
})

const templateConfig = ref({
    defaultAuthor: '系统管理员',
    defaultCompany: '',
    maxFileSize: 10,
    maxFilesPerTemplate: 500
})

const uiConfig = ref({
    theme: 'light',
    language: 'zh-CN',
    pageSize: 20,
    editorTheme: 'vs-dark',
    editorFontSize: 14
})

// 选项数据
const tabOptions = [
    { label: '系统配置', key: 'system', icon: () => h(NIcon, null, { default: () => h(SettingsOutline) }) },
    { label: 'AI配置', key: 'ai', icon: () => h(NIcon, null, { default: () => h(CloudOutline) }) },
    { label: '模板配置', key: 'template', icon: () => h(NIcon, null, { default: () => h(DocumentTextOutline) }) },
    { label: '界面配置', key: 'ui', icon: () => h(NIcon, null, { default: () => h(ColorPaletteOutline) }) }
]

const providerOptions = [
    { label: 'OpenAI', value: 'openai' },
    { label: 'Claude (Anthropic)', value: 'claude' },
    { label: '自定义', value: 'custom' }
]

const openaiModelOptions = [
    { label: 'GPT-4', value: 'gpt-4' },
    { label: 'GPT-4 Turbo', value: 'gpt-4-turbo-preview' },
    { label: 'GPT-3.5 Turbo', value: 'gpt-3.5-turbo' }
]

const claudeModelOptions = [
    { label: 'Claude 3 Sonnet', value: 'claude-3-sonnet-20240229' },
    { label: 'Claude 3 Haiku', value: 'claude-3-haiku-20240307' },
    { label: 'Claude 3 Opus', value: 'claude-3-opus-20240229' }
]

const themeOptions = [
    { label: '浅色主题', value: 'light' },
    { label: '深色主题', value: 'dark' },
    { label: '跟随系统', value: 'auto' }
]

const languageOptions = [
    { label: '简体中文', value: 'zh-CN' },
    { label: 'English', value: 'en-US' }
]

const editorThemeOptions = [
    { label: 'VS Code Dark', value: 'vs-dark' },
    { label: 'VS Code Light', value: 'vs-light' },
    { label: 'Dracula', value: 'dracula' }
]

// 验证规则
const aiRules = {
    'openai.apiKey': [
        { required: true, message: '请输入OpenAI API Key', trigger: 'blur' }
    ],
    'claude.apiKey': [
        { required: true, message: '请输入Claude API Key', trigger: 'blur' }
    ]
}

const systemRules = {
    name: [
        { required: true, message: '请输入系统名称', trigger: 'blur' }
    ]
}

const templateRules = {}
const uiRules = {}

// 方法
async function loadConfigs() {
    try {
        loading.value = true

        const [aiRes, systemRes, templateRes, uiRes] = await Promise.all([
            getAIConfig(),
            getSystemConfig(),
            getTemplateConfig(),
            getUIConfig()
        ])

        // 处理AI配置
        if (aiRes.data.data.configs) {
            const configs = aiRes.data.data.configs
            aiConfig.value = {
                enabled: configs['ai.enabled'] || false,
                provider: configs['ai.provider'] || 'openai',
                openai: {
                    apiKey: configs['ai.openai.api_key'] || '',
                    baseUrl: configs['ai.openai.base_url'] || 'https://api.openai.com',
                    model: configs['ai.openai.model'] || 'gpt-4',
                    maxTokens: configs['ai.openai.max_tokens'] || 4000,
                    temperature: parseFloat(configs['ai.openai.temperature']) || 0.7
                },
                claude: {
                    apiKey: configs['ai.claude.api_key'] || '',
                    baseUrl: configs['ai.claude.base_url'] || 'https://api.anthropic.com',
                    model: configs['ai.claude.model'] || 'claude-3-sonnet-20240229'
                },
                features: {
                    templateGeneration: configs['ai.template_generation.enabled'] || true,
                    codeReview: configs['ai.code_review.enabled'] || false,
                    variableSuggestion: configs['ai.variable_suggestion.enabled'] || true,
                    autoDocumentation: configs['ai.auto_documentation.enabled'] || false
                },
                rateLimit: {
                    requestsPerHour: configs['ai.rate_limit.requests_per_hour'] || 100,
                    requestsPerDay: configs['ai.rate_limit.requests_per_day'] || 1000
                }
            }
        }

        // 处理其他配置
        if (systemRes.data.data.configs) {
            const configs = systemRes.data.data.configs
            systemConfig.value = {
                name: configs['system.name'] || 'Template Starter',
                description: configs['system.description'] || '智能代码模板生成平台',
                version: configs['system.version'] || '1.0.0',
                logoUrl: configs['system.logo_url'] || '/logo.png'
            }
        }

        if (templateRes.data.data.configs) {
            const configs = templateRes.data.data.configs
            templateConfig.value = {
                defaultAuthor: configs['template.default_author'] || '系统管理员',
                defaultCompany: configs['template.default_company'] || '',
                maxFileSize: configs['template.max_file_size'] || 10,
                maxFilesPerTemplate: configs['template.max_files_per_template'] || 500
            }
        }

        if (uiRes.data.data.configs) {
            const configs = uiRes.data.data.configs
            uiConfig.value = {
                theme: configs['ui.theme'] || 'light',
                language: configs['ui.language'] || 'zh-CN',
                pageSize: configs['ui.page_size'] || 20,
                editorTheme: configs['ui.editor_theme'] || 'vs-dark',
                editorFontSize: configs['ui.editor_font_size'] || 14
            }
        }

    } catch (error) {
        console.error('加载配置失败:', error)
        message.error('加载配置失败')
    } finally {
        loading.value = false
    }
}

function handleTabChange(key) {
    activeTab.value = key
}

async function saveConfig() {
    try {
        saving.value = true

        let configs = []

        if (activeTab.value === 'ai') {
            configs = [
                { configKey: 'ai.enabled', configValue: String(aiConfig.value.enabled) },
                { configKey: 'ai.provider', configValue: aiConfig.value.provider },
                { configKey: 'ai.openai.api_key', configValue: aiConfig.value.openai.apiKey },
                { configKey: 'ai.openai.base_url', configValue: aiConfig.value.openai.baseUrl },
                { configKey: 'ai.openai.model', configValue: aiConfig.value.openai.model },
                { configKey: 'ai.openai.max_tokens', configValue: String(aiConfig.value.openai.maxTokens) },
                { configKey: 'ai.openai.temperature', configValue: String(aiConfig.value.openai.temperature) },
                { configKey: 'ai.claude.api_key', configValue: aiConfig.value.claude.apiKey },
                { configKey: 'ai.claude.base_url', configValue: aiConfig.value.claude.baseUrl },
                { configKey: 'ai.claude.model', configValue: aiConfig.value.claude.model },
                { configKey: 'ai.template_generation.enabled', configValue: String(aiConfig.value.features.templateGeneration) },
                { configKey: 'ai.code_review.enabled', configValue: String(aiConfig.value.features.codeReview) },
                { configKey: 'ai.variable_suggestion.enabled', configValue: String(aiConfig.value.features.variableSuggestion) },
                { configKey: 'ai.auto_documentation.enabled', configValue: String(aiConfig.value.features.autoDocumentation) },
                { configKey: 'ai.rate_limit.requests_per_hour', configValue: String(aiConfig.value.rateLimit.requestsPerHour) },
                { configKey: 'ai.rate_limit.requests_per_day', configValue: String(aiConfig.value.rateLimit.requestsPerDay) }
            ]
        } else if (activeTab.value === 'system') {
            configs = [
                { configKey: 'system.name', configValue: systemConfig.value.name },
                { configKey: 'system.description', configValue: systemConfig.value.description },
                { configKey: 'system.version', configValue: systemConfig.value.version },
                { configKey: 'system.logo_url', configValue: systemConfig.value.logoUrl }
            ]
        } else if (activeTab.value === 'template') {
            configs = [
                { configKey: 'template.default_author', configValue: templateConfig.value.defaultAuthor },
                { configKey: 'template.default_company', configValue: templateConfig.value.defaultCompany },
                { configKey: 'template.max_file_size', configValue: String(templateConfig.value.maxFileSize) },
                { configKey: 'template.max_files_per_template', configValue: String(templateConfig.value.maxFilesPerTemplate) }
            ]
        } else if (activeTab.value === 'ui') {
            configs = [
                { configKey: 'ui.theme', configValue: uiConfig.value.theme },
                { configKey: 'ui.language', configValue: uiConfig.value.language },
                { configKey: 'ui.page_size', configValue: String(uiConfig.value.pageSize) },
                { configKey: 'ui.editor_theme', configValue: uiConfig.value.editorTheme },
                { configKey: 'ui.editor_font_size', configValue: String(uiConfig.value.editorFontSize) }
            ]
        }

        await batchUpdateConfig(configs)
        message.success('配置保存成功')
    } catch (error) {
        console.error('保存配置失败:', error)
        message.error('保存配置失败')
    } finally {
        saving.value = false
    }
}

function resetConfig() {
    loadConfigs()
    message.info('配置已重置')
}

async function testConnection() {
    if (!aiConfig.value.enabled) {
        message.warning('请先启用AI功能')
        return
    }

    if (aiConfig.value.provider === 'openai' && !aiConfig.value.openai.apiKey) {
        message.warning('请先填写OpenAI API Key')
        return
    }

    if (aiConfig.value.provider === 'claude' && !aiConfig.value.claude.apiKey) {
        message.warning('请先填写Claude API Key')
        return
    }

    testing.value = true
    try {
        const response = await fetch('/api/v1/ai/testConnection', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })

        const result = await response.json()

        if (result.code === 0 && result.data.success) {
            message.success(`连接成功！服务商: ${result.data.provider}, 模型: ${result.data.model}, 延迟: ${result.data.latency}ms`)
        } else {
            message.error(`连接失败: ${result.data.message}`)
        }
    } catch (error) {
        console.error('测试连接失败:', error)
        message.error('测试连接失败，请检查网络和配置')
    } finally {
        testing.value = false
    }
}

// 生命周期
onMounted(() => {
    loadConfigs()
})
</script>

<style scoped>
.settings {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.page-header {
    text-align: center;
    margin-bottom: 24px;
}

.page-header h1 {
    margin: 0 0 8px 0;
    font-size: 28px;
    color: #333;
}

.page-header p {
    margin: 0;
    color: #666;
    font-size: 16px;
}

.settings-menu {
    border: none;
}

.config-section {
    min-height: 600px;
}

.section-header h2 {
    margin: 0 0 8px 0;
    font-size: 20px;
    color: #333;
}

.section-header p {
    margin: 0 0 16px 0;
    color: #666;
}

.config-form {
    margin-top: 24px;
}

.form-tip {
    margin-left: 12px;
    color: #999;
    font-size: 12px;
}

.config-actions {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid #f0f0f0;
}

:deep(.n-form-item) {
    margin-bottom: 24px;
}

:deep(.n-divider) {
    margin: 24px 0;
}
</style>