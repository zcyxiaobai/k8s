<template>
  <div class="welcome-container">
    <!-- 顶部平台名 -->
    <div class="app-name">
      <span v-if="!isCollapsed">k8s资源监控管理平台</span>
      <span v-else>K8S</span>
    </div>

    <!-- 顶部 header -->
    <div class="header">
      <div class="header-img">
        <img src="/img.png" alt="">
      </div>
      <div class="header-night">
        <div class="name">
          <p>欢迎您, {{ username }}</p>
        </div>
        <div class="divider">|</div>
        <div class="tuichu">
          <button @click="logout">退出登录</button>
        </div>
      </div>
    </div>

    <!-- 主体内容 -->
    <div class="banner">
      <!-- 左侧功能菜单 -->
      <div class="function" :class="{ collapsed: isCollapsed }">
        <div class="collapse-btn" @click="toggleCollapse">
          <span v-if="!isCollapsed">◀</span>
          <span v-else>▶</span>
        </div>
        <div class="function-gb" v-if="!isCollapsed">功能选择</div>

        <!-- 大屏展示 -->
        <div class="menu-item" @click="navigate('dashboard')" :title="isCollapsed ? '大屏展示' : ''">
          <span v-if="!isCollapsed">大屏展示</span>
          <span v-else>D</span>
        </div>

        <!-- 资源清单 -->
        <div class="menu-item" :title="isCollapsed ? '资源清单' : ''">
          <span v-if="!isCollapsed">资源清单</span>
          <span v-else>R</span>
          <el-select v-if="!isCollapsed" v-model="selectedResource" placeholder="请选择资源" size="small" @change="onSelect('Resource')">
            <el-option label="Pod" value="pod"></el-option>
            <el-option label="Deployment" value="deployment"></el-option>
            <el-option label="Service" value="service"></el-option>
            <el-option label="Ingress" value="ingress"></el-option>
          </el-select>
        </div>

        <!-- 创建资源 -->
        <div class="menu-item" :title="isCollapsed ? '创建资源' : ''">
          <span v-if="!isCollapsed">资源操作</span>
          <span v-else>C</span>
          <el-select v-if="!isCollapsed" v-model="selectedCreate" placeholder="请选择资源" size="small" @change="onSelect('Create')" :disabled="!isAdmin">
            <el-option label="Pod" value="pod"></el-option>
            <el-option label="Deployment" value="deployment"></el-option>
            <el-option label="Service" value="service"></el-option>
            <el-option label="Ingress" value="ingress"></el-option>
          </el-select>
        </div>

        <!-- CRD资源 -->
        <div class="menu-item" :title="isCollapsed ? 'CRD资源' : ''">
          <span v-if="!isCollapsed">CRD资源</span>
          <span v-else>CRD</span>
          <el-select v-if="!isCollapsed" v-model="selectedCRD" placeholder="请选择操作" size="small" @change="onSelect('CRD')">
            <el-option label="查看" value="view"></el-option>
            <el-option label="创建" value="create" :disabled="!isAdmin"></el-option>
          </el-select>
        </div>

        <!-- 我的 -->
        <div class="menu-item" @click="navigate('my')" :title="isCollapsed ? '我的' : ''">
          <span v-if="!isCollapsed">我的</span>
          <span v-else>M</span>
        </div>
      </div>

      <!-- 右侧展示区 -->
      <div class="Resource">
        <router-view />
      </div>
    </div>

    <!-- token 显示 -->
    <!-- <p class="token-display">
      您的 token 是:
      <el-input v-model="token" readonly></el-input>
    </p> -->
  </div>
</template>

<script>
/* eslint-disable-next-line no-unused-vars */
import { ref, onMounted, nextTick,computed  } from 'vue'
// import { useRoute, useRouter } from 'vue-router'
import { useRouter } from 'vue-router'
import { onUnmounted } from 'vue'

