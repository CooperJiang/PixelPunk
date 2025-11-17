<script setup lang="ts">
  import { ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    processing: boolean
    selectedCount: number
    paused?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'trigger', maxImages: number): void
    (e: 'reset-stuck', timeThreshold: number): void
    (e: 'retry-all-failed'): void
    (e: 'ignore-selected', reason: string): void
    (e: 'unignore-selected'): void
  }>()

  const maxImages = ref(200)
  const timeThreshold = ref(30)
  const showConfirmDialog = ref(false)
  const confirmMessage = ref('')
  const actionType = ref('')
  const ignoreReason = ref('')
  const showAdvanced = ref(false)

  const triggerTagging = () => {
    emit('trigger', maxImages.value)
  }

  const resetStuckTask = () => {
    if (timeThreshold.value < 5) {
      timeThreshold.value = 5
    }

    confirmMessage.value = $t('admin.tagging.confirm.resetStuck').replace('{minutes}', String(timeThreshold.value))
    actionType.value = 'reset-stuck'
    showConfirmDialog.value = true
  }

  const confirmRetryFailed = () => {
    confirmMessage.value = $t('admin.tagging.confirm.retryFailed')
    actionType.value = 'retry-all-failed'
    showConfirmDialog.value = true
  }

  const confirmIgnore = () => {
    confirmMessage.value = $t('admin.tagging.confirm.ignore').replace('{count}', String(props.selectedCount))
    actionType.value = 'ignore-selected'
    ignoreReason.value = ''
    showConfirmDialog.value = true
  }

  const confirmUnignore = () => {
    confirmMessage.value = $t('admin.tagging.confirm.unignore').replace('{count}', String(props.selectedCount))
    actionType.value = 'unignore-selected'
    showConfirmDialog.value = true
  }

  const confirmAction = () => {
    showConfirmDialog.value = false
    if (actionType.value === 'reset-stuck') {
      emit('reset-stuck', timeThreshold.value)
    } else if (actionType.value === 'retry-all-failed') {
      emit('retry-all-failed')
    } else if (actionType.value === 'ignore-selected') {
      emit('ignore-selected', ignoreReason.value)
    } else if (actionType.value === 'unignore-selected') {
      emit('unignore-selected')
    }
  }
</script>

