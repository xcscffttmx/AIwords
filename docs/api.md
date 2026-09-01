# API 接口文档

## 一、接口概述

- 服务名称：AI 智能单词本后端 API
- 后端框架：Gin
- 默认后端端口：`8080`
- Docker 生产访问入口：通过前端 Nginx 统一访问
- API 基础路径：`/api`
- 健康检查路径：`/health`

生产环境中，浏览器访问：

```text
http://localhost/api/...
```

开发环境中，前端访问 `/api/...`，由 Vite Proxy 转发到：

```text
http://localhost:8080/api/...
```

## 二、鉴权说明

除注册和登录接口外，其余 `/api` 接口都需要 JWT 鉴权。

请求头格式：

```http
Authorization: Bearer <JWT Token>
```

Token 由登录接口返回，前端存储在 `localStorage`，并由 Axios 请求拦截器自动携带。

## 三、通用响应与错误码

### 常见错误响应格式

```json
{
  "error": "错误原因"
}
```

### HTTP 状态码说明

| 状态码 | 含义 | 常见场景 |
|---|---|---|
| 200 | OK | 查询成功、删除成功、登录成功 |
| 201 | Created | 注册成功、保存单词成功 |
| 400 | Bad Request | 请求参数缺失、参数格式错误 |
| 401 | Unauthorized | Token 缺失、Token 无效、用户名或密码错误 |
| 404 | Not Found | 用户不存在、单词不存在 |
| 409 | Conflict | 用户名重复、单词重复保存 |
| 500 | Internal Server Error | 数据库写入失败、Token 生成失败等服务端错误 |

## 四、接口列表

| 模块 | 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|---|
| 健康检查 | GET | `/health` | 否 | 检查后端服务状态 |
| 用户 | POST | `/api/register` | 否 | 用户注册 |
| 用户 | POST | `/api/login` | 否 | 用户登录 |
| 用户 | GET | `/api/user/info` | 是 | 获取当前登录用户信息 |
| 单词 | GET | `/api/words` | 是 | 分页获取当前用户单词列表 |
| 单词 | POST | `/api/words/query` | 是 | 查询单词 |
| 单词 | POST | `/api/words/save` | 是 | 保存单词 |
| 单词 | DELETE | `/api/words/:id` | 是 | 删除单词 |

---

## 1. 健康检查

### 基本信息

- 请求方法：`GET`
- 请求路径：`/health`
- 是否鉴权：否

### 请求参数

无。

### 成功响应

状态码：`200 OK`

```json
{
  "status": "ok",
  "message": "AI Wordbook API is running"
}
```

### 失败情况

如果服务不可用，通常无法获得正常响应。Docker Compose 会使用该接口进行后端健康检查。

---

## 2. 用户注册

### 基本信息

- 请求方法：`POST`
- 请求路径：`/api/register`
- 是否鉴权：否
- Content-Type：`application/json`

### 请求体

```json
{
  "username": "student",
  "password": "123456"
}
```

### 参数说明

| 参数 | 类型 | 必填 | 约束 | 说明 |
|---|---|---|---|---|
| username | string | 是 | 长度 3-50 | 用户名，必须唯一 |
| password | string | 是 | 最少 6 位 | 用户密码，后端使用 bcrypt Hash 存储 |

### 成功响应

状态码：`201 Created`

```json
{
  "message": "User registered successfully",
  "user": 1
}
```

### 失败响应示例

用户名重复：

状态码：`409 Conflict`

```json
{
  "error": "Username already exists"
}
```

参数错误：

状态码：`400 Bad Request`

```json
{
  "error": "Key: 'Username' Error:Field validation for 'Username' failed on the 'required' tag"
}
```

---

## 3. 用户登录

### 基本信息

- 请求方法：`POST`
- 请求路径：`/api/login`
- 是否鉴权：否
- Content-Type：`application/json`

### 请求体

