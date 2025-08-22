<template>
  <div class="permission-management">
    <n-card title="权限管理" :bordered="false">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <!-- 权限管理 -->
        <n-tab-pane name="permissions" tab="权限管理">
          <div class="tab-content">
            <!-- 工具栏 -->
            <div class="toolbar">
              <div class="search-bar">
                <n-input
                  v-model:value="permissionSearchForm.search"
                  placeholder="搜索权限名称、代码或描述"
                  clearable
                  style="width: 300px"
                  @keyup.enter="loadPermissions"
                >
                  <template #prefix>
                    <n-icon>
                      <SearchOutline />
                    </n-icon>
                  </template>
                </n-input>
                <n-select
                  v-model:value="permissionSearchForm.resource"
                  placeholder="选择资源类型"
                  clearable
                  style="width: 150px; margin-left: 12px"
                  :options="resourceOptions"
                  @update:value="loadPermissions"
                />
                <n-button @click="loadPermissions" style="margin-left: 12px">
                  <template #icon>
                    <n-icon>
                      <SearchOutline />
                    </n-icon>
                  </template>
                  搜索
                </n-button>
              </div>
              <div class="actions">
                <n-button type="primary" @click="showPermissionModal()">
                  <template #icon>
                    <n-icon>
                      <AddOutline />
                    </n-icon>
                  </template>
                  新增权限
                </n-button>
              </div>
            </div>

            <!-- 权限表格 -->
            <n-data-table
              :columns="permissionColumns"
              :data="permissionList"
              :loading="permissionLoading"
              :pagination="permissionPagination"
              :row-key="row => row.id"
              @update:page="handlePermissionPageChange"
            />
          </div>
        </n-tab-pane>

        <!-- 角色管理 -->
        <n-tab-pane name="roles" tab="角色管理">
          <div class="tab-content">
            <!-- 工具栏 -->
            <div class="toolbar">
              <div class="search-bar">
                <n-input
                  v-model:value="roleSearchForm.search"
                  placeholder="搜索角色名称、代码或描述"
                  clearable
                  style="width: 300px"
                  @keyup.enter="loadRoles"
                >
                  <template #prefix>
                    <n-icon>
                      <SearchOutline />
                    </n-icon>
                  </template>
                </n-input>
                <n-select
                  v-model:value="roleSearchForm.status"
                  placeholder="选择状态"
                  clearable
                  style="width: 120px; margin-left: 12px"
                  :options="statusOptions"
                  @update:value="loadRoles"
                />
                <n-button @click="loadRoles" style="margin-left: 12px">
                  <template #icon>
                    <n-icon>
                      <SearchOutline />
                    </n-icon>
                  </template>
                  搜索
                </n-button>
              </div>
              <div class="actions">
                <n-button type="primary" @click="showRoleModal()">
                  <template #icon>
                    <n-icon>
                      <AddOutline />
                    </n-icon>
                  </template>
                  新增角色
                </n-button>
              </div>
            </div>

            <!-- 角色表格 -->
            <n-data-table
              :columns="roleColumns"
              :data="roleList"
              :loading="roleLoading"
              :pagination="rolePagination"
              :row-key="row => row.id"
              @update:page="handleRolePageChange"
            />
          </div>
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- 权限编辑弹窗 -->
    <n-modal v-model:show="permissionModalVisible" preset="dialog" title="权限信息">
      <template #header>
        <div>{{ permissionForm.id ? '编辑权限' : '新增权限' }}</div>
      </template>
      
      <n-form
        ref="permissionFormRef"
        :model="permissionForm"
        :rules="permissionRules"
        label-placement="left"
        label-width="80px"
        style="margin-top: 16px"
      >
        <n-form-item label="权限名称" path="name">
          <n-input v-model:value="permissionForm.name" placeholder="请输入权限名称" />
        </n-form-item>

        <n-form-item label="权限代码" path="code">
          <n-input v-model:value="permissionForm.code" placeholder="请输入权限代码，如：user:read" />
        </n-form-item>

        <n-form-item label="资源类型" path="resource">
          <n-input v-model:value="permissionForm.resource" placeholder="请输入资源类型，如：user" />
        </n-form-item>

        <n-form-item label="操作类型" path="action">
          <n-input v-model:value="permissionForm.action" placeholder="请输入操作类型，如：read" />
        </n-form-item>

        <n-form-item label="权限描述" path="description">
          <n-input
            v-model:value="permissionForm.description"
            type="textarea"
            placeholder="请输入权限描述"
            :rows="3"
          />
        </n-form-item>
      </n-form>

      <template #action>
        <n-button @click="permissionModalVisible = false">取消</n-button>
        <n-button type="primary" @click="savePermission" :loading="permissionSaving">
          保存
        </n-button>
      </template>
    </n-modal>

    <!-- 角色编辑弹窗 -->
    <n-modal v-model:show="roleModalVisible" preset="dialog" title="角色信息" style="width: 600px">
      <template #header>
        <div>{{ roleForm.id ? '编辑角色' : '新增角色' }}</div>
      </template>
      
      <n-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleRules"
        label-placement="left"
        label-width="80px"
        style="margin-top: 16px"
      >
        <n-form-item label="角色名称" path="name">
          <n-input v-model:value="roleForm.name" placeholder="请输入角色名称" />
        </n-form-item>

        <n-form-item label="角色代码" path="code">
          <n-input v-model:value="roleForm.code" placeholder="请输入角色代码，如：admin" />
        </n-form-item>

        <n-form-item label="角色状态" path="status">
          <n-select
            v-model:value="roleForm.status"
            placeholder="请选择状态"
            :options="statusOptions"
          />
        </n-form-item>

        <n-form-item label="角色描述" path="description">
          <n-input
            v-model:value="roleForm.description"
            type="textarea"
            placeholder="请输入角色描述"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="角色权限" path="permissions">
          <div class="permission-tree">
            <n-tree
              ref="permissionTreeRef"
              :data="permissionTreeData"
              :default-checked-keys="roleForm.permissions"
              checkable
              cascade
              check-on-click
              expand-on-click
              @update:checked-keys="handlePermissionCheck"
            />
          </div>
        </n-form-item>
      </n-form>

      <template #action>
        <n-button @click="roleModalVisible = false">取消</n-button>
        <n-button type="primary" @click="saveRole" :loading="roleSaving">
          保存
        </n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { SearchOutline, AddOutline, CreateOutline, TrashBinOutline } from '@vicons/ionicons5'
