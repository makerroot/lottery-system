<template>
  <div class="prizelevels-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title font-display">🏆 奖品管理</h1>
        <p class="page-subtitle font-body">管理抽奖奖项和库存</p>
      </div>
      <div class="header-right">
        <a-button type="primary" @click="showAddModal" class="add-btn neon-button">
          <PlusOutlined /> 添加奖项
        </a-button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar glass-card">
      <div class="filter-left">
        <div class="filter-item">
          <label class="filter-label font-body">公司筛选：</label>
          <a-select
            v-model:value="selectedCompanyId"
            placeholder="全部公司"
            style="width: 240px"
            @change="handleCompanyChange"
            :loading="companiesLoading"
            allowClear
          >
            <a-select-option :value="null" v-if="isSuperAdmin">
              全部公司
            </a-select-option>
            <a-select-option v-for="company in filteredCompanies" :key="company.id" :value="company.id">
              <div class="company-option">
                <span class="company-name">{{ company.name }}</span>
                <span class="company-code">({{ company.code }})</span>
              </div>
            </a-select-option>
          </a-select>
        </div>
        <div class="filter-stats">
          <span class="stat-item font-body">
            共 <strong>{{ filteredPrizeLevels.length }}</strong> 个奖项
          </span>
          <span v-if="selectedCompany" class="stat-item font-body">
            来自 <strong :style="{ color: selectedCompany.theme_color || 'var(--neon-cyan)' }">{{ selectedCompany.name }}</strong>
          </span>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card glass-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, var(--neon-purple) 0%, #764ba2 100%);">
          🏆
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ filteredPrizeLevels.length }}</div>
          <div class="stat-label font-body">奖项总数</div>
        </div>
      </div>
      <div class="stat-card glass-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, var(--success-color) 0%, #73d13d 100%);">
          ✅
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ filteredActiveLevelsCount }}</div>
          <div class="stat-label font-body">已启用</div>
        </div>
      </div>
      <div class="stat-card glass-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, var(--warning-color) 0%, #ffc53d 100%);">
          📦
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ filteredTotalStock }}</div>
          <div class="stat-label font-body">总库存</div>
        </div>
      </div>
      <div class="stat-card glass-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, var(--neon-cyan) 0%, #1890ff 100%);">
          ✨
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ filteredRemainingStock }}</div>
          <div class="stat-label font-body">剩余库存</div>
        </div>
      </div>
      <div class="stat-card glass-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, var(--info-color) 0%, #40a9ff 100%);">
          🎁
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ filteredUsedStock }}</div>
          <div class="stat-label font-body">已发放</div>
        </div>
      </div>
    </div>

    <!-- 视图切换 -->
    <div class="view-controls">
      <div class="view-switcher">
        <a-radio-group v-model:value="viewMode" button-style="solid" size="large">
          <a-radio-button value="table">
            <TableOutlined /> 表格视图
          </a-radio-button>
          <a-radio-button value="card">
            <AppstoreOutlined /> 卡片视图
          </a-radio-button>
        </a-radio-group>
      </div>
    </div>

    <!-- 奖项列表 -->
    <div v-if="viewMode === 'card'" class="prizelevels-list">
      <div
        v-for="level in filteredPrizeLevels"
        :key="level.id"
        class="level-card"
        :class="{ 'level-card--inactive': !level.is_active }"
      >
        <div class="level-card-header" :style="{ background: getLevelGradient(level.name) }">
          <div class="level-icon">{{ getLevelIcon(level.name) }}</div>
          <div class="level-info">
            <h3 class="level-name">{{ level.name }}</h3>
            <p class="level-description">{{ level.description || '暂无描述' }}</p>
          </div>
          <a-tag :color="level.is_active ? 'cyan' : 'default'" class="level-status">
            {{ level.is_active ? '启用' : '禁用' }}
          </a-tag>
        </div>
        <div class="level-card-body">
          <div class="level-stats">
            <div class="level-stat">
              <span class="level-stat-label font-body">库存</span>
              <span class="level-stat-value">{{ getLevelPrizeStock(level).remaining }}/{{ getLevelPrizeStock(level).total }}</span>
            </div>
            <div class="level-stat">
              <span class="level-stat-label font-body">进度</span>
              <div class="level-progress">
                <a-progress
                  :percent="getStockPercent(level)"
                  :stroke-color="getStockColor(level)"
                  :show-info="false"
                  size="small"
                />
              </div>
            </div>
          </div>
          <div class="level-company">
            <a-tag v-if="level.company" :color="level.company.theme_color || 'blue'">
              {{ level.company.name }}
            </a-tag>
            <span v-else class="no-company">-</span>
          </div>
          <div class="level-actions">
            <a-button type="link" size="small" @click="managePrizes(level)">
              <GiftOutlined /> 管理奖品
            </a-button>
            <a-button type="link" size="small" @click="editLevel(level)">
              <EditOutlined /> 编辑
            </a-button>
            <a-popconfirm
              title="确定要删除这个奖项吗？"
              @confirm="deleteLevel(level.id)"
            >
              <a-button type="link" size="small" danger>
                <DeleteOutlined /> 删除
              </a-button>
            </a-popconfirm>
          </div>
        </div>
      </div>
      <a-empty v-if="filteredPrizeLevels.length === 0" description="暂无奖项数据" />
    </div>

    <!-- 表格视图 -->
    <div v-else class="prizelevels-table">
      <a-table
        :columns="tableColumns"
        :data-source="filteredPrizeLevels"
        :loading="loading"
        :pagination="{ pageSize: 10, showSizeChanger: true, showQuickJumper: true }"
        row-key="id"
        :row-class-name="(record) => record.is_active ? '' : 'table-row-inactive'"
        class="glass-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="table-level-name">
              <span class="level-icon-small">{{ getLevelIcon(record.name) }}</span>
              <div>
                <div class="name-text">{{ record.name }}</div>
                <div class="description-text">{{ record.description || '暂无描述' }}</div>
              </div>
            </div>
          </template>
          <template v-else-if="column.key === 'company'">
            <a-tag v-if="record.company" :color="record.company.theme_color || 'blue'">
              {{ record.company.name }}
            </a-tag>
            <span v-else class="no-company">-</span>
          </template>
          <template v-else-if="column.key === 'stock'">
            <div class="stock-info">
              <div>{{ getLevelPrizeStock(record).remaining }} / {{ getLevelPrizeStock(record).total }}</div>
              <a-progress
                :percent="getStockPercent(record)"
                :stroke-color="getStockColor(record)"
                :show-info="false"
                size="small"
                :stroke-width="4"
              />
            </div>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.is_active ? 'cyan' : 'default'">
              {{ record.is_active ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-button type="link" size="small" @click="managePrizes(record)">
              <GiftOutlined /> 管理奖品
            </a-button>
            <a-button type="link" size="small" @click="editLevel(record)">
              <EditOutlined /> 编辑
            </a-button>
            <a-popconfirm
              title="确定要删除这个奖项吗？"
              @confirm="deleteLevel(record.id)"
            >
              <a-button type="link" size="small" danger>
                <DeleteOutlined /> 删除
              </a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
      <a-empty v-if="filteredPrizeLevels.length === 0 && !loading" description="暂无奖项数据" />
    </div>

    <!-- 添加/编辑奖项模态框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingLevel ? '编辑奖项' : '添加奖项'"
      :maskClosable="false"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item v-if="isSuperAdmin" label="所属公司">
          <a-select
            v-model:value="form.company_id"
            placeholder="请选择公司"
            style="width: 100%"
          >
            <a-select-option v-for="company in companies" :key="company.id" :value="company.id">
              {{ company.name }} ({{ company.code }})
            </a-select-option>
          </a-select>
          <template #extra>
            <div style="color: #ffffff; font-size: 12px; margin-top: 4px;">
              💡 普通管理员只能为本公司创建奖项
            </div>
          </template>
        </a-form-item>

        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">🏆</span>
            奖项名称
          </label>
          <a-input v-model:value="form.name" placeholder="如：一等奖" class="neon-input" />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">📝</span>
            描述
          </label>
          <a-input v-model:value="form.description" placeholder="奖项描述" class="neon-input" />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">🔢</span>
            排序
          </label>
          <a-input-number
            v-model:value="form.sort_order"
            :min="0"
            style="width: 100%"
            class="neon-input"
          />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="form.is_active" checked-children="启用" un-checked-children="禁用" />
        </a-form-item>
        <a-alert
          message="💡 库存管理提示"
          description="奖项的库存由其包含的所有奖品决定。请在添加该奖项后，点击'管理奖品'按钮为该奖项添加奖品并设置每个奖品的库存数量。"
          type="info"
          show-icon
          closable
          style="margin-bottom: 16px;"
        />
      </a-form>
    </a-modal>

    <!-- 奖品管理模态框 -->
    <a-modal
      v-model:open="prizesModalVisible"
      title="管理奖品"
      width="800px"
      :footer="null"
      @cancel="closePrizesModal"
    >
      <template #title>
        <div style="display: flex; align-items: center; gap: 12px;">
          <GiftOutlined style="font-size: 20px;" />
          <span>
            {{ currentLevelForPrizes ? `奖品列表 - ${currentLevelForPrizes.name}` : '奖品列表' }}
          </span>
        </div>
      </template>

      <a-spin :spinning="prizesLoading">
        <div class="prizes-management">
          <div class="prizes-header">
            <span class="prizes-count">共 {{ prizes.length }} 个奖品</span>
          </div>

          <div v-if="prizes.length === 0 && !prizesLoading" class="empty-prizes">
            <a-empty description="暂无奖品，在下方添加">
            </a-empty>
          </div>

          <div v-else class="prizes-list">
            <div
              v-for="prize in prizes"
              :key="prize.id"
              class="prize-item"
            >
              <div v-if="editingPrize && editingPrize.id === prize.id" class="prize-edit-form">
                <a-form layout="inline" style="width: 100%;">
                  <a-form-item style="flex: 2; margin-bottom: 0;">
                    <label class="form-label font-body">
                      <span class="label-icon">🎁</span>
                      奖品名称
                    </label>
                    <a-input
                      v-model:value="prizeForm.name"
                      placeholder="奖品名称"
                      @pressEnter="handlePrizeSubmit"
                      class="neon-input"
                    />
                  </a-form-item>
                  <a-form-item label="总库存" style="margin-bottom: 0;">
                    <a-input-number
                      v-model:value="prizeForm.total_stock"
                      :min="prizeForm.used_stock || 0"
                      :max="9999"
                      style="width: 100px"
                      class="neon-input"
                    />
                  </a-form-item>
                  <a-form-item label="已发放" style="margin-bottom: 0;">
                    <a-input-number
                      v-model:value="prizeForm.used_stock"
                      :min="0"
                      :max="prizeForm.total_stock || 9999"
                      style="width: 100px"
                      class="neon-input"
                    />
                  </a-form-item>
                  <a-form-item style="margin-bottom: 0;">
                    <a-tag color="blue">
                      剩余: {{ (prizeForm.total_stock || 0) - (prizeForm.used_stock || 0) }}
                    </a-tag>
                  </a-form-item>
                  <a-space>
                    <a-button type="primary" size="small" @click="handlePrizeSubmit">
                      保存
                    </a-button>
                    <a-button size="small" @click="handlePrizeCancel">
                      取消
                    </a-button>
                  </a-space>
                </a-form>
              </div>

              <div v-else class="prize-info">
                <div class="prize-details">
                  <div class="prize-name">
                    <GiftOutlined style="margin-right: 8px;" />
                    {{ prize.name }}
                  </div>
                  <div class="prize-stock">
                    <a-tag :color="getPrizeStockColor(prize)">
                      库存: {{ (prize.total_stock || 0) - (prize.used_stock || 0) }}/{{ prize.total_stock || 0 }}
                    </a-tag>
                  </div>
                </div>
                <div class="prize-actions">
                  <a-button type="link" size="small" @click="editPrize(prize)">
                    <EditOutlined /> 编辑
                  </a-button>
                  <a-popconfirm
                    title="确定要删除这个奖品吗？"
                    @confirm="deletePrize(prize.id)"
                  >
                    <a-button type="link" size="small" danger>
                      <DeleteOutlined /> 删除
                    </a-button>
                  </a-popconfirm>
                </div>
              </div>
            </div>
          </div>

          <div v-if="!editingPrize" class="add-prize-form">
            <a-divider>添加新奖品</a-divider>
            <a-alert
              message="💡 添加奖品"
              description="设置奖品的总库存数量。新添加的奖品初始已发放数量为 0。"
              type="info"
              show-icon
              closable
              style="margin-bottom: 12px;"
            />
            <a-form layout="inline" style="width: 100%;">
              <a-form-item style="flex: 2; margin-bottom: 0;">
                <label class="form-label font-body">
                  <span class="label-icon">🎁</span>
                  奖品名称
                </label>
                <a-input
                  v-model:value="prizeForm.name"
                  placeholder="输入奖品名称"
                  @pressEnter="handlePrizeSubmit"
                  class="neon-input"
                />
              </a-form-item>
              <a-form-item label="总库存" style="margin-bottom: 0;">
                <a-input-number
                  v-model:value="prizeForm.total_stock"
                  :min="0"
                  :max="9999"
                  placeholder="总库存"
                  style="width: 120px"
                  class="neon-input"
                />
              </a-form-item>
              <a-button type="primary" @click="handlePrizeSubmit">
                <PlusOutlined /> 添加
              </a-button>
            </a-form>
          </div>
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, h } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, EditOutlined, DeleteOutlined, AppstoreOutlined, TableOutlined, GiftOutlined, PlusCircleOutlined } from '@ant-design/icons-vue'
import request from '../../utils/request'
import { trimObject } from '../../utils/form'

