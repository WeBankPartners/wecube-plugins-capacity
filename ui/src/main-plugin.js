import router from './router-plugin'
import '@/assets/css/local.bootstrap.css'
import 'bootstrap/dist/js/bootstrap.min.js'
import 'font-awesome/css/font-awesome.css'
import './plugins/iview.js'
import capacityRequestEntrance from '@/assets/js/capacityRequestEntrance.js'
import capacityApiCenter from '@/assets/config/capacity-api-center.json'

window.addOptions({
  $capacityRequestEntrance: capacityRequestEntrance,
  capacityApiCenter: capacityApiCenter,
})

window.addRoutes(router, 'capacity')

import en_local from '@/assets/locale/lang/capacity-en.json'
import zh_local from '@/assets/locale/lang/capacity-zh-CN.json'
window.locale('en-US',en_local)
window.locale('zh-CN',zh_local)
