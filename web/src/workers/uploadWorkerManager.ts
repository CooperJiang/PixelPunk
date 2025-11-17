import * as Comlink from 'comlink'
import type { UploadWorkerAPI } from './upload.worker'
/* 使用Vite的Worker导入语法 */
import UploadWorker from './upload.worker?worker'

/**
 * Worker池配置
 */
interface WorkerPoolConfig {
  maxWorkers?: number
  workerPath?: string
}

/**
 * Worker任务
 */
interface WorkerTask<T = unknown> {
  id: string
  method: keyof UploadWorkerAPI
  args: unknown[]
  resolve: (value: T) => void
  reject: (_error: Error) => void
}

/**
 * Worker池项
 */
interface WorkerPoolItem {
  worker: Worker
  api: Comlink.Remote<UploadWorkerAPI>
  busy: boolean
  taskCount: number
}

/**
 * 上传Worker管理器
 * 管理Worker池，提供负载均衡和错误恢复
 */
export class UploadWorkerManager {
  private pool: WorkerPoolItem[] = []
  private taskQueue: WorkerTask[] = []
  private maxWorkers: number
  private initialized = false

  constructor(config: WorkerPoolConfig = {}) {
    this.maxWorkers = config.maxWorkers || 2 // 默认2个Worker
  }

  async initialize(): Promise<void> {
    if (this.initialized) {
      return
    }

    for (let i = 0; i < this.maxWorkers; i++) {
      try {
        const worker = new UploadWorker()

        const api = Comlink.wrap<UploadWorkerAPI>(worker)

        this.pool.push({
          worker,
          api,
          busy: false,
          taskCount: 0,
        })
      } catch {}
    }

    this.initialized = true
  }

  private getIdleWorker(): WorkerPoolItem | null {
    let selected = this.pool.find((w) => !w.busy)

    if (!selected) {
      selected = this.pool.reduce((prev, curr) => (prev.taskCount < curr.taskCount ? prev : curr))
    }

    return selected
  }

  private async executeTask<T>(task: WorkerTask<T>): Promise<void> {
    const worker = this.getIdleWorker()

    if (!worker) {
      this.taskQueue.push(task)
      return
    }

    worker.busy = true
    worker.taskCount++

    try {
      const result = await (worker.api[task.method] as (...args: unknown[]) => Promise<unknown>)(...task.args)
      task.resolve(result)
    } catch (error) {
      task.reject(error as Error)
    } finally {
      worker.busy = false
      worker.taskCount--

      this.processQueue()
    }
  }

  private processQueue(): void {
    if (this.taskQueue.length === 0) {
      return
    }

    const task = this.taskQueue.shift()
    if (task) {
      this.executeTask(task)
    }
  }

  async calculateMD5(file: File, onProgress?: (progress: number) => void): Promise<string> {
    if (!this.initialized) {
      await this.initialize()
    }

    return new Promise((resolve, reject) => {
      const task: WorkerTask<string> = {
        id: `md5-${Date.now()}`,
        method: 'calculateMD5',
        args: [file, onProgress ? Comlink.proxy(onProgress) : undefined],
        resolve,
        reject,
      }

      this.executeTask(task)
    })
  }

  async createChunks(file: File, chunkSize: number): Promise<Blob[]> {
    if (!this.initialized) {
      await this.initialize()
    }

    return new Promise((resolve, reject) => {
      const task: WorkerTask = {
        id: `chunk-${Date.now()}`,
        method: 'createChunks',
        args: [file, chunkSize],
        resolve,
        reject,
      }

      this.executeTask(task)
    })
  }

  async calculateChunkMD5(chunk: Blob): Promise<string> {
    if (!this.initialized) {
      await this.initialize()
    }

    return new Promise((resolve, reject) => {
      const task: WorkerTask<string> = {
        id: `chunk-md5-${Date.now()}`,
        method: 'calculateChunkMD5',
        args: [chunk],
        resolve,
        reject,
      }

      this.executeTask(task)
    })
  }

  async verifyFileIntegrity(file: File, expectedHash: string): Promise<boolean> {
    if (!this.initialized) {
      await this.initialize()
    }

    return new Promise((resolve, reject) => {
      const task: WorkerTask<boolean> = {
        id: `verify-${Date.now()}`,
        method: 'verifyFileIntegrity',
        args: [file, expectedHash],
        resolve,
        reject,
      }

      this.executeTask(task)
    })
  }

  destroy(): void {
    for (const item of this.pool) {
      item.worker.terminate()
    }

    this.pool = []
    this.taskQueue = []
    this.initialized = false
  }

  getStatus() {
    return {
      initialized: this.initialized,
      workers: this.pool.length,
      busyWorkers: this.pool.filter((w) => w.busy).length,
      queuedTasks: this.taskQueue.length,
      totalTasks: this.pool.reduce((sum, w) => sum + w.taskCount, 0),
    }
  }
}

export const uploadWorkerManager = new UploadWorkerManager({
  maxWorkers: navigator.hardwareConcurrency > 4 ? 3 : 2, // 根据CPU核心数决定Worker数量
})
