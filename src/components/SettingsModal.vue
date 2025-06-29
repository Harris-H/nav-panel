<template>
  <div v-if="store.isSettingsModalOpen" class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>设置</h2>
        <button class="close-btn" @click="closeModal">&times;</button>
      </div>

      <div class="modal-body">
        <!-- 背景设置 -->
        <div class="settings-section">
          <h3>背景设置</h3>

          <div class="form-group">
            <label>背景类型</label>
            <div class="radio-group">
              <label class="radio-item">
                <input
                  type="radio"
                  value="gradient"
                  v-model="localSettings.background.type"
                  @change="updateBackgroundType"
                />
                <span>渐变</span>
              </label>
              <label class="radio-item">
                <input
                  type="radio"
                  value="color"
                  v-model="localSettings.background.type"
                  @change="updateBackgroundType"
                />
                <span>纯色</span>
              </label>
              <label class="radio-item">
                <input
                  type="radio"
                  value="image"
                  v-model="localSettings.background.type"
                  @change="updateBackgroundType"
                />
                <span>图片</span>
              </label>
            </div>
          </div>

          <div class="form-group" v-if="localSettings.background?.type === 'color'">
            <label>背景颜色</label>
            <input type="color" v-model="localSettings.background.value" />
          </div>

          <div class="form-group" v-if="localSettings.background?.type === 'gradient'">
            <label>渐变样式</label>
            <select v-model="localSettings.background.value">
              <option value="linear-gradient(135deg, #667eea 0%, #764ba2 100%)">蓝紫渐变</option>
              <option value="linear-gradient(135deg, #f093fb 0%, #f5576c 100%)">粉红渐变</option>
              <option value="linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)">蓝青渐变</option>
              <option value="linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)">绿青渐变</option>
              <option value="linear-gradient(135deg, #fa709a 0%, #fee140 100%)">粉黄渐变</option>
              <option value="linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)">淡青粉渐变</option>
            </select>
          </div>

          <div class="form-group" v-if="localSettings.background?.type === 'image'">
            <label>背景图片</label>
            <div class="file-upload-area">
              <input
                ref="backgroundFileInput"
                type="file"
                accept="image/*"
                @change="handleBackgroundImageUpload"
                style="display: none"
              />
              <button @click="backgroundFileInput?.click()" class="btn-upload">选择图片文件</button>
              <span v-if="uploadedImageName" class="uploaded-file-name">
                {{ uploadedImageName }}
              </span>
            </div>
          </div>
        </div>

        <!-- 卡片样式设置 -->
        <div class="settings-section">
          <h3>卡片样式</h3>

          <div class="form-group">
            <label>圆角大小: {{ localSettings.cardStyle.borderRadius }}px</label>
            <input
              type="range"
              min="0"
              max="20"
              v-model.number="localSettings.cardStyle.borderRadius"
            />
          </div>

          <div class="form-group">
            <label>透明度: {{ Math.round(localSettings.cardStyle.opacity * 100) }}%</label>
            <input
              type="range"
              min="0.5"
              max="1"
              step="0.05"
              v-model.number="localSettings.cardStyle.opacity"
            />
          </div>

          <div class="form-group">
            <label class="checkbox-item">
              <input type="checkbox" v-model="localSettings.cardStyle.shadow" />
              <span>显示阴影</span>
            </label>
          </div>
        </div>

        <!-- 布局设置 -->
        <div class="settings-section">
          <h3>布局设置</h3>

          <div class="form-group">
            <label>每行列数: {{ localSettings.layout.columns }}</label>
            <input type="range" min="3" max="8" v-model.number="localSettings.layout.columns" />
          </div>

          <div class="form-group">
            <label>卡片间距: {{ localSettings.layout.gap }}px</label>
            <input type="range" min="10" max="40" v-model.number="localSettings.layout.gap" />
          </div>
        </div>

        <!-- 搜索引擎设置 -->
        <div class="settings-section">
          <h3>搜索引擎设置</h3>

          <div class="form-group">
            <label class="checkbox-item">
              <input type="checkbox" v-model="localSettings.search.enabled" />
              <span>启用搜索功能</span>
            </label>
          </div>

          <div v-if="localSettings.search.enabled" class="search-engines-config">
            <div class="form-group">
              <label>默认搜索引擎</label>
              <select v-model="localSettings.search.defaultEngineId">
                <option
                  v-for="engine in localSettings.search.engines"
                  :key="engine.id"
                  :value="engine.id"
                >
                  {{ engine.name }}
                </option>
              </select>
            </div>

            <div class="form-group">
              <label class="checkbox-item">
                <input type="checkbox" v-model="localSettings.search.openInNewTab" />
                <span>在新标签页中打开搜索结果</span>
              </label>
            </div>

            <div class="form-group">
              <label>搜索引擎列表</label>
              <div class="search-engines-list">
                <div
                  v-for="(engine, index) in localSettings.search.engines"
                  :key="engine.id"
                  class="search-engine-item"
                >
                  <div class="engine-info">
                    <img
                      v-if="engine.icon"
                      :src="engine.icon"
                      :alt="engine.name"
                      class="engine-icon-small"
                    />
                    <div class="engine-details">
                      <div class="engine-name">{{ engine.name }}</div>
                      <div class="engine-url">{{ engine.url }}</div>
                    </div>
                  </div>
                  <div class="engine-actions">
                    <button @click="editSearchEngine(index)" class="btn-icon" title="编辑">
                      ✏️
                    </button>
                    <button
                      @click="removeSearchEngine(index)"
                      class="btn-icon"
                      title="删除"
                      :disabled="localSettings.search.engines.length <= 1"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
              </div>

              <button @click="openAddEngineModal" class="btn-secondary btn-full-width">
                + 添加搜索引擎
              </button>
            </div>
          </div>
        </div>

        <!-- 数据管理 -->
        <div class="settings-section">
          <h3>数据管理</h3>

          <div class="form-group">
            <button @click="exportData" class="btn-secondary">导出数据</button>
            <button @click="importData" class="btn-secondary">导入数据</button>
            <input
              ref="fileInput"
              type="file"
              accept=".json"
              @change="handleFileImport"
              style="display: none"
            />
          </div>

          <div class="form-group">
            <button @click="resetData" class="btn-danger">重置所有数据</button>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="modal-footer">
        <button @click="cancelSettings" class="btn-secondary">取消</button>
        <button @click="saveSettings" class="btn-primary">保存设置</button>
      </div>
    </div>

    <!-- 搜索引擎添加/编辑模态框 -->
    <div
      v-if="showAddEngineModal || editingEngineIndex !== -1"
      class="modal-overlay"
      @click="closeEngineModal"
    >
      <div class="modal-content engine-modal" @click.stop>
        <div class="modal-header">
          <h3>{{ editingEngineIndex === -1 ? '添加搜索引擎' : '编辑搜索引擎' }}</h3>
          <button @click="closeEngineModal" class="close-btn">×</button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label for="engine-name">名称 *</label>
            <input
              id="engine-name"
              v-model="engineForm.name"
              type="text"
              placeholder="例如：Google"
              required
            />
          </div>

          <div class="form-group">
            <label for="engine-url">搜索URL *</label>
            <input
              id="engine-url"
              v-model="engineForm.url"
              type="url"
              placeholder="例如：https://www.google.com/search?q={}"
              required
            />
            <small>使用 {} 作为搜索关键词占位符</small>
          </div>

          <div class="form-group">
            <label for="engine-placeholder">占位符文本</label>
            <input
              id="engine-placeholder"
              v-model="engineForm.placeholder"
              type="text"
              placeholder="例如：使用 Google 搜索..."
            />
          </div>

          <div class="form-group">
            <label>搜索引擎图标</label>
            <div class="icon-upload-area">
              <div class="icon-preview">
                <img
                  v-if="engineForm.iconPreview"
                  :src="engineForm.iconPreview"
                  :alt="engineForm.name"
                  class="preview-image"
                />
                <div v-else class="icon-placeholder">
                  {{ engineForm.name ? engineForm.name.charAt(0) : '?' }}
                </div>
              </div>

              <div class="icon-controls">
                <button @click="selectIconFile" class="btn-secondary" type="button">
                  📁 选择文件
                </button>
                <button
                  v-if="engineForm.iconPreview"
                  @click="clearIcon"
                  class="btn-secondary"
                  type="button"
                >
                  🗑️ 清除
                </button>
                <input
                  ref="iconFileInput"
                  type="file"
                  accept="image/*"
                  @change="handleIconFileChange"
                  style="display: none"
                />
              </div>
            </div>
            <small>支持 PNG、JPG、GIF、WebP、SVG 格式，最大 2MB</small>
          </div>

          <div class="form-group">
            <div class="default-engine-setting">
              <div class="setting-info">
                <span class="setting-title">设为默认搜索引擎</span>
                <small class="setting-description">将此搜索引擎设为默认选项</small>
              </div>
              <label class="switch">
                <input v-model="engineForm.isDefault" type="checkbox" class="switch-input" />
                <span class="switch-slider"></span>
              </label>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeEngineModal" class="btn-secondary">取消</button>
          <button @click="saveEngine" class="btn-primary" :disabled="!isEngineFormValid">
            {{ editingEngineIndex === -1 ? '添加' : '更新' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useAppStore } from '@/stores/app'
import type { AppSettings } from '@/types'

const store = useAppStore()
const fileInput = ref<HTMLInputElement>()
const backgroundFileInput = ref<HTMLInputElement>()
const iconFileInput = ref<HTMLInputElement>()
const uploadedImageName = ref<string>('')

// 搜索引擎编辑相关状态
const showAddEngineModal = ref(false)
const editingEngineIndex = ref(-1)
const isEditingEngine = ref(false) // 标记是否正在编辑搜索引擎
const engineForm = ref({
  name: '',
  url: '',
  placeholder: '',
  isDefault: false,
  iconFile: null as File | null,
  iconPreview: '' as string,
})

// 保存原始设置用于取消操作
let originalSettings: AppSettings

// 提供默认值，避免访问未定义属性
const localSettings = ref<AppSettings>({
  theme: 'light',
  background: {
    type: 'gradient',
    value: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  },
  cardStyle: {
    borderRadius: 12,
    opacity: 0.9,
    shadow: true,
  },
  layout: {
    columns: 6,
    cardSize: 'medium',
    showLabels: true,
    gap: 20,
  },
  search: {
    enabled: true,
    defaultEngineId: 'google',
    openInNewTab: true,
    engines: store.defaultSearchEngines,
  },
})

// 监听设置变化，同步到本地状态，并保存原始设置
watch(
  () => store.settings,
  (settings) => {
    if (settings && !isEditingEngine.value) {
      // 只有在不编辑搜索引擎时才更新本地设置
      // 深度克隆设置，确保搜索引擎数据正确同步
      const clonedSettings = JSON.parse(JSON.stringify(settings))

      // 确保搜索引擎配置存在且有效
      if (
        !clonedSettings.search ||
        !clonedSettings.search.engines ||
        clonedSettings.search.engines.length === 0
      ) {
        clonedSettings.search = {
          enabled: true,
          defaultEngineId: 'google',
          openInNewTab: true,
          engines: store.defaultSearchEngines,
        }
      }

      localSettings.value = clonedSettings
      originalSettings = JSON.parse(JSON.stringify(clonedSettings))

      // 检查是否有已上传的图片（base64格式）
      if (
        settings.background?.type === 'image' &&
        settings.background?.value &&
        settings.background.value.startsWith('data:image/')
      ) {
        uploadedImageName.value = '已上传的图片'
      } else {
        uploadedImageName.value = ''
      }
    }
  },
  { immediate: true, deep: true },
)

const closeModal = () => {
  store.closeSettingsModal()
}

const updateBackgroundType = () => {
  if (!localSettings.value.background) return

  if (localSettings.value.background.type === 'color') {
    localSettings.value.background.value = '#667eea'
  } else if (localSettings.value.background.type === 'gradient') {
    localSettings.value.background.value = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  } else if (localSettings.value.background.type === 'image') {
    localSettings.value.background.value = ''
    uploadedImageName.value = '' // 清除上传的文件名
  }
}

// 保存设置
const saveSettings = async () => {
  try {
    await store.updateSettings(localSettings.value)
    store.closeSettingsModal()
  } catch (error) {
    console.error('Failed to save settings:', error)
    alert('保存设置失败，请重试')
  }
}

// 取消设置，恢复原始设置
const cancelSettings = () => {
  if (originalSettings) {
    localSettings.value = JSON.parse(JSON.stringify(originalSettings))
  }
  store.closeSettingsModal()
}

const exportData = () => {
  const data = {
    websites: store.sites,
    settings: store.settings,
    exportTime: new Date().toISOString(),
  }

  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `nav-panel-backup-${new Date().toISOString().split('T')[0]}.json`
  a.click()
  URL.revokeObjectURL(url)
}

const importData = () => {
  fileInput.value?.click()
}

const handleFileImport = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const data = JSON.parse(e.target?.result as string)
      if (data.websites && data.settings) {
        if (confirm('导入数据将覆盖当前所有数据，确定继续吗？')) {
          // 更新网站数据
          store.sites.length = 0
          store.sites.push(...data.websites)
          store.updateSettings(data.settings)
          store.loadData() // 重新加载数据
          alert('数据导入成功！')
        }
      } else {
        alert('文件格式不正确！')
      }
    } catch {
      alert('文件解析失败！')
    }
  }
  reader.readAsText(file)
}

