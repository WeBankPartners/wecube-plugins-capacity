import axios from 'axios'
import {capacityURL_config} from './capacityURL'
export default function ajax (options) {
  const ajaxObj = {
    method: options.method,
    baseURL: capacityURL_config,
    url: options.url,
    timeout: 30000,
    params: options.params,
    // params: options.params || '',
    headers: {
      'Content-type': 'application/json;charset=UTF-8'
    },
    // data: JSON.stringify(options.data || '')
    data: JSON.stringify(options.data)
  }
  return window.request ? window.request(ajaxObj) : axios(ajaxObj)
}
