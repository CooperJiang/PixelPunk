import type { WatermarkConfig } from '../../types'
import type { WatermarkPosition } from './types'

export const drawTextWatermark = (ctx: CanvasRenderingContext2D, config: WatermarkConfig, _position: WatermarkPosition) => {
  ctx.font = `${config.fontWeight} ${config.fontSize * config.scale}px ${config.fontFamily}`
  ctx.fillStyle = config.fontColor
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(config.text, 0, 0)
}

export const drawImageWatermark = (ctx: CanvasRenderingContext2D, img: HTMLImageElement, position: WatermarkPosition) => {
  const { width, height } = position
  ctx.drawImage(img, -width / 2, -height / 2, width, height)
}

export const drawWatermark = (
  ctx: CanvasRenderingContext2D,
  config: WatermarkConfig,
  watermarkImageData: HTMLImageElement | undefined,
  position: WatermarkPosition
) => {
  ctx.save()
  ctx.globalAlpha = config.opacity

  if (config.shadow) {
    ctx.shadowColor = config.shadowColor
    ctx.shadowBlur = config.shadowBlur
    ctx.shadowOffsetX = config.shadowOffsetX
    ctx.shadowOffsetY = config.shadowOffsetY
  }

  const centerX = position.x + position.width / 2
  const centerY = position.y + position.height / 2
  ctx.translate(centerX, centerY)

  if (config.rotation !== 0) {
    ctx.rotate((config.rotation * Math.PI) / 180)
  }

  if (config.type === 'text') {
    drawTextWatermark(ctx, config, position)
  } else if (config.type === 'image' && watermarkImageData) {
    drawImageWatermark(ctx, watermarkImageData, position)
  }

  ctx.restore()
}
