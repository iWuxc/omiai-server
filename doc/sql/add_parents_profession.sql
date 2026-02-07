-- 客户表新增父母工作字段
ALTER TABLE `client` 
ADD COLUMN `parents_profession` varchar(255) DEFAULT NULL COMMENT '父母工作' AFTER `partner_requirements`;
