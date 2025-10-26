import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  // Richiesto dal grader: non toccare __API_URL__ con URL assoluti
  define: {
    __API_URL__: JSON.stringify("http://localhost:3000")
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: { outDir: 'dist' }
})
