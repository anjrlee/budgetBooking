<script setup lang="ts">
export interface DataTableColumn {
  key: string
  label: string
  align?: 'left' | 'right' | 'center'
}

defineProps<{
  columns: DataTableColumn[]
  rows: Record<string, any>[]
  rowKey?: string
}>()
</script>

<template>
  <div class="scroll-soft overflow-x-auto">
    <table class="w-full min-w-[480px] text-left text-sm">
      <thead>
        <tr>
          <th
            v-for="col in columns"
            :key="col.key"
            class="whitespace-nowrap px-3 py-2 text-xs font-bold uppercase tracking-wide text-ink-soft"
            :class="{ 'text-right': col.align === 'right', 'text-center': col.align === 'center' }"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(row, i) in rows"
          :key="rowKey ? row[rowKey] : i"
          class="border-t border-outline/15 transition hover:bg-surface-container-low"
        >
          <td
            v-for="col in columns"
            :key="col.key"
            class="whitespace-nowrap px-3 py-3"
            :class="{ 'text-right': col.align === 'right', 'text-center': col.align === 'center' }"
          >
            <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
              {{ row[col.key] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
