import { createApp, h, type Directive, type DirectiveBinding } from 'vue'
import { useTexts } from '@/composables/useTexts'

/* 存储每个元素的loading实例 */
const loadingInstances = new Map<HTMLElement, any>()

/* 获取翻译函数（延迟调用，避免 Pinia 未初始化错误） */
const getTranslation = (key: string) => {
  try {
    const { $t } = useTexts()
    return $t(key)
  } catch {
    return 'Loading...' // 降级处理
  }
}

const LoadingComponent = {
  name: 'VLoadingComponent',
  props: {
    text: {
      type: String,
      default: '',
    },
    background: {
      type: String,
      default: 'rgba(var(--color-background-800-rgb), 0.85)',
    },
    textColor: {
      type: String,
      default: 'var(--color-cyber-blue)',
    },
    type: {
      type: String,
      default: 'default', // default, progress, pulse
      validator: (value: string) => ['default', 'progress', 'pulse'].includes(value),
    },
    progress: {
      type: Number,
      default: 0, // 0-100
      validator: (value: number) => value >= 0 && value <= 100,
    },
    loading: {
      type: Boolean,
      default: true,
    },
  },
  computed: {
    displayText() {
      return this.text || getTranslation('status.loading')
    },
  },
  data() {
    return {
      animatedProgress: this.type === 'progress' ? Math.max(5, this.progress) : Math.floor(Math.random() * 100),
      progressTimer: null,
    }
  },
  mounted() {
    if (this.type === 'progress') {
      this.animateProgress()
    }
  },
  watch: {
    loading(newValue) {
      if (this.type === 'progress' && !newValue) {
        this.completeProgress()
      }
    },
  },
  beforeUnmount() {
    if (this.progressTimer) {
      clearTimeout(this.progressTimer)
    }
  },
  methods: {
    animateProgress() {
      if (!this.loading) return

      const increment = () => {
        if (this.loading && this.animatedProgress < 85) {
          const rate = this.animatedProgress > 70 ? 0.3 : 2
          this.animatedProgress = Math.min(85, this.animatedProgress + Math.random() * rate)
          this.progressTimer = setTimeout(increment, Math.random() * 400 + 300)
        }
      }
      increment()
    },
    completeProgress() {
      if (this.progressTimer) {
        clearTimeout(this.progressTimer)
      }
      const complete = () => {
        if (this.animatedProgress < 100) {
          this.animatedProgress = Math.min(100, this.animatedProgress + 8)
          setTimeout(complete, 50)
        }
      }
      complete()
    },
  },
  render() {
    const renderText = (size = null, forceHide = false) => {
      if (!this.displayText || forceHide) return null

      const textSize = size?.fontSize || '14px'
      const shouldShow = size?.showText !== false

      if (!shouldShow) return null

      return h(
        'div',
        {
          class: 'loading-text',
          style: {
            fontSize: textSize,
            fontWeight: '500',
            color: 'var(--color-brand-400)',
            textAlign: 'center',
            maxWidth: '100%',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
            marginTop: size?.containerSize === 'tiny' ? '1px' : '2px',
            letterSpacing: '0.3px',
            textShadow: '0 0 8px rgba(var(--color-brand-500-rgb), 0.6)',
            opacity: '0.9',
            animation: 'textFade 1.5s ease-in-out infinite',
          },
        },
        this.displayText
      )
    }

    const getResponsiveSize = () => {
      const containerRect = this.$el?.parentElement?.getBoundingClientRect() || {
        width: 300,
        height: 200,
      }
      const containerWidth = Math.min(containerRect.width, containerRect.height)

      let containerSize: string
      let iconSize: number
      let fontSize: string
      let showText: boolean

      if (containerWidth < 80) {
        containerSize = 'tiny'
        iconSize = Math.max(30, containerWidth * 0.6)
        fontSize = '10px'
        showText = false
      } else if (containerWidth < 120) {
        containerSize = 'small'
        iconSize = Math.max(32, containerWidth * 0.4)
        fontSize = '10px'
        showText = containerWidth > 100
      } else if (containerWidth < 200) {
        containerSize = 'medium'
        iconSize = Math.max(45, containerWidth * 0.35)
        fontSize = '12px'
        showText = true
      } else {
        containerSize = 'normal'
        iconSize = Math.max(60, Math.min(90, containerWidth * 0.35))
        fontSize = '14px'
        showText = true
      }

      return {
        containerSize,
        iconSize,
        fontSize,
        showText,
        animationSize: `${iconSize}px`,
        hexagonSize: `${Math.floor(iconSize * 0.6)}px`,
        gap: containerSize === 'tiny' ? '0.125rem' : containerSize === 'small' ? '0.25rem' : '0.375rem',
      }
    }

    const renderDefaultType = () => {
      const size = getResponsiveSize()

      return h(
        'div',
        {
          class: 'loading-container',
          style: {
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            gap: size.gap,
            padding: '0.25rem',
          },
        },
        [
          h(
            'div',
            {
              class: 'simple-loading-core',
              style: {
                width: size.animationSize,
                height: size.animationSize,
                position: 'relative',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              },
            },
            [
              h('div', {
                class: 'spinner-ring',
                style: {
                  width: '100%',
                  height: '100%',
                  border: `${size.containerSize === 'tiny' ? '2px' : '3px'} solid rgba(var(--color-brand-500-rgb), 0.15)`,
                  borderTopColor: 'var(--color-brand-500)',
                  borderRadius: '50%',
                  position: 'absolute',
                  animation: 'smoothRotate 1.2s cubic-bezier(0.4, 0.0, 0.2, 1) infinite',
                },
              }),

              h('div', {
                class: 'spinner-ring-middle',
                style: {
                  width: '60%',
                  height: '60%',
                  border: `${size.containerSize === 'tiny' ? '2px' : '2px'} solid rgba(var(--color-brand-400-rgb), 0.2)`,
                  borderRightColor: 'var(--color-brand-400)',
                  borderRadius: '50%',
                  position: 'absolute',
                  animation: 'smoothRotate 1.8s cubic-bezier(0.4, 0.0, 0.2, 1) infinite reverse',
                  opacity: 0.7,
                },
              }),

              h('div', {
                class: 'spinner-center',
                style: {
                  position: 'absolute',
                  width: size.containerSize === 'tiny' ? '6px' : '8px',
                  height: size.containerSize === 'tiny' ? '6px' : '8px',
                  background: 'var(--color-brand-500)',
                  borderRadius: '50%',
                  boxShadow: '0 0 12px rgba(var(--color-brand-500-rgb), 0.6)',
                  animation: 'centerGlow 2s ease-in-out infinite',
                },
              }),
            ]
          ),
          renderText(size),
        ]
      )
    }

    const renderProgressType = () => {
      const size = getResponsiveSize()
      const progressHeight = size.containerSize === 'tiny' ? '8px' : size.containerSize === 'small' ? '10px' : '12px'
      const progressWidth = size.containerSize === 'tiny' ? '120px' : size.containerSize === 'small' ? '150px' : '250px'

      return h(
        'div',
        {
          class: 'loading-container',
          style: {
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            gap: size.gap,
            maxWidth: '100%',
            padding: '0.25rem',
          },
        },
        [
          h(
            'div',
            {
              class: 'cyber-progress-container',
              style: {
                width: '100%',
                maxWidth: progressWidth,
                position: 'relative',
                marginBottom: '8px', // 为百分比文字留出空间
              },
            },
            [
              h(
                'div',
                {
                  class: 'progress-track',
                  style: {
                    width: '100%',
                    height: progressHeight,
                    background: 'rgba(30, 41, 59, 0.8)',
                    borderRadius: '8px',
                    position: 'relative',
                    overflow: 'hidden',
                    border: '2px solid #00d4ff',
                    boxShadow: '0 0 15px rgba(0, 212, 255, 0.4)',
                    backdropFilter: 'blur(10px)',
                  },
                },
                [
                  h('div', {
                    class: 'progress-fill',
                    style: {
                      width: `${Math.max(5, this.animatedProgress || this.progress)}%`,
                      height: '100%',
                      background: 'linear-gradient(90deg, #00d4ff, #ff0080)',
                      borderRadius: '6px',
                      position: 'relative',
                      transition: 'width 0.5s ease-out',
                      boxShadow: '0 0 20px #00d4ff, inset 0 0 10px rgba(255,255,255,0.2)',
                      minWidth: '8px', // 确保进度条有最小可见宽度
                    },
                  }),
                  h('div', {
                    class: 'progress-scanner',
                    style: {
                      position: 'absolute',
                      top: '0',
                      left: '0',
                      width: '40px',
                      height: '100%',
                      background: 'linear-gradient(90deg, transparent, rgba(255,255,255,0.8), transparent)',
                      animation: 'progressScan 2.5s infinite linear',
                      borderRadius: '6px',
                      pointerEvents: 'none',
                    },
                  }),
                ]
              ),
              h(
                'div',
                {
                  class: 'progress-text',
                  style: {
                    position: 'absolute',
                    top: `${parseInt(progressHeight) + 12}px`,
                    left: '50%',
                    transform: 'translateX(-50%)',
                    color: '#00d4ff',
                    fontSize: size.fontSize,
                    fontWeight: '700',
                    textShadow: '0 0 10px #00d4ff, 0 0 20px #00d4ff',
                    fontFamily: "'Orbitron', sans-serif",
                    letterSpacing: '1px',
                  },
                },
                `${Math.floor(this.animatedProgress || this.progress)}%`
              ),
            ]
          ),
        ]
      )
    }

    const renderPulseType = () => {
      const size = getResponsiveSize()

      return h(
        'div',
        {
          class: 'loading-container',
          style: {
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            gap: size.gap,
          },
        },
        [
          h(
            'div',
            {
              class: 'pulse-rings',
              style: {
                position: 'relative',
                width: `${size.iconSize}px`,
                height: `${size.iconSize}px`,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              },
            },
            [
              h('div', { class: 'pulse-ring pulse-ring-1' }),
              h('div', { class: 'pulse-ring pulse-ring-2' }),
              h('div', { class: 'pulse-ring pulse-ring-3' }),
              h('div', {
                style: {
                  width: `${Math.max(8, size.iconSize * 0.1)}px`,
                  height: `${Math.max(8, size.iconSize * 0.1)}px`,
                  borderRadius: '50%',
                  background: this.textColor,
                  boxShadow: `0 0 ${Math.max(10, size.iconSize * 0.17)}px ${this.textColor}`,
                  animation: 'centerPulse 1.5s infinite',
                },
              }),
            ]
          ),
          renderText(size),
        ]
      )
    }

    const renderContent = () => {
      switch (this.type) {
        case 'progress':
          return renderProgressType()
        case 'pulse':
          return renderPulseType()
        default:
          return renderDefaultType()
      }
    }

    return h(
      'div',
      {
        class: 'v-loading-mask',
        style: {
          position: 'absolute',
          top: '0',
          left: '0',
          width: '100%',
          height: '100%',
          zIndex: 'var(--z-index-toast)',
          backdropFilter: 'blur(3px)',
          background: 'transparent',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        },
      },
      [renderContent()]
    )
  },
}

