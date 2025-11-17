import type { WatermarkPosition, ImageRect } from './types'
import { getBrandColor, getDangerColor } from './themeColors'

function drawDistanceLabel(ctx: CanvasRenderingContext2D, text: string, x: number, y: number) {
  const padding = 4
  const metrics = ctx.measureText(text)
  const textWidth = metrics.width
  const textHeight = 11

  ctx.fillStyle = getBrandColor(0.9)
  ctx.fillRect(x - textWidth / 2 - padding, y - textHeight / 2 - padding, textWidth + padding * 2, textHeight + padding * 2)

  ctx.fillStyle = '#ffffff'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(text, x, y)
}

export const drawGuideLines = (
  ctx: CanvasRenderingContext2D,
  watermarkPos: WatermarkPosition,
  imgRect: ImageRect,
  position: string,
  isDragging: boolean
) => {
  if (!isDragging) return

  ctx.save()
  ctx.setLineDash([4, 4])
  ctx.strokeStyle = getBrandColor(1)
  ctx.lineWidth = 1
  ctx.fillStyle = getBrandColor(1)
  ctx.font = 'bold 11px Arial'

  const wmX = watermarkPos.x
  const wmY = watermarkPos.y
  const wmW = watermarkPos.width
  const wmH = watermarkPos.height

  const showLeft = position.includes('left')
  const showRight = position.includes('right')
  const showTop = position.includes('top')
  const showBottom = position.includes('bottom')

  if (showLeft) {
    const dist = Math.round(wmX - imgRect.x)
    ctx.beginPath()
    ctx.moveTo(imgRect.x, wmY + wmH / 2)
    ctx.lineTo(wmX, wmY + wmH / 2)
    ctx.stroke()

    const labelX = imgRect.x + (wmX - imgRect.x) / 2
    const labelY = wmY + wmH / 2 - 5
    drawDistanceLabel(ctx, `${dist}px`, labelX, labelY)
  }

  if (showRight) {
    const dist = Math.round(imgRect.x + imgRect.width - (wmX + wmW))
    ctx.beginPath()
    ctx.moveTo(wmX + wmW, wmY + wmH / 2)
    ctx.lineTo(imgRect.x + imgRect.width, wmY + wmH / 2)
    ctx.stroke()

    const labelX = wmX + wmW + dist / 2
    const labelY = wmY + wmH / 2 - 5
    drawDistanceLabel(ctx, `${dist}px`, labelX, labelY)
  }

  if (showTop) {
    const dist = Math.round(wmY - imgRect.y)
    ctx.beginPath()
    ctx.moveTo(wmX + wmW / 2, imgRect.y)
    ctx.lineTo(wmX + wmW / 2, wmY)
    ctx.stroke()

    const labelX = wmX + wmW / 2 + 5
    const labelY = imgRect.y + (wmY - imgRect.y) / 2
    drawDistanceLabel(ctx, `${dist}px`, labelX, labelY)
  }

  if (showBottom) {
    const dist = Math.round(imgRect.y + imgRect.height - (wmY + wmH))
    ctx.beginPath()
    ctx.moveTo(wmX + wmW / 2, wmY + wmH)
    ctx.lineTo(wmX + wmW / 2, imgRect.y + imgRect.height)
    ctx.stroke()

    const labelX = wmX + wmW / 2 + 5
    const labelY = wmY + wmH + dist / 2
    drawDistanceLabel(ctx, `${dist}px`, labelX, labelY)
  }

  ctx.restore()
}

