-- 添加角色字段到users和admins表
-- 执行时间: 2026-01-24

-- 1. 为users表添加role字段
ALTER TABLE `users`
ADD COLUMN `role` VARCHAR(50) NOT NULL DEFAULT 'user'
COMMENT '用户角色: user'
AFTER `phone`;

-- 2. 为admins表添加role字段
ALTER TABLE `admins`
ADD COLUMN `role` VARCHAR(50) NOT NULL DEFAULT 'admin'
COMMENT '管理员角色: admin, super_admin'
AFTER `password`;

-- 3. 为现有用户设置默认角色
UPDATE `users` SET `role` = 'user' WHERE `role` IS NULL OR `role` = '';

-- 4. 为现有管理员设置角色（根据is_super_admin字段）
UPDATE `admins` SET `role` = 'super_admin' WHERE `is_super_admin` = TRUE;
UPDATE `admins` SET `role` = 'admin' WHERE `is_super_admin` = FALSE AND (`role` IS NULL OR `role` = '');

-- 5. 添加索引以提升查询性能
ALTER TABLE `users` ADD INDEX `idx_role` (`role`);
ALTER TABLE `admins` ADD INDEX `idx_role` (`role`);

-- 6. 验证数据
SELECT 'Users table:' as info;
SELECT id, username, role FROM users LIMIT 5;

SELECT 'Admins table:' as info;
SELECT id, username, role, is_super_admin FROM admins LIMIT 5;
