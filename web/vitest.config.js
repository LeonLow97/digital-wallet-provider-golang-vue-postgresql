import { fileURLToPath } from 'node:url'
import { mergeConfig } from 'vite'
import { configDefaults, defineConfig } from 'vitest/config'
import viteConfig from './vite.config'
import vue from '@vitejs/plugin-vue'

export default mergeConfig(
  viteConfig,
  defineConfig({
    plugins: [], // removed Vue() plugin to resolve <template> and <script> Syntax Error
    test: {
      globals: true,
      environment: 'jsdom',
      exclude: [...configDefaults.exclude, 'e2e/*', 'packages/template/*'],
      root: fileURLToPath(new URL('./', import.meta.url)),
      transformMode: {
        web: [/\.[jt]sx$/]
      },
      deps: {
        inline: ['vuetify'] // added this to resolve file extension .css issue
      }
    },
    root: '.'
  })
)
