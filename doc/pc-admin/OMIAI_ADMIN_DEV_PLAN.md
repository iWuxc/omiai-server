# Omiai Admin 后台开发计划

## 项目概述

基于 Vben Admin 5.x (TDesign 版本) 开发的红娘助手后台管理系统，与 omiai-server Go 后端服务配套使用。

---

## 一、后端 API 接口清单

### 1. 认证模块 (Auth)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| H5登录 | POST | /api/auth/login/h5 | 手机号+验证码登录 |
| 微信登录 | POST | /api/auth/login/wx | 微信授权登录 |
| 发送短信 | POST | /api/auth/send_sms | 发送验证码 |
| 获取用户信息 | GET | /api/user/info | 获取当前登录用户信息 |
| 修改密码 | POST | /api/user/change_password | 修改密码 |
| 退出登录 | GET | /api/logout | 退出登录 |

### 2. 客户管理模块 (Clients)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 客户列表 | GET | /api/clients/list | 分页查询，支持筛选 |
| 创建客户 | POST | /api/clients/create | 创建新客户 |
| 更新客户 | POST | /api/clients/update | 更新客户信息 |
| 删除客户 | DELETE | /api/clients/delete/:id | 删除客户 |
| 客户详情 | GET | /api/clients/detail/:id | 获取客户详情 |
| 客户统计 | GET | /api/clients/stats | 获取客户统计数据 |
| 导入分析 | POST | /api/clients/import/analyze | 分析导入文件 |
| 批量导入 | POST | /api/clients/import/batch | 批量导入客户 |
| 认领客户 | POST | /api/clients/claim | 认领公海客户 |
| 释放客户 | POST | /api/clients/release | 释放客户到公海 |
| 候选人推荐 | GET | /api/clients/:id/candidates | 获取匹配候选人 |
| 对比详情 | GET | /api/clients/:id/compare/:candidateId | 对比两个客户 |

### 3. 匹配/情侣模块 (Couples)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 情侣列表 | GET | /api/couples/list | 分页查询匹配记录 |
| 创建匹配 | POST | /api/couples/create | 创建新匹配 |
| 确认匹配 | POST | /api/couples/confirm | 直接确认匹配(V2) |
| 解除匹配 | POST | /api/couples/dissolve | 解除匹配关系 |
| 更新状态 | POST | /api/couples/update_status | 更新匹配状态 |
| 跟进列表 | GET | /api/couples/followup/list | 获取跟进记录 |
| 创建跟进 | POST | /api/couples/followup/create | 创建跟进记录 |
| 状态历史 | GET | /api/couples/status/history | 获取状态变更历史 |
| 匹配统计 | GET | /api/couples/stats | 获取匹配统计数据 |

### 4. 提醒模块 (Reminders)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 提醒列表 | GET | /api/reminders/list | 分页查询提醒 |
| 今日提醒 | GET | /api/reminders/today | 获取今日提醒 |
| 待处理 | GET | /api/reminders/pending | 获取待处理提醒 |
| 提醒统计 | GET | /api/reminders/stats | 获取提醒统计 |
| 标记已读 | POST | /api/reminders/read | 标记提醒已读 |
| 标记完成 | POST | /api/reminders/done | 标记提醒已完成 |
| 删除提醒 | DELETE | /api/reminders/delete | 删除提醒 |

### 5. Banner 模块

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| Banner列表 | GET | /api/banner/list | 获取Banner列表 |
| Banner详情 | GET | /api/banner/detail | 获取Banner详情 |
| 创建Banner | POST | /api/banner/create | 创建Banner |
| 更新Banner | POST | /api/banner/update | 更新Banner |
| 删除Banner | DELETE | /api/banner/delete/:id | 删除Banner |

### 6. 仪表盘模块 (Dashboard)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 统计数据 | GET | /api/dashboard/stats | 获取Dashboard统计数据 |
| 待办事项 | GET | /api/dashboard/todos | 获取待办事项列表 |
| 客户趋势 | GET | /api/dashboard/chart/client?days=30 | 获取客户增长趋势 |
| 撮合趋势 | GET | /api/dashboard/chart/match?days=30 | 获取撮合增长趋势 |

### 7. AI 模块

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 匹配分析 | POST | /api/ai/analyze | AI分析匹配度 |
| 破冰话题 | POST | /api/ai/ice-breaker | 获取破冰话题 |

