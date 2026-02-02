-- 增加 Phase 1 所需字段
ALTER TABLE `client` 
ADD COLUMN `manager_id` bigint(20) unsigned DEFAULT 0 COMMENT '归属红娘ID',
ADD COLUMN `is_public` tinyint(1) DEFAULT 1 COMMENT '是否公海 1是 0否',
ADD COLUMN `tags` text COMMENT '标签列表(JSON)';

-- 增加索引
ALTER TABLE `client` ADD INDEX `idx_manager_id` (`manager_id`);
ALTER TABLE `client` ADD INDEX `idx_is_public` (`is_public`);
