import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import WelcomeUser from '@/views/WelcomeUser.vue'

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

const routes = [
  { path: '/', name: 'login', component: Login },
  {
    path: '/welcome',
    name: 'welcome',
    component: WelcomeUser,
    children: [
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
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
