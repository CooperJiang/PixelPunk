<script setup lang="ts">
  import { computed, nextTick, onMounted, ref } from 'vue'
  import { useSettingsStore } from '@/store/settings'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'HeroTitle',
  })

  const settingsStore = useSettingsStore()
  const { $t } = useTexts()

  const DEFAULT_HERO_TITLE = '让图片管理从繁琐到简单，让文件分享从等待到极速'
  const DEFAULT_FEATURES_TEXT = 'AI自动识别 · 智能分类整理 · 秒传极速分享 · 开源社区驱动'

  const dynamicSiteName = computed(() => {
    const siteName = settingsStore.siteName
    if (siteName) {
      const words = siteName.split(' ').filter((word) => word.trim())
      if (words.length >= 2) {
        return {
          firstWord: words[0],
          otherWords: words.slice(1),
        }
      }
      return {
        firstWord: siteName,
        otherWords: [],
      }
    }
    return {
      firstWord: 'PIXEL',
      otherWords: ['PUNK'],
    }
  })

  const dynamicHeroTitle = computed(() => {
    const configValue = settingsStore.siteHeroTitle?.trim()
    if (!configValue) {
      return ''
    }
    if (configValue === DEFAULT_HERO_TITLE) {
      return $t('home.hero.defaultTitle')
    }
    return configValue
  })

  const dynamicFeaturesText = computed(() => {
    const configValue = settingsStore.siteFeaturesText?.trim()
    if (!configValue) {
      return ''
    }
    if (configValue === DEFAULT_FEATURES_TEXT) {
      return $t('home.hero.defaultFeatures')
    }
    return configValue
  })

  const glitchTitle = ref<HTMLElement>()

  const initTitleAnimations = async () => {
    if (typeof window === 'undefined') {
      return
    }

    try {
      const { default: gsap } = await import('gsap')

      if (!glitchTitle.value) {
        return
      }

      const tl = gsap.timeline()

      tl.to(glitchTitle.value, {
        opacity: 1,
        y: 0,
        scale: 1,
        duration: 0.6,
        ease: 'back.out(1.2)',
        delay: 0.1,
      })
        .to(
          glitchTitle.value,
          {
            x: 1,
            duration: 0.02,
            repeat: 3,
            yoyo: true,
            ease: 'none',
          },
          '-=0.1'
        )
        .set(glitchTitle.value, { x: 0 })
        .to(
          glitchTitle.value,
          {
            textShadow: '0 0 20px rgba(var(--color-brand-500-rgb), 0.6), 0 0 40px rgba(var(--color-brand-500-rgb), 0.3)',
            duration: 3,
            repeat: -1,
            yoyo: true,
            ease: 'sine.inOut',
          },
          '-=0.2'
        )
    } catch (_error) {
      if (glitchTitle.value) {
        glitchTitle.value.classList.add('css-fallback')
      }
    }
  }

  onMounted(() => {
    nextTick(() => {
      initTitleAnimations()
    })
  })
</script>

<template>
  <div class="hero-title-wrapper">
    <div class="cyber-glitch-wrapper mb-6">
      <h1 ref="glitchTitle" class="cyber-glitch-title">
        <span class="cyber-glitch-word cyber-glitch-color-1" :data-text="dynamicSiteName.firstWord">{{
          dynamicSiteName.firstWord
        }}</span>
        <template v-for="(word, index) in dynamicSiteName.otherWords" :key="index">
          <span class="cyber-glitch-separator" />
          <span :class="`cyber-glitch-word cyber-glitch-color-${((index + 1) % 5) + 1}`" :data-text="word">{{ word }}</span>
        </template>
      </h1>
    </div>
    <p v-if="dynamicHeroTitle" class="subtitle mx-auto mb-6 max-w-5xl text-2xl font-light tracking-wider md:text-3xl">
      {{ dynamicHeroTitle }}
    </p>
    <p
      v-if="dynamicFeaturesText"
      class="description mx-auto mb-8 max-w-3xl text-base leading-loose text-content-muted md:text-lg"
      v-html="dynamicFeaturesText"
    />
  </div>
