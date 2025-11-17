import type { WatermarkConfig } from '@/components/WatermarkConfig/types'
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

/* Canvas水印生成配置接口 */
export interface CanvasWatermarkConfig extends WatermarkConfig {
  targetWidth: number // 目标文件宽度
  targetHeight: number // 目标文件高度
}

/* 水印生成结果 */
export interface WatermarkGenerateResult {
  success: boolean
  base64Data?: string // data:image/png;base64,xxx
  width: number
  height: number
  error?: string
  generateTime?: number // 生成耗时(ms)
}

/* Canvas水印生成器类 */
export class CanvasWatermarkGenerator {
  private canvas: HTMLCanvasElement
  private ctx: CanvasRenderingContext2D

  constructor() {
    this.canvas = document.createElement('canvas')
    const ctx = this.canvas.getContext('2d', { alpha: true })
    if (!ctx) {
      const store = useTextThemeStore()
      throw new Error(store.getText('utils.file.watermark.errors.canvasContextFailed'))
    }
    this.ctx = ctx
  }

  async generateWatermark(config: CanvasWatermarkConfig): Promise<WatermarkGenerateResult> {
    const startTime = Date.now()

    try {
      if (!config.enabled) {
        return {
          success: true,
          base64Data: '',
          width: 0,
          height: 0,
          generateTime: 0,
        }
      }

      let watermarkWidth: number, watermarkHeight: number

      const store = useTextThemeStore()

      if (config.type === 'text') {
        const textMetrics = this.calculateTextSize(config)
        watermarkWidth = Math.ceil(textMetrics.width)
        watermarkHeight = Math.ceil(textMetrics.height)
      } else if (config.type === 'image') {
        const img = await this.loadWatermarkImage(config)
        watermarkWidth = Math.ceil(img.width * config.scale)
        watermarkHeight = Math.ceil(img.height * config.scale)
      } else {
        throw new Error(store.getText('utils.file.watermark.errors.unsupportedWatermarkType'))
      }

      let canvasWidth = watermarkWidth
      let canvasHeight = watermarkHeight
      if (config.rotation !== 0) {
        const rad = (Math.abs(config.rotation) * Math.PI) / 180
        const cos = Math.abs(Math.cos(rad))
        const sin = Math.abs(Math.sin(rad))
        canvasWidth = Math.ceil(watermarkWidth * cos + watermarkHeight * sin)
        canvasHeight = Math.ceil(watermarkWidth * sin + watermarkHeight * cos)
      }

      this.setupCanvas(canvasWidth, canvasHeight)

      this.ctx.clearRect(0, 0, canvasWidth, canvasHeight)

      if (config.type === 'text') {
        await this.generateTextWatermarkContent(config, canvasWidth, canvasHeight)
      } else if (config.type === 'image') {
        await this.generateImageWatermarkContent(config, canvasWidth, canvasHeight)
      }

      const base64Data = this.canvas.toDataURL('image/png')

      return {
        success: true,
        base64Data,
        width: canvasWidth, // 返回旋转后的包围盒尺寸
        height: canvasHeight,
        generateTime: Date.now() - startTime,
      }
    } catch (error) {
      console.error('[Canvas水印生成器] 生成失败:', error)
      const store = useTextThemeStore()
      return {
        success: false,
        width: 0,
        height: 0,
        error: error instanceof Error ? error.message : store.getText('utils.file.watermark.errors.watermarkGenerationFailed'),
        generateTime: Date.now() - startTime,
      }
    }
  }

  private setupCanvas(width: number, height: number) {
    this.canvas.width = width
    this.canvas.height = height

    this.ctx.imageSmoothingEnabled = true
    this.ctx.imageSmoothingQuality = 'high'

    this.ctx.setTransform(1, 0, 0, 1, 0, 0)
  }

  private calculateTextSize(config: CanvasWatermarkConfig): { width: number; height: number } {
    const store = useTextThemeStore()
    const tempCanvas = document.createElement('canvas')
    const tempCtx = tempCanvas.getContext('2d')
    if (!tempCtx) {
      throw new Error(store.getText('utils.file.watermark.errors.tempCanvasContextFailed'))
    }

    const actualFontSize = config.fontSize * config.scale
    tempCtx.font = `${config.fontWeight} ${actualFontSize}px ${config.fontFamily}`

    const metrics = tempCtx.measureText(config.text)
    const textWidth = metrics.width
    const textHeight = actualFontSize

    const padding = config.shadow
      ? Math.max(config.shadowBlur, Math.abs(config.shadowOffsetX), Math.abs(config.shadowOffsetY))
      : 0

    return {
      width: textWidth + padding * 2,
      height: textHeight + padding * 2,
    }
  }

