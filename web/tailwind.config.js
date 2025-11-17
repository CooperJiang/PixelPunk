/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      /* ================================
       * 🎨 核心颜色系统 - 映射 CSS 变量
       * ================================ */
      colors: {
        brand: {
          50: 'var(--color-brand-50)',
          100: 'var(--color-brand-100)',
          200: 'var(--color-brand-200)',
          300: 'var(--color-brand-300)',
          400: 'var(--color-brand-400)',
          500: 'var(--color-brand-500)',
          600: 'var(--color-brand-600)',
          700: 'var(--color-brand-700)',
          800: 'var(--color-brand-800)',
          900: 'var(--color-brand-900)',
          DEFAULT: 'var(--color-brand-500)',
        },

        background: {
          50: 'var(--color-background-50)',
          100: 'var(--color-background-100)',
          200: 'var(--color-background-200)',
          300: 'var(--color-background-300)',
          400: 'var(--color-background-400)',
          500: 'var(--color-background-500)',
          600: 'var(--color-background-600)',
          700: 'var(--color-background-700)',
          800: 'var(--color-background-800)',
          900: 'var(--color-background-900)',
          DEFAULT: 'var(--color-background-900)',
        },

        content: {
          DEFAULT: 'var(--color-content-default)',
          heading: 'var(--color-content-heading)',
          muted: 'var(--color-content-muted)',
          disabled: 'var(--color-content-disabled)',
        },

        success: {
          50: 'var(--color-success-50)',
          100: 'var(--color-success-100)',
          200: 'var(--color-success-200)',
          300: 'var(--color-success-300)',
          400: 'var(--color-success-400)',
          500: 'var(--color-success-500)',
          600: 'var(--color-success-600)',
          700: 'var(--color-success-700)',
          800: 'var(--color-success-800)',
          900: 'var(--color-success-900)',
          DEFAULT: 'var(--color-success-500)',
        },
        warning: {
          50: 'var(--color-warning-50)',
          100: 'var(--color-warning-100)',
          200: 'var(--color-warning-200)',
          300: 'var(--color-warning-300)',
          400: 'var(--color-warning-400)',
          500: 'var(--color-warning-500)',
          600: 'var(--color-warning-600)',
          700: 'var(--color-warning-700)',
          800: 'var(--color-warning-800)',
          900: 'var(--color-warning-900)',
          DEFAULT: 'var(--color-warning-500)',
        },
        error: {
          50: 'var(--color-error-50)',
          100: 'var(--color-error-100)',
          200: 'var(--color-error-200)',
          300: 'var(--color-error-300)',
          400: 'var(--color-error-400)',
          500: 'var(--color-error-500)',
          600: 'var(--color-error-600)',
          700: 'var(--color-error-700)',
          800: 'var(--color-error-800)',
          900: 'var(--color-error-900)',
          DEFAULT: 'var(--color-error-500)',
        },
        info: {
          50: 'var(--color-info-50)',
          100: 'var(--color-info-100)',
          200: 'var(--color-info-200)',
          300: 'var(--color-info-300)',
          400: 'var(--color-info-400)',
          500: 'var(--color-info-500)',
          600: 'var(--color-info-600)',
          700: 'var(--color-info-700)',
          800: 'var(--color-info-800)',
          900: 'var(--color-info-900)',
          DEFAULT: 'var(--color-info-500)',
        },

        /* ================================
         * 🎭 交互状态系统
         * ================================ */
        hover: {
          bg: 'var(--color-hover-bg)',
          'bg-neutral': 'var(--color-hover-bg-neutral)',
          border: 'var(--color-hover-border)',
          text: 'var(--color-hover-text)',
        },
        active: {
          bg: 'var(--color-active-bg)',
          border: 'var(--color-active-border)',
        },
        focus: {
          ring: 'var(--color-focus-ring)',
          border: 'var(--color-focus-border)',
        },
        selected: {
          bg: 'var(--color-selected-bg)',
          border: 'var(--color-selected-border)',
        },
        disabled: {
          bg: 'var(--color-disabled-bg)',
          text: 'var(--color-disabled-text)',
          border: 'var(--color-disabled-border)',
        },

        /* ================================
         * 📝 表单专用系统
         * ================================ */
        input: {
          bg: 'var(--color-input-bg)',
          'bg-focus': 'var(--color-input-bg-focus)',
          'bg-disabled': 'var(--color-input-bg-disabled)',
          text: 'var(--color-input-text)',
          'text-disabled': 'var(--color-input-text-disabled)',
          border: 'var(--color-input-border)',
          'border-hover': 'var(--color-input-border-hover)',
          'border-focus': 'var(--color-input-border-focus)',
          'border-error': 'var(--color-input-border-error)',
        },
        placeholder: 'var(--color-placeholder)',
        'placeholder-focus': 'var(--color-placeholder-focus)',
        label: 'var(--color-label)',
        'label-required': 'var(--color-label-required)',

        /* ================================
         * 🎨 特效系统
         * ================================ */
        overlay: {
          light: 'var(--color-overlay-light)',
          medium: 'var(--color-overlay-medium)',
          heavy: 'var(--color-overlay-heavy)',
        },
        skeleton: {
          base: 'var(--color-skeleton-base)',
          highlight: 'var(--color-skeleton-highlight)',
        },

        /* ================================
         * 🎯 徽章系统
         * ================================ */
        badge: {
          'primary-bg': 'var(--color-badge-primary-bg)',
          'primary-text': 'var(--color-badge-primary-text)',
          'primary-border': 'var(--color-badge-primary-border)',
          'accent-bg': 'var(--color-badge-accent-bg)',
          'accent-text': 'var(--color-badge-accent-text)',
          'accent-border': 'var(--color-badge-accent-border)',
          'success-bg': 'var(--color-badge-success-bg)',
          'success-text': 'var(--color-badge-success-text)',
          'success-border': 'var(--color-badge-success-border)',
          'warning-bg': 'var(--color-badge-warning-bg)',
          'warning-text': 'var(--color-badge-warning-text)',
          'warning-border': 'var(--color-badge-warning-border)',
          'error-bg': 'var(--color-badge-error-bg)',
          'error-text': 'var(--color-badge-error-text)',
          'error-border': 'var(--color-badge-error-border)',
          'neutral-bg': 'var(--color-badge-neutral-bg)',
          'neutral-text': 'var(--color-badge-neutral-text)',
          'neutral-border': 'var(--color-badge-neutral-border)',
        },
      },

      /* ================================
       * 🔲 边框系统
       * ================================ */
      borderColor: ({ theme }) => ({
        ...theme('colors'),
        subtle: 'var(--color-border-subtle)',
        default: 'var(--color-border-default)',
        strong: 'var(--color-border-strong)',
        DEFAULT: 'var(--color-border-default)',
      }),

      /* ================================
       * ➗ 分割线系统
       * ================================ */
      divideColor: ({ theme }) => ({
        ...theme('colors'),
        light: 'var(--color-divider-light)',
        medium: 'var(--color-divider-medium)',
        strong: 'var(--color-divider-strong)',
        DEFAULT: 'var(--color-divider-light)',
      }),

      /* ================================
       * ✨ 阴影系统
       * ================================ */
      boxShadow: {
        // 常规阴影
        xs: 'var(--shadow-xs)',
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
        lg: 'var(--shadow-lg)',
        xl: 'var(--shadow-xl)',
        inner: 'var(--shadow-inner)',

        // 赛博朋克发光阴影
        'glow-sm': 'var(--shadow-glow-sm)',
        'glow-md': 'var(--shadow-glow-md)',
        'glow-lg': 'var(--shadow-glow-lg)',
        'glow-xl': 'var(--shadow-glow-xl)',

        // 赛博朋克复合阴影
        'cyber-sm': 'var(--shadow-cyber-sm)',
        'cyber-md': 'var(--shadow-cyber-md)',
        'cyber-lg': 'var(--shadow-cyber-lg)',
        'cyber-xl': 'var(--shadow-cyber-xl)',
      },

      /* ================================
       * 📐 间距系统
       * ================================ */
      spacing: {
        xs: 'var(--space-xs)',
        sm: 'var(--space-sm)',
        md: 'var(--space-md)',
        lg: 'var(--space-lg)',
        xl: 'var(--space-xl)',
        '2xl': 'var(--space-2xl)',
        '3xl': 'var(--space-3xl)',
      },

      /* ================================
       * 🔘 圆角系统
       * ================================ */
      borderRadius: {
        none: '0',
        sm: 'var(--radius-sm)',
        DEFAULT: 'var(--radius-md)' /* rounded 默认使用 md (4px) */,
        md: 'var(--radius-md)',
        lg: 'var(--radius-sm)',
        xl: 'var(--radius-sm)',
        '2xl': 'var(--radius-2xl)',
        full: 'var(--radius-full)',
      },

      /* ================================
       * ⏱️ 过渡系统
       * ================================ */
      transitionDuration: {
        fast: 'var(--transition-fast)',
        normal: 'var(--transition-normal)',
        slow: 'var(--transition-slow)',
      },

      /* ================================
       * 🎬 动画系统
       * ================================ */
      animation: {
        glitch: 'glitch 1s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate',
      },
      keyframes: {
        glitch: {
          '0%, 100%': { transform: 'translate(0)' },
          '33%': { transform: 'translate(-5px, 2px)' },
          '66%': { transform: 'translate(5px, -2px)' },
        },
        glow: {
          '0%': { 'text-shadow': '0 0 5px #ff00ff, 0 0 10px #ff00ff' },
          '100%': { 'text-shadow': '0 0 10px #ff0066, 0 0 20px #ff0066' },
        },
      },

      /* ================================
       * 📱 响应式断点
       * ================================ */
      screens: {
        '3xl': '2560px',
        '4xl': '3840px',
      },

      /* ================================
       * 🔤 字体系统
       * ================================ */
      fontFamily: {
        orbitron: ['IBM Plex Sans', 'Arial', 'sans-serif'],
        'roboto-mono': ['IBM Plex Mono', 'Monaco', 'Consolas', 'monospace'],
        heading: ['IBM Plex Sans', 'Arial', 'sans-serif'],
        mono: ['IBM Plex Mono', 'Monaco', 'Consolas', 'monospace'],
      },
    },
  },
  plugins: [],
}
