<script setup lang="ts">
const authStore = useAuthStore()
const { t } = useI18n()
const config = useRuntimeConfig()
const loginError = ref(false)

async function onLoginSuccess(idToken: string) {
  loginError.value = false
  try {
    await authStore.login(idToken)
  } catch {
    loginError.value = true
  }
}

function onLoginError() {
  loginError.value = true
}

const cards = computed(() => [
  { to: '/setBudget', icon: 'pie_chart', color: 'mint', title: t('home.cards.setBudget.title'), description: t('home.cards.setBudget.description') },
  { to: '/booking', icon: 'receipt_long', color: 'peach', title: t('home.cards.booking.title'), description: t('home.cards.booking.description') },
  { to: '/statistics', icon: 'monitoring', color: 'sky', title: t('home.cards.statistics.title'), description: t('home.cards.statistics.description') },
  { to: '/setting', icon: 'settings', color: 'lavender', title: t('home.cards.setting.title'), description: t('home.cards.setting.description') },
])
</script>

<template>
  <div>
    <!-- Logged in: 4 navigation pageCards -->
    <div v-if="authStore.isLoggedIn">
      <h1 class="mb-1 font-display text-2xl font-bold text-ink sm:text-3xl">
        {{ $t('home.greeting', { name: authStore.user?.name || '' }) }}
      </h1>
      <p class="mb-6 text-ink-soft">{{ $t('home.welcomeSubtitle') }}</p>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <PageCard
          v-for="card in cards"
          :key="card.to"
          :to="card.to"
          :icon="card.icon"
          :color="card.color"
          :title="card.title"
          :description="card.description"
        />
      </div>

      <div class="mt-8 flex justify-center gap-4">
        <NuxtLink to="/about" class="text-xs text-ink-soft underline-offset-2 hover:underline">{{ $t('home.about') }}</NuxtLink>
        <a v-if="config.public.sponsorUrl" :href="config.public.sponsorUrl" target="_blank" rel="noopener" class="text-xs text-ink-soft underline-offset-2 hover:underline">{{ $t('home.sponsor') }}</a>
        <a v-if="config.public.reportIssueUrl" :href="config.public.reportIssueUrl" target="_blank" rel="noopener" class="text-xs text-ink-soft underline-offset-2 hover:underline">{{ $t('home.reportIssue') }}</a>
      </div>
    </div>

    <!-- Logged out: welcome hero + Google login -->
    <div v-else class="flex flex-col items-center pt-8 text-center sm:pt-16">
      <h1 class="font-display text-3xl font-bold text-primary-dark sm:text-4xl">
        {{ $t('home.welcomeTitle') }}
      </h1>
      <p class="mt-2 text-ink-soft">{{ $t('home.welcomeSubtitle') }}</p>

      <div class="relative my-8 h-56 w-56 sm:h-64 sm:w-64">
        <img src="/images/birds/confused-bubbles.jpg" alt="" class="h-full w-full rounded-card object-cover shadow-soft">

        <div class="float-bubble absolute -left-4 -top-4 flex h-12 w-12 items-center justify-center rounded-full bg-mint shadow-soft" style="animation-delay: 0s">
          <span class="material-symbols-outlined text-primary-dark">account_balance_wallet</span>
        </div>
        <div class="float-bubble absolute -right-6 top-2 flex h-12 w-12 items-center justify-center rounded-full bg-peach shadow-soft" style="animation-delay: 0.6s">
          <span class="material-symbols-outlined text-primary-dark">attach_money</span>
        </div>
        <div class="float-bubble absolute -left-6 bottom-8 flex h-12 w-12 items-center justify-center rounded-full bg-sky shadow-soft" style="animation-delay: 1.2s">
          <span class="material-symbols-outlined text-primary-dark">cottage</span>
        </div>
        <div class="float-bubble absolute -right-4 -bottom-4 flex h-12 w-12 items-center justify-center rounded-full bg-lavender shadow-soft" style="animation-delay: 1.8s">
          <span class="material-symbols-outlined text-primary-dark">egg_alt</span>
        </div>
      </div>

      <p class="mb-6 font-display text-lg font-bold text-ink">{{ $t('home.birdMessage') }}</p>

      <GoogleLoginButton @success="onLoginSuccess" @error="onLoginError" />
      <p v-if="loginError" class="mt-3 text-sm text-error">{{ $t('home.loginError') }}</p>

      <div class="mt-10 flex justify-center gap-4">
        <NuxtLink to="/about" class="text-xs text-ink-soft underline-offset-2 hover:underline">{{ $t('home.about') }}</NuxtLink>
        <a v-if="config.public.sponsorUrl" :href="config.public.sponsorUrl" target="_blank" rel="noopener" class="text-xs text-ink-soft underline-offset-2 hover:underline">{{ $t('home.sponsor') }}</a>
        <a v-if="config.public.reportIssueUrl" :href="config.public.reportIssueUrl" target="_blank" rel="noopener" class="text-xs text-ink-soft underline-offset-2 hover:underline">{{ $t('home.reportIssue') }}</a>
      </div>
    </div>
  </div>
</template>
