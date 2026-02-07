# 全国行政区划表使用文档

## 表结构说明

### china_region 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| code | varchar(20) | 行政区划代码（主键） |
| name | varchar(100) | 地区名称 |
| parent_code | varchar(20) | 父级代码 |
| level | tinyint | 层级：1=省/直辖市，2=地级市，3=区/县 |
| pinyin | varchar(200) | 拼音（用于拼音搜索） |
| initial | char(1) | 首字母（A-Z快速索引） |
| is_hot | tinyint | 是否热门城市：0=否，1=是 |
| sort_order | int | 排序权重 |

## 行政区划代码规则

- **省级**：前2位 + 0000（如：110000=北京市）
- **市级**：前4位 + 00（如：110100=北京市市辖区）
- **县级**：6位完整编码（如：110101=东城区）

## 常用查询示例

### 1. 获取所有省份（第一级选择）
```sql
SELECT * FROM china_region WHERE level = 1 ORDER BY code;
```

### 2. 获取某省下的所有城市（级联查询）
```sql
-- 获取广东省下的城市
SELECT * FROM china_region 
WHERE parent_code = '440000' AND level = 2 
ORDER BY code;
```

### 3. 获取某城市下的所有区县（级联查询）
```sql
-- 获取深圳市下的区县
SELECT * FROM china_region 
WHERE parent_code = '440300' AND level = 3 
ORDER BY code;
```

### 4. 查询热门城市（快捷选择）
```sql
SELECT * FROM china_region 
WHERE is_hot = 1 AND level = 2 
ORDER BY sort_order ASC, code ASC;
```

### 5. 按名称模糊搜索
```sql
SELECT * FROM china_region 
WHERE name LIKE '%深圳%' 
ORDER BY level, code;
```

### 6. 按拼音搜索
```sql
-- 支持完整拼音或首字母
SELECT * FROM china_region 
WHERE pinyin LIKE '%shenzhen%' OR initial = 'S'
ORDER BY level, code;
```

### 7. 查询完整地区路径（省-市-区）
```sql
-- 查询深圳市的完整路径
WITH RECURSIVE region_path AS (
    SELECT * FROM china_region WHERE code = '440300'
    UNION ALL
    SELECT r.* FROM china_region r
    INNER JOIN region_path rp ON r.code = rp.parent_code
)
SELECT * FROM region_path ORDER BY level;
```

### 8. 按首字母分组查询（A-Z索引）
```sql
SELECT 
    initial, 
    COUNT(*) as city_count,
    GROUP_CONCAT(name ORDER BY sort_order) as cities
FROM china_region 
WHERE level = 2 AND is_hot = 1
GROUP BY initial 
ORDER BY initial;
```

## API接口

### 1. 获取所有省份
```http
GET /api/regions/provinces
```

### 2. 获取某省下的城市
```http
GET /api/regions/provinces/{code}/cities
```

### 3. 获取某城市下的区县
```http
GET /api/regions/cities/{code}/districts
```

### 4. 获取热门城市
```http
GET /api/regions/cities/hot
```

### 5. 搜索地区
```http
GET /api/regions/search?keyword=深圳
```

### 6. 获取地区详情
```http
GET /api/regions/{code}
```

## 在前端使用

### 级联选择器数据格式
```javascript
// 省份列表
const provinces = await getProvinces();

// 选中省份后加载城市
const cities = await getCitiesByProvince(selectedProvinceCode);

// 选中城市后加载区县
const districts = await getDistrictsByCity(selectedCityCode);
```

### 快捷选择热门城市
```javascript
const hotCities = await getHotCities();
// 显示：北京、上海、广州、深圳、杭州...
```

### 搜索功能
```javascript
const results = await searchRegions('shenzhen');
// 返回：深圳市（广东省）
```

## 数据初始化

1. 执行建表SQL：
```bash
mysql -u root -p omiai < doc/sql/china_region.sql
```

2. 插入基础数据（需要先运行生成脚本）：
```bash
cd scripts
go run generate_region_sql.go
```

3. 导入数据：
```bash
mysql -u root -p omiai < doc/sql/china_region_data.sql
```

## 扩展建议

1. **添加拼音数据**：使用 `github.com/mozillazg/go-pinyin` 库生成拼音
2. **添加经纬度**：便于地图展示和距离计算
3. **添加邮编**：便于邮寄资料
4. **添加区号**：便于电话沟通

## 与现有系统整合

### 在客户表中使用
```sql
-- 为客户表添加工作城市字段（已添加）
ALTER TABLE client ADD COLUMN work_city VARCHAR(128);

-- 查询某城市的所有客户
SELECT c.*, r.name as city_name
FROM client c
LEFT JOIN china_region r ON c.work_city = r.code
WHERE r.name LIKE '%深圳%';
```

### 筛选功能
前端筛选面板已添加工作城市输入框，支持模糊匹配城市名称。
