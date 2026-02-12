# 登录跳转问题调试指南

## 问题现象
登录成功后无法跳转到首页，停留在登录页面。

## 调试步骤

### 第一步：打开浏览器开发者工具

1. 按 `F12` 打开开发者工具
2. 切换到 **Console** 标签页
3. 切换到 **Network** 标签页
4. 勾选 **Preserve log**（保留日志）

### 第二步：执行登录操作

1. 输入手机号: `13800138000`
2. 输入密码: `123456`
3. 点击登录按钮
4. 观察以下内容：

#### 2.1 检查 Network 标签页

**应该看到以下 3 个请求：**

1. `POST /api/auth/login/h5` - 登录接口
   - Status Code: 200
   - Response 应该包含 `accessToken`

2. `GET /api/user/info` - 获取用户信息
   - Status Code: 200
   - Response 应该包含用户信息

3. `GET /api/auth/codes` - 获取权限码
   - Status Code: 200
   - Response 应该包含权限码数组

**如果看到红色（失败）的请求：**
- 点击失败的请求
- 查看 **Response** 标签页
- 记录错误信息

#### 2.2 检查 Console 标签页

应该看到以下日志：
```
登录成功，获取到 accessToken: eyJhbGciOiJIUz...
获取用户信息成功: {userId: '1', username: '管理员', ...}
获取权限码成功: ["*"]
```

**如果看到错误信息：**
- 记录完整的错误堆栈
- 特别关注 "获取用户信息或权限码失败" 相关的错误

### 第三步：检查 Local Storage

1. 切换到 **Application** 标签页
2. 左侧菜单选择 **Local Storage**
3. 选择当前域名（如 `http://localhost:5666`）
4. 检查是否有 `access-token` 键
5. 值应该是一个 JWT token 格式（长字符串）

**如果没有 `access-token`：**
- 说明登录接口返回的 token 没有被正确存储
- 检查登录接口的 Response 格式

### 第四步：手动测试后端接口

#### 4.1 测试登录接口

打开浏览器开发者工具的 **Console** 标签页，执行：

```javascript
fetch('http://localhost:10131/api/auth/login/h5', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    phone: '13800138000',
    password: '123456'
  })
})
.then(res => res.json())
.then(data => console.log('登录响应:', data))
.catch(err => console.error('登录错误:', err))
```

**期望输出：**
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "nickname": "管理员",
      "role": "admin"
    }
  }
}
```

**检查点：**
- ✅ `code` 为 0
- ✅ `data.accessToken` 存在
- ✅ `data.user` 存在

**如果字段名不是 `accessToken`：**
- 说明后端返回格式不正确
- 需要修改后端代码

#### 4.2 测试用户信息接口

```javascript
fetch('http://localhost:10131/api/user/info', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer ' + '你的accessToken'
  }
})
.then(res => res.json())
.then(data => console.log('用户信息响应:', data))
.catch(err => console.error('用户信息错误:', err))
```

**期望输出：**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "nickname": "管理员",
    "avatar": "",
    "role": "admin",
    "created_at": "2026-02-11T...",
    "updated_at": "2026-02-11T..."
  }
}
```

**检查点：**
- ✅ `code` 为 0
- ✅ `data.id` 存在
- ✅ `data.nickname` 存在
- ✅ `data.role` 存在

#### 4.3 测试权限码接口

```javascript
fetch('http://localhost:10131/api/auth/codes', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer ' + '你的accessToken'
  }
})
.then(res => res.json())
.then(data => console.log('权限码响应:', data))
.catch(err => console.error('权限码错误:', err))
```

**期望输出：**
```json
{
  "code": 0,
  "message": "ok",
  "data": ["*"]
}
```

**检查点：**
- ✅ `code` 为 0
- ✅ `data` 是数组
- ✅ 数组包含权限码（如 `["*"]` 或 `["client:view", ...]`）

## 常见问题排查

### 问题1: 登录接口返回 404

**可能原因：**
- 后端服务未启动
- 路由配置错误

**解决方案：**
1. 检查后端服务是否启动：`ps aux | grep omiai-server`
2. 检查后端端口：`lsof -i :10131`
3. 重启后端服务

### 问题2: 登录接口返回 "手机号或密码错误"

**可能原因：**
- 数据库中没有该用户
- 密码未使用 MD5 加密
- 密码格式不正确

**解决方案：**
1. 连接数据库检查用户是否存在：
```sql
SELECT * FROM user WHERE phone = '13800138000';
```

2. 如果不存在，创建用户：
```sql
INSERT INTO user (nickname, phone, password, role, created_at, updated_at)
VALUES ('管理员', '13800138000', MD5('123456'), 'admin', NOW(), NOW());
```

3. 验证密码：
```sql
SELECT phone, password, MD5('123456') as encrypted
FROM user 
WHERE phone = '13800138000';
```
两个值应该一致。

### 问题3: 获取用户信息接口返回 401 Unauthorized

**可能原因：**
- Token 未正确传递
- Token 格式错误
- Token 已过期