const prizeLevels = ref([])
const companies = ref([])
const allPrizes = ref([]) // 存储所有奖品数据（用于计算库存）
const modalVisible = ref(false)
const editingLevel = ref(null)
const isSuperAdmin = ref(false)
const currentCompanyId = ref(null)
const selectedCompanyId = ref(null)
const companiesLoading = ref(false)
const loading = ref(false)
const viewMode = ref('table') // 'card' or 'table' - 默认表格视图
const form = ref({
  company_id: undefined,
  name: '',
  description: '',
  sort_order: 0,
  is_active: true
})

// 奖品管理相关
const prizesModalVisible = ref(false)
const currentLevelForPrizes = ref(null)
const prizes = ref([])
const prizesLoading = ref(false)
const prizeForm = ref({
  name: '',
  level_id: null,
  total_stock: 1,
  used_stock: 0
})
const editingPrize = ref(null)

// 表格列配置
const tableColumns = [
  {
    title: '奖项名称',
    key: 'name',
    width: 250,
    dataIndex: 'name'
  },
  {
    title: '所属公司',
    key: 'company',
    width: 150,
    dataIndex: 'company'
  },
  {
    title: '库存使用',
    key: 'stock',
    width: 180
  },
  {
    title: '排序',
    key: 'sort_order',
    width: 100,
    dataIndex: 'sort_order',
    align: 'center'
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    align: 'center'
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    align: 'center',
    fixed: 'right'
  }
]

