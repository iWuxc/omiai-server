# 情侣关系全生命周期管理系统 - API 接口文档 (v2)

## 1. 概述
Base URL: `/api/v1`

**变更说明**:
- v2: 移除匹配申请流程，改为管理员直接匹配确认；新增候选人推荐与对比功能；支持管理员单向解除匹配。
- v1 (已废弃): 包含匹配申请、审批流程。

## 2. 接口列表

### 2.1 匹配候选与对比 (新增)

#### 获取匹配候选列表
- **GET** `/clients/{id}/candidates`
- **Summary**: 获取指定客户的一键匹配候选列表（基于算法初筛缓存）
- **Parameters**:
  - `id` (path): 客户ID
- **Response** (200 OK):
  ```json
  {
    "code": 0,
    "message": "success",
    "data": [
      {
        "candidate_id": 1002,
        "name": "Jane Doe",
        "avatar": "http://...",
        "match_score": 95,
        "tags": ["本科", "年薪30-50万"]
      },
      ...
    ]
  }
  ```

#### 获取匹配对比详情
- **GET** `/clients/{id}/compare/{candidateId}`
- **Summary**: 获取客户与候选人的详细维度对比数据
- **Parameters**:
  - `id` (path): 客户ID
  - `candidateId` (path): 候选人ID
- **Response** (200 OK):
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "basic_info": {
        "age": { "client": 28, "candidate": 26, "diff": "合适" },
        "height": { "client": 180, "candidate": 165, "diff": "合适" },
        "education": { "client": "硕士", "candidate": "本科", "match": true },
        "job": { "client": "程序员", "candidate": "教师", "match": true },
        "income": { "client": "50w+", "candidate": "20w+", "match": true }
      },
      "personality_radar": {
        "openness": { "client": 80, "candidate": 75 },
        "conscientiousness": { "client": 60, "candidate": 85 },
        "extraversion": { "client": 70, "candidate": 65 },
        "agreeableness": { "client": 90, "candidate": 90 },
        "neuroticism": { "client": 40, "candidate": 30 }
      },
      "interests": {
        "overlap_percentage": 0.8,
        "common_list": ["旅行", "摄影", "美食"]
      },
      "values": {
        "match_percentage": 0.9,
        "details": [...]
      },
      "relationship_expectations": {
        "short_term": { "client": 2, "candidate": 1 }, // 1-5分
        "long_term": { "client": 5, "candidate": 5 }
      }
    }
  }
  ```

### 2.2 匹配确认 (修改)

#### 确认匹配 (管理员)
- **POST** `/couples/confirm`
- **Summary**: 管理员直接确认匹配，生成情侣关系。系统会自动进行去重检查（分布式锁+数据库唯一索引）。
- **Request Body**:
  ```json
  {
    "client_id": 1001,
    "candidate_id": 1002,
    "remark": "双方意向匹配"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "code": 0,
    "message": "匹配成功",
    "data": {
      "couple_id": 200,
      "status": 1, // 1: 相识
      "matched_at": "2023-10-27T10:00:00Z"
    }
  }
  ```

### 2.3 解除匹配 (新增)

#### 解除匹配关系
- **POST** `/couples/dissolve`
- **Summary**: 管理员解除客户当前的匹配关系。操作将记录在审计日志中，双方状态重置为"单身" (1)。
- **Request Body**:
  ```json
  {
    "client_id": 1001,
    "reason": "双方性格不合"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "code": 0,
    "message": "解除匹配成功",
    "data": null
  }
  ```

### 2.4 回访记录 (新增)

#### 创建回访记录
- **POST** `/couples/followup/create`
- **Summary**: 创建情侣关系的回访记录
- **Request Body**:
  ```json
  {
    "match_record_id": 200,
    "follow_up_date": "2026-02-05", // 支持 YYYY-MM-DD 或 RFC3339
    "method": "电话",
    "content": "双方相处融洽",
    "feedback": "希望能多组织线下活动",
    "satisfaction": 5, // 1-5分
    "attachments": "[\"http://.../img1.jpg\"]", // JSON string
    "next_follow_up_at": "2026-02-12"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "code": 0,
    "message": "保存成功",
    "data": { ... }
  }
  ```

#### 获取回访记录列表
- **GET** `/couples/followup/list`
- **Summary**: 获取指定匹配记录的所有回访历史
- **Parameters**:
  - `match_record_id` (query): 匹配记录ID
- **Response** (200 OK):
  ```json
  {
    "code": 0,
    "message": "ok",
    "data": [ ... ]
  }
  ```
