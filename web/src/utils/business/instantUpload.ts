import { calculateFileMD5 } from '../file/fileHash'
import imageAPI from '@/api/file'
import type { CheckDuplicateResponse, InstantUploadResponse } from '@/api/file/instant'

/**
 */
export class InstantUploadUtil {
  static async checkInstantUpload(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<{ canInstant: boolean; md5: string; checkResult?: CheckDuplicateResponse }> {
    try {
      const md5 = await calculateFileMD5(file, (progress) => {
        onProgress?.(progress * 0.8)
      })

      const checkResult = await imageAPI.checkDuplicate({
        md5,
        filename: file.name,
        file_size: file.size,
      })

      onProgress?.(100)

      return {
        canInstant: checkResult.exists,
        md5,
        checkResult,
      }
    } catch (error) {
      console.error('检查秒传失败:', error)
      throw error
    }
  }

  static async executeInstantUpload(
    file: File,
    md5: string,
    options?: {
      folder_id?: string
      access_level?: 'public' | 'private' | 'protected'
      optimize?: boolean
    }
  ): Promise<InstantUploadResponse> {
    try {
      const result = await imageAPI.instantUpload({
        md5,
        filename: file.name,
        file_size: file.size,
        folder_id: options?.folder_id,
        access_level: options?.access_level,
        optimize: options?.optimize,
      })

      return result
    } catch (error) {
      console.error('秒传失败:', error)
      throw error
    }
  }

  static async attemptInstantUpload(
    file: File,
    options?: {
      folder_id?: string
      access_level?: 'public' | 'private' | 'protected'
      optimize?: boolean
    },
    onProgress?: (progress: number) => void
  ): Promise<InstantUploadResponse | null> {
    try {
      const checkResult = await this.checkInstantUpload(file, (progress) => {
        onProgress?.(progress * 0.5)
      })

      if (!checkResult.canInstant) {
        return null
      }

      onProgress?.(75) // 开始秒传
      const uploadResult = await this.executeInstantUpload(file, checkResult.md5, options)
      onProgress?.(100) // 秒传完成

      return uploadResult
    } catch (error) {
      console.error('秒传流程失败:', error)
      throw error
    }
  }
}

export default InstantUploadUtil
