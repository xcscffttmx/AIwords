import { BASE_URL, TOKEN_KEY, USER_KEY } from '../config'

function toLogin() {
  uni.removeStorageSync(TOKEN_KEY)
  uni.removeStorageSync(USER_KEY)
  uni.reLaunch({ url: '/pages/login/login' })
}

export default function request(options) {
  const token = uni.getStorageSync(TOKEN_KEY)

  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + '/api' + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: token ? { Authorization: 'Bearer ' + token } : {},
      timeout: 30000,
      success(res) {
        if (res.statusCode === 401) {
          uni.showToast({ title: '登录已过期', icon: 'none' })
          toLogin()
          reject(res)
          return
        }

        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)
          return
        }

        const message = (res.data && (res.data.error || res.data.message)) || '请求失败，请稍后重试'
        uni.showToast({ title: message, icon: 'none' })
        reject(res)
      },
      fail(err) {
        uni.showToast({ title: '网络异常，请检查后端服务', icon: 'none' })
        reject(err)
      }
    })
  })
}
