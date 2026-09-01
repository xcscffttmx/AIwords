# AI 智能单词本 · 小程序端

基于 uni-app（Vue3）实现的微信小程序端，复用现有 Go + Gin 后端接口，功能与 Web 端保持一致。

## 功能

- 注册、登录，token 存储在本地 Storage
- 输入英文单词，选择 DeepSeek 或通义千问进行查询
- 查询结果手动保存到个人单词本
- 单词本列表支持下拉刷新和触底加载更多
- 删除单词记录

## 目录结构

```text
miniapp/
├── src/
│   ├── api/index.js         # 接口封装
│   ├── utils/request.js     # 请求封装与鉴权、401 处理
│   ├── config.js            # 后端地址与 Storage key
│   ├── pages/
│   │   ├── login/           # 登录
│   │   ├── register/        # 注册
│   │   ├── index/           # 单词查询
│   │   └── words/           # 我的单词本
│   ├── App.vue
│   ├── main.js
│   ├── pages.json
│   └── manifest.json
└── vite.config.js
```

## 后端地址配置

小程序不能使用 Vite Proxy，必须直连后端地址。请在 `src/config.js` 中修改：

```js
export const BASE_URL = 'http://localhost:8080'
```

正式发布时需替换为已备案的 HTTPS 域名，并在微信公众平台配置 request 合法域名。

## 本地运行

安装依赖：

```bash
npm install
```

构建小程序产物：

```bash
npm run dev:mp-weixin
```

然后打开微信开发者工具，导入目录：

```text
miniapp/dist/dev/mp-weixin
```

如果执行生产构建：

```bash
npm run build:mp-weixin
```

导入目录为：

```text
miniapp/dist/build/mp-weixin
```

本地调试时需在微信开发者工具中勾选「不校验合法域名、web-view（业务域名）、TLS 版本以及 HTTPS 证书」。

## 注意事项

- 需要先启动后端服务，默认端口 `8080`
- `manifest.json` 中的 `appid` 需填写自己的小程序 AppID
- 未配置 AI Key 时后端返回模拟数据，可完整走通全部流程
