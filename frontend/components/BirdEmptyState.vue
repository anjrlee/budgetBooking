<script setup lang="ts">
withDefaults(
  defineProps<{
    image?: string
    message: string
  }>(),
  { image: '/images/birds/bird-magnifier.png' },
)
</script>

<template>
  <div class="flex flex-col items-center gap-3 py-8 text-center">
    <img :src="image" alt="" class="h-32 w-32 object-contain">
    <p class="text-sm text-ink-soft">{{ message }}</p>
  </div>
</template>
