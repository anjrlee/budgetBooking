<script setup lang="ts">
withDefaults(defineProps<{ title?: string }>(), { title: '' })
</script>

<template>
  <section class="rounded-card bg-gradient-to-br from-surface to-surface-container-low p-5 shadow-soft sm:p-8">
    <h2 v-if="title" class="mb-4 font-display text-xl font-bold text-ink">{{ title }}</h2>
    <slot />
  </section>
</template>
