<template>
    <div v-if="hasAccess">
        <slot />
    </div>
    <div v-else-if="showFallback" class="no-permission">
        <slot name="fallback">
            <n-empty description="权限不足" />
        </slot>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { NEmpty } from 'naive-ui'

const props = defineProps({
    // 需要的权限列表
    permissions: {
        type: Array,
        default: () => []
    },
    
    // 需要的角色列表
    roles: {
        type: Array,
        default: () => []
    },
    
    // 权限检查模式: 'all' 需要所有权限, 'any' 需要任意一个权限
    mode: {
        type: String,
        default: 'all',
        validator: (value) => ['all', 'any'].includes(value)
    },
    
    // 无权限时是否显示占位内容
    showFallback: {
        type: Boolean,
        default: false
    }
})

const authStore = useAuthStore()

const hasAccess = computed(() => {
    // 检查角色权限
    if (props.roles.length > 0) {
        const hasRole = authStore.hasAnyRole(props.roles)
        if (!hasRole) return false
    }
    
    // 检查功能权限
    if (props.permissions.length > 0) {
        if (props.mode === 'all') {
            return authStore.hasAllPermissions(props.permissions)
        } else {
            return props.permissions.some(permission => authStore.hasPermission(permission))
        }
    }
    
    return true
})
</script>

<style scoped>
.no-permission {
    text-align: center;
    padding: 20px;
    color: #999;
}
</style>