export default {
  name: "WelcomeUser",
  setup() {
    // const route = useRoute()
    const router = useRouter()
    const username = ref('')
    const token = ref('')
    const selectedResource = ref('')
    const selectedCreate = ref('')
    const selectedCRD = ref('')
    const isCollapsed = ref(false)

    const isAdmin = computed(() => username.value === 'admin')

    const toggleCollapse = () => {
      isCollapsed.value = !isCollapsed.value
    }


    const navigate = (feature) => {
      router.push({ name: feature })
      selectedResource.value = ''
      selectedCreate.value = ''
      selectedCRD.value = ''
    }

    const onSelect = (type) => {
      switch(type){
        case 'Resource':
          if(selectedResource.value) router.push({ name: selectedResource.value })
          selectedCreate.value = ''
          selectedCRD.value = ''
          break
        case 'Create':
          if(selectedCreate.value) router.push({ name: `create-${selectedCreate.value}` })
          selectedResource.value = ''
          selectedCRD.value = ''
          break
        case 'CRD':
          if(selectedCRD.value) router.push({ name: `crd-${selectedCRD.value}` })
          selectedResource.value = ''
          selectedCreate.value = ''
          break
      }
    }

    const logout = () => {
      localStorage.removeItem('token')
      router.push({ path: '/' })
    }

    onMounted(() => {
      // username.value = route.query.username || '用户'
        // ✅ 从 localStorage 读取用户名，优先级高于 route.query
      username.value = localStorage.getItem('username') || '用户'
      token.value = localStorage.getItem('token') || ''

      // ✅ 全局屏蔽 ResizeObserver 警告
      const resizeObserverErr = /ResizeObserver loop limit exceeded|ResizeObserver loop completed with undelivered notifications/;
      const errorHandler = (e) => {
        if (resizeObserverErr.test(e.message)) {
          e.stopImmediatePropagation();
          console.warn('ResizeObserver warning ignored.')
        }
      }
        // 监听 Pod 修改事件
  const clearDropdownHandler = () => {
    selectedResource.value = ''
    selectedCreate.value = ''
    selectedCRD.value = ''
  }
      window.addEventListener('error', errorHandler)
      window.addEventListener('clear-welcome-dropdowns', clearDropdownHandler)

      // 组件卸载时移除监听
      const cleanup = () => window.removeEventListener('error', errorHandler)
      window.addEventListener('beforeunload', cleanup)
       // 组件卸载时移除监听
  onUnmounted(() => {
    window.removeEventListener('clear-welcome-dropdowns', clearDropdownHandler)
  })
    })
  
    return { username, token, selectedResource, selectedCreate, selectedCRD, isCollapsed, toggleCollapse, navigate, onSelect, logout,isAdmin }
  }
}
</script>

<style scoped>
/* 样式保持原样 */
.welcome-container { font-family: 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5; min-height: 100vh; }
.app-name { width: 100%; height: 60px; background-color: #fff; display: flex; justify-content: center; align-items: center; color: #1e3c72; font-weight: bold; font-size: 22px; letter-spacing: 2px; border-bottom: 1px solid #eee; }
.header { display: flex; width: 95%; max-width: 1600px; height: 60px; background-color: #fff; justify-content: space-between; align-items: center; margin: 20px auto 0; padding: 0 20px; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,0.05); }
.header-img { width: 50px; height: 50px; background-color: #6149c8; border-radius: 50%; overflow: hidden; }
.header-img img {
  width: 100%;   /* 图片宽度自适应容器 */
  height: 100%;  /* 图片高度自适应容器 */
  object-fit: cover; /* 保持比例裁剪填满容器 */
}
.header-night { display: flex; align-items: center; gap: 10px; }
.header-night .name { font-size: 17px; color: #333; }
.header-night .divider { font-size: 20px; color: #aaa; }
.header-night .tuichu button { all: unset; cursor: pointer; color: #ff4d4f; font-weight: bold; font-size: 16px; }
.banner { width: 95%; max-width: 1600px; display: flex; gap: 20px; margin: 20px auto; min-height: 800px; }
.function { width: 280px; background-color: #fff; border-radius: 10px; padding: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); display: flex; flex-direction: column; gap: 20px; transition: width 0.3s ease; overflow: hidden; }
.function.collapsed { width: 60px; }
.collapse-btn { cursor: pointer; margin-bottom: 10px; font-weight: bold; color: #1e3c72; text-align: center; }
.function-gb { font-weight: bold; font-size: 18px; color: #1e3c72; }
.menu-item { font-size: 16px; color: #333; cursor: pointer; display: flex; flex-direction: column; }
.menu-item:hover { color: #1e3c72; }
.el-select { width: 100%; margin-top: 5px; cursor: pointer; }
.Resource { flex: 1; background-color: #fff; padding: 25px; border-radius: 10px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); transition: all 0.3s ease; }
.token-display { width: 95%; max-width: 1600px; margin: 20px auto; font-size: 16px; }
</style>
