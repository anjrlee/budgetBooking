<script setup lang="ts">
const route = useRoute()

const items = [
  { to: '/', icon: 'home', label: 'nav.home' },
  { to: '/booking', icon: 'receipt_long', label: 'nav.booking' },
  { to: '/setBudget', icon: 'pie_chart', label: 'nav.setBudget' },
  { to: '/statistics', icon: 'monitoring', label: 'nav.statistics' },
  { to: '/setting', icon: 'settings', label: 'nav.setting' },
]

function isActive(path: string) {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <nav
    class="fixed inset-x-0 bottom-0 z-30 border-t border-outline/20 bg-surface/95 backdrop-blur md:static md:mx-auto md:mt-2 md:max-w-5xl md:border-none md:bg-transparent md:px-6"
  >
    <ul class="mx-auto flex max-w-5xl items-center justify-between px-2 py-1 md:justify-center md:gap-2 md:py-0">
      <li v-for="item in items" :key="item.to" class="flex-1 md:flex-none">
        <NuxtLink
          :to="item.to"
          class="flex flex-col items-center gap-0.5 rounded-2xl px-2 py-2 text-xs font-bold transition md:flex-row md:gap-2 md:rounded-full md:px-4 md:py-2 md:text-sm"
          :class="isActive(item.to)
            ? 'text-primary md:bg-primary md:text-white'
            : 'text-ink-soft hover:text-primary'"
        >
          <span class="material-symbols-outlined !text-[22px] md:!text-[20px]">{{ item.icon }}</span>
          <span>{{ $t(item.label) }}</span>
        </NuxtLink>
      </li>
    </ul>
  </nav>
</template>
