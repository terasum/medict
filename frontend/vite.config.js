import {defineConfig} from 'vite'
import { createVuePlugin as vue } from 'vite-plugin-vue2'
import Markdown from 'vite-plugin-md'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue({
    include: [/\.vue$/, /\.md$/],
  }), Markdown()],
  publicDir: "assets"
})