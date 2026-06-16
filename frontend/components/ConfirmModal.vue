<script setup lang="ts">
withDefaults(
  defineProps<{
    open: boolean
    title: string
    body: string
    confirmLabel?: string
    danger?: boolean
  }>(),
  { confirmLabel: '', danger: false },
)

const emit = defineEmits<{ confirm: []; close: [] }>()
</script>

<template>
  <Modal :open="open" :title="title" @close="emit('close')">
    <p class="text-sm text-ink-soft">{{ body }}</p>
    <div class="mt-6 flex justify-end gap-3">
      <button
        type="button"
        class="rounded-full px-5 py-2 text-sm font-bold text-ink-soft transition hover:bg-surface-container-low"
        @click="emit('close')"
      >
        {{ $t('common.cancel') }}
      </button>
      <button
        type="button"
        class="rounded-full px-5 py-2 text-sm font-bold text-white shadow-soft-primary transition hover:scale-105"
        :class="danger ? 'bg-error' : 'bg-primary'"
        @click="emit('confirm')"
      >
        {{ confirmLabel || $t('common.confirm') }}
      </button>
    </div>
  </Modal>
</template>
