/**
 * HEIC/HEIF 格式工具
 * 注意：前端不对 HEIC 进行转换，后端会自动处理
 */

/**
 * 检查文件是否为 HEIC/HEIF 格式
 */
export function isHEICFile(file: File): boolean {
  const fileName = file.name.toLowerCase()
  return fileName.endsWith('.heic') || fileName.endsWith('.heif')
}

function getHEICPlaceholder(): string {
  const svg = `<svg width="200" height="200" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#1e1e1e"/>
        <stop offset="100%" style="stop-color:#2d2d2d"/>
      </linearGradient>
      <linearGradient id="icon" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#6366f1"/>
        <stop offset="100%" style="stop-color:#8b5cf6"/>
      </linearGradient>
    </defs>
    <rect width="200" height="200" fill="url(#bg)" rx="8"/>
    <!-- 照片图标：山+太阳 -->
    <g transform="translate(100, 100)">
      <!-- 外框 -->
      <rect x="-45" y="-35" width="90" height="70" rx="6" fill="none" stroke="url(#icon)" stroke-width="4"/>
      <!-- 太阳 -->
      <circle cx="-20" cy="-15" r="10" fill="url(#icon)" opacity="0.8"/>
      <!-- 山峰 -->
      <path d="M -45 35 L -15 -5 L 10 20 L 45 -10 L 45 35 Z" fill="url(#icon)" opacity="0.6"/>
      <!-- 前景山 -->
      <path d="M -45 35 L -10 5 L 20 35 Z" fill="url(#icon)"/>
    </g>
  </svg>`

  return `data:image/svg+xml;base64,${btoa(unescape(encodeURIComponent(svg)))}`
}

export async function createFileObjectURL(file: File): Promise<string> {
  if (isHEICFile(file)) {
    return getHEICPlaceholder()
  }
  return URL.createObjectURL(file)
}

export async function createFileDataURL(file: File): Promise<string> {
  if (isHEICFile(file)) {
    return getHEICPlaceholder()
  }

  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      resolve(e.target?.result as string)
    }
    reader.onerror = (e) => {
      reject(e)
    }
    reader.readAsDataURL(file)
  })
}
