export default defineNuxtConfig({
  compatibilityDate: '2024-09-01',
  devtools: { enabled: true },
  // JWT auth lives in localStorage; render fully client-side to avoid SSR/auth-state mismatches.
  ssr: false,

  // Works around a vite-node dev-server entry resolution bug when ssr is disabled
  // (see @nuxt/vite-builder's ViteNodePlugin, which otherwise resolves the SSR entry
  // from the client build config).
  experimental: {
    viteEnvironmentApi: true,
  },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@nuxtjs/i18n',
  ],

  css: ['~/assets/css/main.css'],

  app: {
    head: {
      title: 'NestEgg · BudgetBooking',
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Quicksand:wght@500;600;700&family=Nunito+Sans:wght@400;600;700&display=swap',
        },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200',
        },
      ],
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
      googleClientId: process.env.NUXT_PUBLIC_GOOGLE_CLIENT_ID || '',
    },
  },

  i18n: {
    locales: [
      { code: 'zh-Hant', name: '繁體中文', file: 'zh-Hant.json' },
      { code: 'en', name: 'English', file: 'en.json' },
    ],
    defaultLocale: 'zh-Hant',
    langDir: 'locales/',
    strategy: 'no_prefix',
    detectBrowserLanguage: false,
  },
})
