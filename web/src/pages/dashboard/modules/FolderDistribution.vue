<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import * as echarts from 'echarts'
  import { getUserFolders } from '@/api/user/index'
  import { useTheme } from '@/composables/useTheme'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'FolderDistribution',
  })

  const { $t } = useTexts()

  const props = defineProps<{
    refreshKey?: number
  }>()

  interface Folder {
    id: string
    name: string
    file_count: number
    total_size: number
    total_size_formatted: string
  }

  const { appliedTheme, isDark } = useTheme()
  const folders = ref<Folder[]>([])
  const chartContainer = ref<HTMLElement | null>(null)
  const chartInstance = ref<echarts.ECharts | null>(null)
  const lastOption = ref<any>(null)
  const resizeTimeout = ref<any>(null)
  const isInitializing = ref(false)
  const resizeObserver = ref<ResizeObserver | null>(null)

  const sortedFolders = computed(() => {
    if (!folders.value || !Array.isArray(folders.value)) {
      return []
    }

    return [...folders.value]
      .filter((folder) => folder && typeof folder.name === 'string')
      .sort((a, b) => {
        if (a.total_size === 0 && b.total_size === 0) {
          return a.name.localeCompare(b.name)
        }
        return (b.total_size || 0) - (a.total_size || 0)
      })
  })

  const maxSize = computed(() => {
    if (!sortedFolders.value.length) {
      return 1
    }
    return Math.max(...sortedFolders.value.map((folder) => folder.total_size)) || 1
  })

  const getBrandColorFromCSS = (): { r: number; g: number; b: number } => {
    if (typeof window === 'undefined' || !document.documentElement) {
      return { r: 96, g: 165, b: 250 }
    }

    const style = getComputedStyle(document.documentElement)
    const brandRgb = style.getPropertyValue('--color-brand-500-rgb').trim()

    if (brandRgb) {
      const [r, g, b] = brandRgb.split(',').map((v) => parseInt(v.trim(), 10))
      if (!isNaN(r) && !isNaN(g) && !isNaN(b)) {
        return { r, g, b }
      }
    }

    return { r: 96, g: 165, b: 250 }
  }

  const rgbToHsl = (r: number, g: number, b: number): { h: number; s: number; l: number } => {
    r /= 255
    g /= 255
    b /= 255

    const max = Math.max(r, g, b)
    const min = Math.min(r, g, b)
    let h = 0
    let s = 0
    const l = (max + min) / 2

    if (max !== min) {
      const d = max - min
      s = l > 0.5 ? d / (2 - max - min) : d / (max + min)

      switch (max) {
        case r:
          h = ((g - b) / d + (g < b ? 6 : 0)) / 6
          break
        case g:
          h = ((b - r) / d + 2) / 6
          break
        case b:
          h = ((r - g) / d + 4) / 6
          break
      }
    }

    return { h: h * 360, s: s * 100, l: l * 100 }
  }

  const hslToRgb = (h: number, s: number, l: number): { r: number; g: number; b: number } => {
    h /= 360
    s /= 100
    l /= 100

    let r, g, b

    if (s === 0) {
      r = g = b = l
    } else {
      const hue2rgb = (p: number, q: number, t: number) => {
        if (t < 0) t += 1
        if (t > 1) t -= 1
        if (t < 1 / 6) return p + (q - p) * 6 * t
        if (t < 1 / 2) return q
        if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6
        return p
      }

      const q = l < 0.5 ? l * (1 + s) : l + s - l * s
      const p = 2 * l - q

      r = hue2rgb(p, q, h + 1 / 3)
      g = hue2rgb(p, q, h)
      b = hue2rgb(p, q, h - 1 / 3)
    }

    return {
      r: Math.round(r * 255),
      g: Math.round(g * 255),
      b: Math.round(b * 255),
    }
  }

  const generateDynamicColor = (index: number, ratio: number): string => {
    const brandColor = getBrandColorFromCSS()
    const hsl = rgbToHsl(brandColor.r, brandColor.g, brandColor.b)

    const hueShift = (index * 25) % 360
    const newHue = (hsl.h + hueShift) % 360

    const saturation = Math.max(40, Math.min(100, hsl.s + (ratio - 0.5) * 20))
    const lightness = isDark.value
      ? Math.max(35, Math.min(65, hsl.l + (ratio - 0.5) * 15))
      : Math.max(40, Math.min(70, hsl.l + (ratio - 0.5) * 10))

    const rgb = hslToRgb(newHue, saturation, lightness)

    const baseAlpha = 0.5 + ratio * 0.3
    const alpha = Math.max(0.5, Math.min(0.85, baseAlpha))

    return `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, ${alpha})`
  }

  const getTooltipColors = (): { bg: string; border: string } => {
    const brandColor = getBrandColorFromCSS()

    if (isDark.value) {
      return {
        bg: 'rgba(17, 24, 39, 0.95)',
        border: `rgba(${brandColor.r}, ${brandColor.g}, ${brandColor.b}, 0.5)`,
      }
    } else {
      return {
        bg: 'rgba(255, 255, 255, 0.98)',
        border: `rgba(${brandColor.r}, ${brandColor.g}, ${brandColor.b}, 0.5)`,
      }
    }
  }

  const getBrandColorRgbString = (): string => {
    const { r, g, b } = getBrandColorFromCSS()
    return `${r}, ${g}, ${b}`
  }

  const fetchFolders = async () => {
    try {
      const response = await getUserFolders({ with_stats: true })
      folders.value = response.data && response.data.folders ? response.data.folders : []
    } catch (_error) {
      folders.value = []
    }
  }

  const getCyberColor = (size: number, isEmpty: boolean, _name: string, index: number) => {
    if (isEmpty) {
      return isDark.value ? 'rgba(44, 43, 60, 0.4)' : 'rgba(229, 231, 235, 0.5)'
    }

    if (typeof size !== 'number' || size <= 0 || maxSize.value <= 0) {
      return isDark.value ? 'rgba(20, 25, 35, 0.5)' : 'rgba(243, 244, 246, 0.6)'
    }

    const ratio = size / maxSize.value
    return generateDynamicColor(index, ratio)
  }

  const prepareChartData = () => {
    if (!sortedFolders.value.length) {
      return []
    }

    const allEmpty = sortedFolders.value.every((folder) => !folder.total_size || folder.total_size <= 0)

    return sortedFolders.value.map((folder, index) => {
      let actualValue

      if (allEmpty) {
        actualValue = 100
      } else if (folder.total_size > 0) {
        actualValue = folder.total_size
      } else {
        const maxFolderSize = Math.max(...sortedFolders.value.map((f) => f.total_size))
        actualValue = Math.max(maxFolderSize * 0.01, 1)
      }

      const folderName = folder.name || $t('dashboard.folderDistribution.unnamed')
      const sizeFormatted = folder.total_size_formatted || '0 B'
      const fileCount = folder.file_count || 0
      const isEmpty = folder.total_size <= 0

      return {
        name: folderName,
        value: actualValue,
        sizeFormatted,
        fileCount,
        isEmpty,
        itemStyle: {
          color: getCyberColor(folder.total_size, isEmpty, folderName, index),
          borderColor: 'rgba(0, 130, 153, 0.25)',
          borderWidth: 1,
        },
        label: {
          show: true,
          color: isDark.value ? 'rgba(255, 255, 255, 0.95)' : 'rgba(31, 41, 55, 0.95)',
          fontSize: 12,
          fontWeight: 'bold',
        },
      }
    })
  }

  const getTotalSize = () => sortedFolders.value.reduce((total, folder) => total + (folder.total_size || 0), 0) || 1

  const initChart = async () => {
    if (!chartContainer.value || !sortedFolders.value.length) {
      return
    }
    if (isInitializing.value) {
      return
    }
    isInitializing.value = true

    try {
      await nextTick()

      const existing = echarts.getInstanceByDom(chartContainer.value)
      if (existing) {
        chartInstance.value = existing
        chartInstance.value.clear()
      } else {
        chartInstance.value = echarts.init(chartContainer.value, isDark.value ? 'dark' : null, { renderer: 'svg' })
      }

      const brandRgb = getBrandColorRgbString()
      const chartData = prepareChartData()
      const tooltipColors = getTooltipColors()

      const option = {
        backgroundColor: 'transparent',
        tooltip: {
          trigger: 'item',
          formatter(params: any) {
            const { data } = params
            if (!data) {
              return ''
            }

            const totalSize = getTotalSize()
            const percentage = totalSize > 0 ? ((data.value / totalSize) * 100).toFixed(1) : '0.0'

            const textColor = isDark.value ? '#ffffff' : '#1f2937'
            const subTextColor = isDark.value ? 'rgba(255, 255, 255, 0.85)' : 'rgba(107, 114, 128, 0.9)'

            return `
          <div style="padding: 10px 12px; background: ${tooltipColors.bg}; border-radius: var(--radius-sm); border: 2px solid ${tooltipColors.border}; backdrop-filter: blur(12px); box-shadow: 0 8px 16px -4px rgba(0, 0, 0, 0.3), 0 0 20px ${tooltipColors.border};">
            <div style="color: ${textColor}; font-weight: 600; margin-bottom: 6px; font-size: 13px;">${data.name || 'unknown'}</div>
            <div style="color: ${subTextColor}; font-size: 12px; line-height: 1.6;">${$t('dashboard.folderDistribution.tooltipSize', { size: data.sizeFormatted || '0 B' })}</div>
            <div style="color: ${subTextColor}; font-size: 12px; line-height: 1.6;">${$t('dashboard.folderDistribution.tooltipFiles', { count: data.fileCount || 0 })}</div>
            <div style="color: ${subTextColor}; font-size: 12px; line-height: 1.6;">${$t('dashboard.folderDistribution.tooltipPercent', { percent: percentage })}</div>
            </div>
          `
          },
          backgroundColor: 'transparent',
          borderColor: 'transparent',
          confine: true,
          appendToBody: false,
          textStyle: {
            color: '#ffffff',
          },
        },
        series: [
          {
            type: 'treemap',
            width: '100%',
            height: '100%',
            top: 0,
            left: 0,
            bottom: 0,
            right: 0,
            roam: false,
            nodeClick: false,
            breadcrumb: {
              show: false,
            },
            label: {
              show: true,
              position: 'inside',
              align: 'center',
              verticalAlign: 'middle',
              lineHeight: 16,
              fontSize: 12,
              fontWeight: 'bold',
              color: isDark.value ? 'rgba(255, 255, 255, 0.95)' : 'rgba(31, 41, 55, 0.95)',
              formatter(params: any) {
                const { data } = params
                if (!data) {
                  return ''
                }
                return `${data.name || ''}\n${data.sizeFormatted || ''} · ${$t('dashboard.folderDistribution.tooltipFiles', { count: data.fileCount || 0 })}`
              },
            },
            labelLayout(params: any) {
              const r = params.rect
              if (!r) {
                return {}
              }
              return {
                x: r.x + r.width / 2,
                y: r.y + r.height / 2,
                align: 'center',
                verticalAlign: 'middle',
              }
            },
            upperLabel: {
              show: false,
            },
            itemStyle: {
              borderColor: 'rgba(0, 130, 153, 0.25)',
              borderWidth: 1,
              gapWidth: 1,
            },
            emphasis: {
              focus: 'none',
              itemStyle: {
                borderColor: `rgba(${brandRgb}, 0.8)`,
                borderWidth: 2,
                shadowBlur: 10,
                shadowColor: `rgba(${brandRgb}, 0.3)`,
              },
              label: {
                color: 'rgba(255, 255, 255, 1)',
                fontSize: 13,
                fontWeight: 'bold',
              },
            },
            levels: [
              {
                itemStyle: {
                  borderColor: 'rgba(0, 130, 153, 0.25)',
                  borderWidth: 1,
                  gapWidth: 1,
                },
              },
            ],
            leaves: {
              label: {
                show: true,
                position: 'inside',
                align: 'center',
                verticalAlign: 'middle',
                lineHeight: 16,
                fontSize: 12,
                fontWeight: 'bold',
                color: isDark.value ? 'rgba(255, 255, 255, 0.95)' : 'rgba(31, 41, 55, 0.95)',
                formatter(params: any) {
                  const { data } = params
                  if (!data) {
                    return ''
                  }
                  return `${data.name || ''}\n${data.sizeFormatted || ''} · ${$t('dashboard.folderDistribution.tooltipFiles', { count: data.fileCount || 0 })}`
                },
              },
              upperLabel: { show: false },
            },
            data: chartData,
          },
        ],
      }

      chartInstance.value.setOption(option)
      lastOption.value = option

      chartInstance.value.off('finished')
      chartInstance.value.on('finished', () => {
        setTimeout(() => {
          addGridBackgrounds()
        }, 50)
      })

      setTimeout(() => {
        addGridBackgrounds()
      }, 200)
    } finally {
      isInitializing.value = false
    }
  }

  const addGridBackgrounds = () => {
    if (!chartContainer.value) {
      return
    }

    const treemapGroups = chartContainer.value.querySelectorAll('g[clip-path]')

    treemapGroups.forEach((group: Element) => {
      const rects = group.querySelectorAll('rect')

      rects.forEach((rect: Element) => {
        if (rect.hasAttribute('data-grid-added')) {
          return
        }

        const svgRect = rect as SVGRectElement
        const x = parseFloat(svgRect.getAttribute('x') || '0')
        const y = parseFloat(svgRect.getAttribute('y') || '0')
        const width = parseFloat(svgRect.getAttribute('width') || '0')
        const height = parseFloat(svgRect.getAttribute('height') || '0')

        if (width > 40 && height > 40) {
          const defs =
            group.ownerDocument?.querySelector('defs') ||
            group.ownerDocument?.createElementNS('http://www.w3.org/2000/svg', 'defs')

          if (!group.ownerDocument?.querySelector('defs')) {
            group.ownerDocument?.querySelector('svg')?.appendChild(defs!)
          }

          const brandRgb = getBrandColorRgbString()
          const patternId = `grid-pattern-${Math.random().toString(36).substr(2, 9)}`
          const pattern = group.ownerDocument?.createElementNS('http://www.w3.org/2000/svg', 'pattern')

          if (pattern && defs) {
            pattern.setAttribute('id', patternId)
            pattern.setAttribute('patternUnits', 'userSpaceOnUse')
            pattern.setAttribute('width', '20')
            pattern.setAttribute('height', '20')

            const gridLines = `
            <line x1="0" y1="20" x2="20" y2="20" stroke="rgba(${brandRgb}, 0.25)" stroke-width="1.2"/>
            <line x1="20" y1="0" x2="20" y2="20" stroke="rgba(${brandRgb}, 0.25)" stroke-width="1.2"/>
          `
            pattern.innerHTML = gridLines
            defs.appendChild(pattern)

            const gridOverlay = group.ownerDocument?.createElementNS('http://www.w3.org/2000/svg', 'rect')
            if (gridOverlay) {
              gridOverlay.setAttribute('x', x.toString())
              gridOverlay.setAttribute('y', y.toString())
              gridOverlay.setAttribute('width', width.toString())
              gridOverlay.setAttribute('height', height.toString())
              gridOverlay.setAttribute('fill', `url(#${patternId})`)
              gridOverlay.setAttribute('pointer-events', 'none')

              rect.parentNode?.insertBefore(gridOverlay, rect.nextSibling)
            }
          }
        }

        rect.setAttribute('data-grid-added', 'true')
      })
    })
  }

  const calculateChartHeight = () => {
    try {
      const viewportHeight = window.innerHeight || 900
      if (viewportHeight < 600) {
        return 250
      }
      if (viewportHeight < 900) {
        return 300
      }
      return 350
    } catch (_error) {
      return 300
    }
  }

  const handleResize = () => {
    clearTimeout(resizeTimeout.value)
    resizeTimeout.value = setTimeout(() => {
      if (chartInstance.value) {
        chartInstance.value.resize()
      }
    }, 200)
  }

  const setupContainerObserver = () => {
    if (!chartContainer.value || resizeObserver.value) {
      return
    }

    resizeObserver.value = new ResizeObserver(() => {
      clearTimeout(resizeTimeout.value)
      resizeTimeout.value = setTimeout(() => {
        if (chartInstance.value && !isInitializing.value) {
          requestAnimationFrame(() => {
            if (!chartInstance.value) return
            chartInstance.value!.resize()
            if (lastOption.value) {
              chartInstance.value!.setOption(lastOption.value, true)
            }
          })
        }
      }, 100)
    })

    resizeObserver.value.observe(chartContainer.value)
  }

  onMounted(async () => {
    await fetchFolders()
    await initChart()

    await nextTick()
    setupContainerObserver()

    window.addEventListener('resize', handleResize)
  })

  onUnmounted(() => {
    if (resizeObserver.value) {
      resizeObserver.value.disconnect()
      resizeObserver.value = null
    }

    if (chartInstance.value) {
      chartInstance.value.dispose()
      chartInstance.value = null
    }

    window.removeEventListener('resize', handleResize)

    if (resizeTimeout.value) {
      clearTimeout(resizeTimeout.value)
    }
  })

  watch(
    () => props.refreshKey,
    async () => {
      await fetchFolders()
      await nextTick()
      await initChart()
    }
  )

  watch(
    () => appliedTheme.value,
    async () => {
      await initChart()
    }
  )

  watch(
    () => folders.value,
    async (newVal) => {
      if (newVal && Array.isArray(newVal) && newVal.length > 0) {
        await nextTick()
        await initChart()
      }
    },
    { deep: true }
  )
