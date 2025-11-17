<script setup lang="ts">
  import MessageList from '@/components/Message/MessageList.vue'

  defineOptions({
    name: 'MessagesPage',
  })
</script>

<template>
  <div class="messages-page">
    <MessageList />
  </div>
</template>

<style scoped>
  .messages-page {
    @apply relative min-h-screen py-6;
  }

  @media (max-width: 768px) {
    .messages-page {
      @apply py-6;
    }
  }
</style>
