# AI 智能单词本

AI 智能单词本是一个前后端分离的英语学习 Web 应用。用户注册登录后，可以输入英文单词并选择 AI 模型进行智能查询，系统会返回单词释义和英文例句及中文翻译。

查询结果不会自动写入数据库，只有用户点击“保存到我的单词本”后，系统才会将该单词绑定到当前用户并持久化保存。用户可以分页查看自己的单词本，也可以删除不需要的单词记录。

## 项目特点

- 前端使用 Vue3 + Vite + Element Plus
- 后端使用 Go + Gin + GORM
- 使用 JWT 完成登录鉴权
- 使用 bcrypt 加密保存用户密码
- 支持 DeepSeek 与通义千问两种 AI 提供商
- 支持 MySQL 数据持久化和 Docker Compose 全栈编排
- 通过 Vite Proxy 与 Nginx 反向代理处理跨域访问

## 主要功能

1. 用户注册与登录
2. 输入英文单词并调用 AI 生成释义和例句
3. 将查询结果保存到我的单词本
4. 查看和分页管理已保存单词
5. 删除不需要的单词记录

## 技术栈

- 前端：Vite、Vue3、Vue Router、Pinia、Axios、Element Plus
- 后端：Go、Gin、GORM、JWT、bcrypt、godotenv
- 数据库与部署：MySQL 8.0、Docker、Docker Compose、Nginx

## 本地运行

使用 Docker Compose 启动完整服务：

```bash
cp .env.example .env
docker compose up -d --build
```

启动完成后访问：

```text
http://localhost
```

如果只进行前端本地开发：

```bash
cd frontend
npm install
npm run dev
```

前端开发地址：

```text
http://localhost:5173
```

## 远程仓库

```text
git@github.com:xcscffttmx/AIwords.git
```
