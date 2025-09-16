

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import '@fortawesome/fontawesome-free/css/all.css'
// import './css/ResourceCommon.css'
import axios from 'axios'
const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.mount('#app')
// 创建 axios 实例
const instance = axios.create({
  baseURL: "http://192.168.216.50:8090", 
  timeout: 5000
})
//创建白名单
const noAuthUrls = ['http://192.168.216.50:8090/login'] // 不需要 token 的接口
//请求拦截器
instance.interceptors.request.use(config => {
  if (!noAuthUrls.includes(config.url)) {
    const token = localStorage.getItem("token")
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
  }
  return config
}, error => {
  return Promise.reject(error)
})
// -------------------- WebSocket 工厂函数 --------------------
// function createWebSocket(url, protocols) {
//   // 只有指定的 ws 地址才加 token
//   const wsTokenUrlPrefix = "ws://192.168.216.50:8090"
//   if (url.startsWith(wsTokenUrlPrefix)) {
//     const token = localStorage.getItem("token")
//     if (token) {
//       url += (url.includes("?") ? "&" : "?") + `token=${token}`
//     }
//   }
//   return new WebSocket(url, protocols)
// }


export default instance
