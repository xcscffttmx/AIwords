# AI 智能单词本

AI 智能单词本是一个前后端分离的英语学习 Web 应用。用户登录后可以输入英文单词，选择 AI 模型生成中文释义、英文例句和中文翻译，并将查询结果保存到个人单词本中进行管理。

项目同时提供 Docker Compose 部署方式和前后端本地开发方式，适合用于完整演示注册、登录、AI 查询、单词保存、分页查看和删除记录等流程。

## 功能特性

- 用户注册、登录和 JWT 鉴权
- bcrypt 密码加密存储
- 支持 DeepSeek 与通义千问两种 AI 提供商
- 未配置 AI Key 时返回模拟数据，便于本地演示
- 查询结果手动保存到个人单词本
- 单词列表支持分页查看
- 单词删除采用软删除方式
- 前端通过 Vite Proxy 处理开发环境跨域
- 生产环境通过 Nginx 反向代理统一入口

## 技术栈

- 前端：Vue3、Vite、Vue Router、Pinia、Axios、Element Plus
- 后端：Go、Gin、GORM、JWT、bcrypt、godotenv
- 数据库：MySQL 8.0
- 部署：Docker、Docker Compose、Nginx

## 项目结构

```text
.
├── backend/              # Go + Gin 后端服务
├── frontend/             # Vue3 + Vite 前端应用
├── docs/                 # 接口文档和数据库说明
├── docker-compose.yml    # Docker Compose 编排配置
├── .env.example          # 环境变量示例
└── README.md             # 项目说明文档
```

## 环境要求

- Go 1.21+
- Node.js 18+
- npm
- MySQL 8.0
- Docker Desktop 与 Docker Compose，可选，用于容器化启动

## 快速启动

推荐使用 Docker Compose 一次性启动前端、后端和数据库：

```bash
cp .env.example .env
docker compose up -d --build
```

启动完成后访问：

```text
http://localhost
```

健康检查地址：

```text
http://localhost/health
```

如果本机 80 端口被占用，可以在 `.env` 中修改：

```env
FRONTEND_PORT=8088
```

然后访问：

```text
http://localhost:8088
```

## 本地开发

后端本地启动：

```bash
cd backend
go mod download
go run .
```

后端默认运行在：

```text
http://localhost:8080
```

前端本地启动：

```bash
cd frontend
npm install
npm run dev
```

前端默认运行在：

```text
http://localhost:5173
```

开发环境下，前端请求 `/api` 会通过 Vite Proxy 转发到后端 `http://localhost:8080`。

## 环境变量

项目根目录提供 `.env.example`，复制为 `.env` 后可按需修改：

```env
MYSQL_ROOT_PASSWORD=root123
DB_NAME=ai_wordbook
MYSQL_USER=wordbook
MYSQL_PASSWORD=wordbook123
JWT_SECRET=please-change-this-jwt-secret
DEEPSEEK_API_KEY=
QIANWEN_API_KEY=
FRONTEND_PORT=80
```

`DEEPSEEK_API_KEY` 和 `QIANWEN_API_KEY` 可以留空，后端会使用模拟数据返回查询结果。

## 使用流程

1. 打开前端页面
2. 注册账号并登录
3. 输入英文单词
4. 选择 DeepSeek 或通义千问
5. 点击查询并查看释义和例句
6. 点击“保存到我的单词本”
7. 在单词列表中分页查看或删除记录

## 截图说明

建议将项目截图放在 `docs/images/` 目录下，并在 README 中引用：

```markdown
![登录页](docs/images/login.png)
![单词查询页](docs/images/dashboard.png)
```

推荐截图内容：

- 登录页或注册页
- 单词查询结果
- 我的单词本列表
- Docker Compose 启动成功后的服务页面

## 文档

- API 文档：`docs/api.md`
- 数据库设计：`docs/db.md`
- 数据库初始化脚本：`docs/init.sql`

## 远程仓库

```text
git@github.com:xcscffttmx/AIwords.git
```
