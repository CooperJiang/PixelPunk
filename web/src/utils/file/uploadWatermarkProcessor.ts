import type { WatermarkConfig } from '@/components/WatermarkConfig/types'
import { type CanvasWatermarkConfig, generateWatermark } from './canvasWatermarkGenerator'
import { useTextThemeStore } from '@/store/textTheme'

/**
 * 格式化翻译文本，替换参数
 */
function formatText(template: string, params: Record<string, string | number>): string {
  let result = template
  Object.keys(params).forEach((key) => {
    result = result.replace(new RegExp(`\\{${key}\\}`, 'g'), String(params[key]))
  })
  return result
}

export async function getFileDimensions(file: File): Promise<{ width: number; height: number }> {
  const { createFileObjectURL } = await import('./heicConverter')
  const store = useTextThemeStore()

  return new Promise((resolve, reject) => {
    const img = new Image()

    createFileObjectURL(file)
      .then((objectUrl) => {
        img.onload = () => {
          resolve({
            width: img.naturalWidth,
            height: img.naturalHeight,
          })
          URL.revokeObjectURL(objectUrl)
        }

        img.onerror = () => {
          reject(new Error(store.getText('utils.file.watermark.errors.cannotLoadImage')))
          URL.revokeObjectURL(objectUrl)
        }

        img.src = objectUrl
      })
      .catch((error) => {
        reject(error)
      })
  })
}

export async function processFileWatermark(file: File, config: WatermarkConfig): Promise<WatermarkConfig> {
  if (!config.enabled || !config.generateMode || config.generateMode !== 'frontend') {
    return config
  }

  try {
    const dimensions = await getFileDimensions(file)

    const canvasConfig: CanvasWatermarkConfig = {
      ...config,
      targetWidth: dimensions.width,
      targetHeight: dimensions.height,
    }

    const result = await generateWatermark(canvasConfig)

    if (result.success && result.base64Data) {
      const updatedConfig: WatermarkConfig = {
        ...config,
        type: 'image', // 统一转为文件类型
        generateMode: 'frontend',
        generatedImage: result.base64Data,
        generatedWidth: result.width,
        generatedHeight: result.height,
        fileBase64: result.base64Data.replace(/^data:image\/\w+;base64,/, ''), // 去除前缀
      }

      return updatedConfig
    }
    return {
      ...config,
      enabled: false,
    }
  } catch (_error) {
    console.error(`[水印处理] 文件 ${file.name} 处理出错:`, _error)
    return {
      ...config,
      enabled: false,
    }
  }
}

export async function processFilesWatermark(
  files: File[],
  config: WatermarkConfig
): Promise<{
  processedConfig: WatermarkConfig
  processedCount: number
  failedCount: number
}> {
  if (!config.enabled || !config.generateMode || config.generateMode !== 'frontend') {
    return {
      processedConfig: config,
      processedCount: 0,
      failedCount: 0,
    }
  }

  const startTime = Date.now()

  let processedCount = 0
  let failedCount = 0
  let lastSuccessConfig = config

  for (const file of files) {
    if (file.type.startsWith('image/')) {
      try {
        const processedConfig = await processFileWatermark(file, config)
        if (processedConfig.enabled && processedConfig.generatedImage) {
          processedCount++
          lastSuccessConfig = processedConfig
          break // 使用第一个成功的配置应用到所有文件
        } else {
          failedCount++
        }
      } catch (_error) {
        failedCount++
      }
    }
  }

  const _totalTime = Date.now() - startTime

  return {
    processedConfig: processedCount > 0 ? lastSuccessConfig : { ...config, enabled: false },
    processedCount,
    failedCount,
  }
}

export function shouldGenerateWatermark(config: WatermarkConfig): boolean {
  const result = Boolean(config.enabled && config.generateMode === 'frontend' && config.type === 'image' && config.fileBase64)

  return result
}

export function formatWatermarkProcessMessage(processedCount: number, failedCount: number, totalFiles: number): string {
  if (processedCount === 0 && failedCount === 0) {
    return ''
  }

  const store = useTextThemeStore()

  if (failedCount === 0) {
    return formatText(store.getText('utils.file.watermark.status.completed'), { processedCount, totalFiles })
  } else if (processedCount === 0) {
    return formatText(store.getText('utils.file.watermark.status.failed'), { failedCount, totalFiles })
  }
  return formatText(store.getText('utils.file.watermark.status.partialSuccess'), { processedCount, failedCount })
}
