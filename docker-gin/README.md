# AI 智能单词本

## 一、项目基本信息

| 信息 | 内容 |
| --- | --- |
| **项目名称** | AI 智能单词本 |
| **开发者** | 郑康月 |
| **学校** | 中国地质大学（武汉） |
| **学号** | 1202411197 |


## 二、项目简介

AI 智能单词本是一个前后端分离的英语学习 Web 应用。用户注册登录后，可以输入英文单词并选择 AI 模型（DeepSeek 或通义千问）进行智能查询，后端会返回单词释义和 3 条英文例句及中文翻译。

查询结果不会自动写入数据库，只有用户点击“保存到我的单词本”后，系统才会将该单词绑定到当前用户并持久化保存。用户可以分页查看自己的单词本，也可以删除不需要的单词记录。删除操作采用软删除方式，保证数据可追溯。

本项目重点实现：

- 前后端分离架构
- JWT 登录鉴权
- bcrypt 密码加密
- AI 大模型接口调用
- MySQL 8.0 数据持久化
- GORM 数据访问
- Docker Compose 全栈编排
- Vite Proxy 与 Nginx 反向代理跨域治理

## 三、技术栈

### 前端

- Vite
- Vue3
- Vue Router
- Pinia
- Axios
- Element Plus
- Nginx

### 后端

- Go
- Gin
- GORM
- MySQL Driver
- JWT
- bcrypt
- godotenv

### 数据库与部署

- MySQL 8.0
- Docker
- Docker Compose
- Nginx 反向代理

## 四、系统架构说明

```text
浏览器
  |
  | http://localhost
  v
frontend 容器（Nginx，唯一对外入口）
  |
  |-- /                 -> Vue 静态资源
  |-- /api/*            -> proxy_pass http://backend:8080/api/*
  |-- /health           -> proxy_pass http://backend:8080/health
                         |
                         v
backend 容器（Gin API）
  |
  | GORM + MySQL Driver
  v
db 容器（MySQL 8.0）
```

### 跨域处理说明

本项目严格遵守作业要求：**后端 Go 代码中没有配置 CORS 中间件**。

- 开发环境：前端通过 `frontend/vite.config.js` 中的 Vite Proxy 将 `/api` 转发到 `http://localhost:8080`。
- 生产环境：前端容器使用 Nginx，通过 `frontend/nginx.conf` 将 `/api/` 反向代理到内部 `backend:8080`。
- Docker 部署时，浏览器只访问 Nginx 暴露的端口，前端页面与 API 对浏览器表现为同源。

## 五、已完成任务清单

- 用户注册接口
- 用户登录接口
- bcrypt 密码 Hash 存储
- JWT Token 签发与校验
- 前端 Token 持久化到 `localStorage`
- Axios 请求拦截器自动携带 `Authorization: Bearer <Token>`
- 路由守卫：未登录用户自动跳转登录页
- AI 单词查询：支持 `deepseek` 和 `qianwen`
- 已保存单词优先从数据库读取
- 未保存单词调用 AI 接口，未配置 Key 时返回模拟数据
- 手动保存单词
- 分页获取当前用户单词列表
- 删除单词，使用软删除机制
- MySQL 初始化脚本 `docs/init.sql`
- 禁止使用 GORM `AutoMigrate`
- Vite 开发代理配置
- Nginx 生产反向代理配置
- 后端 Dockerfile 多阶段构建
- 前端 Dockerfile 多阶段构建
- Docker Compose 编排 `db`、`backend`、`frontend` 三个服务
- API 文档与数据库设计文档

## 六、目录结构

```text
week07/homework/docker-gin
├── backend/
│   ├── Dockerfile
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── vite.config.js
│   ├── package.json
│   ├── index.html
│   └── src/
│       ├── api/
│       ├── components/
│       ├── router/
│       ├── stores/
│       ├── views/
│       ├── App.vue
│       ├── main.js
│       └── styles.css
├── docs/
│   ├── api.md
│   ├── db.md
│   └── init.sql
├── .env.example
├── docker-compose.yml
├── PLAN.md
├── profile.md
└── README.md
```

## 七、从零启动指南

### 1. 前置依赖

请先安装并启动：

- Docker Desktop
- Docker Compose

如果需要本地开发前端，还需要安装：

- Node.js 18+
- npm

如果需要本地开发后端，还需要安装：

- Go 1.21+
- MySQL 8.0 或使用 Docker MySQL

### 2. 配置环境变量

Docker Compose 会读取 `week07/homework/docker-gin/.env`。请先复制根目录示例配置：

```bash
cp .env.example .env
```

Windows PowerShell 可以执行：

