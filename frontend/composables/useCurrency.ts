export function useCurrency() {
  const { locale } = useI18n()

  function formatCurrency(amount: number) {
    const safeAmount = Number.isFinite(amount) ? amount : 0
    return new Intl.NumberFormat(locale.value === 'zh-Hant' ? 'zh-TW' : 'en-US', {
      style: 'currency',
      currency: 'TWD',
      maximumFractionDigits: 0,
    }).format(safeAmount)
  }

  return { formatCurrency }
}
