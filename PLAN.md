# AI 智能单词本 - 开发计划与完成记录

## 📋 项目概述

**项目名称**：AI 智能单词本  
**技术栈**：Vite + Vue3（前端） | Gin + GORM（后端） | MySQL 8.0（数据库） | Docker Compose（部署）  
**架构特点**：前后端分离、Nginx 反向代理、JWT 认证、AI 大模型集成  
**当前状态**：✅ 全部阶段已完成，并已通过 Docker Compose 联调验证

---

## ✅ 当前状态分析

### 已完成部分

1. ✅ **后端基础架构**：`backend/main.go` 已实现完整 API 逻辑
   - 用户注册/登录（JWT 认证）
   - 密码 bcrypt Hash 存储
   - JWT 鉴权中间件
   - 当前用户信息接口
   - 单词查询（AI 集成 + 数据库优先读取）
   - 单词手动保存
   - 单词分页列表
   - 单词软删除
   - 健康检查接口 `/health`

2. ✅ **数据库设计**：`docs/init.sql` 和 `docs/db.md` 已完成
   - `users` 用户表
   - `words` 单词本表
   - 外键关联 `words.user_id -> users.id`
   - `deleted_at` 软删除字段
   - 用户名单列唯一索引
   - 用户单词联合索引
   - 数据库通过 MySQL 容器初始化脚本创建
   - 后端未使用 `AutoMigrate`

3. ✅ **前端项目**：`frontend/` 已完成
   - Vite + Vue3 项目结构
   - Vue Router 路由
   - Pinia 用户状态管理
   - Axios 请求封装
   - Token 自动注入
   - 登录页
   - 注册页
   - Dashboard 主页面
   - 单词查询、保存、分页列表、删除
   - Element Plus UI
   - 响应式样式

4. ✅ **Nginx 配置**：`frontend/nginx.conf` 已完成
   - `/` 服务前端静态资源
   - `/api/` 反向代理到 `backend:8080/api/`
   - `/health` 反向代理到 `backend:8080/health`
   - 生产环境通过 Nginx 实现同源访问

5. ✅ **Docker 配置**：Docker Compose 和 Dockerfile 已完成
   - `backend/Dockerfile` 多阶段构建 Go 服务
   - `frontend/Dockerfile` 多阶段构建 Vite + Nginx 镜像
   - `frontend/.dockerignore` 已补充
   - `docker-compose.yml` 定义 `db`、`backend`、`frontend` 三个服务
   - `db` 和 `backend` 不直接映射宿主机端口
   - 仅 `frontend` 暴露宿主机端口
   - 健康检查已配置
   - 自定义 bridge 网络已配置

6. ✅ **项目文档**：文档已补齐并更新
   - `README.md`：项目说明、架构、运行指南、验收说明
   - `docs/api.md`：完整 API 接口文档
   - `docs/db.md`：数据库设计文档
   - `docs/init.sql`：数据库初始化脚本

7. ✅ **测试与验证**：已完成 Docker 联调测试
   - `docker compose build` 构建成功
   - `docker compose up -d` 启动成功
   - 三个容器均为 healthy
   - Nginx 访问前端成功
   - `/health` 反向代理成功
   - 注册、登录、JWT 获取用户信息测试成功
   - 单词查询、保存、分页列表、删除测试成功

### 待完成部分

当前按作业要求和本计划拆分的任务均已完成。  
仅剩提交前人工确认项：

1. ⚠️ `README.md` 中的学校和学号需要按真实信息填写。
2. ⚠️ 如需真实 AI 调用，请配置 `DEEPSEEK_API_KEY` 或 `QIANWEN_API_KEY`；未配置时系统使用模拟数据，功能流程仍可演示。
3. ⚠️ 若 80 端口被占用，可在根目录 `.env` 中设置 `FRONTEND_PORT=8088` 后重新启动。

---

## 🎯 开发任务清单

### 阶段一：后端完善（状态：✅ 已完成）

#### 1.1 补充后端依赖文件