import {
  getPermissions,
  createPermission,
  updatePermission,
  deletePermission,
  getRoles,
  createRole,
  updateRole,
  deleteRole
} from '@/api/permission'

const message = useMessage()
const dialog = useDialog()

// 当前活动的tab
const activeTab = ref('permissions')

// ============================================================================
// 权限管理
// ============================================================================

const permissionList = ref([])
const permissionLoading = ref(false)
const permissionModalVisible = ref(false)
const permissionSaving = ref(false)
const permissionFormRef = ref(null)

// 权限搜索表单
const permissionSearchForm = reactive({
  search: '',
  resource: null,
  page: 1,
  size: 20
})

// 权限表单
const permissionForm = reactive({
  id: null,
  name: '',
  code: '',
  resource: '',
  action: '',
  description: ''
})

// 权限分页
const permissionPagination = reactive({
  page: 1,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page) => {
    permissionPagination.page = page
    permissionSearchForm.page = page
    loadPermissions()
  },
  onUpdatePageSize: (pageSize) => {
    permissionPagination.pageSize = pageSize
    permissionPagination.page = 1
    permissionSearchForm.size = pageSize
    permissionSearchForm.page = 1
    loadPermissions()
  }
})

// 权限表格列定义
const permissionColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '权限名称', key: 'name' },
  { title: '权限代码', key: 'code' },
  { title: '资源类型', key: 'resource' },
  { title: '操作类型', key: 'action' },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  { title: '创建时间', key: 'createdAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row) {
      return [
        h('n-button', {
          size: 'small',
          type: 'primary',
          style: { marginRight: '8px' },
          onClick: () => showPermissionModal(row)
        }, { default: () => '编辑' }),
        h('n-button', {
          size: 'small',
          type: 'error',
          onClick: () => handleDeletePermission(row)
        }, { default: () => '删除' })
      ]
    }
  }
]