const createLoadingInstance = (el: HTMLElement, options: any = {}) => {
  const app = createApp(LoadingComponent, {
    text: options.text || '',
    type: options.type || 'default',
    progress: options.progress || 0,
    loading: options.loading !== false, // 默认为true，除非明确传入false
    background: options.background,
    textColor: options.textColor,
  })

  const container = document.createElement('div')
  const instance = app.mount(container)

  if (!document.querySelector('#v-loading-styles')) {
    const style = document.createElement('style')
    style.id = 'v-loading-styles'
    style.textContent = `
      .hexagon::before,
      .hexagon::after,
      .hex-inner::before,
      .hex-inner::after {
        content: '';
        position: absolute;
        width: 100%;
        height: 100%;
        background: transparent;
        border: 2px solid var(--color-cyber-blue);
        animation: hexPulse 1.5s infinite alternate, hexGlitch 3s infinite;
      }

      .hexagon::before {
        clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
      }

      .hexagon::after {
        clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
        border-color: var(--color-cyber-pink);
        opacity: 0.7;
        animation-delay: 0.3s;
      }

      .hex-inner {
        position: absolute;
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
      }

      .hex-inner::before {
        clip-path: polygon(50% 10%, 90% 30%, 90% 70%, 50% 90%, 10% 70%, 10% 30%);
        border-color: var(--color-cyber-blue);
        opacity: 0.3;
        transform: scale(0.8);
        animation-delay: 0.6s;
      }

      .hex-inner::after {
        clip-path: polygon(50% 20%, 80% 35%, 80% 65%, 50% 80%, 20% 65%, 20% 35%);
        border-color: var(--color-cyber-pink);
        opacity: 0.5;
        transform: scale(0.6);
        animation-delay: 0.9s;
      }

      .circle {
        position: absolute;
        border-radius: 50%;
        border: 1px solid transparent;
        animation: circlePulse 2s infinite;
      }

      .circle-1 {
        width: 50px;
        height: 50px;
        border-top-color: var(--color-cyber-blue);
        border-right-color: var(--color-cyber-blue);
        animation-delay: 0s;
      }

      .circle-2 {
        width: 70px;
        height: 70px;
        border-right-color: var(--color-cyber-pink);
        border-left-color: var(--color-cyber-pink);
        animation-delay: 0.3s;
      }

      .circle-3 {
        width: 90px;
        height: 90px;
        border-bottom-color: var(--color-primary);
        border-left-color: var(--color-primary);
        animation-delay: 0.6s;
      }

      .scan-line {
        position: absolute;
        width: 100%;
        height: 2px;
        background: linear-gradient(to right, transparent, var(--color-cyber-blue), var(--color-cyber-pink), transparent);
        animation: scanLine 1.5s infinite;
      }

      .text-glitch {
        position: relative;
        animation: textFlicker 3s infinite;
      }

      .text-glitch::before,
      .text-glitch::after {
        content: attr(data-text);
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        opacity: 0.8;
      }

      .text-glitch::before {
        left: 2px;
        color: var(--color-cyber-pink);
        text-shadow: -1px 0 var(--color-cyber-blue);
        clip: rect(44px, 450px, 56px, 0);
        animation: textGlitch 5s infinite linear alternate-reverse;
      }

      .text-glitch::after {
        left: -2px;
        color: var(--color-cyber-blue);
        text-shadow: 2px 0 var(--color-cyber-pink);
        clip: rect(44px, 450px, 56px, 0);
        animation: textGlitch 6s infinite linear alternate-reverse;
      }

      @keyframes hexPulse {
        0% { transform: scale(1); opacity: 1; }
        100% { transform: scale(1.1); opacity: 0.8; }
      }

      @keyframes hexGlitch {
        0%, 100% { transform: translateX(0); }
        25% { transform: translateX(-2px); }
        50% { transform: translateX(2px); }
        75% { transform: translateX(-1px); }
      }

      @keyframes circlePulse {
        0% { transform: rotate(0deg) scale(1); opacity: 0.8; }
        50% { transform: rotate(180deg) scale(1.2); opacity: 0.5; }
        100% { transform: rotate(360deg) scale(1); opacity: 0.8; }
      }

      @keyframes scanLine {
        0% { top: 0%; opacity: 0; }
        25% { opacity: 1; }
        75% { opacity: 1; }
        100% { top: 100%; opacity: 0; }
      }

      @keyframes textFlicker {
        0%, 100% { opacity: 1; }
        8%, 10% { opacity: 0.8; }
        9%, 50%, 70% { opacity: 1; }
        59.5%, 60%, 69.5% { opacity: 0.6; }
      }

      @keyframes textGlitch {
        0% { clip: rect(68px, 350px, 29px, 0); }
        5.88235% { clip: rect(52px, 350px, 16px, 0); }
        11.76471% { clip: rect(22px, 350px, 76px, 0); }
        17.64706% { clip: rect(60px, 350px, 56px, 0); }
        23.52941% { clip: rect(3px, 350px, 86px, 0); }
        29.41176% { clip: rect(87px, 350px, 47px, 0); }
        35.29412% { clip: rect(67px, 350px, 70px, 0); }
        41.17647% { clip: rect(84px, 350px, 23px, 0); }
        47.05882% { clip: rect(89px, 350px, 17px, 0); }
        52.94118% { clip: rect(38px, 350px, 92px, 0); }
        58.82353% { clip: rect(95px, 350px, 3px, 0); }
        64.70588% { clip: rect(12px, 350px, 5px, 0); }
        70.58824% { clip: rect(23px, 350px, 50px, 0); }
        76.47059% { clip: rect(18px, 350px, 30px, 0); }
        82.35294% { clip: rect(2px, 350px, 29px, 0); }
        88.23529% { clip: rect(42px, 350px, 15px, 0); }
        94.11765% { clip: rect(10px, 350px, 49px, 0); }
        100% { clip: rect(31px, 350px, 13px, 0); }
      }

      @keyframes progressScan {
        0% { transform: translateX(-100%); }
        100% { transform: translateX(100%); }
      }

      @keyframes energyRing {
        0% {
          transform: rotate(0deg) scale(1);
          filter: drop-shadow(0 0 10px #00d4ff) brightness(1);
        }
        33% {
          transform: rotate(120deg) scale(1.05);
          filter: drop-shadow(0 0 20px #ff0080) brightness(1.2);
        }
        66% {
          transform: rotate(240deg) scale(0.95);
          filter: drop-shadow(0 0 15px #00ff88) brightness(1.1);
        }
        100% {
          transform: rotate(360deg) scale(1);
          filter: drop-shadow(0 0 10px #00d4ff) brightness(1);
        }
      }

      @keyframes hexagonRotate {
        0% {
          transform: rotate(0deg) scale(1);
          opacity: 0.8;
        }
        25% {
          transform: rotate(90deg) scale(1.1);
          opacity: 1;
        }
        50% {
          transform: rotate(180deg) scale(0.9);
          opacity: 0.9;
        }
        75% {
          transform: rotate(270deg) scale(1.05);
          opacity: 1;
        }
        100% {
          transform: rotate(360deg) scale(1);
          opacity: 0.8;
        }
      }

      @keyframes coreGlow {
        0% {
          transform: scale(1) rotate(0deg);
          box-shadow: 0 0 30px #00d4ff, 0 0 60px #ff0080, inset 0 0 20px rgba(255,255,255,0.4);
          filter: brightness(1);
        }
        50% {
          transform: scale(1.1) rotate(180deg);
          box-shadow: 0 0 40px #ff0080, 0 0 80px #00ff88, inset 0 0 30px rgba(255,255,255,0.6);
          filter: brightness(1.3);
        }
        100% {
          transform: scale(1) rotate(360deg);
          box-shadow: 0 0 30px #00d4ff, 0 0 60px #ff0080, inset 0 0 20px rgba(255,255,255,0.4);
          filter: brightness(1);
        }
      }

      @keyframes dataFlow {
        0% { transform: rotate(0deg); opacity: 0.6; }
        50% { opacity: 1; }
        100% { transform: rotate(360deg); opacity: 0.6; }
      }

      @keyframes dataNodePulse {
        0% {
          transform: scale(1);
          opacity: 0.7;
        }
        100% {
          transform: scale(1.5);
          opacity: 1;
        }
      }

      @keyframes scanBeam {
        0% {
          transform: rotate(0deg);
          opacity: 0.8;
        }
        25% { opacity: 1; }
        50% {
          transform: rotate(180deg);
          opacity: 0.9;
        }
        75% { opacity: 1; }
        100% {
          transform: rotate(360deg);
          opacity: 0.8;
        }
      }

      @keyframes energyCore {
        0% {
          transform: scale(1);
          box-shadow: 0 0 20px #ffffff, 0 0 40px #00d4ff, 0 0 60px rgba(0, 212, 255, 0.3);
        }
        100% {
          transform: scale(1.3);
          box-shadow: 0 0 30px #ffffff, 0 0 60px #00d4ff, 0 0 90px rgba(0, 212, 255, 0.5);
        }
      }

      @keyframes gridShift {
        0% {
          transform: rotate(0deg);
          opacity: 0.3;
        }
        50% {
          opacity: 0.1;
        }
        100% {
          transform: rotate(360deg);
          opacity: 0.3;
        }
      }

      .pulse-ring {
        position: absolute;
        border: 2px solid var(--color-cyber-blue);
        border-radius: 50%;
        animation: pulseRing 2s infinite;
      }

      .pulse-ring-1 {
        width: 40px;
        height: 40px;
        animation-delay: 0s;
      }

      .pulse-ring-2 {
        width: 60px;
        height: 60px;
        border-color: var(--color-cyber-pink);
        animation-delay: 0.3s;
      }

      .pulse-ring-3 {
        width: 80px;
        height: 80px;
        border-color: var(--color-cyber-purple);
        animation-delay: 0.6s;
      }

      @keyframes pulseRing {
        0% {
          transform: scale(0.8);
          opacity: 1;
        }
        50% {
          transform: scale(1.2);
          opacity: 0.7;
        }
        100% {
          transform: scale(1.6);
          opacity: 0;
        }
      }

      @keyframes centerPulse {
        0%, 100% {
          transform: scale(1);
          opacity: 1;
        }
        50% {
          transform: scale(1.5);
          opacity: 0.8;
        }
      }

      @keyframes smoothRotate {
        0% {
          transform: rotate(0deg);
        }
        100% {
          transform: rotate(360deg);
        }
      }

      @keyframes centerGlow {
        0%, 100% {
          transform: scale(1);
          opacity: 0.8;
          box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.6);
        }
        50% {
          transform: scale(1.1);
          opacity: 1;
          box-shadow: 0 0 16px rgba(var(--color-brand-500-rgb), 0.8);
        }
      }

      @keyframes textFade {
        0%, 100% {
          opacity: 0.7;
        }
        50% {
          opacity: 1;
        }
      }

    `
    document.head.appendChild(style)
  }

  return {
    instance,
    element: container.firstChild as HTMLElement,
    app,
  }
}

