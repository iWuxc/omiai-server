-- 客户表新增工作城市字段
ALTER TABLE `client` 
ADD COLUMN `work_city` varchar(128) DEFAULT NULL COMMENT '工作城市' AFTER `profession`;
