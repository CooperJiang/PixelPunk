<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { StorageUtil } from '@/utils/storage/storage'
  import router from '@/router'
  import { getHealthCheck } from '@/api/common'
  import { AccessRestrictionCodes, AccessRestrictionKeys } from '@/utils/network/http'
  import { useTexts } from '@/composables/useTexts'

  const { IP_NOT_IN_WHITELIST, IP_IN_BLACKLIST, DOMAIN_NOT_IN_WHITELIST, DOMAIN_IN_BLACKLIST, USER_ACCOUNT_DISABLED } =
    AccessRestrictionCodes

  const { IP_RESTRICTED_KEY, IP_ADDRESS_KEY, IP_ERROR_CODE_KEY, IP_ERROR_MSG_KEY, DOMAIN_KEY, USER_DISABLED_KEY } =
    AccessRestrictionKeys

  const { $t } = useTexts()

  const message = ref($t('refuse.defaultMessage'))
  const ipAddress = ref($t('refuse.unknown'))
  const errorCode = ref('')
  const domainName = ref($t('refuse.unknown'))
  const restrictionType = ref('IP') // 'IP', 'DOMAIN', 或 'USER'

  /* 动态标题和徽章 */
  const mainTitle = computed(() => {
    if (restrictionType.value === 'USER') {
      return $t('refuse.titles.accountDisabled')
    }
    return $t('refuse.titles.accessDenied')
  })

  const analysisTitle = computed(() => {
    if (restrictionType.value === 'USER') {
      return $t('refuse.analysisTitle.account')
    }
    return $t('refuse.analysisTitle.system')
  })

  const errorBadgeText = computed(() => {
    if (restrictionType.value === 'USER') {
      return 'ACCOUNT DISABLED'
    }
    if (restrictionType.value === 'DOMAIN') {
      return 'DOMAIN RESTRICTED'
    }
    return 'IP RESTRICTED'
  })

  const errorBadgeClass = computed(() => {
    if (restrictionType.value === 'USER') {
      return 'bg-error-500 text-content-heading'
    }
    return 'bg-warning-500 text-background-900'
  })

  const formattedTimestamp = computed(() => {
    const now = new Date()
    return `${now.getFullYear()}.${String(now.getMonth() + 1).padStart(2, '0')}.${String(now.getDate()).padStart(2, '0')} / ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`
  })

  const clearIpRestriction = () => {
    StorageUtil.remove(IP_RESTRICTED_KEY)
    StorageUtil.remove(IP_ADDRESS_KEY)
    StorageUtil.remove(IP_ERROR_CODE_KEY)
    StorageUtil.remove(IP_ERROR_MSG_KEY)
    StorageUtil.remove(DOMAIN_KEY)
    StorageUtil.remove(USER_DISABLED_KEY)
  }

  const checkIpStatus = async () => {
    const isUserDisabled = StorageUtil.get(USER_DISABLED_KEY) === 'true'
    if (isUserDisabled) {
      return
    }

    try {
      await getHealthCheck()
      clearIpRestriction()
      router.push('/')
    } catch {}
  }

  const generateDigitalRain = () => {
    const binaryCode = document.querySelector('.binary-code')
    if (binaryCode) {
      let code = ''
      for (let i = 0; i < 20; i++) {
        code += Math.floor(Math.random() * 2)
      }
      binaryCode.textContent = code
      setTimeout(generateDigitalRain, 500)
    }
  }

  onMounted(() => {
    const storedErrorCode = StorageUtil.get<string>(IP_ERROR_CODE_KEY)
    const errorMessage = StorageUtil.get<string>(IP_ERROR_MSG_KEY)
    const storedIp = StorageUtil.get<string>(IP_ADDRESS_KEY)
    const storedDomain = StorageUtil.get<string>(DOMAIN_KEY)

    errorCode.value = storedErrorCode || ''

    if (storedErrorCode === String(USER_ACCOUNT_DISABLED)) {
      restrictionType.value = 'USER'
    } else if (storedErrorCode === String(DOMAIN_NOT_IN_WHITELIST) || storedErrorCode === String(DOMAIN_IN_BLACKLIST)) {
      restrictionType.value = 'DOMAIN'
      domainName.value = storedDomain || $t('refuse.unknown')
    } else {
      restrictionType.value = 'IP'
      ipAddress.value = storedIp || $t('refuse.unknown')
    }

    if (storedErrorCode === String(IP_NOT_IN_WHITELIST)) {
      message.value = $t('refuse.messages.ipNotInWhitelist')
    } else if (storedErrorCode === String(IP_IN_BLACKLIST)) {
      message.value = $t('refuse.messages.ipInBlacklist')
    } else if (storedErrorCode === String(DOMAIN_NOT_IN_WHITELIST)) {
      message.value = $t('refuse.messages.domainNotInWhitelist')
    } else if (storedErrorCode === String(DOMAIN_IN_BLACKLIST)) {
      message.value = $t('refuse.messages.domainInBlacklist')
    } else if (storedErrorCode === String(USER_ACCOUNT_DISABLED)) {
      message.value = $t('refuse.messages.accountDisabled')
    } else if (errorMessage) {
      message.value = errorMessage
    }

    checkIpStatus()
    generateDigitalRain()
  })