- [x] 检查 `backend/go.sum` 文件
- [x] 验证所有依赖包版本兼容性
- [x] 测试后端编译

完成记录：

- `backend/go.sum` 已存在。
- 已执行后端编译验证：`go build ./...`。
- 编译通过。

#### 1.2 优化后端代码

- [x] 检查 JWT 密钥配置（从环境变量读取）
- [x] 验证 AI 接口调用逻辑
- [x] 测试数据库连接配置
- [x] 检查无 CORS 中间件
- [x] 检查无 GORM `AutoMigrate`
- [x] 修复 GORM 字段映射问题：`AIProvider` 显式映射到 `ai_provider`

完成记录：

- JWT Secret 从 `JWT_SECRET` 环境变量读取。
- 数据库配置从 `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME` 读取。
- DeepSeek / 通义千问 API Key 从环境变量读取。
- 未配置 AI Key 时自动返回模拟数据，便于演示完整流程。
- 后端未配置 CORS。
- 后端未使用 `AutoMigrate`。
- 已修复 `AIProvider` 字段默认映射为 `ai_provider`，与 `docs/init.sql` 表结构保持一致。

---

### 阶段二：前端开发（状态：✅ 已完成）

#### 2.1 初始化 Vite + Vue3 项目

- [x] 创建 `frontend/package.json`
- [x] 配置 Vite + Vue3
- [x] 引入 Axios、Vue Router、Pinia、Element Plus

#### 2.2 创建前端项目结构

```text
frontend/
├── src/
│   ├── api/
│   │   ├── request.js
│   │   └── index.js
│   ├── components/
│   │   └── WordCard.vue
│   ├── views/
│   │   ├── Login.vue
│   │   ├── Register.vue
│   │   └── Dashboard.vue
│   ├── router/
│   │   └── index.js
│   ├── stores/
│   │   └── user.js
│   ├── App.vue
│   ├── main.js
│   └── styles.css
├── index.html
├── vite.config.js
├── nginx.conf
├── Dockerfile
└── .dockerignore
```

- [x] API 请求封装
- [x] 路由配置
- [x] 用户状态管理
- [x] 页面组件
- [x] 单词卡片组件
- [x] 全局样式

#### 2.3 实现核心功能模块

**2.3.1 API 请求层**

- [x] 配置 Axios 实例（BaseURL、超时、请求/响应拦截器）
- [x] 实现 Token 自动注入（Authorization: Bearer `<Token>`）
- [x] 实现统一错误处理

**2.3.2 用户认证模块**

- [x] 登录页面（用户名 + 密码表单）
- [x] 注册页面（用户名 + 密码确认）
- [x] Token 存储（localStorage）
- [x] 路由守卫（未登录跳转登录页）
- [x] 退出登录

**2.3.3 单词学习模块**

- [x] 单词查询页面
  - [x] 输入框（单词）
  - [x] 下拉选择（AI 提供商：DeepSeek / 通义千问）
  - [x] 查询按钮
  - [x] 结果显示区域（释义 + 例句）
  - [x] 保存按钮
- [x] 单词列表页面
  - [x] 表格展示（单词、释义、AI 来源、操作）
  - [x] 分页器（page / page_size）
  - [x] 删除按钮（二次确认）

**2.3.4 UI 样式**

- [x] 使用 Element Plus 组件库
- [x] 响应式布局
- [x] 加载状态（Loading）
- [x] 消息提示（Message）
- [x] 现代化卡片布局与渐变背景

#### 2.4 配置文件

**2.4.1 `vite.config.js`（开发环境）**

- [x] 配置 `/api` 代理到 `http://localhost:8080`

**2.4.2 `nginx.conf`（生产环境）**

- [x] 配置前端静态资源
- [x] 配置 `/api/` 反向代理到 `backend:8080/api/`
- [x] 配置 `/health` 反向代理到 `backend:8080/health`

**2.4.3 `frontend/Dockerfile`**

- [x] Node 构建阶段
- [x] Nginx 运行阶段
- [x] 复制 Vite build 产物
- [x] 替换自定义 Nginx 配置

