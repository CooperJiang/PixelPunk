import { reactive, ref, watch } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import {
  createChannel,
  getChannelConfigs,
  updateChannel,
  getConfigTemplates,
  getChannelList,
  type StorageConfigTemplate,
} from '@/api/storage'
import type { CreateStorageChannelRequest, StorageChannel } from '@/api/types/index'
import Validator, { getValidationRules } from '@/utils/validation/validator'

export function useChannelForm() {
  const toast = useToast()
  const { $t } = useTexts()
  const R = getValidationRules($t)

  const currentChannel = ref<StorageChannel | null>(null)
  const showAddChannelDialog = ref(false)

  const channelForm = reactive<CreateStorageChannelRequest & Record<string, unknown>>({
    name: '',
    type: '',
    status: 1,
    is_default: false,
    remark: '',
  })

  const resetForm = () => {
    const baseKeys = new Set(['name', 'type', 'status', 'is_default', 'remark'])
    const formData = channelForm as Record<string, unknown>
    Object.keys(formData).forEach((k) => {
      if (!baseKeys.has(k)) delete formData[k]
    })
    channelForm.name = ''
    channelForm.type = ''
    channelForm.status = 1
    channelForm.is_default = false
    channelForm.remark = ''
  }

  const editChannel = async (channel: StorageChannel) => {
    currentChannel.value = channel
    resetForm()
    channelForm.name = channel.name
    channelForm.type = channel.type
    channelForm.status = channel.status
    channelForm.is_default = channel.is_default
    channelForm.remark = channel.remark || ''

    try {
      const result = await getChannelConfigs(channel.id)

      if (result.success) {
        const configs = result.data || []
        configs.forEach((config) => {
          if (config.key_name === 'access_key' || config.key_name === 'secret_key' || config.key_name === 'secret_id') {
            channelForm[config.key_name as keyof typeof channelForm] = config.is_secret ? '********' : config.value
          } else {
            try {
              const formData = channelForm as Record<string, unknown>
              if (config.type === 'int') {
                let intValue = 0
                if (config.value !== null && config.value !== undefined && config.value !== '') {
                  const parsed = parseInt(String(config.value), 10)
                  intValue = isNaN(parsed) ? 0 : parsed
                }
                formData[config.key_name] = intValue
              } else if (config.type === 'bool') {
                let boolValue = false
                const strValue = String(config.value || '').toLowerCase()
                if (strValue === 'true' || strValue === '1') {
                  boolValue = true
                }
                formData[config.key_name] = boolValue
              } else {
                formData[config.key_name] = config.value || ''
              }
            } catch {
              const formData = channelForm as Record<string, unknown>
              formData[config.key_name] = config.value
            }
          }
        })

        showAddChannelDialog.value = true
      }
    } catch (error: unknown) {
      toast.error(error.message)
    }
  }

  const saveChannel = async (onSuccess?: () => void) => {
    const nameCheck = Validator.validate(channelForm.name || '', [R.required], $t)
    if (!nameCheck.valid) return toast.error($t('admin.channels.form.nameRequired'))
    const typeCheck = Validator.validate(String(channelForm.type || ''), [R.required], $t)
    if (!typeCheck.valid) return toast.error($t('admin.channels.form.typeRequired'))

    try {
      const tmplRes = await getConfigTemplates(channelForm.type)
      const tmplData = tmplRes as Record<string, unknown>
      if (tmplData.success) {
        const templates: StorageConfigTemplate[] = (tmplData.data as StorageConfigTemplate[]) || []
        const formData = channelForm as Record<string, unknown>
        for (const t of templates) {
          const key = t.key_name
          const rawVal = formData[key]

          if (t.required) {
            const isEmpty = rawVal === undefined || rawVal === null || String(rawVal).trim() === ''
            if (isEmpty) {
              toast.error($t('admin.channels.form.fieldRequired', { field: t.name }))
              return
            }
          }

          if (rawVal !== undefined && rawVal !== null && String(rawVal) !== '') {
            if (t.type === 'int') {
              const n = Number(rawVal)
              if (!Number.isFinite(n) || !Number.isInteger(n)) {
                toast.error($t('admin.channels.form.mustBeInt', { field: t.name }))
                return
              }
            } else if (t.type === 'bool') {
              const isBool = typeof rawVal === 'boolean'
              if (!isBool) {
                toast.error($t('admin.channels.form.mustBeBool', { field: t.name }))
                return
              }
            } else if (t.type === 'string' || t.type === 'password') {
              if (typeof rawVal !== 'string') {
                toast.error($t('admin.channels.form.mustBeString', { field: t.name }))
                return
              }
            }
          }
        }
      }
    } catch {
      toast.error($t('admin.channels.form.loadTemplateFailed'))
      return
    }

    if (currentChannel.value && channelForm.is_default === false) {
      try {
        const listRes = await getChannelList()
        if (listRes.success) {
          const currentId = currentChannel.value.id
          const others = (listRes.data || []).filter((c) => c.id !== currentId && c.is_default)
          if (others.length === 0) {
            toast.error($t('admin.channels.form.needDefaultChannel'))
            return
          }
        }
      } catch {}
    }

    const channelData: Record<string, unknown> = {
      name: channelForm.name,
      type: channelForm.type,
      status: channelForm.status,
      is_default: channelForm.is_default,
      remark: channelForm.remark || '',
    }

    const configs: Record<string, unknown> = {}

    const baseKeys = new Set(['name', 'type', 'status', 'is_default', 'remark'])
    const formData = channelForm as Record<string, unknown>
    Object.keys(formData).forEach((key) => {
      if (!baseKeys.has(key)) {
        const value = formData[key]
        if (value !== '********' || typeof value === 'boolean' || typeof value === 'number') {
          configs[key] = value
        }
      }
    })

    if (Object.keys(configs).length > 0) {
      channelData.configs = configs
    }

    try {
      if (currentChannel.value) {
        const result = await updateChannel(currentChannel.value.id, channelData)
        if (result.success) {
          toast.success($t('admin.channels.form.channelUpdated'))
          showAddChannelDialog.value = false
          onSuccess?.()
        }
      } else {
        const result = await createChannel(channelData)
        if (result.success) {
          toast.success($t('admin.channels.form.channelAdded'))
          showAddChannelDialog.value = false
          onSuccess?.()
        }
      }
    } catch (error: unknown) {
      toast.error(error.message)
    }
  }

  watch(showAddChannelDialog, (newVal) => {
    if (!newVal) {
      currentChannel.value = null
      resetForm()
    }
  })

  return {
    currentChannel,
    showAddChannelDialog,
    channelForm,

    resetForm,
    editChannel,
    saveChannel,
  }
}