### 8. 通用模块

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 文件上传 | POST | /api/common/upload | 上传文件到COS |

---

## 二、前端开发功能点

### Phase 1: 基础架构搭建 (Day 1-2)

#### 2.1 项目初始化
- [ ] 确认使用 `web-tdesign` 作为基础版本
- [ ] 配置 API 基础路径和代理
- [ ] 配置路由模式 (history)
- [ ] 配置状态管理 (Pinia)

#### 2.2 API 层封装
- [ ] 创建 `api/omiai/` 目录结构
- [ ] 封装通用请求拦截器 (添加token)
- [ ] 封装响应拦截器 (统一错误处理)
- [ ] 按模块封装 API:
  - `api/omiai/auth.ts` - 认证相关
  - `api/omiai/client.ts` - 客户管理
  - `api/omiai/couple.ts` - 匹配/情侣
  - `api/omiai/reminder.ts` - 提醒
  - `api/omiai/banner.ts` - Banner
  - `api/omiai/dashboard.ts` - 仪表盘
  - `api/omiai/ai.ts` - AI功能
  - `api/omiai/common.ts` - 通用

#### 2.3 类型定义
- [ ] 创建 `types/omiai.ts` 定义所有业务类型
  - Client 客户类型
  - MatchRecord 匹配记录类型
  - Candidate 候选人类型
  - Reminder 提醒类型
  - Banner 类型
  - 枚举类型定义

#### 2.4 路由和菜单配置
- [ ] 创建 `router/routes/modules/omiai.ts`
- [ ] 配置菜单结构:
  ```
  仪表盘
  客户管理
    ├── 客户列表
    ├── 客户详情
    └── 导入客户
  匹配管理
    ├── 情侣档案
    ├── 匹配跟进
    └── 状态历史
  提醒中心
    ├── 待办提醒
    └── 历史记录
  Banner管理
  ```

---

### Phase 2: 仪表盘模块 (Day 3-4)

#### 2.5 页面功能
- [ ] **数据概览卡片**
  - 客户总数
  - 今日新增客户
  - 本月新增客户
  - 总撮合数
  - 本月新增撮合
  - 待跟进数

- [ ] **图表组件**
  - 客户增长趋势图 (近30天)
  - 撮合增长趋势图 (近30天)
  - 匹配状态分布饼图

- [ ] **待办事项列表**
  - 今日需跟进客户
  - 客户生日提醒
  - 匹配状态变更提醒
  - 支持快速跳转

- [ ] **快捷操作区**
  - 快速添加客户按钮
  - 查看今日提醒按钮
  - 进入客户列表按钮

---

### Phase 3: 客户管理模块 (Day 5-10)

#### 2.6 客户列表页
- [ ] **高级筛选器**
  - 基础筛选: 性别、年龄段、身高范围、学历
  - 高级筛选: 婚姻状况、房产情况、车辆情况、工作城市
  - 状态筛选: 单身、匹配中、已匹配、停止服务
  - 时间筛选: 录入时间、更新时间
  - 关键词搜索: 姓名、手机号

- [ ] **列表展示**
  - 表格展示: 头像、姓名、性别、年龄、身高、学历、状态
  - 分页功能
  - 排序功能
  - 自定义列显示

- [ ] **操作按钮**
  - 查看详情
  - 编辑客户
  - 删除客户 (确认弹窗)
  - 为客户匹配 (跳转到匹配页)
  - 认领/释放 (公海客户)

- [ ] **批量操作**
  - 批量导入按钮
  - 导出数据按钮

#### 2.7 客户详情页
- [ ] **基础信息卡片**
  - 头像、姓名、性别、年龄
  - 联系电话、微信号
  - 当前状态标签

- [ ] **详细资料**
  - 基础资料: 生日、属相、身高、体重
  - 教育职业: 学历、职业、月收入、工作城市
  - 家庭情况: 房产、车辆、家庭住址、家庭成员
  - 择偶要求: 对另一半的要求
  - 红娘备注: 可编辑的备注信息
  - 照片墙: 客户上传的照片

- [ ] **智能匹配区域**
  - 一键匹配按钮
  - 候选人列表 (姓名、匹配度、基本信息)
  - 对比按钮 (跳转到对比页)
  - 确认匹配按钮

