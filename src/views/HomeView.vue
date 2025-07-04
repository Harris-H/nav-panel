<template>
  <div class="home-container" :style="containerStyle">
    <header class="app-header">
      <h1>导航面板</h1>
      <div class="header-actions">
        <button @click="store.openGroupModal()" class="btn-secondary">管理分组</button>
        <button @click="store.openAddCardModal()" class="btn-primary">添加网站</button>
        <button @click="store.openSettingsModal()" class="btn-secondary">设置</button>
      </div>
    </header>

    <!-- 搜索框组件 -->
    <SearchBox />

    <main class="main-content">
      <!-- 分组展示 -->
      <template v-for="group in store.groupsWithWebsites" :key="group.id">
        <GroupSection
          :group="group"
          :websites="group.websites"
          :grid-style="gridStyle"
          :card-style="cardStyle"
          @toggle-collapse="store.toggleGroupCollapse"
          @edit-group="store.openEditGroupModal"
          @delete-group="handleDeleteGroup"
          @site-right-click="handleRightClick"
          @site-mouse-down="handleSiteMouseDown"
          @drop-website="handleDropWebsite"
        />
      </template>

      <!-- 未分组的网站 -->
      <div v-if="ungroupedSites.length > 0" class="ungrouped-section">
        <div class="section-header">
          <h3>未分组</h3>
        </div>
        <div class="sites-grid" :style="gridStyle">
          <template v-for="(site, index) in ungroupedSites" :key="site.id">
            <!-- 插入指示器 - 根据移动方向显示 -->
            <div
              v-if="
                dragState.isDragging &&
                (dragState.movingRight
                  ? dragState.insertIndex + 1 === index
                  : dragState.insertIndex === index)
              "
              class="insert-indicator"
            ></div>

            <div
              class="site-card"
              :style="cardStyle"
              :class="{
                dragging: dragState.isDragging && dragState.dragIndex === index,
                'drag-placeholder': dragState.isDragging && dragState.dragIndex === index,
              }"
              @contextmenu.prevent="handleRightClick($event, site)"
              @mousedown="handleUngroupedSiteMouseDown($event, index)"
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

          <!-- 末尾插入指示器 - 根据移动方向显示 -->
          <div
            v-if="
              dragState.isDragging &&
              (dragState.movingRight
                ? dragState.insertIndex === ungroupedSites.length - 1
                : dragState.insertIndex === ungroupedSites.length)
            "
            class="insert-indicator"
          ></div>

          <div class="site-card add-card" :style="cardStyle" @click="store.openAddCardModal()">
            <div class="site-icon">+</div>
            <div class="site-name">添加网站</div>
          </div>
        </div>
      </div>

      <!-- 如果没有分组也没有网站 -->
      <div
        v-if="store.groupsWithWebsites.length === 0 && ungroupedSites.length === 0"
        class="empty-state"
      >
        <div class="empty-icon">🌟</div>
        <div class="empty-text">还没有添加任何网站</div>
        <button @click="store.openAddCardModal()" class="btn-primary">添加第一个网站</button>
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

    <!-- 分组管理模态框 -->
    <GroupModal
      :show="store.isGroupModalOpen"
      :group="store.editingGroup"
      @close="store.closeGroupModal"
      @submit="handleGroupSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, computed, ref, nextTick } from 'vue'
import { useAppStore } from '@/stores/app'
import AddCardModal from '@/components/AddCardModal.vue'
import SettingsModal from '@/components/SettingsModal.vue'
import SearchBox from '@/components/SearchBox.vue'
import GroupModal from '@/components/GroupModal.vue'
import GroupSection from '@/components/GroupSection.vue'
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

// 拖拽状态
const dragState = ref({
  isDragging: false,
  dragIndex: -1,
  insertIndex: -1, // 改为插入位置而不是覆盖位置
  startX: 0,
  startY: 0,
  currentX: 0,
  currentY: 0,
  lastX: 0, // 上一次的X位置，用于计算移动方向
  movingRight: true, // 是否向右移动
  dragPreview: null as HTMLElement | null,
})

// 右键菜单样式
const contextMenuStyle = computed(() => ({
  left: `${contextMenu.value.x}px`,
  top: `${contextMenu.value.y}px`,
}))

const openSite = (url: string) => {
  window.open(url, '_blank')
}

