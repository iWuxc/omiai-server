-- =============================================
-- 全国行政区划基础数据（主要省市）
-- 用于快速测试和演示
-- 完整数据请运行 generate_region_sql.go 生成
-- =============================================

INSERT INTO `china_region` (`code`, `name`, `parent_code`, `level`, `pinyin`, `initial`, `sort_order`) VALUES
-- 直辖市
('110000', '北京市', NULL, 1, 'beijing', 'B', 0),
('120000', '天津市', NULL, 1, 'tianjin', 'T', 0),
('310000', '上海市', NULL, 1, 'shanghai', 'S', 0),
('500000', '重庆市', NULL, 1, 'chongqing', 'C', 0),

-- 河北省
('130000', '河北省', NULL, 1, 'hebei', 'H', 0),
('130100', '石家庄市', '130000', 2, 'shijiazhuang', 'S', 0),
('130200', '唐山市', '130000', 2, 'tangshan', 'T', 0),
('130300', '秦皇岛市', '130000', 2, 'qinhuangdao', 'Q', 0),

-- 山西省
('140000', '山西省', NULL, 1, 'shanxi', 'S', 0),
('140100', '太原市', '140000', 2, 'taiyuan', 'T', 0),

-- 辽宁省
('210000', '辽宁省', NULL, 1, 'liaoning', 'L', 0),
('210100', '沈阳市', '210000', 2, 'shenyang', 'S', 0),
('210200', '大连市', '210000', 2, 'dalian', 'D', 0),

-- 吉林省
('220000', '吉林省', NULL, 1, 'jilin', 'J', 0),
('220100', '长春市', '220000', 2, 'changchun', 'C', 0),

-- 黑龙江省
('230000', '黑龙江省', NULL, 1, 'heilongjiang', 'H', 0),
('230100', '哈尔滨市', '230000', 2, 'haerbin', 'H', 0),

-- 江苏省
('320000', '江苏省', NULL, 1, 'jiangsu', 'J', 0),
('320100', '南京市', '320000', 2, 'nanjing', 'N', 0),
('320200', '无锡市', '320000', 2, 'wuxi', 'W', 0),
('320400', '常州市', '320000', 2, 'changzhou', 'C', 0),
('320500', '苏州市', '320000', 2, 'suzhou', 'S', 0),
('320600', '南通市', '320000', 2, 'nantong', 'N', 0),

-- 浙江省
('330000', '浙江省', NULL, 1, 'zhejiang', 'Z', 0),
('330100', '杭州市', '330000', 2, 'hangzhou', 'H', 0),
('330200', '宁波市', '330000', 2, 'ningbo', 'N', 0),
('330300', '温州市', '330000', 2, 'wenzhou', 'W', 0),

-- 安徽省
('340000', '安徽省', NULL, 1, 'anhui', 'A', 0),
('340100', '合肥市', '340000', 2, 'hefei', 'H', 0),

-- 福建省
('350000', '福建省', NULL, 1, 'fujian', 'F', 0),
('350100', '福州市', '350000', 2, 'fuzhou', 'F', 0),
('350200', '厦门市', '350000', 2, 'xiamen', 'X', 0),

-- 江西省
('360000', '江西省', NULL, 1, 'jiangxi', 'J', 0),
('360100', '南昌市', '360000', 2, 'nanchang', 'N', 0),

-- 山东省
('370000', '山东省', NULL, 1, 'shandong', 'S', 0),
('370100', '济南市', '370000', 2, 'jinan', 'J', 0),
('370200', '青岛市', '370000', 2, 'qingdao', 'Q', 0),
('370600', '烟台市', '370000', 2, 'yantai', 'Y', 0),

-- 河南省
('410000', '河南省', NULL, 1, 'henan', 'H', 0),
('410100', '郑州市', '410000', 2, 'zhengzhou', 'Z', 0),

-- 湖北省
('420000', '湖北省', NULL, 1, 'hubei', 'H', 0),
('420100', '武汉市', '420000', 2, 'wuhan', 'W', 0),

