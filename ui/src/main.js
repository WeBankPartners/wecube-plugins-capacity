import 'babel-polyfill'
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import '@/assets/css/local.bootstrap.css'
import 'bootstrap/dist/js/bootstrap.min.js'
import 'font-awesome/css/font-awesome.css'
import './plugins/iview.js'
import capacityRequestEntrance from '@/assets/js/capacityRequestEntrance.js'
import capacityApiCenter from '@/assets/config/capacity-api-center.json'
Vue.prototype.$capacityRequestEntrance = capacityRequestEntrance
Vue.prototype.capacityApiCenter = capacityApiCenter


Vue.config.productionTip = false

import VueI18n from 'vue-i18n'
import en from 'view-design/dist/locale/en-US'
import zh from 'view-design/dist/locale/zh-CN'

import en_local from '@/assets/locale/lang/capacity-en.json'
import zh_local from '@/assets/locale/lang/capacity-zh-CN.json'
Vue.use(VueI18n)
Vue.locale = () => {}
const messages = {
  'en-US': Object.assign(en_local, en),
  'zh-CN': Object.assign(zh_local, zh)
}
const i18n = new VueI18n({
  locale:
    localStorage.getItem('lang') ||
    (navigator.language || navigator.userLanguage === 'zh-CN'
      ? 'zh-CN'
      : 'en-US'),
  messages
})

new Vue({
  render: h => h(App),
  router,
  i18n
}).$mount('#app')