**解决方案：**
1. 检查 Request Headers：
   - 应该包含 `Authorization: Bearer eyJhbGciOiJIUzI1NiIs...`
   - 注意 `Bearer` 后面有个空格

2. 检查 Token 是否过期：
   - JWT Token 默认 24 小时过期
   - 如果过期，需要重新登录

3. 检查后端中间件：
   - 确认 Authorization 中间件正确配置
   - 检查 Token 解析逻辑

### 问题4: 获取用户信息接口返回 "用户不存在"

**可能原因：**
- Token 中的 user_id 不正确
- 用户已被删除

**解决方案：**
1. 解析 Token 查看 user_id：
```javascript
const token = '你的accessToken';
const payload = JSON.parse(atob(token.split('.')[1]));
console.log('Token payload:', payload);
```

2. 检查数据库中是否存在该 user_id：
```sql
SELECT * FROM user WHERE id = '你的user_id';
```

### 问题5: 获取权限码接口返回错误

**可能原因：**
- 后端未实现该接口
- 路由未注册

**解决方案：**
1. 检查路由注册：
```bash
grep -n "auth/codes" internal/server/router.go
```

2. 确认路由在认证组内：
```go
authGroup := g.Group("", middleware.Authorization(r.DB, r.Redis))
{
    authGroup.GET("/auth/codes", r.AuthController.GetAccessCodes)
}
```

### 问题6: 前端报错 "Cannot read property 'homePath' of undefined"

**可能原因：**
- 用户信息转换失败
- `userInfo` 为 null

**解决方案：**
1. 检查 `getUserInfoApi()` 返回值
2. 确认转换后的 `userInfo` 包含 `homePath` 字段

## 完整的测试流程

### 准备工作

1. 确保后端服务启动：
```bash
cd /Users/edy/apps/go/src/github.com/iwuxc/omiai-server
go run ./cmd/server
```

2. 确保前端服务启动：
```bash
cd /Users/edy/apps/go/src/github.com/iwuxc/omiai-server/omiai-admin/apps/web-antd
pnpm dev
```

3. 确保数据库中有测试用户：
```sql
-- 连接数据库
mysql -u root -p omiai

-- 检查用户
SELECT * FROM user WHERE phone = '13800138000';

-- 如果不存在，创建
INSERT INTO user (nickname, phone, password, role, created_at, updated_at)
VALUES ('管理员', '13800138000', MD5('123456'), 'admin', NOW(), NOW());
```

### 执行测试

1. 清空浏览器缓存：
   - 按 `Ctrl+Shift+Delete` (Windows/Linux)
   - 或 `Cmd+Shift+Delete` (Mac)
   - 选择 "缓存的图片和文件"
   - 点击"清除数据"

2. 打开浏览器访问：
   - `http://localhost:5666`

3. 打开开发者工具（F12）

4. 执行登录操作并观察

5. 根据观察结果，对照上述常见问题排查

## 日志查看

### 后端日志

```bash
# 查看实时日志
tail -f runtime/logs/omiai-server.log

# 查看错误日志
tail -f runtime/logs/error.log
```

### 前端日志

在浏览器开发者工具的 **Console** 标签页查看

## 代码检查清单

### 后端检查

- [ ] `H5Login()` 返回字段名是 `accessToken`（不是 `token`）
- [ ] `WxLogin()` 返回字段名是 `accessToken`
- [ ] `GetUserInfo()` 接口已实现
- [ ] `GetAccessCodes()` 接口已实现
- [ ] 路由 `/api/auth/codes` 已注册
- [ ] 路由在认证组内（需要登录）

### 前端检查

- [ ] `loginApi()` 调用 `/api/auth/login/h5`
- [ ] `getUserInfoApi()` 调用 `/api/user/info`
- [ ] `getAccessCodesApi()` 调用 `/api/auth/codes`
- [ ] `getUserInfoApi()` 返回值包含 `homePath`
- [ ] `getUserInfoApi()` 返回值包含 `realName`
- [ ] 响应拦截器配置正确（`codeField: 'code'`, `dataField: 'data'`, `successCode: 0`）

### 路由检查

- [ ] `/dashboard` 路由已注册
- [ ] `/dashboard` 路由在认证组外（不需要登录）
- [ ] `/dashboard` 组件正确导入

## 联系支持

如果以上步骤都无法解决问题，请提供以下信息：

1. **浏览器控制台的完整错误日志**
2. **Network 标签页中所有请求的 Response**
3. **后端日志的错误信息**
4. **数据库中用户表的记录**
5. **前端和后端的版本号**

## 附录：快速测试脚本

已提供 `test_login.sh` 脚本，可以直接测试后端接口：

```bash
cd /Users/edy/apps/go/src/github.com/iwuxc/omiai-server
./test_login.sh
```

该脚本会依次测试：
1. 登录接口
2. 获取用户信息接口
3. 获取权限码接口

---

**文档版本**: v1.0
**最后更新**: 2026-02-11
**维护者**: CodeBuddy Code