const handleBackgroundImageUpload = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    alert('请选择图片文件！')
    return
  }

  // 检查文件大小（限制为5MB）
  if (file.size > 5 * 1024 * 1024) {
    alert('图片文件不能超过5MB！')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    const result = e.target?.result as string
    localSettings.value.background.value = result
    uploadedImageName.value = file.name
  }
  reader.readAsDataURL(file)
}

const resetData = () => {
  if (confirm('确定要重置所有数据吗？此操作不可恢复！')) {
    localStorage.clear()
    location.reload()
  }
}

// 搜索引擎管理方法
// 表单验证
const isEngineFormValid = computed(() => {
  return (
    engineForm.value.name.trim() !== '' &&
    engineForm.value.url.trim() !== '' &&
    engineForm.value.url.includes('{}')
  )
})

// 打开添加搜索引擎模态框
const openAddEngineModal = () => {
  isEditingEngine.value = true
  showAddEngineModal.value = true
}

// 关闭搜索引擎模态框
const closeEngineModal = () => {
  showAddEngineModal.value = false
  editingEngineIndex.value = -1
  isEditingEngine.value = false
  engineForm.value = {
    name: '',
    url: '',
    placeholder: '',
    isDefault: false,
    iconFile: null,
    iconPreview: '',
  }
}

