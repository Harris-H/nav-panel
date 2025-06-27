<template>
  <div class="home-container" :style="containerStyle">
    <div class="debug-info">调试: 组件已加载</div>

    <header class="app-header">
      <h1>导航面板</h1>
      <div class="header-actions">
        <button @click="store.openAddCardModal()" class="btn-primary">添加网站</button>
        <button @click="store.openSettingsModal()" class="btn-secondary">设置</button>
      </div>
    </header>

    <!-- 搜索框组件 -->
    <SearchBox />

    <main class="main-content">
      <div class="sites-grid" :style="gridStyle">
        <div
          v-for="site in store.sites"
          :key="site.id"
          class="site-card"
          :style="cardStyle"
          @click="openSite(site.url)"
          @contextmenu.prevent="handleRightClick($event, site)"
        >
          <div class="site-icon">
            <img v-if="site.icon" :src="site.icon" :alt="site.name" />
            <span v-else>{{ site.name.charAt(0).toUpperCase() }}</span>
          </div>
          <div class="site-name">{{ site.name }}</div>
        </div>

        <div class="site-card add-card" :style="cardStyle" @click="store.openAddCardModal()">
          <div class="site-icon">+</div>
          <div class="site-name">添加网站</div>
        </div>
      </div>
    </main>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenu.show"
      ref="contextMenuRef"
      class="context-menu"
      :style="contextMenuStyle"
      @click.stop
    >
      <div class="context-menu-item" @click="editSite">
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
        编辑网站
      </div>
      <div class="context-menu-item delete" @click="deleteSite">
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
        删除网站
      </div>
    </div>

    <!-- 添加网站模态框 -->
    <AddCardModal />

    <!-- 设置模态框 -->
    <SettingsModal />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, computed, ref, nextTick } from 'vue'
import { useAppStore } from '@/stores/app'
import AddCardModal from '@/components/AddCardModal.vue'
import SettingsModal from '@/components/SettingsModal.vue'
import SearchBox from '@/components/SearchBox.vue'
import type { Website } from '@/types'

const store = useAppStore()
const contextMenuRef = ref<HTMLElement>()

// 右键菜单状态
const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  site: null as Website | null,
})

// 右键菜单样式
const contextMenuStyle = computed(() => ({
  left: `${contextMenu.value.x}px`,
  top: `${contextMenu.value.y}px`,
}))

const openSite = (url: string) => {
  window.open(url, '_blank')
}

// 处理右键点击
const handleRightClick = (event: MouseEvent, site: Website) => {
  event.preventDefault()

  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    site,
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

// 编辑网站
const editSite = () => {
  if (contextMenu.value.site) {
    store.openEditCardModal(contextMenu.value.site)
  }
  hideContextMenu()
}

// 删除网站
const deleteSite = () => {
  if (contextMenu.value.site) {
    if (confirm(`确定要删除"${contextMenu.value.site.name}"吗？`)) {
      store.deleteSite(contextMenu.value.site.id)
    }
  }
  hideContextMenu()
}

// 隐藏右键菜单
const hideContextMenu = () => {
  contextMenu.value.show = false
  contextMenu.value.site = null
}

// 点击其他地方关闭菜单
const handleClickOutside = (event: Event) => {
  if (contextMenu.value.show) {
    const target = event.target as Element
    if (!contextMenuRef.value?.contains(target)) {
      hideContextMenu()
    }
  }
}

// 动态背景样式
const containerStyle = computed(() => {
  const { background } = store.settings
  let backgroundStyle = ''

  if (background.type === 'color') {
    backgroundStyle = background.value
  } else if (background.type === 'gradient') {
    backgroundStyle = background.value
  } else if (background.type === 'image') {
    backgroundStyle = `url(${background.value}) center/cover no-repeat`
  } else {
    backgroundStyle = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  }

  return {
    background: backgroundStyle,
  }
})

// 动态网格样式
const gridStyle = computed(() => {
  const { layout } = store.settings
  return {
    gridTemplateColumns: `repeat(${layout.columns}, 1fr)`,
    gap: `${layout.gap}px`,
  }
})

// 动态卡片样式
const cardStyle = computed(() => {
  const { cardStyle } = store.settings
  return {
    borderRadius: `${cardStyle.borderRadius}px`,
    opacity: cardStyle.opacity,
    boxShadow: cardStyle.shadow
      ? '0 8px 32px rgba(0, 0, 0, 0.1), inset 0 1px 0 rgba(255, 255, 255, 0.2)'
      : 'none',
  }
})

onMounted(() => {
  console.log('HomeView mounted')
  store.loadData()
  console.log('Data loaded, sites count:', store.sites?.length || 0)

  // 添加全局点击事件监听，用于关闭右键菜单
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  // 清理事件监听器
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  padding: 20px;
  color: white;
  transition: background 0.3s ease;
}

.debug-info {
  position: fixed;
  top: 10px;
  right: 10px;
  background: red;
  color: white;
  padding: 5px 10px;
  border-radius: 4px;
  z-index: 9999;
  font-size: 12px;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.app-header h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 600;
  background: linear-gradient(135deg, #ffffff 0%, #f0f0f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.btn-primary,
.btn-secondary {
  padding: 12px 24px;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
  position: relative;
  overflow: hidden;
}

.btn-primary {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.25) 0%, rgba(255, 255, 255, 0.15) 100%);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.btn-primary:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.35) 0%, rgba(255, 255, 255, 0.25) 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.btn-secondary {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.btn-secondary:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, rgba(255, 255, 255, 0.1) 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
}

.sites-grid {
  display: grid;
  /* 动态列数和间距通过内联样式设置 */
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 12px;
  transition: gap 0.3s ease;
}

.site-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.08) 0%, rgba(255, 255, 255, 0.03) 100%);
  /* 动态圆角、透明度、阴影通过内联样式设置 */
  border-radius: 12px;
  padding: 12px 8px;
  text-align: center;
  cursor: pointer;
  transition:
    all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275),
    border-radius 0.3s ease,
    opacity 0.3s ease,
    box-shadow 0.3s ease;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow:
    0 4px 16px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
  position: relative;
  overflow: hidden;
}