- [ ] **操作记录**
  - 客户资料变更历史
  - 匹配记录
  - 跟进记录

#### 2.8 客户编辑/创建页
- [ ] **分步表单**
  - 步骤1: 基础信息 (姓名、性别、电话、生日、身高、体重)
  - 步骤2: 教育职业 (学历、职业、收入、工作城市)
  - 步骤3: 家庭情况 (房产、车辆、住址、家庭成员)
  - 步骤4: 择偶要求 (年龄范围、身高范围、学历要求等)
  - 步骤5: 照片上传

- [ ] **表单验证**
  - 手机号格式验证
  - 必填项验证
  - 年龄范围验证

- [ ] **照片管理**
  - 多图上传组件
  - 拖拽排序
  - 删除确认
  - 设置头像

#### 2.9 客户导入功能
- [ ] **导入向导**
  - 步骤1: 下载模板
  - 步骤2: 上传文件 (支持 Excel/CSV)
  - 步骤3: 数据预览和校验
  - 步骤4: 确认导入

- [ ] **数据校验**
  - 手机号重复检测
  - 必填字段校验
  - 数据格式校验
  - 错误行高亮显示

---

### Phase 4: 匹配管理模块 (Day 11-16)

#### 2.10 情侣档案列表
- [ ] **筛选功能**
  - 按状态筛选: 相识、交往、稳定、订婚、结婚、分手
  - 按时间筛选: 匹配时间范围
  - 关键词搜索: 客户姓名

- [ ] **列表展示**
  - 男女双方信息 (头像、姓名、年龄)
  - 匹配日期
  - 当前状态 (带颜色标签)
  - 匹配得分
  - 操作按钮

- [ ] **操作功能**
  - 查看详情
  - 更新状态
  - 添加跟进
  - 解除匹配

#### 2.11 情侣详情页
- [ ] **双方信息对比**
  - 并排展示男女双方资料
  - 关键指标对比 (年龄、身高、学历、收入)
  - 匹配度评分

- [ ] **状态时间轴**
  - 展示状态变更历史
  - 每个状态的持续时长
  - 操作人和备注

- [ ] **跟进记录**
  - 跟进记录列表
  - 添加跟进按钮
  - 跟进详情 (时间、方式、内容、反馈)

- [ ] **操作区**
  - 状态更新按钮组
  - 添加跟进按钮
  - AI分析按钮

#### 2.12 匹配对比页
- [ ] **雷达图对比**
  - 性格维度对比 (开放性、责任心、外向性等)
  - 兴趣爱好重叠度
  - 价值观匹配度

- [ ] **详细对比**
  - 基础信息对比表
  - 差异高亮显示
  - AI分析建议

- [ ] **确认匹配**
  - 确认按钮
  - 备注输入
  - 成功提示

#### 2.13 跟进管理
- [ ] **跟进记录列表**
  - 筛选: 按时间、按匹配对
  - 列表展示: 时间、方式、内容、满意度
  - 详情查看

- [ ] **添加跟进弹窗**
  - 跟进时间选择
  - 跟进方式选择 (电话/面谈/线上)
  - 跟进内容文本域
  - 客户反馈文本域
  - 满意度评分
  - 下次跟进时间
  - 附件上传

---

### Phase 5: 提醒中心模块 (Day 17-18)

#### 2.14 提醒列表
- [ ] **筛选标签**
  - 全部提醒
  - 今日提醒
  - 待处理
  - 已完成

- [ ] **列表展示**
  - 提醒类型图标 (回访/生日/纪念日/预警)
  - 优先级标识 (高/中/低)
  - 提醒标题和内容
  - 提醒时间
  - 关联客户

- [ ] **操作功能**
  - 标记已读
  - 标记完成
  - 删除提醒
  - 跳转到客户详情

#### 2.15 提醒统计
- [ ] **数据卡片**
  - 总提醒数
  - 待处理数
  - 今日提醒数
  - 高优先级数

---

### Phase 6: Banner 管理模块 (Day 19)

#### 2.16 Banner 列表
- [ ] **列表展示**
  - 缩略图预览
  - 标题
  - 排序号
  - 状态 (启用/禁用)
  - 操作按钮

