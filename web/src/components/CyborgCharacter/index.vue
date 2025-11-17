<script setup lang="ts">
  interface Props {
    size?: 'small' | 'medium' | 'large'
    position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
    coverEyes?: boolean
  }

  withDefaults(defineProps<Props>(), {
    size: 'medium',
    position: 'top-left',
    coverEyes: false,
  })

  defineOptions({
    name: 'CyborgCharacter',
  })
</script>

<template>
  <div class="cyborg-character" :class="[`size-${size}`, `position-${position}`, { 'cover-eyes': coverEyes }]">
    <div class="cyborg-head">
      <div class="cyborg-eyes">
        <div class="cyborg-eye left" />
        <div class="cyborg-eye right" />
      </div>
      <div class="circuit-lines">
        <div class="circuit line1" />
      </div>
    </div>

    <div class="cyborg-body">
      <div class="cyborg-arm left-arm">
        <div class="arm-segment upper" />
        <div class="arm-segment lower" />
        <div class="arm-hand" />
      </div>

      <div class="cyborg-arm right-arm">
        <div class="arm-segment upper" />
        <div class="arm-segment lower" />
        <div class="arm-hand" />
      </div>

      <div class="power-core" />
      <div class="body-panel left-panel" />
      <div class="body-panel right-panel" />
    </div>
  </div>
</template>

