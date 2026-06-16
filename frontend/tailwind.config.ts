import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      fontFamily: {
        display: ['Quicksand', 'sans-serif'],
        sans: ['"Nunito Sans"', 'sans-serif'],
      },
      colors: {
        background: '#FFF9F0',
        surface: '#FFFFFF',
        'surface-container': '#F3EDE4',
        'surface-container-low': '#F9F3EA',
        'surface-container-high': '#EDE7DF',
        ink: '#191919',
        'ink-soft': '#444651',
        primary: {
          DEFAULT: '#2D4B9B',
          dark: '#0E3383',
          light: '#AEC0FF',
        },
        secondary: {
          DEFAULT: '#F9BC00',
          dark: '#795900',
          light: '#FFDF9E',
        },
        mint: '#E8F5E9',
        lavender: '#F3E5F5',
        peach: '#FFF3E0',
        sky: '#E1F5FE',
        success: '#2E7D32',
        error: {
          DEFAULT: '#BA1A1A',
          container: '#FFDAD6',
        },
        outline: '#C4C6D3',
      },
      borderRadius: {
        card: '1.5rem',
        field: '1rem',
      },
      boxShadow: {
        soft: '0 10px 30px rgba(0,0,0,0.04)',
        'soft-primary': '0 10px 30px rgba(45,75,155,0.12)',
        'soft-lift': '0 16px 40px rgba(0,0,0,0.08)',
      },
    },
  },
}
