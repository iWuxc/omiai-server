package cron

import (
	"context"
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"
	"time"
)

// ReminderService 提醒服务
type ReminderService struct {
	db           *data.DB
	reminderRepo biz_omiai.ReminderInterface
	clientRepo   biz_omiai.ClientInterface
	matchRepo    biz_omiai.MatchInterface
}

func NewReminderService(db *data.DB, reminderRepo biz_omiai.ReminderInterface, clientRepo biz_omiai.ClientInterface, matchRepo biz_omiai.MatchInterface) *ReminderService {
	return &ReminderService{
		db:           db,
		reminderRepo: reminderRepo,
		clientRepo:   clientRepo,
		matchRepo:    matchRepo,
	}
}

// GenerateDailyReminders 生成每日提醒
func (s *ReminderService) GenerateDailyReminders(ctx context.Context) error {
	// log.Info("开始生成每日提醒...")

	// 1. 生成回访提醒（7天未联系的客户）
	if err := s.generateFollowUpReminders(ctx); err != nil {
		// log.Errorf("生成回访提醒失败: %v", err)
	}

	// 2. 生成生日提醒（未来3天内生日的客户）
	if err := s.generateBirthdayReminders(ctx); err != nil {
		// log.Errorf("生成生日提醒失败: %v", err)
	}

	// 3. 生成纪念日提醒（已匹配客户的相识纪念日）
	if err := s.generateAnniversaryReminders(ctx); err != nil {
		// log.Errorf("生成纪念日提醒失败: %v", err)
	}

	// 4. 生成流失预警（30天未联系的客户）
	if err := s.generateChurnRiskReminders(ctx); err != nil {
		// log.Errorf("生成流失预警失败: %v", err)
	}

	// log.Info("每日提醒生成完成")
	return nil
}

// generateFollowUpReminders 生成回访提醒（7天未联系）
func (s *ReminderService) generateFollowUpReminders(ctx context.Context) error {
	// 查询所有需要回访的客户
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	// 查询7天内没有跟进记录的客户（通过客户创建时间或最后更新时间判断）
	// 这里简化处理：查询所有状态为单身或已匹配的客户
	var clients []*biz_omiai.Client
	err := s.db.DB.WithContext(ctx).Model(&biz_omiai.Client{}).
		Where("status IN ?", []int8{1, 3}). // ClientStatusSingle=1, ClientStatusMatched=3
		Where("updated_at <= ?", sevenDaysAgo).
		Find(&clients).Error
	if err != nil {
		return err
	}

	for _, client := range clients {
		// 获取用户的 manager_id 作为 user_id
		userID := client.ManagerID
		if userID == 0 {
			userID = 1 // 默认用户ID
		}

		// 检查是否已存在今天的提醒
		// TODO: ExistsByClientAndType 需要在 ReminderRepo 中实现，暂时跳过
		/*
			exists, err := s.reminderRepo.ExistsByClientAndType(ctx, client.ID, "FollowUp",
				getTodayStart(), getTodayEnd())
			if err != nil {
				log.Errorf("检查提醒是否存在失败: %v", err)
				continue
			}
			if exists {
				continue
			}
		*/

		daysSinceUpdate := int(now.Sub(client.UpdatedAt).Hours() / 24)
		// priority := int8(2) // Medium
		if daysSinceUpdate > 14 {
			// priority = int8(3) // High
		}

		task := &biz_omiai.ReminderTask{
			ClientID:    int64(client.ID),
			RuleID:      0, // 系统规则
			Content:     fmt.Sprintf("该客户已%d天未联系，建议回访维护关系", daysSinceUpdate),
			ScheduledAt: now,
			Status:      "pending",
		}

		if err := s.reminderRepo.CreateTask(task); err != nil {
			// log.Errorf("创建回访提醒失败: %v", err)
		}
	}

	return nil
}

