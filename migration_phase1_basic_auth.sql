-- ================================================================================================
-- Template Starter 认证系统迁移 - 第一阶段：基础认证
-- 执行前请备份数据库！
-- ================================================================================================

-- 1. 用户表 (users)
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=正常',
  `email_verified` tinyint(1) NOT NULL DEFAULT '0' COMMENT '邮箱验证状态',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
  `login_count` int(11) NOT NULL DEFAULT '0' COMMENT '登录次数',
  `organization_id` bigint(20) DEFAULT NULL COMMENT '所属组织ID（暂时保留，第三阶段使用）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_organization` (`organization_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 角色表 (roles)
CREATE TABLE `roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` text COMMENT '角色描述',
  `is_system` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否系统角色',
  `organization_id` bigint(20) DEFAULT NULL COMMENT '所属组织ID（NULL表示系统级角色）',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=正常',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code_org` (`code`, `organization_id`),
  KEY `idx_organization` (`organization_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 3. 权限表 (permissions)
CREATE TABLE `permissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `code` varchar(100) NOT NULL COMMENT '权限编码',
  `resource` varchar(50) NOT NULL COMMENT '资源类型',
  `action` varchar(50) NOT NULL COMMENT '操作类型',
  `description` text COMMENT '权限描述',
  `parent_id` bigint(20) DEFAULT NULL COMMENT '父权限ID',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_resource_action` (`resource`, `action`),
  KEY `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 4. 用户角色关联表 (user_roles)
CREATE TABLE `user_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `assigned_by` bigint(20) DEFAULT NULL COMMENT '分配者ID',
  `assigned_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_role` (`role_id`),
  KEY `idx_assigned_by` (`assigned_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 5. 角色权限关联表 (role_permissions)
CREATE TABLE `role_permissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `permission_id` bigint(20) NOT NULL COMMENT '权限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 6. 用户会话表 (user_sessions)
CREATE TABLE `user_sessions` (
  `id` varchar(128) NOT NULL COMMENT '会话ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `ip_address` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COMMENT '用户代理',
  `data` json DEFAULT NULL COMMENT '会话数据',
  `expires_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话表';

-- ================================================================================================
-- 初始化系统权限数据
-- ================================================================================================

-- 插入系统权限
INSERT INTO `permissions` (`name`, `code`, `resource`, `action`, `description`) VALUES
-- 模板管理权限
('查看模板', 'template:read', 'template', 'read', '查看模板列表和详情'),
('使用模板', 'template:use', 'template', 'use', '使用模板生成代码'),
('创建模板', 'template:create', 'template', 'create', '创建新模板'),
('编辑模板', 'template:edit', 'template', 'edit', '编辑模板内容'),
('删除模板', 'template:delete', 'template', 'delete', '删除模板'),
('分享模板', 'template:share', 'template', 'share', '分享模板给其他用户'),
('管理模板', 'template:manage', 'template', 'manage', '管理所有模板'),

-- 分类管理权限  
('查看分类', 'category:read', 'category', 'read', '查看分类列表'),
('管理分类', 'category:manage', 'category', 'manage', '创建、编辑、删除分类'),

-- 语言管理权限
('查看语言', 'language:read', 'language', 'read', '查看编程语言列表'),
('管理语言', 'language:manage', 'language', 'manage', '创建、编辑、删除编程语言'),

-- 标签管理权限
('查看标签', 'tag:read', 'tag', 'read', '查看标签列表'),
('管理标签', 'tag:manage', 'tag', 'manage', '创建、编辑、删除标签'),

-- 用户管理权限
('查看用户', 'user:read', 'user', 'read', '查看用户信息'),
('管理用户', 'user:manage', 'user', 'manage', '管理用户账户'),
('分配角色', 'user:assign_role', 'user', 'assign_role', '为用户分配角色'),

-- 组织管理权限（为第三阶段准备）
('查看组织', 'organization:read', 'organization', 'read', '查看组织信息'),
('管理组织', 'organization:manage', 'organization', 'manage', '管理组织设置'),
('邀请成员', 'organization:invite', 'organization', 'invite', '邀请新成员加入组织'),

-- 角色管理权限
('查看角色', 'role:read', 'role', 'read', '查看角色列表'),
('管理角色', 'role:manage', 'role', 'manage', '创建、编辑、删除角色'),
('分配权限', 'role:assign_permission', 'role', 'assign_permission', '为角色分配权限'),

-- 系统管理权限
('系统配置', 'system:config', 'system', 'config', '系统配置管理'),
('审计日志', 'system:audit', 'system', 'audit', '查看审计日志'),
('统计分析', 'system:analytics', 'system', 'analytics', '查看系统统计分析');

-- ================================================================================================
-- 初始化系统角色
-- ================================================================================================

INSERT INTO `roles` (`name`, `code`, `description`, `is_system`) VALUES
('超级管理员', 'super_admin', '拥有所有权限的超级管理员', 1),
('系统管理员', 'system_admin', '系统管理员，可管理系统配置和所有组织', 1),
('组织管理员', 'org_admin', '组织管理员，可管理组织内所有资源', 1),
('模板管理员', 'template_admin', '模板管理员，可管理组织内所有模板', 1),
('开发者', 'developer', '开发者角色，可创建和管理自己的模板', 1),
('普通用户', 'user', '普通用户，可使用和创建模板', 1),
('访客', 'guest', '访客用户，只能查看公开模板', 1);

-- ================================================================================================
-- 分配角色权限
-- ================================================================================================

-- 超级管理员拥有所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'super_admin';

-- 系统管理员权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'system_admin' AND p.code IN (
  'template:read', 'template:use', 'template:create', 'template:edit', 'template:delete', 'template:manage',
  'category:read', 'category:manage',
  'language:read', 'language:manage',
  'tag:read', 'tag:manage',
  'user:read', 'user:manage', 'user:assign_role',
  'organization:read', 'organization:manage',
  'role:read', 'role:manage', 'role:assign_permission',
  'system:config', 'system:audit', 'system:analytics'
);

-- 组织管理员权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'org_admin' AND p.code IN (
  'template:read', 'template:use', 'template:create', 'template:edit', 'template:delete', 'template:share',
  'category:read', 'language:read', 'tag:read',
  'user:read', 'user:assign_role',
  'organization:read', 'organization:invite',
  'role:read'
);

-- 模板管理员权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'template_admin' AND p.code IN (
  'template:read', 'template:use', 'template:create', 'template:edit', 'template:delete', 'template:share',
  'category:read', 'language:read', 'tag:read'
);

-- 开发者权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'developer' AND p.code IN (
  'template:read', 'template:use', 'template:create', 'template:edit', 'template:share',
  'category:read', 'language:read', 'tag:read'
);

-- 普通用户权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'user' AND p.code IN (
  'template:read', 'template:use', 'template:create', 'template:edit',
  'category:read', 'language:read', 'tag:read'
);

-- 访客权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id 
FROM `roles` r, `permissions` p 
WHERE r.code = 'guest' AND p.code IN (
  'template:read', 'category:read', 'language:read', 'tag:read'
);

-- ================================================================================================
-- 创建默认系统管理员用户
-- ================================================================================================

-- 插入默认管理员用户（密码：admin123，请在生产环境中修改）
-- 密码哈希是 bcrypt("admin123") 的结果
INSERT INTO `users` (`username`, `email`, `password_hash`, `nickname`, `status`, `email_verified`) VALUES
('admin', 'admin@example.com', '$2a$10$rQ7Q1VxJ2QJ2Q1VxJ2QJ2O1VxJ2QJ2Q1VxJ2QJ2Q1VxJ2QJ2Q1VxJ2', '系统管理员', 1, 1);

-- 为默认管理员分配超级管理员角色
INSERT INTO `user_roles` (`user_id`, `role_id`, `assigned_by`) 
SELECT u.id, r.id, u.id
FROM `users` u, `roles` r 
WHERE u.username = 'admin' AND r.code = 'super_admin';

-- ================================================================================================
-- 完成第一阶段迁移
-- ================================================================================================

-- 验证迁移结果
SELECT '=== 迁移完成统计 ===' as info;
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'roles' as table_name, COUNT(*) as count FROM roles
UNION ALL  
SELECT 'permissions' as table_name, COUNT(*) as count FROM permissions
UNION ALL
SELECT 'user_roles' as table_name, COUNT(*) as count FROM user_roles
UNION ALL
SELECT 'role_permissions' as table_name, COUNT(*) as count FROM role_permissions;

SELECT '=== 默认管理员信息 ===' as info;
SELECT u.username, u.email, u.nickname, r.name as role_name 
FROM users u 
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.username = 'admin';

-- 注意事项：
-- 1. 默认管理员密码为 admin123，请立即修改！
-- 2. 邮箱地址请修改为真实邮箱
-- 3. 建议在生产环境中重新生成密码哈希