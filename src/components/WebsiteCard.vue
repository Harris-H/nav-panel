<template>
  <div
    class="website-card"
    :style="cardStyle"
    @click="openWebsite"
    @contextmenu.prevent="showContextMenu"
  >
    <div class="card-icon">
      <img v-if="website.icon" :src="website.icon" :alt="website.name" @error="handleIconError" />
      <div v-else class="default-icon">
        {{ website.name.charAt(0).toUpperCase() }}
      </div>
    </div>
    <div class="card-content">
      <h3 class="card-title">{{ website.name }}</h3>
      <p v-if="website.description" class="card-description">
        {{ website.description }}
      </p>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="showMenu"
      class="context-menu"
      :style="{ top: menuY + 'px', left: menuX + 'px' }"
      @click.stop
    >
      <button @click="editCard">编辑</button>
      <button @click="deleteCard" class="delete-btn">删除</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAppStore } from '@/stores/app'
import type { WebsiteCard } from '@/types'

interface Props {
  website: WebsiteCard
}

const props = defineProps<Props>()
const store = useAppStore()

const showMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)

const cardStyle = computed(() => ({
  borderRadius: `${store.settings.cardStyle.borderRadius}px`,
  boxShadow: store.settings.cardStyle.shadow ? '0 4px 12px rgba(0, 0, 0, 0.15)' : 'none',
  opacity: store.settings.cardStyle.opacity,
}))

const openWebsite = () => {
  if (props.website.url) {
    window.open(props.website.url, '_blank')
  }
  hideContextMenu()
}

const showContextMenu = (event: MouseEvent) => {
  event.preventDefault()
  menuX.value = event.clientX
  menuY.value = event.clientY
  showMenu.value = true

  // 点击其他地方隐藏菜单
  document.addEventListener('click', hideContextMenu, { once: true })
}

const hideContextMenu = () => {
  showMenu.value = false
}

const editCard = () => {
  store.openEditCardModal(props.website)
  hideContextMenu()
}

const deleteCard = () => {
  if (confirm(`确定要删除 "${props.website.name}" 吗？`)) {
    store.deleteSite(props.website.id)
  }
  hideContextMenu()
}

const handleIconError = (event: Event) => {
  const target = event.target as HTMLImageElement
  target.style.display = 'none'
}
</script>

<style scoped>
.website-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.9) 100%);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 20px;
  padding: 24px 20px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  min-height: 140px;
  overflow: hidden;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
}

/* 微光效果 */
.website-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
  z-index: 0;
}

.website-card:hover::before {
  left: 100%;
}

.website-card:hover {
  transform: translateY(-8px) scale(1.02);
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.98) 0%, rgba(255, 255, 255, 0.93) 100%);
  box-shadow:
    0 16px 48px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
  border-color: rgba(255, 255, 255, 0.5);
}

.card-icon {
  width: 56px;
  height: 56px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
  border: none;
}

.website-card:hover .card-icon {
  transform: scale(1.1);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.card-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 14px;
  transition: all 0.3s ease;
  opacity: 0.85;
  filter: brightness(1.1) contrast(0.9);
  border: none;
}

.website-card:hover .card-icon img {
  opacity: 1;
  filter: brightness(1.2) contrast(1);
}

.default-icon {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(103, 126, 234, 0.75) 0%, rgba(118, 75, 162, 0.75) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  border-radius: 14px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(15px);
  border: none;
}

/* 额外的透明图标效果类（可选） */
.card-icon.extra-transparent img {
  opacity: 0.7;
  filter: brightness(1.3) contrast(0.8) saturate(0.9);
}

.website-card:hover .card-icon.extra-transparent img {
  opacity: 0.9;
  filter: brightness(1.4) contrast(1) saturate(1);
}

.card-content {
  flex: 1;
  position: relative;
  z-index: 1;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  margin: 0 0 6px 0;
  color: #2d3748;
  line-height: 1.3;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.card-description {
  font-size: 12px;
  color: #718096;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-shadow: 0 1px 1px rgba(255, 255, 255, 0.6);
}

.context-menu {
  position: fixed;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.9) 100%);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  box-shadow:
    0 12px 40px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  z-index: 1000;
  min-width: 120px;
  overflow: hidden;
  animation: contextMenuFadeIn 0.2s ease-out;
}