// generateBirthdayReminders 生成生日提醒（未来3天）
func (s *ReminderService) generateBirthdayReminders(ctx context.Context) error {
	now := time.Now()

	// 查询所有有生日信息的客户
	var clients []*biz_omiai.Client
	err := s.db.DB.WithContext(ctx).Model(&biz_omiai.Client{}).
		Where("birthday IS NOT NULL AND birthday != ''").
		Find(&clients).Error
	if err != nil {
		return err
	}

	for _, client := range clients {
		if len(client.Birthday) < 10 {
			continue
		}

		// 解析生日 MM-DD
		birthMonth := client.Birthday[5:7]
		birthDay := client.Birthday[8:10]
		currentYear := now.Year()

		birthdayThisYear, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%s-%s", currentYear, birthMonth, birthDay))
		if err != nil {
			continue
		}

		// 如果今年的生日已过，看明年的
		if birthdayThisYear.Before(now) {
			birthdayThisYear = birthdayThisYear.AddDate(1, 0, 0)
		}

		daysUntil := int(birthdayThisYear.Sub(now).Hours() / 24)
		if daysUntil > 3 {
			continue
		}

		// 检查是否已存在提醒
		// TODO: 实现检查逻辑
		/*
			exists, err := s.reminderRepo.ExistsByClientAndType(ctx, client.ID, "Birthday",
				getTodayStart(), getTodayEnd().AddDate(0, 0, 3))
			if err != nil {
				continue
			}
			if exists {
				continue
			}
		*/

		// userID := client.ManagerID

		// var title string
		// if daysUntil == 0 {
		// 	title = fmt.Sprintf("🎂 %s 今天生日！", client.Name)
		// } else {
		// 	title = fmt.Sprintf("🎂 %s %d天后生日", client.Name, daysUntil)
		// }

		task := &biz_omiai.ReminderTask{
			ClientID:    int64(client.ID),
			RuleID:      0,
			Content:     fmt.Sprintf("记得给%s发送生日祝福哦，维护客户关系的好时机", client.Name),
			ScheduledAt: birthdayThisYear.AddDate(0, 0, -daysUntil), // 今天提醒
			Status:      "pending",
		}

		if err := s.reminderRepo.CreateTask(task); err != nil {
			// log.Errorf("创建生日提醒失败: %v", err)
		}
	}

	return nil
}

// generateAnniversaryReminders 生成相识纪念日提醒（每月提醒）
func (s *ReminderService) generateAnniversaryReminders(ctx context.Context) error {
	now := time.Now()

	// 查询所有活跃的匹配记录（状态不是分手）
	var matchList []*biz_omiai.MatchRecord
	err := s.db.DB.WithContext(ctx).Model(&biz_omiai.MatchRecord{}).
		Where("status != ?", 6). // MatchStatusBroken=6
		Order("match_date desc").
		Find(&matchList).Error
	if err != nil {
		return err
	}

	for _, match := range matchList {
		// 计算相识月数
		months := int(now.Sub(match.MatchDate).Hours() / 24 / 30)
		if months < 1 {
			continue
		}

		// 检查是否是这个月的纪念日（match_date 的日期 == 今天的日期）
		if match.MatchDate.Day() != now.Day() {
			continue
		}

		// 检查是否已存在今天的提醒
		// TODO: 实现检查逻辑
		/*
			exists, err := s.reminderRepo.ExistsByClientAndType(ctx, match.MaleClientID, "Anniversary",
				getTodayStart(), getTodayEnd())
			if err != nil {
				continue
			}
			if exists {
				continue
			}
		*/

		// 获取客户信息
		var maleClient, femaleClient biz_omiai.Client
		s.db.DB.First(&maleClient, match.MaleClientID)
		s.db.DB.First(&femaleClient, match.FemaleClientID)

		maleName := maleClient.Name
		femaleName := femaleClient.Name

		task := &biz_omiai.ReminderTask{
			ClientID:    int64(match.MaleClientID), // 关联到男方
			RuleID:      0,
			Content:     fmt.Sprintf("💕 %s & %s 相识%d个月纪念日，建议跟进了解进展", maleName, femaleName, months),
			ScheduledAt: now,
			Status:      "pending",
		}

		if err := s.reminderRepo.CreateTask(task); err != nil {
			// log.Errorf("创建纪念日提醒失败: %v", err)
		}
	}

	return nil
}

// generateChurnRiskReminders 生成流失预警（30天未联系）
func (s *ReminderService) generateChurnRiskReminders(ctx context.Context) error {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	var clients []*biz_omiai.Client
	err := s.db.DB.WithContext(ctx).Model(&biz_omiai.Client{}).
		Where("updated_at <= ?", thirtyDaysAgo).
		Find(&clients).Error
	if err != nil {
		return err
	}

	for _, client := range clients {
		// 检查是否已存在今天的提醒
		// TODO: 实现检查逻辑
		/*
			exists, err := s.reminderRepo.ExistsByClientAndType(ctx, client.ID, "ChurnRisk",
				getTodayStart(), getTodayEnd())
			if err != nil {
				continue
			}
			if exists {
				continue
			}
		*/

		// userID := client.ManagerID
		daysSinceUpdate := int(now.Sub(client.UpdatedAt).Hours() / 24)

		task := &biz_omiai.ReminderTask{
			ClientID:    int64(client.ID),
			RuleID:      0,
			Content:     fmt.Sprintf("该客户已%d天未联系，存在流失风险，请尽快回访", daysSinceUpdate),
			ScheduledAt: now,
			Status:      "pending",
		}

		if err := s.reminderRepo.CreateTask(task); err != nil {
			// log.Errorf("创建流失预警失败: %v", err)
		}
	}

	return nil
}

// 辅助函数
func getTodayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func getTodayEnd() time.Time {
	return getTodayStart().Add(24 * time.Hour)
}