/* 添加微光效果 */
.site-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left 0.5s ease;
}

.site-card:hover::before {
  left: 100%;
}

.site-card:hover {
  transform: translateY(-4px) scale(1.02);
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.15) 0%, rgba(255, 255, 255, 0.08) 100%);
  box-shadow:
    0 8px 24px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.add-card {
  border: 2px dashed rgba(255, 255, 255, 0.4);
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.08) 0%,
    rgba(255, 255, 255, 0.03) 100%
  ) !important;
  position: relative;
}

.add-card::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 18px;
  background: linear-gradient(
    145deg,
    transparent 0%,
    rgba(255, 255, 255, 0.05) 50%,
    transparent 100%
  );
  pointer-events: none;
}

.add-card:hover {
  border-color: rgba(255, 255, 255, 0.6);
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.12) 0%,
    rgba(255, 255, 255, 0.06) 100%
  ) !important;
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
  transition: all 0.3s ease;
  position: relative;
  border: none;
}

/* 当没有图标时显示背景 */
.site-icon:not(:has(img)) {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.12) 0%, rgba(255, 255, 255, 0.06) 100%);
  backdrop-filter: blur(15px);
  border: none;
  border-radius: 12px;
}

.site-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  background: transparent;
  opacity: 0.85;
  transition: all 0.3s ease;
  filter: brightness(1.1) contrast(0.9);
  border: none;
}

.site-card:hover .site-icon {
  transform: scale(1.1);
}

.site-card:hover .site-icon img {
  opacity: 1;
  filter: brightness(1.2) contrast(1);
}

/* 额外的透明图标效果类（可选） */
.site-icon.extra-transparent img {
  opacity: 0.7;
  filter: brightness(1.3) contrast(0.8) saturate(0.9);
}

.site-card:hover .site-icon.extra-transparent img {
  opacity: 0.9;
  filter: brightness(1.4) contrast(1) saturate(1);
}

.site-name {
  font-size: 11px;
  font-weight: 600;
  word-break: break-word;
  line-height: 1.2;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-container {
    padding: 16px;
  }

  .app-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .app-header h1 {
    font-size: 1.8rem;
  }

  .header-actions {
    flex-wrap: wrap;
    justify-content: center;
  }

  .sites-grid {
    grid-template-columns: repeat(auto-fill, minmax(70px, 1fr));
    gap: 10px;
  }

  .site-card {
    padding: 10px 6px;
    border-radius: 10px;
  }

  .site-icon {
    width: 28px;
    height: 28px;
    margin-bottom: 6px;
  }

  .site-icon:not(:has(img)) {
    border-radius: 7px;
  }

  .site-name {
    font-size: 10px;
  }
}

/* 平板设备优化 */
@media (min-width: 769px) and (max-width: 1024px) {
  .sites-grid {
    grid-template-columns: repeat(auto-fill, minmax(90px, 1fr));
    gap: 14px;
  }
}

/* 大屏幕优化 */
@media (min-width: 1400px) {
  .sites-grid {
    grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
    gap: 18px;
  }

  .site-card {
    padding: 18px 14px;
  }

  .site-icon {
    width: 44px;
    height: 44px;
    margin-bottom: 14px;
  }

  .site-name {
    font-size: 13px;
  }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .search-input {
    background: rgba(255, 255, 255, 0.1);
    color: white;
  }

  .search-input::placeholder {
    color: rgba(255, 255, 255, 0.6);
  }

  .search-input:focus {
    background: rgba(255, 255, 255, 0.15);
  }
}

/* 动画优化 */
@media (prefers-reduced-motion: reduce) {
  .site-card,
  .btn-primary,
  .btn-secondary,
  .search-input,
  .site-icon {
    transition: none;
  }

  .site-card::before {
    display: none;
  }
}

/* 高对比度模式 */
@media (prefers-contrast: high) {
  .site-card {
    border-width: 2px;
    border-color: rgba(255, 255, 255, 0.8);
  }

  .btn-primary,
  .btn-secondary {
    border-width: 2px;
  }
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

/* 移动端右键菜单优化 */
@media (max-width: 768px) {
  .context-menu {
    min-width: 160px;
  }

  .context-menu-item {
    padding: 14px 18px;
    font-size: 15px;
  }
}
</style>
