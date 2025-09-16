<template>
  <div class="login-wrapper">
    <el-card class="login-card">
      <h2 class="title">用户登录</h2>
      <el-form :model="form" ref="loginForm" class="login-form">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" prefix-icon="el-icon-user" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input type="password" v-model="form.password" placeholder="请输入密码" prefix-icon="el-icon-lock" />
        </el-form-item>
        <el-form-item>
          <el-row :gutter="10">
            <el-col :span="14">
              <el-input v-model="form.captcha" placeholder="请输入验证码" />
            </el-col>
            <el-col :span="10">
              <div class="captcha-number" @click="refreshCaptcha">{{ captcha }}</div>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="login-btn" @click="handleLogin" style="width:100%">登录</el-button>
        </el-form-item>
        <el-alert v-if="errorMsg" :title="errorMsg" type="error" show-icon />
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

export default {
  name: "LoginUser",
  setup() {
    const router = useRouter()
    const form = ref({
      username: '',
      password: '',
      captcha: ''
    })
    const errorMsg = ref('')
    const captcha = ref('')

    const generateCaptcha = () => {
      // 生成 4 位随机数字
      captcha.value = Math.floor(1000 + Math.random() * 9000).toString()
    }

    const refreshCaptcha = () => generateCaptcha()
    generateCaptcha() // 初始化验证码

    const handleLogin = async () => {
      if (!form.value.username || !form.value.password || !form.value.captcha) {
        errorMsg.value = '请填写完整信息'
        return
      }

      // 前端校验验证码
      if (form.value.captcha !== captcha.value) {
        errorMsg.value = '验证码错误'
        refreshCaptcha()
        form.value.captcha = ''
        return
      }

      try {
        const res = await axios.post('http://192.168.216.50:8090/login', {
          username: form.value.username,
          password: form.value.password
        })
        if (res.status === 200 && res.data.message === '登陆成功') {
          // 保存 token 到 localStorage
          localStorage.setItem('token', res.data.token) // 假设 token 以 token 字段返回
            localStorage.setItem('username', res.data.username) // ✅ 新增
          // 跳转到 welcome 页面
          router.push({ path: '/welcome', query: { username: res.data.username } })
        } else {
          errorMsg.value = res.data.message || '登录失败'
          refreshCaptcha()
          form.value.captcha = ''
        }
      } catch (err) {
        errorMsg.value = err.response?.data?.message || '登录异常'
        refreshCaptcha()
        form.value.captcha = ''
      }
    }

    return { form, errorMsg, captcha, refreshCaptcha, handleLogin }
  }
}
</script>


<style scoped>
.login-wrapper {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5; /* 浅灰白背景 */
}

.login-card {
  width: 360px;
  padding: 30px 20px;
  border-radius: 12px; /* 圆角更柔和 */
  background-color: #fff; /* 卡片白色 */
  box-shadow: 0 4px 12px rgba(0,0,0,0.08); /* 柔和阴影 */
}

.title {
  text-align: center;
  margin-bottom: 25px;
  color: #1e3c72; /* 与系统主色统一 */
  font-weight: 600;
  font-size: 22px;
}

.el-input .el-input__inner {
  border-radius: 8px;
  border: 1px solid #ddd;
  padding: 6px 10px;
  background-color: #fafafa;
}

.captcha-number {
  height: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #fafafa;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-weight: bold;
  cursor: pointer;
  user-select: none;
  color: #555;
}

.login-btn {
  margin-top: 10px;
  border-radius: 8px;
  background-color: #42a5f5; /* 系统主色 */
  color: #fff;
  font-weight: 500;
}

.el-alert {
  border-radius: 6px;
  font-size: 14px;
}
</style>


