-- MySQL 数据库初始化脚本
-- 抽奖系统数据库结构

-- 创建数据库
CREATE DATABASE IF NOT EXISTS lottery_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE lottery_db;

-- 公司表
CREATE TABLE IF NOT EXISTS companies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL COMMENT '公司代码，用于URL识别',
    name VARCHAR(100) NOT NULL COMMENT '公司名称',
    logo VARCHAR(255) COMMENT 'Logo URL',
    theme_color VARCHAR(20) DEFAULT '#667eea' COMMENT '主题颜色',

    -- 文案配置
    title VARCHAR(100) DEFAULT '幸运抽奖' COMMENT '系统标题',
    subtitle VARCHAR(200) COMMENT '副标题',
    welcome_text TEXT COMMENT '欢迎语',
    rules_text TEXT COMMENT '规则说明',
    draw_button_text VARCHAR(50) DEFAULT '点击抽奖' COMMENT '抽奖按钮文字',
    success_text VARCHAR(200) DEFAULT '恭喜中奖！' COMMENT '中奖恭喜语',

    -- 联系信息
    contact_name VARCHAR(100) COMMENT '联系人',
    contact_phone VARCHAR(20) COMMENT '联系电话',
    contact_email VARCHAR(100) COMMENT '联系邮箱',

    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    draw_count INT DEFAULT 1 COMMENT '每次抽奖人数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_code (code),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公司表';

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    company_id INT COMMENT '所属公司ID，NULL表示超级管理员',
    is_super_admin BOOLEAN DEFAULT FALSE COMMENT '是否超级管理员',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL,
    INDEX idx_username (username),
    INDEX idx_company_id (company_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL COMMENT '所属公司ID',
    phone VARCHAR(20) UNIQUE NOT NULL COMMENT '手机号',
    name VARCHAR(100) COMMENT '姓名',
    has_drawn BOOLEAN DEFAULT FALSE COMMENT '是否已抽奖',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    INDEX idx_company_id (company_id),
    INDEX idx_phone (phone),
    INDEX idx_has_drawn (has_drawn)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 奖项等级表
CREATE TABLE IF NOT EXISTS prize_levels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL COMMENT '所属公司ID',
    name VARCHAR(50) NOT NULL COMMENT '等级名称（如：一等奖、二等奖）',
    description VARCHAR(200) COMMENT '等级描述',
    probability DOUBLE NOT NULL COMMENT '中奖概率（0-1之间）',
    total_stock INT NOT NULL COMMENT '总库存',
    used_stock INT DEFAULT 0 COMMENT '已使用库存',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    INDEX idx_company_id (company_id),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='奖项等级表';

-- 奖品表
CREATE TABLE IF NOT EXISTS prizes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    level_id INT NOT NULL COMMENT '所属等级ID',
    name VARCHAR(100) NOT NULL COMMENT '奖品名称',
    image VARCHAR(255) COMMENT '奖品图片URL',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (level_id) REFERENCES prize_levels(id) ON DELETE CASCADE,
    INDEX idx_level_id (level_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='奖品表';

-- 抽奖记录表
CREATE TABLE IF NOT EXISTS draw_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL COMMENT '所属公司ID',
    user_id INT NOT NULL COMMENT '用户ID',
    level_id INT COMMENT '奖项等级ID',
    prize_id INT COMMENT '奖品ID',
    ip VARCHAR(50) COMMENT '抽奖IP地址',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (level_id) REFERENCES prize_levels(id) ON DELETE SET NULL,
    FOREIGN KEY (prize_id) REFERENCES prizes(id) ON DELETE SET NULL,
    INDEX idx_company_id (company_id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='抽奖记录表';

-- 插入默认超级管理员
-- 密码: admin123 (需要在应用程序中通过 bcrypt 加密后使用)
-- 注意：这里的密码是明文，实际使用时需要先注册或使用工具生成 bcrypt hash
INSERT INTO admins (username, password, company_id, is_super_admin)
VALUES ('admin', '$2a$10$YourBcryptHashedPasswordHere', NULL, TRUE)
ON DUPLICATE KEY UPDATE username=username;

-- 显示创建结果
SELECT 'Database initialization completed!' AS message;
SELECT COUNT(*) AS company_count FROM companies;
SELECT COUNT(*) AS admin_count FROM admins;
