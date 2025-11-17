<template>
  <main class="admin-layout-content" :class="{ 'sidebar-layout': sidebarLayout }">
    <template v-if="sidebarLayout">
      <aside v-if="$slots.sidebar" class="admin-layout-sidebar">
        <slot name="sidebar" />
      </aside>

      <div class="sidebar-content">
        <div v-if="$slots.toolbar" class="content-toolbar">
          <slot name="toolbar" />
        </div>

        <div class="main-content" :class="{ compact: compact }">
          <slot name="default">
            <slot name="content" />
          </slot>
        </div>
      </div>
    </template>

    <template v-else>
      <div v-if="$slots.toolbar" class="content-toolbar">
        <slot name="toolbar" />
      </div>

      <div class="main-content" :class="{ compact: compact }">
        <slot name="default">
          <slot name="content" />
        </slot>
      </div>
    </template>
  </main>
</template>

<script setup lang="ts">
  interface AdminWrapperContentProps {
    sidebarLayout?: boolean
    contentPadding?: string
    compact?: boolean
  }

  const { sidebarLayout = false, contentPadding = '1.5rem', compact = false } = defineProps<AdminWrapperContentProps>()
</script>

<style scoped>
  .admin-layout-content {
    @apply relative flex flex-col bg-background-900;
    flex: 1;
    min-height: 0;
    height: 0;

    overflow-y: auto;
    overflow-x: hidden;
    scrollbar-gutter: stable;
    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.5) transparent;
  }

  .content-toolbar {
    @apply shrink-0 px-6 pt-4;
  }

  .main-content {
    @apply flex flex-1 flex-col;
    padding: v-bind(contentPadding);
    min-height: 0;
    height: 100%;
  }

  .admin-layout-content.sidebar-layout {
    @apply flex min-h-0 flex-row;

    overflow: hidden !important;
  }

  .admin-layout-sidebar {
    @apply relative w-48 shrink-0 overflow-y-auto overflow-x-hidden;
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.03) 0%,
      rgba(var(--color-error-rgb), 0.02) 50%,
      rgba(var(--color-error-rgb), 0.01) 100%
    );
    border-right: 1px solid var(--color-border-subtle);
  }

  .admin-layout-sidebar::-webkit-scrollbar {
    @apply w-1.5 bg-transparent;
  }

  .admin-layout-sidebar::-webkit-scrollbar-track {
    @apply rounded-sm bg-[rgba(var(--color-brand-500-rgb),0.1)];
  }

  .admin-layout-sidebar::-webkit-scrollbar-thumb {
    @apply rounded-sm bg-[var(--color-border-default)] transition-all duration-300;
  }

  .admin-layout-sidebar::-webkit-scrollbar-thumb:hover {
    @apply bg-[var(--color-brand-500)];
  }

  .sidebar-content {
    @apply relative flex flex-1 flex-col overflow-hidden;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.85) 100%
    );
    min-height: 0;
  }

  .sidebar-content .main-content {
    padding: v-bind(contentPadding);

    overflow-y: auto;
    overflow-x: hidden;
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .main-content.compact {
    @apply px-6 py-5;
  }

  @media (max-width: 768px) {
    .content-toolbar {
      @apply px-4 pt-3.5;
    }

    .main-content {
      @apply px-4 py-5;
    }
  }

  @media (max-width: 480px) {
    .main-content {
      @apply px-3.5 py-4;
    }
  }
</style>