</script>

<template>
  <div class="relative flex min-h-screen items-center justify-center overflow-hidden p-4">
    <CyberParticleBackground theme="default" :max-particles="80" />

    <div class="pointer-events-none absolute left-0 top-0 z-[1] h-full w-full overflow-hidden opacity-10">
      <div class="animate-pulse-slow border-brand-500/30 absolute left-10 top-10 h-48 w-48 rotate-45 border" />
      <div class="animate-float border-brand-500/20 absolute bottom-20 right-10 h-36 w-36 rotate-12 border" />
      <div class="animate-spin-slow border-brand-500/20 absolute right-[25%] top-[30%] h-24 w-24 -rotate-12 border-2" />
      <div class="via-brand-500/20 absolute left-0 top-1/3 h-0.5 w-full bg-gradient-to-r from-transparent to-transparent" />
      <div class="via-brand-500/20 absolute right-0 top-2/3 h-0.5 w-full bg-gradient-to-l from-transparent to-transparent" />

      <div class="digital-rain absolute inset-0 opacity-10" />

      <div class="circuit-lines absolute left-0 top-0 h-full w-full">
        <div class="circuit-1" />
        <div class="circuit-2" />
        <div class="circuit-3" />
      </div>
    </div>

    <div class="shake-subtle relative z-10 w-full max-w-2xl">
      <div
        class="cyberpunk-card bg-background-800/90 overflow-hidden rounded-xl border border-subtle shadow-[0_0_25px_rgba(var(--color-brand-500-rgb),_0.15)] backdrop-blur-lg"
      >
        <div
          class="flex items-center justify-between border-b border-default bg-gradient-to-r from-background-800 to-background-700 px-6 py-3"
        >
          <div class="flex items-center">
            <i class="fas fa-exclamation-triangle pulse mr-2 text-warning-400" />
            <span class="typewriter font-mono text-sm tracking-wider text-brand-400">{{ $t('refuse.protocolHeader') }}</span>
          </div>
          <div class="cyber-dots flex space-x-1">
            <div class="h-2 w-2 animate-pulse rounded-full bg-error-500" />
            <div class="h-2 w-2 animate-pulse rounded-full bg-warning-500 delay-100" />
            <div class="h-2 w-2 animate-pulse rounded-full bg-brand-500 delay-200" />
          </div>
        </div>

        <div class="relative p-8">
          <div class="binary-code text-brand-400/20 absolute right-4 top-4 font-mono text-xs" />

          <div class="mb-8 text-center">
            <div class="relative mb-4 inline-block">
              <div class="bg-error-500/10 absolute inset-0 animate-pulse rounded-full" />
              <div class="bg-error-500/5 absolute inset-0 animate-ping rounded-full" />
              <i class="fas fa-ban flicker-animation relative z-10 text-7xl text-error-400" />
            </div>
            <h1 class="glitch-text mb-2 text-3xl font-bold uppercase tracking-wider text-content-heading">
              {{ mainTitle }}
            </h1>
            <div class="from-brand-500/30 via-brand-500/60 to-brand-500/30 mx-auto h-1 w-32 bg-gradient-to-r" />
          </div>

          <div class="holo-card bg-background-700/50 relative mb-6 rounded-lg border border-default p-6">
            <div class="absolute right-0 top-0 rounded-bl-lg px-3 py-1 font-mono text-xs font-bold" :class="errorBadgeClass">
              {{ errorBadgeText }}
            </div>

            <p class="glow-text mb-4 text-lg font-medium text-content-heading">{{ message }}</p>
            <div class="space-y-2 text-sm text-content-muted">
              <p class="mb-2 flex items-center text-xs font-bold uppercase tracking-wider text-brand-400">
                <i class="fas fa-cog fa-spin-pulse mr-2" />{{ analysisTitle }}
              </p>
              <ul class="border-brand-500/30 list-none space-y-3 border-l-2 pl-4">
                <template v-if="restrictionType === 'USER'">
                  <li class="slide-in flex items-start" style="--delay: 0s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.accountDetails.disabled') }}</span>
                  </li>
                  <li class="slide-in flex items-start" style="--delay: 0.2s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.accountDetails.reason') }}</span>
                  </li>
                  <li class="slide-in flex items-start" style="--delay: 0.4s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.accountDetails.contact') }}</span>
                  </li>
                </template>

                <template v-else-if="restrictionType === 'IP'">
                  <li class="slide-in flex items-start" style="--delay: 0s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.ip.notInWhitelist') }}</span>
                  </li>
                  <li class="slide-in flex items-start" style="--delay: 0.2s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.ip.inBlacklist') }}</span>
                  </li>
                  <li class="slide-in flex items-start" style="--delay: 0.4s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.ip.contactAdmin') }}</span>
                  </li>
                </template>

                <template v-else>
                  <li class="slide-in flex items-start" style="--delay: 0s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.domain.notInWhitelist') }}</span>
                  </li>
                  <li class="slide-in flex items-start" style="--delay: 0.2s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.domain.inBlacklist') }}</span>
                  </li>
                  <li class="slide-in flex items-start" style="--delay: 0.4s">
                    <i class="fas fa-chevron-right mr-2 mt-1 text-xs text-brand-400" />
                    <span>{{ $t('refuse.domain.contactAdmin') }}</span>
                  </li>
                </template>
              </ul>
            </div>
          </div>

          <div
            class="cyber-terminal bg-background-700/30 flex flex-col items-center justify-between rounded-lg border border-default p-4 md:flex-row"
          >
            <div class="mb-4 flex items-center md:mb-0">
              <div class="mr-2 h-2 w-2 animate-ping rounded-full bg-brand-500" />
              <template v-if="restrictionType === 'USER'">
                <span class="text-sm text-content-muted">{{ $t('refuse.status.accountStatus') }} </span>
                <span class="typing-effect ml-2 font-mono text-error-400">{{ $t('refuse.status.disabled') }}</span>
              </template>
              <template v-else-if="restrictionType === 'IP'">
                <span class="text-sm text-content-muted">{{ $t('refuse.status.currentIp') }} </span>
                <span class="typing-effect ml-2 font-mono text-brand-400">{{ ipAddress }}</span>
              </template>
              <template v-else>
                <span class="text-sm text-content-muted">{{ $t('refuse.status.currentDomain') }} </span>
                <span class="typing-effect ml-2 font-mono text-brand-400">{{ domainName }}</span>
              </template>
            </div>
            <div class="text-sm text-content-disabled">
              <span>{{ $t('refuse.status.detectionTime') }}</span>
              <span class="scanning-text font-mono text-success-400">{{ formattedTimestamp }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .circuit-lines div {
    position: absolute;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.15), transparent);
    height: 1px;
  }

  .circuit-1 {
    top: 25%;
    left: 0;
    width: 100%;
    animation: circuit-move-1 8s var(--ease-in-out) infinite;
  }

  .circuit-2 {
    top: 45%;
    right: 0;
    width: 70%;
    animation: circuit-move-2 12s var(--ease-in-out) infinite;
  }

  .circuit-3 {
    top: 75%;
    left: 0;
    width: 85%;
    animation: circuit-move-3 15s var(--ease-in-out) infinite;
  }

  @keyframes circuit-move-1 {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  @keyframes circuit-move-2 {
    0% {
      transform: translateX(100%);
    }
    100% {
      transform: translateX(-100%);
    }
  }

  @keyframes circuit-move-3 {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  .cyberpunk-card {
    position: relative;
    transition: all var(--transition-normal) var(--ease-out);
  }

  .holo-card {
    position: relative;
    overflow: hidden;
  }

  .holo-card::after {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 50%;
    height: 100%;
    background: linear-gradient(to right, transparent 0%, rgba(var(--color-brand-500-rgb), 0.05) 50%, transparent 100%);
    transform: skewX(-25deg);
    animation: holo-shine 8s var(--ease-in-out) infinite;
  }

  @keyframes holo-shine {
    0% {
      left: -100%;
    }
    20% {
      left: 100%;
    }
    100% {
      left: 100%;
    }
  }

  .flicker-animation {
    animation: flicker 4s var(--ease-in-out) infinite;
  }

  @keyframes flicker {
    0% {
      opacity: 1;
    }
    92% {
      opacity: 1;
    }
    93% {
      opacity: var(--opacity-hover);
    }
    94% {
      opacity: 1;
    }
    95% {
      opacity: var(--opacity-disabled);
    }
    96% {
      opacity: 1;
    }
    97% {
      opacity: var(--opacity-loading);
    }
    98% {
      opacity: 1;
    }
    99% {
      opacity: var(--opacity-active);
    }
    100% {
      opacity: 1;
    }
  }

  .glow-text {
    text-shadow: var(--shadow-glow-sm);
  }

  .typing-effect {
    border-right: 2px solid rgba(var(--color-brand-500-rgb), 0.5);
    animation: cursor-blink 1s step-end infinite;
  }

  @keyframes cursor-blink {
    0%,
    100% {
      border-color: rgba(var(--color-brand-500-rgb), 0.5);
    }
    50% {
      border-color: transparent;
    }
  }

  .scanning-text {
    position: relative;
  }

  .scanning-text::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: var(--space-xs);
    background-color: rgba(var(--color-success-rgb), 0.5);
    animation: scan 3s var(--ease-in-out) infinite;
  }

  @keyframes scan {
    0% {
      left: 0;
    }
    100% {
      left: 100%;
    }
  }

  .pulse {
    animation: pulse-anim 2s var(--ease-in-out) infinite;
  }

  @keyframes pulse-anim {
    0% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.2);
      opacity: var(--opacity-loading);
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  .shake-subtle {
    animation: shake 60s var(--ease-in-out) infinite;
  }

  @keyframes shake {
    0%,
    100% {
      transform: translateY(0);
    }
    10% {
      transform: translateY(-1px) rotate(-0.5deg);
    }
    20% {
      transform: translateY(0.5px) rotate(0.5deg);
    }
    30% {
      transform: translateY(-0.5px) rotate(-0.5deg);
    }
    40% {
      transform: translateY(0.5px) rotate(0.5deg);
    }
    50% {
      transform: translateY(-1px) rotate(-0.5deg);
    }
    60% {
      transform: translateY(1px) rotate(0.5deg);
    }
    70% {
      transform: translateY(-0.5px) rotate(-0.5deg);
    }
    80% {
      transform: translateY(0.5px) rotate(0.5deg);
    }
    90% {
      transform: translateY(-0.5px) rotate(-0.5deg);
    }
  }

  .animate-pulse-slow {
    animation: pulse 6s var(--ease-in-out) infinite;
  }

  .animate-spin-slow {
    animation: spin 8s linear infinite;
  }

  .animate-float {
    animation: float 10s var(--ease-in-out) infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: var(--opacity-hover);
    }
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0) rotate(12deg);
    }
    50% {
      transform: translateY(-20px) rotate(12deg);
    }
  }

  .slide-in {
    animation: slide-from-left 0.5s var(--ease-out) forwards;
    animation-delay: var(--delay, 0s);
    opacity: 0;
    transform: translateX(-20px);
  }

  @keyframes slide-from-left {
    0% {
      opacity: 0;
      transform: translateX(-20px);
    }
    100% {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .cyber-terminal {
    box-shadow: inset 0 0 var(--space-sm) rgba(var(--color-brand-500-rgb), 0.1);
    position: relative;
  }

  .glitch-text {
    position: relative;
    text-shadow: var(--shadow-glow-sm);
  }

  .delay-100 {
    animation-delay: var(--transition-fast);
  }

  .delay-200 {
    animation-delay: var(--transition-normal);
  }

  .digital-rain {
    background: linear-gradient(180deg, transparent 0%, rgba(var(--color-success-rgb), 0.05) 50%, transparent 100%);
    background-size: 100% 200%;
    animation: rain-fall 10s linear infinite;
  }

  @keyframes rain-fall {
    0% {
      background-position: 0 -100%;
    }
    100% {
      background-position: 0 100%;
    }
  }
</style>
