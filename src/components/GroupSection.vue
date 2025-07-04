<template>
  <div class="group-section" :class="{ collapsed: group.isCollapsed }">
    <!-- 分组头部 -->
    <div
      class="group-header"
      :style="{ borderLeftColor: group.color || '#667eea' }"
      @click="toggleCollapse"
      @contextmenu.prevent="handleRightClick"
    >
      <div class="group-info">
        <div class="group-icon">{{ group.icon || '📁' }}</div>
        <div class="group-name">{{ group.name }}</div>
        <div class="group-count">({{ websites?.length || 0 }})</div>
      </div>

      <div class="group-actions">
        <button
          class="toggle-btn"
          :class="{ collapsed: group.isCollapsed }"
          @click.stop="toggleCollapse"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="6,9 12,15 18,9"></polyline>
          </svg>
        </button>
      </div>
    </div>

    <!-- 分组内容 -->
    <div class="group-content" :class="{ collapsed: group.isCollapsed }">
      <div
        class="drop-zone"
        :class="{ 'drop-active': isDropTarget }"
        @dragover.prevent="handleDragOver"
        @dragleave="handleDragLeave"
        @drop="handleDrop"
      >
        <div v-if="!websites || websites.length === 0" class="empty-group">
          <div class="empty-icon">📂</div>
          <div class="empty-text">暂无网站</div>
          <div class="empty-hint">拖拽网站到这里</div>
        </div>

        <div v-else class="websites-grid" :style="gridStyle">
          <template v-for="(site, index) in websites" :key="site.id">
            <div
              class="site-card"
              :style="cardStyle"
              @contextmenu.prevent="handleSiteRightClick($event, site)"
              @mousedown="handleSiteMouseDown($event, index)"
            >
              <div class="site-icon" style="pointer-events: none">
                <img
                  v-if="site.icon"
                  :src="site.icon"
                  :alt="site.name"
                  style="pointer-events: none"
                />
                <span v-else style="pointer-events: none">{{
                  site.name.charAt(0).toUpperCase()
                }}</span>
              </div>
              <div class="site-name" style="pointer-events: none">{{ site.name }}</div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- 右键菜单 - 使用 Teleport 渲染到 body -->
    <Teleport to="body">
      <div
        v-if="contextMenu.show"
        ref="contextMenuRef"
        class="context-menu"
        :style="contextMenuStyle"
        @click.stop
      >
        <div class="context-menu-item" @click="editGroup">
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="m18.5 2.5 a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
          编辑分组
        </div>
        <div class="context-menu-item delete" @click="deleteGroup">
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="3,6 5,6 21,6"></polyline>
            <path
              d="m19,6v14a2,2 0 0 1-2,2H7a2,2 0 0 1-2-2V6m3,0V4a2,2 0 0 1 2-2h4a2,2 0 0 1 2,2v2"
            ></path>
            <line x1="10" y1="11" x2="10" y2="17"></line>
            <line x1="14" y1="11" x2="14" y2="17"></line>
          </svg>
          删除分组
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
import type { Group, Website } from '@/types'

interface Props {
  group: Group
  websites: Website[] | null | undefined
  gridStyle: Record<string, string | number>
  cardStyle: Record<string, string | number>
}

interface Emits {
  (e: 'toggle-collapse', groupId: string): void
  (e: 'edit-group', group: Group): void
  (e: 'delete-group', groupId: string): void
  (e: 'site-right-click', event: MouseEvent, site: Website): void
  (e: 'site-mouse-down', event: MouseEvent, index: number, groupId: string): void
  (e: 'drop-website', websiteId: string, targetGroupId: string, position?: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const contextMenuRef = ref<HTMLElement>()
const isDropTarget = ref(false)

// 右键菜单状态
const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
})

// 右键菜单样式
const contextMenuStyle = computed(() => ({
  left: `${contextMenu.value.x}px`,
  top: `${contextMenu.value.y}px`,
}))

// 切换折叠状态
const toggleCollapse = () => {
  emit('toggle-collapse', props.group.id)
}

// 处理分组右键点击
const handleRightClick = (event: MouseEvent) => {
  event.preventDefault()
  event.stopPropagation() // 阻止事件冒泡

  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
  }

  // 下一帧调整菜单位置，防止超出屏幕
  nextTick(() => {
    if (contextMenuRef.value) {
      const menuRect = contextMenuRef.value.getBoundingClientRect()
      const viewportWidth = window.innerWidth
      const viewportHeight = window.innerHeight

      // 调整水平位置
      if (contextMenu.value.x + menuRect.width > viewportWidth) {
        contextMenu.value.x = viewportWidth - menuRect.width - 10
      }

      // 调整垂直位置
      if (contextMenu.value.y + menuRect.height > viewportHeight) {
        contextMenu.value.y = viewportHeight - menuRect.height - 10
      }
    }
  })
}

// 编辑分组
const editGroup = () => {
  emit('edit-group', props.group)
  hideContextMenu()
}

// 删除分组
const deleteGroup = () => {
  emit('delete-group', props.group.id)
  hideContextMenu()
}

// 隐藏右键菜单
const hideContextMenu = () => {
  contextMenu.value.show = false
}

// 处理网站右键点击
const handleSiteRightClick = (event: MouseEvent, site: Website) => {
  emit('site-right-click', event, site)
}

