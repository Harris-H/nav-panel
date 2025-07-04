export interface Website {
  id: string
  name: string
  url: string
  icon?: string
  description?: string
  category?: string
  groupId?: string | null
  sortOrder?: number
  createdAt?: string
  updatedAt?: string
}

// 为了兼容性，导出 WebsiteCard 作为 Website 的别名
export type WebsiteCard = Website

export interface SearchEngine {
  id: string
  name: string
  url: string
  icon?: string // 图标DataURL
  placeholder?: string
  isDefault?: boolean
}

export interface AppSettings {
  theme: 'light' | 'dark'
  layout: {
    columns: number
    cardSize: 'small' | 'medium' | 'large'
    showLabels: boolean
    gap: number
  }
  background: {
    type: 'color' | 'gradient' | 'image'
    value: string
  }
  cardStyle: {
    borderRadius: number
    opacity: number
    shadow: boolean
  }
  search: {
    enabled: boolean
    engines: SearchEngine[]
    defaultEngineId: string
    openInNewTab: boolean
  }
}

export interface Category {
  id: string
  name: string
  color?: string
  icon?: string
}

export interface Group {
  id: string
  name: string
  color?: string
  icon?: string
  sortOrder: number
  isCollapsed: boolean
  createdAt?: string
  updatedAt?: string
}

export interface GroupWithWebsites extends Group {
  websites: Website[]
}
