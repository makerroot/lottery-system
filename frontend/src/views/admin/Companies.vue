<template>
  <div class="companies-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title font-display">🏢 公司管理</h1>
        <p class="page-subtitle font-body">管理参与抽奖的公司信息</p>
      </div>
      <div class="header-right">
        <a-button
          v-if="isSuperAdmin"
          type="primary"
          @click="showAddModal"
          class="add-btn neon-button"
        >
          <PlusOutlined /> 添加公司
        </a-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
          🏢
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ companies.length }}</div>
          <div class="stat-label font-body">总公司数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);">
          ✅
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ activeCompaniesCount }}</div>
          <div class="stat-label font-body">已启用</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #faad14 0%, #ffc53d 100%);">
          ⏸️
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ inactiveCompaniesCount }}</div>
          <div class="stat-label font-body">已禁用</div>
        </div>
      </div>
    </div>

    <!-- 视图切换 -->
    <div class="view-controls">
      <div class="view-switcher">
        <a-radio-group v-model:value="viewMode" button-style="solid">
          <a-radio-button value="table">
            <AppstoreOutlined /> 表格视图
          </a-radio-button>
          <a-radio-button value="card">
            <BarsOutlined /> 卡片视图
          </a-radio-button>
        </a-radio-group>
      </div>
      <a-input-search
        v-model:value="searchText"
        placeholder="搜索公司名称或代码"
        style="width: 300px"
        allow-clear
      />
    </div>

    <!-- 表格视图 -->
    <div v-if="viewMode === 'table'" class="table-view">
      <a-table
        :columns="columns"
        :data-source="filteredCompanies"
        :pagination="pagination"
        row-key="id"
        class="companies-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'logo'">
            <div v-if="record.logo" class="company-logo-wrapper">
              <img :src="record.logo" class="company-logo" />
            </div>
            <div v-else class="company-logo-placeholder">
              {{ record.name.charAt(0) }}
            </div>
          </template>
          <template v-else-if="column.key === 'name'">
            <div class="company-name-cell">
              <div class="company-name font-body">{{ record.name }}</div>
              <div class="company-code">@{{ record.code }}</div>
            </div>
          </template>
          <template v-else-if="column.key === 'theme_color'">
            <div class="color-display">
              <span class="color-preview" :style="{ background: record.theme_color }"></span>
              <span class="color-code">{{ record.theme_color }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'is_active'">
            <a-tag :color="record.is_active ? 'success' : 'default'">
              {{ record.is_active ? '✓ 启用' : '✗ 禁用' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="editCompany(record)">
                <EditOutlined /> 编辑
              </a-button>
              <a-popconfirm
                v-if="isSuperAdmin"
                title="确定要删除这个公司吗？"
                @confirm="deleteCompany(record.id)"
              >
                <a-button type="link" size="small" danger>
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
      <div class="companies-grid">
        <div
          v-for="company in filteredCompanies"
          :key="company.id"
          class="company-card"
          :class="{ 'company-card--inactive': !company.is_active }"
        >
          <div class="company-card-header" :style="{ background: company.theme_color }">
            <div class="company-logo-large">
              {{ company.name.charAt(0) }}
            </div>
            <div class="company-status">
              <a-tag :color="company.is_active ? 'cyan' : 'default'" style="margin: 0;">
                {{ company.is_active ? '启用' : '禁用' }}
              </a-tag>
            </div>
          </div>
          <div class="company-card-body">
            <h3 class="company-card-name">{{ company.name }}</h3>
            <p class="company-card-code">@{{ company.code }}</p>

            <div class="company-card-stats">
              <div class="card-stat">
                <span class="card-stat-label font-body">主题色</span>
                <div class="color-preview-large" :style="{ background: company.theme_color }"></div>
              </div>
            </div>


            <div class="company-card-actions">
              <a-button type="primary" size="small" @click="editCompany(company)">
                <EditOutlined /> 编辑
              </a-button>
              <a-popconfirm
                v-if="isSuperAdmin"
                title="确定要删除这个公司吗？"
                @confirm="deleteCompany(company.id)"
              >
                <a-button size="small" danger>
                  <DeleteOutlined /> 删除
                </a-button>
              </a-popconfirm>
            </div>
          </div>
        </div>
      </div>
      <a-empty v-if="filteredCompanies.length === 0" description="暂无公司数据" />
    </div>

    <!-- 添加/编辑公司弹窗 -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingCompany ? '编辑公司' : '添加公司'"
      width="700px"
      :maskClosable="false"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form :model="form" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">🏢</span>
                公司代码
              </label>
              <a-input v-model:value="form.code" placeholder="如：default, acme" :disabled="editingCompany" class="neon-input" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">📛</span>
                公司名称
              </label>
              <a-input v-model:value="form.name" placeholder="公司名称" class="neon-input" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="主题颜色" name="theme_color">
          <a-select
            v-model:value="form.theme_color"
            placeholder="请选择主题颜色"
            style="width: 100%"
            :dropdownStyle="{ maxHeight: 400, overflow: 'auto' }"
          >
            <a-select-opt-group label="🔥 热门推荐">
              <a-select-option v-for="color in popularColors" :key="color.value" :value="color.value">
                <div class="color-option">
                  <span class="color-preview-inline" :style="{ background: color.value }"></span>
                  <span class="color-name">{{ color.name }}</span>
                  <span class="color-code">{{ color.value }}</span>
                </div>
              </a-select-option>
            </a-select-opt-group>

            <a-select-opt-group label="🌈 全部颜色">
              <a-select-option v-for="color in allColors" :key="color.value" :value="color.value">
                <div class="color-option">
                  <span class="color-preview-inline" :style="{ background: color.value }"></span>
                  <span class="color-name">{{ color.name }}</span>
                  <span class="color-code">{{ color.value }}</span>
                </div>
              </a-select-option>
            </a-select-opt-group>

            <a-select-opt-group label="✏️ 自定义">
              <a-select-option value="custom">
                <div class="color-option">
                  <span class="color-preview-inline" style="background: linear-gradient(45deg, #ff0000, #00ff00, #0000ff);"></span>
                  <span class="color-name">自定义颜色</span>
                  <span class="color-code">手动输入</span>
                </div>
              </a-select-option>
            </a-select-opt-group>
          </a-select>

          <!-- 自定义颜色输入 -->
          <div v-if="form.theme_color === 'custom' || !isPresetColor(form.theme_color)" style="margin-top: 8px;">
            <a-input
              v-model:value="customColor"
              placeholder="#00fff5"
              @change="handleCustomColorChange"
              class="neon-input"
            >
              <template #prefix>
                <div class="color-preview" :style="{ background: customColor || form.theme_color || '#00fff5' }"></div>
              </template>
            </a-input>
          </div>

          <template #extra>
            <div style="color: #ffffff; font-size: 12px; margin-top: 4px;">
              💡 选择预设颜色快速配置，或选择"自定义"手动输入颜色代码
            </div>
          </template>
        </a-form-item>

        <a-form-item label="背景颜色" name="bg_color">
          <a-select
            v-model:value="form.bg_color"
            placeholder="请选择背景颜色"
            style="width: 100%"
            :dropdownStyle="{ maxHeight: 400, overflow: 'auto' }"
          >
            <a-select-opt-group label="🌑 深色背景（推荐）">
              <a-select-option value="#0a0f14">深空黑 (#0a0f14)</a-select-option>
              <a-select-option value="#0d1117">午夜黑 (#0d1117)</a-select-option>
              <a-select-option value="#141420">暗夜蓝 (#141420)</a-select-option>
              <a-select-option value="#1a1a2e">深空蓝 (#1a1a2e)</a-select-option>
              <a-select-option value="#1e1e2e">系统灰 (#1e1e2e)</a-select-option>
            </a-select-opt-group>

            <a-select-opt-group label="🌈 彩色背景">
              <a-select-option value="#1a0b2e">深邃蓝 (#1a0b2e)</a-select-option>
              <a-select-option value="#1a0a1f">午夜紫 (#1a0a1f)</a-select-option>
              <a-select-option value="#1f1510">暗红色 (#1f1510)</a-select-option>
              <a-select-option value="#0f1a15">森林绿 (#0f1a15)</a-select-option>
            </a-select-opt-group>

            <a-select-opt-group label="✏️ 自定义">
              <a-select-option value="custom">
                <div class="color-option">
                  <span class="color-preview-inline" style="background: linear-gradient(45deg, #1a1a2e, #2d2d3a);"></span>
                  <span class="color-name">自定义背景色</span>
                  <span class="color-code">手动输入</span>
                </div>
              </a-select-option>
            </a-select-opt-group>
          </a-select>

          <!-- 自定义背景颜色输入 -->
          <div v-if="form.bg_color === 'custom'" style="margin-top: 8px;">
            <a-input
              v-model:value="customBgColor"
              placeholder="#0a0f14"
              @change="handleCustomBgColorChange"
              class="neon-input"
            >
              <template #prefix>
                <div class="color-preview" :style="{ background: customBgColor || form.bg_color || '#0a0f14' }"></div>
              </template>
            </a-input>
          </div>

          <template #extra>
            <div style="color: #ffffff; font-size: 12px; margin-top: 4px;">
              💡 背景颜色用于抽奖页面的整体背景，建议使用深色以突出主题色
            </div>
          </template>
        </a-form-item>

        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">🎯</span>
            系统标题
          </label>
          <a-input v-model:value="form.title" placeholder="如：🎉 幸运大抽奖" class="neon-input" />
        </a-form-item>

        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">📝</span>
            副标题
          </label>
          <a-input v-model:value="form.subtitle" placeholder="副标题" class="neon-input" />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">🎰</span>
                抽奖按钮文字
              </label>
              <a-input v-model:value="form.draw_button_text" placeholder="点击抽奖" class="neon-input" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">🎉</span>
                成功提示文字
              </label>
              <a-input v-model:value="form.success_text" placeholder="恭喜中奖！" class="neon-input" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="欢迎语" name="welcome_text">
          <a-textarea v-model:value="form.welcome_text" :rows="2" placeholder="欢迎参与抽奖活动！" />
        </a-form-item>

        <a-form-item label="规则说明" name="rules_text">
          <a-textarea v-model:value="form.rules_text" :rows="2" placeholder="每人只能抽一次，祝您好运！" />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">👤</span>
                联系人
              </label>
              <a-input v-model:value="form.contact_name" class="neon-input" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">📞</span>
                联系电话
              </label>
              <a-input v-model:value="form.contact_phone" class="neon-input" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item>
              <label class="form-label font-body">
                <span class="label-icon">✉️</span>
                联系邮箱
              </label>
              <a-input v-model:value="form.contact_email" class="neon-input" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="状态" name="is_active">
          <a-switch v-model:checked="form.is_active" checked-children="启用" un-checked-children="禁用" />
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
  EditOutlined,
  DeleteOutlined,
  AppstoreOutlined,
  BarsOutlined
} from '@ant-design/icons-vue'
import request from '../../utils/request'

const companies = ref([])
const modalVisible = ref(false)
const editingCompany = ref(null)
const customColor = ref('')
const customBgColor = ref('')
const viewMode = ref('table')
const searchText = ref('')
const isSuperAdmin = ref(false)

// 预设颜色列表
const popularColors = [
  { name: '霓虹青', value: '#00fff5' },
  { name: '热情红', value: '#ff6b6b' },
  { name: '活力橙', value: '#ffa94d' },
  { name: '阳光黄', value: '#ffd93d' },
  { name: '清新绿', value: '#6bcb77' },
  { name: '天空蓝', value: '#4d96ff' }
]

const allColors = [
  { name: '霓虹青', value: '#00fff5' },
  { name: '玫瑰红', value: '#e74c3c' },
  { name: '珊瑚橙', value: '#ff7675' },
  { name: '金橘色', value: '#fdcb6e' },
  { name: '柠檬黄', value: '#ffeaa7' },
  { name: '薄荷绿', value: '#00b894' },
  { name: '青绿色', value: '#55efc4' },
  { name: '天空蓝', value: '#0984e3' },
  { name: '海洋蓝', value: '#74b9ff' },
  { name: '薰衣草', value: '#a29bfe' },
  { name: '紫水晶', value: '#6c5ce7' },
  { name: '粉玫瑰', value: '#fd79a8' },
  { name: '樱桃红', value: '#e84393' },
  { name: '深紫色', value: '#6c5ce7' },
  { name: '皇家蓝', value: '#4834d4' },
  { name: '翡翠绿', value: '#00cec9' },
  { name: '孔雀蓝', value: '#0984e3' }
]

const form = ref({
  code: '',
  name: '',
  logo: '',
  theme_color: '#00fff5',
  bg_color: '#0a0f14',
  title: '',
  subtitle: '',
  welcome_text: '',
  rules_text: '',
  draw_button_text: '点击抽奖',
  success_text: '恭喜中奖！',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  is_active: true
})

// 分页配置
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

// 表格列配置
const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: 'Logo', key: 'logo', width: 80 },
  { title: '公司名称', key: 'name', width: 200 },
  { title: '主题色', key: 'theme_color', width: 150 },
  { title: '状态', key: 'is_active', width: 100 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' }
]

// 计算属性
const activeCompaniesCount = computed(() => {
  return companies.value.filter(c => c.is_active).length
})

const inactiveCompaniesCount = computed(() => {
  return companies.value.filter(c => !c.is_active).length
})

const filteredCompanies = computed(() => {
  if (!searchText.value) {
    return companies.value
  }
  const search = searchText.value.toLowerCase()
  return companies.value.filter(company =>
    company.name.toLowerCase().includes(search) ||
    company.code.toLowerCase().includes(search)
  )
})

// 检查是否是预设颜色
const isPresetColor = (color) => {
  const allPresetColors = [...popularColors, ...allColors]
  return allPresetColors.some(c => c.value === color)
}

// 处理自定义颜色变化
const handleCustomColorChange = (e) => {
  const value = e.target.value
  if (value && value.startsWith('#')) {
    form.value.theme_color = value
  }
}

// 处理自定义背景颜色变化
const handleCustomBgColorChange = (e) => {
  const value = e.target.value
  if (value && value.startsWith('#')) {
    form.value.bg_color = value
  }
}

const fetchCompanies = async () => {
  try {
    const data = await request.get('/admin/companies')
    companies.value = data
    pagination.value.total = data.length
  } catch (error) {
    message.error('获取公司列表失败')
  }
}

const fetchCurrentAdmin = async () => {
  try {
    const data = await request.get('/admin/info')
    isSuperAdmin.value = data.is_super_admin
  } catch (error) {
    console.error('获取管理员信息失败', error)
  }
}

const showAddModal = () => {
  editingCompany.value = null
  form.value = {
    code: '',
    name: '',
    logo: '',
    theme_color: '#00fff5',
    title: '',
    subtitle: '',
    welcome_text: '',
    rules_text: '',
    draw_button_text: '点击抽奖',
    success_text: '恭喜中奖！',
    contact_name: '',
    contact_phone: '',
    contact_email: '',
    is_active: true
  }
  customColor.value = ''
  customBgColor.value = ''
  modalVisible.value = true
}

const editCompany = (company) => {
  editingCompany.value = company
  form.value = { ...company }
  customColor.value = company.theme_color
  customBgColor.value = company.bg_color || '#0a0f14'
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.value.code || !form.value.name) {
    message.warning('请填写公司代码和名称')
    return
  }

  try {
    if (editingCompany.value) {
      await request.put(`/admin/companies/${editingCompany.value.id}`, form.value)
      message.success('更新成功')
    } else {
      await request.post('/admin/companies', form.value)
      message.success('添加成功')
    }
    modalVisible.value = false
    await fetchCompanies()
  } catch (error) {
    message.error(error.response?.data?.error || '操作失败')
  }
}

const handleCancel = () => {
  modalVisible.value = false
}

const deleteCompany = async (id) => {
  try {
    await request.delete(`/admin/companies/${id}`)
    message.success('删除成功')
    await fetchCompanies()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(async () => {
  await fetchCurrentAdmin()
  await fetchCompanies()
})
</script>

<style scoped>
.companies-page {
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

/* 视图控制 */
.view-controls {
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

.companies-table {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: all var(--transition-base);
}

.companies-table:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan);
}

.company-logo-wrapper {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-base);
  overflow: hidden;
}

.company-logo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.company-logo-placeholder {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-base);
  background: var(--primary-gradient);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
}

.company-name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.company-name {
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.company-code {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

.color-display {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: var(--radius-base);
  border: 2px solid rgba(0, 0, 0, 0.1);
}

.color-code {
  font-size: var(--font-size-xs);
  font-family: var(--font-mono);
  color: var(--text-primary);
}

.draw-count-label {
  margin-right: var(--spacing-xs);
  color: var(--text-primary);
}

/* 卡片视图 */
.card-view {
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.companies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--spacing-lg);
}

.company-card {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-1);
  transition: all var(--transition-bounce);
  animation: slideInUp 0.5s ease-out;
}

.company-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-3);
  transform: translateY(-4px);
}

.company-card--inactive {
  opacity: 0.6;
}

.company-card-header {
  padding: var(--spacing-xl);
  text-align: center;
  color: white;
  position: relative;
}

.company-logo-large {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0 auto var(--spacing-md);
  border: 3px solid rgba(255, 255, 255, 0.5);
}

.company-status {
  position: absolute;
  top: var(--spacing-md);
  right: var(--spacing-md);
}

.company-card-body {
  padding: var(--spacing-lg);
}

.company-card-name {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  text-align: center;
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
}

.company-card-code {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  text-align: center;
  margin: 0 0 var(--spacing-lg) 0;
}

.company-card-stats {
  display: flex;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
}

.card-stat {
  flex: 1;
  text-align: center;
}

.card-stat-label {
  display: block;
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  margin-bottom: var(--spacing-xs);
}

.color-preview-large {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-base);
  margin: 0 auto;
  border: 2px solid rgba(0, 0, 0, 0.1);
}

.card-stat-value {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--primary-color);
}

.company-card-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.company-card-actions .ant-btn {
  flex: 1;
}

/* 模态框样式保持不变 */
.color-preview {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}

.color-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
}

.color-preview-inline {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.color-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.color-code {
  font-size: 12px;
  color: var(--text-tertiary);
  font-family: var(--font-mono);
}

.ant-select-dropdown .ant-select-item-option:hover .color-option {
  background: #f5f5f5;
  border-radius: 4px;
}

.ant-select-dropdown .ant-select-item-option-selected .color-option {
  background: #e6f7ff;
  border-radius: 4px;
}

/* 响应式 */
@media (max-width: 768px) {
  .companies-page {
    padding: var(--spacing-md);
  }

  .page-title {
    font-size: var(--font-size-2xl);
  }

  .stats-cards {
    grid-template-columns: 1fr;
  }

  .view-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .companies-grid {
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
</style>
