<template>
  <div class="admin-management-page">
    <div class="page-header">
      <h1 class="page-title font-display">🔐 管理员管理</h1>
      <a-button type="primary" @click="showAddModal" :disabled="!isSuperAdmin" class="neon-button">
        <PlusOutlined /> 添加管理员
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="admins"
      :pagination="false"
      row-key="id"
      class="admins-table"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'company'">
          <a-tag v-if="record.company" :color="record.company.theme_color">
            {{ record.company.name }}
          </a-tag>
          <a-tag v-else color="purple">超级管理员</a-tag>
        </template>
        <template v-else-if="column.key === 'is_super_admin'">
          <a-tag :color="record.is_super_admin ? 'purple' : 'blue'">
            {{ record.is_super_admin ? '超级管理员' : '普通管理员' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="editAdmin(record)" :disabled="!isSuperAdmin && record.id !== currentAdminId">
              编辑
            </a-button>
            <a-popconfirm
              v-if="isSuperAdmin"
              title="确定要删除这个管理员吗？"
              @confirm="deleteAdmin(record.id)"
            >
              <a-button type="link" danger :disabled="record.id === currentAdminId">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <!-- 添加/编辑管理员弹窗 -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingAdmin ? '编辑管理员' : '添加管理员'"
      width="600px"
      :maskClosable="false"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">👤</span>
            用户名
          </label>
          <a-input
            v-model:value="form.username"
            placeholder="请输入用户名"
            :disabled="editingAdmin != null"
            class="neon-input"
          />
        </a-form-item>

        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">🔑</span>
            密码
          </label>
          <a-input-password
            v-model:value="form.password"
            :placeholder="editingAdmin ? '留空则不修改密码' : '请输入密码（至少6位）'"
            class="neon-input"
          />
        </a-form-item>

        <a-form-item v-if="isSuperAdmin" label="管理员身份" name="admin_type">
          <a-radio-group v-model:value="form.admin_type" @change="handleAdminTypeChange">
            <a-radio value="super_admin">
              <span style="color: #a855f7;">👑 超级管理员</span>
            </a-radio>
            <a-radio value="admin">
              <span style="color: #3b82f6;">👨‍💼 普通管理员</span>
            </a-radio>
          </a-radio-group>
          <template #extra>
            <div style="color: #9ca3af; font-size: 12px; margin-top: 4px;">
              <div v-if="form.admin_type === 'super_admin'">
                💡 超级管理员可以管理所有公司和所有管理员
              </div>
              <div v-else>
                💡 普通管理员只能管理所属公司的用户和抽奖
              </div>
            </div>
          </template>
        </a-form-item>

        <a-form-item v-if="isSuperAdmin && form.admin_type === 'admin'" label="所属公司" name="company_id">
          <a-select
            v-model:value="form.company_id"
            placeholder="请选择所属公司"
            allowClear
            style="width: 100%"
          >
            <a-select-option v-for="company in companies" :key="company.id" :value="company.id">
              {{ company.name }} ({{ company.code }})
            </a-select-option>
          </a-select>
          <template #extra>
            <div style="color: #f59e0b; font-size: 12px; margin-top: 4px;">
              ⚠️ 普通管理员必须选择所属公司
            </div>
          </template>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import request from '../../utils/request'
import { trimObject } from '../../utils/form'

const admins = ref([])
const companies = ref([])
const modalVisible = ref(false)
const editingAdmin = ref(null)
const currentAdminId = ref(null)
const isSuperAdmin = ref(false)

const form = ref({
  username: '',
  password: '',
  admin_type: 'admin', // 'super_admin' or 'admin'
  company_id: undefined
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户名', dataIndex: 'username', width: 150 },
  { title: '所属公司', key: 'company', width: 200 },
  { title: '类型', key: 'is_super_admin', width: 120 },
  { title: '操作', key: 'action', width: 150 }
]

const fetchAdmins = async () => {
  try {
    const data = await request.get('/admin/admins')
    admins.value = data
  } catch (error) {
    message.error('获取管理员列表失败')
  }
}

const fetchCompanies = async () => {
  try {
    const data = await request.get('/admin/companies')
    companies.value = data
  } catch (error) {
    message.error('获取公司列表失败')
  }
}

const fetchCurrentAdmin = async () => {
  try {
    const data = await request.get('/admin/info')
    currentAdminId.value = data.id
    isSuperAdmin.value = data.is_super_admin
  } catch (error) {
  }
}

const showAddModal = () => {
  editingAdmin.value = null
  form.value = {
    username: '',
    password: '',
    admin_type: 'admin',
    company_id: undefined
  }
  modalVisible.value = true
}

const handleAdminTypeChange = (e) => {
  if (e.target.value === 'super_admin') {
    // 切换到超级管理员时，清空公司ID
    form.value.company_id = undefined
  }
}

const editAdmin = (admin) => {
  editingAdmin.value = admin
  form.value = {
    username: admin.username,
    password: '',
    admin_type: admin.is_super_admin ? 'super_admin' : 'admin',
    company_id: admin.company_id
  }
  modalVisible.value = true
}

const handleSubmit = async () => {
  // 去除前后空格
  const trimmedForm = trimObject(form.value)

  if (!trimmedForm.username) {
    message.warning('请输入用户名')
    return
  }
  if (!editingAdmin.value && !trimmedForm.password) {
    message.warning('请输入密码')
    return
  }
  if (isSuperAdmin.value && trimmedForm.admin_type === 'admin' && !trimmedForm.company_id) {
    message.warning('普通管理员必须选择所属公司')
    return
  }

  try {
    if (editingAdmin.value) {
      // 更新
      const updateData = {
        username: trimmedForm.username,
        is_super_admin: trimmedForm.admin_type === 'super_admin'
      }
      if (trimmedForm.password) {
        updateData.password = trimmedForm.password
      }
      if (isSuperAdmin.value) {
        // 只有普通管理员才设置 company_id
        if (trimmedForm.admin_type === 'admin') {
          updateData.company_id = trimmedForm.company_id
        } else {
          updateData.company_id = null
        }
      }
      await request.put(`/admin/admins/${editingAdmin.value.id}`, updateData)
      message.success('更新成功')
    } else {
      // 创建
      const createData = {
        username: trimmedForm.username,
        password: trimmedForm.password,
        is_super_admin: trimmedForm.admin_type === 'super_admin'
      }
      // 只有普通管理员才设置 company_id
      if (trimmedForm.admin_type === 'admin') {
        if (!trimmedForm.company_id) {
          message.warning('普通管理员必须选择所属公司')
          return
        }
        createData.company_id = trimmedForm.company_id
      }
      await request.post('/admin/admins', createData)
      message.success('添加成功')
    }
    modalVisible.value = false
    await fetchAdmins()
  } catch (error) {
    message.error(error.response?.data?.error || '操作失败')
  }
}

const handleCancel = () => {
  modalVisible.value = false
}

const deleteAdmin = async (id) => {
  try {
    await request.delete(`/admin/admins/${id}`)
    message.success('删除成功')
    await fetchAdmins()
  } catch (error) {
    message.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(async () => {
  await fetchCurrentAdmin()
  await fetchAdmins()
  if (isSuperAdmin.value) {
    await fetchCompanies()
  }
})
</script>

<style scoped>
.admin-management-page {
  padding: var(--spacing-lg);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0;
  color: var(--text-primary);
}

.admins-table {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: all var(--transition-base);
}

.admins-table:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan);
}

/* 表单标签样式 */
.form-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: var(--spacing-sm);
  font-weight: var(--font-weight-semibold);
}

.label-icon {
  color: var(--neon-cyan);
  font-size: var(--font-size-lg);
}

/* 输入框统一样式 */
.neon-input :deep(.ant-input),
.neon-input :deep(.ant-input-password input) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
  color: #1a1a1a !important;
  transition: all var(--transition-base);
}

.neon-input :deep(.ant-input::placeholder),
.neon-input :deep(.ant-input-password input::placeholder) {
  color: #8c8c8c !important;
}

.neon-input :deep(.ant-input:focus),
.neon-input :deep(.ant-input-focused),
.neon-input :deep(.ant-input-password:focus),
.neon-input :deep(.ant-input-password-focused) {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 0 2px rgba(0, 255, 245, 0.2);
  background: rgba(255, 255, 255, 1) !important;
}

.neon-input :deep(.ant-input-password) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
}

.neon-input :deep(.ant-input-password:hover) {
  border-color: var(--neon-cyan);
  background: rgba(255, 255, 255, 1) !important;
}

.neon-input :deep(.ant-input-password .ant-input) {
  background: transparent !important;
  color: #1a1a1a !important;
}
</style>
