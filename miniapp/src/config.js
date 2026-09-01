// 后端接口地址
// H5 预览：走 vite.config.js 中的代理，避免浏览器跨域
// 微信小程序：必须直连后端地址，开发者工具需勾选“不校验合法域名”
// 正式发布时请替换为已备案的 HTTPS 域名
export const BASE_URL = process.env.UNI_PLATFORM === 'h5' ? '' : 'http://localhost:8080'

export const TOKEN_KEY = 'token'
export const USER_KEY = 'user'
