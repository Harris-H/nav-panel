import type { Website, AppSettings, SearchEngine } from '@/types'

const API_BASE_URL = 'http://localhost:8080/api'

class ApiClient {
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    try {
      const response = await fetch(url, config)

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return data.data || data
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error)
      throw error
    }
  }

  // 多部分表单请求（用于文件上传）
  private async requestWithFile<T>(endpoint: string, formData: FormData): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    try {
      const response = await fetch(url, {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return data.data || data
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error)
      throw error
    }
  }

  // 网站相关API
  async getWebsites(): Promise<Website[]> {
    return this.request<Website[]>('/websites')
  }

  async createWebsite(website: Omit<Website, 'id' | 'createdAt' | 'updatedAt'>): Promise<Website> {
    return this.request<Website>('/websites', {
      method: 'POST',
      body: JSON.stringify(website),
    })
  }

  async updateWebsite(id: string, updates: Partial<Website>): Promise<Website> {
    return this.request<Website>(`/websites/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    })
  }

  async deleteWebsite(id: string): Promise<void> {
    return this.request<void>(`/websites/${id}`, {
      method: 'DELETE',
    })
  }

  // 搜索引擎相关API
  async getSearchEngines(): Promise<SearchEngine[]> {
    return this.request<SearchEngine[]>('/search-engines')
  }

  async createSearchEngine(engine: Omit<SearchEngine, 'id'>): Promise<SearchEngine> {
    return this.request<SearchEngine>('/search-engines', {
      method: 'POST',
      body: JSON.stringify(engine),
    })
  }

  // 创建带图片上传的搜索引擎
  async createSearchEngineWithIcon(
    engine: { name: string; url: string; placeholder?: string; isDefault?: boolean },
    iconFile?: File,
  ): Promise<SearchEngine> {
    const formData = new FormData()

    // 生成ID
    const id = engine.name.toLowerCase().replace(/[^a-z0-9]/g, '') + '_' + Date.now()

    formData.append('id', id)
    formData.append('name', engine.name)
    formData.append('url', engine.url)

    if (engine.placeholder) {
      formData.append('placeholder', engine.placeholder)
    }

    if (engine.isDefault) {
      formData.append('isDefault', 'true')
    }

    if (iconFile) {
      formData.append('icon', iconFile)
    }

    return this.requestWithFile<SearchEngine>('/search-engines/with-icon', formData)
  }

  async updateSearchEngine(id: string, updates: Partial<SearchEngine>): Promise<SearchEngine> {
    return this.request<SearchEngine>(`/search-engines/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    })
  }

  // 更新带图片上传的搜索引擎
  async updateSearchEngineWithIcon(
    id: string,
    updates: { name?: string; url?: string; placeholder?: string; isDefault?: boolean },
    iconFile?: File,
  ): Promise<SearchEngine> {
    const formData = new FormData()

    if (updates.name) formData.append('name', updates.name)
    if (updates.url) formData.append('url', updates.url)
    if (updates.placeholder) formData.append('placeholder', updates.placeholder)
    if (updates.isDefault !== undefined) formData.append('isDefault', updates.isDefault.toString())

    if (iconFile) {
      formData.append('icon', iconFile)
    }

    const url = `${API_BASE_URL}/search-engines/${id}/with-icon`

    try {
      const response = await fetch(url, {
        method: 'PUT',
        body: formData,
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return data.data || data
    } catch (error) {
      console.error(`API request failed: /search-engines/${id}/with-icon`, error)
      throw error
    }
  }

  async deleteSearchEngine(id: string): Promise<void> {
    return this.request<void>(`/search-engines/${id}`, {
      method: 'DELETE',
    })
  }

  // 设置相关API
  async getSettings(): Promise<AppSettings> {
    return this.request<AppSettings>('/settings')
  }

  async updateSettings(settings: Partial<AppSettings>): Promise<AppSettings> {
    return this.request<AppSettings>('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    })
  }

  // 导入导出API
  async exportData(): Promise<{
    websites: Website[]
    searchEngines: SearchEngine[]
    settings: AppSettings
  }> {
    return this.request<{
      websites: Website[]
      searchEngines: SearchEngine[]
      settings: AppSettings
    }>('/export')
  }

  async importData(data: {
    websites?: Website[]
    searchEngines?: SearchEngine[]
    settings?: AppSettings
  }): Promise<void> {
    return this.request<void>('/import', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }
}

export const api = new ApiClient()
