package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// 常见城市的拼音映射（简化版）
var pinyinMap = map[string]string{
	"北京市": "beijing", "天津市": "tianjin", "上海市": "shanghai", "重庆市": "chongqing",
	"河北省": "hebei", "山西省": "shanxi", "辽宁省": "liaoning", "吉林省": "jilin",
	"黑龙江省": "heilongjiang", "江苏省": "jiangsu", "浙江省": "zhejiang", "安徽省": "anhui",
	"福建省": "fujian", "江西省": "jiangxi", "山东省": "shandong", "河南省": "henan",
	"湖北省": "hubei", "湖南省": "hunan", "广东省": "guangdong", "海南省": "hainan",
	"四川省": "sichuan", "贵州省": "guizhou", "云南省": "yunnan", "陕西省": "shanxi",
	"甘肃省": "gansu", "青海省": "qinghai", "台湾省": "taiwan", "内蒙古自治区": "neimenggu",
	"广西壮族自治区": "guangxi", "西藏自治区": "xizang", "宁夏回族自治区": "ningxia",
	"新疆维吾尔自治区": "xinjiang", "香港特别行政区": "xianggang", "澳门特别行政区": "aomen",
}

// 热门城市列表
var hotCities = []string{
	"110000", "120000", "310000", "500000", // 直辖市
	"130100", "140100", "210100", "210200", "220100", "230100", // 东北省会
	"320100", "320200", "320400", "320500", "320600", // 江苏
	"330100", "330200", "330300", "330600", // 浙江
	"340100", "350100", "350200", "360100", "370100", "370200", // 华东
	"410100", "420100", "430100", // 华中
	"440100", "440300", "440400", "440600", "441300", "441900", "442000", // 华南
	"450100", "460100", // 西南
	"510100", "520100", "530100", "610100", "620100", // 西部省会
	"630100", "640100", "650100",
}

