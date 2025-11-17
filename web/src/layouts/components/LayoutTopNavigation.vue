<script setup lang="ts">
  import { useRouter } from 'vue-router'

  const router = useRouter()

  const redirectToAuth = () => {
    router.push('/auth')
  }
</script>

<template>
  <cyberNavbar @login="redirectToAuth" />
</template>
