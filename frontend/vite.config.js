import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
// import Markdown from 'vite-plugin-md'

export default {
  publicDir: "assets",
  plugins: [
    vue({
      template: {
        compilerOptions: {
          compatConfig: {
            MODE: 2
          }
        }
      }
    })
  ]
}