export const drawPreviewAreaBorder = (
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  previewWidth: number,
  previewHeight: number,
  labelText?: string
) => {
  ctx.save()

  ctx.strokeStyle = getBrandColor(0.8)
  ctx.lineWidth = 3
  ctx.setLineDash([10, 5])
  ctx.strokeRect(x, y, width, height)
  ctx.setLineDash([])

  const cornerSize = 8
  ctx.fillStyle = getBrandColor(0.9)
  ctx.fillRect(x - cornerSize / 2, y - cornerSize / 2, cornerSize, cornerSize)
  ctx.fillRect(x + width - cornerSize / 2, y - cornerSize / 2, cornerSize, cornerSize)
  ctx.fillRect(x - cornerSize / 2, y + height - cornerSize / 2, cornerSize, cornerSize)
  ctx.fillRect(x + width - cornerSize / 2, y + height - cornerSize / 2, cornerSize, cornerSize)

  ctx.fillStyle = getBrandColor(0.95)
  ctx.font = 'bold 13px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif'
  ctx.textAlign = 'left'
  ctx.textBaseline = 'top'

  const finalLabelText = labelText || `Preview Area: ${previewWidth} × ${previewHeight}px`
  const labelPadding = 8
  const labelMetrics = ctx.measureText(finalLabelText)
  const labelWidth = labelMetrics.width + labelPadding * 2
  const labelHeight = 24

  ctx.fillStyle = getBrandColor(0.95)
  ctx.fillRect(x + 10, y + 10, labelWidth, labelHeight)

  ctx.fillStyle = '#ffffff'
  ctx.fillText(finalLabelText, x + 10 + labelPadding, y + 10 + 6)

  ctx.restore()
}

export const drawAnchorPoint = (ctx: CanvasRenderingContext2D, x: number, y: number) => {
  ctx.save()

  const gradient = ctx.createRadialGradient(x, y, 0, x, y, 12)
  gradient.addColorStop(0, getDangerColor(0.4))
  gradient.addColorStop(1, getDangerColor(0))
  ctx.fillStyle = gradient
  ctx.beginPath()
  ctx.arc(x, y, 12, 0, Math.PI * 2)
  ctx.fill()

  ctx.fillStyle = getDangerColor(1)
  ctx.beginPath()
  ctx.arc(x, y, 6, 0, Math.PI * 2)
  ctx.fill()

  ctx.strokeStyle = '#ffffff'
  ctx.lineWidth = 2
  ctx.stroke()

  ctx.fillStyle = 'rgba(255, 255, 255, 0.5)'
  ctx.beginPath()
  ctx.arc(x - 1, y - 1, 2, 0, Math.PI * 2)
  ctx.fill()

  ctx.restore()
}

export const drawAnchorLine = (
  ctx: CanvasRenderingContext2D,
  anchorX: number,
  anchorY: number,
  watermarkCenterX: number,
  watermarkCenterY: number
) => {
  ctx.save()

  ctx.strokeStyle = getDangerColor(0.5)
  ctx.lineWidth = 2
  ctx.setLineDash([8, 6])

  ctx.beginPath()
  ctx.moveTo(anchorX, anchorY)
  ctx.lineTo(watermarkCenterX, watermarkCenterY)
  ctx.stroke()

  ctx.setLineDash([])
  ctx.restore()
}

export const drawResizeHandles = (
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  hoveredCorner: string | null = null
) => {
  const handleSize = 12
  const handles = [
    { name: 'top-left', x, y },
    { name: 'top-right', x: x + width, y },
    { name: 'bottom-left', x, y: y + height },
    { name: 'bottom-right', x: x + width, y: y + height },
  ]

  ctx.save()

  handles.forEach((handle) => {
    const isHovered = hoveredCorner === handle.name

    if (isHovered) {
      const gradient = ctx.createRadialGradient(handle.x, handle.y, 0, handle.x, handle.y, handleSize * 1.5)
      gradient.addColorStop(0, getBrandColor(0.4))
      gradient.addColorStop(1, getBrandColor(0))
      ctx.fillStyle = gradient
      ctx.beginPath()
      ctx.arc(handle.x, handle.y, handleSize * 1.5, 0, Math.PI * 2)
      ctx.fill()
    }

    ctx.fillStyle = isHovered ? getBrandColor(0.8) : getBrandColor(1)
    ctx.fillRect(handle.x - handleSize / 2, handle.y - handleSize / 2, handleSize, handleSize)

    ctx.strokeStyle = '#ffffff'
    ctx.lineWidth = 2
    ctx.strokeRect(handle.x - handleSize / 2, handle.y - handleSize / 2, handleSize, handleSize)

    if (isHovered) {
      ctx.fillStyle = 'rgba(255, 255, 255, 0.4)'
      ctx.fillRect(handle.x - handleSize / 2 + 2, handle.y - handleSize / 2 + 2, handleSize - 4, handleSize - 4)
    }
  })

  ctx.restore()
}
