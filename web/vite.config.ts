import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { visualizer } from 'rollup-plugin-visualizer'

function filterWarningsPlugin() {
  return {
    name: 'filter-warnings',
    configResolved(config: Record<string, unknown>) {
      const originalWarn = config.logger.warn
      config.logger.warn = (msg: string, options?: Record<string, unknown>) => {
        if (msg.includes('A nested style rule cannot start with')) {
          return
        }
        if (msg.includes('css-syntax-error')) {
          return
        }
        if (msg.includes('is dynamically imported by') && msg.includes('but also statically imported by')) {
          return
        }
        if (msg.includes('dynamic import will not move module into another chunk')) {
          return
        }
        if (msg.includes('Default import of CSS without') && msg.includes('?inline')) {
          return
        }
        originalWarn(msg, options)
      }
    },
  }
}

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const envDir = __dirname
  const env = loadEnv(mode, envDir, '')

  return {
    define: {
      __VITE_API_BASE_URL__: JSON.stringify(env.VITE_API_BASE_URL || '/api/v1'),
    },
    envDir,
    base: '/',
    plugins: [
      vue(),
      filterWarningsPlugin(),
      AutoImport({
        imports: ['vue', 'vue-router', 'pinia', '@vueuse/core'],
        dts: true,
        eslintrc: {
          enabled: true,
        },
      }),
      Components({
        dts: true,
        include: [],
        resolvers: [],
      }),
      visualizer({
        open: false,
        filename: 'dist/stats.html',
        gzipSize: true,
        brotliSize: false,
      }),
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
      },
    },
    server: {
      port: 3800,
      host: true,
    },
    optimizeDeps: {
      include: [
        'gsap',
        'gsap/ScrollTrigger',
        'comlink',
        'spark-md5',
        'vue',
        '@vue/runtime-dom',
        '@vue/runtime-core',
        '@vue/reactivity',
        '@vue/shared',
        'md-editor-v3',
      ],
    },
    worker: {
      format: 'es',
      plugins: [vue()],
    },
    esbuild: {
      target: 'es2020',
      minifyIdentifiers: false,
      minifySyntax: true,
      minifyWhitespace: true,
    },
    css: {
      preprocessorOptions: {
        sass: {
          api: 'modern',
          silenceDeprecations: ['legacy-js-api'],
          quietDeps: true,
        },
        scss: {
          api: 'modern',
          silenceDeprecations: ['legacy-js-api'],
          quietDeps: true,
        },
      },
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
      target: 'es2020',
      cssTarget: 'chrome61',
      reportCompressedSize: true,
      chunkSizeWarningLimit: 500,
      cssMinify: 'esbuild',
      rollupOptions: {
        maxParallelFileOps: 5,
        onwarn(warning, warn) {
          if (warning.code === 'DYNAMIC_IMPORT_WILL_NOT_MOVE_MODULE') {
            return
          }
          if (warning.code === 'CIRCULAR_DEPENDENCY') {
            return
          }
          warn(warning)
        },
        output: {
          manualChunks(id) {
            if (id.includes('node_modules')) {
              if (id.includes('node_modules/vue/') || id.includes('node_modules/@vue/')) {
                return 'vue-vendor'
              }
              if (id.includes('node_modules/vue-router')) {
                return 'vue-vendor'
              }
              if (id.includes('node_modules/pinia')) {
                return 'vue-vendor'
              }
              if (id.includes('node_modules/@vueuse/')) {
                return 'vueuse'
              }

              if (id.includes('echarts')) {
                return 'echarts'
              }
              if (id.includes('highlight.js')) {
                if (id.includes('highlight.js/lib/languages/')) {
                  const match = id.match(/languages\/([^/]+)/)
                  if (match) {
                    const commonLangs = ['javascript', 'typescript', 'python', 'java', 'css', 'html', 'json', 'markdown']
                    const lang = match[1].replace('.js', '')
                    if (commonLangs.includes(lang)) {
                      return 'highlight-common'
                    }
                    return 'highlight-other'
                  }
                }
                return 'highlight-core'
              }
              if (id.includes('quill')) {
                return 'editor-quill'
              }
              if (id.includes('@vueup/vue-quill')) {
                return 'vue-quill'
              }
              if (id.includes('md-editor-v3')) {
                return 'editor-markdown'
              }
              if (id.includes('pinyin-pro')) {
                return 'pinyin'
              }
              if (id.includes('axios')) {
                return 'axios'
              }
              if (id.includes('spark-md5')) {
                return 'spark-md5'
              }
              if (id.includes('comlink')) {
                return 'comlink'
              }
              if (id.includes('jszip')) {
                return 'jszip'
              }
              if (id.includes('gsap')) {
                return 'gsap'
              }
              if (id.includes('date-fns')) {
                return 'date-fns'
              }
              if (id.includes('@tanstack/vue-table')) {
                return 'vue-table'
              }
              if (id.includes('vue-easy-lightbox')) {
                return 'lightbox'
              }
              if (id.includes('marked')) {
                return 'marked'
              }
              if (id.includes('dompurify')) {
                return 'dompurify'
              }
              if (id.includes('html-entities')) {
                return 'html-entities'
              }
              if (id.includes('vue-virtual-scroller')) {
                return 'virtual-scroller'
              }
              if (id.includes('vue-draggable')) {
                return 'draggable'
              }
            }
          },
          chunkFileNames: 'js/[name]-[hash].js',
          entryFileNames: 'js/[name]-[hash].js',
          assetFileNames: (assetInfo) => {
            const info = assetInfo.name.split('.')
            const ext = info[info.length - 1]
            if (/\.(png|jpe?g|gif|svg|ico)$/i.test(assetInfo.name)) {
              return `images/[name]-[hash].${ext}`
            }
            if (/\.(woff2?|eot|ttf|otf)$/i.test(assetInfo.name)) {
              return `fonts/[name]-[hash].${ext}`
            }
            return `assets/[name]-[hash].${ext}`
          },
        },
      },
    },
  }
})
