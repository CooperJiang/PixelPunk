import { getImageRect } from './calculations'

export const drawGrid = (ctx: CanvasRenderingContext2D, width: number, height: number) => {
  ctx.strokeStyle = 'rgba(139, 92, 246, 0.1)'
  ctx.lineWidth = 1

  const gridSize = 20

  for (let x = 0; x <= width; x += gridSize) {
    ctx.beginPath()
    ctx.moveTo(x, 0)
    ctx.lineTo(x, height)
    ctx.stroke()
  }

  for (let y = 0; y <= height; y += gridSize) {
    ctx.beginPath()
    ctx.moveTo(0, y)
    ctx.lineTo(width, y)
    ctx.stroke()
  }
}

export const drawPlaceholder = (
  ctx: CanvasRenderingContext2D,
  width: number,
  height: number,
  text1 = 'Preview Area',
  text2 = 'Drag file here to preview watermark effect'
) => {
  ctx.fillStyle = 'rgba(75, 85, 99, 0.3)'
  ctx.fillRect(0, 0, width, height)

  ctx.fillStyle = 'rgba(156, 163, 175, 0.8)'
  ctx.font = '16px Arial'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(text1, width / 2, height / 2)
  ctx.fillText(text2, width / 2, height / 2 + 25)
}

export const drawPreviewImage = (
  ctx: CanvasRenderingContext2D,
  img: HTMLImageElement,
  canvasWidth: number,
  canvasHeight: number
) => {
  const rect = getImageRect(img, canvasWidth, canvasHeight)
  ctx.drawImage(img, rect.x, rect.y, rect.width, rect.height)
}
