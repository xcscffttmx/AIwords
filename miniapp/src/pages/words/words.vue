<template>
  <view class="page">
    <view class="card brand-card">
      <view class="brand-emoji emoji-bounce">🧸</view>
      <view class="title">我的单词本</view>
      <view class="subtitle">已收藏 {{ total }} 个小可爱 💫</view>
    </view>

    <view class="card" v-for="item in words" :key="item.id">
      <view class="word-head">
        <view class="title" style="font-size: 34rpx">{{ item.word }}</view>
        <view class="tag">{{ item.ai_provider === 'qianwen' ? '🌸 通义千问' : '🐳 DeepSeek' }}</view>
      </view>
      <view class="text">{{ item.definition }}</view>
      <view class="btn btn-danger" hover-class="btn-hover" @tap="confirmDelete(item)">🗑️ 删除</view>
    </view>

    <view class="card empty" v-if="!loading && !words.length">
      <view class="emoji-bounce" style="font-size: 64rpx">🌱</view>
      还没有收藏的单词，先去查询并保存几个吧
    </view>

    <view class="empty" v-if="loading"><text class="emoji-spin">🌀</text> 加载中...</view>
    <view class="empty" v-else-if="noMore && words.length">已经到底啦 🎀</view>
  </view>
</template>

<script>
import { deleteWord, getWords } from '../../api/index'

export default {
  data() {
    return {
      words: [],
      page: 1,
      pageSize: 10,
      total: 0,
      loading: false,
      noMore: false
    }
  },
  onShow() {
    this.reload()
  },
  onPullDownRefresh() {
    this.reload().finally(() => uni.stopPullDownRefresh())
  },
  onReachBottom() {
    if (this.noMore || this.loading) return
    this.page += 1
    this.loadWords()
  },
  methods: {
    reload() {
      this.page = 1
      this.noMore = false
      this.words = []
      return this.loadWords()
    },
    async loadWords() {
      this.loading = true
      try {
        const res = await getWords({ page: this.page, page_size: this.pageSize })
        const data = res.data || {}
        const list = data.list || []

        this.words = this.page === 1 ? list : this.words.concat(list)
        this.total = data.total || 0
        this.noMore = this.words.length >= this.total
      } catch (e) {
        // 请求层已统一提示错误
      } finally {
        this.loading = false
      }
    },
    confirmDelete(item) {
      uni.showModal({
        title: '要和这个单词说再见吗？🥺',
        success: (res) => {
          if (res.confirm) {
            this.handleDelete(item.id)
          }
        }
      })
    },
    async handleDelete(id) {
      try {
        await deleteWord(id)
        uni.showToast({ title: '已删除 🧹' })
        await this.reload()
      } catch (e) {
        // 请求层已统一提示错误
      }
    }
  }
}
</script>

<style>
.brand-card {
  text-align: center;
  background: linear-gradient(135deg, rgba(255, 224, 236, 0.95), rgba(214, 240, 255, 0.95));
}

.brand-emoji {
  font-size: 78rpx;
}

.word-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tag {
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  background: #f3efff;
  color: #7c5cf0;
  font-size: 22rpx;
}

.text {
  margin-top: 16rpx;
  font-size: 28rpx;
  line-height: 1.7;
  color: #5b5370;
}
</style>
