<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { RouterLink } from 'vue-router'
  import { TIMING } from '@/constants/ui'
  import { useNotFoundCountdown } from './composables'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'NotFoundPage',
  })

  const { $t } = useTexts()

  const glitchTextRef = ref<HTMLElement | null>(null)
  let glitchTimer: number | null = null

  function restoreShadow(target: HTMLElement) {
    target.style.textShadow = '2px 2px 0px var(--color-error-500), -2px -2px 0px var(--color-brand-500)'
  }

  function applyRandomShadow(target: HTMLElement) {
    const firstX = Math.random() * 10 - 5
    const firstY = Math.random() * 10 - 5
    const secondX = Math.random() * 10 - 5
    const secondY = Math.random() * 10 - 5
    target.style.textShadow =
      `${firstX}px ${firstY}px rgba(var(--color-error-rgb), 0.7), ` +
      `${secondX}px ${secondY}px rgba(var(--color-brand-500-rgb), 0.7)`
    window.setTimeout(() => {
      restoreShadow(target)
    }, 100)
  }

  function clearGlitch() {
    if (glitchTimer) {
      clearInterval(glitchTimer)
      glitchTimer = null
    }
  }

  function startGlitch(target: HTMLElement) {
    clearGlitch()
    applyRandomShadow(target)
    glitchTimer = window.setInterval(() => {
      applyRandomShadow(target)
    }, TIMING.TOAST.DEFAULT_DURATION)
  }

  const { countdown, progress, start, clear } = useNotFoundCountdown({ seconds: 7 })

  const _progressStyle = computed(() => {
    const value = Math.min(Math.max(progress.value, 0), 100)
    return {
      width: `${value}%`,
    }
  })

  onMounted(() => {
    start()
    const target = glitchTextRef.value
    if (target) {
      startGlitch(target)
    }
  })

  onUnmounted(() => {
    clearGlitch()
    clear()
  })
</script>

<template>
  <div class="not-found-page">
    <div class="atmospheric-layer">
      <div class="fog-effect" />
      <div class="light-beam beam-1" />
      <div class="light-beam beam-2" />
      <div class="floating-elements">
        <div class="hex-element hex-1" />
        <div class="hex-element hex-2" />
        <div class="hex-element hex-3" />
      </div>
    </div>
    <div class="page-content">
      <div class="glitch-container">
        <h1 class="error-code">
          <span ref="glitchTextRef" class="glitch-text" data-text="404">404</span>
        </h1>
        <div class="error-message">
          <h2 class="error-title">{{ $t('notFound.title') }}</h2>
          <p class="error-description">{{ $t('notFound.description') }}</p>
          <div class="countdown-container">
            <div class="countdown-terminal">
              <div class="terminal-header">
                <div class="terminal-dots">
                  <span class="terminal-dot red" />
                  <span class="terminal-dot yellow" />
                  <span class="terminal-dot green" />
                </div>
                <span class="terminal-title">{{ $t('notFound.terminal.title') }}</span>
                <div class="terminal-status">
                  <span class="status-indicator active" />
                  <span class="terminal-id">TERM_0x404</span>
                </div>
              </div>
              <div class="terminal-content">
                <p class="terminal-line">
                  <span class="terminal-prompt">$</span>
                  <span class="terminal-text">{{ $t('notFound.terminal.connecting') }}</span>
                </p>
                <p class="terminal-line">
                  <span class="terminal-prompt">$</span>
                  <span class="terminal-text">{{ $t('notFound.terminal.countdown', { seconds: countdown }) }}</span>
                </p>
              </div>
            </div>
          </div>
          <div class="action-buttons">
            <RouterLink to="/" class="cyber-capsule-btn">
              <div class="btn-inner">
                <div class="btn-icon">
                  <i class="fas fa-rocket" />
                </div>
                <span class="btn-text">{{ $t('notFound.action') }}</span>
              </div>
              <div class="btn-glitch-top" />
              <div class="btn-glitch-bottom" />
            </RouterLink>
          </div>
        </div>
      </div>
      <div class="corner-decoration top-left" />
      <div class="corner-decoration top-right" />
      <div class="corner-decoration bottom-left" />
      <div class="corner-decoration bottom-right" />
      <div class="corner-line top-line" />
      <div class="corner-line bottom-line" />
      <div class="corner-line left-line" />
      <div class="corner-line right-line" />
      <div class="robot-container">
        <div class="robot-head">
          <div class="robot-eye left" />
          <div class="robot-eye right" />
          <div class="robot-mouth" />
        </div>
        <div class="robot-body">
          <div class="robot-core" />
        </div>
        <div class="robot-arm left" />
        <div class="robot-arm right" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./index.scss"></style>