// 处理分组表单提交
const handleGroupSubmit = async (data: { name: string; color?: string; icon?: string }) => {
  try {
    if (store.editingGroup) {
      // 编辑分组
      await store.updateGroup(store.editingGroup.id, data)
    } else {
      // 创建分组
      await store.createGroup(data)
    }
    // 操作成功后关闭模态框
    store.closeGroupModal()
  } catch (error) {
    console.error('Group operation failed:', error)
  }
}

// 处理删除分组
const handleDeleteGroup = async (groupId: string) => {
  if (confirm('确定要删除这个分组吗？分组内的网站将移动到未分组。')) {
    try {
      await store.deleteGroup(groupId)
    } catch (error) {
      console.error('Delete group failed:', error)
    }
  }
}

// 处理分组内网站的鼠标按下
const handleSiteMouseDown = (event: MouseEvent, index: number, groupId: string) => {
  // 这里可以实现分组内网站的拖拽功能
  console.log('Site mouse down in group:', { index, groupId })
}

// 处理未分组网站的鼠标按下
const handleUngroupedSiteMouseDown = (event: MouseEvent, index: number) => {
  // 使用原有的拖拽逻辑
  handleMouseDown(event, index)
}

// 处理网站拖放到分组
const handleDropWebsite = async (websiteId: string, targetGroupId: string, position?: number) => {
  try {
    await store.moveWebsiteToGroup(websiteId, targetGroupId, position)
  } catch (error) {
    console.error('Drop website failed:', error)
  }
}

// 拖拽辅助函数
const createDragPreview = (sourceElement: HTMLElement, index: number) => {
  const site = store.sites[index]
  if (!site) {
    return
  }

  // 找到实际的卡片元素
  const cardElement = sourceElement.closest('.site-card') as HTMLElement
  if (!cardElement) {
    return
  }

  // 创建拖拽预览元素，克隆原卡片的结构和样式
  const preview = document.createElement('div')
  preview.className = 'drag-preview site-card'

  // 复制原卡片的内容，保持完整的卡片结构
  const iconElement = site.icon
    ? `<img src="${site.icon}" alt="${site.name}" />`
    : `<span>${site.name.charAt(0).toUpperCase()}</span>`

  preview.innerHTML = `
    <div class="site-icon">
      ${iconElement}
    </div>
    <div class="site-name">${site.name}</div>
  `

  // 设置初始位置（鼠标在卡片中心）
  const currentX = dragState.value.currentX || dragState.value.startX
  const currentY = dragState.value.currentY || dragState.value.startY

  // 卡片尺寸
  const cardWidth = 80
  const cardHeight = 80 // 准确高度：14px(上padding) + 32px(图标) + 8px(间距) + 15px(文字) + 14px(下padding) ≈ 83px

  preview.style.position = 'fixed'
  preview.style.left = `${currentX - cardWidth / 2}px` // 鼠标在卡片水平中心
  preview.style.top = `${currentY - cardHeight / 2}px` // 鼠标在卡片垂直中心
  preview.style.zIndex = '9999'
  preview.style.pointerEvents = 'none'
  preview.style.width = `${cardWidth}px` // 固定宽度，保持卡片比例

  document.body.appendChild(preview)
  dragState.value.dragPreview = preview
}

const updateDragPreview = (x: number, y: number) => {
  if (dragState.value.dragPreview) {
    // 卡片尺寸（与创建时保持一致）
    const cardWidth = 80
    const cardHeight = 80

    // 让鼠标位置在卡片中心
    dragState.value.dragPreview.style.left = `${x - cardWidth / 2}px`
    dragState.value.dragPreview.style.top = `${y - cardHeight / 2}px`
  }
}