func main() {
	// 读取JSON文件
	data, err := os.ReadFile("/Users/edy/Downloads/data.json")
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 解析JSON
	var regionMap map[string]string
	if err := json.Unmarshal(data, &regionMap); err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		return
	}

	// 按code排序
	var codes []string
	for code := range regionMap {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	// 生成SQL文件
	file, err := os.Create("/Users/edy/apps/go/src/github.com/iwuxc/omiai-server/doc/sql/china_region_full.sql")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 写入文件头
	file.WriteString("-- =============================================\n")
	file.WriteString("-- 全国行政区划完整数据\n")
	file.WriteString("-- 数据来源：国家统计局最新行政区划代码\n")
	file.WriteString(fmt.Sprintf("-- 生成时间：%s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString("-- 数据条数：" + fmt.Sprintf("%d", len(codes)) + "\n")
	file.WriteString("-- =============================================\n\n")

	file.WriteString("-- 清空现有数据\n")
	file.WriteString("TRUNCATE TABLE `china_region`;\n\n")
	file.WriteString("-- 插入数据\n")
	file.WriteString("INSERT INTO `china_region` (`code`, `name`, `parent_code`, `level`, `pinyin`, `initial`, `is_hot`, `sort_order`) VALUES\n")

	// 生成INSERT语句
	var values []string
	hotCityMap := make(map[string]int)
	for i, code := range hotCities {
		hotCityMap[code] = (i + 1) * 10
	}

	for _, code := range codes {
		name := regionMap[code]
		parentCode := getParentCode(code)
		level := getLevel(code)
		pinyin := getPinyin(name)
		initial := getInitial(pinyin)

		isHot := 0
		sortOrder := 0
		if order, ok := hotCityMap[code]; ok {
			isHot = 1
			sortOrder = order
		}

		value := fmt.Sprintf("('%s', '%s', %s, %d, '%s', '%s', %d, %d)",
			code,
			escapeSQL(name),
			parentCode,
			level,
			pinyin,
			initial,
			isHot,
			sortOrder,
		)
		values = append(values, value)
	}

	// 写入所有值
	for i, v := range values {
		if i < len(values)-1 {
			file.WriteString(v + ",\n")
		} else {
			file.WriteString(v + ";\n")
		}
	}

	// 写入统计信息
	file.WriteString("\n-- =============================================\n")
	file.WriteString("-- 数据统计\n")
	file.WriteString("-- =============================================\n")

	// 统计各级别数量
	level1Count := 0
	level2Count := 0
	level3Count := 0
	for _, code := range codes {
		level := getLevel(code)
		switch level {
		case 1:
			level1Count++
		case 2:
			level2Count++
		case 3:
			level3Count++
		}
	}

	file.WriteString(fmt.Sprintf("-- 省级行政区：%d 个\n", level1Count))
	file.WriteString(fmt.Sprintf("-- 地级行政区：%d 个\n", level2Count))
	file.WriteString(fmt.Sprintf("-- 县级行政区：%d 个\n", level3Count))
	file.WriteString(fmt.Sprintf("-- 总计：%d 个\n", len(codes)))
	file.WriteString(fmt.Sprintf("-- 热门城市：%d 个\n", len(hotCities)))
	file.WriteString("-- =============================================\n")

	fmt.Printf("✅ 成功生成SQL文件！\n")
	fmt.Printf("📊 数据统计：\n")
	fmt.Printf("   省级：%d 个\n", level1Count)
	fmt.Printf("   地级：%d 个\n", level2Count)
	fmt.Printf("   县级：%d 个\n", level3Count)
	fmt.Printf("   总计：%d 个\n", len(codes))
	fmt.Printf("📁 文件位置：doc/sql/china_region_full.sql\n")
}

// getParentCode 获取父级行政区划代码
func getParentCode(code string) string {
	if len(code) != 6 {
		return "NULL"
	}

	// 省级（后四位为0000）
	if code[2:] == "0000" {
		return "NULL"
	}

	// 市级（后两位为00）
	if code[4:] == "00" {
		return fmt.Sprintf("'%s0000'", code[:2])
	}

	// 县级（直辖市下的区，前两位是11,12,31,50，且第3-4位为00）
	if (code[:2] == "11" || code[:2] == "12" || code[:2] == "31" || code[:2] == "50") && code[2:4] == "00" {
		return fmt.Sprintf("'%s0000'", code[:2])
	}

	// 县级（普通地级市下的区县）
	return fmt.Sprintf("'%s00'", code[:4])
}

// getLevel 获取层级
func getLevel(code string) int {
	if len(code) != 6 {
		return 3
	}

	// 省级
	if code[2:] == "0000" {
		return 1
	}

	// 市级
	if code[4:] == "00" {
		return 2
	}

	// 县级
	return 3
}

// getPinyin 获取拼音
func getPinyin(name string) string {
	// 先查映射表
	if py, ok := pinyinMap[name]; ok {
		return py
	}

	// 去掉"市"、"县"、"区"等后缀，查基础名称
	simpleName := strings.TrimSuffix(name, "市")
	simpleName = strings.TrimSuffix(simpleName, "县")
	simpleName = strings.TrimSuffix(simpleName, "区")
	simpleName = strings.TrimSuffix(simpleName, "省")
	simpleName = strings.TrimSuffix(simpleName, "自治区")
	simpleName = strings.TrimSuffix(simpleName, "特别行政区")

	if py, ok := pinyinMap[simpleName]; ok {
		return py
	}

	// 返回空字符串，让应用层处理或使用拼音库
	return ""
}

// getInitial 获取首字母
func getInitial(pinyin string) string {
	if len(pinyin) == 0 {
		return ""
	}
	return strings.ToUpper(pinyin[:1])
}

// escapeSQL 转义SQL字符串
func escapeSQL(s string) string {
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\\", "\\\\")
	return s
}
