<script setup lang="ts">
  import { useParticleBackground } from '@/composables/useLayoutManager'

  const { getParticleStyle } = useParticleBackground()
</script>

<template>
  <div class="cyber-background cyber-optimized">
    <div class="cyber-grid-base" />

    <div class="cyber-glow-base" />

    <div class="cyber-particles-container cyber-particles-mobile-hidden">
      <div v-for="n in 20" :key="n" class="cyber-particle" :style="getParticleStyle(n)" />
    </div>

    <div class="cyber-scanlines" />
  </div>
</template>
