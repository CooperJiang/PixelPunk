import { ref, nextTick } from 'vue'
import type { Density } from '../types'

export function useDisplayDensity(defaultDensity: Density = 'normal') {
  const currentDensity = ref<Density>(defaultDensity)

  const changeDensity = (density: Density) => {
    currentDensity.value = density
    try {
      localStorage.setItem('gallery_density', density)
    } catch {}
    nextTick(() => {
      window.dispatchEvent(new Event('resize'))
    })
  }

  const loadSavedDensity = () => {
    try {
      const saved = localStorage.getItem('gallery_density')
      const allowed = ['compact', 'normal', 'comfortable']
      if (saved && allowed.includes(saved)) currentDensity.value = saved as Density
    } catch {}
  }

  return { currentDensity, changeDensity, loadSavedDensity }
}
