-- 提醒记录表
CREATE TABLE `reminder` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID（红娘ID）',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '提醒类型：1=回访提醒，2=生日提醒，3=纪念日提醒，4=流失预警',
  `client_id` bigint unsigned DEFAULT NULL COMMENT '关联客户ID',
  `match_record_id` bigint unsigned DEFAULT NULL COMMENT '关联匹配记录ID（纪念日提醒时用）',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '提醒标题',
  `content` text COMMENT '提醒内容',
  `remind_at` datetime NOT NULL COMMENT '提醒时间',
  `is_read` tinyint NOT NULL DEFAULT '0' COMMENT '是否已读：0=未读，1=已读',
  `is_done` tinyint NOT NULL DEFAULT '0' COMMENT '是否已完成：0=未完成，1=已完成',
  `priority` tinyint NOT NULL DEFAULT '2' COMMENT '优先级：1=低，2=中，3=高',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_remind_at` (`remind_at`),
  KEY `idx_is_done` (`is_done`),
  KEY `idx_type` (`type`),
  KEY `idx_user_status` (`user_id`, `is_done`, `remind_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提醒记录表';
