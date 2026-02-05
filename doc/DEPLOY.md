# 部署手册 (v2.3.0)

本版本对情侣管理系统进行了重大重构，移除了“匹配申请”流程，改为管理员直接匹配模式。

## 1. 数据库迁移

请使用 Flyway 执行以下 SQL 脚本：

- 脚本路径: `doc/sql/V2.3.0__omiai_couple_refactor.sql`
- 版本: `2.3.0`
- 描述: Couple Refactor (Remove Requests, Add Direct Match)

### 手动执行步骤 (如不使用 Flyway)

1. 备份数据库 `omiai_db`。
2. 执行脚本内容，主要变更包括：
   - 删除表: `match_request`, `match_approval`
   - 重建表: `match_record` (结构变更), `match_status_history` (新增), `follow_up_record` (新增)
   - 修改表 `client`: 
     - 新增 `candidate_cache_json` (LONGTEXT)
     - 新增/修改 `partner_id` (BIGINT, NULL, UNIQUE INDEX) - 用于防重复匹配
     - 新增 `manager_id` (BIGINT, INDEX)

## 2. 后端服务部署

1. 拉取最新代码 (branch: `feature/couple-refactor-v2`).
2. 编译: `go build -o omiai-server cmd/omiai-server/main.go`
3. 配置文件更新: 确保 `conf/config.yaml` 中包含 Redis 配置 (用于 Cron 分布式锁)。
4. 重启服务。

## 3. 定时任务

本版本新增 `CandidatePreFilterService` 定时任务。
- 频率: 每日凌晨 2:00 (自动执行)
- 功能: 预计算所有单身客户的匹配候选人并缓存至数据库。
- 验证: 部署后可观察日志 `CandidatePreFilterService start`。

## 4. 回滚方案

若需回滚至 v2.2.0：
1. 恢复数据库备份 (由于涉及删表操作，建议直接恢复备份)。
2. 回退代码至上一版本 tag。
3. 重启服务。
