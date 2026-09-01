# 数据库设计文档

## 一、数据库概述

- 数据库名称：`ai_wordbook`
- 数据库版本：MySQL 8.0
- 默认字符集：`utf8mb4`
- 默认排序规则：`utf8mb4_unicode_ci`
- 存储引擎：InnoDB
- 初始化方式：Docker Compose 挂载并执行 `docs/init.sql`

本项目严格遵守作业要求：数据库表通过 `docs/init.sql` 初始化，后端 Go 代码中不使用 GORM `AutoMigrate` 自动建表。

## 二、初始化脚本说明

Docker Compose 中将初始化脚本挂载到 MySQL 容器：

```text
./docs/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
```

MySQL 容器首次创建数据目录时会自动执行该脚本，完成以下操作：

1. 创建数据库 `ai_wordbook`。
2. 创建用户表 `users`。
3. 创建单词本表 `words`。
4. 创建主键、唯一索引、普通索引、联合索引和外键约束。

如果需要重新执行初始化脚本，需要清空数据库卷：

```bash
docker compose down -v
docker compose up -d --build
```

## 三、表关系设计

```text
users 1 -------- N words
```

- 一个用户可以保存多个单词。
- 一个单词记录只属于一个用户。
- `words.user_id` 外键关联 `users.id`。
- 当用户被删除时，其单词记录通过 `ON DELETE CASCADE` 级联删除。
- 业务删除单词时采用软删除，通过 `words.deleted_at` 标记。

## 四、用户表：users

### 1. 表用途

`users` 表用于存储系统用户账号信息和登录认证凭据。

### 2. 字段设计

| 字段名 | 数据类型 | 是否为空 | 主键/外键 | 索引 | 默认值 | 业务含义 |
|---|---|---|---|---|---|---|
| id | BIGINT UNSIGNED | 否 | 主键 | PRIMARY KEY | AUTO_INCREMENT | 用户 ID |
| username | VARCHAR(50) | 否 | 否 | UNIQUE KEY `uk_username` | 无 | 用户名，注册时唯一 |
| password_hash | VARCHAR(255) | 否 | 否 | 无 | 无 | bcrypt 加密后的密码哈希 |
| created_at | TIMESTAMP | 否 | 否 | 无 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 否 | 否 | 无 | CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

### 3. 建表语句

```sql
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名（唯一）',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希（bcrypt 加密）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

### 4. 索引说明

| 索引名 | 字段 | 类型 | 作用 |
|---|---|---|---|
| PRIMARY | id | 主键索引 | 根据用户 ID 查询用户 |
| uk_username | username | 唯一索引 | 保证用户名唯一，加速登录查询 |

### 5. 安全说明

- 密码严禁明文存储。
- 后端使用 `bcrypt.GenerateFromPassword` 生成密码哈希。
- 登录时使用 `bcrypt.CompareHashAndPassword` 校验密码。
- `password_hash` 字段在接口 JSON 响应中不会返回给前端。

## 五、单词本表：words

### 1. 表用途

`words` 表用于存储用户手动保存的单词查询结果，包括单词、释义、例句、AI 来源和软删除状态。

### 2. 字段设计

| 字段名 | 数据类型 | 是否为空 | 主键/外键 | 索引 | 默认值 | 业务含义 |
|---|---|---|---|---|---|---|
| id | BIGINT UNSIGNED | 否 | 主键 | PRIMARY KEY | AUTO_INCREMENT | 单词记录 ID |
| user_id | BIGINT UNSIGNED | 否 | 外键，关联 `users.id` | `idx_user_id`、`idx_user_word` | 无 | 当前单词所属用户 ID |
| word | VARCHAR(100) | 否 | 否 | `idx_word`、`idx_user_word` | 无 | 英文单词 |
| definition | TEXT | 否 | 否 | 无 | 无 | 单词释义文本 |
| examples | JSON | 否 | 否 | 无 | 无 | 例句数组，包含英文句子和中文翻译 |
| ai_provider | VARCHAR(50) | 否 | 否 | 无 | 无 | AI 提供商，值为 `deepseek` 或 `qianwen` |
| created_at | TIMESTAMP | 否 | 否 | 无 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 否 | 否 | 无 | CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 是 | 否 | 无 | NULL | 删除时间，NULL 表示未删除 |

### 3. 建表语句

```sql
CREATE TABLE `words` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '单词记录 ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID（外键）',
    `word` VARCHAR(100) NOT NULL COMMENT '单词',
    `definition` TEXT NOT NULL COMMENT '单词释义（JSON 格式）',
    `examples` JSON NOT NULL COMMENT '例句列表（JSON 数组，包含 3 条例句）',
    `ai_provider` VARCHAR(50) NOT NULL COMMENT 'AI 提供商（deepseek/qianwen）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间（软删除）',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_word` (`word`),
    KEY `idx_user_word` (`user_id`, `word`),
    CONSTRAINT `fk_words_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户单词本表';
