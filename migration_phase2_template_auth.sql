-- ================================================================================================
-- Template Starter 认证系统迁移 - 第二阶段：模板权限控制
-- 执行前请备份数据库！
-- 前置条件：必须先执行 migration_phase1_basic_auth.sql
-- ================================================================================================

-- 1. 扩展现有模板表，增加权限控制字段
ALTER TABLE `templates` ADD COLUMN `visibility` ENUM('public', 'private', 'organization', 'shared') NOT NULL DEFAULT 'public' COMMENT '可见性：public=公开，private=私有，organization=组织内，shared=指定分享';
ALTER TABLE `templates` ADD COLUMN `owner_id` bigint(20) DEFAULT NULL COMMENT '模板拥有者ID';
ALTER TABLE `templates` ADD COLUMN `organization_id` bigint(20) DEFAULT NULL COMMENT '所属组织ID';

-- 添加索引
ALTER TABLE `templates` ADD INDEX `idx_owner` (`owner_id`);
ALTER TABLE `templates` ADD INDEX `idx_organization` (`organization_id`);
ALTER TABLE `templates` ADD INDEX `idx_visibility` (`visibility`);

-- 2. 创建模板使用记录表 (template_usage_logs)
CREATE TABLE `template_usage_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `template_id` bigint(20) NOT NULL COMMENT '模板ID',
  `user_id` bigint(20) DEFAULT NULL COMMENT '使用者ID（可为空，支持匿名用户）',
  `action` varchar(50) NOT NULL COMMENT '操作类型：view,use,download,copy,fork',
  `ip_address` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COMMENT '用户代理',
  `request_data` json DEFAULT NULL COMMENT '请求数据（如变量值）',
  `response_size` int(11) DEFAULT NULL COMMENT '响应大小（字节）',
  `duration_ms` int(11) DEFAULT NULL COMMENT '执行时长（毫秒）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_template` (`template_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created` (`created_at`),
  KEY `idx_template_user` (`template_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模板使用记录表';

-- 3. 创建系统审计日志表 (audit_logs)
CREATE TABLE `audit_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL COMMENT '操作用户ID',
  `organization_id` bigint(20) DEFAULT NULL COMMENT '组织ID',
  `action` varchar(100) NOT NULL COMMENT '操作类型',
  `resource_type` varchar(50) NOT NULL COMMENT '资源类型',
  `resource_id` varchar(100) DEFAULT NULL COMMENT '资源ID',
  `old_data` json DEFAULT NULL COMMENT '变更前数据',
  `new_data` json DEFAULT NULL COMMENT '变更后数据',
  `ip_address` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COMMENT '用户代理',
  `result` ENUM('success', 'failure') NOT NULL COMMENT '操作结果',
  `error_message` text COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_organization` (`organization_id`),
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_action` (`action`),
  KEY `idx_result` (`result`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统审计日志表';

-- ================================================================================================
-- 数据迁移：将现有模板设置合理的默认值
-- ================================================================================================

-- 将所有现有模板设置为公开，这样不会影响现有的使用
UPDATE `templates` SET 
  `visibility` = 'public',
  `owner_id` = NULL,
  `organization_id` = NULL
WHERE `visibility` IS NULL OR `owner_id` IS NULL;

-- ================================================================================================
-- 添加外键约束
-- ================================================================================================

-- 添加外键约束（如果需要严格的数据完整性）
-- 注意：如果有历史数据不一致，可能需要先清理数据
-- ALTER TABLE `templates` ADD CONSTRAINT `fk_templates_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;
-- ALTER TABLE `template_usage_logs` ADD CONSTRAINT `fk_usage_template` FOREIGN KEY (`template_id`) REFERENCES `templates` (`id`) ON DELETE CASCADE;
-- ALTER TABLE `template_usage_logs` ADD CONSTRAINT `fk_usage_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;
-- ALTER TABLE `audit_logs` ADD CONSTRAINT `fk_audit_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;

-- ================================================================================================
-- 创建视图和存储过程（可选）
-- ================================================================================================

-- 创建用户权限视图，方便查询用户拥有的所有权限
CREATE VIEW `user_permissions_view` AS
SELECT 
    u.id as user_id,
    u.username,
    u.email,
    r.code as role_code,
    r.name as role_name,
    p.code as permission_code,
    p.name as permission_name,
    p.resource,
    p.action
FROM users u
JOIN user_roles ur ON u.id = ur.user_id AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
JOIN roles r ON ur.role_id = r.id AND r.status = 1
JOIN role_permissions rp ON r.id = rp.role_id  
JOIN permissions p ON rp.permission_id = p.id
WHERE u.status = 1;

-- 创建模板访问权限视图
CREATE VIEW `template_access_view` AS
SELECT 
    t.id as template_id,
    t.name as template_name,
    t.visibility,
    t.owner_id,
    t.organization_id,
    CASE 
        WHEN t.visibility = 'public' THEN 'public'
        WHEN t.visibility = 'private' THEN CONCAT('user:', t.owner_id)
        WHEN t.visibility = 'organization' THEN CONCAT('org:', t.organization_id)
        ELSE 'shared'
    END as access_rule
FROM templates t;

-- ================================================================================================
-- 插入示例审计日志（展示数据结构）
-- ================================================================================================

-- 记录迁移操作
INSERT INTO `audit_logs` (
    `user_id`, 
    `action`, 
    `resource_type`, 
    `resource_id`, 
    `new_data`, 
    `ip_address`, 
    `result`
) VALUES (
    1, -- 假设管理员用户ID为1
    'database_migration',
    'system',
    'phase2',
    JSON_OBJECT(
        'migration', 'phase2_template_auth',
        'tables_added', JSON_ARRAY('template_usage_logs', 'audit_logs'),
        'columns_added', JSON_OBJECT('templates', JSON_ARRAY('visibility', 'owner_id', 'organization_id')),
        'views_created', JSON_ARRAY('user_permissions_view', 'template_access_view')
    ),
    '127.0.0.1',
    'success'
);

-- ================================================================================================
-- 创建示例数据（用于测试）
-- ================================================================================================

-- 创建测试用户（可选，用于测试）
INSERT INTO `users` (`username`, `email`, `password_hash`, `nickname`, `status`, `email_verified`) VALUES
('testuser', 'test@example.com', '$2a$10$rQ7Q1VxJ2QJ2Q1VxJ2QJ2O1VxJ2QJ2Q1VxJ2QJ2Q1VxJ2QJ2Q1VxJ2', '测试用户', 1, 1),
('developer1', 'dev1@example.com', '$2a$10$rQ7Q1VxJ2QJ2Q1VxJ2QJ2O1VxJ2QJ2Q1VxJ2QJ2Q1VxJ2QJ2Q1VxJ2', '开发者1', 1, 1);

-- 为测试用户分配角色
INSERT INTO `user_roles` (`user_id`, `role_id`, `assigned_by`) 
SELECT u.id, r.id, 1
FROM `users` u, `roles` r 
WHERE u.username = 'testuser' AND r.code = 'user'
UNION ALL
SELECT u.id, r.id, 1
FROM `users` u, `roles` r 
WHERE u.username = 'developer1' AND r.code = 'developer';

-- 如果有现有模板，可以为一些模板分配测试所有者
-- UPDATE `templates` SET `owner_id` = (SELECT id FROM users WHERE username = 'developer1' LIMIT 1), `visibility` = 'private' WHERE id = 1;

-- ================================================================================================
-- 验证迁移结果
-- ================================================================================================

SELECT '=== 第二阶段迁移完成统计 ===' as info;

-- 检查表结构
SELECT '检查模板表新字段' as step;
SELECT 
    COUNT(*) as total_templates,
    SUM(CASE WHEN visibility = 'public' THEN 1 ELSE 0 END) as public_templates,
    SUM(CASE WHEN visibility = 'private' THEN 1 ELSE 0 END) as private_templates,
    SUM(CASE WHEN owner_id IS NOT NULL THEN 1 ELSE 0 END) as owned_templates
FROM templates;

-- 检查新表
SELECT 'template_usage_logs' as table_name, COUNT(*) as count FROM template_usage_logs
UNION ALL
SELECT 'audit_logs' as table_name, COUNT(*) as count FROM audit_logs;

-- 检查视图
SELECT '检查权限视图' as step;
SELECT COUNT(DISTINCT user_id) as users_with_permissions FROM user_permissions_view;
SELECT COUNT(*) as templates_in_access_view FROM template_access_view;

-- 检查用户权限示例
SELECT '=== 用户权限示例 ===' as info;
SELECT 
    upv.username,
    upv.role_name,
    COUNT(DISTINCT upv.permission_code) as permission_count,
    GROUP_CONCAT(DISTINCT upv.resource ORDER BY upv.resource) as resources
FROM user_permissions_view upv 
GROUP BY upv.user_id, upv.username, upv.role_name
LIMIT 5;

COMMIT;

-- ================================================================================================
-- 第二阶段迁移完成提示
-- ================================================================================================

SELECT '
=== 第二阶段迁移完成！===

新增功能：
1. ✅ 模板权限控制（可见性：public, private, organization, shared）
2. ✅ 模板所有权管理
3. ✅ 模板使用记录跟踪  
4. ✅ 系统审计日志
5. ✅ 用户权限查询视图

注意事项：
1. 所有现有模板默认设置为 public，保持向前兼容
2. 建议在应用层实现权限检查逻辑
3. 审计日志会记录所有关键操作
4. 可以通过 user_permissions_view 视图快速查询用户权限

下一步：
- 在后端实现权限检查中间件
- 在前端添加权限控制组件
- 考虑是否需要执行第三阶段迁移（多租户功能）
' as completion_message;