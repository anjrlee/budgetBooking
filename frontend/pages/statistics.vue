<script setup lang="ts">
const { t } = useI18n()

const cards = computed(() => [
  { to: '/report', icon: 'description', color: 'sky', title: t('statistics.cards.report.title'), description: t('statistics.cards.report.description') },
  { to: '/budget', icon: 'account_balance_wallet', color: 'mint', title: t('statistics.cards.budget.title'), description: t('statistics.cards.budget.description') },
])
</script>

<template>
  <div>
    <h1 class="mb-6 font-display text-2xl font-bold text-ink sm:text-3xl">{{ $t('statistics.title') }}</h1>
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
  </div>
</template>
