/* 上传状态本地存储管理 */
export interface UploadSessionData {
  id: string
  sessionId: string
  fileName: string
  fileSize: number
  fileMD5: string
  mimeType: string
  chunkSize: number
  totalChunks: number
  uploadedChunks: number[] // 已上传的分片序号
  folderId?: string
  accessLevel: 'public' | 'private' | 'protected'
  optimize: boolean
  createdAt: number
  lastActivity: number
}

const STORAGE_KEY = 'chunked_upload_sessions'
const SESSION_EXPIRE_TIME = 24 * 60 * 60 * 1000 // 24小时

export class UploadStorageManager {
  static saveSession(sessionData: UploadSessionData): void {
    try {
      const sessions = this.getAllSessions()
      sessions[sessionData.id] = {
        ...sessionData,
        lastActivity: Date.now(),
      }
      localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions))
    } catch (error) {
      console.error('保存上传会话失败:', error)
    }
  }

  static getSession(id: string): UploadSessionData | null {
    try {
      const sessions = this.getAllSessions()
      const session = sessions[id]

      if (!session) {
        return null
      }

      if (Date.now() - session.createdAt > SESSION_EXPIRE_TIME) {
        this.removeSession(id)
        return null
      }

      return session
    } catch (error) {
      console.error('获取上传会话失败:', error)
      return null
    }
  }

  static getAllSessions(): Record<string, UploadSessionData> {
    try {
      const data = localStorage.getItem(STORAGE_KEY)
      return data ? JSON.parse(data) : {}
    } catch (error) {
      console.error('获取所有会话失败:', error)
      return {}
    }
  }

  static updateSessionProgress(id: string, uploadedChunks: number[]): void {
    try {
      const sessions = this.getAllSessions()
      if (sessions[id]) {
        sessions[id].uploadedChunks = uploadedChunks
        sessions[id].lastActivity = Date.now()
        localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions))
      }
    } catch (error) {
      console.error('更新会话进度失败:', error)
    }
  }

  static addUploadedChunk(id: string, chunkNumber: number): void {
    try {
      const session = this.getSession(id)
      if (session) {
        if (!session.uploadedChunks.includes(chunkNumber)) {
          session.uploadedChunks.push(chunkNumber)
          session.uploadedChunks.sort((a, b) => a - b)
          this.saveSession(session)
        }
      }
    } catch (error) {
      console.error('添加已上传分片失败:', error)
    }
  }

  static removeSession(id: string): void {
    try {
      const sessions = this.getAllSessions()
      delete sessions[id]
      localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions))
    } catch (error) {
      console.error('移除会话失败:', error)
    }
  }

  static cleanExpiredSessions(): void {
    try {
      const sessions = this.getAllSessions()
      const now = Date.now()
      let hasChanges = false

      Object.keys(sessions).forEach((id) => {
        if (now - sessions[id].createdAt > SESSION_EXPIRE_TIME) {
          delete sessions[id]
          hasChanges = true
        }
      })

      if (hasChanges) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions))
      }
    } catch (error) {
      console.error('清理过期会话失败:', error)
    }
  }

  static getResumableSessions(): UploadSessionData[] {
    try {
      this.cleanExpiredSessions()
      const sessions = this.getAllSessions()

      return Object.values(sessions)
        .filter((session) => session.uploadedChunks.length < session.totalChunks)
        .sort((a, b) => b.lastActivity - a.lastActivity) // 按最后活动时间排序
    } catch (error) {
      console.error('获取可恢复会话失败:', error)
      return []
    }
  }

  static findExistingSession(fileMD5: string, fileSize: number): UploadSessionData | null {
    try {
      const sessions = this.getAllSessions()

      for (const session of Object.values(sessions)) {
        if (session.fileMD5 === fileMD5 && session.fileSize === fileSize) {
          if (Date.now() - session.createdAt <= SESSION_EXPIRE_TIME) {
            return session
          }
        }
      }

      return null
    } catch (error) {
      console.error('查找现有会话失败:', error)
      return null
    }
  }

  static clearAllSessions(): void {
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch (error) {
      console.error('清空所有会话失败:', error)
    }
  }

  static getStorageInfo(): {
    sessionCount: number
    storageSize: number
    oldestSession?: UploadSessionData
    newestSession?: UploadSessionData
  } {
    try {
      const sessions = this.getAllSessions()
      const sessionList = Object.values(sessions)

      const storageData = localStorage.getItem(STORAGE_KEY) || '{}'
      const storageSize = new Blob([storageData]).size

      let oldestSession: UploadSessionData | undefined
      let newestSession: UploadSessionData | undefined

      if (sessionList.length > 0) {
        oldestSession = sessionList.reduce((oldest, current) => (current.createdAt < oldest.createdAt ? current : oldest))
        newestSession = sessionList.reduce((newest, current) => (current.lastActivity > newest.lastActivity ? current : newest))
      }

      return {
        sessionCount: sessionList.length,
        storageSize,
        oldestSession,
        newestSession,
      }
    } catch (error) {
      console.error('获取存储信息失败:', error)
      return { sessionCount: 0, storageSize: 0 }
    }
  }

  static cleanExpiredSessionsWithCount(): number {
    const sessions = this.getAllSessions()
    const expiredTime = Date.now() - 24 * 60 * 60 * 1000 // 24小时前
    let cleanedCount = 0

    const validSessions: Record<string, UploadSessionData> = {}

    Object.entries(sessions).forEach(([key, session]) => {
      if (session.lastActivity < expiredTime) {
        cleanedCount++
      } else {
        validSessions[key] = session
      }
    })

    if (cleanedCount > 0) {
      this.saveSessions(validSessions)
    }

    return cleanedCount
  }

  static getSessionStats(): {
    totalSessions: number
    activeSessions: number
    expiredSessions: number
    totalSize: number
  } {
    const sessions = this.getAllSessions()
    const expiredTime = Date.now() - 24 * 60 * 60 * 1000 // 24小时前

    let activeSessions = 0
    let expiredSessions = 0
    let totalSize = 0

    Object.values(sessions).forEach((session) => {
      totalSize += session.fileSize
      if (session.lastActivity < expiredTime) {
        expiredSessions++
      } else {
        activeSessions++
      }
    })

    return {
      totalSessions: Object.keys(sessions).length,
      activeSessions,
      expiredSessions,
      totalSize,
    }
  }

  static startPeriodicCleanup(): () => void {
    this.cleanExpiredSessions()

    const intervalId = setInterval(
      () => {
        this.cleanExpiredSessions()
      },
      60 * 60 * 1000
    ) // 1小时

    return () => {
      clearInterval(intervalId)
    }
  }
}
