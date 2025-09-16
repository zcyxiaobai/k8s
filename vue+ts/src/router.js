import { createRouter, createWebHistory } from 'vue-router'
import LoginUser from './views/LoginUser.vue'
import WelcomeUser from './views/WelcomeUser.vue'

import AppDashboard from '@/components/AppDashboard.vue'
import AppMy from '@/components/AppMy.vue'

// 资源
import AppPod from '@/components/resources/AppPod.vue'
import AppDeployment from '@/components/resources/AppDeployment.vue'
import AppService from '@/components/resources/AppService.vue'
import AppIngress from '@/components/resources/AppIngress.vue'

// 创建
import CreatePod from '@/components/create/AppPod.vue'
import CreateDeployment from '@/components/create/AppDeployment.vue'
import CreateService from '@/components/create/AppService.vue'
import CreateIngress from '@/components/create/AppIngress.vue'

// CRD
import ViewCRD from '@/components/crd/AppViewCRD.vue'
import CreateCRD from '@/components/crd/AppCreateCRD.vue'

// 路由配置
const routes = [
  { path: '/', redirect: '/login' }, // 默认进入登录页
  { path: '/login', component: LoginUser },
  // { 
  //   path: '/welcome', 
  //   component: WelcomeUser,
  //   meta: { requiresAuth: true } // 标记需要鉴权
  // }
  {
      path: '/welcome',
      name: 'welcome',
      component: WelcomeUser,
      children: [
        { path: '', redirect: '/welcome/dashboard' },// 默认进入 dashboard
        { path: 'dashboard', name: 'dashboard', component: AppDashboard },
        { path: 'my', name: 'my', component: AppMy },
  
        // 资源
        { path: 'pod', name: 'pod', component: AppPod },
        { path: 'deployment', name: 'deployment', component: AppDeployment },
        { path: 'service', name: 'service', component: AppService },
        { path: 'ingress', name: 'ingress', component: AppIngress },
  
        // 创建
        { path: 'create-pod', name: 'create-pod', component: CreatePod },
        { path: 'create-deployment', name: 'create-deployment', component: CreateDeployment },
        { path: 'create-service', name: 'create-service', component: CreateService },
        { path: 'create-ingress', name: 'create-ingress', component: CreateIngress },
  
        // CRD
        { path: 'crd-view', name: 'crd-view', component: ViewCRD },
        { path: 'crd-create', name: 'crd-create', component: CreateCRD },
        //修改的路由
        { path: 'update/pod',
          name: 'UpdatePod',
          component: () => import('@/components/update/AppPod.vue'),
          props: route => {
             const podStr = route.query.pod || null;
             let podData = null;
             if (podStr) {
             try { podData = JSON.parse(decodeURIComponent(podStr)); } 
                  catch(e){ console.error('解析 Pod 数据失败', e); }
              }
         return { podData };
         }
        },
        {
          path: 'update/deployment',
          name: 'UpdateDeployment',
          component: () => import('@/components/update/AppDeployment.vue'),
          props: route => {
              const depStr = route.query.deployment || null;
              let depData = null;
              if (depStr) {
                  try { depData = JSON.parse(decodeURIComponent(depStr)); }
                  catch(e){ console.error('解析 Deployment 数据失败', e); }
              }
               return { depData };
              }
        },
        // Service
{
  path: 'update/service',
  name: 'UpdateService',
  component: () => import('@/components/update/AppService.vue'),
  props: route => {
    const svcStr = route.query.service || null;
    let svcData = null;
    if (svcStr) {
      try { svcData = JSON.parse(decodeURIComponent(svcStr)); }
      catch(e){ console.error('解析 Service 数据失败', e); }
    }
    return { svcData };
  }
},
// Ingress
{
  path: 'update/ingress',
  name: 'UpdateIngress',
  component: () => import('@/components/update/AppIngress.vue'),
  props: route => {
    const ingStr = route.query.ingress || null;
    let ingData = null;
    if (ingStr) {
      try { ingData = JSON.parse(decodeURIComponent(ingStr)); }
      catch(e){ console.error('解析 Ingress 数据失败', e); }
    }
    return { ingData };
  }
},
{
  path: 'update/bookapp',
  name: 'UpdateBookapp',
  component: () => import('@/components/update/AppBook.vue'),
  props: route => {
    const itemStr = route.query.bookapp || null;
    let itemData = null;
    if (itemStr) {
      try { itemData = JSON.parse(decodeURIComponent(itemStr)); } 
      catch(e){ console.error('解析 Bookapp 数据失败', e); }
    }
    return { bookappData: itemData };
  }
}
      ],
      meta: { requiresAuth: true } // 标记需要鉴权
  },
  // //修改的路由
  // {
  //    path: '/update/pod',
  //    name:'UpdatePod',
  //    component: () => import('@/components/update/AppPod.vue'),
  //     props: route => ({
  //   podData: route.state?.podData || null,  // 使用 router.push 的 state
  //   namespace: route.query.namespace || '',
  //   name: route.query.name || ''
  // })
  // }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：检查 token
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
