-- AI 智能单词本 - 数据库初始化脚本
-- 执行顺序：docker-compose 启动时自动执行 /docker-entrypoint-initdb.d/ 目录下的所有 .sql 文件

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `ai_wordbook` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `ai_wordbook`;

-- 用户表
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名（唯一）',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希（bcrypt 加密）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 单词本表
DROP TABLE IF EXISTS `words`;
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

-- 插入测试数据（可选）
-- INSERT INTO `users` (`username`, `password_hash`) VALUES 
-- ('test', '$2a$10$test_hash_placeholder');