```json
{
  "username": "student",
  "password": "123456"
}
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

### 成功响应

状态码：`200 OK`

```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxxx",
  "user": {
    "id": 1,
    "username": "student"
  }
}
```

### 失败响应示例

用户名或密码错误：

状态码：`401 Unauthorized`

```json
{
  "error": "Invalid username or password"
}
```

---

## 4. 获取用户信息

### 基本信息

- 请求方法：`GET`
- 请求路径：`/api/user/info`
- 是否鉴权：是

### 请求头

```http
Authorization: Bearer <JWT Token>
```

### 请求参数

无。

### 成功响应

状态码：`200 OK`

```json
{
  "id": 1,
  "username": "student"
}
```

### 失败响应示例

Token 缺失：

状态码：`401 Unauthorized`

```json
{
  "error": "Missing authorization header"
}
```

Token 无效或过期：

状态码：`401 Unauthorized`

```json
{
  "error": "Invalid or expired token"
}
```

---

## 5. 查询单词

### 基本信息

- 请求方法：`POST`
- 请求路径：`/api/words/query`
- 是否鉴权：是
- Content-Type：`application/json`

### 请求头

```http
Authorization: Bearer <JWT Token>
```

### 请求体

```json
{
  "word": "innovation",
  "ai_provider": "deepseek"
}
```

### 参数说明

| 参数 | 类型 | 必填 | 约束 | 说明 |
|---|---|---|---|---|
| word | string | 是 | 非空 | 要查询的英文单词 |
| ai_provider | string | 是 | `deepseek` 或 `qianwen` | AI 提供商 |

### 业务逻辑

1. 后端读取当前登录用户 ID。
2. 根据 `user_id + word` 查询数据库中是否已有未删除记录。
3. 如果已经保存过该单词，则直接返回数据库记录，`source` 为 `database`。
4. 如果未保存，则根据 `ai_provider` 调用对应 AI 接口。
5. 如果没有配置对应 AI Key，则返回模拟数据，便于演示。
6. 查询接口只返回结果，不自动保存到数据库。

### 成功响应：来自 AI

状态码：`200 OK`

```json
{
  "message": "Query successful",
  "source": "ai",
  "word": "innovation",
  "definition": "n. 创新；革新；新方法",
  "examples": "[{\"sentence\":\"Innovation drives the growth of modern companies.\",\"translation\":\"创新推动现代公司的增长。\"},{\"sentence\":\"The school encourages innovation in teaching.\",\"translation\":\"学校鼓励教学创新。\"},{\"sentence\":\"This product is a major innovation.\",\"translation\":\"这个产品是一项重大创新。\"}]",
  "ai_provider": "deepseek"
}
```

### 成功响应：来自数据库

状态码：`200 OK`

```json
{
  "message": "Word found in your wordbook",
  "source": "database",
  "word": "innovation",
  "definition": "n. 创新；革新；新方法",
  "examples": "[{\"sentence\":\"Innovation drives the growth of modern companies.\",\"translation\":\"创新推动现代公司的增长。\"}]",
  "ai_provider": "deepseek"
}
```

### 失败响应示例

AI Provider 不合法：

状态码：`400 Bad Request`

```json
{
  "error": "Key: 'AIProvider' Error:Field validation for 'AIProvider' failed on the 'oneof' tag"
}
```

---

## 6. 保存单词

### 基本信息

- 请求方法：`POST`
- 请求路径：`/api/words/save`
- 是否鉴权：是
- Content-Type：`application/json`

### 请求头

```http
Authorization: Bearer <JWT Token>
```

### 请求体

```json
{
  "word": "innovation",
  "definition": "n. 创新；革新；新方法",
  "examples": "[{\"sentence\":\"Innovation drives the growth of modern companies.\",\"translation\":\"创新推动现代公司的增长。\"},{\"sentence\":\"The school encourages innovation in teaching.\",\"translation\":\"学校鼓励教学创新。\"},{\"sentence\":\"This product is a major innovation.\",\"translation\":\"这个产品是一项重大创新。\"}]",
  "ai_provider": "deepseek"
}
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| word | string | 是 | 单词 |
| definition | string | 是 | 单词释义 |
| examples | string | 是 | JSON 字符串格式的例句数组 |
| ai_provider | string | 是 | AI 提供商，如 `deepseek` 或 `qianwen` |

