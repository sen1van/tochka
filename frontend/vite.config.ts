import { defineConfig } from 'vite'
import { viteSingleFile } from 'vite-plugin-singlefile'
import vue from '@vitejs/plugin-vue'
import vueRouter from 'vue-router/vite'
import ui from '@nuxt/ui/vite'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: '../builds'
  },

  plugins: [
    viteSingleFile(),
    vueRouter({
      dts: 'src/route-map.d.ts'
    }),
    vue(),
    ui({
      ui: {
        colors: {
          primary: 'black',
          neutral: 'zinc'
        }
      }
    })
  ]
})
