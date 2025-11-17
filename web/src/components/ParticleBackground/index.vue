<script setup lang="ts">
  import { onMounted, onUnmounted } from 'vue'
  import type { ParticleBackgroundProps } from './types'

  defineOptions({
    name: 'ParticleBackground',
  })

  const props = withDefaults(defineProps<ParticleBackgroundProps>(), {
    maxParticles: 100,
    colors: () => ['var(--color-brand-500)', 'var(--color-brand-600)', 'var(--color-warning-500)', 'var(--color-success-500)'],
    theme: 'default',
    connectionDistance: 100,
    connectionOpacity: 0.15,
    showConnections: true,
    particleSpeed: 1,
    showGrid: true,
  })

  /* 粒子类型 */
  interface Particle {
    x: number
    y: number
    size: number
    speedX: number
    speedY: number
    color: string
    alpha: number
    direction: number
  }

  function getCSSVariableRGB(variableName: string): string {
    if (typeof window === 'undefined') return '14, 165, 233'
    const value = getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
    return value || '14, 165, 233'
  }

  function getCSSVariableColor(variableName: string): string {
    if (typeof window === 'undefined') return 'rgb(14, 165, 233)'
    const value = getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
    return value || 'rgb(14, 165, 233)'
  }

  const particlesArray: Particle[] = []
  let animationId: number | null = null
  let canvas: HTMLCanvasElement | null = null
  let ctx: CanvasRenderingContext2D | null = null
  let containerWidth = 0
  let containerHeight = 0
  let resolvedColors: string[] = []
  let brandColor = ''

  function createParticle(): Particle {
    return {
      x: Math.random() * containerWidth,
      y: Math.random() * containerHeight,
      size: Math.random() * 4 + 1,
      speedX: (Math.random() - 0.5) * 0.5 * props.particleSpeed,
      speedY: (Math.random() - 0.5) * 0.5 * props.particleSpeed,
      color: resolvedColors[Math.floor(Math.random() * resolvedColors.length)],
      alpha: Math.random() * 0.6 + 0.2,
      direction: Math.random() * Math.PI * 2,
    }
  }

  function initParticles() {
    for (let i = 0; i < props.maxParticles; i++) {
      particlesArray.push(createParticle())
    }
  }

  function updateParticles() {
    for (let i = 0; i < particlesArray.length; i++) {
      const p = particlesArray[i]
      p.x += p.speedX
      p.y += p.speedY
      if (p.x < 0) p.x = containerWidth
      if (p.x > containerWidth) p.x = 0
      if (p.y < 0) p.y = containerHeight
      if (p.y > containerHeight) p.y = 0

      if (Math.random() < 0.01) {
        p.direction += (Math.random() - 0.5) * 0.2
        p.speedX = Math.cos(p.direction) * (Math.random() * 0.5 + 0.1) * props.particleSpeed
        p.speedY = Math.sin(p.direction) * (Math.random() * 0.5 + 0.1) * props.particleSpeed
      }
    }
  }

  function drawParticles() {
    if (!ctx) {
      return
    }

    ctx.clearRect(0, 0, containerWidth, containerHeight)

    for (let i = 0; i < particlesArray.length; i++) {
      const p = particlesArray[i]
      ctx.beginPath()
      ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2)
      ctx.fillStyle = p.color
      ctx.globalAlpha = p.alpha
      ctx.fill()
    }

    if (props.showConnections) {
      ctx.globalAlpha = props.theme === 'docs' ? props.connectionOpacity * 0.5 : props.connectionOpacity
      ctx.strokeStyle = brandColor
      ctx.lineWidth = 0.5

      for (let i = 0; i < particlesArray.length; i++) {
        for (let j = i + 1; j < particlesArray.length; j++) {
          const dx = particlesArray[i].x - particlesArray[j].x
          const dy = particlesArray[i].y - particlesArray[j].y
          const distance = Math.sqrt(dx * dx + dy * dy)

          if (distance < props.connectionDistance) {
            ctx.beginPath()
            ctx.moveTo(particlesArray[i].x, particlesArray[i].y)
            ctx.lineTo(particlesArray[j].x, particlesArray[j].y)
            ctx.stroke()
          }
        }
      }
    }

    const gradient = ctx.createRadialGradient(
      containerWidth / 2,
      containerHeight / 2,
      0,
      containerWidth / 2,
      containerHeight / 2,
      containerWidth / 1.5
    )

    const cyberBlueRGB = getCSSVariableRGB('--color-brand-500-rgb')
    const cyberPinkRGB = getCSSVariableRGB('--color-brand-600-rgb')

    gradient.addColorStop(0, `rgba(${cyberBlueRGB}, 0)`)
    gradient.addColorStop(0.5, props.theme === 'docs' ? `rgba(${cyberBlueRGB}, 0.01)` : `rgba(${cyberBlueRGB}, 0.03)`)
    gradient.addColorStop(1, props.theme === 'docs' ? `rgba(${cyberPinkRGB}, 0.02)` : `rgba(${cyberPinkRGB}, 0.05)`)

    ctx.globalAlpha = props.theme === 'docs' ? 0.15 : 0.3
    ctx.fillStyle = gradient
    ctx.fillRect(0, 0, containerWidth, containerHeight)

    if (props.showGrid) {
      ctx.globalAlpha = props.theme === 'docs' ? 0.03 : 0.08
      ctx.strokeStyle = brandColor
      ctx.lineWidth = 0.2

      const gridSize = 40

      for (let x = 0; x <= containerWidth; x += gridSize) {
        ctx.beginPath()
        ctx.moveTo(x, 0)
        ctx.lineTo(x, containerHeight)
        ctx.stroke()
      }

      for (let y = 0; y <= containerHeight; y += gridSize) {
        ctx.beginPath()
        ctx.moveTo(0, y)
        ctx.lineTo(containerWidth, y)
        ctx.stroke()
      }
    }
  }

  function animate() {
    updateParticles()
    drawParticles()
    animationId = requestAnimationFrame(animate)
  }

  function handleResize() {
    const container = document.getElementById('particle-background')
    if (!container || !canvas) {
      return
    }

    containerWidth = container.offsetWidth
    containerHeight = container.offsetHeight

    canvas.width = containerWidth
    canvas.height = containerHeight
  }

  function initCanvas() {
    const container = document.getElementById('particle-background')
    if (!container) {
      return
    }

    resolvedColors = props.colors.map((color) => {
      if (color.startsWith('var(')) {
        const varName = color.match(/var\((--[^)]+)\)/)?.[1]
        return varName ? getCSSVariableColor(varName) : color
      }
      return color
    })

    brandColor = getCSSVariableColor('--color-brand-500')

    canvas = document.createElement('canvas')
    ctx = canvas.getContext('2d')

    if (!ctx) {
      return
    }

    containerWidth = container.offsetWidth
    containerHeight = container.offsetHeight
    canvas.width = containerWidth
    canvas.height = containerHeight

    container.appendChild(canvas)

    initParticles()
    animate()

    window.addEventListener('resize', handleResize)
  }

  onMounted(() => {
    setTimeout(initCanvas, 100)
  })

  onUnmounted(() => {
    if (animationId !== null) {
      cancelAnimationFrame(animationId)
    }
    window.removeEventListener('resize', handleResize)

    particlesArray.length = 0
  })
</script>

<template>
  <div class="particle-background" :class="{ 'docs-theme': theme === 'docs' }">
    <div id="particle-background" class="particle-container" />
  </div>
</template>

<style scoped>
  .particle-background {
    @apply pointer-events-none fixed left-0 top-0 z-0 h-full w-full overflow-hidden;
  }

  .particle-container {
    @apply h-full w-full;
    background: linear-gradient(135deg, var(--color-background-darker) 0%, var(--color-background-dark) 100%);
  }

  .docs-theme .particle-container {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.98) 50%,
      rgba(var(--color-background-900-rgb), 0.95) 100%
    );
  }

  .particle-container::after {
    @apply pointer-events-none absolute left-0 top-0 h-full w-full;
    content: '';
    background: radial-gradient(
      circle at center,
      rgba(var(--color-brand-500-rgb), 0.1) 0%,
      rgba(var(--color-brand-600-rgb), 0.05) 50%,
      transparent 70%
    );
  }

  .docs-theme .particle-container::after {
    background: radial-gradient(
      circle at center,
      rgba(var(--color-brand-500-rgb), 0.05) 0%,
      rgba(var(--color-brand-600-rgb), 0.02) 50%,
      transparent 70%
    );
  }
</style>