const calculateInsertPosition = (x: number): number => {
  const sitesGrid = document.querySelector('.sites-grid')
  if (!sitesGrid) return -1

  const allCards = Array.from(document.querySelectorAll('.site-card:not(.add-card)'))
  if (allCards.length === 0) return 0

  // 如果只有一个卡片（被拖拽的），返回0
  if (allCards.length === 1) return 0

  // 获取所有卡片的位置信息，排除被拖拽的卡片
  const visibleCardInfos = allCards
    .map((card, index) => ({
      element: card,
      originalIndex: index, // 在原始sites数组中的索引
      rect: card.getBoundingClientRect(),
      isDragTarget: index === dragState.value.dragIndex,
    }))
    .filter((card) => !card.isDragTarget) // 只保留可见的卡片

  // 计算插入位置的逻辑：
  // 1. 遍历所有可见卡片（排除被拖拽的）
  // 2. 找到鼠标当前悬停的卡片区域
  // 3. 根据移动方向决定插入位置
  // 4. 最后转换回原始数组的索引

  let insertIndex = store.sites.length // 默认插入到末尾

  for (let i = 0; i < visibleCardInfos.length; i++) {
    const card = visibleCardInfos[i]
    const rect = card.rect

    const leftEdge = rect.left
    const rightEdge = rect.right

    // 检查鼠标是否在当前卡片的水平范围内
    if (x >= leftEdge && x <= rightEdge) {
      // 根据移动方向使用不同的策略
      if (dragState.value.movingRight) {
        // 向右移动：基于前一个卡片的位置计算threshold
        let threshold = leftEdge
        if (i > 0) {
          // 如果有前一个卡片，基于前一个卡片的右侧向左偏移30%宽度
          const prevCard = visibleCardInfos[i - 1]
          threshold = prevCard.rect.right - prevCard.rect.width * 0.05
        } else {
          // 如果是第一个卡片，使用卡片自身的10%位置
          threshold = leftEdge + rect.width * 0.1
        }

        if (x >= threshold) {
          // 插入到这个可见卡片之后
          // 需要转换回原始索引：找到下一个卡片的位置，或插入到末尾
          if (i + 1 < visibleCardInfos.length) {
            insertIndex = visibleCardInfos[i + 1].originalIndex
          } else {
            insertIndex = store.sites.length
          }
        } else {
          // 插入到这个可见卡片之前
          insertIndex = card.originalIndex
        }
      } else {
        // 向左移动：只要鼠标还没到达卡片右边缘的30%，就插入到左侧
        const threshold = rightEdge - rect.width * 0.05
        if (x <= threshold) {
          // 插入到这个可见卡片之前
          insertIndex = card.originalIndex
        } else {
          // 插入到这个可见卡片之后
          if (i + 1 < visibleCardInfos.length) {
            insertIndex = visibleCardInfos[i + 1].originalIndex
          } else {
            insertIndex = store.sites.length
          }
        }
      }
      break
    }

    // 如果鼠标在当前卡片左边，插入到这个位置
    if (x < leftEdge) {
      insertIndex = card.originalIndex
      break
    }
  }

  // 处理被拖拽卡片的索引调整
  // 如果插入位置在被拖拽卡片之后，需要减1（因为被拖拽的卡片会被移除）
  if (insertIndex > dragState.value.dragIndex) {
    insertIndex = insertIndex - 1
  }

  // 确保索引在有效范围内，允许插入到末尾
  insertIndex = Math.max(0, Math.min(insertIndex, store.sites.length))

  return insertIndex
}

const cleanupDragPreview = () => {
  if (dragState.value.dragPreview) {
    document.body.removeChild(dragState.value.dragPreview)
    dragState.value.dragPreview = null
  }
}

