<script setup lang="ts">
declare global {
  interface Window {
    google?: any
  }
}

const config = useRuntimeConfig()
const emit = defineEmits<{ success: [string]; error: [] }>()
const buttonEl = ref<HTMLElement | null>(null)

function loadScript(): Promise<void> {
  return new Promise((resolve, reject) => {
    if (window.google?.accounts?.id) {
      resolve()
      return
    }
    const script = document.createElement('script')
    script.src = 'https://accounts.google.com/gsi/client'
    script.async = true
    script.defer = true
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Failed to load Google Identity Services'))
    document.head.appendChild(script)
  })
}

onMounted(async () => {
  if (!config.public.googleClientId) return
  try {
    await loadScript()
    window.google.accounts.id.initialize({
      client_id: config.public.googleClientId,
      callback: (response: { credential?: string }) => {
        if (response?.credential) emit('success', response.credential)
        else emit('error')
      },
    })
    if (buttonEl.value) {
      window.google.accounts.id.renderButton(buttonEl.value, {
        theme: 'outline',
        size: 'large',
        shape: 'pill',
        text: 'signin_with',
        width: 280,
      })
    }
  } catch {
    emit('error')
  }
})
</script>

<template>
  <div class="flex flex-col items-center gap-2">
    <div ref="buttonEl" />
    <p v-if="!config.public.googleClientId" class="text-xs text-error">
      Google Client ID is not configured (NUXT_PUBLIC_GOOGLE_CLIENT_ID).
    </p>
  </div>
</template>