- [ ] **操作功能**
  - 新增Banner
  - 编辑Banner
  - 删除Banner
  - 调整排序

#### 2.17 Banner 编辑
- [ ] **表单字段**
  - 标题输入
  - 图片上传
  - 链接URL
  - 排序号
  - 状态开关

---

### Phase 7: 系统优化 (Day 20-21)

#### 2.18 性能优化
- [ ] 列表页虚拟滚动 (客户量大时)
- [ ] 图片懒加载
- [ ] 接口缓存策略

#### 2.19 用户体验
- [ ] 操作成功提示
- [ ] 加载状态优化
- [ ] 空状态页面
- [ ] 错误页面处理

#### 2.20 权限控制
- [ ] 基于角色的菜单显示控制
- [ ] 按钮级权限控制
- [ ] 数据权限 (只能看自己管理的客户)

---

## 三、数据结构定义

### 3.1 枚举常量

```typescript
// 客户状态
enum ClientStatus {
  Single = 1,    // 单身
  Matching = 2,  // 匹配中
  Matched = 3,   // 已匹配
  Stopped = 4    // 停止服务
}

// 性别
enum Gender {
  Male = 1,      // 男
  Female = 2     // 女
}

// 婚姻状况
enum MaritalStatus {
  Unmarried = 1, // 未婚
  Married = 2,   // 已婚
  Divorced = 3,  // 离异
  Widowed = 4    // 丧偶
}

// 学历
enum Education {
  HighSchool = 1,    // 高中及以下
  JuniorCollege = 2, // 大专
  Bachelor = 3,      // 本科
  Master = 4,        // 硕士
  Doctor = 5         // 博士
}

// 匹配状态
enum MatchStatus {
  Acquaintance = 1,  // 相识
  Dating = 2,        // 交往
  Stable = 3,        // 稳定
  Engagement = 4,    // 订婚
  Married = 5,       // 结婚
  Broken = 6         // 分手
}

// 提醒类型
enum ReminderType {
  FollowUp = 1,      // 回访提醒
  Birthday = 2,      // 生日提醒
  Anniversary = 3,   // 纪念日提醒
  ChurnRisk = 4      // 流失预警
}

// 提醒优先级
enum ReminderPriority {
  Low = 1,           // 低
  Medium = 2,        // 中
  High = 3           // 高
}
```

### 3.2 主要数据类型

```typescript
// 客户
interface Client {
  id: number;
  name: string;
  gender: Gender;
  phone: string;
  birthday: string;
  avatar: string;
  age: number;
  zodiac: string;
  height: number;
  weight: number;
  education: Education;
  marital_status: MaritalStatus;
  address: string;
  family_description: string;
  income: number;
  profession: string;
  work_city: string;
  house_status: number;
  house_address: string;
  car_status: number;
  status: ClientStatus;
  partner_id?: number;
  partner?: Client;
  manager_id: number;
  is_public: boolean;
  tags: string[];
  partner_requirements: PartnerRequirements;
  parents_profession: string;
  remark: string;
  photos: string[];
  created_at: string;
  updated_at: string;
}

// 择偶要求
interface PartnerRequirements {
  age_min?: number;
  age_max?: number;
  height_min?: number;
  height_max?: number;
  education?: Education[];
  marital_status?: MaritalStatus[];
  income_min?: number;
  work_city?: string[];
  other_requirements?: string;
}

// 匹配记录
interface MatchRecord {
  id: number;
  male_client_id: number;
  female_client_id: number;
  male_client?: Client;
  female_client?: Client;
  match_date: string;
  match_score: number;
  status: MatchStatus;
  remark: string;
  admin_id: string;
  created_at: string;
  updated_at: string;
}

// 候选人
interface Candidate {
  candidate_id: number;
  name: string;
  avatar: string;
  match_score: number;
  tags: string[];
  age: number;
  height: number;
  education: number;
}

// 提醒
interface Reminder {
  id: number;
  user_id: number;
  type: ReminderType;
  client_id?: number;
  match_record_id?: number;
  client?: Client;
  match_record?: MatchRecord;
  title: string;
  content: string;
  remind_at: string;
  is_read: number;
  is_done: number;
  priority: ReminderPriority;
  created_at: string;
  updated_at: string;
}

// Banner
interface Banner {
  id: number;
  title: string;
  image_url: string;
  sort_order: number;
  status: number;
  link_url: string;
  created_at: string;
  updated_at: string;
}

// 跟进记录
interface FollowUpRecord {
  id: number;
  match_record_id: number;
  follow_up_date: string;
  method: string;
  content: string;
  feedback: string;
  satisfaction: number;
  attachments: string[];
  next_follow_up_at?: string;
  created_at: string;
  updated_at: string;
}
```

