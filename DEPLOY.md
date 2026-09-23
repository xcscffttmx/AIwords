# AI 智能单词本 - 最省钱部署上线计划

目标：用最低成本把 Web 应用部署到公网（HTTPS + 域名）。AI Key 先留空走 mock，小程序作为可选轨道。

## 成本一览（只有两项固定花费）

- 云服务器：选**轻量应用服务器**（阿里云/腾讯云轻量），比 ECS 便宜，1 核 2G 足够。优先用学生优惠机型，通常几十元/月甚至更低。
- 域名：选便宜后缀（`.cn` / `.top` / `.xyz`），首年常见几元到几十元。
- HTTPS 证书：用 Caddy 自动申请 Let's Encrypt，**免费**。
- 备案：**免费**（大陆服务器需要，走流程约几天）。
- AI Key / 小程序发布：本次都不花（Key 留空、小程序只做 H5 演示）。
- 具体服务器/域名报价请以云厂商价格计算器为准。

## 需要你先准备

- 轻量服务器一台（装好 Docker + Docker Compose），拿到公网 IP
- 一个便宜域名，DNS A 记录指向该 IP，大陆服务器完成 ICP 备案

## 改造步骤

### 1. 修配置（已完成，改仓库文件）

- [docker-compose.yml](docker-compose.yml)：backend.environment 增加 `DB_DRIVER: mysql`（否则容器里实际走 SQLite，见 [backend/main.go](backend/main.go) 第 92 行默认值）
- `DEEPSEEK_API_KEY` / `QIANWEN_API_KEY` 保持空（走 mock）
- 更新 [.env.example](.env.example) 注释，标明必须替换 JWT_SECRET 与 MySQL 密码

### 2. 加免费 HTTPS（已完成，Caddy 自动证书）

- 在 [docker-compose.yml](docker-compose.yml) 新增 `caddy` 服务，对外开 80/443，反代到 `frontend:80`
- 新增 [Caddyfile](Caddyfile)，写上你的域名，Caddy 自动申请并续期 Let's Encrypt 证书，HTTP 自动跳 HTTPS
- frontend 服务取消对外端口映射，只走 Caddy 入口

### 3. 服务器部署（最省钱路径，待你在服务器执行）

1. 服务器上 `git clone git@github.com:xcscffttmx/AIwords.git`
2. 进入 `week07/homework/`，按 `.env.example` 新建真实 `.env`：填 `DOMAIN`、`ACME_EMAIL`、强随机 `JWT_SECRET`、强 MySQL 密码；AI Key 留空
3. `docker compose up -d --build`
4. 验证：`docker compose ps` 全 healthy；浏览器访问 `https://你的域名/health` 返回 ok；注册/登录/查询/保存/删除全流程通
5. 确认数据落在 MySQL（不是 SQLite），`.env` 未提交到 Git

### 4. 小程序（可选，先只做 H5，省掉发布成本）

- [miniapp/src/config.js](miniapp/src/config.js) 把 `PROD_DOMAIN` 改成线上 HTTPS 域名
- 真机发布需微信 appid + 域名白名单 + 备案，个人主体受限，本阶段先跳过，只做 H5 演示

## 省钱要点

- 用轻量服务器而非 ECS；挑学生优惠。
- 域名选便宜后缀，注意续费价可能比首年高。
- HTTPS 用 Caddy 免费自动证书，不买付费证书。
- AI 释义想变真实时，只需在服务器 `.env` 填 Key 后重启 backend，按调用量付费，用多少花多少。
