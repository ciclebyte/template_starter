<template>
    <div class="dashboard">
        <n-spin :show="loading">
            <!-- 基础统计卡片 -->
            <n-grid :cols="5" :x-gap="20" :y-gap="20">
                <n-gi>
                    <n-card>
                        <n-statistic label="模板总数" :value="overview.totalTemplates">
                            <template #prefix>
                                <n-icon color="#18a058">
                                    <DocumentTextOutline />
                                </n-icon>
                            </template>
                        </n-statistic>
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card>
                        <n-statistic label="分类数量" :value="overview.totalCategories">
                            <template #prefix>
                                <n-icon color="#2080f0">
                                    <LayersOutline />
                                </n-icon>
                            </template>
                        </n-statistic>
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card>
                        <n-statistic label="语言支持" :value="overview.totalLanguages">
                            <template #prefix>
                                <n-icon color="#f0a020">
                                    <LanguageOutline />
                                </n-icon>
                            </template>
                        </n-statistic>
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card>
                        <n-statistic label="文件总数" :value="overview.totalFiles">
                            <template #prefix>
                                <n-icon color="#d03050">
                                    <FolderOutline />
                                </n-icon>
                            </template>
                        </n-statistic>
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card>
                        <n-statistic label="变量总数" :value="overview.totalVariables">
                            <template #prefix>
                                <n-icon color="#722ed1">
                                    <CodeOutline />
                                </n-icon>
                            </template>
                        </n-statistic>
                    </n-card>
                </n-gi>
            </n-grid>

            <!-- 图表区域 -->
            <n-grid :cols="2" :x-gap="20" :y-gap="20" class="charts-section">
                <!-- 分类分布饼图 -->
                <n-gi>
                    <n-card title="分类分布">
                        <v-chart class="chart" :option="categoryChartOption" />
                    </n-card>
                </n-gi>

                <!-- 语言流行度柱状图 -->
                <n-gi>
                    <n-card title="语言流行度">
                        <v-chart class="chart" :option="languageChartOption" />
                    </n-card>
                </n-gi>

                <!-- 模板复杂度分析 -->
                <n-gi>
                    <n-card title="模板复杂度分析">
                        <v-chart class="chart" :option="complexityChartOption" />
                    </n-card>
                </n-gi>

                <!-- 使用趋势线图 -->
                <n-gi>
                    <n-card title="创建趋势（30天）">
                        <v-chart class="chart" :option="trendsChartOption" />
                    </n-card>
                </n-gi>
            </n-grid>

            <!-- 快速操作 -->
            <n-card title="快速操作" class="quick-actions">
                <n-space>
                    <n-button type="primary" @click="goToTemplateCreate">
                        <template #icon>
                            <n-icon>
                                <AddOutline />
                            </n-icon>
                        </template>
                        创建模板
                    </n-button>
                    <n-button @click="goToCategoryCreate">
                        <template #icon>
                            <n-icon>
                                <AddOutline />
                            </n-icon>
                        </template>
                        创建分类
                    </n-button>
                    <n-button @click="goToLanguageCreate">
                        <template #icon>
                            <n-icon>
                                <AddOutline />
                            </n-icon>
                        </template>
                        创建语言
                    </n-button>
                    <n-button @click="goToSettings">
                        <template #icon>
                            <n-icon>
                                <SettingsOutline />
                            </n-icon>
                        </template>
                        系统设置
                    </n-button>
                </n-space>
            </n-card>

            <!-- 最近活动 -->
            <n-card title="最近活动" class="recent-activity">
                <n-list>
                    <n-list-item v-for="activity in recentActivities" :key="activity.id">
                        <n-thing :title="activity.title" :description="activity.description">
                            <template #avatar>
                                <n-avatar :style="{ backgroundColor: activity.color }">
                                    <n-icon>
                                        <component :is="activity.icon" />
                                    </n-icon>
                                </n-avatar>
                            </template>
                            <template #footer>
                                <n-time :time="activity.time" relative />
                            </template>
                        </n-thing>
                    </n-list-item>
                </n-list>
            </n-card>
        </n-spin>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import { use } from 'echarts/core'
import { PieChart, BarChart, LineChart } from 'echarts/charts'
import {
    GridComponent,
    TooltipComponent,
    LegendComponent,
    TitleComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import {
    NGrid,
    NGi,
    NCard,
    NStatistic,
    NIcon,
    NSpace,
    NButton,
    NList,
    NListItem,
    NThing,
    NAvatar,
    NTime,
    NSpin
} from 'naive-ui'
import {
    DocumentTextOutline,
    LayersOutline,
    LanguageOutline,
    FolderOutline,
    AddOutline,
    SettingsOutline,
    CreateOutline,
    TrashOutline,
    CodeOutline
} from '@vicons/ionicons5'
import {
    getOverview,
    getCategoryDistribution,
    getLanguagePopularity,
    getTemplateComplexity,
    getUsageTrends
} from '@/api/statistics'

// 注册ECharts组件
use([
    PieChart,
    BarChart,
    LineChart,
    GridComponent,
    TooltipComponent,
    LegendComponent,
    TitleComponent,
    CanvasRenderer
])

const router = useRouter()

// 数据状态
const loading = ref(true)
const overview = ref({
    totalTemplates: 0,
    totalCategories: 0,
    totalLanguages: 0,
    totalFiles: 0,
    totalVariables: 0,
    featuredTemplates: 0,
    templatesWithVariables: 0,
    templatesWithDescription: 0,
    avgFilesPerTemplate: 0
})
const categoryDistribution = ref([])
const languagePopularity = ref([])
const templateComplexity = ref({})
const usageTrends = ref([])

// 图表配置
const categoryChartOption = computed(() => ({
    tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
        bottom: '10%',
        left: 'center'
    },
    series: [
        {
            name: '分类分布',
            type: 'pie',
            radius: ['40%', '70%'],
            center: ['50%', '40%'],
            data: categoryDistribution.value.map(item => ({
                value: item.templateCount,
                name: item.categoryName
            })),
            emphasis: {
                itemStyle: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
            }
        }
    ]
}))