// 计算过滤后的公司列表（只显示激活的公司）
const filteredCompanies = computed(() => {
  return companies.value.filter(c => c.is_active)
})

// 计算选中的公司
const selectedCompany = computed(() => {
  if (!selectedCompanyId.value) return null
  return companies.value.find(c => c.id === selectedCompanyId.value)
})

// 计算过滤后的奖项列表
const filteredPrizeLevels = computed(() => {
  if (!selectedCompanyId.value) {
    // 未选择公司，返回所有奖项
    return prizeLevels.value
  }
  // 返回选定公司的奖项
  return prizeLevels.value.filter(level => level.company_id === selectedCompanyId.value)
})

const fetchCurrentAdmin = async () => {
  try {
    const data = await request.get('/admin/info')
    isSuperAdmin.value = data.is_super_admin
    currentCompanyId.value = data.company_id
  } catch (error) {
  }
}

const fetchCompanies = async () => {
  try {
    companiesLoading.value = true
    const data = await request.get('/admin/companies')
    companies.value = data.filter(c => c.is_active)
  } catch (error) {
    message.error('获取公司列表失败')
  } finally {
    companiesLoading.value = false
  }
}

const handleCompanyChange = (value) => {
  selectedCompanyId.value = value
}

const fetchPrizeLevels = async () => {
  try {
    loading.value = true
    const data = await request.get('/admin/prize-levels')
    prizeLevels.value = data

    // 同时获取所有奖品（用于库存计算）
    await fetchAllPrizes()
  } catch (error) {
    message.error('获取奖项列表失败')
  } finally {
    loading.value = false
  }
}