</script>

<template>
  <div v-if="folders && folders.length > 0" class="folder-distribution cyber-card">
    <h3 class="section-title">
      <i class="fas fa-folder" />
      {{ $t('dashboard.folderDistribution.title') }}
      <span class="folder-count">({{ folders.length }})</span>
    </h3>

    <div ref="chartContainer" class="treemap-container" :style="{ height: `${calculateChartHeight()}px` }" />

    <div class="chart-info">
      <div class="info-text">
        <i class="fas fa-info-circle" />
        {{ $t('dashboard.folderDistribution.description') }}
      </div>
      <div class="folder-total">{{ $t('dashboard.folderDistribution.totalFolders', { count: folders.length }) }}</div>
    </div>
  </div>

  <div v-else class="folder-distribution cyber-card empty">
    <h3 class="section-title">
      <i class="fas fa-folder" />
      {{ $t('dashboard.folderDistribution.title') }}
    </h3>
    <div class="empty-state">
      <div class="empty-icon">
        <i class="fas fa-folder-open" />
      </div>
      <div class="empty-text">{{ $t('dashboard.folderDistribution.empty') }}</div>
    </div>
  </div>
</template>

<style scoped>
  .folder-distribution {
    padding: 20px;
  }

  .folder-distribution.empty .empty-state {
    padding: 40px 20px;
    text-align: center;
  }

  .section-title {
    margin: 0 0 16px 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 16px;
  }

  .folder-count {
    font-size: 14px;
    color: rgba(var(--color-content-rgb), 0.6);
    font-weight: 400;
  }

  .treemap-container {
    position: relative;
    overflow: hidden;
    background:
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.15) 2px, transparent 2px),
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.15) 2px, transparent 2px),
      rgba(var(--color-background-800-rgb), 0.3);
    background-size: 20px 20px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-border-default-rgb), 0.3);
    margin-bottom: 12px;
    width: 100%;
  }

  .chart-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
    color: rgba(var(--color-content-rgb), 0.6);
    padding-top: 8px;
    border-top: 1px solid rgba(var(--color-border-default-rgb), 0.2);
  }

  .info-text {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .info-text i {
    color: var(--color-brand-500);
    font-size: 10px;
  }

  .folder-total {
    font-weight: 500;
  }

  .empty-icon {
    font-size: 36px;
    margin-bottom: 12px;
    color: rgba(var(--color-content-rgb), 0.3);
  }

  .empty-text {
    font-size: 14px;
    color: rgba(var(--color-content-rgb), 0.5);
  }

  @media (max-width: 768px) {
    .folder-distribution {
      padding: 16px;
    }
  }
</style>
