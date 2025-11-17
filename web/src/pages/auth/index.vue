<script setup lang="ts">
  import { ref } from 'vue'
  import AuthPanel from './components/AuthPanel/index.vue'

  const isPasswordFocused = ref(false)

  const handlePasswordFocus = (isFocused: boolean) => {
    isPasswordFocused.value = isFocused
  }
</script>

<template>
  <div class="auth-page">
    <CyberHomeBackground :show-enhanced-effects="true" :show-digital-rain="false" :show-red-orb="true" :show-aurora="true" />

    <div class="floating-particles">
      <div v-for="i in 8" :key="i" class="particle" :style="{ '--delay': i * 0.5 + 's', '--duration': 8 + i * 2 + 's' }" />
    </div>

    <div class="page-content">
      <div class="main-content">
        <div class="auth-panel-container">
          <div class="panel-glow-effect" />

          <div class="corner-lines">
            <div class="corner-line top-left" />
            <div class="corner-line top-right" />
            <div class="corner-line bottom-left" />
            <div class="corner-line bottom-right" />
          </div>

          <CyborgCharacter size="medium" position="top-left" :cover-eyes="isPasswordFocused" />

          <div class="floating-tech-element left-top" />
          <div class="floating-tech-element right-bottom" />

          <div class="cyber-header-decoration">
            <span class="cyber-line" />
            <span class="cyber-circle" />
            <span class="cyber-square" />
          </div>

          <div class="scan-line-container">
            <div class="scan-line" />
          </div>

          <AuthPanel @password-focus="handlePasswordFocus" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./index.scss"></style>
