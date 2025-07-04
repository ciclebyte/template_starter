const { resolve } = require('path');
const vue = require('@vitejs/plugin-vue');
const AutoImport = require('unplugin-auto-import/vite');
const Components = require('unplugin-vue-components/vite');
const { NaiveUiResolver } = require('unplugin-vue-components/resolvers');
const { Languages, esbuildPluginMonacoEditorNls } = require('./plugins/nls/index.cjs');
const zh_hans = require('./plugins/nls/zh-hans.json');
const nlsPlugin = require('./plugins/nls/index.cjs').default;

// https://vite.dev/config/
module.exports = {
  plugins: [
    vue(),
    AutoImport({
      imports: ['vue', 'vue-router', 'pinia'],
      dts: 'src/auto-imports.d.ts',
      resolvers: [NaiveUiResolver()]
    }),
    Components({
      resolvers: [NaiveUiResolver()]
    }),
    ...(process.env.NODE_ENV !== 'development'
      ? [
          nlsPlugin({
            locale: Languages.zh_hans,
            localeData: zh_hans,
          }),
        ]
      : []),
  ],
  optimizeDeps: {
    esbuildOptions: {
      plugins: [
        esbuildPluginMonacoEditorNls({
          locale: Languages.zh_hans,
          localeData: zh_hans,
        }),
      ],
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
};
