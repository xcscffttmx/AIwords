<template>
  <view class="page">
    <view class="card brand-card">
      <view class="brand-emoji emoji-bounce">🌱</view>
      <view class="title">创建账号</view>
      <view class="subtitle">注册后即可拥有自己的单词小窝 🏡</view>
    </view>

    <view class="card">
      <view class="label">👤 用户名</view>
      <input class="input" v-model="form.username" placeholder="3-50 个字符" />

      <view class="label" style="margin-top: 24rpx">🔒 密码</view>
      <input class="input" v-model="form.password" password placeholder="至少 6 位" />

      <view class="label" style="margin-top: 24rpx">✅ 确认密码</view>
      <input class="input" v-model="form.confirmPassword" password placeholder="请再次输入密码" />

      <view class="btn" hover-class="btn-hover" @tap="handleRegister">
        <text v-if="loading" class="emoji-spin">🌀</text>
        <text v-else>🎈</text>
        {{ loading ? ' 注册中...' : ' 注册' }}
      </view>
      <view class="link" hover-class="link-hover" @tap="goLogin">已有账号？返回登录 👋</view>
    </view>
  </view>
</template>

<script>
import { register } from '../../api/index'

export default {
  data() {
    return {
      loading: false,
      form: {
        username: '',
        password: '',
        confirmPassword: ''
      }
    }
  },
  methods: {
    async handleRegister() {
      if (this.loading) return

      const username = this.form.username.trim()
      const { password, confirmPassword } = this.form

      if (username.length < 3 || username.length > 50) {
        uni.showToast({ title: '用户名长度应为 3-50 个字符', icon: 'none' })
        return
      }
      if (password.length < 6) {
        uni.showToast({ title: '密码至少 6 位', icon: 'none' })
        return
      }
      if (password !== confirmPassword) {
        uni.showToast({ title: '两次输入的密码不一致', icon: 'none' })
        return
      }

      this.loading = true
      try {
        await register({ username, password })
        uni.showToast({ title: '注册成功，请登录 🎉' })
        setTimeout(() => {
          uni.reLaunch({ url: '/pages/login/login' })
        }, 800)
      } catch (e) {
        // 请求层已统一提示错误
      } finally {
        this.loading = false
      }
    },
    goLogin() {
      uni.reLaunch({ url: '/pages/login/login' })
    }
  }
}
</script>

<style>
.brand-card {
  text-align: center;
  background: linear-gradient(135deg, rgba(214, 245, 227, 0.95), rgba(226, 236, 255, 0.95));
}

.brand-emoji {
  font-size: 78rpx;
}
</style>
