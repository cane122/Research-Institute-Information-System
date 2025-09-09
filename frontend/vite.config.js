import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    rollupOptions: {
      external: ['./wailsjs/go/models', './wailsjs/go/main/App', './wailsjs/runtime/runtime'],
      output: {
        globals: {
          './wailsjs/go/models': 'models',
          './wailsjs/go/main/App': 'App',
          './wailsjs/runtime/runtime': 'runtime'
        }
      }
    }
  },
  server: {
    fs: {
      allow: ['..']
    }
  }
})