// 获取所有奖品
const fetchAllPrizes = async () => {
  try {
    const data = await request.get('/admin/prizes/all')
    allPrizes.value = data || []
  } catch (error) {
    console.error('获取奖品列表失败:', error)
  }
}

const showAddModal = () => {
  editingLevel.value = null
  form.value = {
    company_id: isSuperAdmin.value ? undefined : currentCompanyId.value,
    name: '',
    description: '',
    sort_order: 0,
    is_active: true
  }
  modalVisible.value = true
}

const editLevel = (level) => {
  editingLevel.value = level
  form.value = {
    company_id: level.company_id,
    name: level.name,
    description: level.description,
    sort_order: level.sort_order,
    is_active: level.is_active
  }
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    // 去除所有字符串字段的前后空格
    const payload = trimObject({
      ...form.value
    })

    if (editingLevel.value) {
      await request.put(`/admin/prize-levels/${editingLevel.value.id}`, payload)
      message.success('更新成功')
    } else {
      await request.post('/admin/prize-levels', payload)
      message.success('添加成功')
    }
    modalVisible.value = false
    fetchPrizeLevels()
  } catch (error) {
    message.error('操作失败')
  }
}

const handleCancel = () => {
  modalVisible.value = false
}

const deleteLevel = async (id) => {
  try {
    await request.delete(`/admin/prize-levels/${id}`)
    message.success('删除成功')
    fetchPrizeLevels()
  } catch (error) {
    message.error('删除失败')
  }
}