```powershell
Copy-Item .env.example .env
```

然后编辑 `.env`，填写真实配置：

```bash
JWT_SECRET=请替换为复杂随机字符串
DEEPSEEK_API_KEY=你的 DeepSeek API Key
QIANWEN_API_KEY=你的通义千问 API Key
FRONTEND_PORT=80
MYSQL_ROOT_PASSWORD=root123
DB_NAME=ai_wordbook
MYSQL_USER=wordbook
MYSQL_PASSWORD=wordbook123
```

注意：`docker-gin/.env` 用于 Docker Compose 部署，包含真实密钥，不要提交到 Git；`backend/.env.example` 仅用于后端本地运行配置参考。

如果没有 AI API Key，可以留空：

```bash
DEEPSEEK_API_KEY=
QIANWEN_API_KEY=
```

此时后端会返回模拟单词数据，仍可完整演示注册、登录、查询、保存、分页列表和删除流程。

后端本地运行时也可以使用 `backend/.env` 或 `backend/.env.example` 中的配置。

### 3. 一键构建并启动

进入项目目录：

```bash
cd week07/homework/docker-gin
```

执行：

```bash
docker compose up -d --build
```

首次启动时：

1. MySQL 容器启动。
2. Docker Compose 将 `docs/init.sql` 挂载到 MySQL 的 `/docker-entrypoint-initdb.d/init.sql`。
3. MySQL 自动创建 `ai_wordbook` 数据库、`users` 表和 `words` 表。
4. 后端等待数据库健康检查通过后启动。
5. 前端 Nginx 等待后端健康检查通过后启动。

### 4. 访问系统

- 前端页面：`http://localhost`
- 后端健康检查：`http://localhost/health`
- API 入口示例：`http://localhost/api/login`

Docker 部署模式下，只有 `frontend` 服务映射宿主机端口。`backend` 和 `db` 不直接暴露到宿主机，符合安全要求。

如果 80 端口被占用，可以在项目根目录 `.env` 中修改：

```bash
FRONTEND_PORT=8088
```

然后访问：

```text
http://localhost:8088
```

### 5. 查看服务状态与日志

查看容器状态：

```bash
docker compose ps
```

查看全部日志：

```bash
docker compose logs -f
```

只查看后端日志：

```bash
docker compose logs -f backend
```

### 6. 停止服务

停止容器但保留数据库数据卷：

```bash
docker compose down
```

停止容器并清空数据库数据卷：

```bash
docker compose down -v
```

如果修改了 `docs/init.sql` 并希望重新初始化数据库，需要执行 `docker compose down -v` 后再启动。

## 八、本地开发方式

### 1. 后端本地开发

确保本地 MySQL 已启动，并创建了 `ai_wordbook` 数据库和对应表，或者通过 Docker 单独启动数据库。

进入后端目录：

```bash
cd week07/homework/docker-gin/backend
```

安装依赖并运行：

```bash
go mod download
go run .
```

默认后端地址：

```text
http://localhost:8080
```

### 2. 前端本地开发

进入前端目录：

```bash
cd week07/homework/docker-gin/frontend
```

安装依赖并启动：

```bash
npm install
npm run dev
```

默认前端地址：

```text
http://localhost:5173
```

开发环境下，前端请求 `/api` 会通过 `vite.config.js` 代理到：

```text
http://localhost:8080
```

## 九、主要功能使用流程

1. 打开前端页面。
2. 注册一个新用户。
3. 登录获取 Token。
4. 在主页面输入英文单词。
5. 选择 DeepSeek 或通义千问。
6. 点击查询，查看 AI 返回的释义和例句。
7. 点击“保存到我的单词本”。
8. 在右侧列表查看已保存单词。
9. 使用分页器翻页。
10. 点击删除按钮软删除单词。

## 十、接口与数据库文档

- API 接口文档：`docs/api.md`
- 数据库设计文档：`docs/db.md`
- 数据库初始化脚本：`docs/init.sql`

## 十一、验收说明

本项目满足以下作业关键要求：

- 前端基于 Vite 构建。
- 后端使用 Gin + GORM。
- 数据库使用 MySQL 8.0。
- 用户认证使用 JWT。
- 密码使用 bcrypt 哈希保存。
- 数据库表通过 `docs/init.sql` 初始化。
- 后端没有使用 `AutoMigrate` 自动建表。
- 后端没有配置 CORS。
- 开发环境通过 Vite Proxy 处理跨域。
- 生产环境通过 Nginx 反向代理处理跨域。
- Docker Compose 定义 `db`、`backend`、`frontend` 三个服务。
- 对外只暴露前端 Nginx 服务。