@keyframes contextMenuFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-5px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.context-menu button {
  display: block;
  width: 100%;
  padding: 12px 16px;
  border: none;
  background: none;
  text-align: left;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #2d3748;
  transition: all 0.2s ease;
  position: relative;
}

.context-menu button::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, rgba(103, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.context-menu button:hover::before {
  opacity: 1;
}

.context-menu button:hover {
  color: #1a202c;
  background: linear-gradient(90deg, rgba(103, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.context-menu button:first-child {
  border-radius: 12px 12px 0 0;
}

.context-menu button:last-child {
  border-radius: 0 0 12px 12px;
}

.delete-btn {
  color: #e53e3e !important;
}

.delete-btn:hover {
  background: linear-gradient(
    90deg,
    rgba(229, 62, 62, 0.1) 0%,
    rgba(229, 62, 62, 0.05) 100%
  ) !important;
  color: #c53030 !important;
}

.delete-btn:hover::before {
  background: linear-gradient(90deg, rgba(229, 62, 62, 0.1) 0%, rgba(229, 62, 62, 0.05) 100%);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .website-card {
    padding: 20px 16px;
    min-height: 120px;
    border-radius: 16px;
  }

  .card-icon {
    width: 48px;
    height: 48px;
    margin-bottom: 12px;
    border-radius: 12px;
  }

  .default-icon {
    font-size: 16px;
    border-radius: 10px;
  }

  .card-title {
    font-size: 13px;
  }

  .card-description {
    font-size: 11px;
  }

  .context-menu {
    min-width: 100px;
  }

  .context-menu button {
    padding: 10px 14px;
    font-size: 13px;
  }
}

/* 平板设备优化 */
@media (min-width: 769px) and (max-width: 1024px) {
  .website-card {
    padding: 22px 18px;
    min-height: 130px;
  }

  .card-icon {
    width: 52px;
    height: 52px;
    margin-bottom: 14px;
  }
}

/* 大屏幕优化 */
@media (min-width: 1400px) {
  .website-card {
    padding: 28px 24px;
    min-height: 160px;
    border-radius: 24px;
  }

  .card-icon {
    width: 64px;
    height: 64px;
    margin-bottom: 20px;
    border-radius: 18px;
  }

  .default-icon {
    font-size: 24px;
    border-radius: 16px;
  }

  .card-title {
    font-size: 15px;
  }

  .card-description {
    font-size: 13px;
  }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .website-card {
    background: linear-gradient(
      145deg,
      rgba(255, 255, 255, 0.15) 0%,
      rgba(255, 255, 255, 0.08) 100%
    );
    border-color: rgba(255, 255, 255, 0.2);
  }

  .website-card:hover {
    background: linear-gradient(
      145deg,
      rgba(255, 255, 255, 0.25) 0%,
      rgba(255, 255, 255, 0.15) 100%
    );
  }

  .card-title {
    color: #f7fafc;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  }

  .card-description {
    color: #e2e8f0;
    text-shadow: 0 1px 1px rgba(0, 0, 0, 0.2);
  }

  .context-menu {
    background: linear-gradient(145deg, rgba(45, 55, 72, 0.95) 0%, rgba(45, 55, 72, 0.9) 100%);
    border-color: rgba(255, 255, 255, 0.2);
  }

  .context-menu button {
    color: #f7fafc;
  }

  .context-menu button:hover {
    color: #ffffff;
  }
}

/* 动画优化 */
@media (prefers-reduced-motion: reduce) {
  .website-card,
  .card-icon,
  .context-menu button {
    transition: none;
  }

  .website-card::before {
    display: none;
  }

  .context-menu {
    animation: none;
  }
}

/* 高对比度模式 */
@media (prefers-contrast: high) {
  .website-card {
    border-width: 2px;
    border-color: rgba(0, 0, 0, 0.5);
  }

  .card-title {
    color: #000000;
    text-shadow: none;
  }

  .card-description {
    color: #333333;
    text-shadow: none;
  }

  .context-menu {
    border-width: 2px;
    border-color: rgba(0, 0, 0, 0.5);
  }
}
</style>
