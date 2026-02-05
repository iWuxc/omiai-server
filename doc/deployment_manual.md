# 部署手册 (Deployment Manual)

## 1. 数据库迁移 (Database Migration)

本项目使用 SQL 脚本进行数据库版本管理。V2 版本引入了 admin-direct matching 模式，移除了匹配申请流程，并增加了算法初筛缓存。

### 1.1 前置条件
- Flyway (Version 2.3.0) 已安装
- MySQL 数据库已运行

### 1.2 迁移脚本
迁移脚本位置: 
- `doc/sql/V2.3.0__omiai_couple_refactor.sql` (V2 核心架构重构)
- `doc/sql/V2.3.1__add_missing_followup_columns.sql` (V2.3.1 补丁: 修复回访记录表缺失字段)

如果使用 Flyway，请确保将脚本重命名为符合 Flyway 命名规范的格式。

### 1.3 执行迁移
运行以下命令执行迁移：

```bash
flyway -url=jdbc:mysql://<host>:<port>/<database> -user=<user> -password=<password> migrate
```

### 1.4 验证迁移
执行成功后，检查数据库表结构：
- `match_request`, `match_approval` 表应被移除 (如果执行了 drop 操作) 或者不再使用。
- `client` 表应包含 `candidate_cache_json` 和 `partner_id` 字段。
- `client` 表的 `partner_id` 应有唯一索引 (Unique Index) 以支持防重复匹配。
- `match_record` 表应包含 `match_score`, `admin_id` 等新字段。
- `follow_up_record` 表应包含 `feedback`, `satisfaction`, `attachments` 字段。

## 2. 回滚步骤 (Rollback)

如果 V2 部署失败需要回滚，请执行以下步骤：

### 2.1 数据库回滚
如果使用 Flyway Pro，可以使用 `undo` 命令。否则，请手动执行以下 SQL：

```sql
-- 1. 删除新增字段 (注意数据丢失风险)
ALTER TABLE `client` DROP COLUMN `candidate_cache_json`;
ALTER TABLE `match_record` DROP COLUMN `match_score`, DROP COLUMN `admin_id`, DROP COLUMN `remark`;

-- 2. 恢复旧表 (如果已备份)
-- source V1_backup.sql;
```

**注意**: 生产环境回滚前务必进行全量备份。

## 3. 服务部署

### 3.1 编译
```bash
go build -o omiai-server ./cmd/server
```

### 3.2 启动
```bash
./omiai-server -conf ./configs/config.yaml
```

### 3.3 定时任务验证
V2 版本新增了 `CandidatePreFilterService` 定时任务 (每日 2:00 AM 执行)。
启动后，查看日志确认任务注册成功：
`[cron] CandidatePreFilterService registered`

可以通过手动触发或修改 cron 表达式进行测试。
