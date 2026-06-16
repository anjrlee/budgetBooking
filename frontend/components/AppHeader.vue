<script setup lang="ts">
const { locale, locales, setLocale } = useI18n()
const authStore = useAuthStore()
const api = useApi()

const otherLocale = computed(() => (locale.value === 'zh-Hant' ? 'en' : 'zh-Hant'))

async function toggleLocale() {
  const next = otherLocale.value
  await setLocale(next)
  if (authStore.isLoggedIn) {
    authStore.setLanguage(next)
    try {
      await api('/user/setting', { method: 'POST', body: { language: next } })
    } catch {
      // non-critical; UI already switched
    }
  }
}
</script>

<template>
  <header class="sticky top-0 z-30 border-b border-outline/20 bg-background/80 backdrop-blur">
    <div class="mx-auto flex max-w-5xl items-center justify-between px-4 py-3 sm:px-6">
      <NuxtLink to="/" class="flex items-center gap-2">
        <span class="material-symbols-outlined text-primary">savings</span>
        <span class="font-display text-xl font-bold text-primary-dark">{{ $t('common.appName') }}</span>
        <span class="rounded-full bg-primary/15 px-1.5 py-0.5 text-[10px] font-bold uppercase tracking-wide text-primary">beta</span>
      </NuxtLink>

      <div class="flex items-center gap-3">
        <button
          type="button"
          class="rounded-full bg-surface px-3 py-1.5 text-sm font-bold text-ink-soft shadow-soft transition hover:scale-105"
          @click="toggleLocale"
        >
          {{ $t(`setting.languages.${otherLocale}`) }}
        </button>

        <NuxtLink
          v-if="authStore.isLoggedIn"
          to="/setting"
          class="flex h-9 w-9 items-center justify-center overflow-hidden rounded-full bg-surface shadow-soft"
        >
          <img
            v-if="authStore.user?.picture"
            :src="authStore.user.picture"
            :alt="authStore.user?.name"
            class="h-full w-full object-cover"
          >
          <span v-else class="material-symbols-outlined text-primary">account_circle</span>
        </NuxtLink>
      </div>
    </div>
  </header>
</template>
