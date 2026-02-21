/*
 Navicat Premium Dump SQL

 Source Server         : 本地dcoker
 Source Server Type    : MySQL
 Source Server Version : 80034 (8.0.34)
 Source Host           : 127.0.0.1:3306
 Source Schema         : omiai

 Target Server Type    : MySQL
 Target Server Version : 80034 (8.0.34)
 File Encoding         : 65001

 Date: 21/02/2026 18:36:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for banner
-- ----------------------------
DROP TABLE IF EXISTS `banner`;
CREATE TABLE `banner` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `image_url` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `sort_order` bigint unsigned DEFAULT NULL,
  `status` tinyint DEFAULT NULL,
  `link_url` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of banner
-- ----------------------------
BEGIN;
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (1, '春季相亲大会', 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000', 1, 1, '/pages/activity/detail?id=1', '2026-02-01 08:04:25.157', '2026-02-01 08:04:25.157');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (2, '精英男士专场', 'https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000', 2, 1, '/pages/activity/detail?id=2', '2026-02-01 08:04:25.157', '2026-02-01 08:04:25.157');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (3, '牵手成功案例分享', 'https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000', 3, 1, '/pages/activity/detail?id=3', '2026-02-01 08:04:25.157', '2026-02-01 08:04:25.157');
COMMIT;

-- ----------------------------
-- Table structure for client
-- ----------------------------
DROP TABLE IF EXISTS `client`;
CREATE TABLE `client` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '姓名',
  `gender` tinyint DEFAULT NULL COMMENT '性别 1男 2女',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '联系电话',
  `birthday` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '出生年月',
  `zodiac` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '属相',
  `height` bigint DEFAULT NULL COMMENT '身高cm',
  `weight` bigint DEFAULT NULL COMMENT '体重kg',
  `education` tinyint DEFAULT NULL COMMENT '学历',
  `marital_status` tinyint DEFAULT NULL COMMENT '婚姻状况 1未婚 2已婚 3离异 4丧偶',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '家庭住址',
  `family_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '家庭成员描述',
  `income` bigint DEFAULT NULL COMMENT '月收入',
  `profession` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '具体工作',
  `work_city` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '工作城市',
  `house_status` tinyint DEFAULT NULL COMMENT '房产情况 1无房 2已购房 3贷款购房',
  `car_status` tinyint DEFAULT NULL COMMENT '车辆情况 1无车 2有车',
  `partner_requirements` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '对另一半要求(JSON)',
  `parents_profession` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '父母工作',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '红娘备注',
  `photos` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '照片URL列表(JSON)',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '头像URL',
  `status` tinyint DEFAULT '1' COMMENT '状态 1单身 2匹配中 3已匹配 4停止服务',
  `age` bigint DEFAULT NULL COMMENT '年龄',
  `work_unit` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作单位',
  `position` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '职位',
  `house_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '买房地址',
  `candidate_cache_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '匹配候选缓存',
  `partner_id` bigint unsigned DEFAULT NULL COMMENT '当前匹配对象ID',
  `manager_id` bigint unsigned DEFAULT '0' COMMENT '归属红娘ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_client_partner` (`partner_id`),
  KEY `idx_phone` (`phone`),
  KEY `idx_client_phone` (`phone`),
  KEY `idx_client_manager` (`manager_id`)
) ENGINE=InnoDB AUTO_INCREMENT=357 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户档案表';

-- ----------------------------
-- Records of client
-- ----------------------------
BEGIN;
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `work_city`, `house_status`, `car_status`, `partner_requirements`, `parents_profession`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (356, '伍六一', 2, '18612571940', '2001-01', '猴', 178, 75, 3, 1, '河北省张家口', '爸爸妈妈和我', 20000, '程序员', '北京', 2, 2, '性格温和大方，能够经得起共同郝菲尔那个号否', '都是普通工人', '阿斯顿发顺丰撒上档次闲杂v代发说法中山西正擦', '[\"https://omiai-server-1252902619.cos.ap-nanjing.myqcloud.com/uploads/20260221/92ab36bf-c507-4463-9a34-aba41aacbfb8.jpg\",\"https://omiai-server-1252902619.cos.ap-nanjing.myqcloud.com/uploads/20260221/c279b7d5-b0ea-453f-ba1b-dd1ff0cf15cc.jpg\",\"https://omiai-server-1252902619.cos.ap-nanjing.myqcloud.com/uploads/20260221/b03a9f53-49a6-476f-8744-9f121d82b62b.jpg\",\"https://omiai-server-1252902619.cos.ap-nanjing.myqcloud.com/uploads/20260221/ab2e7119-491f-42a0-bcdb-23e097fc3823.jpg\"]', '2026-02-21 17:53:23.619', '2026-02-21 17:53:23.619', 'https://omiai-server-1252902619.cos.ap-nanjing.myqcloud.com/uploads/20260221/69daa181-5dad-4692-838b-3cd6954833c7.jpg', 1, 25, '互联网行业', '普通职工', '河北省张家口', '', NULL, 0);
COMMIT;

-- ----------------------------
-- Table structure for follow_up_record
-- ----------------------------
DROP TABLE IF EXISTS `follow_up_record`;
CREATE TABLE `follow_up_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `match_record_id` bigint unsigned NOT NULL COMMENT '情侣档案ID',
  `follow_up_date` datetime NOT NULL COMMENT '回访日期',
  `method` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '回访方式(电话/微信/面谈)',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '回访内容',
  `feedback` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '客户反馈',
  `satisfaction` tinyint DEFAULT NULL COMMENT '满意度 1-5',
  `attachments` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '附件列表(JSON)',
  `next_follow_up_at` datetime DEFAULT NULL COMMENT '下次回访时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_match_record_follow` (`match_record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='回访记录表';

-- ----------------------------
-- Records of follow_up_record
-- ----------------------------
BEGIN;
COMMIT;

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
  `match_score` int DEFAULT '0' COMMENT '匹配得分',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注',
  `admin_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作管理员ID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_male_client` (`male_client_id`),
  KEY `idx_female_client` (`female_client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='情侣档案表';

-- ----------------------------
-- Records of match_record
-- ----------------------------
BEGIN;
COMMIT;

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
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '变更原因',
  `operator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作人',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_match_record` (`match_record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='状态变更历史表';

-- ----------------------------
-- Records of match_status_history
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for reminder
-- ----------------------------
DROP TABLE IF EXISTS `reminder`;
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
  KEY `idx_user_status` (`user_id`,`is_done`,`remind_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='提醒记录表';

-- ----------------------------
-- Records of reminder
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '头像',
  `role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'operator' COMMENT '角色 admin/operator',
  `wx_openid` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '微信OpenID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_phone` (`phone`),
  KEY `idx_user_wx_open_id` (`wx_openid`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
BEGIN;
INSERT INTO `user` (`id`, `phone`, `password`, `nickname`, `avatar`, `role`, `wx_openid`, `created_at`, `updated_at`) VALUES (1, '18612571940', 'e10adc3949ba59abbe56e057f20f883e', '管理员', '', 'admin', '', '2026-02-01 18:05:02.005', '2026-02-01 18:05:02.005');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
