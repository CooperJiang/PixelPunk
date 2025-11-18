<script setup lang="ts">
  import { ref, onMounted, computed, watch, reactive, nextTick } from 'vue'
  import type { CreateStorageChannelRequest } from '@/api/types/index'
  import { getSupportedTypes, getConfigTemplates, type StorageConfigTemplate } from '@/api/storage'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    form: CreateStorageChannelRequest & Record<string, unknown>
    isEdit: boolean
  }>()

  const emit = defineEmits<{
    'update:form': [value: CreateStorageChannelRequest & Record<string, unknown>]
  }>()

  const localForm = reactive<CreateStorageChannelRequest & Record<string, unknown>>({ ...props.form })

  /* 监听 props.form 变化，同步到本地 */
  watch(
    () => props.form,
    (newForm) => {
      Object.assign(localForm, newForm)
    },
    { deep: true }
  )

  /* 监听本地表单变化，emit 更新 */
  watch(
    localForm,
    (newForm) => {
      emit('update:form', { ...newForm })
    },
    { deep: true }
  )

  const channelTypeOptions = ref<{ label: string; value: string }[]>([])
  const mapTypeLabel = (t: string) => {
    const m: Record<string, string> = {
      local: $t('admin.channels.types.local'),
      oss: $t('admin.channels.types.oss'),
      cos: $t('admin.channels.types.cos'),
      qiniu: $t('admin.channels.types.qiniu'),
      upyun: $t('admin.channels.types.upyun'),
      s3: $t('admin.channels.types.s3'),
      rainyun: $t('admin.channels.types.rainyun'),
      webdav: 'WebDAV',
      r2: 'Cloudflare R2',
      azureblob: 'Azure Blob',
      sftp: 'SFTP',
      ftp: 'FTP/FTPS',
    }
    return m[t] || t
  }

  const ranking = ['oss', 'cos', 'qiniu', 'upyun', 's3', 'rainyun', 'webdav', 'r2', 'azureblob', 'local', 'sftp', 'ftp']
  const applyRanking = (opts: { label: string; value: string }[]) => {
    const order = new Map<string, number>()
    ranking.forEach((v, i) => order.set(v, i))
    return opts.slice().sort((a, b) => {
      const ia = order.has(a.value) ? (order.get(a.value) as number) : Number.MAX_SAFE_INTEGER
      const ib = order.has(b.value) ? (order.get(b.value) as number) : Number.MAX_SAFE_INTEGER
      if (ia !== ib) return ia - ib
      return a.label.localeCompare(b.label, 'zh-Hans-CN')
    })
  }
  onMounted(async () => {
    try {
      const res = await getSupportedTypes()
      const resData = res as Record<string, unknown>

      const isSuccess = resData.success === true || resData.code === 200
      const dataArray = resData.data

      if (isSuccess && Array.isArray(dataArray)) {
        const arr = dataArray as unknown[]
        if (arr.length > 0 && typeof arr[0] === 'object' && arr[0] && 'value' in arr[0]) {
          const opts = arr.map((it: unknown) => {
            const item = it as Record<string, string>
            return { label: item.label || mapTypeLabel(item.value), value: item.value }
          })
          channelTypeOptions.value = applyRanking(opts)
        } else {
          const opts = (arr as string[]).map((t: string) => ({ label: mapTypeLabel(t), value: t }))
          channelTypeOptions.value = applyRanking(opts)
        }
      } else if (Array.isArray(res)) {
        const opts = (res as string[]).map((t: string) => ({ label: mapTypeLabel(t), value: t }))
        channelTypeOptions.value = applyRanking(opts)
      }
    } catch (_error) {
      channelTypeOptions.value = []
    }
  })

  const isEdit = computed(() => props.isEdit)

  const dynamicFields = ref<StorageConfigTemplate[]>([])
  const fieldTips: Record<string, string> = {
    use_path_style: $t('admin.channels.tips.usePathStyle'),
  }
  const getFieldTip = (key: string) => fieldTips[key] || ''
  watch(
    () => localForm.type,
    async (newType) => {
      if (!newType || String(newType).trim() === '') {
        dynamicFields.value = []
        return
      }
      try {
        const res = await getConfigTemplates(newType)
        const resData = res as Record<string, unknown>

        const isSuccess = resData.success === true || resData.code === 200

        if (isSuccess) {
          const data = (resData.data || []) as StorageConfigTemplate[]

          dynamicFields.value.splice(0, dynamicFields.value.length, ...data)

          for (const field of data) {
            const key = field.key_name
            if (localForm[key] === undefined) {
              if (field.default !== undefined && field.default !== null) {
                switch (field.type) {
                  case 'bool':
                    localForm[key] = field.default === 'true' || field.default === true
                    break
                  case 'int':
                    localForm[key] = parseInt(String(field.default), 10) || 0
                    break
                  case 'string':
                  case 'password':
                    localForm[key] = String(field.default)
                    break
                  default:
                    localForm[key] = field.default
                }
              } else {
                switch (field.type) {
                  case 'bool':
                    localForm[key] = false
                    break
                  case 'int':
                    localForm[key] = 0
                    break
                  case 'string':
                  case 'password':
                    localForm[key] = ''
                    break
                  default:
                    localForm[key] = ''
                }
              }
            }
          }

          await nextTick()

          if (!props.isEdit) {
            const has = (k: string) => data.some((f: StorageConfigTemplate) => f.key_name === k)
            const setIfUndef = (k: string, v: unknown) => {
              const cur = localForm[k]
              if (cur === undefined || cur === null || cur === '') localForm[k] = v
            }
            if (has('use_https')) setIfUndef('use_https', true)
            if (has('hide_remote_url')) setIfUndef('hide_remote_url', false)
            if (newType === 's3' && has('use_path_style')) setIfUndef('use_path_style', false)
          }
          showBaseGroup.value = true
          showLinkGroup.value = false
          showAdvancedGroup.value = false
        } else {
          dynamicFields.value = []
        }
      } catch (_error) {
        dynamicFields.value = []
      }
    },
    { immediate: true }
  )

  const s3ProvidersTip = $t('admin.channels.tips.s3Providers')

  const showBaseGroup = ref(true)
  const showLinkGroup = ref(false)
  const showAdvancedGroup = ref(false)
  const linkKeys = new Set(['custom_domain', 'use_https', 'hide_remote_url', 'use_path_style', 'access_control', 'domain'])
  const grouped = computed(() => {
    const base: StorageConfigTemplate[] = []
    const link: StorageConfigTemplate[] = []
    const adv: StorageConfigTemplate[] = []
    for (const f of dynamicFields.value) {
      if (f.required) {
        base.push(f)
        continue
      }
      if (linkKeys.has(f.key_name)) {
        link.push(f)
        continue
      }
      adv.push(f)
    }
    return { base, link, adv }
  })
