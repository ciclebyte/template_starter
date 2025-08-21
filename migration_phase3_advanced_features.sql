-- ================================================================================================
-- Template Starter 认证系统迁移 - 第三阶段：高级功能（多租户、分享机制）
-- 执行前请备份数据库！
-- 前置条件：必须先执行 migration_phase1_basic_auth.sql 和 migration_phase2_template_auth.sql
-- ================================================================================================

-- 1. 创建组织/租户表 (organizations)
CREATE TABLE `organizations` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '组织ID',
  `name` varchar(100) NOT NULL COMMENT '组织名称',
  `code` varchar(50) NOT NULL COMMENT '组织编码',
  `description` text COMMENT '组织描述',
  `logo` varchar(255) DEFAULT NULL COMMENT '组织logo',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=正常',
  `owner_id` bigint(20) NOT NULL COMMENT '组织拥有者ID',
  `member_limit` int(11) NOT NULL DEFAULT '100' COMMENT '成员上限',
  `template_limit` int(11) NOT NULL DEFAULT '1000' COMMENT '模板上限',
  `storage_limit` bigint(20) NOT NULL DEFAULT '1073741824' COMMENT '存储限制（字节，默认1GB）',
  `api_call_limit` int(11) NOT NULL DEFAULT '10000' COMMENT 'API调用限制（每月）',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间',
  `settings` json DEFAULT NULL COMMENT '组织设置',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_owner` (`owner_id`),
  KEY `idx_status` (`status`),
  KEY `idx_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';

-- 2. 创建组织成员表 (organization_members)
CREATE TABLE `organization_members` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) NOT NULL COMMENT '组织ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role` ENUM('owner', 'admin', 'member') NOT NULL DEFAULT 'member' COMMENT '组织内角色',
  `status` ENUM('active', 'pending', 'suspended') NOT NULL DEFAULT 'active' COMMENT '成员状态',
  `invited_by` bigint(20) DEFAULT NULL COMMENT '邀请者ID',
  `invited_at` datetime DEFAULT NULL COMMENT '邀请时间',
  `joined_at` datetime DEFAULT NULL COMMENT '加入时间',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间',
  `permissions` json DEFAULT NULL COMMENT '额外权限配置',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_org_user` (`organization_id`, `user_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_invited_by` (`invited_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织成员表';

-- 3. 创建模板分享表 (template_shares)
CREATE TABLE `template_shares` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `template_id` bigint(20) NOT NULL COMMENT '模板ID',
  `share_type` ENUM('user', 'role', 'organization', 'public_link') NOT NULL COMMENT '分享类型',
  `target_id` bigint(20) DEFAULT NULL COMMENT '目标ID（用户/角色/组织）',
  `share_code` varchar(32) DEFAULT NULL COMMENT '分享码（用于公开链接）',
  `permission` ENUM('read', 'use', 'copy', 'edit') NOT NULL DEFAULT 'read' COMMENT '权限类型',
  `shared_by` bigint(20) NOT NULL COMMENT '分享者ID',
  `access_count` int(11) NOT NULL DEFAULT '0' COMMENT '访问次数',
  `max_access_count` int(11) DEFAULT NULL COMMENT '最大访问次数限制',
  `password` varchar(255) DEFAULT NULL COMMENT '访问密码（可选）',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_template_share` (`template_id`, `share_type`, `target_id`),
  UNIQUE KEY `uk_share_code` (`share_code`),
  KEY `idx_target` (`share_type`, `target_id`),
  KEY `idx_shared_by` (`shared_by`),
  KEY `idx_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模板分享表';

