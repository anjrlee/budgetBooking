<script setup lang="ts">
defineProps<{ open: boolean; title?: string }>()
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-ink/40 p-4"
        @click.self="emit('close')"
      >
        <div class="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-card bg-surface p-6 shadow-soft-lift">
          <div class="mb-4 flex items-center justify-between gap-4">
            <h3 class="font-display text-lg font-bold text-ink">{{ title }}</h3>
            <button type="button" class="text-ink-soft transition hover:text-ink" @click="emit('close')">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