---

### 阶段三：Docker 部署优化（状态：✅ 已完成）

#### 3.1 优化 `docker-compose.yml`

- [x] 确认服务依赖关系（`depends_on` + `condition: service_healthy`）
- [x] 配置健康检查（healthcheck）
- [x] 设置环境变量（AI API Keys、JWT、数据库配置、前端端口）
- [x] 配置网络（自定义 bridge 网络）
- [x] `db` 不映射宿主机端口
- [x] `backend` 不映射宿主机端口，仅 `expose: 8080`
- [x] 仅 `frontend` 对外暴露端口
- [x] `docs/init.sql` 只读挂载到 MySQL 初始化目录

#### 3.2 测试 Docker 构建

- [x] 执行 `docker compose config` 验证配置语法
- [x] 执行 `docker compose build` 构建所有服务
- [x] 执行 `docker compose up -d` 启动服务
- [x] 执行 `docker compose ps` 查看健康状态

完成记录：

- 已成功构建 `docker-gin-backend:latest`。
- 已成功构建 `docker-gin-frontend:latest`。
- 当前服务状态：
  - `ai-wordbook-db`：healthy
  - `ai-wordbook-backend`：healthy
  - `ai-wordbook-frontend`：healthy

---

### 阶段四：文档编写（状态：✅ 已完成）

#### 4.1 `README.md`（项目说明）

- [x] 项目基本信息（姓名、学校、学号占位）
- [x] 项目简介与架构图 / 架构说明
- [x] 技术栈说明
- [x] 运行指南
  - [x] 前置依赖（Docker、Docker Compose、Node.js、Go）
  - [x] 环境变量配置（`.env` 文件创建）
  - [x] AI API Key 配置说明
  - [x] 一键启动命令
  - [x] 访问地址说明
  - [x] 日志查看与停止服务说明
- [x] 本地开发方式
- [x] 主要功能使用流程
- [x] 验收说明

#### 4.2 `docs/api.md`（API 接口文档）

- [x] 接口列表
  - [x] `GET /health` - 健康检查
  - [x] `POST /api/register` - 用户注册
  - [x] `POST /api/login` - 用户登录
  - [x] `GET /api/user/info` - 获取用户信息
  - [x] `GET /api/words` - 获取单词列表（分页）
  - [x] `POST /api/words/query` - 查询单词
  - [x] `POST /api/words/save` - 保存单词
  - [x] `DELETE /api/words/:id` - 删除单词
- [x] 每个接口包含：
  - [x] 请求方法 + 路径
  - [x] 鉴权说明
  - [x] 请求参数（Query / Body / Path）
  - [x] 成功返回示例
  - [x] 失败错误码及含义
  - [x] 业务逻辑说明

#### 4.3 更新 `docs/db.md`

- [x] 确认数据库设计与 `init.sql` 一致
- [x] 补充索引说明
- [x] 补充查询优化建议
- [x] 补充表关系说明
- [x] 补充软删除机制说明
- [x] 补充 GORM 模型字段映射说明

---

### 阶段五：测试与验证（状态：✅ 已完成）

#### 5.1 功能测试

- [x] 用户注册/登录
- [x] JWT Token 验证
- [x] 获取当前用户信息
- [x] 单词查询（未配置 AI Key 时使用模拟数据）
- [x] 单词保存
- [x] 单词列表分页
- [x] 单词删除（软删除）

验证结果：

- 注册成功：`User registered successfully`
- 登录成功：`Login successful`
- JWT 获取用户信息成功
- 查询单词成功：`source = ai`
- 保存单词成功，返回保存 ID
- 分页列表保存后总数为 1
- 删除单词成功：`Word deleted successfully`
- 删除后列表总数为 0

#### 5.2 跨域验证

- [x] 开发环境：Vite Proxy 配置已完成
- [x] 生产环境：Nginx 反向代理测试通过
- [x] 后端未配置 CORS 中间件

验证结果：

- `http://localhost/` 前端首页返回 200
- `http://localhost/health` 经 Nginx 代理后返回：

