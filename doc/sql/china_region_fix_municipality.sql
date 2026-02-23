-- =============================================
-- 修复直辖市层级断层问题
-- 北京、天津、上海、重庆等直辖市在标准行政区划中是省级（Level 1）
-- 但其下辖的区县（如朝阳区）是县级（Level 3）
-- 导致前端级联选择器无法找到中间的“市级”（Level 2）节点
-- 本脚本插入虚拟的“市辖区”节点来连接 Level 1 和 Level 3
-- =============================================

-- 1. 插入 Level 2 的市辖区节点
INSERT INTO `china_region` (`code`, `name`, `parent_code`, `level`, `pinyin`, `initial`, `sort_order`) VALUES
('110100', '北京市', '110000', 2, 'beijingshi', 'B', 0),
('120100', '天津市', '120000', 2, 'tianjinshi', 'T', 0),
('310100', '上海市', '310000', 2, 'shanghaishi', 'S', 0),
('500100', '重庆市', '500000', 2, 'chongqingshi', 'C', 0);

-- 2. 更新直辖市下属区县的父节点
-- 将原本直接挂在省级（如110000）下的区县（Level 3），挂载到新插入的市级（110100）下
-- 注意：标准区县代码通常前4位就是 1101，所以逻辑上是匹配的

-- 更新北京区县 (1101xx -> parent 110100)
UPDATE `china_region` SET `parent_code` = '110100' WHERE `parent_code` = '110000' AND `level` = 3;

-- 更新天津区县 (1201xx -> parent 120100)
UPDATE `china_region` SET `parent_code` = '120100' WHERE `parent_code` = '120000' AND `level` = 3;

-- 更新上海区县 (3101xx -> parent 310100)
UPDATE `china_region` SET `parent_code` = '310100' WHERE `parent_code` = '310000' AND `level` = 3;

-- 更新重庆区县 (5001xx -> parent 500100)
UPDATE `china_region` SET `parent_code` = '500100' WHERE `parent_code` = '500000' AND `level` = 3;

-- =============================================
-- 验证查询
-- =============================================
-- SELECT * FROM china_region WHERE parent_code = '110000'; -- 应该查到 110100
-- SELECT * FROM china_region WHERE parent_code = '110100'; -- 应该查到朝阳区、海淀区等
