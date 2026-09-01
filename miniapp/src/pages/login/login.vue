<template>
  <view class="page">
    <view class="card brand-card">
      <view class="brand-emoji emoji-bounce">📚</view>
      <view class="title">AI 智能单词本</view>
      <view class="subtitle">登录后开始查询和收藏单词 ✨</view>
    </view>

    <view class="card">
      <view class="label">👤 用户名</view>
      <input class="input" v-model="form.username" placeholder="请输入用户名" />

      <view class="label" style="margin-top: 24rpx">🔒 密码</view>
      <input class="input" v-model="form.password" password placeholder="请输入密码" />

      <view class="btn" hover-class="btn-hover" @tap="handleLogin">
        <text v-if="loading" class="emoji-spin">🌀</text>
        <text v-else>🚀</text>
        {{ loading ? ' 登录中...' : ' 登录' }}
      </view>
      <view class="link" hover-class="link-hover" @tap="goRegister">还没有账号？立即注册 💖</view>
    </view>
  </view>
</template>

<script>
import { login } from '../../api/index'
import { TOKEN_KEY, USER_KEY } from '../../config'

export default {
  data() {
    return {
      loading: false,
      form: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    async handleLogin() {
      if (this.loading) return

      const username = this.form.username.trim()
      const password = this.form.password

      if (!username || !password) {
        uni.showToast({ title: '请输入用户名和密码', icon: 'none' })
        return
      }

      this.loading = true
      try {
        const res = await login({ username, password })
        uni.setStorageSync(TOKEN_KEY, res.token)
        uni.setStorageSync(USER_KEY, res.user)
        uni.showToast({ title: '登录成功 🎉' })
        uni.reLaunch({ url: '/pages/index/index' })
      } catch (e) {
        // 请求层已统一提示错误
      } finally {
        this.loading = false
      }
    },
    goRegister() {
      uni.navigateTo({ url: '/pages/register/register' })
    }
  }
}
</script>

<style>
.brand-card {
  text-align: center;
  background: linear-gradient(135deg, rgba(255, 224, 236, 0.95), rgba(233, 226, 255, 0.95));
}

.brand-emoji {
  font-size: 78rpx;
}
</style>