-- 湖南省
('430000', '湖南省', NULL, 1, 'hunan', 'H', 0),
('430100', '长沙市', '430000', 2, 'changsha', 'C', 0),

-- 广东省
('440000', '广东省', NULL, 1, 'guangdong', 'G', 0),
('440100', '广州市', '440000', 2, 'guangzhou', 'G', 0),
('440300', '深圳市', '440000', 2, 'shenzhen', 'S', 0),
('440400', '珠海市', '440000', 2, 'zhuhai', 'Z', 0),
('440600', '佛山市', '440000', 2, 'foshan', 'F', 0),
('441900', '东莞市', '440000', 2, 'dongguan', 'D', 0),

-- 海南省
('460000', '海南省', NULL, 1, 'hainan', 'H', 0),
('460100', '海口市', '460000', 2, 'haikou', 'H', 0),

-- 四川省
('510000', '四川省', NULL, 1, 'sichuan', 'S', 0),
('510100', '成都市', '510000', 2, 'chengdu', 'C', 0),

-- 贵州省
('520000', '贵州省', NULL, 1, 'guizhou', 'G', 0),
('520100', '贵阳市', '520000', 2, 'guiyang', 'G', 0),

-- 云南省
('530000', '云南省', NULL, 1, 'yunnan', 'Y', 0),
('530100', '昆明市', '530000', 2, 'kunming', 'K', 0),

-- 陕西省
('610000', '陕西省', NULL, 1, 'shanxi', 'S', 0),
('610100', '西安市', '610000', 2, 'xian', 'X', 0);

-- 更新热门城市标记
UPDATE china_region SET is_hot = 1, sort_order = 10 WHERE code IN ('110000', '120000', '310000', '500000');
UPDATE china_region SET is_hot = 1, sort_order = 20 WHERE code = '130100'; -- 石家庄
UPDATE china_region SET is_hot = 1, sort_order = 30 WHERE code = '210100'; -- 沈阳
UPDATE china_region SET is_hot = 1, sort_order = 40 WHERE code = '210200'; -- 大连
UPDATE china_region SET is_hot = 1, sort_order = 50 WHERE code = '320100'; -- 南京
UPDATE china_region SET is_hot = 1, sort_order = 60 WHERE code = '320200'; -- 无锡
UPDATE china_region SET is_hot = 1, sort_order = 70 WHERE code = '320500'; -- 苏州
UPDATE china_region SET is_hot = 1, sort_order = 80 WHERE code = '330100'; -- 杭州
UPDATE china_region SET is_hot = 1, sort_order = 90 WHERE code = '330200'; -- 宁波
UPDATE china_region SET is_hot = 1, sort_order = 100 WHERE code = '350200'; -- 厦门
UPDATE china_region SET is_hot = 1, sort_order = 110 WHERE code = '370100'; -- 济南
UPDATE china_region SET is_hot = 1, sort_order = 120 WHERE code = '370200'; -- 青岛
UPDATE china_region SET is_hot = 1, sort_order = 130 WHERE code = '410100'; -- 郑州
UPDATE china_region SET is_hot = 1, sort_order = 140 WHERE code = '420100'; -- 武汉
UPDATE china_region SET is_hot = 1, sort_order = 150 WHERE code = '430100'; -- 长沙
UPDATE china_region SET is_hot = 1, sort_order = 160 WHERE code = '440100'; -- 广州
UPDATE china_region SET is_hot = 1, sort_order = 170 WHERE code = '440300'; -- 深圳
UPDATE china_region SET is_hot = 1, sort_order = 180 WHERE code = '440600'; -- 佛山
UPDATE china_region SET is_hot = 1, sort_order = 190 WHERE code = '441900'; -- 东莞
UPDATE china_region SET is_hot = 1, sort_order = 200 WHERE code = '510100'; -- 成都
UPDATE china_region SET is_hot = 1, sort_order = 210 WHERE code = '520100'; -- 贵阳
UPDATE china_region SET is_hot = 1, sort_order = 220 WHERE code = '530100'; -- 昆明
UPDATE china_region SET is_hot = 1, sort_order = 230 WHERE code = '610100'; -- 西安
