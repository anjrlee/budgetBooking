<script setup lang="ts">
import { Line } from 'vue-chartjs'

const props = defineProps<{
  labels: string[]
  data: number[]
  color?: string
}>()

const color = computed(() => props.color || '#2D4B9B')

const chartData = computed(() => ({
  labels: props.labels,
  datasets: [
    {
      data: props.data,
      borderColor: color.value,
      backgroundColor: `${color.value}26`,
      pointBackgroundColor: color.value,
      tension: 0.35,
      fill: true,
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: {
      grid: { display: false },
      ticks: { font: { family: 'Nunito Sans', size: 11 }, color: '#444651' },
    },
    y: {
      grid: { color: '#EDE7DF' },
      ticks: { font: { family: 'Nunito Sans', size: 11 }, color: '#444651' },
    },
  },
}
</script>

<template>
  <div class="relative h-56 w-full sm:h-72">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
