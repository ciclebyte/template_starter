import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import router from './router'
import naive from 'naive-ui'
import 'vfonts/Lato.css'
import 'vfonts/FiraCode.css'

// Monaco Editor Worker 配置
window.MonacoEnvironment = {
  getWorkerUrl: function (moduleId, label) {
    if (label === 'json') {
      return '/monaco/json.worker.js'
    }
    if (label === 'css' || label === 'scss' || label === 'less') {
      return '/monaco/css.worker.js'
    }
    if (label === 'html' || label === 'handlebars' || label === 'razor') {
      return '/monaco/html.worker.js'
    }
    if (label === 'typescript' || label === 'javascript') {
      return '/monaco/ts.worker.js'
    }
    return '/monaco/editor.worker.js'
  }
}

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(naive)
app.mount('#app')
