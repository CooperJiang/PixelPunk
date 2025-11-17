<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed } from 'vue'
  import { formatFileSize } from '@/utils/formatting/format'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    modelValue: boolean
    file: any
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
  }>()

  const toast = useToast()
  const { $t } = useTexts()

  const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val),
  })

  /* 格式化日期 */
  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    try {
      return new Date(dateStr).toLocaleString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      })
    } catch {
      return dateStr
    }
  }

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        toast.success($t('components.fileDetailDrawer.messages.copySuccess', { label }))
      })
      .catch(() => {
        toast.error($t('components.fileDetailDrawer.messages.copyFailed'))
      })
  }

  const getAccessLevelText = (level: string) => {
    const map: Record<string, string> = {
      public: $t('components.fileDetailDrawer.accessLevel.public'),
      private: $t('components.fileDetailDrawer.accessLevel.private'),
      password: $t('components.fileDetailDrawer.accessLevel.password'),
    }
    return map[level] || level
  }

  const getStorageDurationText = (duration: string) => {
    const map: Record<string, string> = {
      permanent: $t('components.fileDetailDrawer.storageDuration.permanent'),
      temporary: $t('components.fileDetailDrawer.storageDuration.temporary'),
    }
    return map[duration] || duration
  }

  const hasEXIFData = (exif: any) => {
    if (!exif) return false
    return !!(
      exif.make ||
      exif.model ||
      exif.lens_model ||
      exif.f_number ||
      exif.exposure_time ||
      exif.iso ||
      exif.focal_length ||
      exif.gps_latitude ||
      exif.gps_longitude ||
      exif.date_time_original
    )
  }

  const getCameraInfo = (exif: any) => {
    if (!exif) return ''
    if (exif.make && exif.model) {
      if (exif.model.toLowerCase().startsWith(exif.make.toLowerCase())) {
        return exif.model
      }
      return `${exif.make} ${exif.model}`
    }
    return exif.make || exif.model || ''
  }

  const copyGPS = (exif: any) => {
    if (!exif.gps_latitude || !exif.gps_longitude) return
    const coords = `${exif.gps_latitude.toFixed(6)}, ${exif.gps_longitude.toFixed(6)}`
    copyToClipboard(coords, $t('components.fileDetailDrawer.labels.gpsCoords'))
  }
</script>