// 选择图标文件
const selectIconFile = () => {
  iconFileInput.value?.click()
}

// 处理图标文件变化
const handleIconFileChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    alert('请选择图片文件！')
    return
  }

  // 检查文件大小（限制2MB）
  if (file.size > 2 * 1024 * 1024) {
    alert('图片文件不能超过2MB！')
    return
  }

  engineForm.value.iconFile = file

  // 创建预览
  const reader = new FileReader()
  reader.onload = (e) => {
    engineForm.value.iconPreview = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

// 清除图标
const clearIcon = () => {
  engineForm.value.iconFile = null
  engineForm.value.iconPreview = ''
  if (iconFileInput.value) {
    iconFileInput.value.value = ''
  }
}

// 保存搜索引擎
const saveEngine = async () => {
  if (!isEngineFormValid.value) return

  try {
    isEditingEngine.value = true // 开始编辑

    if (editingEngineIndex.value === -1) {
      // 添加新搜索引擎
      if (engineForm.value.iconFile) {
        await store.addSearchEngineWithIcon(
          {
            name: engineForm.value.name,
            url: engineForm.value.url,
            placeholder: engineForm.value.placeholder,
            isDefault: engineForm.value.isDefault,
          },
          engineForm.value.iconFile,
        )
      } else {
        await store.addSearchEngine({
          name: engineForm.value.name,
          url: engineForm.value.url,
          placeholder: engineForm.value.placeholder,
          isDefault: engineForm.value.isDefault,
        })
      }
    } else {
      // 更新现有搜索引擎
      const existingEngine = localSettings.value.search.engines[editingEngineIndex.value]

      if (engineForm.value.iconFile) {
        await store.updateSearchEngineWithIcon(
          existingEngine.id,
          {
            name: engineForm.value.name,
            url: engineForm.value.url,
            placeholder: engineForm.value.placeholder,
            isDefault: engineForm.value.isDefault,
          },
          engineForm.value.iconFile,
        )
      } else {
        await store.updateSearchEngine(existingEngine.id, {
          name: engineForm.value.name,
          url: engineForm.value.url,
          placeholder: engineForm.value.placeholder,
          isDefault: engineForm.value.isDefault,
        })
      }
    }

    // 在 store 操作完成后，同步更新本地设置
    localSettings.value.search.engines = [...store.settings.search.engines]

    closeEngineModal()
    alert(editingEngineIndex.value === -1 ? '搜索引擎添加成功！' : '搜索引擎更新成功！')
  } catch (error) {
    console.error('Failed to save search engine:', error)
    alert('保存搜索引擎失败，请重试')
  } finally {
    isEditingEngine.value = false // 结束编辑
  }
}

const editSearchEngine = (index: number) => {
  const engine = localSettings.value.search.engines[index]
  if (!engine) return

  isEditingEngine.value = true
  editingEngineIndex.value = index
  engineForm.value = {
    name: engine.name,
    url: engine.url,
    placeholder: engine.placeholder || '',
    isDefault: engine.isDefault || false,
    iconFile: null,
    iconPreview: engine.icon || '',
  }
}

const removeSearchEngine = async (index: number) => {
  if (localSettings.value.search.engines.length <= 1) {
    alert('至少需要保留一个搜索引擎！')
    return
  }

  const engine = localSettings.value.search.engines[index]
  if (!confirm(`确定要删除"${engine.name}"搜索引擎吗？`)) {
    return
  }

  try {
    isEditingEngine.value = true // 开始编辑

    await store.deleteSearchEngine(engine.id)

    // 在 store 操作完成后，同步更新本地设置
    localSettings.value.search.engines = [...store.settings.search.engines]

    // 同步默认搜索引擎ID
    localSettings.value.search.defaultEngineId = store.settings.search.defaultEngineId

    alert('搜索引擎删除成功！')
  } catch (error) {
    console.error('Failed to delete search engine:', error)
    alert('删除搜索引擎失败，请重试')
  } finally {
    isEditingEngine.value = false // 结束编辑
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #eee;
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f5f5f5;
  color: #666;
}

.modal-body {
  padding: 24px;
}

.settings-section {
  margin-bottom: 32px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.settings-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  border-bottom: 1px solid #eee;
  padding-bottom: 8px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group input[type='range'] {
  padding: 0;
  height: 6px;
  background: #ddd;
  border-radius: 3px;
  border: none;
}

.form-group input[type='color'] {
  width: 60px;
  height: 40px;
  padding: 0;
  border: 1px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
}

.radio-group {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.radio-item {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-weight: normal !important;
  margin-bottom: 0 !important;
}

.radio-item input[type='radio'] {
  width: auto;
  margin: 0;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-weight: normal !important;
  margin-bottom: 0 !important;
}

.checkbox-item input[type='checkbox'] {
  width: auto;
  margin: 0;
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-right: 8px;
  margin-bottom: 8px;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5a6fd8;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-secondary {
  background: #f5f5f5;
  color: #666;
}

.btn-secondary:hover {
  background: #e8e8e8;
}

.btn-danger {
  background: #e74c3c;
  color: white;
}

.btn-danger:hover {
  background: #d32f2f;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(211, 47, 47, 0.3);
}

.file-upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.btn-upload {
  padding: 8px 16px;
  border: 1px solid #667eea;
  border-radius: 6px;
  background: white;
  color: #667eea;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-upload:hover {
  background: #667eea;
  color: white;
}

.uploaded-file-name {
  font-size: 13px;
  color: #666;
  background: #e8f4f8;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #d0e8f0;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid #eee;
  background: #fafafa;
  border-bottom-left-radius: 12px;
  border-bottom-right-radius: 12px;
}

/* 搜索引擎设置样式 */
.search-engines-config {
  margin-top: 16px;
}

.search-engines-list {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  margin-bottom: 12px;
  overflow: hidden;
}

.search-engine-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s ease;
}

.search-engine-item:last-child {
  border-bottom: none;
}

.search-engine-item:hover {
  background-color: #f8f9fa;
}

.engine-info {
  display: flex;
  align-items: center;
  flex: 1;
}

.engine-icon-small {
  width: 24px;
  height: 24px;
  margin-right: 12px;
  border-radius: 4px;
  object-fit: cover;
}

.engine-icon-fallback-small {
  width: 24px;
  height: 24px;
  margin-right: 12px;
  border-radius: 4px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}

.engine-details {
  flex: 1;
  min-width: 0;
}

.engine-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.engine-url {
  font-size: 12px;
  color: #666;
  word-break: break-all;
}

.engine-actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  background: none;
  border: none;
  padding: 6px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-size: 14px;
}

.btn-icon:hover:not(:disabled) {
  background-color: #f0f0f0;
}

.btn-icon:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.btn-full-width {
  width: 100%;
  justify-content: center;
}

/* 搜索引擎模态框样式 */
.engine-modal {
  max-width: 500px;
}

.icon-upload-area {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 2px dashed #ddd;
  border-radius: 8px;
  background: #f9f9f9;
}

.icon-preview {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border: 1px solid #ddd;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.icon-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
}

.icon-controls {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

/* 默认搜索引擎设置样式 */
.default-engine-setting {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background: #f8f9fa;
  transition: all 0.3s ease;
}

.default-engine-setting:hover {
  border-color: #667eea;
  background: #f0f2ff;
}

.setting-info {
  flex: 1;
}

.setting-title {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.setting-description {
  color: #666;
  font-size: 12px;
  line-height: 1.4;
}

/* 开关样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 24px;
  cursor: pointer;
}

.switch-input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  border-radius: 24px;
  transition: 0.3s;
}

.switch-slider:before {
  position: absolute;
  content: '';
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  border-radius: 50%;
  transition: 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.switch-input:checked + .switch-slider {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.switch-input:focus + .switch-slider {
  box-shadow: 0 0 4px rgba(102, 126, 234, 0.5);
}

.switch-input:checked + .switch-slider:before {
  transform: translateX(24px);
}

/* 响应式设计 */
@media (max-width: 600px) {
  .modal-content {
    margin: 10px;
    width: calc(100% - 20px);
  }

  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 16px;
  }

  .settings-section {
    margin-bottom: 24px;
  }

  .search-engine-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .engine-info {
    width: 100%;
  }

  .engine-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .icon-upload-area {
    flex-direction: column;
    text-align: center;
  }

  .icon-controls {
    flex-direction: row;
    justify-content: center;
  }

  .default-engine-setting {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    text-align: left;
  }

  .setting-info {
    width: 100%;
  }

  .switch {
    align-self: flex-end;
  }
}
</style>