// 处理网站鼠标按下
const handleSiteMouseDown = (event: MouseEvent, index: number) => {
  emit('site-mouse-down', event, index, props.group.id)
}

// 拖拽相关
const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDropTarget.value = true
}

const handleDragLeave = (event: DragEvent) => {
  // 只有当鼠标真正离开drop-zone时才隐藏效果
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const x = event.clientX
  const y = event.clientY

  if (x < rect.left || x > rect.right || y < rect.top || y > rect.bottom) {
    isDropTarget.value = false
  }
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDropTarget.value = false

  const websiteId = event.dataTransfer?.getData('text/website-id')
  if (websiteId) {
    emit('drop-website', websiteId, props.group.id)
  }
}

// 处理点击外部区域关闭右键菜单
const handleClickOutside = (event: Event) => {
  if (contextMenu.value.show) {
    const target = event.target as Element
    if (!contextMenuRef.value?.contains(target)) {
      hideContextMenu()
    }
  }
}

// 添加和移除全局事件监听器
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.group-section {
  margin-bottom: 24px;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  border-radius: 16px;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  overflow: hidden;
  transition: transform 0.2s ease;
  will-change: transform;
}

.group-section:hover {
  transform: translateY(-1px);
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.08);
  border-left: 4px solid #667eea;
  cursor: pointer;
  transition: background 0.15s ease;
  user-select: none;
  will-change: background;
}

.group-header:hover {
  background: rgba(255, 255, 255, 0.12);
}

.group-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-icon {
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
}

.group-name {
  font-size: 18px;
  font-weight: 600;
  color: white;
}

.group-count {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.group-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toggle-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  will-change: transform;
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.toggle-btn.collapsed {
  transform: rotate(-90deg);
}

.group-content {
  padding: 8px 20px 16px;
  overflow: hidden;
  transition:
    max-height 0.3s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.25s ease,
    padding 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  opacity: 1;
  max-height: 500px; /* 足够容纳大多数分组内容 */
  will-change: max-height, opacity;
}

.group-content.collapsed {
  opacity: 0;
  max-height: 0;
  padding: 0 20px;
  pointer-events: none;
}

.drop-zone {
  min-height: 80px;
  border: 2px dashed transparent;
  border-radius: 12px;
  transition: all 0.3s ease;
  padding: 12px;
}

.drop-zone.drop-active {
  border-color: rgba(255, 255, 255, 0.5);
  background: rgba(255, 255, 255, 0.05);
}

.empty-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.6;
}

.empty-text {
  color: rgba(255, 255, 255, 0.8);
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.empty-hint {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
}

.websites-grid {
  display: grid;
  gap: 12px;
  transition: gap 0.3s ease;
}

/* 网站卡片样式 */
.site-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.08) 0%, rgba(255, 255, 255, 0.03) 100%);
  border-radius: 12px;
  padding: 12px 8px;
  text-align: center;
  cursor: grab;
  transition:
    transform 0.25s cubic-bezier(0.25, 0.46, 0.45, 0.94),
    box-shadow 0.25s ease,
    background 0.25s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow:
    0 4px 16px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
  position: relative;
  overflow: hidden;
  will-change: transform;
}

/* 简化微光效果，减少性能消耗 */
.site-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left 0.3s ease;
  opacity: 0;
}

.site-card:hover::before {
  left: 100%;
  opacity: 1;
}

.site-card:hover {
  transform: translateY(-4px) scale(1.02);
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.15) 0%, rgba(255, 255, 255, 0.08) 100%);
  box-shadow:
    0 8px 24px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.site-card:active {
  cursor: grabbing;
}

.site-icon {
  width: 32px;
  height: 32px;
  margin: 0 auto 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  overflow: hidden;
  transition: transform 0.2s ease;
  position: relative;
  border: none;
  will-change: transform;
}

.site-icon:not(:has(img)) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.12) 0%, rgba(255, 255, 255, 0.06) 100%);
  backdrop-filter: blur(8px);
  border: none;
  border-radius: 12px;
}

.site-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  background: transparent;
  opacity: 0.9;
  transition: opacity 0.2s ease;
  border: none;
}

.site-card:hover .site-icon {
  transform: scale(1.08);
}

.site-card:hover .site-icon img {
  opacity: 1;
}

.site-name {
  font-size: 11px;
  font-weight: 600;
  word-break: break-word;
  line-height: 1.2;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
  color: white;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  z-index: 9999;
  min-width: 140px;
  overflow: hidden;
  animation: contextMenuSlide 0.15s ease-out;
}

@keyframes contextMenuSlide {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-5px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  cursor: pointer;
  transition: all 0.2s ease;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.context-menu-item:last-child {
  border-bottom: none;
}

.context-menu-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.context-menu-item.delete:hover {
  background: rgba(231, 76, 60, 0.1);
  color: #e74c3c;
}

.context-menu-item svg {
  flex-shrink: 0;
  opacity: 0.7;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .group-section {
    margin-bottom: 16px;
  }

  .group-header {
    padding: 12px 16px;
  }

  .group-content {
    padding: 0 16px 16px;
  }

  .websites-grid {
    grid-template-columns: repeat(auto-fill, minmax(70px, 1fr));
    gap: 10px;
  }

  .site-card {
    padding: 10px 6px;
  }

  .site-icon {
    width: 28px;
    height: 28px;
    margin-bottom: 6px;
  }

  .site-name {
    font-size: 10px;
  }
}
</style>