```json
{"message":"AI Wordbook API is running","status":"ok"}
```

#### 5.3 Docker 部署测试

- [x] 从零构建：`docker compose build`
- [x] 一键启动：`docker compose up -d`
- [x] 服务健康检查：`docker compose ps`
- [x] 数据库初始化验证
- [x] 业务接口通过 Nginx 入口访问验证

当前容器状态：

```text
ai-wordbook-db         healthy
ai-wordbook-backend    healthy
ai-wordbook-frontend   healthy
```

---

## 📅 时间规划与实际完成情况

| 阶段 | 任务 | 原预计时间 | 当前状态 |
|------|------|------------|----------|
| 阶段一 | 后端完善 | 1 小时 | ✅ 已完成 |
| 阶段二 | 前端开发 | 3-4 小时 | ✅ 已完成 |
| 阶段三 | Docker 部署 | 1 小时 | ✅ 已完成 |
| 阶段四 | 文档编写 | 1-2 小时 | ✅ 已完成 |
| 阶段五 | 测试验证 | 1 小时 | ✅ 已完成 |
| **总计** | | **7-9 小时** | ✅ 全部完成 |

---

## 🎯 核心考察点完成情况

1. ✅ **前后端分离架构**：`backend/` 与 `frontend/` 分离，目录结构清晰。
2. ✅ **跨域处理**：
   - 开发环境：Vite Proxy。
   - 生产环境：Nginx 反向代理。
   - 后端无 CORS。
3. ✅ **身份认证**：JWT Token + 前端路由守卫 + Axios 自动携带 Token。
4. ✅ **数据库设计**：MySQL 8.0 + GORM + 用户表 + 单词本表 + 软删除。
5. ✅ **Docker 部署**：多阶段构建 + Compose 服务编排 + 健康检查。
6. ✅ **AI 接口集成**：支持 DeepSeek / 通义千问；未配置 Key 时使用模拟数据保证演示稳定。
7. ✅ **数据库初始化规范**：使用 `docs/init.sql`，未使用 `AutoMigrate`。
8. ✅ **接口文档与数据库文档**：已补齐。

---

## ⚠️ 注意事项

1. **严禁在后端配置 CORS**：当前已满足，跨域问题通过代理解决。
2. **密码必须 Hash 加密**：当前使用 bcrypt。
3. **数据库表必须通过 `init.sql` 初始化**：当前已满足，未使用 GORM `AutoMigrate`。
4. **Token 必须前端持久化**：当前使用 localStorage + 请求头携带。
5. **单词本必须支持分页**：当前后端接收 `page` 和 `page_size` 参数。
6. **软删除机制**：当前使用 `deleted_at` 字段标记。
7. **学校和学号**：`README.md` 中仍为占位文本，提交前请手动替换。
8. **AI Key**：如果需要真实 AI 响应，请配置环境变量；否则使用模拟数据。

---

## 🚀 当前运行方式

进入项目目录：

```bash
cd week07/homework/docker-gin
```

构建并启动：

```bash
docker compose up -d --build
```

查看状态：

```bash
docker compose ps
```

访问页面：

```text
http://localhost
```

健康检查：

```text
http://localhost/health
```

停止服务：

```bash
docker compose down
```

如需清空数据库数据：

```bash
docker compose down -v
```

---

## 📝 下一步行动

✅ 开发、文档、部署与测试任务均已完成。  
提交前建议执行以下人工检查：

1. 打开 `http://localhost`，确认页面可正常访问。
2. 注册一个测试账号，例如：用户名 `test001`，密码 `123456`。
3. 登录后查询单词并保存。
4. 检查分页列表和删除功能。
5. 将 `README.md` 中的学校和学号替换为真实信息。
6. 如老师要求真实 AI 调用，请填写真实 `DEEPSEEK_API_KEY` 或 `QIANWEN_API_KEY`。

---

**制定时间**：2026-05-01  
**更新时间**：2026-05-03  
**制定人**：AI 助手  
**当前状态**：✅ 已完成