</script>

<template>
  <div class="channel-form">
    <div class="form-group">
      <label><span class="required-star">*</span>{{ $t('admin.channels.form.name') }}</label>
      <CyberInput v-model="localForm.name" :placeholder="$t('admin.channels.form.namePlaceholder')" />
    </div>

    <div class="form-group">
      <label class="label-with-tip">
        <span class="label-text"> <span class="required-star">*</span>{{ $t('admin.channels.form.type') }} </span>
        <CyberTooltip v-if="localForm.type === 's3'" :content="s3ProvidersTip" placement="top">
          <i class="fas fa-info-circle tip-icon text-content" />
        </CyberTooltip>
      </label>
      <CyberDropdown
        v-model="localForm.type"
        :disabled="!!isEdit"
        :options="channelTypeOptions"
        :placeholder="$t('admin.channels.form.typePlaceholder')"
      />
    </div>

    <div class="group-header" @click="showBaseGroup = !showBaseGroup">
      <i :class="['fas', showBaseGroup ? 'fa-chevron-down' : 'fa-chevron-right']" />
      <span>{{ $t('admin.channels.form.basicRequired') }}</span>
      <span class="count">{{ grouped.base.length }}</span>
    </div>
    <div v-show="showBaseGroup">
      <div v-for="f in grouped.base" :key="f.key_name" class="form-group">
        <label class="label-with-tip">
          <span class="label-text">
            <span v-if="f.required" class="required-star">*</span>
            {{ f.name }}
          </span>
          <CyberTooltip v-if="getFieldTip(f.key_name)" :content="getFieldTip(f.key_name)" placement="top">
            <i class="fas fa-info-circle tip-icon text-content" />
          </CyberTooltip>
        </label>
        <CyberDropdown
          v-if="Array.isArray(f.options) && f.options.length > 0"
          v-model="localForm[f.key_name]"
          :options="f.options.map((o: string) => ({ label: o, value: o }))"
          :placeholder="f.description || $t('admin.channels.form.pleaseSelect')"
        />
        <CyberInput
          v-else-if="f.type === 'string' || f.type === 'password'"
          v-model="localForm[f.key_name]"
          :type="f.type === 'password' ? 'password' : 'text'"
          :placeholder="f.description || ''"
        />
        <CyberInput
          v-else-if="f.type === 'int'"
          v-model.number="localForm[f.key_name]"
          type="number"
          :placeholder="f.description || ''"
        />
        <div v-else-if="f.type === 'bool'" class="radio-group">
          <CyberRadio v-model="localForm[f.key_name]" :value="true">{{ $t('admin.channels.form.yes') }}</CyberRadio>
          <CyberRadio v-model="localForm[f.key_name]" :value="false">{{ $t('admin.channels.form.no') }}</CyberRadio>
        </div>
        <CyberInput v-else v-model="localForm[f.key_name]" :placeholder="f.description || ''" />
      </div>
    </div>

    <div class="group-header" @click="showLinkGroup = !showLinkGroup">
      <i :class="['fas', showLinkGroup ? 'fa-chevron-down' : 'fa-chevron-right']" />
      <span>{{ $t('admin.channels.form.linkPermission') }}</span>
      <span class="count">{{ grouped.link.length }}</span>
    </div>
    <div v-show="showLinkGroup">
      <div v-for="f in grouped.link" :key="f.key_name" class="form-group">
        <label class="label-with-tip">
          <span class="label-text">
            <span v-if="f.required" class="required-star">*</span>
            {{ f.name }}
          </span>
          <CyberTooltip v-if="getFieldTip(f.key_name)" :content="getFieldTip(f.key_name)" placement="top">
            <i class="fas fa-info-circle tip-icon text-content" />
          </CyberTooltip>
        </label>
        <CyberDropdown
          v-if="Array.isArray(f.options) && f.options.length > 0"
          v-model="localForm[f.key_name]"
          :options="f.options.map((o: string) => ({ label: o, value: o }))"
          :placeholder="f.description || $t('admin.channels.form.pleaseSelect')"
        />
        <CyberInput
          v-else-if="f.type === 'string' || f.type === 'password'"
          v-model="localForm[f.key_name]"
          :type="f.type === 'password' ? 'password' : 'text'"
          :placeholder="f.description || ''"
        />
        <CyberInput
          v-else-if="f.type === 'int'"
          v-model.number="localForm[f.key_name]"
          type="number"
          :placeholder="f.description || ''"
        />
        <div v-else-if="f.type === 'bool'" class="radio-group">
          <CyberRadio v-model="localForm[f.key_name]" :value="true">{{ $t('admin.channels.form.yes') }}</CyberRadio>
          <CyberRadio v-model="localForm[f.key_name]" :value="false">{{ $t('admin.channels.form.no') }}</CyberRadio>
        </div>
        <CyberInput v-else v-model="localForm[f.key_name]" :placeholder="f.description || ''" />
      </div>
    </div>

    <div class="group-header" @click="showAdvancedGroup = !showAdvancedGroup">
      <i :class="['fas', showAdvancedGroup ? 'fa-chevron-down' : 'fa-chevron-right']" />
      <span>{{ $t('admin.channels.form.advanced') }}</span>
      <span class="count">{{ grouped.adv.length }}</span>
    </div>
    <div v-show="showAdvancedGroup">
      <div v-for="f in grouped.adv" :key="f.key_name" class="form-group">
        <label class="label-with-tip">
          <span class="label-text">
            <span v-if="f.required" class="required-star">*</span>
            {{ f.name }}
          </span>
          <CyberTooltip v-if="getFieldTip(f.key_name)" :content="getFieldTip(f.key_name)" placement="top">
            <i class="fas fa-info-circle tip-icon text-content" />
          </CyberTooltip>
        </label>
        <CyberDropdown
          v-if="Array.isArray(f.options) && f.options.length > 0"
          v-model="localForm[f.key_name]"
          :options="f.options.map((o: string) => ({ label: o, value: o }))"
          :placeholder="f.description || $t('admin.channels.form.pleaseSelect')"
        />
        <CyberInput
          v-else-if="f.type === 'string' || f.type === 'password'"
          v-model="localForm[f.key_name]"
          :type="f.type === 'password' ? 'password' : 'text'"
          :placeholder="f.description || ''"
        />
        <CyberInput
          v-else-if="f.type === 'int'"
          v-model.number="localForm[f.key_name]"
          type="number"
          :placeholder="f.description || ''"
        />
        <div v-else-if="f.type === 'bool'" class="radio-group">
          <CyberRadio v-model="localForm[f.key_name]" :value="true">{{ $t('admin.channels.form.yes') }}</CyberRadio>
          <CyberRadio v-model="localForm[f.key_name]" :value="false">{{ $t('admin.channels.form.no') }}</CyberRadio>
        </div>
        <CyberInput v-else v-model="localForm[f.key_name]" :placeholder="f.description || ''" />
      </div>
    </div>

    <div class="form-group">
      <label>{{ $t('admin.channels.form.isDefault') }}</label>
      <div class="radio-group">
        <CyberRadio v-model="localForm.is_default" :value="true"> {{ $t('admin.channels.form.yes') }} </CyberRadio>
        <CyberRadio v-model="localForm.is_default" :value="false"> {{ $t('admin.channels.form.no') }} </CyberRadio>
      </div>
    </div>

    <div class="form-group">
      <label>{{ $t('admin.channels.form.status') }}</label>
      <div class="radio-group">
        <CyberRadio v-model="localForm.status" :value="1" :disabled="isEdit && localForm.type === 'local'">
          {{ $t('admin.channels.form.enabled') }}
        </CyberRadio>
        <CyberRadio v-model="localForm.status" :value="0" :disabled="isEdit && localForm.type === 'local'">
          {{ $t('admin.channels.form.disabled') }}
        </CyberRadio>

        <div v-if="isEdit && localForm.type === 'local'" class="status-local-tip">
          <i class="fas fa-info-circle" />{{ $t('admin.channels.form.localStatusTip') }}
        </div>
      </div>
    </div>

    <div class="form-group">
      <label>{{ $t('admin.channels.form.remark') }}</label>
      <CyberInput v-model="localForm.remark" :placeholder="$t('admin.channels.form.remarkPlaceholder')" />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .channel-form {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
  }

  .form-group label {
    font-size: var(--text-sm);
    color: rgba(255, 255, 255, 0.7);
  }

  .required-star {
    margin-right: var(--space-xs);
    color: var(--color-error-500);
  }

  .label-with-tip {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
  }
  .label-text {
    display: inline-flex;
    align-items: center;
  }
  .tip-icon {
    font-size: var(--text-xs);
  }

  .form-row {
    display: flex;
    gap: var(--space-md);
  }

  .flex-1 {
    flex: 1;
  }

  .radio-group {
    display: flex;
    gap: var(--space-md);
  }

  .status-local-tip {
    display: inline-flex;
    align-items: center;
    border: none;
    background: transparent;
    padding: 0 0 var(--space-xs);
    font-size: var(--text-xs);
    color: rgba(255, 255, 255, 0.7);
  }
  .group-header {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    margin-top: var(--space-sm);
    margin-bottom: var(--space-xs);
    cursor: pointer;
    user-select: none;
    color: var(--color-cyber-light);
  }
  .group-header .count {
    margin-left: var(--space-sm);
    font-size: var(--text-xs);
    opacity: var(--opacity-disabled);
  }

  .status-local-tip i {
    margin-right: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-brand-500);
  }

  .form-tip {
    display: flex;
    align-items: center;
    margin-top: var(--space-sm);
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    background: rgba(255, 255, 255, 0.05);
    font-size: var(--text-xs);
    color: rgba(255, 255, 255, 0.6);
  }

  .form-tip i {
    margin-right: var(--space-sm);
    font-size: var(--text-xs);
    color: var(--color-brand-500);
  }

  .advanced-settings {
    margin-top: var(--space-xl);
    padding-top: var(--space-xl);
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }

  .advanced-settings .form-group {
    margin-bottom: var(--space-md);
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    margin-bottom: var(--space-xl);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-content-default);
  }

  .section-title i {
    font-size: var(--text-xs);
  }

  .subsection-title {
    margin-top: var(--space-md);
    margin-bottom: var(--space-sm);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    color: rgba(255, 255, 255, 0.6);
  }
</style>