// ==================== 奖品管理功能 ====================

// 打开奖品管理模态框
const managePrizes = async (level) => {
  currentLevelForPrizes.value = level
  prizeForm.value.level_id = level.id
  editingPrize.value = null
  await fetchPrizes(level.id)
  prizesModalVisible.value = true
}

// 获取奖品列表
const fetchPrizes = async (levelId) => {
  try {
    prizesLoading.value = true
    const data = await request.get(`/admin/prizes/${levelId}`)
    prizes.value = data || []
  } catch (error) {
    message.error('获取奖品列表失败')
  } finally {
    prizesLoading.value = false
  }
}

// 显示添加奖品表单
const showAddPrizeForm = () => {
  editingPrize.value = null
  prizeForm.value = {
    name: '',
    level_id: currentLevelForPrizes.value.id,
    total_stock: 1,
    used_stock: 0
  }
}

// 编辑奖品
const editPrize = (prize) => {
  editingPrize.value = prize
  prizeForm.value = {
    name: prize.name,
    level_id: prize.level_id,
    total_stock: prize.total_stock || 0,
    used_stock: prize.used_stock || 0
  }
}

// 删除奖品
const deletePrize = async (prizeId) => {
  try {
    await request.delete(`/admin/prizes/${prizeId}`)
    message.success('删除奖品成功')

    // 刷新模态框内的奖品列表
    await fetchPrizes(currentLevelForPrizes.value.id)

    // 直接更新全局奖品数据（删除对应的奖品）
    allPrizes.value = allPrizes.value.filter(p => p.id !== prizeId)
  } catch (error) {
    message.error('删除奖品失败')
  }
}

// 提交奖品表单
const handlePrizeSubmit = async () => {
  try {
    if (!prizeForm.value.name || prizeForm.value.name.trim() === '') {
      message.error('请输入奖品名称')
      return
    }

    // 验证库存：总库存必须 >= 已发放
    const totalStock = prizeForm.value.total_stock || 0
    const usedStock = prizeForm.value.used_stock || 0

    if (totalStock < usedStock) {
      message.error(`总库存 (${totalStock}) 不能小于已发放 (${usedStock})`)
      return
    }

    let updatedPrize = null

    if (editingPrize.value) {
      // 更新奖品
      const res = await request.put(`/admin/prizes/${editingPrize.value.id}`, prizeForm.value)
      updatedPrize = res
      message.success('更新奖品成功')
    } else {
      // 添加奖品
      const res = await request.post('/admin/prizes', prizeForm.value)
      updatedPrize = res
      message.success('添加奖品成功')
    }

    // 重置表单
    prizeForm.value = {
      name: '',
      level_id: currentLevelForPrizes.value.id,
      total_stock: 1,
      used_stock: 0
    }
    editingPrize.value = null

    // 刷新模态框内的奖品列表
    await fetchPrizes(currentLevelForPrizes.value.id)

    // 直接更新全局奖品数据
    if (updatedPrize) {
      const existingIndex = allPrizes.value.findIndex(p => p.id === updatedPrize.id)
      if (existingIndex >= 0) {
        // 更新现有奖品
        allPrizes.value[existingIndex] = updatedPrize
      } else {
        // 添加新奖品
        allPrizes.value.push(updatedPrize)
      }
    }
  } catch (error) {
    message.error('操作失败')
  }
}

