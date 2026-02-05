ALTER TABLE `follow_up_record` ADD COLUMN `feedback` text COLLATE utf8mb4_general_ci COMMENT '客户反馈' AFTER `content`;
ALTER TABLE `follow_up_record` ADD COLUMN `satisfaction` tinyint DEFAULT NULL COMMENT '满意度 1-5' AFTER `feedback`;
ALTER TABLE `follow_up_record` ADD COLUMN `attachments` text COLLATE utf8mb4_general_ci COMMENT '附件列表(JSON)' AFTER `satisfaction`;
