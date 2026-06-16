<script setup lang="ts">
defineProps<{
  modelValue: string
  options: { value: string; label: string }[]
}>()

const emit = defineEmits<{ 'update:modelValue': [string] }>()
</script>

<template>
  <div class="inline-flex rounded-full bg-surface-container p-1">
    <button
      v-for="opt in options"
      :key="opt.value"
      type="button"
      class="rounded-full px-4 py-2 text-sm font-bold transition"
      :class="
        modelValue === opt.value
          ? 'bg-primary text-white shadow-soft-primary'
          : 'text-ink-soft hover:text-primary'
      "
      @click="emit('update:modelValue', opt.value)"
    >
      {{ opt.label }}
    </button>
  </div>
</template>