// 取消奖品表单
const handlePrizeCancel = () => {
  prizeForm.value = {
    name: '',
    level_id: currentLevelForPrizes.value.id,
    total_stock: 1,
    used_stock: 0
  }
  editingPrize.value = null
}

// 关闭奖品模态框
const closePrizesModal = () => {
  prizesModalVisible.value = false
  prizes.value = []
  currentLevelForPrizes.value = null
}

// 获取奖品库存颜色（基于剩余库存）
const getPrizeStockColor = (prize) => {
  const total = prize.total_stock || 0
  const used = prize.used_stock || 0
  const remaining = total - used

  if (total === 0) return 'default'
  if (remaining === 0) return 'error'        // 红色：库存耗尽
  if (remaining < total * 0.2) return 'warning'  // 黄色：库存不足（< 20%）
  return 'success'                          // 绿色：库存充足
}

// 获取等级图标
const getLevelIcon = (name) => {
  if (!name) return '🏆'
  if (name.includes('一等')) return '🥇'
  if (name.includes('二等')) return '🥈'
  if (name.includes('三等')) return '🥉'
  if (name.includes('参与')) return '🎁'
  return '🏆'
}

// 获取等级渐变色
const getLevelGradient = (name) => {
  if (!name) return 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  if (name.includes('一等')) return 'linear-gradient(135deg, #f5222d 0%, #cf1322 100%)'
  if (name.includes('二等')) return 'linear-gradient(135deg, #faad14 0%, #d48806 100%)'
  if (name.includes('三等')) return 'linear-gradient(135deg, #52c41a 0%, #389e0d 100%)'
  if (name.includes('参与')) return 'linear-gradient(135deg, #1890ff 0%, #096dd9 100%)'
  return 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
}

// 获取库存颜色（基于剩余百分比）
const getStockColor = (level) => {
  const percent = getStockPercent(level)
  // 剩余百分比，100%表示库存充足，0%表示库存耗尽
  if (percent === 0) return '#f5222d'  // 红色：库存耗尽
  if (percent < 20) return '#faad14'   // 黄色：库存不足（< 20%）
  return '#52c41a'                     // 绿色：库存充足（>= 20%）
}

// 获取奖项的奖品库存总和
const getLevelPrizeStock = (level) => {
  // 使用全局奖品数据，而不是模态框的奖品数据
  const prizesOfLevel = allPrizes.value.filter(p => p.level_id === level.id)

  // 如果没有加载奖品数据，使用奖项等级自己的库存作为备用
  if (prizesOfLevel.length === 0) {
    return {
      total: level.total_stock || 0,
      used: level.used_stock || 0,
      remaining: (level.total_stock || 0) - (level.used_stock || 0)
    }
  }

  const total = prizesOfLevel.reduce((sum, p) => sum + (p.total_stock || 0), 0)
  const used = prizesOfLevel.reduce((sum, p) => sum + (p.used_stock || 0), 0)
  const remaining = total - used
  return { total, used, remaining }
}

// 获取奖项库存百分比（基于奖品库存，显示剩余百分比）
const getStockPercent = (level) => {
  const { total, used } = getLevelPrizeStock(level)
  if (total === 0) return 0
  // 显示剩余百分比（而不是已用百分比）
  return Math.round(((total - used) / total) * 100)
}

// 计算统计数据（基于过滤后的列表）
const filteredActiveLevelsCount = computed(() => {
  return filteredPrizeLevels.value.filter(l => l.is_active).length
})

const filteredTotalStock = computed(() => {
  return filteredPrizeLevels.value.reduce((sum, level) => {
    // 计算该奖项下所有奖品的总库存（使用全局奖品数据）
    const prizesOfLevel = allPrizes.value.filter(p => p.level_id === level.id)
    const prizeStock = prizesOfLevel.reduce((s, p) => s + (p.total_stock || 0), 0)
    return sum + prizeStock
  }, 0)
})

