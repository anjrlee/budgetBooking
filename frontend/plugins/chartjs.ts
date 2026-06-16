import {
  ArcElement,
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
} from 'chart.js'

ChartJS.register(ArcElement, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, Filler)

export default defineNuxtPlugin(() => {})
