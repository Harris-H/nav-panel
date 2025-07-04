<template>
  <Teleport to="body">
    <div v-if="show" class="modal-overlay" @click="close">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>{{ isEdit ? '编辑分组' : '创建分组' }}</h2>
          <button @click="close" class="close-btn">×</button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="groupName">分组名称 *</label>
              <input
                id="groupName"
                v-model="formData.name"
                type="text"
                placeholder="请输入分组名称"
                required
                maxlength="20"
              />
            </div>

            <div class="form-group">
              <label for="groupColor">分组颜色</label>
              <div class="color-picker">
                <input id="groupColor" v-model="formData.color" type="color" class="color-input" />
                <div class="color-presets">
                  <button
                    v-for="color in colorPresets"
                    :key="color"
                    type="button"
                    class="color-preset"
                    :style="{ backgroundColor: color }"
                    @click="formData.color = color"
                  ></button>
                </div>
              </div>
            </div>

            <div class="form-group">
              <label for="groupIcon">分组图标</label>
              <div class="icon-picker">
                <input
                  id="groupIcon"
                  v-model="formData.icon"
                  type="text"
                  placeholder="输入图标字符或emoji"
                  maxlength="2"
                />
                <div class="icon-presets">
                  <button
                    v-for="icon in iconPresets"
                    :key="icon"
                    type="button"
                    class="icon-preset"
                    @click="formData.icon = icon"
                  >
                    {{ icon }}
                  </button>
                </div>
              </div>
            </div>

            <div class="modal-actions">
              <button type="button" @click="close" class="btn-secondary">取消</button>
              <button type="submit" class="btn-primary" :disabled="!formData.name.trim()">
                {{ isEdit ? '更新' : '创建' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import type { Group } from '@/types'

interface Props {
  show: boolean
  group?: Group | null
}

interface Emits {
  (e: 'close'): void
  (e: 'submit', data: { name: string; color?: string; icon?: string }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isEdit = computed(() => !!props.group)

const formData = reactive({
  name: '',
  color: '#667eea',
  icon: '📁',
})

const colorPresets = [
  '#667eea',
  '#764ba2',
  '#f093fb',
  '#f5576c',
  '#4facfe',
  '#00f2fe',
  '#43e97b',
  '#38f9d7',
  '#ffecd2',
  '#fcb69f',
  '#a8edea',
  '#fed6e3',
  '#fad0c4',
  '#ffd1ff',
  '#c2e9fb',
  '#a1c4fd',
]

const iconPresets = [
  '📁',
  '🏠',
  '💼',
  '🎮',
  '📚',
  '🛒',
  '⚡',
  '🎨',
  '🔧',
  '📊',
  '🌟',
  '🔥',
  '💡',
  '🎯',
  '🚀',
  '💎',
]

// 监听 props.group 变化，更新表单数据
watch(
  () => props.group,
  (newGroup) => {
    if (newGroup) {
      formData.name = newGroup.name
      formData.color = newGroup.color || '#667eea'
      formData.icon = newGroup.icon || '📁'
    } else {
      // 重置表单
      formData.name = ''
      formData.color = '#667eea'
      formData.icon = '📁'
    }
  },
  { immediate: true },
)

const close = () => {
  emit('close')
}

const handleSubmit = () => {
  if (!formData.name.trim()) return

  emit('submit', {
    name: formData.name.trim(),
    color: formData.color,
    icon: formData.icon || undefined,
  })

  // 让父组件决定何时关闭模态框，而不是在这里自动关闭
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.9) 100%);
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  width: 90%;
  max-width: 400px;
  max-height: 90vh;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 24px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 2rem;
  color: #666;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.1);
  color: #333;
}

.modal-body {
  padding: 0 24px 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.form-group input[type='text'] {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  font-size: 14px;
  transition: all 0.2s ease;
  background: rgba(255, 255, 255, 0.8);
}

.form-group input[type='text']:focus {
  outline: none;
  border-color: #667eea;
  background: rgba(255, 255, 255, 1);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.color-picker {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.color-input {
  width: 60px;
  height: 40px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  cursor: pointer;
  background: none;
}

.color-presets {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
}

.color-preset {
  width: 32px;
  height: 32px;
  border: 2px solid rgba(255, 255, 255, 0.8);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.color-preset:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.icon-picker {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.icon-picker input {
  width: 80px;
  text-align: center;
  font-size: 18px;
}

.icon-presets {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
}

.icon-preset {
  width: 36px;
  height: 36px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.icon-preset:hover {
  background: rgba(102, 126, 234, 0.1);
  border-color: #667eea;
  transform: scale(1.1);
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

.btn-primary,
.btn-secondary {
  padding: 12px 24px;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(0, 0, 0, 0.1);
  color: #333;
}

.btn-secondary:hover {
  background: rgba(0, 0, 0, 0.2);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 0 auto;
  }

  .color-presets,
  .icon-presets {
    grid-template-columns: repeat(6, 1fr);
  }

  .modal-header {
    padding: 20px 20px 0;
  }

  .modal-body {
    padding: 0 20px 20px;
  }
}
</style>
