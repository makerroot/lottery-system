<template>
  <div class="users-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title font-display">👥 用户管理</h1>
        <p class="page-subtitle font-body">管理参与抽奖的用户信息</p>
      </div>
      <div class="header-right">
        <a-space :wrap="true">
          <a-select
            v-model:value="selectedCompanyId"
            style="width: 200px"
            placeholder="选择公司"
            @change="fetchUsers"
            class="company-select neon-select"
          >
            <a-select-option v-for="company in companies" :key="company.id" :value="company.id">
              {{ company.name }}
            </a-select-option>
          </a-select>
          <a-button type="primary" @click="showAddModal" class="add-btn neon-button" size="large">
            <PlusOutlined /> 添加用户
          </a-button>
          <a-button @click="showBatchImportModal" class="neon-button-secondary" size="large">
            <UploadOutlined /> 批量导入
          </a-button>
        </a-space>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
          👥
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.total_users }}</div>
          <div class="stat-label font-body">总用户数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);">
          🎰
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.drawn_users }}</div>
          <div class="stat-label font-body">已抽奖</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #1890ff 0%, #40a9ff 100%);">
          ⏳
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.undrawn_users }}</div>
          <div class="stat-label font-body">未抽奖</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #faad14 0%, #ffc53d 100%);">
          📊
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ drawRate }}%</div>
          <div class="stat-label font-body">抽奖率</div>
        </div>
      </div>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="view-switcher">
        <a-radio-group v-model:value="viewMode" button-style="solid">
          <a-radio-button value="table">
            <TableOutlined /> 表格
          </a-radio-button>
          <a-radio-button value="card">
            <AppstoreOutlined /> 卡片
          </a-radio-button>
        </a-radio-group>
      </div>
      <a-input-search
        v-model:value="searchText"
        placeholder="搜索姓名或用户名"
        style="width: 300px"
        allow-clear
      />
    </div>

    <!-- 表格视图 -->
    <div v-if="viewMode === 'table'" class="table-view">
      <a-table
        :columns="columns"
        :data-source="filteredUsers"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        class="users-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="user-cell">
              <div class="user-avatar">
                {{ (record.name || '未')[0] }}
              </div>
              <div class="user-info">
                <div class="user-name font-body">{{ record.name || '未设置' }}</div>
                <div class="user-phone">@{{ record.username }}</div>
              </div>
            </div>
          </template>
          <template v-else-if="column.key === 'has_drawn'">
            <a-tag :color="record.has_drawn ? 'success' : 'default'">
              <span class="status-dot" :class="{ 'active': record.has_drawn }"></span>
              {{ record.has_drawn ? '已抽奖' : '未抽奖' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="editUser(record)">
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除这个用户吗？"
                @confirm="deleteUser(record.id)"
              >
                <a-button type="link" danger size="small">
                  <DeleteOutlined /> 删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>

    <!-- 卡片视图 -->
    <div v-else class="card-view">
      <div class="users-grid">
        <div
          v-for="user in filteredUsers"
          :key="user.id"
          class="user-card"
          :class="{ 'user-card--drawn': user.has_drawn }"
        >
          <div class="user-card-header">
            <div class="user-avatar-large">
              {{ (user.name || '未')[0] }}
            </div>
            <a-tag :color="user.has_drawn ? 'success' : 'default'" class="user-status">
              {{ user.has_drawn ? '已抽奖' : '未抽奖' }}
            </a-tag>
          </div>
          <div class="user-card-body">
            <h3 class="user-card-name">{{ user.name || '未设置姓名' }}</h3>
            <p class="user-card-phone">@{{ user.username }}</p>
            <div class="user-card-actions">
              <a-popconfirm
                title="确定要删除这个用户吗？"
                @confirm="deleteUser(user.id)"
              >
                <a-button type="text" danger size="small">
                  <DeleteOutlined /> 删除
                </a-button>
              </a-popconfirm>
            </div>
          </div>
        </div>
      </div>
      <a-empty v-if="filteredUsers.length === 0" description="暂无用户数据" />
    </div>

    <!-- 添加用户弹窗 -->
    <a-modal
      v-model:open="addModalVisible"
      title="添加用户"
      :maskClosable="false"
      @ok="handleAddUser"
      @cancel="addModalVisible = false"
    >
      <a-form :model="addForm" layout="vertical">
        <a-form-item label="公司">
          <a-select v-model:value="addForm.company_id" placeholder="选择公司">
            <a-select-option v-for="company in companies" :key="company.id" :value="company.id">
              {{ company.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">📝</span>
            姓名 <span style="color: red;">*</span>
          </label>
          <a-input v-model:value="addForm.name" placeholder="请输入姓名" class="neon-input" />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">📱</span>
            手机号（选填）
          </label>
          <a-input v-model:value="addForm.phone" placeholder="请输入手机号" class="neon-input" />
        </a-form-item>
        <a-alert
          message="提示"
          description="添加的用户仅用于抽奖，无法登录。用户可通过扫码注册自己设置账号密码。"
          type="info"
          show-icon
          style="margin-bottom: 16px;"
        />
      </a-form>
    </a-modal>

    <!-- 编辑用户弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑用户"
      :maskClosable="false"
      @ok="handleUpdateUser"
      @cancel="editModalVisible = false"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">👤</span>
            用户名
          </label>
          <a-input v-model:value="editForm.username" disabled class="neon-input" />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">📝</span>
            姓名
          </label>
          <a-input v-model:value="editForm.name" placeholder="请输入姓名" class="neon-input" />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">📱</span>
            手机号（选填）
          </label>
          <a-input v-model:value="editForm.phone" placeholder="请输入手机号" class="neon-input" />
        </a-form-item>
        <a-form-item label="抽奖状态">
          <a-switch v-model:checked="editForm.has_drawn" checked-children="已抽奖" un-checked-children="未抽奖" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 批量导入弹窗 -->
    <a-modal
      v-model:open="batchImportVisible"
      title="批量导入用户"
      @ok="handleBatchImport"
      @cancel="batchImportVisible = false"
      width="700px"
    >
      <a-alert
        message="⚠️ 重要提示"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
      >
        <template #description>
          <div style="line-height: 1.8;">
            <p style="color: #ff4d4f; font-weight: 600;">导入的用户仅用于抽奖池，无法登录系统！</p>
            <p style="color: var(--text-secondary);">• 用户无法登录（只有管理员可登录）</p>
            <p style="color: var(--text-secondary);">• 管理员在抽奖页面代为用户执行抽奖操作</p>
            <p style="color: var(--text-secondary);">• 创建的用户将出现在"未抽奖用户"列表中供管理员选择</p>
          </div>
        </template>
      </a-alert>
      <a-alert
        message="导入格式说明"
        type="info"
        show-icon
        style="margin-bottom: 16px"
      >
        <template #description>
          <div style="line-height: 1.8;">
            <p><strong>每行一个用户，格式：姓名,手机号（可选）</strong></p>
            <p><strong>注意：</strong>使用英文逗号（,）分隔，不是分号</p>
            <div style="background: rgba(255,255,255,0.05); padding: 12px; border-radius: 6px; margin-top: 8px;">
              <p style="margin: 0; color: var(--text-secondary);">示例：</p>
              <pre style="margin: 8px 0; padding: 12px; background: rgba(0,0,0,0.3); border-radius: 4px; font-size: 13px; color: var(--neon-cyan);">张三,13800138000
李四
王五,13900139000</pre>
              <p style="margin: 8px 0 0 0; color: #999; font-size: 12px;">
                💡 提示：
              </p>
              <ul style="margin: 0; padding-left: 20px; color: var(--text-secondary); font-size: 12px;">
                <li>姓名必填，手机号选填</li>
                <li>每行一个用户，用逗号分隔</li>
                <li>如果只有姓名，可以不加逗号</li>
                <li>相同姓名和手机号的用户会被视为重复</li>
              </ul>
            </div>
          </div>
        </template>
      </a-alert>
      <a-form :model="batchForm" layout="vertical">
        <a-form-item label="选择公司">
          <a-select v-model:value="batchForm.company_id" placeholder="请选择公司" style="width: 100%;">
            <a-select-option v-for="company in companies" :key="company.id" :value="company.id">
              {{ company.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="用户列表">
          <a-textarea
            v-model:value="batchForm.users"
            :rows="12"
            placeholder="张三,13800138000&#10;李四&#10;王五,13900139000"
            style="font-family: 'Courier New', monospace; font-size: 14px;"
          />
          <template #extra>
            <div style="color: var(--text-tertiary); font-size: 12px;">
              每行一个用户，格式：姓名,手机号（手机号可选）
            </div>
          </template>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  UploadOutlined,
  TableOutlined,
  AppstoreOutlined,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import request from '../../utils/request'

const viewMode = ref('table')
const searchText = ref('')
const loading = ref(false)
const users = ref([])
const companies = ref([])
const selectedCompanyId = ref(null)
const stats = ref({
  total_users: 0,
  undrawn_users: 0,
  drawn_users: 0
})

const addModalVisible = ref(false)
const addForm = ref({
  company_id: null,
  name: '',
  phone: ''
})

const editModalVisible = ref(false)
const editForm = ref({
  id: null,
  username: '',
  name: '',
  phone: '',
  has_drawn: false
})

const batchImportVisible = ref(false)
const batchForm = ref({
  company_id: null,
  users: ''
})

// 分页配置
const pagination = ref({
  current: 1,
  pageSize: 12,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

// 表格列配置
const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '用户信息', key: 'name', width: 300 },
  { title: '用户名', dataIndex: 'username', key: 'username', width: 150 },
  { title: '抽奖状态', key: 'has_drawn', width: 120 },
  { title: '操作', key: 'action', width: 120, fixed: 'right' }
]

// 计算属性
const drawRate = computed(() => {
  if (stats.value.total_users === 0) return 0
  return ((stats.value.drawn_users / stats.value.total_users) * 100).toFixed(1)
})

const filteredUsers = computed(() => {
  if (!searchText.value) {
    return users.value
  }
  const search = searchText.value.toLowerCase()
  return users.value.filter(user =>
    (user.name && user.name.toLowerCase().includes(search)) ||
    (user.phone && user.phone.includes(search))
  )
})

// 获取公司列表
const fetchCompanies = async () => {
  try {
    const data = await request.get('/admin/companies')
    companies.value = data
    if (data.length > 0) {
      selectedCompanyId.value = data[0].id
      await fetchUsers()
    }
  } catch (error) {
    message.error('获取公司列表失败')
  }
}

// 获取用户列表
const fetchUsers = async () => {
  if (!selectedCompanyId.value) return

  loading.value = true
  try {
    const data = await request.get(`/admin/users?company_id=${selectedCompanyId.value}`)
    users.value = data
    pagination.value.total = data.length

    // 计算统计数据
    stats.value = {
      total_users: data.length,
      undrawn_users: data.filter(u => !u.has_drawn).length,
      drawn_users: data.filter(u => u.has_drawn).length
    }
  } catch (error) {
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 显示添加用户弹窗
const showAddModal = () => {
  addForm.value = {
    company_id: selectedCompanyId.value,
    username: '',
    password: '',
    name: '',
    phone: ''
  }
  addModalVisible.value = true
}

// 添加用户
const handleAddUser = async () => {
  if (!addForm.value.company_id || !addForm.value.name) {
    message.warning('请填写公司和姓名')
    return
  }

  // 去除前后空格
  addForm.value.name = addForm.value.name.trim()
  if (addForm.value.phone) {
    addForm.value.phone = addForm.value.phone.trim()
  }

  try {
    // 只发送 name 和 phone，不发送 username 和 password
    await request.post('/admin/users', {
      company_id: addForm.value.company_id,
      name: addForm.value.name,
      phone: addForm.value.phone
    })
    message.success('添加成功')
    addModalVisible.value = false
    await fetchUsers()
  } catch (error) {
    message.error(error.response?.data?.error || '添加失败')
  }
}

// 显示批量导入弹窗
const showBatchImportModal = () => {
  batchForm.value = {
    company_id: selectedCompanyId.value,
    users: ''
  }
  batchImportVisible.value = true
}

// 批量导入
const handleBatchImport = async () => {
  if (!batchForm.value.company_id || !batchForm.value.users.trim()) {
    message.warning('请选择公司并填写用户数据')
    return
  }

  // 解析用户列表并去除前后空格
  const lines = batchForm.value.users.trim().split('\n')
  const users = lines
    .map(line => line.trim())
    .filter(line => line.length > 0)

  if (users.length === 0) {
    message.warning('请输入至少一个用户')
    return
  }

  // 验证格式
  const invalidLines = []
  const validUsers = []

  users.forEach((line, index) => {
    const parts = line.split(',')
    // 格式：姓名,手机号（可选）
    const name = parts[0].trim()
    const phone = parts[1] ? parts[1].trim() : ''

    if (!name) {
      invalidLines.push(`第${index + 1}行: ${line}（姓名为空）`)
    } else {
      // 重新组合为去除空格后的格式
      validUsers.push(phone ? `${name},${phone}` : name)
    }
  })

  if (invalidLines.length > 0) {
    message.error(`格式错误：\n${invalidLines.join('\n')}`)
    return
  }

  try {
    const result = await request.post('/admin/users/batch', {
      company_id: batchForm.value.company_id,
      users: validUsers
    })

    batchImportVisible.value = false

    // 显示详细结果
    if (result.failed === 0) {
      message.success(`✅ 成功导入 ${result.created} 个用户`)
    } else {
      message.warning(`⚠️ 成功 ${result.created} 个，失败 ${result.failed} 个`)
    }

    // 如果有失败的，显示详细信息
    if (result.errors && result.errors.length > 0) {
      console.log('导入失败的用户:', result.errors)
      // 可以考虑在界面上显示这些错误
    }

    await fetchUsers()
  } catch (error) {
    const errorMsg = error.response?.data?.error || '导入失败'
    message.error(`❌ ${errorMsg}`)
    console.error('批量导入失败:', error)
  }
}

// 删除用户
const deleteUser = async (id) => {
  try {
    await request.delete(`/admin/users/${id}`)
    message.success('删除成功')
    await fetchUsers()
  } catch (error) {
    message.error('删除失败')
  }
}

const editUser = (user) => {
  editForm.value = {
    id: user.id,
    username: user.username,
    name: user.name || '',
    phone: user.phone || '',
    has_drawn: user.has_drawn
  }
  editModalVisible.value = true
}

const handleUpdateUser = async () => {
  if (!editForm.value.name) {
    message.warning('请填写姓名')
    return
  }

  try {
    await request.put(`/admin/users/${editForm.value.id}`, {
      name: editForm.value.name,
      phone: editForm.value.phone,
      has_drawn: editForm.value.has_drawn
    })
    message.success('更新成功')
    editModalVisible.value = false
    await fetchUsers()
  } catch (error) {
    message.error('更新失败')
  }
}

onMounted(() => {
  fetchCompanies()
})
</script>

<style scoped>
.users-page {
  padding: var(--spacing-lg);
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-xl);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
}

.page-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-primary);
  margin: 0;
}

.header-right {
  flex-shrink: 0;
}

.add-btn {
  height: 40px;
  font-weight: var(--font-weight-semibold);
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-1);
  transition: all var(--transition-bounce);
  animation: fadeInUp 0.6s ease-out;
}

.stat-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-2);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: white;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  line-height: 1;
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

/* 表格视图 */
.table-view {
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.users-table {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: all var(--transition-base);
}

.users-table:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.user-avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  background: var(--primary-gradient);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  flex-shrink: 0;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: 2px;
}

.user-phone {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  font-family: var(--font-mono);
}

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--gray-400);
  margin-right: var(--spacing-xs);
}

.status-dot.active {
  background: var(--success-color);
}

/* 卡片视图 */
.card-view {
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.user-card {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-1);
  transition: all var(--transition-bounce);
  animation: slideInUp 0.5s ease-out;
}

.user-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-3);
  transform: translateY(-4px);
}

