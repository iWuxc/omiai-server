/*
 Navicat Premium Dump SQL

 Source Server         : 本地docker
 Source Server Type    : MySQL
 Source Server Version : 80034 (8.0.34)
 Source Host           : 127.0.0.1:3306
 Source Schema         : omiai

 Target Server Type    : MySQL
 Target Server Version : 80034 (8.0.34)
 File Encoding         : 65001

 Date: 02/02/2026 16:31:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- V2 Cleanup: Remove deprecated tables
-- ----------------------------
DROP TABLE IF EXISTS `match_request`;
DROP TABLE IF EXISTS `match_approval`;

-- ----------------------------
-- Table structure for match_record
-- ----------------------------
DROP TABLE IF EXISTS `match_record`;
CREATE TABLE `match_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `male_client_id` bigint unsigned NOT NULL COMMENT '男方ID',
  `female_client_id` bigint unsigned NOT NULL COMMENT '女方ID',
  `match_date` datetime DEFAULT NULL COMMENT '确认匹配时间',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '关系状态 1:相识 2:交往 3:稳定 4:订婚 5:结婚',
  `match_score` int DEFAULT 0 COMMENT '匹配得分',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `admin_id` varchar(64) DEFAULT NULL COMMENT '操作管理员ID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_male_client` (`male_client_id`),
  KEY `idx_female_client` (`female_client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='情侣档案表';

-- ----------------------------
-- Table structure for match_status_history
-- ----------------------------
DROP TABLE IF EXISTS `match_status_history`;
CREATE TABLE `match_status_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `match_record_id` bigint unsigned NOT NULL COMMENT '情侣档案ID',
  `previous_status` tinyint NOT NULL COMMENT '变更前状态',
  `current_status` tinyint NOT NULL COMMENT '变更后状态',
  `change_time` datetime NOT NULL COMMENT '变更时间',
  `reason` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '变更原因',
  `operator` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作人',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_match_record` (`match_record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='状态变更历史表';

-- ----------------------------
-- Table structure for follow_up_record
-- ----------------------------
DROP TABLE IF EXISTS `follow_up_record`;
CREATE TABLE `follow_up_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `match_record_id` bigint unsigned NOT NULL COMMENT '情侣档案ID',
  `follow_up_date` datetime NOT NULL COMMENT '回访日期',
  `method` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '回访方式(电话/微信/面谈)',
  `content` text COLLATE utf8mb4_general_ci COMMENT '回访内容',
  `feedback` text COLLATE utf8mb4_general_ci COMMENT '客户反馈',
  `satisfaction` tinyint DEFAULT NULL COMMENT '满意度 1-5',
  `attachments` text COLLATE utf8mb4_general_ci COMMENT '附件列表(JSON)',
  `next_follow_up_at` datetime DEFAULT NULL COMMENT '下次回访时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_match_record_follow` (`match_record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='回访记录表';

-- ----------------------------
-- Alter table client (Schema Updates)
-- ----------------------------
-- Add candidate_cache_json
ALTER TABLE `client` ADD COLUMN `candidate_cache_json` longtext COMMENT '匹配候选缓存';

-- Add/Modify partner_id to be nullable and unique
-- Note: Running both ADD and MODIFY to cover different initial states (idempotency relies on operator handling errors or using a tool that handles this)
-- Here we assume it might not exist or needs modification. 
-- In pure SQL script, we usually check existence, but standard SQL doesn't support IF NOT EXISTS for columns easily without procedure.
-- For Flyway/Migration tools, we expect this to be a new version that runs once.

-- Ensure partner_id exists and is nullable
ALTER TABLE `client` ADD COLUMN `partner_id` bigint unsigned DEFAULT NULL COMMENT '当前匹配对象ID';
ALTER TABLE `client` MODIFY COLUMN `partner_id` bigint unsigned DEFAULT NULL COMMENT '当前匹配对象ID';

-- Ensure manager_id exists
ALTER TABLE `client` ADD COLUMN `manager_id` bigint unsigned DEFAULT '0' COMMENT '归属红娘ID';

-- Add Unique Index on partner_id (Anti-duplicate match)
-- We use CREATE UNIQUE INDEX IF NOT EXISTS syntax if supported (MySQL 8.0+ supports IF NOT EXISTS for CREATE INDEX? No, mostly DROP first or checking).
-- We will try to create it. If it fails, it might exist.
CREATE UNIQUE INDEX `idx_client_partner` ON `client`(`partner_id`);

-- Add Index on manager_id
CREATE INDEX `idx_client_manager` ON `client`(`manager_id`);
