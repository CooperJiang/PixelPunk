<script setup lang="ts">
  import { ref, computed, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { CategoryTemplate } from '@/api/types/category'

  defineOptions({
    name: 'TemplateForm',
  })

  const props = defineProps<{
    modelValue: CategoryTemplate | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: CategoryTemplate | null]
  }>()

  const { $t } = useTexts()
  const formData = ref<CategoryTemplate | null>(null)

  watch(
    () => props.modelValue,
    (newValue) => {
      formData.value = newValue ? { ...newValue } : null
    },
    { immediate: true }
  )

  watch(
    formData,
    (newValue) => {
      emit('update:modelValue', newValue)
    },
    { deep: true }
  )

  /* 将 Font Awesome class 转换为 locale key */
  const iconClassToKey = (iconClass: string): string => {
    const name = iconClass.replace(/^(fas|far) fa-/, '')
    return name.replace(/-([a-z])/g, (_, letter) => letter.toUpperCase())
  }

  const iconValues = [
    'fas fa-user',
    'fas fa-users',
    'fas fa-baby',
    'fas fa-child',
    'fas fa-mars',
    'fas fa-venus',
    'fas fa-user-circle',
    'far fa-smile',
    'fas fa-mask',
    'fas fa-mountain',
    'fas fa-tree',
    'fas fa-leaf',
    'fas fa-spa',
    'fas fa-seedling',
    'fas fa-sun',
    'fas fa-moon',
    'fas fa-star',
    'fas fa-cloud',
    'fas fa-cloud-rain',
    'fas fa-snowflake',
    'fas fa-water',
    'fas fa-wave-square',
    'fas fa-tint',
    'fas fa-home',
    'fas fa-building',
    'fas fa-city',
    'fas fa-bridge',
    'fas fa-broadcast-tower',
    'fas fa-church',
    'fas fa-school',
    'fas fa-hospital',
    'fas fa-store',
    'fas fa-industry',
    'fas fa-dog',
    'fas fa-cat',
    'fas fa-fish',
    'fas fa-dove',
    'fas fa-horse',
    'fas fa-feather',
    'fas fa-spider',
    'fas fa-bug',
    'fas fa-dragon',
    'fas fa-paw',
    'fas fa-utensils',
    'fas fa-hamburger',
    'fas fa-pizza-slice',
    'fas fa-apple-alt',
    'fas fa-carrot',
    'fas fa-bread-slice',
    'fas fa-cheese',
    'fas fa-coffee',
    'fas fa-wine-glass',
    'fas fa-birthday-cake',
    'fas fa-ice-cream',
    'fas fa-candy-cane',
    'fas fa-car',
    'fas fa-truck',
    'fas fa-motorcycle',
    'fas fa-bicycle',
    'fas fa-plane',
    'fas fa-rocket',
    'fas fa-ship',
    'fas fa-subway',
    'fas fa-train',
    'fas fa-bus',
    'fas fa-book',
    'fas fa-mobile-alt',
    'fas fa-laptop',
    'fas fa-camera',
    'fas fa-music',
    'fas fa-gamepad',
    'fas fa-tools',
    'fas fa-key',
    'fas fa-wallet',
    'fas fa-shopping-bag',
    'fas fa-tshirt',
    'fas fa-shoe-prints',
    'fas fa-glasses',
    'fas fa-clock',
    'fas fa-paint-brush',
    'fas fa-palette',
    'fas fa-cut',
    'fas fa-pencil-alt',
    'fas fa-icons',
    'fas fa-magic',
    'fas fa-drafting-compass',
    'fas fa-image',
    'fas fa-images',
    'fas fa-video',
    'fas fa-football-ball',
    'fas fa-basketball-ball',
    'fas fa-volleyball-ball',
    'fas fa-swimmer',
    'fas fa-running',
    'fas fa-biking',
    'fas fa-skiing',
    'fas fa-dumbbell',
    'fas fa-trophy',
    'fas fa-medal',
    'fas fa-tag',
    'fas fa-folder',
    'fas fa-heart',
    'fas fa-fire',
    'fas fa-bolt',
    'fas fa-eye',
    'fas fa-gift',
    'fas fa-gem',
    'fas fa-crown',
    'fas fa-shield-alt',
    'fas fa-question',
    'fas fa-exclamation',
    'fas fa-ellipsis-h',
  ]

  const categoryIconOptions = computed(() =>
    iconValues.map((iconClass) => ({
      label: $t(`admin.categories.iconOptions.${iconClassToKey(iconClass)}`),
      value: iconClass,
    }))
  )
</script>

<template>
  <div v-if="formData" class="template-form">
    <div class="form-container">
      <div class="form-row">
        <label class="form-label">{{ $t('admin.categories.form.nameLabel') }}</label>
        <CyberInput
          v-model="formData.name"
          :placeholder="$t('admin.categories.form.namePlaceholder')"
          :maxlength="50"
          show-count
        />
      </div>

      <div class="form-row">
        <label class="form-label">{{ $t('admin.categories.form.descriptionLabel') }}</label>
        <CyberInput
          v-model="formData.description"
          type="textarea"
          :placeholder="$t('admin.categories.form.descriptionPlaceholder')"
          :maxlength="500"
          :rows="3"
          show-count
        />
      </div>

      <div class="form-row">
        <label class="form-label">{{ $t('admin.categories.form.iconLabel') }}</label>
        <CyberIconPicker v-model="formData.icon" :options="categoryIconOptions" />
      </div>

      <div class="form-row">
        <label class="form-label">{{ $t('admin.categories.form.sortOrderLabel') }}</label>
        <CyberInput
          v-model.number="formData.sort_order"
          type="number"
          :placeholder="$t('admin.categories.form.sortOrderPlaceholder')"
          :min="0"
          :max="9999"
        />
        <div class="form-tip">{{ $t('admin.categories.form.sortOrderTip') }}</div>
      </div>

      <div class="form-row">
        <label class="form-label">{{ $t('admin.categories.form.isPopularLabel') }}</label>
        <CyberSwitch v-model="formData.is_popular" :true-value="true" :false-value="false">
          <template #checked>{{ $t('admin.categories.form.popularChecked') }}</template>
          <template #unchecked>{{ $t('admin.categories.form.popularUnchecked') }}</template>
        </CyberSwitch>
        <div class="form-tip">{{ $t('admin.categories.form.popularTip') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .template-form {
    .form-container {
      display: flex;
      flex-direction: column;
      gap: var(--space-xl);
      padding: var(--space-md);
    }

    .form-row {
      display: flex;
      flex-direction: column;
      gap: var(--space-sm);
    }

    .form-label {
      color: var(--color-cyber-light);
      font-size: var(--text-sm);
      font-weight: var(--font-medium);
    }

    .form-tip {
      color: var(--color-content-default-content-muted);
      font-size: var(--text-xs);
      margin-top: var(--space-xs);
    }
  }
</style>
