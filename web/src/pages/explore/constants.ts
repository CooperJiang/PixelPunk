import type { SortOption, DensityOption } from './types'
import { useTexts } from '@/composables/useTexts'

export const getSortOptions = (): SortOption[] => {
  const { $t } = useTexts()
  return [
    { value: 'newest', label: $t('explore.sort.newest'), icon: 'fas fa-clock' },
    { value: 'oldest', label: $t('explore.sort.oldest'), icon: 'fas fa-history' },
    { value: 'size', label: $t('explore.sort.size'), icon: 'fas fa-sort-amount-down' },
    { value: 'name', label: $t('explore.sort.name'), icon: 'fas fa-sort-alpha-down' },
  ]
}

export const SORT_OPTIONS: SortOption[] = getSortOptions()

export const getDensityOptions = (): DensityOption[] => {
  const { $t } = useTexts()
  return [
    { value: 'compact', label: $t('explore.density.compact'), icon: 'fas fa-compress' },
    { value: 'normal', label: $t('explore.density.normal'), icon: 'fas fa-arrows-alt' },
    { value: 'comfortable', label: $t('explore.density.comfortable'), icon: 'fas fa-expand' },
  ]
}

export const DENSITY_OPTIONS: DensityOption[] = getDensityOptions()
