-- =============================================
-- 全国行政区划表 (适用于婚恋项目)
-- 支持省市区三级查询、拼音搜索、热门城市标记
-- =============================================

DROP TABLE IF EXISTS `china_region`;

CREATE TABLE `china_region` (
  `code` varchar(20) NOT NULL COMMENT '行政区划代码',
  `name` varchar(100) NOT NULL COMMENT '地区名称',
  `parent_code` varchar(20) DEFAULT NULL COMMENT '父级行政区划代码',
  `level` tinyint NOT NULL DEFAULT 1 COMMENT '层级: 1=省/直辖市, 2=地级市, 3=区/县',
  `pinyin` varchar(200) DEFAULT NULL COMMENT '拼音（用于拼音搜索）',
  `initial` char(1) DEFAULT NULL COMMENT '首字母（A-Z快速索引）',
  `is_hot` tinyint NOT NULL DEFAULT 0 COMMENT '是否热门城市: 0=否, 1=是（用于快捷选择）',
  `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序权重（热门城市排在前面）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`code`),
  KEY `idx_parent_code` (`parent_code`),
  KEY `idx_level` (`level`),
  KEY `idx_name` (`name`),
  KEY `idx_pinyin` (`pinyin`),
  KEY `idx_initial` (`initial`),
  KEY `idx_is_hot` (`is_hot`),
  KEY `idx_level_hot` (`level`, `is_hot`),
  KEY `idx_parent_level` (`parent_code`, `level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='全国行政区划表';

-- =============================================
-- 插入热门城市标记（便于快捷选择）
-- =============================================

-- 直辖市
UPDATE china_region SET is_hot = 1, sort_order = 10 WHERE code IN ('110000', '120000', '310000', '500000');

-- 热门省会城市及计划单列市
UPDATE china_region SET is_hot = 1, sort_order = 20 WHERE code IN (
  '130100', -- 石家庄
  '140100', -- 太原
  '210100', -- 沈阳
  '210200', -- 大连
  '220100', -- 长春
  '230100', -- 哈尔滨
  '320100', -- 南京
  '320200', -- 无锡
  '320500', -- 苏州
  '330100', -- 杭州
  '330200', -- 宁波
  '340100', -- 合肥
  '350100', -- 福州
  '350200', -- 厦门
  '360100', -- 南昌
  '370100', -- 济南
  '370200', -- 青岛
  '410100', -- 郑州
  '420100', -- 武汉
  '430100', -- 长沙
  '440100', -- 广州
  '440300', -- 深圳
  '440400', -- 珠海
  '440600', -- 佛山
  '441900', -- 东莞
  '442000', -- 中山
  '450100', -- 南宁
  '460100', -- 海口
  '510100', -- 成都
  '520100', -- 贵阳
  '530100', -- 昆明
  '610100', -- 西安
  '620100', -- 兰州
  '630100', -- 西宁
  '640100', -- 银川
  '650100'  -- 乌鲁木齐
);

-- 其他热门地级市
UPDATE china_region SET is_hot = 1, sort_order = 30 WHERE code IN (
  '320400', -- 常州
  '320600', -- 南通
  '321000', -- 扬州
  '321100', -- 镇江
  '330300', -- 温州
  '330400', -- 嘉兴
  '330600', -- 绍兴
  '331000', -- 台州
  '350500', -- 泉州
  '370600', -- 烟台
  '370700', -- 潍坊
  '371000', -- 威海
  '371300', -- 临沂
  '410300', -- 洛阳
  '420600', -- 襄阳
  '430200', -- 株洲
  '430600', -- 岳阳
  '440500', -- 汕头
  '440700', -- 江门
  '440800', -- 湛江
  '441300', -- 惠州
  '441800', -- 清远
  '500100', -- 重庆主城
  '510700', -- 绵阳
  '511000'  -- 内江
);

-- =============================================
-- 常用查询SQL示例
-- =============================================

-- 1. 查询所有省份（用于第一级选择）
-- SELECT * FROM china_region WHERE level = 1 ORDER BY code;

-- 2. 查询某省份下的所有城市（用于第二级级联）
-- SELECT * FROM china_region WHERE parent_code = '440000' ORDER BY code; -- 广东省下的城市

-- 3. 查询某城市下的所有区县（用于第三级级联）
-- SELECT * FROM china_region WHERE parent_code = '440300' ORDER BY code; -- 深圳市下的区县

-- 4. 查询热门城市（快捷选择）
-- SELECT * FROM china_region WHERE is_hot = 1 AND level = 2 ORDER BY sort_order, code;

-- 5. 按名称模糊搜索（支持中文）
-- SELECT * FROM china_region WHERE name LIKE '%深圳%' ORDER BY level, code;

-- 6. 按拼音搜索（支持拼音首字母）
-- SELECT * FROM china_region WHERE pinyin LIKE '%shenzhen%' OR initial = 'S' ORDER BY level, code;

-- 7. 查询完整地区路径（省-市-区）
-- SELECT 
--   c.name as city_name,
--   p.name as province_name
-- FROM china_region c
-- LEFT JOIN china_region p ON c.parent_code = p.code
-- WHERE c.code = '440300';

-- 8. 查询所有城市（排除省直辖县）
-- SELECT * FROM china_region WHERE level = 2 ORDER BY parent_code, code;

-- 9. 按首字母分组查询城市（A-Z索引）
-- SELECT initial, COUNT(*) as count, GROUP_CONCAT(name) as cities
-- FROM china_region 
-- WHERE level = 2 AND is_hot = 1
-- GROUP BY initial 
-- ORDER BY initial;
