import router from './router-plugin'
import '@/assets/css/local.bootstrap.css'
import 'bootstrap/dist/js/bootstrap.min.js'
import 'font-awesome/css/font-awesome.css'
import './plugins/iview.js'
import capacityRequestEntrance from '@/assets/js/capacityRequestEntrance.js'
import {tableUtil} from '@/assets/js/tableUtil.js'
import {validate} from '@/assets/js/validate.js'
import capacityApiCenter from '@/assets/config/api-center.json'

window.addOptions({
  $capacityRequestEntrance: capacityRequestEntrance,
  $validate: validate,
  $tableUtil: tableUtil,
  capacityApiCenter: capacityApiCenter,
})

window.addRoutes(router, 'capacity')

import en_local from '@/assets/locale/lang/en.json'
import zh_local from '@/assets/locale/lang/zh-CN.json'
window.locale('en-US',en_local)
window.locale('zh-CN',zh_local)
