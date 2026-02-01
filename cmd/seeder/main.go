package main

import (
	"fmt"
	"math/rand"
	"time"

	biz_omiai "omiai-server/internal/biz/omiai"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config
const DSN = "root:onex(#)666@tcp(127.0.0.1:3306)/omiai?charset=utf8mb4&parseTime=True&loc=Local"

// Data Sources
var (
	lastNames        = []string{"赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋", "沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张"}
	firstNamesMale   = []string{"伟", "强", "磊", "洋", "勇", "军", "杰", "涛", "超", "明", "刚", "平", "辉", "伟", "志强", "建国", "志明"}
	firstNamesFemale = []string{"芳", "娜", "敏", "静", "秀", "娟", "英", "华", "慧", "巧", "美", "静", "丽", "霞", "燕", "琳", "雪"}
	professions      = []string{"互联网/IT", "金融/银行", "教育/教师", "医疗/医生", "公务员", "自由职业", "企业高管", "销售", "设计师", "律师"}
	addresses        = []string{"朝阳区", "海淀区", "西城区", "东城区", "丰台区", "通州区", "昌平区"}
	zodiacs          = []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	educations       = []int8{3, 4, 3, 4, 3, 2, 5} // 2:大专 3:本科 4:硕士 5:博士
)

func main() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate
	db.AutoMigrate(&biz_omiai.Client{}, &biz_omiai.Banner{}, &biz_omiai.MatchRecord{}, &biz_omiai.FollowUpRecord{})

	fmt.Println("Start seeding...")

	rand.Seed(time.Now().UnixNano())

	// Generate 50 Clients
	for i := 0; i < 50; i++ {
		gender := int8(rand.Intn(2) + 1) // 1 or 2
		name := generateName(gender)

		client := &biz_omiai.Client{
			Name:              name,
			Gender:            gender,
			Phone:             fmt.Sprintf("13%d%08d", rand.Intn(10), rand.Intn(100000000)),
			Birthday:          generateBirthday(),
			Zodiac:            zodiacs[rand.Intn(len(zodiacs))],
			Height:            generateHeight(gender),
			Weight:            generateWeight(gender),
			Education:         educations[rand.Intn(len(educations))],
			MaritalStatus:     int8(rand.Intn(2) + 1), // 1:未婚 2:离异
			Address:           addresses[rand.Intn(len(addresses))] + "某小区",
			Income:            (rand.Intn(40) + 5) * 1000, // 5000 - 45000
			Profession:        professions[rand.Intn(len(professions))],
			HouseStatus:       int8(rand.Intn(2) + 1),
			CarStatus:         int8(rand.Intn(2) + 1),
			FamilyDescription: "父母退休，家庭和睦",
			Remark:            "系统自动生成测试数据",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := db.Create(client).Error; err != nil {
			fmt.Printf("Failed to create client %s: %v\n", client.Name, err)
		} else {
			fmt.Printf("Created client: %s (ID: %d)\n", client.Name, client.ID)
		}
	}

	fmt.Println("Seeding completed!")

	// Generate 3 Banners
	banners := []*biz_omiai.Banner{
		{
			Title:     "春季相亲大会",
			ImageURL:  "https://images.unsplash.com/photo-1511795409834-ef04bbd61622?auto=format&fit=crop&q=80&w=1000",
			SortOrder: 1,
			Status:    biz_omiai.BannerStatusEnable,
			LinkUrl:   "/pages/activity/detail?id=1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "精英男士专场",
			ImageURL:  "https://images.unsplash.com/photo-1516589174184-c68526673fd0?auto=format&fit=crop&q=80&w=1000",
			SortOrder: 2,
			Status:    biz_omiai.BannerStatusEnable,
			LinkUrl:   "/pages/activity/detail?id=2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "牵手成功案例分享",
			ImageURL:  "https://images.unsplash.com/photo-1519741497674-611481863552?auto=format&fit=crop&q=80&w=1000",
			SortOrder: 3,
			Status:    biz_omiai.BannerStatusEnable,
			LinkUrl:   "/pages/activity/detail?id=3",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, b := range banners {
		if err := db.Create(b).Error; err != nil {
			fmt.Printf("Failed to create banner %s: %v\n", b.Title, err)
		} else {
			fmt.Printf("Created banner: %s (ID: %d)\n", b.Title, b.ID)
		}
	}

	fmt.Println("All seeding completed!")
}

func generateName(gender int8) string {
	ln := lastNames[rand.Intn(len(lastNames))]
	var fn string
	if gender == 1 {
		fn = firstNamesMale[rand.Intn(len(firstNamesMale))]
	} else {
		fn = firstNamesFemale[rand.Intn(len(firstNamesFemale))]
	}
	return ln + fn
}

func generateBirthday() string {
	min := time.Date(1990, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2000, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format("2006-01-02")
}

func generateHeight(gender int8) int {
	if gender == 1 {
		return 170 + rand.Intn(15) // 170-185
	}
	return 158 + rand.Intn(12) // 158-170
}

func generateWeight(gender int8) int {
	if gender == 1 {
		return 65 + rand.Intn(20)
	}
	return 45 + rand.Intn(15)
}
