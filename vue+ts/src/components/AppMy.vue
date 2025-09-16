<template>
  <div class="my-page-container">
    <h2 class="title">
      <i class="fa fa-user icon"></i> 我的信息
    </h2>

    <!-- 用户信息 -->
    <div class="info-section">
      <p><strong>用户名：</strong> {{ username }}</p>
      <p><strong>权限：</strong> {{ permission }}</p>
    </div>

    <!-- 操作按钮 -->
    <div class="actions">
      <button class="btn btn-primary" @click="logout">退出登录</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: "AppMy",
  setup() {
    const router = useRouter()
    const username = ref('')
    const permission = ref('')

    onMounted(() => {
      // ✅ 从 localStorage 读取用户名，优先级高于 route.query
      username.value = localStorage.getItem('username') || '用户'

      // 判断权限
      if (username.value === 'admin') {
        permission.value = '全部权限'
      } else {
        permission.value = '查看权限'
      }
    })

    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      router.push({ path: '/' })
    }

    return { username, permission, logout }
  }
}
</script>

<style scoped>
.my-page-container {
  width: 95%;
  max-width: 800px;
  margin: 40px auto;
  padding: 25px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  font-family: 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 22px;
  font-weight: bold;
  color: #1e3c72;
  margin-bottom: 20px;
}

.title .icon {
  color: #42a5f5;
}

.info-section {
  font-size: 16px;
  color: #333;
  margin-bottom: 25px;
}

.info-section p {
  margin: 10px 0;
}

.actions {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary {
  background: #42a5f5;
  color: #fff;
}
</style>