```

### 4. 索引说明

| 索引名 | 字段 | 类型 | 作用 |
|---|---|---|---|
| PRIMARY | id | 主键索引 | 根据单词记录 ID 查询、删除 |
| idx_user_id | user_id | 普通索引 | 加速查询当前用户的所有单词 |
| idx_word | word | 普通索引 | 加速根据单词进行查询 |
| idx_user_word | user_id, word | 联合索引 | 加速判断当前用户是否已经保存某个单词 |

### 5. 外键说明

| 外键名 | 子表字段 | 父表字段 | 删除策略 | 说明 |
|---|---|---|---|---|
| fk_words_user | words.user_id | users.id | ON DELETE CASCADE | 用户删除后，其单词记录级联删除 |

### 6. 软删除说明

`words.deleted_at` 用于软删除：

- `deleted_at IS NULL`：正常记录。
- `deleted_at IS NOT NULL`：已删除记录。

后端删除单词时调用 GORM 的 `Delete`，由于模型包含 `gorm.DeletedAt` 字段，GORM 会执行软删除，将 `deleted_at` 设置为当前时间。

业务查询中会限制：

```sql
WHERE user_id = ? AND deleted_at IS NULL
```

因此已删除记录不会出现在普通单词列表中。

## 六、JSON 字段结构

### 1. examples 字段

`examples` 字段为 MySQL JSON 类型。后端当前以 JSON 字符串的方式接收和返回，内容结构如下：

```json
[
  {
    "sentence": "Innovation drives the growth of modern companies.",
    "translation": "创新推动现代公司的增长。"
  },
  {
    "sentence": "The school encourages innovation in teaching.",
    "translation": "学校鼓励教学创新。"
  },
  {
    "sentence": "This product is a major innovation.",
    "translation": "这个产品是一项重大创新。"
  }
]
```

### 2. definition 字段

`definition` 字段在初始化脚本中为 `TEXT` 类型，用于保存 AI 返回的释义文本，例如：

```text
n. 创新；革新；新方法
```

说明：`init.sql` 中该字段注释为“单词释义（JSON 格式）”，但当前后端和前端实际保存的是释义文本。字段类型为 `TEXT`，可以兼容纯文本或未来扩展为 JSON 字符串。

## 七、核心查询场景与优化

### 1. 用户注册查重

```sql
SELECT * FROM users WHERE username = ? LIMIT 1;
```

使用索引：`uk_username`

### 2. 用户登录

```sql
SELECT * FROM users WHERE username = ? LIMIT 1;
```

使用索引：`uk_username`

### 3. 查询当前用户是否已保存单词

```sql
SELECT * FROM words
WHERE user_id = ? AND word = ? AND deleted_at IS NULL
LIMIT 1;
```

使用索引：`idx_user_word`

### 4. 分页获取当前用户单词本

```sql
SELECT * FROM words
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
```

使用索引：`idx_user_id`

### 5. 统计当前用户单词总数

```sql
SELECT COUNT(*) FROM words
WHERE user_id = ? AND deleted_at IS NULL;
```

使用索引：`idx_user_id`

### 6. 删除当前用户的指定单词

```sql
UPDATE words
SET deleted_at = NOW()
WHERE id = ? AND user_id = ? AND deleted_at IS NULL;
```

使用索引：主键 `PRIMARY`，同时校验 `user_id` 保证用户只能删除自己的单词。

## 八、与 GORM 模型的对应关系

### User 模型

| Go 字段 | 数据库字段 | 说明 |
|---|---|---|
| ID | id | 用户 ID |
| Username | username | 用户名 |
| PasswordHash | password_hash | 密码哈希 |
| CreatedAt | created_at | 创建时间 |
| UpdatedAt | updated_at | 更新时间 |

### Word 模型

| Go 字段 | 数据库字段 | 说明 |
|---|---|---|
| ID | id | 单词记录 ID |
| UserID | user_id | 用户 ID |
| Word | word | 单词 |
| Definition | definition | 释义 |
| Examples | examples | 例句 JSON 字符串 |
| AIProvider | ai_provider | AI 提供商 |
| CreatedAt | created_at | 创建时间 |
| UpdatedAt | updated_at | 更新时间 |
| DeletedAt | deleted_at | GORM 软删除字段 |

## 九、扩展性设计

后续可扩展字段：

### users 表

- `email`：邮箱，用于找回密码。
- `avatar`：头像 URL。
- `last_login_at`：最后登录时间。

### words 表

- `phonetic`：音标。
- `difficulty`：单词难度。
- `mastery_level`：掌握程度。
- `last_reviewed_at`：最后复习时间。
- `tags`：自定义标签。

## 十、版本历史

- v1.0，2026-05-01：完成用户表和单词本表设计。
- v1.1，2026-05-02：补充 Docker 初始化、索引、外键、软删除与 GORM 模型对应说明。
