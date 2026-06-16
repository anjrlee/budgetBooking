<script setup lang="ts">
const authStore = useAuthStore()
const api = useApi()
const { locale, setLocale } = useI18n()

const showLogoutConfirm = ref(false)
const saving = ref(false)

const languageOptions = [
  { value: 'zh-Hant', label: 'setting.languages.zh-Hant' },
  { value: 'en', label: 'setting.languages.en' },
]

async function onLanguageChange(lang: string) {
  if (lang === locale.value) return
  saving.value = true
  try {
    await api('/user/setting', { method: 'POST', body: { language: lang } })
    authStore.setLanguage(lang)
    await setLocale(lang as 'zh-Hant' | 'en')
  } finally {
    saving.value = false
  }
}

function logout() {
  authStore.logout()
  showLogoutConfirm.value = false
  navigateTo('/')
}
</script>

<template>
  <div>
    <h1 class="mb-6 font-display text-2xl font-bold text-ink sm:text-3xl">{{ $t('setting.title') }}</h1>

    <FunctionCard :title="$t('setting.profile')">
      <div class="flex items-center gap-4">
        <img
          v-if="authStore.user?.picture"
          :src="authStore.user.picture"
          :alt="authStore.user?.name"
          class="h-16 w-16 rounded-full object-cover shadow-soft"
        >
        <span v-else class="material-symbols-outlined !text-[64px] text-primary">account_circle</span>
        <div>
          <p class="font-display text-lg font-bold text-ink">{{ authStore.user?.name }}</p>
          <p class="text-sm text-ink-soft">{{ authStore.user?.email }}</p>
        </div>
      </div>
    </FunctionCard>

    <FunctionCard :title="$t('setting.language')" class="mt-6">
      <div class="flex flex-wrap gap-3">
        <button
          v-for="opt in languageOptions"
          :key="opt.value"
          type="button"
          :disabled="saving"
          class="rounded-full px-5 py-2 text-sm font-bold transition disabled:opacity-50"
          :class="locale === opt.value ? 'bg-primary text-white shadow-soft-primary' : 'bg-surface-container text-ink-soft hover:bg-surface-container-high'"
          @click="onLanguageChange(opt.value)"
        >
          {{ $t(opt.label) }}
        </button>
      </div>
    </FunctionCard>

    <FunctionCard class="mt-6">
      <button type="button" class="btn-secondary w-full text-error sm:w-auto" @click="showLogoutConfirm = true">
        <span class="material-symbols-outlined mr-2 align-middle !text-[20px]">logout</span>
        {{ $t('setting.logout') }}
      </button>
    </FunctionCard>

    <ConfirmModal
      :open="showLogoutConfirm"
      :title="$t('setting.logoutConfirmTitle')"
      :body="$t('setting.logoutConfirmBody')"
      :confirm-label="$t('setting.logout')"
      danger
      @close="showLogoutConfirm = false"
      @confirm="logout"
    />
  </div>
</template>
