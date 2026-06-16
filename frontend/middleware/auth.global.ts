const PROTECTED_PATHS = ['/setBudget', '/booking', '/statistics', '/report', '/budget', '/setting']

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()
  if (PROTECTED_PATHS.includes(to.path) && !authStore.isLoggedIn) {
    return navigateTo('/')
  }
})