<template>
  <CyberDrawer v-model="visible" width="420px" :mask-closable="true" :show-header="false" :show-keyboard-tip="false">
    <div v-if="file" class="detail-content">
      <div class="detail-section">
        <div class="section-title">
          <i class="fas fa-file-alt" />
          <span>{{ $t('components.fileDetailDrawer.sections.basicInfo') }}</span>
        </div>
        <div class="detail-items">
          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.fileName') }}</div>
            <div class="item-value text-truncate" :title="file.original_name">
              {{ file.original_name || '-' }}
            </div>
            <button
              v-if="file.original_name"
              class="copy-btn"
              @click="copyToClipboard(file.original_name, $t('components.fileDetailDrawer.labels.fileName'))"
            >
              <i class="fas fa-copy" />
            </button>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.fileSize') }}</div>
            <div class="item-value">{{ formatFileSize(file.size) }}</div>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.resolution') }}</div>
            <div class="item-value">
              <span class="resolution-badge">{{ file.width }} × {{ file.height }}</span>
            </div>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.format') }}</div>
            <div class="item-value">
              <span class="format-badge">{{ file.format?.toUpperCase() || '-' }}</span>
            </div>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.accessLevel') }}</div>
            <div class="item-value">
              <span class="access-badge" :class="file.access_level">
                {{ getAccessLevelText(file.access_level) }}
              </span>
            </div>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.storageDuration') }}</div>
            <div class="item-value">{{ getStorageDurationText(file.storage_duration) }}</div>
          </div>

          <div v-if="file.expires_at" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.expiresAt') }}</div>
            <div class="item-value">{{ formatDate(file.expires_at) }}</div>
          </div>
        </div>
      </div>

      <div v-if="file.ai_info" class="detail-section">
        <div class="section-title">
          <i class="fas fa-robot" />
          <span>{{ $t('components.fileDetailDrawer.sections.aiAnalysis') }}</span>
        </div>
        <div class="detail-items">
          <div v-if="file.ai_info.description" class="detail-item column">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.description') }}</div>
            <div class="item-value description">{{ file.ai_info.description }}</div>
          </div>

          <div v-if="file.ai_info.tags && file.ai_info.tags.length" class="detail-item column">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.tags') }}</div>
            <div class="tags-container">
              <span v-for="(tag, index) in file.ai_info.tags" :key="index" class="tag-badge">
                {{ tag }}
              </span>
            </div>
          </div>

          <div v-if="file.ai_info.dominant_color" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.dominantColor') }}</div>
            <div class="item-value">
              <div class="color-preview">
                <div class="color-block" :style="{ background: file.ai_info.dominant_color }" />
                <span class="color-text">{{ file.ai_info.dominant_color }}</span>
              </div>
            </div>
            <button
              class="copy-btn"
              @click="copyToClipboard(file.ai_info.dominant_color, $t('components.fileDetailDrawer.labels.colorCode'))"
            >
              <i class="fas fa-copy" />
            </button>
          </div>

          <div v-if="file.ai_info.resolution" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.resolutionLevel') }}</div>
            <div class="item-value">
              <span class="resolution-badge">{{ file.ai_info.resolution }}</span>
            </div>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.contentSafety') }}</div>
            <div class="item-value">
              <span class="nsfw-badge" :class="{ safe: !file.ai_info.is_nsfw, nsfw: file.ai_info.is_nsfw }">
                {{ file.ai_info.nsfw_evaluation || $t('components.fileDetailDrawer.labels.notEvaluated') }}
              </span>
            </div>
          </div>

          <div v-if="file.ai_info.nsfw_score !== undefined" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.nsfwScore') }}</div>
            <div class="item-value">{{ (file.ai_info.nsfw_score * 100).toFixed(1) }}%</div>
          </div>
        </div>
      </div>

      <div class="detail-section">
        <div class="section-title">
          <i class="fas fa-link" />
          <span>{{ $t('components.fileDetailDrawer.sections.links') }}</span>
        </div>
        <div class="detail-items">
          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.fullUrl') }}</div>
            <div class="item-value text-truncate font-mono text-sm" :title="file.full_url">
              {{ file.full_url || '-' }}
            </div>
            <button
              v-if="file.full_url"
              class="copy-btn"
              @click="copyToClipboard(file.full_url, $t('components.fileDetailDrawer.labels.fullUrl'))"
            >
              <i class="fas fa-copy" />
            </button>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.shortUrl') }}</div>
            <div class="item-value text-truncate font-mono text-sm" :title="file.short_url">
              {{ file.short_url || '-' }}
            </div>
            <button
              v-if="file.short_url"
              class="copy-btn"
              @click="copyToClipboard(file.short_url, $t('components.fileDetailDrawer.labels.shortUrl'))"
            >
              <i class="fas fa-copy" />
            </button>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.thumbUrl') }}</div>
            <div class="item-value text-truncate font-mono text-sm" :title="file.full_thumb_url">
              {{ file.full_thumb_url || '-' }}
            </div>
            <button
              v-if="file.full_thumb_url"
              class="copy-btn"
              @click="copyToClipboard(file.full_thumb_url, $t('components.fileDetailDrawer.labels.thumbUrl'))"
            >
              <i class="fas fa-copy" />
            </button>
          </div>
        </div>
      </div>

      <div v-if="file.exif_info && hasEXIFData(file.exif_info)" class="detail-section">
        <div class="section-title">
          <i class="fas fa-camera" />
          <span>{{ $t('components.fileDetailDrawer.sections.shootingInfo') }}</span>
        </div>
        <div class="detail-items">
          <div v-if="file.exif_info.make || file.exif_info.model" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.camera') }}</div>
            <div class="item-value">{{ getCameraInfo(file.exif_info) }}</div>
          </div>

          <div v-if="file.exif_info.lens_model" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.lens') }}</div>
            <div class="item-value">{{ file.exif_info.lens_model }}</div>
          </div>

          <div v-if="file.exif_info.focal_length" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.focalLength') }}</div>
            <div class="item-value">
              {{ file.exif_info.focal_length }}mm
              <span v-if="file.exif_info.focal_length_in_35mm" class="text-secondary">
                ({{ file.exif_info.focal_length_in_35mm }}mm {{ $t('components.fileDetailDrawer.labels.equivalent') }})
              </span>
            </div>
          </div>

          <div v-if="file.exif_info.f_number" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.aperture') }}</div>
            <div class="item-value">f/{{ file.exif_info.f_number.toFixed(1) }}</div>
          </div>

          <div v-if="file.exif_info.exposure_time" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.shutterSpeed') }}</div>
            <div class="item-value">{{ file.exif_info.exposure_time }}s</div>
          </div>

          <div v-if="file.exif_info.iso" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.iso') }}</div>
            <div class="item-value">ISO {{ file.exif_info.iso }}</div>
          </div>

          <div v-if="file.exif_info.date_time_original" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.shootingTime') }}</div>
            <div class="item-value">{{ formatDate(file.exif_info.date_time_original) }}</div>
          </div>

          <div v-if="file.exif_info.gps_latitude && file.exif_info.gps_longitude" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.gpsCoords') }}</div>
            <div class="item-value font-mono text-sm">
              {{ file.exif_info.gps_latitude.toFixed(6) }}, {{ file.exif_info.gps_longitude.toFixed(6) }}
            </div>
            <button class="copy-btn" @click="copyGPS(file.exif_info)">
              <i class="fas fa-copy" />
            </button>
          </div>

          <div v-if="file.exif_info.gps_altitude" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.altitude') }}</div>
            <div class="item-value">{{ file.exif_info.gps_altitude.toFixed(1) }}m</div>
          </div>

          <div v-if="file.exif_info.software" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.editingSoftware') }}</div>
            <div class="item-value">{{ file.exif_info.software }}</div>
          </div>

          <div v-if="file.exif_info.copyright" class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.copyright') }}</div>
            <div class="item-value">{{ file.exif_info.copyright }}</div>
          </div>
        </div>
      </div>

      <div class="detail-section">
        <div class="section-title">
          <i class="fas fa-clock" />
          <span>{{ $t('components.fileDetailDrawer.sections.time') }}</span>
        </div>
        <div class="detail-items">
          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.createdAt') }}</div>
            <div class="item-value">{{ formatDate(file.created_at) }}</div>
          </div>

          <div class="detail-item">
            <div class="item-label">{{ $t('components.fileDetailDrawer.labels.updatedAt') }}</div>
            <div class="item-value">{{ formatDate(file.updated_at) }}</div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <i class="fas fa-file-alt" />
      <p>{{ $t('components.fileDetailDrawer.empty') }}</p>
    </div>
  </CyberDrawer>