// 拖拽功能
const handleMouseDown = (event: MouseEvent, index: number) => {
  // 阻止右键和中键拖拽
  if (event.button !== 0) return

  // 防止拖拽时触发右键菜单
  hideContextMenu()

  let hasDragged = false
  const startX = event.clientX
  const startY = event.clientY

  const handleMouseMove = (moveEvent: MouseEvent) => {
    const deltaX = moveEvent.clientX - startX
    const deltaY = moveEvent.clientY - startY
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY)

    // 当移动距离超过5px时开始拖拽
    if (!hasDragged && distance > 5) {
      hasDragged = true

      // 先设置拖拽状态
      dragState.value = {
        isDragging: true,
        dragIndex: index,
        insertIndex: -1,
        startX: startX,
        startY: startY,
        currentX: moveEvent.clientX,
        currentY: moveEvent.clientY,
        lastX: moveEvent.clientX,
        movingRight: true,
        dragPreview: null,
      }

      // 然后创建拖拽预览
      createDragPreview(event.target as HTMLElement, index)

      document.body.style.userSelect = 'none'
      document.body.classList.add('dragging')
    }

    if (hasDragged) {
      // 计算移动方向
      const deltaX = moveEvent.clientX - dragState.value.lastX
      if (Math.abs(deltaX) > 2) {
        // 避免微小移动导致频繁切换方向
        dragState.value.movingRight = deltaX > 0
      }

      dragState.value.lastX = dragState.value.currentX
      dragState.value.currentX = moveEvent.clientX
      dragState.value.currentY = moveEvent.clientY

      // 更新拖拽预览位置
      updateDragPreview(moveEvent.clientX, moveEvent.clientY)

      // 计算插入位置
      const insertIndex = calculateInsertPosition(moveEvent.clientX)
      dragState.value.insertIndex = insertIndex
    }
  }

  const handleMouseUp = async () => {
    if (hasDragged) {
      // 执行排序 - 允许插入到末尾
      if (dragState.value.insertIndex !== -1 && dragState.value.insertIndex !== index) {
        await store.reorderSites(index, dragState.value.insertIndex)
      }

      // 清理拖拽预览
      cleanupDragPreview()

      document.body.style.userSelect = ''
      document.body.classList.remove('dragging')
    } else {
      // 如果没有拖拽，执行原来的点击事件
      const site = store.sites[index]
      if (site) {
        openSite(site.url)
      }
    }

    // 重置拖拽状态
    dragState.value = {
      isDragging: false,
      dragIndex: -1,
      insertIndex: -1,
      startX: 0,
      startY: 0,
      currentX: 0,
      currentY: 0,
      lastX: 0,
      movingRight: true,
      dragPreview: null,
    }

    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
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
    // 检查点击的元素不在HomeView的右键菜单内，且不在任何右键菜单内
    if (!contextMenuRef.value?.contains(target) && !target.closest('.context-menu')) {
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

// 未分组的网站
const ungroupedSites = computed(() => {
  // 更准确的过滤条件：groupId 为 null、undefined 或空字符串的网站
  return store.sites.filter(
    (site) =>
      !site.groupId || site.groupId === '' || site.groupId === null || site.groupId === undefined,
  )
})

onMounted(() => {
  store.loadData()

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
  cursor: grab;
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

/* 拖拽相关样式 */
.site-card:active {
  cursor: grabbing;
}

.site-card.dragging {
  opacity: 0.3;
  transform: scale(0.95);
  pointer-events: none;
  transition: all 0.2s ease;
}

.site-card.drag-placeholder {
  opacity: 0.3;
  transform: scale(0.95);
  border: 2px dashed rgba(255, 255, 255, 0.5);
}

/* 拖拽预览样式 - 继承卡片样式并增强 */
.drag-preview.site-card {
  position: fixed !important;
  z-index: 9999 !important;
  pointer-events: none !important;
  transform: rotate(3deg) scale(1.08) !important;
  opacity: 0.92 !important;
  transition: none !important;

  /* 增强的卡片样式 */
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.3) 0%,
    rgba(255, 255, 255, 0.18) 100%
  ) !important;
  backdrop-filter: blur(25px) !important;
  border: 2px solid rgba(255, 255, 255, 0.5) !important;
  border-radius: 15px !important;
  box-shadow:
    0 25px 50px rgba(0, 0, 0, 0.6),
    inset 0 2px 6px rgba(255, 255, 255, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.2) !important;

  /* 确保内容样式正确 */
  color: white !important;
  cursor: grabbing !important;
  padding: 14px 10px !important;
  min-width: 80px !important;
}

/* 拖拽预览内容的特殊样式 */
.drag-preview.site-card .site-icon {
  transform: scale(1.1) !important;
  filter: brightness(1.1) contrast(1.05) !important;
}

.drag-preview.site-card .site-name {
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3) !important;
  font-weight: 700 !important;
}

/* 拖拽预览内的图标和名称会继承原卡片的样式 */

/* 插入指示器样式 */
.insert-indicator {
  width: 3px;
  height: 60px;
  background: linear-gradient(to bottom, #667eea, #764ba2);
  border-radius: 2px;
  opacity: 0.9;
  animation: pulse 0.8s ease-in-out infinite alternate;
  box-shadow: 0 0 8px rgba(102, 126, 234, 0.6);
  margin: 0 6px;
  align-self: center;
  justify-self: center;
}

@keyframes pulse {
  from {
    opacity: 0.6;
    transform: scaleX(1);
  }
  to {
    opacity: 1;
    transform: scaleX(1.2);
  }
}

/* 全局拖拽状态 */
body.dragging {
  cursor: grabbing !important;
}

body.dragging * {
  cursor: grabbing !important;
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

/* 未分组部分样式 */
.ungrouped-section {
  margin-bottom: 24px;
}

.section-header {
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  padding-left: 8px;
  border-left: 3px solid rgba(255, 255, 255, 0.5);
}

/* 空状态样式 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 24px;
  opacity: 0.8;
}

.empty-text {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 32px;
  max-width: 400px;
  line-height: 1.5;
}

.empty-state .btn-primary {
  font-size: 16px;
  padding: 16px 32px;
}
</style>
