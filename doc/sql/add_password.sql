-- 用户表添加password字段
ALTER TABLE `user` 
ADD COLUMN `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码' AFTER `phone`;

-- 为现有用户设置默认密码123456（使用MD5加密）
-- 注意：实际应用中应该使用更安全的加密方式，如bcrypt
UPDATE `user` SET `password` = MD5('123456') WHERE `password` = '' OR `password` IS NULL;
