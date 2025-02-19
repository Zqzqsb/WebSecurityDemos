# SQL注入演示项目

使用Go和SQLite实现的SQL注入漏洞演示和防护技术教学项目。

## 项目概述

本项目提供了多种SQL注入技术的实践演示及其防范措施。通过包含脆弱性和安全性的登录系统实现，帮助理解：
- SQL注入攻击的工作原理
- 不同类型的SQL注入技术
- 如何使用参数化查询防止SQL注入

## 功能特性

### 1. SQL注入攻击类型演示

#### 基础认证绕过
- 示例：`' OR '1'='1`
- 演示如何通过简单的字符串操作绕过登录验证

#### 注释型注入
- 示例：`admin'--`
- 展示如何使用SQL注释修改查询逻辑

#### UNION查询注入
- 示例：`admin' UNION SELECT 1 as id, 'hacker' as username, 'pwned' as password, 'admin' as role --`
- 演示如何合并查询结果与注入数据

#### 布尔盲注
- 示例：`admin' AND (SELECT CASE WHEN (1=1) THEN 1 ELSE 0 END)='1`
- 展示如何通过真/假响应提取信息

#### 时间延迟注入
- 示例：`admin' AND (SELECT CASE WHEN (1=1) THEN sqlite3_sleep(2000) ELSE 1 END)='1`
- 演示基于时间延迟的数据提取技术

#### 报错注入
- 示例：`admin' AND (SELECT CASE WHEN (1=1) THEN CAST('a' AS INTEGER) ELSE 1 END)='1`
- 展示如何利用错误信息提取数据

### 2. 安全特性

- 参数化查询的演示
- 不安全与安全SQL实践的对比
- 实时查询日志记录分析
- 清晰的错误信息展示

## 技术栈

- 后端：Go（Gin框架）
- 数据库：SQLite
- ORM框架：GORM
- 前端：HTML、CSS、JavaScript

## 安装说明

1. 确保已安装Go（版本1.16或更高）
2. 克隆仓库
3. 安装依赖：
```bash
go mod download
```

## 运行演示

1. 启动服务器：
```bash
go run main.go
```

2. 访问演示页面：http://localhost:8080

## 默认测试账号

- 管理员账号：
  - 用户名：admin
  - 密码：123456

- 普通用户账号：
  - 用户名：user1
  - 密码：password1
  
  - 用户名：user2
  - 密码：password2

## 安全提示

本项目仅用于教育目的。脆弱性登录实现展示了常见的安全缺陷，切勿在生产环境中使用。实际应用中应始终使用参数化查询和适当的安全措施。

## 学习目标

1. 理解SQL注入漏洞
2. 学习不同的SQL注入技术
3. 理解参数化查询的重要性
4. 识别不安全的编码实践
5. 实现安全的认证系统

## 最佳实践演示

1. 使用参数化查询
2. 正确的错误处理
3. 输入验证
4. 安全的密码存储（注意：本演示为了简单使用明文密码 - 生产环境切勿这样做！）
5. 清晰的代码组织和文档

## 参与贡献

欢迎通过以下方式参与项目：
1. 添加新的SQL注入技术演示
2. 改进用户界面和用户体验
3. 添加更多安全演示
4. 完善文档

## 开源协议

MIT许可证 - 可自由用于教育目的。

## 免责声明

本工具仅用于教育目的。请勿在未经授权的系统上使用这些技术。在进行安全测试时，确保您有合法的授权。

## 注意事项

1. 本项目仅用于学习和研究SQL注入攻击的原理和防范方法
2. 请勿将演示的攻击技术用于非法用途
3. 在实际项目中应始终采用安全的编码实践
4. 建议在安全的测试环境中运行此演示

## 常见问题解答

1. 为什么看不到注入效果？
   - 检查输入的注入代码是否正确
   - 确保使用的是不安全的登录接口
   - 查看浏览器控制台是否有错误信息

2. 如何确认注入是否成功？
   - 观察登录响应
   - 查看服务器日志输出
   - 检查返回的用户信息

3. 为什么使用明文密码？
   - 仅为了演示方便
   - 实际应用中必须使用加密存储
   - 推荐使用bcrypt等加密算法

## 更新日志

### v1.0.0
- 初始版本发布
- 实现基本的SQL注入演示
- 添加安全登录对比
- 完善文档说明
