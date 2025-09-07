// @ts-check
import { defineConfig } from "astro/config"

import react from "@astrojs/react"
import tailwindcss from "@tailwindcss/vite"

// https://astro.build/config
export default defineConfig({
  devToolbar: {
    enabled: false,
  },

  integrations: [react()],
  vite: {
    plugins: [tailwindcss()],
    server: {
      proxy: {
        "/api": "http://127.0.0.1:3000",
      },
    },
  },
})
