<template>
  <div v-if="store.isAddCardModalOpen" class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>{{ isEditing ? '编辑网站' : '添加网站' }}</h2>
        <button class="close-btn" @click="closeModal">×</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-form">
        <div class="form-group">
          <label for="url">网站地址 *</label>
          <div class="url-input-wrapper">
            <input
              id="url"
              v-model="form.url"
              type="url"
              placeholder="https://example.com"
              required
              @blur="handleUrlBlur"
              @input="handleUrlInput"
            />
            <button
              v-if="form.url && !isLoadingTitle"
              type="button"
              @click="fetchWebsiteTitle"
              class="fetch-title-btn"
              title="获取网站标题"
            >
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path
                  fill="currentColor"
                  d="M17.65,6.35C16.2,4.9 14.21,4 12,4A8,8 0 0,0 4,12A8,8 0 0,0 12,20C15.73,20 18.84,17.45 19.73,14H17.65C16.83,16.33 14.61,18 12,18A6,6 0 0,1 6,12A6,6 0 0,1 12,6C13.66,6 15.14,6.69 16.22,7.78L13,11H20V4L17.65,6.35Z"
                />
              </svg>
            </button>
            <div v-if="isLoadingTitle" class="loading-indicator">
              <div class="spinner"></div>
            </div>
          </div>
          <small>请输入完整的网站地址，系统将自动获取网站名称</small>
        </div>

        <div class="form-group">
          <label for="name">网站名称 *</label>
          <div class="name-input-wrapper">
            <input
              id="name"
              v-model="form.name"
              type="text"
              placeholder="请输入网站名称"
              required
            />
            <span v-if="autoFetched" class="auto-fetched-indicator" title="已自动获取">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path
                  fill="currentColor"
                  d="M21,7L9,19L3.5,13.5L4.91,12.09L9,16.17L19.59,5.59L21,7Z"
                />
              </svg>
            </span>
          </div>
          <small v-if="autoFetched" class="auto-fetched-hint">已自动获取，您可以手动修改</small>
        </div>

        <div class="form-group">
          <label>网站图标</label>
          <div class="icon-input-section">
            <!-- 图标上传区域 -->
            <div
              class="icon-upload-area"
              :class="{
                'drag-over': isDragOver,
                'has-image': iconPreview,
              }"
              @drop="handleDrop"
              @dragover.prevent="handleDragOver"
              @dragenter.prevent="handleDragEnter"
              @dragleave.prevent="handleDragLeave"
              @click="triggerFileInput"
            >
              <img v-if="iconPreview" :src="iconPreview" alt="图标预览" class="preview-image" />
              <div v-else class="upload-placeholder">
                <svg viewBox="0 0 24 24" width="32" height="32">
                  <path
                    fill="currentColor"
                    d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"
                  />
                </svg>
                <div class="upload-text">
                  <span class="primary-text">{{
                    isDragOver ? '释放图片' : '点击或拖拽上传图标'
                  }}</span>
                  <span class="secondary-text">支持 PNG、JPG、GIF、SVG、WebP</span>
                </div>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                @change="handleFileUpload"
                hidden
              />
            </div>

            <!-- URL输入和操作按钮 -->
            <div class="icon-options">
              <div class="url-input-wrapper">
                <input
                  v-model="form.icon"
                  type="url"
                  placeholder="或输入图标URL地址"
                  @input="handleIconUrlInput"
                  class="icon-url-input"
                />
              </div>
              <button
                v-if="iconPreview"
                type="button"
                @click="clearIcon"
                class="clear-icon-btn"
                title="清除图标"
              >
                <svg viewBox="0 0 24 24" width="16" height="16">
                  <path
                    fill="currentColor"
                    d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label for="description">网站描述</label>
          <textarea
            id="description"
            v-model="form.description"
            placeholder="简单描述一下这个网站..."
            rows="3"
          ></textarea>
        </div>

        <div class="form-group">
          <label for="category">分类</label>
          <input
            id="category"
            v-model="form.category"
            type="text"
            placeholder="如：工具、娱乐、学习等"
          />
        </div>

        <div class="modal-actions">
          <button type="button" class="btn-cancel" @click="closeModal">取消</button>
          <button type="submit" class="btn-submit">
            {{ isEditing ? '保存修改' : '添加网站' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useAppStore } from '@/stores/app'

const store = useAppStore()

const form = ref({
  name: '',
  url: '',
  icon: '',
  description: '',
  category: '',
})

const fileInput = ref<HTMLInputElement>()
const iconPreview = ref<string>('')
const isDragOver = ref(false)
const isLoadingTitle = ref(false)
const autoFetched = ref(false)

const isEditing = computed(() => !!store.editingCard)

// 先定义 resetForm 函数
const resetForm = () => {
  form.value = {
    name: '',
    url: '',
    icon: '',
    description: '',
    category: '',
  }
  iconPreview.value = ''
  autoFetched.value = false
}

// 然后再使用 watch
watch(
  () => store.editingCard,
  (card) => {
    if (card) {
      form.value = {
        name: card.name,
        url: card.url,
        icon: card.icon || '',
        description: card.description || '',
        category: card.category || '',
      }
      iconPreview.value = card.icon || ''
      autoFetched.value = false
    } else {
      resetForm()
    }
  },
  { immediate: true },
)

// 获取网站标题的函数
const fetchWebsiteTitle = async () => {
  if (!form.value.url || isLoadingTitle.value) return

  let url = form.value.url.trim()
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = `https://${url}`
  }

  try {
    isLoadingTitle.value = true

    // 使用多个CORS代理服务作为备用方案
    const proxyServices = [
      `https://api.allorigins.win/get?url=${encodeURIComponent(url)}`,
      `https://corsproxy.io/?${encodeURIComponent(url)}`,
      `https://cors-anywhere.herokuapp.com/${url}`,
    ]

    let title = ''

    // 尝试多个代理服务
    for (const proxyUrl of proxyServices) {
      try {
        const response = await fetch(proxyUrl, {
          method: 'GET',
          headers: {
            Accept: 'application/json, text/html',
          },
          // 设置超时
          signal: AbortSignal.timeout(8000),
        })

        if (response.ok) {
          let html = ''

          // 处理不同代理服务的响应格式
          if (proxyUrl.includes('allorigins.win')) {
            const data = await response.json()
            html = data.contents
          } else {
            html = await response.text()
          }

          // 提取网站标题
          const titleMatch = html.match(/<title[^>]*>([^<]+)<\/title>/i)

          if (titleMatch && titleMatch[1]) {
            title = titleMatch[1].trim()

            // 清理标题，移除一些常见的后缀
            title = title.replace(/\s*[-|–—]\s*.*$/, '') // 移除 - 或 | 后的内容
            title = title.replace(/\s*_.*$/, '') // 移除 _ 后的内容
            title = title.replace(/\s*\|.*$/, '') // 移除 | 后的内容
            title = title.replace(/\s*-.*$/, '') // 移除单个 - 后的内容

            // 解码HTML实体
            const tempDiv = document.createElement('div')
            tempDiv.innerHTML = title
            title = tempDiv.textContent || tempDiv.innerText || title

            break // 成功获取标题，跳出循环
          }
        }
      } catch (proxyError) {
        console.log(`代理服务 ${proxyUrl} 失败:`, proxyError)
        continue // 尝试下一个代理服务
      }
    }

    // 如果获取到标题且名称字段为空，则填充
    if (title && !form.value.name.trim()) {
      form.value.name = title
      autoFetched.value = true
    } else if (!title) {
      // 如果所有代理都失败，尝试从URL提取网站名称
      const urlObj = new URL(url)
      const hostname = urlObj.hostname.replace(/^www\./, '')
      const siteName = hostname.split('.')[0]

      if (siteName && !form.value.name.trim()) {
        // 将域名首字母大写作为默认名称
        const defaultName = siteName.charAt(0).toUpperCase() + siteName.slice(1)
        form.value.name = defaultName
        autoFetched.value = true
      }
    }
  } catch (error) {
    console.log('获取网站标题失败:', error)

    // 作为最后的备用方案，从URL提取网站名称
    try {
      const urlObj = new URL(url)
      const hostname = urlObj.hostname.replace(/^www\./, '')
      const siteName = hostname.split('.')[0]

      if (siteName && !form.value.name.trim()) {
        const defaultName = siteName.charAt(0).toUpperCase() + siteName.slice(1)
        form.value.name = defaultName
        autoFetched.value = true
      }
    } catch (urlError) {
      console.log('URL解析失败:', urlError)
    }
  } finally {
    isLoadingTitle.value = false
  }
}