<template>
  <div class="tagging-action-panel mb-4 rounded-lg border border-subtle bg-background-600 p-4">
    <div class="mb-3 flex items-center justify-between">
      <div class="text-content-content-disabled text-xs">{{ $t('admin.tagging.hints.advancedCollapsed') }}</div>
      <CyberButton
        type="secondary"
        class="text-xs"
        :icon="showAdvanced ? 'chevron-up' : 'chevron-down'"
        @click="showAdvanced = !showAdvanced"
      >
        {{ showAdvanced ? $t('admin.tagging.actions.hideAdvanced') : $t('admin.tagging.actions.showAdvanced') }}
      </CyberButton>
    </div>
    <div class="flex flex-wrap gap-3">
      <div class="action-item min-w-48 flex-1" v-show="showAdvanced">
        <h3 class="mb-2 text-sm text-content"><i class="fas fa-tools mr-1.5" />{{ $t('admin.tagging.actions.showAdvanced') }}</h3>
        <div class="flex items-center gap-2">
          <CyberInput
            v-model="maxImages"
            type="number"
            :min="0"
            :max="500"
            class="w-24"
            :placeholder="$t('admin.tagging.placeholders.maxImages')"
          />
          <CyberTooltip :content="$t('admin.tagging.hints.triggerTooltip')" placement="top" theme="cyber">
            <CyberButton type="warning" icon="play" :disabled="processing || paused" @click="triggerTagging">
              {{ $t('admin.tagging.actions.trigger') }}
            </CyberButton>
          </CyberTooltip>
        </div>
        <div class="mt-1 text-xs" :class="paused ? 'text-error-500' : 'text-content-content-disabled'">
          <template v-if="paused">{{ $t('admin.tagging.hints.queuePaused') }}</template>
          <template v-else>{{
            maxImages
              ? $t('admin.tagging.hints.maxImages').replace('{count}', String(maxImages))
              : $t('admin.tagging.hints.maxImagesDefault')
          }}</template>
        </div>
      </div>

      <div class="action-item min-w-48 flex-1" v-show="showAdvanced">
        <h3 class="mb-2 text-sm text-content">
          <i class="fas fa-sync-alt mr-1.5" />{{ $t('admin.tagging.actions.resetStuck') }}
        </h3>
        <div class="flex items-center gap-2">
          <CyberInput
            v-model="timeThreshold"
            type="number"
            :min="5"
            :max="240"
            class="w-24"
            :placeholder="$t('admin.tagging.placeholders.minutes')"
          />
          <CyberTooltip :content="$t('admin.tagging.hints.triggerTooltip')" placement="top" theme="cyber">
            <CyberButton type="warning" icon="sync-alt" :disabled="processing || paused" @click="resetStuckTask">
              {{ $t('admin.tagging.actions.resetStuck') }}
            </CyberButton>
          </CyberTooltip>
        </div>
        <div class="text-content-content-disabled mt-1 text-xs">
          {{ $t('admin.tagging.hints.maxImages').replace('{count}', String(timeThreshold)) }}
        </div>
      </div>

      <div class="action-item min-w-48 flex-1" v-show="showAdvanced">
        <h3 class="mb-2 text-sm text-content"><i class="fas fa-redo mr-1.5" />{{ $t('admin.tagging.actions.retryFailed') }}</h3>
        <div class="flex items-center">
          <CyberTooltip :content="$t('admin.tagging.hints.triggerTooltip')" placement="top" theme="cyber">
            <CyberButton type="danger" icon="exclamation-triangle" :disabled="processing || paused" @click="confirmRetryFailed">
              {{ $t('admin.tagging.actions.retryFailed') }}
            </CyberButton>
          </CyberTooltip>
        </div>
        <div class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.tagging.hints.advancedCollapsed') }}</div>
      </div>

      <div class="action-item min-w-[220px] flex-1" v-show="selectedCount > 0">
        <h3 class="mb-2 text-sm text-content"><i class="fas fa-ban mr-1.5" />{{ $t('admin.tagging.actions.ignore') }}</h3>
        <div class="flex items-center gap-2">
          <CyberButton type="warning" icon="ban" :disabled="processing || paused || selectedCount <= 0" @click="confirmIgnore">
            {{ $t('admin.tagging.actions.ignore') }}
          </CyberButton>
        </div>
        <div class="text-content-content-disabled mt-1 text-xs">
          {{ $t('admin.tagging.confirm.ignore').replace('{count}', String(selectedCount)) }}
        </div>
      </div>

      <div class="action-item min-w-[220px] flex-1" v-show="selectedCount > 0">
        <h3 class="mb-2 text-sm text-content"><i class="fas fa-undo mr-1.5" />{{ $t('admin.tagging.actions.unignore') }}</h3>
        <div class="flex items-center gap-2">
          <CyberButton
            type="secondary"
            icon="undo"
            :disabled="processing || paused || selectedCount <= 0"
            @click="confirmUnignore"
          >
            {{ $t('admin.tagging.actions.unignore') }}
          </CyberButton>
        </div>
        <div class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.tagging.statusDesc.ignored') }}</div>
      </div>
    </div>
  </div>

  <teleport to="body">
    <div v-if="showConfirmDialog" class="bg-overlay fixed inset-0 z-50 flex items-center justify-center">
      <CyberDialog
        v-model="showConfirmDialog"
        width="460px"
        :title="$t('admin.tagging.confirm.resetStuck').split('?')[0]"
        :show-close="false"
      >
        <template #icon>
          <i class="fas fa-exclamation-triangle text-yellow-500" />
        </template>
        <div class="space-y-3 text-sm text-content-heading">
          <div>{{ confirmMessage }}</div>
          <div v-if="actionType === 'ignore-selected'">
            <label class="mb-1 block text-xs text-content-muted">{{ $t('admin.tagging.placeholders.reason') }}</label>
            <textarea
              v-model="ignoreReason"
              rows="3"
              class="w-full rounded border border-subtle bg-input-bg p-2 text-sm text-content-heading"
              :placeholder="$t('admin.tagging.placeholders.reason')"
            />
          </div>
        </div>
        <template #footer>
          <div class="flex justify-end space-x-3 p-4">
            <CyberButton type="text" @click="showConfirmDialog = false"> {{ $t('admin.channels.actions.cancel') }} </CyberButton>
            <CyberButton type="danger" icon="check" @click="confirmAction">{{ $t('admin.channels.actions.save') }}</CyberButton>
          </div>
        </template>
      </CyberDialog>
    </div>
  </teleport>
</template>

<style scoped lang="scss">
  .tagging-action-panel {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
  }
</style>
