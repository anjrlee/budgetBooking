<script setup lang="ts">
import { Pie } from 'vue-chartjs'

const props = defineProps<{
  labels: string[]
  data: number[]
  colors: string[]
}>()

const chartData = computed(() => ({
  labels: props.labels,
  datasets: [
    {
      data: props.data,
      backgroundColor: props.colors,
      borderColor: '#FFFFFF',
      borderWidth: 2,
      hoverOffset: 6,
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom' as const,
      labels: {
        font: { family: 'Nunito Sans', size: 12, weight: 'bold' as const },
        color: '#444651',
        usePointStyle: true,
        padding: 16,
      },
    },
  },
}
</script>

<template>
  <div class="relative h-64 w-full sm:h-80">
    <Pie :data="chartData" :options="chartOptions" />
  </div>
</template>