const filteredUsedStock = computed(() => {
  return filteredPrizeLevels.value.reduce((sum, level) => {
    // 计算该奖项下所有奖品的已用库存（使用全局奖品数据）
    const prizesOfLevel = allPrizes.value.filter(p => p.level_id === level.id)
    const prizeUsed = prizesOfLevel.reduce((s, p) => s + (p.used_stock || 0), 0)
    return sum + prizeUsed
  }, 0)
})

// 计算剩余库存
const filteredRemainingStock = computed(() => {
  return filteredTotalStock.value - filteredUsedStock.value
})

// 保留原来的统计（用于全部数据）
const activeLevelsCount = computed(() => {
  return prizeLevels.value.filter(l => l.is_active).length
})

const totalStock = computed(() => {
  return prizeLevels.value.reduce((sum, level) => sum + (level.total_stock || 0), 0)
})

const usedStock = computed(() => {
  return prizeLevels.value.reduce((sum, level) => sum + (level.used_stock || 0), 0)
})

onMounted(async () => {
  await fetchCurrentAdmin()
  await fetchPrizeLevels()
  await fetchCompanies()

  // 如果不是超级管理员，自动选择自己的公司
  if (!isSuperAdmin.value && currentCompanyId.value) {
    selectedCompanyId.value = currentCompanyId.value
  }
})
</script>

<style scoped>
.prizelevels-page {
  padding: var(--spacing-xl);
  max-width: 1600px;
  margin: 0 auto;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-xl) var(--spacing-2xl);
  background: rgba(26, 26, 36, 0.4);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.page-title {
  font-size: var(--font-size-4xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
}

.page-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-primary);
  margin: 0;
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg) var(--spacing-xl);
  margin-bottom: var(--spacing-xl);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.filter-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-xl);
  flex-wrap: wrap;
  flex: 1;
}

/* 视图控制区域 */
.view-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-md) var(--spacing-lg);
  background: rgba(26, 26, 36, 0.4);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.view-switcher {
  display: flex;
  align-items: center;
}

.view-switcher :deep(.ant-radio-group) {
  background: rgba(255, 255, 255, 0.05);
  border-radius: var(--radius-base);
  padding: 4px;
}

.view-switcher :deep(.ant-radio-button-wrapper) {
  margin-right: 0;
}

.view-switcher :deep(.ant-radio-button) {
  background: transparent;
  border: none;
  color: var(--text-primary);
  transition: all var(--transition-base);
}

.view-switcher :deep(.ant-radio-button:hover) {
  color: var(--text-primary);
}

.view-switcher :deep(.ant-radio-button-checked) {
  background: var(--neon-cyan) !important;
  border-color: var(--neon-cyan) !important;
  color: var(--bg-primary) !important;
  box-shadow: 0 2px 8px rgba(0, 255, 245, 0.3);
}

.view-switcher :deep(.ant-radio-button-checked .anticon) {
  color: var(--bg-primary);
}

.filter-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.filter-label {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  white-space: nowrap;
}

.company-option {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.company-name {
  font-weight: var(--font-weight-medium);
}

.company-code {
  color: var(--text-tertiary);
  font-size: var(--font-size-sm);
}

.filter-stats {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  flex-wrap: wrap;
}

.stat-item {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.stat-item strong {
  color: var(--text-primary);
  font-weight: var(--font-weight-semibold);
}

/* 玻璃卡片 */
.glass-card {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
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

/* 奖项列表 */
.prizelevels-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: var(--spacing-lg);
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.level-card {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-1);
  transition: all var(--transition-bounce);
  animation: slideInUp 0.5s ease-out;
}

.level-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-3);
  transform: translateY(-4px);
}

.level-card--inactive {
  opacity: 0.6;
}

.level-card-header {
  padding: var(--spacing-xl);
  text-align: center;
  color: white;
  position: relative;
}

.level-icon {
  font-size: 48px;
  margin-bottom: var(--spacing-sm);
}

.level-info {
  flex: 1;
}

.level-name {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--spacing-xs) 0;
}

.level-description {
  font-size: var(--font-size-sm);
  opacity: 0.9;
  margin: 0;
}

.level-status {
  position: absolute;
  top: var(--spacing-md);
  right: var(--spacing-md);
}

.level-card-body {
  padding: var(--spacing-lg);
}

.level-stats {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.level-stat {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--border-color-light);
}

.level-stat:last-child {
  border-bottom: none;
}