// URL输入变化处理
const handleUrlInput = () => {
  // 清除自动获取标志
  autoFetched.value = false
}

// URL失去焦点时尝试获取标题
const handleUrlBlur = () => {
  if (form.value.url && !form.value.name.trim() && !isEditing.value) {
    fetchWebsiteTitle()
  }
}

const processImageFile = (file: File) => {
  // 验证文件类型
  const validTypes = [
    'image/png',
    'image/jpeg',
    'image/jpg',
    'image/gif',
    'image/svg+xml',
    'image/webp',
  ]
  if (!validTypes.includes(file.type)) {
    alert('请上传有效的图片文件（PNG、JPG、GIF、SVG、WebP）')
    return
  }

  // 验证文件大小（2MB）
  if (file.size > 2 * 1024 * 1024) {
    alert('图片文件大小不能超过2MB')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    const base64 = e.target?.result as string
    form.value.icon = base64
    iconPreview.value = base64
  }
  reader.readAsDataURL(file)
}

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    processImageFile(file)
  }
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false

  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    processImageFile(files[0])
  }
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
}

const handleDragEnter = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleIconUrlInput = () => {
  if (form.value.icon && form.value.icon.startsWith('http')) {
    iconPreview.value = form.value.icon
  } else if (!form.value.icon) {
    iconPreview.value = ''
  }
}