### 业务逻辑

1. 后端读取当前登录用户 ID。
2. 检查该用户是否已经保存过同一个未删除单词。
3. 如果未保存，则写入 `words` 表。
4. 如果已保存，则返回冲突错误。

### 成功响应

状态码：`201 Created`

```json
{
  "message": "Word saved successfully",
  "id": 1
}
```

### 失败响应示例

重复保存：

状态码：`409 Conflict`

```json
{
  "error": "Word already exists in your wordbook"
}
```

参数缺失：

状态码：`400 Bad Request`

```json
{
  "error": "Key: 'Definition' Error:Field validation for 'Definition' failed on the 'required' tag"
}
```

---

## 7. 获取单词列表

### 基本信息

- 请求方法：`GET`
- 请求路径：`/api/words`
- 是否鉴权：是

### 请求头

```http
Authorization: Bearer <JWT Token>
```

### Query 参数

| 参数 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|---|---|---|---|---|---|
| page | number | 否 | 1 | 最小 1 | 页码 |
| page_size | number | 否 | 10 | 1-100 | 每页数量，超出范围时后端重置为 10 |

### 请求示例

```text
GET /api/words?page=1&page_size=10
```

### 成功响应

状态码：`200 OK`

```json
{
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "word": "innovation",
        "definition": "n. 创新；革新；新方法",
        "examples": "[{\"sentence\":\"Innovation drives the growth of modern companies.\",\"translation\":\"创新推动现代公司的增长。\"}]",
        "ai_provider": "deepseek",
        "created_at": "2026-05-02T10:00:00+08:00",
        "updated_at": "2026-05-02T10:00:00+08:00"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "pages": 1
  }
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|---|---|---|
| data.list | array | 当前页单词列表 |
| data.total | number | 当前用户未删除单词总数 |
| data.page | number | 当前页码 |
| data.page_size | number | 每页数量 |
| data.pages | number | 总页数 |

---

## 8. 删除单词

### 基本信息

- 请求方法：`DELETE`
- 请求路径：`/api/words/:id`
- 是否鉴权：是

### 请求头

```http
Authorization: Bearer <JWT Token>
```

### Path 参数

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| id | number | 是 | 单词记录 ID |

### 请求示例

```text
DELETE /api/words/1
```

### 成功响应

状态码：`200 OK`

```json
{
  "message": "Word deleted successfully"
}
```

### 业务逻辑

删除时会同时校验：

- 单词 ID 是否存在；
- 单词是否属于当前登录用户；
- 单词是否未被删除。

后端使用 GORM 软删除机制，删除后会写入 `deleted_at`，普通列表查询不会再返回该记录。

### 失败响应示例

ID 格式错误：

状态码：`400 Bad Request`

```json
{
  "error": "Invalid word ID"
}
```

单词不存在或不属于当前用户：

状态码：`404 Not Found`

```json
{
  "error": "Word not found"
}
```

---

## 五、调用流程示例

### 1. 注册

```text
POST /api/register
```

### 2. 登录并保存 Token

```text
POST /api/login
```

### 3. 携带 Token 查询单词

```http
POST /api/words/query
Authorization: Bearer <JWT Token>
```

### 4. 手动保存查询结果

```http
POST /api/words/save
Authorization: Bearer <JWT Token>
```

### 5. 分页查看单词本

```http
GET /api/words?page=1&page_size=10
Authorization: Bearer <JWT Token>
```

### 6. 删除单词

```http
DELETE /api/words/1
Authorization: Bearer <JWT Token>
```
