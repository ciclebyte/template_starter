-- 模板订阅预设变量功能数据库迁移脚本
-- 执行前请备份数据库

-- 创建模板订阅预设变量关联表
CREATE TABLE IF NOT EXISTS `template_variable_presets` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关联ID，自增主键',
    `template_id` bigint(20) unsigned NOT NULL COMMENT '所属模板ID',
    `var_preset_id` bigint(20) unsigned NOT NULL COMMENT '预设变量ID',
    `preset_path` varchar(500) NOT NULL DEFAULT '' COMMENT '预设中的变量路径 (如: user.name)',
    `mapped_name` varchar(200) NOT NULL DEFAULT '' COMMENT '映射到模板中的变量名',
    `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用此映射关系',
    `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` datetime DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_var_preset_id` (`var_preset_id`),
    KEY `idx_template_preset` (`template_id`, `var_preset_id`),
    UNIQUE KEY `uk_template_preset_path` (`template_id`, `var_preset_id`, `preset_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模板订阅预设变量关联表';

-- 检查表是否创建成功
SELECT 
    TABLE_NAME, 
    TABLE_COMMENT,
    CREATE_TIME 
FROM 
    information_schema.TABLES 
WHERE 
    TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'template_variable_presets';

-- 显示表结构
DESCRIBE template_variable_presets;