const clearIcon = () => {
  form.value.icon = ''
  iconPreview.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const closeModal = () => {
  store.closeAddCardModal()
  resetForm()
}

const handleSubmit = () => {
  if (!form.value.name.trim() || !form.value.url.trim()) {
    alert('请填写网站名称和地址')
    return
  }

  // 确保URL格式正确
  let url = form.value.url.trim()
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = `https://${url}`
  }

  const siteData = {
    id: store.editingCard?.id || Date.now().toString(),
    name: form.value.name.trim(),
    url: url,
    icon: form.value.icon || '',
    description: form.value.description.trim(),
    category: form.value.category.trim(),
  }

  if (isEditing.value) {
    store.updateSite(siteData)
  } else {
    store.addSite(siteData)
  }

  closeModal()
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
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 20px 0;
  border-bottom: 1px solid #eee;
  margin-bottom: 20px;
}

.modal-header h2 {
  margin: 0;
  color: #333;
  font-size: 1.5rem;
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
  transition: background-color 0.2s;
}

.close-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.modal-form {
  padding: 0 20px 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: #333;
  font-weight: 500;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.form-group small {
  display: block;
  margin-top: 5px;
  color: #666;
  font-size: 12px;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

/* URL输入相关样式 */
.url-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.url-input-wrapper input {
  flex: 1;
  padding-right: 40px;
}

.fetch-title-btn {
  position: absolute;
  right: 8px;
  background: none;
  border: none;
  cursor: pointer;
  color: #667eea;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.fetch-title-btn:hover {
  background: #f0f4ff;
  color: #5a67d8;
}

.loading-indicator {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #f3f3f3;
  border-top: 2px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* 网站名称输入相关样式 */
.name-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.name-input-wrapper input {
  flex: 1;
  padding-right: 40px;
}

.auto-fetched-indicator {
  position: absolute;
  right: 8px;
  color: #28a745;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
}

.auto-fetched-hint {
  color: #28a745 !important;
  font-size: 12px;
  margin-top: 5px;
}

/* 图标上传相关样式 */
.icon-input-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.icon-upload-area {
  width: 100%;
  min-height: 120px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  position: relative;
}

/* 无图标时的样式 */
.icon-upload-area:not(.has-image) {
  border: 2px dashed #d1d5db;
  background: linear-gradient(135deg, #f9fafb 0%, #f3f4f6 100%);
}

.icon-upload-area:not(.has-image):hover {
  border-color: #667eea;
  background: linear-gradient(135deg, #f0f4ff 0%, #e6f0ff 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.15);
}

.icon-upload-area:not(.has-image).drag-over {
  border-color: #667eea !important;
  background: linear-gradient(135deg, #e6f0ff 0%, #dde7ff 100%) !important;
  transform: scale(1.02);
  box-shadow: 0 12px 35px rgba(102, 126, 234, 0.25);
}

/* 有图标时的样式 */
.icon-upload-area.has-image {
  border: none;
  background: transparent;
  min-height: 80px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.icon-upload-area.has-image:hover {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.icon-upload-area.has-image.drag-over {
  opacity: 0.8;
  transform: scale(1.05);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
}

.preview-image {
  max-width: 100%;
  max-height: 80px;
  width: auto;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
  display: block;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
  padding: 20px;
  pointer-events: none;
}

.upload-placeholder svg {
  color: #9ca3af;
  opacity: 0.8;
  transition: all 0.3s ease;
}

.icon-upload-area:hover .upload-placeholder svg {
  color: #667eea;
  opacity: 1;
  transform: translateY(-2px);
}

.upload-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.primary-text {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.secondary-text {
  font-size: 12px;
  color: #9ca3af;
}

.icon-upload-area:hover .primary-text {
  color: #667eea;
}

/* URL输入和操作按钮 */
.icon-options {
  display: flex;
  align-items: center;
  gap: 8px;
}

.url-input-wrapper {
  flex: 1;
}

.icon-url-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  box-sizing: border-box;
  transition: all 0.2s ease;
  background: #f9fafb;
}

.icon-url-input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.icon-url-input::placeholder {
  color: #9ca3af;
  font-style: italic;
}

.clear-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #dc2626;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.clear-icon-btn:hover {
  background: #fecaca;
  border-color: #f87171;
  color: #b91c1c;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.2);
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.btn-cancel,
.btn-submit {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background: #f5f5f5;
  color: #666;
}

.btn-cancel:hover {
  background: #ebebeb;
  color: #333;
}

.btn-submit {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-submit:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

@media (max-width: 768px) {
  .modal-content {
    margin: 0;
    border-radius: 0;
    height: 100vh;
    max-height: none;
  }

  .icon-upload-area {
    min-height: 100px;
  }

  .icon-upload-area:not(.has-image) {
    border-width: 1px;
  }

  .upload-placeholder {
    padding: 16px;
    gap: 8px;
  }

  .upload-placeholder svg {
    width: 28px;
    height: 28px;
  }

  .primary-text {
    font-size: 13px;
  }

  .secondary-text {
    font-size: 11px;
  }

  .icon-options {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .clear-icon-btn {
    align-self: flex-end;
    width: 40px;
    height: 40px;
  }

  .modal-actions {
    flex-direction: column;
  }

  .btn-cancel,
  .btn-submit {
    width: 100%;
  }
}
</style>
