<script setup lang="ts">
withDefaults(
  defineProps<{
    to: string
    icon: string
    title: string
    description: string
    color?: string
  }>(),
  { color: 'mint' },
)

const bgClass: Record<string, string> = {
  mint: 'bg-mint',
  lavender: 'bg-lavender',
  peach: 'bg-peach',
  sky: 'bg-sky',
}
</script>

<template>
  <NuxtLink
    :to="to"
    class="group flex flex-col gap-3 rounded-card bg-surface p-6 shadow-soft transition duration-300 hover:-translate-y-1 hover:shadow-soft-lift"
  >
    <div
      class="flex h-12 w-12 items-center justify-center rounded-2xl text-primary-dark transition group-hover:scale-110"
      :class="bgClass[color] || bgClass.mint"
    >
      <span class="material-symbols-outlined !text-[28px]">{{ icon }}</span>
    </div>
    <div>
      <h3 class="font-display text-lg font-bold text-ink">{{ title }}</h3>
      <p class="mt-1 text-sm text-ink-soft">{{ description }}</p>
    </div>
    <span class="material-symbols-outlined ml-auto mt-auto text-ink-soft transition group-hover:translate-x-1 group-hover:text-primary">
      arrow_forward_ios
    </span>
  </NuxtLink>
</template>