.level-stat-label {
  font-size: var(--font-size-sm);
  color: #ffffff;
}

.level-stat-value {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: #ffffff;
}

.level-progress {
  flex: 1;
  margin-left: var(--spacing-md);
}

.level-company {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
}

.no-company {
  color: var(--text-tertiary);
  font-size: var(--font-size-sm);
}

.level-actions {
  display: flex;
  gap: var(--spacing-sm);
  justify-content: center;
}

/* 表格视图 */
.prizelevels-table {
  animation: fadeInUp 0.6s ease-out 0.1s both;
}

.glass-table {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.glass-table :deep(.ant-table) {
  background: transparent;
}

.glass-table :deep(.ant-table-container) {
  background: transparent;
}

.glass-table :deep(.ant-table-thead > tr > th) {
  background: rgba(0, 255, 245, 0.05);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  font-weight: var(--font-weight-semibold);
  padding: var(--spacing-md) var(--spacing-lg);
}

.glass-table :deep(.ant-table-tbody > tr > td) {
  background: transparent;
  border-bottom: 1px solid var(--border-color-light);
  padding: var(--spacing-md) var(--spacing-lg);
  color: var(--text-primary);
}

.glass-table :deep(.ant-table-tbody > tr:hover > td) {
  background: rgba(0, 255, 245, 0.05);
}

.glass-table :deep(.table-row-inactive) {
  opacity: 0.5;
}

.glass-table :deep(.ant-pagination) {
  color: var(--text-primary);
}

.glass-table :deep(.ant-pagination-item) {
  background: transparent;
  border-color: var(--border-color);
}

.glass-table :deep(.ant-pagination-item a) {
  color: var(--text-primary);
}

.glass-table :deep(.ant-pagination-item-active) {
  background: var(--neon-cyan);
  border-color: var(--neon-cyan);
}

.glass-table :deep(.ant-pagination-item-active a) {
  color: var(--bg-primary);
}

.table-level-name {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.level-icon-small {
  font-size: 32px;
  line-height: 1;
  flex-shrink: 0;
}

.name-text {
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-base);
  color: var(--text-primary);
}

.description-text {
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  margin-top: 2px;
}

.stock-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.stock-info > div:first-child {
  font-weight: var(--font-weight-medium);
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

/* 响应式 */
@media (max-width: 768px) {
  .prizelevels-page {
    padding: var(--spacing-md);
  }

  .page-title {
    font-size: var(--font-size-2xl);
  }

  .filter-left,
  .filter-right {
    width: 100%;
  }

  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .prizelevels-list {
    grid-template-columns: 1fr;
  }
}

/* ==================== 奖品管理样式 ==================== */

.prizes-management {
  padding: var(--spacing-md) 0;
}

.prizes-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.prizes-count {
  font-size: var(--font-size-base);
  color: var(--text-tertiary);
}

.empty-prizes {
  padding: var(--spacing-xxl) 0;
}

.prizes-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.prize-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color-light);
  border-radius: var(--radius-lg);
  transition: all var(--transition-base);
}

.prize-item:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: var(--neon-cyan);
}

.prize-item--editing {
  background: rgba(0, 255, 245, 0.05);
  border-color: var(--neon-cyan);
  border-style: dashed;
}

.prize-edit-form {
  width: 100%;
}

.prize-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex: 1;
}

.prize-name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.prize-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.add-prize-form {
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--border-color-light);
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
.neon-input :deep(.ant-input-password input),
.neon-input :deep(.ant-input-number) {
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
.neon-input :deep(.ant-input-password-focused),
.neon-input :deep(.ant-input-number:focus),
.neon-input :deep(.ant-input-number-focused) {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 0 2px rgba(0, 255, 245, 0.2);
  background: rgba(255, 255, 255, 1) !important;
}

.neon-input :deep(.ant-input-password),
.neon-input :deep(.ant-input-number) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
}

.neon-input :deep(.ant-input-password:hover),
.neon-input :deep(.ant-input-number:hover) {
  border-color: var(--neon-cyan);
  background: rgba(255, 255, 255, 1) !important;
}

.neon-input :deep(.ant-input-password .ant-input) {
  background: transparent !important;
  color: #1a1a1a !important;
}

.neon-input :deep(.ant-input-number-input) {
  background: transparent !important;
  color: #1a1a1a !important;
}
</style>