</template>

<style scoped lang="scss">
  .detail-content {
    padding: 16px;
    background: linear-gradient(
      180deg,
      rgba(var(--color-background-900-rgb), 0.3) 0%,
      rgba(var(--color-background-900-rgb), 0.5) 100%
    );
  }

  .detail-section {
    margin-bottom: 20px;
    padding: 14px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.6) 0%,
      rgba(var(--color-background-700-rgb), 0.4) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    backdrop-filter: blur(8px);
    position: relative;

    &:last-child {
      margin-bottom: 0;
    }

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 20px;
      right: 20px;
      height: 1px;
      background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.5) 50%, transparent);
    }
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
    padding-bottom: 10px;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    font-size: 13px;
    font-weight: 700;
    color: var(--color-brand-300);
    letter-spacing: 0.05em;
    text-transform: uppercase;

    i {
      font-size: 14px;
      color: var(--color-brand-400);
      filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.5));
    }
  }

  .detail-items {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .detail-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 10px;
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;

    &:hover {
      background: rgba(var(--color-background-900-rgb), 0.5);
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      transform: translateX(-2px);
      box-shadow: 2px 0 8px rgba(var(--color-brand-500-rgb), 0.15);
    }

    &.column {
      flex-direction: column;
      align-items: stretch;
      gap: 8px;
    }
  }

  .item-label {
    flex-shrink: 0;
    width: 80px;
    font-size: 11px;
    color: rgba(var(--color-brand-500-rgb), 0.7);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.03em;
  }

  .item-value {
    flex: 1;
    font-size: 12px;
    color: var(--color-content-default);
    word-break: break-word;
    font-weight: 500;

    &.description {
      line-height: 1.6;
      margin-top: 4px;
      color: rgba(var(--color-content-rgb), 0.85);
    }

    &.text-truncate {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    &.font-mono {
      font-family: 'SF Mono', Monaco, Consolas, monospace;
      font-size: 11px;
      padding: 3px 6px;
      background: rgba(var(--color-background-900-rgb), 0.5);
      border-radius: var(--radius-sm);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.12);
    }

    &.text-sm {
      font-size: 10px;
    }
  }

  .copy-btn {
    flex-shrink: 0;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    background: rgba(var(--color-background-800-rgb), 0.4);
    color: rgba(var(--color-content-rgb), 0.5);
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 11px;

    &:hover {
      border-color: rgba(var(--color-brand-500-rgb), 0.5);
      background: rgba(var(--color-brand-500-rgb), 0.15);
      color: var(--color-brand-300);
      box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.2);
    }

    &:active {
      transform: scale(0.95);
    }
  }

  .format-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 8px;
    font-size: 10px;
    font-weight: 700;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.1));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-300);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    box-shadow: 0 1px 4px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .access-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 8px;
    font-size: 10px;
    font-weight: 700;
    border-radius: var(--radius-sm);
    text-transform: uppercase;
    letter-spacing: 0.03em;

    &.public {
      background: rgba(var(--color-success-rgb), 0.15);
      border: 1px solid rgba(var(--color-success-rgb), 0.3);
      color: var(--color-success-300);
      box-shadow: 0 1px 4px rgba(var(--color-success-rgb), 0.2);
    }

    &.private {
      background: rgba(var(--color-warning-rgb), 0.15);
      border: 1px solid rgba(var(--color-warning-rgb), 0.3);
      color: var(--color-warning-300);
      box-shadow: 0 1px 4px rgba(var(--color-warning-rgb), 0.2);
    }

    &.password {
      background: rgba(var(--color-error-rgb), 0.15);
      border: 1px solid rgba(var(--color-error-rgb), 0.3);
      color: var(--color-error-300);
      box-shadow: 0 1px 4px rgba(var(--color-error-rgb), 0.2);
    }
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-top: 6px;
  }

  .tag-badge {
    display: inline-flex;
    align-items: center;
    padding: 4px 8px;
    font-size: 10px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    color: var(--color-brand-300);
    transition: all 0.2s ease;

    &:hover {
      background: rgba(var(--color-brand-500-rgb), 0.25);
      border-color: rgba(var(--color-brand-500-rgb), 0.4);
      transform: translateY(-1px);
      box-shadow: 0 2px 6px rgba(var(--color-brand-500-rgb), 0.25);
    }
  }

  .color-preview {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .color-block {
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  }

  .color-text {
    font-family: 'SF Mono', Monaco, Consolas, monospace;
    font-size: 11px;
    font-weight: 600;
    color: var(--color-brand-300);
    padding: 2px 6px;
    background: rgba(var(--color-background-900-rgb), 0.5);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    letter-spacing: 0.03em;
  }

  .resolution-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 8px;
    font-size: 10px;
    font-weight: 700;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-300);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    box-shadow: 0 1px 4px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .nsfw-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 8px;
    font-size: 10px;
    font-weight: 700;
    border-radius: var(--radius-sm);
    text-transform: uppercase;
    letter-spacing: 0.03em;

    &.safe {
      background: rgba(var(--color-success-rgb), 0.15);
      border: 1px solid rgba(var(--color-success-rgb), 0.3);
      color: var(--color-success-300);
      box-shadow: 0 1px 4px rgba(var(--color-success-rgb), 0.2);
    }

    &.nsfw {
      background: rgba(var(--color-error-rgb), 0.15);
      border: 1px solid rgba(var(--color-error-rgb), 0.3);
      color: var(--color-error-300);
      box-shadow: 0 1px 4px rgba(var(--color-error-rgb), 0.2);
    }
  }

  .text-secondary {
    opacity: 0.6;
    font-size: 11px;
    margin-left: 4px;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 20px;
    color: var(--color-content-muted);
    position: relative;

    &::before {
      content: '';
      position: absolute;
      inset: 0;
      background: repeating-linear-gradient(
        45deg,
        transparent,
        transparent 20px,
        rgba(var(--color-brand-500-rgb), 0.02) 20px,
        rgba(var(--color-brand-500-rgb), 0.02) 40px
      );
    }

    i {
      font-size: 64px;
      margin-bottom: 20px;
      color: var(--color-brand-500);
      opacity: 0.3;
      filter: drop-shadow(0 0 16px rgba(var(--color-brand-500-rgb), 0.4));
      animation: float 3s ease-in-out infinite;
    }

    p {
      font-size: 15px;
      font-weight: 500;
      margin: 0;
      color: rgba(var(--color-content-rgb), 0.6);
    }
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-10px);
    }
  }
</style>
