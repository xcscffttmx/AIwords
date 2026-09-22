// 后端接口地址
// H5 预览：走 vite.config.js 中的代理，避免浏览器跨域
// 微信小程序：必须直连后端地址，且需在微信后台配置合法域名（要求 HTTPS + 备案）
// 正式发布时把 PROD_DOMAIN 替换为你已备案的 HTTPS 域名
const PROD_DOMAIN = 'https://your-domain.com'
export const BASE_URL = process.env.UNI_PLATFORM === 'h5' ? '' : PROD_DOMAIN

export const TOKEN_KEY = 'token'
export const USER_KEY = 'user'