</template>

<style scoped>
  .cyber-glitch-wrapper {
    position: relative;
  }

  .cyber-glitch-title {
    @apply text-5xl font-bold uppercase tracking-wider md:text-6xl lg:text-7xl;
    font-size: 4rem;
    font-weight: 900;
    position: relative;
    letter-spacing: 0.08em;
    margin-bottom: 0.5rem;
    display: inline-flex;
    align-items: baseline;

    opacity: 0;
    transform: translateY(30px) scale(0.9);
  }

  .cyber-glitch-separator {
    margin: 0 0.1em;
  }

  .cyber-glitch-word {
    position: relative;
    animation:
      cyber-glitch-main 8s infinite,
      cyber-glow-pulse 4s ease-in-out infinite alternate;
    will-change: transform, filter, text-shadow;
    transform: translateZ(0);
  }

  .cyber-glitch-color-1 {
    color: var(--color-brand-500);
    text-shadow:
      0 0 15px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 45px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-glitch-color-2 {
    color: var(--color-error-500);
    text-shadow:
      0 0 15px rgba(var(--color-error-rgb), 0.6),
      0 0 30px rgba(var(--color-error-rgb), 0.4),
      0 0 45px rgba(var(--color-error-rgb), 0.2);
  }

  .cyber-glitch-color-3 {
    color: var(--color-success-500);
    text-shadow:
      0 0 15px rgba(0, 255, 136, 0.6),
      0 0 30px rgba(0, 255, 136, 0.4),
      0 0 45px rgba(0, 255, 136, 0.2);
  }

  .cyber-glitch-color-4 {
    color: var(--color-purple-500);
    text-shadow:
      0 0 15px rgba(var(--color-error-rgb), 0.6),
      0 0 30px rgba(var(--color-error-rgb), 0.4),
      0 0 45px rgba(var(--color-error-rgb), 0.2);
  }

  .cyber-glitch-color-5 {
    color: var(--color-warning-500);
    text-shadow:
      0 0 15px rgba(var(--color-warning-rgb), 0.6),
      0 0 30px rgba(var(--color-warning-rgb), 0.4),
      0 0 45px rgba(var(--color-warning-rgb), 0.2);
  }

  @keyframes cyber-glitch-main {
    0%,
    90%,
    92%,
    94%,
    96%,
    98%,
    100% {
      transform: translateY(0) skewX(0deg);
      filter: hue-rotate(0deg) saturate(1);
      text-shadow:
        0 0 15px rgba(var(--color-brand-500-rgb), 0.6),
        0 0 30px rgba(var(--color-brand-500-rgb), 0.4),
        0 0 45px rgba(var(--color-brand-500-rgb), 0.2);
    }

    91% {
      transform: translateY(-0.5px) skewX(2deg);
      filter: hue-rotate(90deg) saturate(1.5);
      text-shadow:
        0 0 20px rgba(var(--color-error-rgb), 0.8),
        0 0 40px rgba(var(--color-error-rgb), 0.6),
        2px 0 0 rgba(var(--color-brand-500-rgb), 1),
        -2px 0 0 rgba(var(--color-error-rgb), 1);
    }

    93% {
      transform: translateY(0.3px) skewX(-1deg);
      filter: hue-rotate(-60deg) saturate(2);
      text-shadow:
        0 0 25px rgba(var(--color-error-rgb), 0.9),
        0 0 50px rgba(var(--color-error-rgb), 0.5),
        -1px 0 0 rgba(var(--color-brand-500-rgb), 1),
        1px 0 0 rgba(var(--color-error-rgb), 1);
    }

    95% {
      transform: translateY(-0.3px) skewX(1.5deg);
      filter: hue-rotate(120deg) saturate(1.2);
      text-shadow:
        0 0 30px rgba(var(--color-brand-500-rgb), 1),
        0 0 60px rgba(var(--color-brand-500-rgb), 0.7),
        3px 0 0 rgba(var(--color-error-rgb), 0.8),
        -3px 0 0 rgba(var(--color-brand-500-rgb), 0.8);
    }

    97% {
      transform: translateY(0.4px) skewX(-2deg);
      filter: hue-rotate(180deg) saturate(1.8);
      text-shadow:
        0 0 35px rgba(var(--color-error-rgb), 0.9),
        0 0 70px rgba(var(--color-error-rgb), 0.6),
        -2px 0 0 rgba(var(--color-error-rgb), 1),
        2px 0 0 rgba(var(--color-brand-500-rgb), 1);
    }

    99% {
      transform: translateY(-0.2px) skewX(3deg);
      filter: hue-rotate(240deg) saturate(2.5);
      text-shadow:
        0 0 40px rgba(var(--color-error-rgb), 1),
        0 0 80px rgba(var(--color-error-rgb), 0.8),
        4px 0 0 rgba(var(--color-brand-500-rgb), 1),
        -4px 0 0 rgba(var(--color-error-rgb), 1);
    }
  }

  @keyframes cyber-glow-pulse {
    0% {
      filter: brightness(1) contrast(1);
    }
    100% {
      filter: brightness(1.2) contrast(1.3);
    }
  }

  .cyber-glitch-title.css-fallback {
    animation: titleEntrance 0.8s cubic-bezier(0.25, 0.46, 0.45, 0.94) 0.1s both;
  }

  @keyframes titleEntrance {
    0% {
      opacity: 0;
      transform: translateY(30px) scale(0.9);
    }
    100% {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }

  .cyber-glitch-word::before,
  .cyber-glitch-word::after {
    content: attr(data-text);
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
  }

  .cyber-glitch-word::after {
    left: -2px;
    text-shadow:
      2px 0 var(--color-brand-500),
      0 0 10px rgba(var(--color-brand-500-rgb), 0.8),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.5);
    clip: rect(85px, 550px, 140px, 0);
    animation:
      glitch-anim2-enhanced 4.5s infinite linear,
      glitch-color-shift2 6s infinite ease-in-out;
    mix-blend-mode: screen;
  }

  @keyframes glitch-anim-enhanced {
    0% {
      clip: rect(4px, 9999px, 12px, 0);
      transform: translateX(0px);
    }
    5% {
      clip: rect(68px, 9999px, 75px, 0);
      transform: translateX(-2px);
    }
    10% {
      clip: rect(24px, 9999px, 30px, 0);
      transform: translateX(3px);
    }
    15% {
      clip: rect(85px, 9999px, 95px, 0);
      transform: translateX(-1px);
    }
    20% {
      clip: rect(12px, 9999px, 25px, 0);
      transform: translateX(2px);
    }
    25% {
      clip: rect(50px, 9999px, 65px, 0);
      transform: translateX(-3px);
    }
    30% {
      clip: rect(90px, 9999px, 105px, 0);
      transform: translateX(1px);
    }
    35% {
      clip: rect(35px, 9999px, 45px, 0);
      transform: translateX(-2px);
    }
    40% {
      clip: rect(70px, 9999px, 85px, 0);
      transform: translateX(4px);
    }
    45% {
      clip: rect(15px, 9999px, 28px, 0);
      transform: translateX(-1px);
    }
    50% {
      clip: rect(55px, 9999px, 72px, 0);
      transform: translateX(2px);
    }
    55% {
      clip: rect(95px, 9999px, 108px, 0);
      transform: translateX(-3px);
    }
    60% {
      clip: rect(8px, 9999px, 18px, 0);
      transform: translateX(1px);
    }
    65% {
      clip: rect(45px, 9999px, 58px, 0);
      transform: translateX(-2px);
    }
    70% {
      clip: rect(80px, 9999px, 92px, 0);
      transform: translateX(3px);
    }
    75% {
      clip: rect(25px, 9999px, 38px, 0);
      transform: translateX(-1px);
    }
    80% {
      clip: rect(62px, 9999px, 78px, 0);
      transform: translateX(2px);
    }
    85% {
      clip: rect(18px, 9999px, 32px, 0);
      transform: translateX(-4px);
    }
    90% {
      clip: rect(88px, 9999px, 102px, 0);
      transform: translateX(1px);
    }
    95% {
      clip: rect(40px, 9999px, 55px, 0);
      transform: translateX(-2px);
    }
    100% {
      clip: rect(75px, 9999px, 88px, 0);
      transform: translateX(0px);
    }
  }

  @keyframes glitch-anim2-enhanced {
    0% {
      clip: rect(73px, 9999px, 85px, 0);
      transform: translateX(0px);
    }
    15% {
      clip: rect(25px, 9999px, 35px, 0);
      transform: translateX(3px);
    }
    30% {
      clip: rect(54px, 9999px, 68px, 0);
      transform: translateX(-2px);
    }
    45% {
      clip: rect(97px, 9999px, 110px, 0);
      transform: translateX(1px);
    }
    60% {
      clip: rect(21px, 9999px, 32px, 0);
      transform: translateX(-4px);
    }
    75% {
      clip: rect(89px, 9999px, 98px, 0);
      transform: translateX(2px);
    }
    90% {
      clip: rect(17px, 9999px, 28px, 0);
      transform: translateX(-3px);
    }
    100% {
      clip: rect(73px, 9999px, 85px, 0);
      transform: translateX(0px);
    }
  }

  @keyframes glitch-color-shift {
    0%,
    100% {
      filter: hue-rotate(0deg) brightness(1);
    }
    25% {
      filter: hue-rotate(90deg) brightness(1.2);
    }
    50% {
      filter: hue-rotate(180deg) brightness(0.8);
    }
    75% {
      filter: hue-rotate(270deg) brightness(1.1);
    }
  }

  @keyframes glitch-color-shift2 {
    0%,
    100% {
      filter: hue-rotate(0deg) brightness(1);
    }
    20% {
      filter: hue-rotate(60deg) brightness(1.3);
    }
    40% {
      filter: hue-rotate(120deg) brightness(0.9);
    }
    60% {
      filter: hue-rotate(200deg) brightness(1.1);
    }
    80% {
      filter: hue-rotate(300deg) brightness(1.2);
    }
  }

  .text-brand-500 {
    color: var(--color-brand-500);
  }

  .text-error-500 {
    color: var(--color-error-500);
  }

  .subtitle {
    color: var(--color-content-heading);
    font-weight: 400;

    text-shadow:
      0 0 10px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.2),
      0 1px 2px rgba(0, 0, 0, 0.2);

    opacity: 1;
    visibility: visible;

    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-content-heading) 50%, var(--color-brand-500) 100%);

    @supports (-webkit-background-clip: text) {
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    @supports not (-webkit-background-clip: text) {
      background: none;
      color: var(--color-content-heading);
    }
  }

  .description {
    color: var(--color-content-muted);
    font-weight: 400;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  }

  @media (max-width: 768px) {
    .cyber-glitch-title {
      font-size: 2.2rem;
    }

    .cyber-glitch-pixel::before,
    .cyber-glitch-pixel::after,
    .cyber-glitch-punk::before,
    .cyber-glitch-punk::after {
      clip: rect(0, 350px, 200px, 0);
    }

    .subtitle {
      font-size: 1.5rem;
    }

    .description {
      font-size: 0.95rem;
      max-width: 90%;
    }
  }

  @media (max-width: 480px) {
    .cyber-glitch-title {
      font-size: 2rem;
    }
  }

  @media (max-width: 360px) {
    .cyber-glitch-title {
      font-size: 1.8rem;
      letter-spacing: 0.02em;
    }
  }
</style>
