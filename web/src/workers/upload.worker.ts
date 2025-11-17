import * as Comlink from 'comlink'
import SparkMD5 from 'spark-md5'

/**
 * Worker消息类型
 */
export interface WorkerProgress {
  type: 'md5' | 'chunk'
  progress: number
  message?: string
}

/**
 * 文件分片信息
 */
export interface ChunkInfo {
  index: number
  start: number
  end: number
  blob: Blob
  hash?: string
}

/**
 * Worker API接口
 */
export interface UploadWorkerAPI {
  calculateMD5: (file: File, onProgress?: (progress: number) => void) => Promise<string>
  createChunks: (file: File, chunkSize: number) => Promise<ChunkInfo[]>
  calculateChunkMD5: (chunk: Blob) => Promise<string>
  verifyFileIntegrity: (file: File, expectedHash: string) => Promise<boolean>
}

/**
 * 上传Worker实现
 */
class UploadWorker implements UploadWorkerAPI {
  async calculateMD5(file: File, onProgress?: (progress: number) => void): Promise<string> {
    return new Promise((resolve, reject) => {
      const chunkSize = 2 * 1024 * 1024 // 2MB chunks for MD5 calculation
      const chunks = Math.ceil(file.size / chunkSize)
      const spark = new SparkMD5.ArrayBuffer()
      let currentChunk = 0

      const fileReader = new FileReader()

      fileReader.onload = (e) => {
        if (e.target?.result) {
          spark.append(e.target.result as ArrayBuffer)
          currentChunk++

          if (onProgress) {
            const progress = Math.round((currentChunk / chunks) * 100)
            onProgress(progress)
          }

          if (currentChunk < chunks) {
            loadNext()
          } else {
            const hash = spark.end()
            resolve(hash)
          }
        }
      }

      fileReader.onerror = () => {
        reject(new Error('MD5 calculation failed'))
      }

      const loadNext = () => {
        const start = currentChunk * chunkSize
        const end = Math.min(start + chunkSize, file.size)
        fileReader.readAsArrayBuffer(file.slice(start, end))
      }

      loadNext()
    })
  }

  async createChunks(file: File, chunkSize: number): Promise<ChunkInfo[]> {
    const chunks: ChunkInfo[] = []
    const totalChunks = Math.ceil(file.size / chunkSize)

    for (let i = 0; i < totalChunks; i++) {
      const start = i * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      const blob = file.slice(start, end)

      chunks.push({
        index: i,
        start,
        end,
        blob,
      })
    }

    for (const chunk of chunks) {
      chunk.hash = await this.calculateChunkMD5(chunk.blob)
    }

    return chunks
  }

  async calculateChunkMD5(chunk: Blob): Promise<string> {
    return new Promise((resolve, reject) => {
      const fileReader = new FileReader()

      fileReader.onload = (e) => {
        if (e.target?.result) {
          const spark = new SparkMD5.ArrayBuffer()
          spark.append(e.target.result as ArrayBuffer)
          resolve(spark.end())
        }
      }

      fileReader.onerror = () => {
        reject(new Error('Chunk MD5 calculation failed'))
      }

      fileReader.readAsArrayBuffer(chunk)
    })
  }

  async verifyFileIntegrity(file: File, expectedHash: string): Promise<boolean> {
    try {
      const actualHash = await this.calculateMD5(file)
      return actualHash === expectedHash
    } catch {
      return false
    }
  }
}

const worker = new UploadWorker()
Comlink.expose(worker)
