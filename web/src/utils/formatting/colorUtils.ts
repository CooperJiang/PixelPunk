/**
 * 颜色处理相关工具函数
 */

/**
 * RGB转HSL颜色
 * @param r 红色通道值 (0-255)
 * @param g 绿色通道值 (0-255)
 * @param b 蓝色通道值 (0-255)
 * @returns HSL颜色值 [h, s, l]
 */
export const rgbToHsl = (r: number, g: number, b: number): [number, number, number] => {
  r /= 255
  g /= 255
  b /= 255

  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0,
    s = 0
  const l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)

    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0)
        break
      case g:
        h = (b - r) / d + 2
        break
      case b:
        h = (r - g) / d + 4
        break
    }

    h /= 6
  }

  return [h, s, l]
}

export const hslToRgb = (h: number, s: number, l: number): [number, number, number] => {
  let r, g, b

  if (s === 0) {
    r = g = b = l
  } else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) {
        t += 1
      }
      if (t > 1) {
        t -= 1
      }
      if (t < 1 / 6) {
        return p + (q - p) * 6 * t
      }
      if (t < 1 / 2) {
        return q
      }
      if (t < 2 / 3) {
        return p + (q - p) * (2 / 3 - t) * 6
      }
      return p
    }

    const q = l < 0.5 ? l * (1 + s) : l + s - l * s
    const p = 2 * l - q

    r = hue2rgb(p, q, h + 1 / 3)
    g = hue2rgb(p, q, h)
    b = hue2rgb(p, q, h - 1 / 3)
  }

  return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)]
}

export const extractImageColor = async (imageUrl: string): Promise<string> => {
  const { useTextThemeStore } = await import('@/store/textTheme')
  const store = useTextThemeStore()

  return new Promise((resolve, reject) => {
    if (!imageUrl) {
      reject(store.getText('utils.file.colorUtils.errors.invalidImageUrl'))
      return
    }

    try {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      if (!ctx) {
        reject(store.getText('utils.file.colorUtils.errors.cannotCreateCanvasContext'))
        return
      }

      const img = new Image()
      img.crossOrigin = 'Anonymous' // 尝试解决跨域问题

      img.onload = () => {
        canvas.width = img.width
        canvas.height = img.height

        ctx.drawImage(img, 0, 0)

        let data
        try {
          data = ctx.getImageData(0, 0, canvas.width, canvas.height).data
        } catch {
          reject(store.getText('utils.file.colorUtils.errors.cannotGetImageData'))
          return
        }

        let r = 0,
          g = 0,
          b = 0
        const sampleSize = Math.floor(data.length / 16) // 只取前1/4的像素
        let count = 0

        for (let i = 0; i < sampleSize; i += 4) {
          r += data[i]
          g += data[i + 1]
          b += data[i + 2]
          count++
        }

        r = Math.floor(r / count)
        g = Math.floor(g / count)
        b = Math.floor(b / count)

        const hsl = rgbToHsl(r, g, b)
        hsl[1] = Math.min(hsl[1], 0.3) // 降低饱和度
        hsl[2] = Math.max(0.1, Math.min(0.25, hsl[2])) // 控制亮度在合适范围

        const [rNew, gNew, bNew] = hslToRgb(hsl[0], hsl[1], hsl[2])

        const backgroundColor = `rgba(${rNew}, ${gNew}, ${bNew}, 0.92)`
        resolve(backgroundColor)
      }

      img.onerror = () => {
        reject(store.getText('utils.file.colorUtils.errors.imageLoadFailed'))
      }

      img.src = imageUrl
    } catch (error) {
      const formatText = (template: string, params: Record<string, string | number>): string => {
        let result = template
        Object.keys(params).forEach((key) => {
          result = result.replace(new RegExp(`\\{${key}\\}`, 'g'), String(params[key]))
        })
        return result
      }
      reject(formatText(store.getText('utils.file.colorUtils.errors.extractColorFailed'), { error: String(error) }))
    }
  })
}
