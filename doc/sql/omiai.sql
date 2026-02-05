/*
 Navicat Premium Dump SQL

 Source Server         : 腾讯云
 Source Server Type    : MySQL
 Source Server Version : 80045 (8.0.45)
 Source Host           : 82.156.241.188:3306
 Source Schema         : omiai

 Target Server Type    : MySQL
 Target Server Version : 80045 (8.0.45)
 File Encoding         : 65001

 Date: 05/02/2026 19:58:49
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
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (4, '春季相亲大会', 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000', 1, 1, '/pages/activity/detail?id=1', '2026-02-01 10:12:53.131', '2026-02-01 10:12:53.131');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (5, '精英男士专场', 'https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000', 2, 1, '/pages/activity/detail?id=2', '2026-02-01 10:12:53.131', '2026-02-01 10:12:53.131');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (6, '牵手成功案例分享', 'https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000', 3, 1, '/pages/activity/detail?id=3', '2026-02-01 10:12:53.131', '2026-02-01 10:12:53.131');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (7, '春季相亲大会', 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000', 1, 1, '/pages/activity/detail?id=1', '2026-02-01 18:04:46.321', '2026-02-01 18:04:46.321');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (8, '精英男士专场', 'https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000', 2, 1, '/pages/activity/detail?id=2', '2026-02-01 18:04:46.321', '2026-02-01 18:04:46.321');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (9, '牵手成功案例分享', 'https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000', 3, 1, '/pages/activity/detail?id=3', '2026-02-01 18:04:46.321', '2026-02-01 18:04:46.321');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (10, '春季相亲大会', 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000', 1, 1, '/pages/activity/detail?id=1', '2026-02-01 18:05:02.001', '2026-02-01 18:05:02.001');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (11, '精英男士专场', 'https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000', 2, 1, '/pages/activity/detail?id=2', '2026-02-01 18:05:02.001', '2026-02-01 18:05:02.001');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (12, '牵手成功案例分享', 'https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000', 3, 1, '/pages/activity/detail?id=3', '2026-02-01 18:05:02.001', '2026-02-01 18:05:02.001');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (13, '春季相亲大会', 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000', 1, 1, '/pages/activity/detail?id=1', '2026-02-01 19:13:24.244', '2026-02-01 19:13:24.244');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (14, '精英男士专场', 'https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000', 2, 1, '/pages/activity/detail?id=2', '2026-02-01 19:13:24.244', '2026-02-01 19:13:24.244');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (15, '牵手成功案例分享', 'https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000', 3, 1, '/pages/activity/detail?id=3', '2026-02-01 19:13:24.244', '2026-02-01 19:13:24.244');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (16, '春季相亲大会', 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000', 1, 1, '/pages/activity/detail?id=1', '2026-02-01 19:25:52.593', '2026-02-01 19:25:52.593');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (17, '精英男士专场', 'https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000', 2, 1, '/pages/activity/detail?id=2', '2026-02-01 19:25:52.593', '2026-02-01 19:25:52.593');
INSERT INTO `banner` (`id`, `title`, `image_url`, `sort_order`, `status`, `link_url`, `created_at`, `updated_at`) VALUES (18, '牵手成功案例分享', 'https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000', 3, 1, '/pages/activity/detail?id=3', '2026-02-01 19:25:52.593', '2026-02-01 19:25:52.593');
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
  `house_status` tinyint DEFAULT NULL COMMENT '房产情况 1无房 2已购房 3贷款购房',
  `car_status` tinyint DEFAULT NULL COMMENT '车辆情况 1无车 2有车',
  `partner_requirements` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '对另一半要求(JSON)',
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
  `candidate_cache_json` longtext COLLATE utf8mb4_general_ci COMMENT '匹配候选缓存',
  `partner_id` bigint unsigned DEFAULT NULL COMMENT '当前匹配对象ID',
  `manager_id` bigint unsigned DEFAULT '0' COMMENT '归属红娘ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_client_partner` (`partner_id`),
  KEY `idx_phone` (`phone`),
  KEY `idx_client_phone` (`phone`),
  KEY `idx_client_manager` (`manager_id`)
) ENGINE=InnoDB AUTO_INCREMENT=355 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户档案表';

-- ----------------------------
-- Records of client
-- ----------------------------
BEGIN;
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (1, '何平', 1, '13811968136', '1991-03-14', '兔', 180, 83, 5, 1, '海淀区某小区', '父母退休，家庭和睦', 8000, '金融/银行', 2, 1, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.570', '2026-01-31 22:01:52.570', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (2, '李涛', 1, '13211861338', '1993-09-28', '猴', 173, 69, 3, 1, '东城区某小区', '父母退休，家庭和睦', 37000, '互联网/IT', 2, 2, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.574', '2026-01-31 22:01:52.574', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (3, '秦超', 1, '13373348298', '1991-01-18', '羊', 173, 68, 4, 1, '通州区某小区', '父母退休，家庭和睦', 21000, '自由职业', 2, 1, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.576', '2026-01-31 22:01:52.576', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (4, '尤秀', 2, '13773302801', '1991-07-01', '猪', 169, 59, 2, 2, '东城区某小区', '父母退休，家庭和睦', 7000, '公务员', 2, 2, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.578', '2026-01-31 22:01:52.578', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (40, '王英', 2, '13828281977', '1991-06-30', '牛', 158, 58, 3, 1, '海淀区某小区', '父母退休，家庭和睦', 36000, '金融/银行', 2, 2, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.625', '2026-01-31 22:01:52.625', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (41, '张平', 1, '13943032315', '1998-08-21', '猴', 178, 79, 3, 1, '海淀区某小区', '父母退休，家庭和睦', 33000, '企业高管', 2, 2, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.626', '2026-01-31 22:01:52.626', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (42, '何芳', 2, '13403239666', '1997-09-10', '龙', 158, 54, 4, 2, '丰台区某小区', '父母退休，家庭和睦', 24000, '互联网/IT', 1, 2, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.627', '2026-01-31 22:01:52.627', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (43, '朱伟', 1, '13016833980', '1995-01-14', '蛇', 174, 81, 2, 2, '海淀区某小区', '父母退休，家庭和睦', 33000, '公务员', 1, 2, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.628', '2026-01-31 22:01:52.628', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (44, '张娜', 2, '13267655110', '1993-05-22', '猪', 162, 45, 3, 1, '昌平区某小区', '父母退休，家庭和睦', 20000, '企业高管', 2, 1, '', '系统自动生成测试数据', '', '2026-01-31 22:01:52.629', '2026-01-31 22:01:52.629', NULL, 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
INSERT INTO `client` (`id`, `name`, `gender`, `phone`, `birthday`, `zodiac`, `height`, `weight`, `education`, `marital_status`, `address`, `family_description`, `income`, `profession`, `house_status`, `car_status`, `partner_requirements`, `remark`, `photos`, `created_at`, `updated_at`, `avatar`, `status`, `age`, `work_unit`, `position`, `house_address`, `candidate_cache_json`, `partner_id`, `manager_id`) VALUES (141, '朱敏', 2, '13590346602', '1998-02-10', '羊', 166, 58, 3, 1, '东城区某小区', '父母退休，家庭和睦', 20000, '公务员', 2, 1, '', '系统自动生成测试数据', '', '2026-02-01 10:12:53.124', '2026-02-01 10:12:53.124', '', 1, NULL, NULL, NULL, NULL, NULL, NULL, 0);
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
  `admin_id` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作管理员ID',
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
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号',
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
INSERT INTO `user` (`id`, `phone`, `nickname`, `avatar`, `role`, `wx_openid`, `created_at`, `updated_at`) VALUES (1, '18612571940', '管理员', '', 'admin', '', '2026-02-01 18:05:02.005', '2026-02-01 18:05:02.005');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
