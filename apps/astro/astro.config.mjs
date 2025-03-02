// @ts-check
import react from '@astrojs/react'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'astro/config'
import tsConfigPaths from 'vite-tsconfig-paths'

// https://astro.build/config
export default defineConfig({
  vite: {
    plugins: [tailwindcss(), tsConfigPaths()],
  },

  integrations: [react()],
})
