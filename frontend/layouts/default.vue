<script setup lang="ts">
const authStore = useAuthStore()
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppHeader />
    <AppNav v-if="authStore.isLoggedIn" />
    <main class="mx-auto w-full max-w-5xl flex-1 px-4 pb-24 pt-6 sm:px-6 md:pb-10">
      <slot />
    </main>
  </div>
</template>