<style scoped>
  .cyborg-character {
    position: absolute;
    z-index: 5;
    pointer-events: none;
    animation: cyborg-float 4s ease-in-out infinite;
  }

  .cyborg-character.size-small {
    width: 50px;
    height: 70px;
  }

  .cyborg-character.size-medium {
    width: 70px;
    height: 100px;
  }

  .cyborg-character.size-large {
    width: 90px;
    height: 130px;
  }

  .cyborg-character.position-top-left {
    top: -50px;
    left: -40px;
    transform: rotate(-10deg);
  }

  .cyborg-character.position-top-right {
    top: -30px;
    right: -20px;
    transform: rotate(10deg) scaleX(-1);
  }

  .cyborg-character.position-bottom-left {
    bottom: -80px;
    left: -40px;
    transform: rotate(10deg);
  }

  .cyborg-character.position-bottom-right {
    bottom: -80px;
    right: -40px;
    transform: rotate(-10deg) scaleX(-1);
  }

  .cyborg-arm {
    position: absolute;
    z-index: 10;
    transform-origin: top center;
  }

  .left-arm {
    top: 15px;
    left: -6px;
    animation: arm-wave-left 3s ease-in-out infinite;
    transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .right-arm {
    top: 15px;
    right: -6px;
    animation: arm-wave-right 3.5s ease-in-out infinite;
    transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .cyborg-character.cover-eyes .left-arm {
    animation: cover-eyes-left 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards;
  }

  .cyborg-character.cover-eyes .right-arm {
    animation: cover-eyes-right 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards;
  }

  .cyborg-character.cover-eyes .cyborg-eye {
    height: 1px !important;
    opacity: 0.3 !important;
    animation: none;
  }

  .arm-segment {
    width: 2px;
    background: linear-gradient(180deg, var(--color-background-700), var(--color-background-800));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  .arm-segment.upper {
    height: 15px;
  }

  .arm-segment.lower {
    height: 12px;
    margin-top: -1px;
  }

  .arm-hand {
    width: 4px;
    height: 4px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    margin-top: -1px;
    margin-left: -1px;
    box-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.6);
  }

  .cyborg-head {
    position: relative;
    width: 40px;
    height: 40px;
    background: var(--color-background-800);
    border-radius: 40% 40% 50% 50%;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    overflow: visible;
    animation: cyborg-head-tilt 6s ease-in-out infinite;
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .size-small .cyborg-head {
    width: 30px;
    height: 30px;
  }

  .size-large .cyborg-head {
    width: 50px;
    height: 50px;
  }

  .cyborg-eyes {
    position: absolute;
    display: flex;
    justify-content: space-around;
    width: 70%;
    top: 40%;
    left: 15%;
  }

  .cyborg-eye {
    width: 8px;
    height: 4px;
    border-radius: var(--radius-sm);
    background: var(--color-brand-500);
    animation: cyborg-blink 4s infinite;
    transition:
      height 0.3s ease,
      opacity 0.3s ease;
  }

  .size-small .cyborg-eye {
    width: 6px;
    height: 3px;
  }

  .size-large .cyborg-eye {
    width: 10px;
    height: 5px;
  }

  .circuit-lines {
    position: absolute;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
  }

  .circuit {
    position: absolute;
    height: 1px;
    background: var(--color-brand-500);
  }

  .circuit.line1 {
    width: 60%;
    top: 70%;
    left: 20%;
    animation: circuit-pulse 4s infinite;
  }

  .cyborg-body {
    position: relative;
    top: -5px;
    left: -5px;
    width: 50px;
    height: 60px;
    border-radius: 30% 30% 40% 40%;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    background: linear-gradient(180deg, var(--color-background-900) 0%, var(--color-background-700) 100%);
    overflow: visible;
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .size-small .cyborg-body {
    width: 38px;
    height: 45px;
  }

  .size-large .cyborg-body {
    width: 62px;
    height: 75px;
  }

  .power-core {
    position: absolute;
    top: 20px;
    left: 18px;
    width: 15px;
    height: 15px;
    border-radius: var(--radius-full);
    background: var(--color-brand-500);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
    animation: core-pulse 3s infinite;
    z-index: 5;
  }

  .body-panel {
    position: absolute;
    width: 2px;
    height: 20px;
    background: rgba(var(--color-brand-500-rgb), 0.3);
    top: 15px;
    z-index: 1;
  }

  .left-panel {
    left: 8px;
    animation: panel-glow-left 4s ease-in-out infinite;
  }

  .right-panel {
    right: 8px;
    animation: panel-glow-right 4.5s ease-in-out infinite;
  }

  .size-small .power-core {
    width: 12px;
    height: 12px;
    top: 15px;
    left: 13px;
  }

  .size-large .power-core {
    width: 18px;
    height: 18px;
    top: 25px;
    left: 22px;
  }

  @keyframes cyborg-float {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-8px);
    }
  }

  .position-top-left {
    animation: cyborg-float-top-left 4s ease-in-out infinite;
  }

  @keyframes cyborg-float-top-left {
    0%,
    100% {
      transform: rotate(-10deg) translateY(0px);
    }
    50% {
      transform: rotate(-10deg) translateY(-8px);
    }
  }

  .position-top-right {
    animation: cyborg-float-top-right 4s ease-in-out infinite;
  }

  @keyframes cyborg-float-top-right {
    0%,
    100% {
      transform: rotate(10deg) scaleX(-1) translateY(0px);
    }
    50% {
      transform: rotate(10deg) scaleX(-1) translateY(-8px);
    }
  }

  .position-bottom-left {
    animation: cyborg-float-bottom-left 4s ease-in-out infinite;
  }

  @keyframes cyborg-float-bottom-left {
    0%,
    100% {
      transform: rotate(10deg) translateY(0px);
    }
    50% {
      transform: rotate(10deg) translateY(8px);
    }
  }

  .position-bottom-right {
    animation: cyborg-float-bottom-right 4s ease-in-out infinite;
  }

  @keyframes cyborg-float-bottom-right {
    0%,
    100% {
      transform: rotate(-10deg) scaleX(-1) translateY(0px);
    }
    50% {
      transform: rotate(-10deg) scaleX(-1) translateY(8px);
    }
  }

  @keyframes cyborg-head-tilt {
    0%,
    100% {
      transform: rotate(0deg);
    }
    25% {
      transform: rotate(-3deg);
    }
    75% {
      transform: rotate(3deg);
    }
  }

  @keyframes cyborg-blink {
    0%,
    90%,
    100% {
      height: 4px;
      opacity: 1;
    }
    92%,
    96% {
      height: 1px;
      opacity: 0.8;
    }
    94% {
      height: 0.5px;
      opacity: 0.6;
    }
  }

  @keyframes circuit-pulse {
    0%,
    100% {
      opacity: 0.3;
      box-shadow: 0 0 3px rgba(var(--color-brand-500-rgb), 0.3);
    }
    50% {
      opacity: 1;
      box-shadow:
        0 0 8px rgba(var(--color-brand-500-rgb), 0.8),
        0 0 12px rgba(var(--color-brand-500-rgb), 0.4);
    }
  }

  @keyframes core-pulse {
    0%,
    100% {
      transform: scale(1);
      box-shadow:
        0 0 8px rgba(var(--color-brand-500-rgb), 0.5),
        0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
    }
    50% {
      transform: scale(1.2);
      box-shadow:
        0 0 15px rgba(var(--color-brand-500-rgb), 0.9),
        0 0 25px rgba(var(--color-brand-500-rgb), 0.6),
        0 0 35px rgba(var(--color-brand-500-rgb), 0.3);
    }
  }

  @keyframes arm-wave-left {
    0%,
    100% {
      transform: rotate(0deg);
    }
    50% {
      transform: rotate(-5deg);
    }
  }

  @keyframes arm-wave-right {
    0%,
    100% {
      transform: rotate(0deg);
    }
    50% {
      transform: rotate(5deg);
    }
  }

  @keyframes panel-glow-left {
    0%,
    100% {
      opacity: 0.3;
      box-shadow: 0 0 3px rgba(var(--color-brand-500-rgb), 0.3);
    }
    50% {
      opacity: 1;
      box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.8);
    }
  }

  @keyframes panel-glow-right {
    0%,
    100% {
      opacity: 0.3;
      box-shadow: 0 0 3px rgba(var(--color-brand-500-rgb), 0.3);
    }
    50% {
      opacity: 1;
      box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.8);
    }
  }

  @keyframes cover-eyes-left {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(-110deg);
    }
  }

  @keyframes cover-eyes-right {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(110deg);
    }
  }

  @media (max-width: 768px) {
    .cyborg-character.size-medium {
      transform: scale(0.8);
    }

    .cyborg-character.position-top-left {
      top: -60px;
      left: -30px;
    }

    .cyborg-character.position-top-right {
      top: -60px;
      right: -30px;
    }
  }

  @media (max-width: 480px) {
    .cyborg-character.size-medium {
      transform: scale(0.7);
    }

    .cyborg-character.position-top-left {
      top: -50px;
      left: -20px;
    }

    .cyborg-character.position-top-right {
      top: -50px;
      right: -20px;
    }
  }
</style>
