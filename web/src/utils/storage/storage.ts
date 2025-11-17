/**
 * 本地存储工具类
 * 提供对localStorage的封装，支持对象存储和过期时间设置

 */
interface StorageData<T> {
  data: T
  expires?: number // 过期时间戳
}

export class StorageUtil {
  private static isBrowser(): boolean {
    return typeof window !== 'undefined' && typeof localStorage !== 'undefined'
  }

  static set<T>(key: string, value: T, expires?: number): void {
    if (!this.isBrowser()) {
      return
    }

    const storageData: StorageData<T> = {
      data: value,
    }

    if (expires !== undefined) {
      const expiresTime = new Date().getTime() + expires * 60 * 60 * 1000
      storageData.expires = expiresTime
    }

    localStorage.setItem(key, JSON.stringify(storageData))
  }

  static get<T>(key: string): T | null {
    if (!this.isBrowser()) {
      return null
    }

    const value = localStorage.getItem(key)
    if (!value) {
      return null
    }

    try {
      const storageData = JSON.parse(value) as StorageData<T>

      if (storageData.expires && storageData.expires < new Date().getTime()) {
        this.remove(key)
        return null
      }

      return storageData.data
    } catch (error) {
      console.error('解析本地存储数据失败:', error)
      this.remove(key)
      return null
    }
  }

  static remove(key: string): void {
    if (!this.isBrowser()) {
      return
    }
    localStorage.removeItem(key)
  }

  static clear(): void {
    if (!this.isBrowser()) {
      return
    }
    localStorage.clear()
  }

  static has(key: string): boolean {
    return this.get(key) !== null
  }

  static keys(): string[] {
    if (!this.isBrowser()) {
      return []
    }

    const keys: string[] = []
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key) {
        keys.push(key)
      }
    }
    return keys
  }
}
