<script>
import { TOKEN_KEY } from './config'

export default {
  onLaunch() {
    const token = uni.getStorageSync(TOKEN_KEY)
    if (!token) {
      uni.reLaunch({ url: '/pages/login/login' })
    }
  }
}
</script>

<style>
page {
  background: linear-gradient(180deg, #fff1f6 0%, #eef4ff 60%, #f2f9ff 100%);
  font-family: -apple-system, "PingFang SC", "Microsoft YaHei", sans-serif;
  color: #3d3450;
}

.page {
  padding: 24rpx;
}

/* 卡片：圆润 + 入场动画 */
.card {
  background: rgba(255, 255, 255, 0.92);
  border-radius: 36rpx;
  padding: 34rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 16rpx 44rpx rgba(160, 130, 200, 0.16);
  animation: pop-in 0.42s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes pop-in {
  0% {
    opacity: 0;
    transform: translateY(24rpx) scale(0.96);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.title {
  font-size: 40rpx;
  font-weight: 700;
}

.subtitle {
  margin-top: 12rpx;
  font-size: 26rpx;
  color: #9188a8;
}

.label {
  font-size: 26rpx;
  color: #9188a8;
  margin-bottom: 12rpx;
}

.input {
  height: 92rpx;
  padding: 0 28rpx;
  border-radius: 26rpx;
  background: #f7f4ff;
  font-size: 30rpx;
  transition: background 0.25s ease;
}

.input:focus {
  background: #efe9ff;
}

/* 按钮：渐变 + 点击回弹 */
.btn {
  margin-top: 28rpx;
  height: 92rpx;
  line-height: 92rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #ff9ac1, #a78bfa);
  color: #ffffff;
  font-size: 30rpx;
  font-weight: 600;
  text-align: center;
  box-shadow: 0 12rpx 26rpx rgba(167, 139, 250, 0.32);
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.btn-hover {
  transform: scale(0.95);
  box-shadow: 0 6rpx 14rpx rgba(167, 139, 250, 0.26);
}

.btn-plain {
  background: #f3efff;
  color: #7c5cf0;
  box-shadow: none;
}

.btn-danger {
  background: #ffeef0;
  color: #f0698a;
  box-shadow: none;
}

.link {
  margin-top: 28rpx;
  text-align: center;
  font-size: 26rpx;
  color: #7c5cf0;
}

.link-hover {
  opacity: 0.6;
}

.empty {
  padding: 56rpx 0;
  text-align: center;
  font-size: 26rpx;
  color: #b3abc4;
}

/* 会呼吸的表情，用于品牌区和空态 */
.emoji-bounce {
  display: inline-block;
  animation: bounce 1.8s ease-in-out infinite;
}

@keyframes bounce {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10rpx);
  }
}

/* 加载中的转圈表情 */
.emoji-spin {
  display: inline-block;
  animation: spin 1.1s linear infinite;
}

@keyframes spin {
  100% {
    transform: rotate(360deg);
  }
}
</style>