// 权限表单验证规则
const permissionRules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限代码', trigger: 'blur' }],
  resource: [{ required: true, message: '请输入资源类型', trigger: 'blur' }],
  action: [{ required: true, message: '请输入操作类型', trigger: 'blur' }]
}

// 资源类型选项
const resourceOptions = [
  { label: '用户', value: 'user' },
  { label: '角色', value: 'role' },
  { label: '权限', value: 'permission' },
  { label: '模板', value: 'template' },
  { label: '分类', value: 'category' },
  { label: '语言', value: 'language' },
  { label: '标签', value: 'tag' },
  { label: '系统', value: 'system' }
]

// 加载权限列表
async function loadPermissions() {
  try {
    permissionLoading.value = true
    const response = await getPermissions(permissionSearchForm)
    
    if (response.data.code === 0) {
      permissionList.value = response.data.data.list || []
      permissionPagination.itemCount = response.data.data.total || 0
    }
  } catch (error) {
    message.error('加载权限列表失败')
  } finally {
    permissionLoading.value = false
  }
}

// 显示权限编辑弹窗
function showPermissionModal(permission = null) {
  if (permission) {
    Object.assign(permissionForm, permission)
  } else {
    Object.assign(permissionForm, {
      id: null,
      name: '',
      code: '',
      resource: '',
      action: '',
      description: ''
    })
  }
  permissionModalVisible.value = true
}

// 保存权限
async function savePermission() {
  try {
    await permissionFormRef.value?.validate()
    permissionSaving.value = true

    const data = { ...permissionForm }
    delete data.id

    if (permissionForm.id) {
      await updatePermission(permissionForm.id, data)
      message.success('权限更新成功')
    } else {
      await createPermission(data)
      message.success('权限创建成功')
    }

    permissionModalVisible.value = false
    loadPermissions()
  } catch (error) {
    if (error.message) {
      message.error(error.message)
    }
  } finally {
    permissionSaving.value = false
  }
}

// 删除权限
function handleDeletePermission(permission) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除权限"${permission.name}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deletePermission(permission.id)
        message.success('权限删除成功')
        loadPermissions()
      } catch (error) {
        message.error('权限删除失败')
      }
    }
  })
}

// 权限分页变化
function handlePermissionPageChange(page) {
  permissionPagination.page = page
  permissionSearchForm.page = page
  loadPermissions()
}

// ============================================================================
// 角色管理
// ============================================================================

const roleList = ref([])
const roleLoading = ref(false)
const roleModalVisible = ref(false)
const roleSaving = ref(false)
const roleFormRef = ref(null)
const permissionTreeRef = ref(null)
const permissionTreeData = ref([])

// 角色搜索表单
const roleSearchForm = reactive({
  search: '',
  status: null,
  page: 1,
  size: 20
})

// 角色表单
const roleForm = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  status: 1,
  permissions: []
})

// 角色分页
const rolePagination = reactive({
  page: 1,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page) => {
    rolePagination.page = page
    roleSearchForm.page = page
    loadRoles()
  },
  onUpdatePageSize: (pageSize) => {
    rolePagination.pageSize = pageSize
    rolePagination.page = 1
    roleSearchForm.size = pageSize
    roleSearchForm.page = 1
    loadRoles()
  }
})

// 状态选项
const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