  private async generateTextWatermarkContent(config: CanvasWatermarkConfig, width: number, height: number) {
    const { text, fontSize, fontFamily, fontColor, fontWeight: _fontWeight, opacity } = config

    await this.ensureFontLoaded(fontFamily)

    const actualFontSize = fontSize * config.scale

    this.ctx.font = `${config.fontWeight} ${actualFontSize}px ${config.fontFamily}`
    this.ctx.textAlign = 'center'
    this.ctx.textBaseline = 'middle'

    const color = this.parseColor(fontColor)
    this.ctx.fillStyle = this.applyOpacity(color, opacity)

    const centerX = width / 2
    const centerY = height / 2

    if (config.rotation !== 0) {
      this.ctx.save()
      this.ctx.translate(centerX, centerY)
      this.ctx.rotate((config.rotation * Math.PI) / 180)
      this.ctx.translate(-centerX, -centerY)
    }

    if (config.shadow) {
      this.ctx.save()
      this.ctx.fillStyle = this.applyOpacity(this.parseColor(config.shadowColor), opacity * 0.5)
      this.ctx.fillText(text, centerX + config.shadowOffsetX, centerY + config.shadowOffsetY)
      this.ctx.restore()
    }

    this.ctx.fillStyle = this.applyOpacity(color, opacity)
    this.ctx.fillText(text, centerX, centerY)

    if (config.rotation !== 0) {
      this.ctx.restore()
    }
  }

  private async generateImageWatermarkContent(config: CanvasWatermarkConfig, width: number, height: number) {
    const watermarkImage = await this.loadWatermarkImage(config)

    const scaledWidth = watermarkImage.width * config.scale
    const scaledHeight = watermarkImage.height * config.scale

    const centerX = width / 2
    const centerY = height / 2

    this.ctx.globalAlpha = config.opacity

    if (config.rotation !== 0) {
      this.ctx.save()
      this.ctx.translate(centerX, centerY)
      this.ctx.rotate((config.rotation * Math.PI) / 180)
      this.ctx.translate(-centerX, -centerY)
    }

    if (config.shadow) {
      this.ctx.save()
      this.ctx.globalAlpha = config.opacity * 0.3
      this.ctx.fillStyle = config.shadowColor
      this.ctx.fillRect(
        centerX - scaledWidth / 2 + config.shadowOffsetX,
        centerY - scaledHeight / 2 + config.shadowOffsetY,
        scaledWidth,
        scaledHeight
      )
      this.ctx.restore()
      this.ctx.globalAlpha = config.opacity
    }

    this.ctx.drawImage(watermarkImage, centerX - scaledWidth / 2, centerY - scaledHeight / 2, scaledWidth, scaledHeight)

    if (config.rotation !== 0) {
      this.ctx.restore()
    }

    this.ctx.globalAlpha = 1
  }

  private async loadWatermarkImage(config: CanvasWatermarkConfig): Promise<HTMLImageElement> {
    if (config.imageFile) {
      return await this.loadImageFromFile(config.imageFile)
    } else if (config.imageUrl) {
      return await this.loadImageFromUrl(config.imageUrl)
    }
    const store = useTextThemeStore()
    throw new Error(store.getText('utils.file.watermark.errors.noWatermarkFile'))
  }

  private async ensureFontLoaded(fontFamily: string): Promise<boolean> {
    if (!document.fonts) {
      return true
    }

    try {
      await document.fonts.load(`16px "${fontFamily}"`)
      return document.fonts.check(`16px "${fontFamily}"`)
    } catch (error) {
      console.warn(`字体加载失败: ${fontFamily}`, error)
      return false
    }
  }

  private async loadImageFromFile(file: File): Promise<HTMLImageElement> {
    const { createFileObjectURL } = await import('./heicConverter')
    const store = useTextThemeStore()

    return new Promise((resolve, reject) => {
      const img = new Image()

      createFileObjectURL(file)
        .then((url) => {
          img.onload = () => {
            URL.revokeObjectURL(url)
            resolve(img)
          }

          img.onerror = () => {
            URL.revokeObjectURL(url)
            reject(new Error(store.getText('utils.file.watermark.errors.fileLoadFailed')))
          }

          img.src = url
        })
        .catch((error) => {
          reject(error)
        })
    })
  }

  private loadImageFromUrl(url: string): Promise<HTMLImageElement> {
    return new Promise((resolve, reject) => {
      const store = useTextThemeStore()
      const img = new Image()

      img.onload = () => resolve(img)
      img.onerror = () => reject(new Error(formatText(store.getText('utils.file.watermark.errors.imageLoadFailed'), { url })))

      img.crossOrigin = 'anonymous'
      img.src = url
    })
  }

  private parseColor(colorStr: string): string {
    if (colorStr.startsWith('#')) {
      return colorStr
    }
    return `#${colorStr}`
  }

  private applyOpacity(color: string, opacity: number): string {
    const hex = color.replace('#', '')
    const r = parseInt(hex.substr(0, 2), 16)
    const g = parseInt(hex.substr(2, 2), 16)
    const b = parseInt(hex.substr(4, 2), 16)

    return `rgba(${r}, ${g}, ${b}, ${opacity})`
  }

  destroy() {
    this.canvas.remove()
  }
}

export const generateWatermark = async (config: CanvasWatermarkConfig): Promise<WatermarkGenerateResult> => {
  const generator = new CanvasWatermarkGenerator()
  try {
    return await generator.generateWatermark(config)
  } finally {
    generator.destroy()
  }
}
