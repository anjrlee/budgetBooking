interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: Record<string, unknown>
  query?: Record<string, unknown>
  /** Whether to attach the JWT Authorization header. Defaults to true. */
  auth?: boolean
}

/**
 * Wraps $fetch with the API base URL, JWT auth header injection, and global
 * 401 handling (clears session + redirects to "/").
 */
export function useApi() {
  const config = useRuntimeConfig()

  return async function apiFetch<T = unknown>(path: string, options: ApiOptions = {}): Promise<T> {
    const { auth: requiresAuth = true, ...fetchOptions } = options
    const authStore = useAuthStore()

    const headers: Record<string, string> = {}
    if (requiresAuth && authStore.token) {
      headers.Authorization = `Bearer ${authStore.token}`
    }

    try {
      return await $fetch<T>(path, {
        baseURL: config.public.apiBase,
        headers,
        ...fetchOptions,
      })
    } catch (error: any) {
      if (error?.response?.status === 401) {
        authStore.logout()
        await navigateTo('/')
      }
      throw error
    }
  }
}
