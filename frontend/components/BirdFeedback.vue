<script setup lang="ts">
import { gsap } from 'gsap'

const props = defineProps<{ type: 'success' | 'error'; message: string }>()
const el = ref<HTMLElement | null>(null)

onMounted(() => {
  if (!el.value) return
  gsap.fromTo(
    el.value,
    { y: 16, opacity: 0, scale: 0.85 },
    { y: 0, opacity: 1, scale: 1, duration: 0.5, ease: 'back.out(1.7)' },
  )
})
</script>

<template>
  <div
    ref="el"
    class="flex items-center gap-3 rounded-card p-4"
    :class="type === 'success' ? 'bg-mint' : 'bg-error-container'"
  >
    <span
      class="material-symbols-outlined !text-[32px]"
      :class="type === 'success' ? 'text-success' : 'text-error'"
    >
      {{ type === 'success' ? 'flutter_dash' : 'priority_high' }}
    </span>
    <p class="text-sm font-bold" :class="type === 'success' ? 'text-success' : 'text-error'">
      {{ message }}
    </p>
  </div>
</template>