.user-card--drawn {
  border: 2px solid var(--success-color);
}

.user-card-header {
  padding: var(--spacing-xl);
  text-align: center;
  background: var(--bg-secondary);
  position: relative;
}

.user-avatar-large {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: var(--primary-gradient);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0 auto var(--spacing-md);
}

.user-status {
  position: absolute;
  top: var(--spacing-md);
  right: var(--spacing-md);
}

.user-card-body {
  padding: var(--spacing-lg);
  text-align: center;
}

.user-card-name {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
}

.user-card-phone {
  font-size: var(--font-size-base);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-lg) 0;
  font-family: var(--font-mono);
}

.user-card-actions {
  display: flex;
  gap: var(--spacing-sm);
  justify-content: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .users-page {
    padding: var(--spacing-md);
  }

  .page-title {
    font-size: var(--font-size-2xl);
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .users-grid {
    grid-template-columns: 1fr;
  }
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

/* 扫码功能样式 */
.scan-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  border: none;
  color: white;
}

.scan-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.scan-modal-content {
  min-height: 400px;
}

.camera-scan-container {
  text-align: center;
}

.camera-placeholder {
  padding: 60px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.placeholder-icon {
  font-size: 64px;
  opacity: 0.5;
}

.camera-view {
  position: relative;
}

.qrcode-scanner {
  width: 100%;
  height: 300px;
  border: 2px solid var(--neon-cyan);
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: #000;
}

.scan-result-preview {
  margin-top: 16px;
}
</style>
