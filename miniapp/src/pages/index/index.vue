<template>
  <view class="page">
    <view class="card brand-card">
      <view class="brand-emoji emoji-bounce">🐣</view>
      <view class="title">智能查询</view>
      <view class="subtitle">你好，{{ username || '同学' }}，今天想学什么词呀 ~</view>
    </view>

    <view class="card">
      <view class="label">🔤 单词</view>
      <input class="input" v-model="word" placeholder="例如：innovation" maxlength="64" />

      <view class="label" style="margin-top: 24rpx">🤖 AI 提供商</view>
      <picker :range="providerLabels" :value="providerIndex" @change="onProviderChange">
        <view class="input picker-value">{{ providerLabels[providerIndex] }} ⌄</view>
      </picker>

      <view class="btn" hover-class="btn-hover" @tap="handleQuery">
        <text v-if="queryLoading" class="emoji-spin">🌀</text>
        <text v-else>🔍</text>
        {{ queryLoading ? ' 查询中...' : ' 查询单词' }}
      </view>
      <view class="btn btn-plain" hover-class="btn-hover" @tap="resetResult">🧹 清空结果</view>
    </view>

    <view class="card" v-if="result">
      <view class="word-title">
        <text class="title" style="font-size: 38rpx">{{ result.word }}</text>
        <text class="emoji-bounce" style="font-size: 34rpx">🌟</text>
      </view>
      <view class="subtitle">{{ sourceText }} · {{ providerText }}</view>

      <view class="label" style="margin-top: 24rpx">💡 释义</view>
      <view class="text">{{ result.definition || '暂无释义' }}</view>

      <view class="label" style="margin-top: 24rpx">✏️ 例句</view>
      <view v-if="examples.length">
        <view class="example" v-for="(item, index) in examples" :key="index">
          <view class="example-en">{{ index + 1 }}. {{ item.sentence }}</view>
          <view class="example-cn">{{ item.translation }}</view>
        </view>
      </view>
      <view class="empty" v-else>暂无例句 🍃</view>

      <view class="btn" v-if="result.source !== 'database'" hover-class="btn-hover" @tap="handleSave">
        <text v-if="saveLoading" class="emoji-spin">🌀</text>
        <text v-else>💖</text>
        {{ saveLoading ? ' 保存中...' : ' 保存到我的单词本' }}
      </view>
    </view>

    <view class="card empty" v-else>
      <view class="emoji-bounce" style="font-size: 64rpx">🔎</view>
      输入一个英文单词，开始今天的学习吧
    </view>

    <view class="card">
      <view class="btn btn-plain" hover-class="btn-hover" @tap="handleLogout">👋 退出登录</view>
    </view>
  </view>
</template>

<script>
import { queryWord, saveWord } from '../../api/index'
import { TOKEN_KEY, USER_KEY } from '../../config'

const PROVIDERS = [
  { label: 'DeepSeek', value: 'deepseek' },
  { label: '通义千问', value: 'qianwen' }
]

export default {
  data() {
    return {
      word: '',
      providerIndex: 0,
      providerLabels: PROVIDERS.map((item) => item.label),
      result: null,
      queryLoading: false,
      saveLoading: false,
      username: ''
    }
  },
  computed: {
    provider() {
      return PROVIDERS[this.providerIndex].value
    },
    providerText() {
      const found = PROVIDERS.find((item) => item.value === (this.result && this.result.ai_provider))
      return found ? found.label : 'AI'
    },
    sourceText() {
      return this.result && this.result.source === 'database' ? '来自我的单词本' : '来自 AI 查询'
    },
    examples() {
      const raw = this.result && this.result.examples
      if (Array.isArray(raw)) return raw
      if (!raw) return []
      try {
        const parsed = JSON.parse(raw)
        return Array.isArray(parsed) ? parsed : []
      } catch (e) {
        return []
      }
    }
  },
  onShow() {
    const user = uni.getStorageSync(USER_KEY)
    this.username = user && user.username ? user.username : ''
  },
  methods: {
    onProviderChange(e) {
      this.providerIndex = Number(e.detail.value)
    },
    async handleQuery() {
      if (this.queryLoading) return

      const word = this.word.trim()
      if (!word) {
        uni.showToast({ title: '请输入要查询的单词', icon: 'none' })
        return
      }
      if (!/^[a-zA-Z][a-zA-Z\-' ]*$/.test(word)) {
        uni.showToast({ title: '请输入英文单词', icon: 'none' })
        return
      }

      this.queryLoading = true
      try {
        this.result = await queryWord({ word, ai_provider: this.provider })
      } catch (e) {
        // 请求层已统一提示错误
      } finally {
        this.queryLoading = false
      }
    },
    async handleSave() {
      if (!this.result || this.saveLoading) return

      this.saveLoading = true
      try {
        const examples =
          typeof this.result.examples === 'string'
            ? this.result.examples
            : JSON.stringify(this.result.examples)

        await saveWord({
          word: this.result.word,
          definition: this.result.definition,
          examples,
          ai_provider: this.result.ai_provider
        })

        this.result.source = 'database'
        uni.showToast({ title: '收藏成功 💖' })
      } catch (e) {
        // 请求层已统一提示错误
      } finally {
        this.saveLoading = false
      }
    },
    resetResult() {
      this.result = null
      this.word = ''
    },
    handleLogout() {
      uni.removeStorageSync(TOKEN_KEY)
      uni.removeStorageSync(USER_KEY)
      uni.reLaunch({ url: '/pages/login/login' })
    }
  }
}
</script>

<style>
.brand-card {
  text-align: center;
  background: linear-gradient(135deg, rgba(255, 236, 214, 0.95), rgba(233, 226, 255, 0.95));
}

.brand-emoji {
  font-size: 78rpx;
}

.picker-value {
  line-height: 92rpx;
}

.word-title {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.text {
  font-size: 28rpx;
  line-height: 1.8;
}

.example {
  padding: 20rpx 0;
  border-top: 2rpx dashed #f0eafc;
  animation: pop-in 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.example-en {
  font-size: 28rpx;
  font-weight: 600;
}

.example-cn {
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #9188a8;
}
</style>
