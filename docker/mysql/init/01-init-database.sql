-- ============================================
-- MySQL 初始化脚本
-- ============================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS lottery_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE lottery_system;

-- 注意：表结构将由GORM自动创建
-- 这个脚本主要用于创建数据库和设置初始配置

-- 设置时区
SET time_zone = '+08:00';

-- 创建一个简单的健康检查表（用于数据库健康检查）
CREATE TABLE IF NOT EXISTS health_check (
    id INT AUTO_INCREMENT PRIMARY KEY,
    status VARCHAR(50) NOT NULL DEFAULT 'ok',
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入初始健康检查记录
INSERT INTO health_check (status) VALUES ('initialized');

-- 授予用户权限
-- 注意：这些变量将由Docker环境变量提供
-- GRANT ALL PRIVILEGES ON lottery_system.* TO '${MYSQL_USER}'@'%';
-- FLUSH PRIVILEGES;
