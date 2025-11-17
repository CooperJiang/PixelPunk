<script setup lang="ts">
  defineOptions({
    name: 'CyberHomeBackground',
  })

  interface Props {
    showEnhancedEffects?: boolean
    showDigitalRain?: boolean
    showRedOrb?: boolean
    showAurora?: boolean
  }

  withDefaults(defineProps<Props>(), {
    showEnhancedEffects: false,
    showDigitalRain: false,
    showRedOrb: false,
    showAurora: false,
  })
</script>

<template>
  <div class="cyber-background">
    <div class="grid-overlay" />

    <div v-if="showAurora" class="aurora-bands">
      <div class="band band-1" />
      <div class="band band-2" />
      <div class="band band-3" />
    </div>

    <div v-if="showEnhancedEffects" class="enhanced-effects">
      <div class="ambient-glow">
        <div class="glow glow-1" />
        <div class="glow glow-2" />
        <div v-if="showRedOrb" class="glow glow-3" />
      </div>

      <div class="elegant-beams">
        <div class="beam-line" />
      </div>

      <div v-if="showDigitalRain" class="digital-rain">
        <div class="rain-column rain-1">01</div>
        <div class="rain-column rain-2">10</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-background {
    position: fixed;
    inset: 0;
    z-index: 0;
    background:
      radial-gradient(ellipse at top left, rgba(var(--color-error-rgb), 0.15) 0%, transparent 50%),
      radial-gradient(ellipse at bottom right, rgba(var(--color-brand-500-rgb), 0.15) 0%, transparent 50%),
      radial-gradient(ellipse at center, rgba(var(--color-error-rgb), 0.1) 0%, transparent 70%),
      linear-gradient(180deg, var(--color-background-900) 0%, var(--color-background-800) 100%);
  }

  .grid-overlay {
    @apply absolute;
    top: -100%;
    left: -100%;
    width: 300%;
    height: 300%;
    background-image:
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.08) 1px, transparent 1px),
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.08) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: gridMove 60s linear infinite;
  }

  @keyframes gridMove {
    0% {
      transform: translate(0, 0);
    }
    100% {
      transform: translate(33.333%, 33.333%);
    }
  }

  .aurora-bands {
    position: absolute;
    inset: 0;
    overflow: hidden;
    z-index: 1;
    pointer-events: none;
  }

  .band {
    position: absolute;
    width: 200%;
    height: 300px;
    filter: blur(80px);
    opacity: 0;
    animation: auroraFlow 25s ease-in-out infinite;
  }

  .band-1 {
    background: linear-gradient(
      90deg,
      transparent 0%,
      rgba(var(--color-brand-500-rgb), 0.15) 20%,
      rgba(var(--color-error-rgb), 0.12) 50%,
      rgba(var(--color-error-rgb), 0.1) 80%,
      transparent 100%
    );
    top: 10%;
    left: -50%;
    animation-delay: 0s;
  }

  .band-2 {
    background: linear-gradient(
      90deg,
      transparent 0%,
      rgba(var(--color-error-rgb), 0.12) 20%,
      rgba(var(--color-brand-500-rgb), 0.15) 50%,
      rgba(var(--color-error-rgb), 0.1) 80%,
      transparent 100%
    );
    top: 40%;
    left: -50%;
    animation-delay: -8s;
    animation-duration: 30s;
  }

  .band-3 {
    background: linear-gradient(
      90deg,
      transparent 0%,
      rgba(var(--color-error-rgb), 0.15) 20%,
      rgba(var(--color-error-rgb), 0.12) 50%,
      rgba(var(--color-brand-500-rgb), 0.1) 80%,
      transparent 100%
    );
    top: 70%;
    left: -50%;
    animation-delay: -16s;
    animation-duration: 35s;
  }

  @keyframes auroraFlow {
    0% {
      transform: translateX(0) translateY(0) rotate(-5deg);
      opacity: 0;
    }
    20% {
      opacity: 0.5;
    }
    80% {
      opacity: 0.5;
    }
    100% {
      transform: translateX(50%) translateY(-20px) rotate(5deg);
      opacity: 0;
    }
  }

  .enhanced-effects {
    position: absolute;
    inset: 0;
    pointer-events: none;
    overflow: hidden;
    z-index: 1;
  }

  .ambient-glow {
    position: absolute;
    inset: 0;
    filter: blur(100px);
  }

  .glow {
    @apply absolute rounded-full;
    animation: glowPulse 20s ease-in-out infinite;
  }

  .glow-1 {
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.15) 0%, transparent 70%);
    top: -300px;
    left: -300px;
  }

  .glow-2 {
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, rgba(var(--color-error-rgb), 0.12) 0%, transparent 70%);
    bottom: -250px;
    left: 50%;
    transform: translateX(-50%);
    animation-delay: -10s;
  }

  .glow-3 {
    width: 500px;
    height: 500px;
    background: radial-gradient(
      circle,
      rgba(var(--color-error-rgb), 0.35) 0%,
      rgba(var(--color-error-rgb), 0.15) 50%,
      transparent 70%
    );
    bottom: -250px;
    right: -250px;
    animation-delay: -5s;
  }

  @keyframes glowPulse {
    0%,
    100% {
      opacity: 0.6;
      transform: scale(1);
    }
    50% {
      opacity: 0.9;
      transform: scale(1.1);
    }
  }

  .elegant-beams {
    @apply absolute inset-0;
  }

  .beam-line {
    position: absolute;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(var(--color-brand-500-rgb), 0.3) 30%,
      rgba(var(--color-brand-500-rgb), 0.5) 50%,
      rgba(var(--color-brand-500-rgb), 0.3) 70%,
      transparent
    );
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
    animation: elegantBeam 25s linear infinite;
  }

  @keyframes elegantBeam {
    0% {
      top: -5%;
      opacity: 0;
    }
    5% {
      opacity: 0.6;
    }
    95% {
      opacity: 0.6;
    }
    100% {
      top: 105%;
      opacity: 0;
    }
  }

  .digital-rain {
    @apply absolute inset-0;
    overflow: hidden;
    font-family: 'JetBrains Mono', 'Fira Code', 'Courier New', monospace;
    font-size: 11px;
    font-weight: 600;
  }

  .rain-column {
    position: absolute;
    color: rgba(var(--color-brand-500-rgb), 0.4);
    writing-mode: vertical-lr;
    letter-spacing: 3px;
    opacity: 0;
    text-shadow: 0 0 10px currentColor;
  }

  .rain-1 {
    left: 15%;
    animation: rainFall 18s linear infinite;
  }

  .rain-2 {
    right: 15%;
    color: rgba(var(--color-error-rgb), 0.35);
    animation: rainFall 22s linear infinite 8s;
  }

  @keyframes rainFall {
    0% {
      top: -120px;
      opacity: 0;
    }
    10% {
      opacity: 0.5;
    }
    90% {
      opacity: 0.5;
    }
    100% {
      top: 100vh;
      opacity: 0;
    }
  }
</style>
