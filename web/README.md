## template_starter
### 项目介绍

开发一个前后台分离的项目模板生成器网站，这是前端部分

#### 功能点

- 语言、分类的展示、查询
- 隐蔽式增删语言、分类、修改
- 模板展示、新增模板、编辑模板（naive ui）
- 使用模板可以根据配置的变量生成项目、并下载projiect_name.zip

### 页面设计
 
仅仅做示意
```
-----------------------------------------------
Template Starter     首页 模板        【搜索】
-----------------------------------------------
快速启动您的项目
从数百个精心设计的项目模板中选择，一键生成完整项目结构
            【搜索】
-------------------------------------------------
                热门分类
                xxxx
------------------------------------------------
    web应用   移动应用    桌面应用   api服务
------------------------------------------------
推荐模板

vue全栈应用
xxx
使用模板
-------------------------------------------------
            footer
```

模板示意
```
-----------------------------------------------
Template Starter     首页 模板        【搜索】
-----------------------------------------------
    分类
    
    web应用 移动应用 桌面应用 api服务 springboot应用 gin应用(tag标签分类)


vue全栈应用   
xxx
使用模板

----------------------------------------------

footer

```



### 技术栈

- vue3
- naive ui
- pinia
- vue-router
- Monaco Editor​​ 

### 数据库结构参考
```sql
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '分类ID，自增主键',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称，唯一',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分类描述',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '分类图标标识或URL',
  `sort` int NOT NULL DEFAULT 0 COMMENT '数字越大越靠前',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  INDEX `idx_category_name`(`name` ASC) USING BTREE COMMENT '分类名称索引'
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '模板分类表，用于组织和管理模板的分类信息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for languages
-- ----------------------------
DROP TABLE IF EXISTS `languages`;
CREATE TABLE `languages`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '语言ID，自增主键',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '语言名称（如JavaScript、Python）',
  `display_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '语言显示名称',
  `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '语言代码（如js、py）',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '语言图标标识或URL',
  `color` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '语言代表色（十六进制）',
  `is_popular` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否热门语言',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `code`(`code` ASC) USING BTREE,
  INDEX `idx_language_name`(`name` ASC) USING BTREE COMMENT '语言名称索引',
  INDEX `idx_language_code`(`code` ASC) USING BTREE COMMENT '语言代码索引'
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '编程语言表，存储支持的编程语言信息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for template_files
-- ----------------------------
DROP TABLE IF EXISTS `template_files`;
CREATE TABLE `template_files`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '文件ID，自增主键',
  `template_id` bigint UNSIGNED NOT NULL COMMENT '所属模板ID',
  `file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件路径（相对路径）',
  `file_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '文件内容',
  `file_size` int UNSIGNED NULL DEFAULT NULL COMMENT '文件大小（字节）',
  `is_directory` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否为目录',
  `md5` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'md5',
  `parent_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '父目录ID，如果是文件则指向所属目录',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_file_template`(`template_id` ASC) USING BTREE COMMENT '模板ID索引',
  INDEX `idx_file_parent`(`parent_id` ASC) USING BTREE COMMENT '父目录ID索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '模板文件表，存储模板的实际文件内容和目录结构' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for template_languages
-- ----------------------------
DROP TABLE IF EXISTS `template_languages`;
CREATE TABLE `template_languages`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID，自增主键',
  `template_id` bigint UNSIGNED NOT NULL COMMENT '关联的模板ID',
  `language_id` int UNSIGNED NOT NULL COMMENT '关联的语言ID',
  `is_primary` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否主要语言',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '模板与语言关联表，记录模板支持的语言' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for template_variables
-- ----------------------------
DROP TABLE IF EXISTS `template_variables`;
CREATE TABLE `template_variables`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '变量ID，自增主键',
  `template_id` bigint UNSIGNED NOT NULL COMMENT '所属模板ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '变量名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '变量描述',
  `default_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '变量默认值',
  `is_required` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否为必填变量',
  `validation_regex` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '变量值验证正则表达式',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_template_variable`(`template_id` ASC, `name` ASC) USING BTREE COMMENT '模板ID和变量名的唯一约束',
  INDEX `idx_variable_template`(`template_id` ASC) USING BTREE COMMENT '模板ID索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '模板变量表，定义模板中可配置的变量参数' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for templates
-- ----------------------------
DROP TABLE IF EXISTS `templates`;
CREATE TABLE `templates`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '模板ID，自增主键',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模板名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '模板详细描述',
  `category_id` int UNSIGNED NOT NULL COMMENT '所属分类ID',
  `is_featured` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否推荐模板',
  `logo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '模板logo图片URL',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_template_name`(`name` ASC) USING BTREE COMMENT '模板名称索引',
  INDEX `idx_template_category`(`category_id` ASC) USING BTREE COMMENT '分类ID索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '模板主表，存储所有模板的基本信息' ROW_FORMAT = DYNAMIC;
```

