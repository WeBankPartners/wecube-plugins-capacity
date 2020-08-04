import 'babel-polyfill'
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import '@/assets/css/local.bootstrap.css'
import 'bootstrap/dist/js/bootstrap.min.js'
import 'font-awesome/css/font-awesome.css'
import './plugins/iview.js'
import httpRequestEntrance from '@/assets/js/httpRequestEntrance.js'
import jquery from 'jquery'
import { tableUtil } from '@/assets/js/tableUtil.js'
import { validate } from '@/assets/js/validate.js'
import VeeValidate from '@/assets/veeValidate/VeeValidate'
import apiCenter from '@/assets/config/api-center.json'
const eventBus = new Vue()
Vue.prototype.$eventBus = eventBus
Vue.prototype.$httpRequestEntrance = httpRequestEntrance
Vue.prototype.JQ = jquery
Vue.prototype.$validate = validate
Vue.prototype.$tableUtil = tableUtil
Vue.prototype.apiCenter = apiCenter

import ModalComponent from '@/components/modal'
Vue.component('ModalComponent', ModalComponent)
Vue.use(VeeValidate)

Vue.config.productionTip = false

import VueI18n from 'vue-i18n'
import en from 'view-design/dist/locale/en-US'
import zh from 'view-design/dist/locale/zh-CN'

import en_local from '@/assets/locale/lang/en.json'
import zh_local from '@/assets/locale/lang/zh-CN.json'
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