// 角色表格列定义
const roleColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '角色名称', key: 'name' },
  { title: '角色代码', key: 'code' },
  { 
    title: '状态', 
    key: 'status',
    width: 100,
    render(row) {
      return h('n-tag', {
        type: row.status === 1 ? 'success' : 'default'
      }, {
        default: () => row.status === 1 ? '启用' : '禁用'
      })
    }
  },
  { title: '用户数量', key: 'userCount', width: 100 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  { title: '创建时间', key: 'createdAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row) {
      return [
        h('n-button', {
          size: 'small',
          type: 'primary',
          style: { marginRight: '8px' },
          onClick: () => showRoleModal(row)
        }, { default: () => '编辑' }),
        h('n-button', {
          size: 'small',
          type: 'error',
          onClick: () => handleDeleteRole(row)
        }, { default: () => '删除' })
      ]
    }
  }
]

// 角色表单验证规则
const roleRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色代码', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

// 加载角色列表
async function loadRoles() {
  try {
    roleLoading.value = true
    const response = await getRoles(roleSearchForm)
    
    if (response.data.code === 0) {
      roleList.value = response.data.data.list || []
      rolePagination.itemCount = response.data.data.total || 0
    }
  } catch (error) {
    message.error('加载角色列表失败')
  } finally {
    roleLoading.value = false
  }
}

// 构建权限树数据
function buildPermissionTree(permissions) {
  const groups = {}
  
  permissions.forEach(permission => {
    if (!groups[permission.resource]) {
      groups[permission.resource] = {
        key: permission.resource,
        label: permission.resource,
        children: []
      }
    }
    
    groups[permission.resource].children.push({
      key: permission.id,
      label: `${permission.name} (${permission.code})`,
      permission: permission
    })
  })
  
  return Object.values(groups)
}

// 显示角色编辑弹窗
async function showRoleModal(role = null) {
  // 加载所有权限用于权限树
  try {
    const response = await getPermissions({ page: 1, size: 1000 })
    if (response.data.code === 0) {
      permissionTreeData.value = buildPermissionTree(response.data.data.list || [])
    }
  } catch (error) {
    message.error('加载权限列表失败')
    return
  }

  if (role) {
    Object.assign(roleForm, {
      ...role,
      permissions: role.permissions ? role.permissions.map(p => p.id) : []
    })
  } else {
    Object.assign(roleForm, {
      id: null,
      name: '',
      code: '',
      description: '',
      status: 1,
      permissions: []
    })
  }
  
  roleModalVisible.value = true
}

// 处理权限选择
function handlePermissionCheck(checkedKeys) {
  roleForm.permissions = checkedKeys
}

// 保存角色
async function saveRole() {
  try {
    await roleFormRef.value?.validate()
    roleSaving.value = true

    const data = { ...roleForm }
    delete data.id

    if (roleForm.id) {
      await updateRole(roleForm.id, data)
      message.success('角色更新成功')
    } else {
      await createRole(data)
      message.success('角色创建成功')
    }

    roleModalVisible.value = false
    loadRoles()
  } catch (error) {
    if (error.message) {
      message.error(error.message)
    }
  } finally {
    roleSaving.value = false
  }
}

// 删除角色
function handleDeleteRole(role) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除角色"${role.name}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteRole(role.id)
        message.success('角色删除成功')
        loadRoles()
      } catch (error) {
        message.error('角色删除失败')
      }
    }
  })
}

// 角色分页变化
function handleRolePageChange(page) {
  rolePagination.page = page
  roleSearchForm.page = page
  loadRoles()
}

// 初始化
onMounted(() => {
  loadPermissions()
  loadRoles()
})
</script>

<style scoped>
.permission-management {
  padding: 16px;
}

.tab-content {
  margin-top: 16px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.search-bar {
  display: flex;
  align-items: center;
}

.permission-tree {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e0e0e6;
  border-radius: 6px;
  padding: 8px;
}
</style>