-- 4. 创建模板分享访问记录表 (template_share_logs)
CREATE TABLE `template_share_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `share_id` bigint(20) NOT NULL COMMENT '分享ID',
  `template_id` bigint(20) NOT NULL COMMENT '模板ID',
  `user_id` bigint(20) DEFAULT NULL COMMENT '访问者ID（可为空）',
  `action` varchar(50) NOT NULL COMMENT '操作类型',
  `ip_address` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COMMENT '用户代理',
  `access_result` ENUM('success', 'denied', 'expired', 'limit_exceeded') NOT NULL COMMENT '访问结果',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_share` (`share_id`),
  KEY `idx_template` (`template_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模板分享访问记录表';

-- 5. 创建组织邀请表 (organization_invitations)
CREATE TABLE `organization_invitations` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) NOT NULL COMMENT '组织ID',
  `email` varchar(100) NOT NULL COMMENT '被邀请者邮箱',
  `role` ENUM('admin', 'member') NOT NULL DEFAULT 'member' COMMENT '邀请角色',
  `invited_by` bigint(20) NOT NULL COMMENT '邀请者ID',
  `invitation_code` varchar(64) NOT NULL COMMENT '邀请码',
  `status` ENUM('pending', 'accepted', 'declined', 'expired') NOT NULL DEFAULT 'pending' COMMENT '邀请状态',
  `message` text COMMENT '邀请消息',
  `expires_at` datetime NOT NULL COMMENT '过期时间',
  `accepted_at` datetime DEFAULT NULL COMMENT '接受时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_invitation_code` (`invitation_code`),
  KEY `idx_organization` (`organization_id`),
  KEY `idx_email` (`email`),
  KEY `idx_invited_by` (`invited_by`),
  KEY `idx_status` (`status`),
  KEY `idx_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织邀请表';

-- 6. 创建API调用统计表 (api_call_statistics)
CREATE TABLE `api_call_statistics` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) DEFAULT NULL COMMENT '组织ID',
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户ID',
  `endpoint` varchar(200) NOT NULL COMMENT 'API端点',
  `method` varchar(10) NOT NULL COMMENT 'HTTP方法',
  `call_date` date NOT NULL COMMENT '调用日期',
  `call_count` int(11) NOT NULL DEFAULT '1' COMMENT '调用次数',
  `success_count` int(11) NOT NULL DEFAULT '0' COMMENT '成功次数',
  `error_count` int(11) NOT NULL DEFAULT '0' COMMENT '错误次数',
  `total_duration_ms` bigint(20) NOT NULL DEFAULT '0' COMMENT '总耗时（毫秒）',
  `avg_duration_ms` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '平均耗时（毫秒）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_stats` (`organization_id`, `user_id`, `endpoint`, `method`, `call_date`),
  KEY `idx_call_date` (`call_date`),
  KEY `idx_endpoint` (`endpoint`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API调用统计表';

-- ================================================================================================
-- 更新现有表的外键关系
-- ================================================================================================

-- 为现有表添加外键约束
ALTER TABLE `users` ADD CONSTRAINT `fk_users_organization` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE SET NULL;
ALTER TABLE `templates` ADD CONSTRAINT `fk_templates_organization` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE SET NULL;

-- ================================================================================================
-- 初始化默认组织
-- ================================================================================================

-- 创建默认系统组织
INSERT INTO `organizations` (`name`, `code`, `description`, `owner_id`, `member_limit`, `template_limit`, `storage_limit`, `api_call_limit`) 
SELECT 
    '系统默认组织', 
    'system_default',
    '系统默认组织，用于管理公共资源',
    u.id,
    1000,
    10000, 
    10737418240, -- 10GB
    100000 -- 10万次/月
FROM users u WHERE u.username = 'admin' LIMIT 1;

-- 将系统管理员添加为默认组织的拥有者
INSERT INTO `organization_members` (`organization_id`, `user_id`, `role`, `status`, `joined_at`)
SELECT o.id, u.id, 'owner', 'active', NOW()
FROM organizations o, users u 
WHERE o.code = 'system_default' AND u.username = 'admin';

-- ================================================================================================
-- 创建高级视图
-- ================================================================================================

-- 组织成员权限视图
CREATE VIEW `organization_member_permissions` AS
SELECT 
    om.organization_id,
    o.name as organization_name,
    om.user_id,
    u.username,
    u.email,
    om.role as org_role,
    om.status as member_status,
    GROUP_CONCAT(DISTINCT r.code) as user_roles,
    GROUP_CONCAT(DISTINCT p.code) as permissions
FROM organization_members om
JOIN organizations o ON om.organization_id = o.id
JOIN users u ON om.user_id = u.id
LEFT JOIN user_roles ur ON u.id = ur.user_id AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
LEFT JOIN roles r ON ur.role_id = r.id AND r.status = 1
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE om.status = 'active' AND u.status = 1
GROUP BY om.organization_id, om.user_id;

-- 模板访问权限详细视图
CREATE VIEW `template_access_control` AS
SELECT 
    t.id as template_id,
    t.name as template_name,
    t.visibility,
    t.owner_id,
    owner.username as owner_username,
    t.organization_id,
    org.name as organization_name,
    CASE 
        WHEN t.visibility = 'public' THEN 'all_users'
        WHEN t.visibility = 'private' THEN owner.username
        WHEN t.visibility = 'organization' THEN org.name
        WHEN t.visibility = 'shared' THEN 'shared_users'
    END as access_description,
    COUNT(ts.id) as share_count,
    COUNT(tul.id) as usage_count
FROM templates t
LEFT JOIN users owner ON t.owner_id = owner.id
LEFT JOIN organizations org ON t.organization_id = org.id
LEFT JOIN template_shares ts ON t.id = ts.template_id AND (ts.expires_at IS NULL OR ts.expires_at > NOW())
LEFT JOIN template_usage_logs tul ON t.id = tul.template_id
GROUP BY t.id;

-- 组织资源使用统计视图
CREATE VIEW `organization_usage_stats` AS
SELECT 
    o.id as organization_id,
    o.name as organization_name,
    o.code as organization_code,
    COUNT(DISTINCT om.user_id) as member_count,
    o.member_limit,
    COUNT(DISTINCT t.id) as template_count,
    o.template_limit,
    COALESCE(SUM(LENGTH(tf.file_content)), 0) as storage_used,
    o.storage_limit,
    COALESCE(stats.monthly_api_calls, 0) as monthly_api_calls,
    o.api_call_limit,
    CASE 
        WHEN o.expires_at IS NULL THEN 'never'
        WHEN o.expires_at > NOW() THEN 'active'
        ELSE 'expired'
    END as status
FROM organizations o
LEFT JOIN organization_members om ON o.id = om.organization_id AND om.status = 'active'
LEFT JOIN templates t ON o.id = t.organization_id
LEFT JOIN template_files tf ON t.id = tf.template_id
LEFT JOIN (
    SELECT 
        organization_id, 
        SUM(call_count) as monthly_api_calls
    FROM api_call_statistics 
    WHERE call_date >= DATE_FORMAT(NOW(), '%Y-%m-01')
    GROUP BY organization_id
) stats ON o.id = stats.organization_id
WHERE o.status = 1
GROUP BY o.id;

-- ================================================================================================
-- 创建存储过程
-- ================================================================================================

DELIMITER ;;

-- 检查用户是否可以访问模板的存储过程
CREATE PROCEDURE `CheckTemplateAccess`(
    IN p_template_id BIGINT,
    IN p_user_id BIGINT,
    IN p_action VARCHAR(50),
    OUT p_has_access BOOLEAN,
    OUT p_access_reason VARCHAR(200)
)
BEGIN
    DECLARE v_visibility VARCHAR(20);
    DECLARE v_owner_id BIGINT;
    DECLARE v_org_id BIGINT;
    DECLARE v_user_org_id BIGINT;
    DECLARE v_share_count INT;
    DECLARE v_is_admin BOOLEAN DEFAULT FALSE;
    
    -- 获取模板信息
    SELECT visibility, owner_id, organization_id 
    INTO v_visibility, v_owner_id, v_org_id
    FROM templates 
    WHERE id = p_template_id;
    
    -- 获取用户组织信息
    SELECT organization_id INTO v_user_org_id FROM users WHERE id = p_user_id;
    
    -- 检查用户是否是管理员
    SELECT COUNT(*) > 0 INTO v_is_admin
    FROM user_roles ur 
    JOIN roles r ON ur.role_id = r.id 
    WHERE ur.user_id = p_user_id AND r.code IN ('super_admin', 'system_admin') 
    AND (ur.expires_at IS NULL OR ur.expires_at > NOW());
    
    SET p_has_access = FALSE;
    SET p_access_reason = 'Access denied';
    
    -- 检查访问权限
    IF v_visibility = 'public' THEN
        SET p_has_access = TRUE;
        SET p_access_reason = 'Public template';
    ELSEIF v_owner_id = p_user_id THEN
        SET p_has_access = TRUE;
        SET p_access_reason = 'Owner access';
    ELSEIF v_is_admin THEN
        SET p_has_access = TRUE;
        SET p_access_reason = 'Admin access';
    ELSEIF v_visibility = 'organization' AND v_org_id = v_user_org_id AND v_user_org_id IS NOT NULL THEN
        SET p_has_access = TRUE;
        SET p_access_reason = 'Organization member access';
    ELSEIF v_visibility = 'shared' THEN
        SELECT COUNT(*) INTO v_share_count
        FROM template_shares ts
        WHERE ts.template_id = p_template_id 
        AND (
            (ts.share_type = 'user' AND ts.target_id = p_user_id) OR
            (ts.share_type = 'organization' AND ts.target_id = v_user_org_id)
        )
        AND (ts.expires_at IS NULL OR ts.expires_at > NOW())
        AND (ts.max_access_count IS NULL OR ts.access_count < ts.max_access_count);
        
        IF v_share_count > 0 THEN
            SET p_has_access = TRUE;
            SET p_access_reason = 'Shared access';
        END IF;
    END IF;
END;;

DELIMITER ;

-- ================================================================================================
-- 插入示例数据
-- ================================================================================================

-- 创建示例组织
INSERT INTO `organizations` (`name`, `code`, `description`, `owner_id`, `member_limit`, `template_limit`) 
SELECT 
    'Demo Company', 
    'demo_company',
    '演示公司，用于展示多租户功能',
    u.id,
    50,
    500
FROM users u WHERE u.username = 'developer1' LIMIT 1;

-- 添加组织成员
INSERT INTO `organization_members` (`organization_id`, `user_id`, `role`, `status`, `joined_at`)
SELECT o.id, u.id, 'owner', 'active', NOW()
FROM organizations o, users u 
WHERE o.code = 'demo_company' AND u.username = 'developer1';

-- 创建示例模板分享
-- 假设存在ID为1的模板，将其分享给测试用户
INSERT INTO `template_shares` (`template_id`, `share_type`, `target_id`, `permission`, `shared_by`, `share_code`, `expires_at`)
SELECT 1, 'user', u.id, 'read', 
       (SELECT id FROM users WHERE username = 'developer1' LIMIT 1),
       CONCAT('share_', FLOOR(RAND() * 1000000)),
       DATE_ADD(NOW(), INTERVAL 30 DAY)
FROM users u WHERE u.username = 'testuser' LIMIT 1;

-- ================================================================================================
-- 更新现有数据
-- ================================================================================================

-- 将现有用户分配到默认组织
UPDATE users u 
JOIN organizations o ON o.code = 'system_default'
SET u.organization_id = o.id 
WHERE u.organization_id IS NULL AND u.username != 'admin';

-- ================================================================================================
-- 验证迁移结果
-- ================================================================================================

SELECT '=== 第三阶段迁移完成统计 ===' as info;

-- 检查新表
SELECT 'organizations' as table_name, COUNT(*) as count FROM organizations
UNION ALL
SELECT 'organization_members' as table_name, COUNT(*) as count FROM organization_members
UNION ALL
SELECT 'template_shares' as table_name, COUNT(*) as count FROM template_shares
UNION ALL
SELECT 'template_share_logs' as table_name, COUNT(*) as count FROM template_share_logs
UNION ALL
SELECT 'organization_invitations' as table_name, COUNT(*) as count FROM organization_invitations
UNION ALL
SELECT 'api_call_statistics' as table_name, COUNT(*) as count FROM api_call_statistics;

-- 检查组织统计
SELECT '=== 组织统计 ===' as info;
SELECT * FROM organization_usage_stats;

-- 检查模板访问控制
SELECT '=== 模板访问控制示例 ===' as info;
SELECT 
    template_name,
    visibility,
    owner_username,
    organization_name,
    access_description,
    share_count,
    usage_count
FROM template_access_control 
LIMIT 5;

-- 测试存储过程
SELECT '=== 访问权限检查测试 ===' as info;
CALL CheckTemplateAccess(1, 1, 'read', @has_access, @reason);
SELECT @has_access as has_access, @reason as access_reason;

COMMIT;

-- ================================================================================================
-- 第三阶段迁移完成提示
-- ================================================================================================

SELECT '
=== 第三阶段迁移完成！===

新增功能：
1. ✅ 多租户组织管理
2. ✅ 组织成员管理
3. ✅ 模板分享机制（用户、组织、公开链接）
4. ✅ 组织邀请系统
5. ✅ API调用统计
6. ✅ 高级权限视图
7. ✅ 模板访问检查存储过程

关键特性：
- 支持组织级别的资源隔离
- 灵活的模板分享机制
- 详细的使用统计和监控
- 完整的邀请和成员管理流程

注意事项：
1. 所有现有用户已分配到系统默认组织
2. 建议在应用层实现组织切换功能
3. 分享链接支持访问次数限制和密码保护
4. API统计可用于计费和限流

下一步：
- 实现组织管理界面
- 添加模板分享功能
- 实现API调用统计和限流
- 考虑添加计费系统集成
' as completion_message;