---

## 四、组件设计

### 4.1 通用组件

```
components/
├── ClientCard/          # 客户卡片
├── ClientSelector/      # 客户选择器
├── MatchScore/          # 匹配度显示
├── StatusTag/           # 状态标签
├── PhotoUpload/         # 照片上传
├── PhotoGallery/        # 照片画廊
├── ImportDialog/        # 导入弹窗
├── ComparePanel/        # 对比面板
├── RadarChart/          # 雷达图
├── TrendChart/          # 趋势图
└── StatCard/            # 统计卡片
```

### 4.2 页面组件

```
views/
├── dashboard/           # 仪表盘
│   ├── index.vue
│   ├── components/
│   │   ├── StatCards.vue
│   │   ├── TrendCharts.vue
│   │   └── TodoList.vue
├── client/              # 客户管理
│   ├── list.vue
│   ├── detail.vue
│   ├── edit.vue
│   └── import.vue
├── couple/              # 匹配管理
│   ├── list.vue
│   ├── detail.vue
│   ├── compare.vue
│   └── followup.vue
├── reminder/            # 提醒中心
│   ├── list.vue
│   └── stats.vue
└── banner/              # Banner管理
    ├── list.vue
    └── edit.vue
```

---

## 五、开发计划表

| 阶段 | 功能模块 | 预计工期 | 优先级 |
|------|----------|----------|--------|
| Phase 1 | 基础架构搭建 | 2天 | P0 |
| Phase 2 | 仪表盘模块 | 2天 | P0 |
| Phase 3 | 客户管理模块 | 6天 | P0 |
| Phase 4 | 匹配管理模块 | 6天 | P0 |
| Phase 5 | 提醒中心模块 | 2天 | P1 |
| Phase 6 | Banner管理模块 | 1天 | P2 |
| Phase 7 | 系统优化 | 2天 | P1 |
| **总计** | | **21天** | |

---

## 六、技术要点

### 6.1 状态管理
- 使用 Pinia 管理全局状态
- 按模块划分 store:
  - `stores/client.ts` - 客户相关状态
  - `stores/couple.ts` - 匹配相关状态
  - `stores/reminder.ts` - 提醒相关状态

### 6.2 路由设计
- 使用动态路由加载
- 路由懒加载优化
- 路由守卫处理权限

### 6.3 样式规范
- 使用 Tailwind CSS
- 遵循 TDesign 设计规范
- 统一色彩、间距、圆角

### 6.4 图表库
- 使用 ECharts 或 AntV G2Plot
- 封装通用图表组件

### 6.5 文件上传
- 集成腾讯云 COS SDK
- 支持图片压缩和预览
- 多图上传支持

---

## 七、后续扩展计划

### 7.1 功能扩展
- [ ] 消息通知中心 (WebSocket)
- [ ] 数据分析报表
- [ ] 客户标签管理
- [ ] 公海池管理
- [ ] 红娘绩效统计

### 7.2 性能优化
- [ ] 大数据量表格优化
- [ ] 图片CDN加速
- [ ] 接口分页优化

### 7.3 移动端适配
- [ ] 响应式布局优化
- [ ] 移动端专属页面
- [ ] 小程序管理端

---

## 八、开发环境配置

### 8.1 启动命令
```bash
# 安装依赖
pnpm install

# 启动开发服务器 (TDesign版本)
pnpm dev:tdesign

# 构建生产环境
pnpm build:tdesign
```

### 8.2 代理配置
在 `vite.config.ts` 中配置代理:
```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:10131',
      changeOrigin: true,
    },
  },
}
```

### 8.3 环境变量
创建 `.env.development`:
```
VITE_API_BASE_URL=/api
VITE_APP_TITLE=红娘助手
```

---

**开发团队**: Omiai Team  
**文档版本**: v1.0  
**最后更新**: 2026-02-09
