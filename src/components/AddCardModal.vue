<template>
  <div v-if="store.isAddCardModalOpen" class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>{{ isEditing ? '编辑网站' : '添加网站' }}</h2>
        <button class="close-btn" @click="closeModal">×</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-form">
        <div class="form-group">
          <label for="name">网站名称 *</label>
          <input id="name" v-model="form.name" type="text" placeholder="请输入网站名称" required />
        </div>

        <div class="form-group">
          <label for="url">网站地址 *</label>
          <input id="url" v-model="form.url" type="url" placeholder="https://example.com" required />
          <small>请输入完整的网站地址，包含 http:// 或 https://</small>
        </div>

        <div class="form-group">
          <label>网站图标</label>
          <div class="icon-input-section">
            <!-- 图标预览 -->
            <div class="icon-preview">
              <div class="icon-preview-box" :class="{
                'drag-over': isDragOver,
                'has-image': iconPreview
              }" @drop="handleDrop" @dragover.prevent="handleDragOver" @dragenter.prevent="handleDragEnter"
                @dragleave.prevent="handleDragLeave" @click="triggerFileInput">
                <img v-if="iconPreview" :src="iconPreview" alt="图标预览" class="preview-image" />
                <div v-else class="preview-placeholder">
                  <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor"
                      d="M21,19V5C21,3.89 20.1,3 19,3H5A2,2 0 0,0 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19M19,5V19H5V5H19M13.96,12.29L11.21,15.83L9.25,13.47L6.5,17H17.5L13.96,12.29Z" />
                  </svg>
                  <span>{{ isDragOver ? '释放图片' : '点击或拖拽上传' }}</span>
                </div>
              </div>
            </div>

            <!-- 上传选项 -->
            <div class="icon-input-options">
              <div class="upload-section">
                <label class="upload-btn" for="icon-upload">
                  <svg viewBox="0 0 24 24" width="16" height="16">
                    <path fill="currentColor"
                      d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
                  </svg>
                  上传图标
                </label>
                <input id="icon-upload" ref="fileInput" type="file" accept="image/*" @change="handleFileUpload"
                  hidden />
              </div>

              <div class="or-divider">或</div>

              <div class="url-section">
                <input v-model="form.icon" type="url" placeholder="输入图标URL地址" @input="handleUrlInput"
                  class="icon-url-input" />
              </div>

              <button v-if="iconPreview" type="button" @click="clearIcon" class="clear-icon-btn">
                清除图标
              </button>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label for="description">网站描述</label>
          <textarea id="description" v-model="form.description" placeholder="简单描述一下这个网站..." rows="3"></textarea>
        </div>

        <div class="form-group">
          <label for="category">分类</label>
          <input id="category" v-model="form.category" type="text" placeholder="如：工具、娱乐、学习等" />
        </div>

        <div class="modal-actions">
          <button type="button" class="btn-cancel" @click="closeModal">
            取消
          </button>
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
  category: ''
})

const fileInput = ref<HTMLInputElement>()
const iconPreview = ref<string>('')
const isDragOver = ref(false)

const isEditing = computed(() => !!store.editingCard)

// 先定义 resetForm 函数
const resetForm = () => {
  form.value = {
    name: '',
    url: '',
    icon: '',
    description: '',
    category: ''
  }
  iconPreview.value = ''
}

// 然后再使用 watch
watch(() => store.editingCard, (card) => {
  if (card) {
    form.value = {
      name: card.name,
      url: card.url,
      icon: card.icon || '',
      description: card.description || '',
      category: card.category || ''
    }
    iconPreview.value = card.icon || ''
  } else {
    resetForm()
  }
}, { immediate: true })

const processImageFile = (file: File) => {
  // 验证文件类型
  const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/gif', 'image/svg+xml', 'image/webp']
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

const handleUrlInput = () => {
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
    category: form.value.category.trim()
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

/* 图标上传相关样式 */
.icon-input-section {
  display: flex;
  gap: 15px;
  align-items: flex-start;
}

.icon-preview {
  flex-shrink: 0;
}

.icon-preview-box {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
  position: relative;
  background: transparent;
  /* 确保背景透明 */
}

/* 有图标时的样式 */
.icon-preview-box.has-image {
  border: none !important;
  background: transparent !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
  /* 移除内边距 */
}

.icon-preview-box.has-image .preview-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 8px;
  display: block;
  max-width: 100%;
  max-height: 100%;
}

/* 无图标时的样式 */
.icon-preview-box:not(.has-image) {
  border: 2px dashed #ddd;
  background: #fafafa;
}

.icon-preview-box:not(.has-image):hover {
  border-color: #667eea;
  background: #f0f4ff;
}

.icon-preview-box:not(.has-image).drag-over {
  border: 2px solid #667eea !important;
  background: #e6f0ff !important;
}

/* 有图标状态 - 隐藏边框 */
.icon-preview-box.has-image {
  border: none;
  background: transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.icon-preview-box.has-image:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.icon-preview-box .preview-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  /* 改为 contain 确保图片完整显示 */
  border-radius: 8px;
  display: block;
  /* 添加 block 显示 */
}

/* 或者可以尝试这个样式 */
.icon-preview-box.has-image .preview-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  object-position: center;
  border-radius: 8px;
  display: block;
  vertical-align: middle;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  text-align: center;
  color: #999;
  font-size: 11px;
  padding: 4px;
  pointer-events: none;
}

.preview-placeholder svg {
  opacity: 0.5;
}

/* 拖拽时的样式 */
.icon-preview-box.has-image.drag-over {
  opacity: 0.8;
  transform: scale(1.02);
}

/* 如果浏览器不支持 :has 选择器，添加一个类来控制 */
.icon-preview-box.has-image {
  border: none !important;
  background: transparent !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.icon-preview-box.has-image:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: none !important;
  background: transparent !important;
}

.icon-input-options {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.upload-section {
  display: flex;
  align-items: center;
}

.upload-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.upload-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.upload-btn svg {
  flex-shrink: 0;
}

.or-divider {
  position: relative;
  text-align: center;
  color: #999;
  font-size: 12px;
  margin: 5px 0;
}

.or-divider::before,
.or-divider::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 30%;
  height: 1px;
  background: #ddd;
}

.or-divider::before {
  left: 0;
}

.or-divider::after {
  right: 0;
}

.url-section {
  flex: 1;
}

.icon-url-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
}

.icon-url-input:focus {
  outline: none;
  border-color: #667eea;
}

.clear-icon-btn {
  padding: 6px 12px;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  color: #666;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.clear-icon-btn:hover {
  background: #ebebeb;
  color: #333;
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

  .icon-input-section {
    flex-direction: column;
    gap: 10px;
  }

  .icon-preview {
    align-self: center;
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