const languageChartOption = computed(() => ({
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'shadow'
        }
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
    },
    xAxis: {
        type: 'category',
        data: languagePopularity.value.map(item => item.languageName),
        axisLabel: {
            rotate: 45
        }
    },
    yAxis: {
        type: 'value'
    },
    series: [
        {
            name: '模板数量',
            type: 'bar',
            data: languagePopularity.value.map(item => item.templateCount),
            itemStyle: {
                color: '#5470c6'
            }
        }
    ]
}))

const complexityChartOption = computed(() => ({
    tooltip: {
        trigger: 'item'
    },
    legend: {
        bottom: '5%',
        left: 'center'
    },
    series: [
        {
            name: '按文件数',
            type: 'pie',
            radius: [0, '30%'],
            center: ['25%', '40%'],
            data: [
                { value: templateComplexity.value.simpleTemplates || 0, name: '简单(1-3文件)' },
                { value: templateComplexity.value.mediumTemplates || 0, name: '中等(4-10文件)' },
                { value: templateComplexity.value.complexTemplates || 0, name: '复杂(10+文件)' }
            ]
        },
        {
            name: '按变量数',
            type: 'pie',
            radius: [0, '30%'],
            center: ['75%', '40%'],
            data: [
                { value: templateComplexity.value.noVariableTemplates || 0, name: '无变量' },
                { value: templateComplexity.value.fewVariableTemplates || 0, name: '少量变量(1-5)' },
                { value: templateComplexity.value.manyVariableTemplates || 0, name: '多变量(5+)' }
            ]
        }
    ]
}))

const trendsChartOption = computed(() => ({
    tooltip: {
        trigger: 'axis'
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
    },
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: usageTrends.value.map(item => item.date)
    },
    yAxis: {
        type: 'value'
    },
    series: [
        {
            name: '创建数量',
            type: 'line',
            stack: 'Total',
            data: usageTrends.value.map(item => item.templateCreated),
            areaStyle: {},
            itemStyle: {
                color: '#91cc75'
            }
        }
    ]
}))

// 加载数据
async function loadStatistics() {
    try {
        loading.value = true

        const [overviewRes, categoryRes, languageRes, complexityRes, trendsRes] = await Promise.all([
            getOverview(),
            getCategoryDistribution(),
            getLanguagePopularity(),
            getTemplateComplexity(),
            getUsageTrends(30)
        ])

        overview.value = overviewRes.data.data
        categoryDistribution.value = categoryRes.data.data.items || []
        languagePopularity.value = languageRes.data.data.items || []
        templateComplexity.value = complexityRes.data.data
        usageTrends.value = trendsRes.data.data.items || []
    } catch (error) {
        console.error('加载统计数据失败:', error)
    } finally {
        loading.value = false
    }
}

// 最近活动
const recentActivities = ref([
    {
        id: 1,
        title: '创建了新模板',
        description: 'Vue3 组件模板',
        time: new Date(Date.now() - 1000 * 60 * 30),
        icon: markRaw(CreateOutline),
        color: '#18a058'
    },
    {
        id: 2,
        title: '更新了分类',
        description: '前端框架分类',
        time: new Date(Date.now() - 1000 * 60 * 60 * 2),
        icon: markRaw(LayersOutline),
        color: '#2080f0'
    },
    {
        id: 3,
        title: '删除了模板',
        description: '旧版本React模板',
        time: new Date(Date.now() - 1000 * 60 * 60 * 5),
        icon: markRaw(TrashOutline),
        color: '#d03050'
    }
])

// 生命周期
onMounted(() => {
    loadStatistics()
})

function goToTemplateCreate() {
    router.push({ name: 'admin-templates-create' })
}

function goToCategoryCreate() {
    router.push({ name: 'admin-categories-create' })
}

function goToLanguageCreate() {
    router.push({ name: 'admin-languages-create' })
}

function goToSettings() {
    router.push({ name: 'admin-settings' })
}
</script>

<style scoped>
.dashboard {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.charts-section {
    margin-top: 24px;
}

.chart {
    height: 300px;
    width: 100%;
}

.quick-actions,
.recent-activity {
    margin-top: 24px;
}
</style>