const showLoading = (el: HTMLElement, options: any = {}) => {
  const position = getComputedStyle(el).position
  if (position === 'static') {
    el.style.position = 'relative'
  }

  const loadingInstance = createLoadingInstance(el, options)
  el.appendChild(loadingInstance.element)
  loadingInstances.set(el, loadingInstance)
}

const hideLoading = (el: HTMLElement) => {
  const loadingInstance = loadingInstances.get(el)
  if (loadingInstance) {
    if (loadingInstance.element && loadingInstance.element.parentNode) {
      loadingInstance.element.parentNode.removeChild(loadingInstance.element)
    }
    loadingInstance.app.unmount()
    loadingInstances.delete(el)
  }
}

const parseLoadingOptions = (value: any): { shouldShow: boolean; options: any } => {
  if (typeof value === 'boolean') {
    return { shouldShow: value, options: {} }
  } else if (typeof value === 'object' && value !== null) {
    return {
      shouldShow: value.loading || value.visible || false,
      options: value,
    }
  } else if (typeof value === 'string') {
    return { shouldShow: true, options: { text: value } }
  } else if (value) {
    return { shouldShow: true, options: {} }
  }
  return { shouldShow: false, options: {} }
}

export const loading: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { shouldShow, options } = parseLoadingOptions(binding.value)

    if (shouldShow) {
      showLoading(el, options)
    }
  },

  updated(el: HTMLElement, binding: DirectiveBinding) {
    const { value, oldValue } = binding

    if (value === oldValue) {
      return
    }

    const isCurrentlyLoading = loadingInstances.has(el)
    const { shouldShow, options } = parseLoadingOptions(value)

    if (shouldShow && !isCurrentlyLoading) {
      showLoading(el, options)
    } else if (!shouldShow && isCurrentlyLoading) {
      hideLoading(el)
    } else if (shouldShow && isCurrentlyLoading) {
      const currentInstance = loadingInstances.get(el)
      if (currentInstance && currentInstance.app) {
        const component = currentInstance.app._instance
        if (component && component.props) {
          Object.assign(component.props, {
            text: options.text || '',
            type: options.type || 'default',
            progress: options.progress || 0,
            loading: options.loading !== false,
            background: options.background,
            textColor: options.textColor,
          })
          component.update?.()
        } else {
          hideLoading(el)
          showLoading(el, options)
        }
      } else {
        hideLoading(el)
        showLoading(el, options)
      }
    }
  },

  beforeUnmount(el: HTMLElement) {
    hideLoading(el)
  },
}

export default {
  install(app: any) {
    app.directive('loading', loading)
  },
}
