import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import mixin from './mixin'

createApp(App).mixin(mixin).use(router).mount('#app')
