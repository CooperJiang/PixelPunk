/**
 * 浏览器指纹生成工具
 * 用于游客上传时的身份识别

 */
export interface FingerprintOptions {
  enableCanvas?: boolean
  enableWebGL?: boolean
  enableAudio?: boolean
  enableScreen?: boolean
}

/**
 * 生成简单的浏览器指纹
 * 基于用户代理、屏幕分辨率、时区、语言等基本信息
 */
export function generateSimpleFingerprint(): string {
  const components: string[] = []

  components.push(navigator.userAgent)

  components.push(`${screen.width}x${screen.height}x${screen.colorDepth}`)

  components.push(Intl.DateTimeFormat().resolvedOptions().timeZone)

  components.push(navigator.language)

  components.push(navigator.platform)

  components.push(navigator.hardwareConcurrency?.toString() || 'unknown')

  components.push((navigator as any).deviceMemory?.toString() || 'unknown')

  components.push(navigator.maxTouchPoints > 0 ? 'touch' : 'notouch')

  const fingerprint = components.join('|')

  let hash = 0
  for (let i = 0; i < fingerprint.length; i++) {
    const char = fingerprint.charCodeAt(i)
    hash = (hash << 5) - hash + char
    hash = hash & hash // 转换为32位整数
  }

  return Math.abs(hash).toString(36)
}

export function generateEnhancedFingerprint(options?: FingerprintOptions): string {
  const baseFingerprint = generateSimpleFingerprint()
  const enhancedComponents: string[] = [baseFingerprint]

  if (options?.enableCanvas !== false) {
    try {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.textBaseline = 'top'
        ctx.font = '14px Arial'
        ctx.fillStyle = '#f60'
        ctx.fillRect(125, 1, 62, 20)
        ctx.fillStyle = '#069'
        ctx.fillText('🤖', 2, 15)
        ctx.fillStyle = 'rgba(102, 204, 0, 0.7)'
        ctx.fillText('🤖', 4, 17)
        enhancedComponents.push(canvas.toDataURL())
      }
    } catch {}
  }

  if (options?.enableWebGL !== false) {
    try {
      const canvas = document.createElement('canvas')
      const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl')
      if (gl) {
        const debugInfo = gl.getExtension('WEBGL_debug_renderer_info')
        if (debugInfo) {
          const vendor = gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL)
          const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL)
          enhancedComponents.push(`${vendor}|${renderer}`)
        }
      }
    } catch {}
  }

  const combined = enhancedComponents.join('|')
  let hash = 0
  for (let i = 0; i < combined.length; i++) {
    const char = combined.charCodeAt(i)
    hash = (hash << 5) - hash + char
    hash = hash & hash
  }

  return Math.abs(hash).toString(36)
}

export default generateSimpleFingerprint
