# Web安全漏洞演示项目

这是一个用于学习和研究Web安全漏洞的综合性演示项目。通过实际的代码示例和详细的说明，帮助开发者理解常见的Web安全漏洞以及相应的防护措施。

## 项目结构

```
WebSecurity/
├── SQL_Inject/            # SQL注入漏洞演示
│   ├── main.go           # SQL注入示例代码
│   ├── go.mod           # Go模块依赖
│   └── README.md        # SQL注入项目说明
├── XSS_Inject/           # 跨站脚本攻击演示
│   ├── main.go          # XSS攻击示例代码
│   ├── go.mod          # Go模块依赖
│   └── README.md       # XSS攻击项目说明
├── CSRF_Attack/          # 跨站请求伪造演示
│   ├── main.go          # CSRF攻击示例代码
│   ├── go.mod          # Go模块依赖
│   └── README.md       # CSRF攻击项目说明
└── README.md            
```

## 已实现的漏洞演示

### 1. SQL注入（SQL_Inject）
- 多种SQL注入技术演示：
  - 基础认证绕过
  - 注释型注入
  - UNION查询注入
  - 布尔盲注
  - 时间延迟注入
  - 报错注入
- 安全实现对比
- [查看详情](SQL_Inject/README.md)

### 2. 跨站脚本攻击（XSS_Inject）
- 三种主要XSS攻击类型：
  - 反射型XSS
  - 存储型XSS
  - DOM型XSS
- 完整的防护示例
- [查看详情](XSS_Inject/README.md)

### 3. 跨站请求伪造（CSRF_Attack）
- 模拟银行转账系统
- CSRF攻击演示：
  - 基本CSRF攻击
  - 自动提交表单
  - 隐藏表单攻击
- Token验证防护
- [查看详情](CSRF_Attack/README.md)

## 技术栈

- 后端：Go + Gin框架
- 数据库：SQLite
- 前端：原生HTML/CSS/JavaScript
- 开发环境：Go 1.16+

## 快速开始

1. 克隆仓库：
```bash
git clone [repository-url]
cd WebSecurity
```

2. 运行SQL注入演示：
```bash
cd SQL_Inject
go mod tidy
go run main.go
# 访问 http://localhost:8080
```

3. 运行XSS攻击演示：
```bash
cd ../XSS_Inject
go mod tidy
go run main.go
# 访问 http://localhost:8080
```

4. 运行CSRF攻击演示：
```bash
cd ../CSRF_Attack
go mod tidy
go run main.go
# 访问 http://localhost:8080
```

## 学习路径建议

1. **基础知识**
   - 了解Web应用基本架构
   - 掌握HTTP协议基础
   - 学习基本的安全概念

2. **SQL注入**
   - 从基础认证绕过开始
   - 理解不同类型的注入
   - 掌握防护措施

3. **XSS攻击**
   - 学习三种XSS类型
   - 理解攻击原理
   - 实践防护方法

4. **CSRF攻击**
   - 了解CSRF原理
   - 掌握Token验证
   - 学习防护策略

## 安全建议

1. **开发环境**
   - 使用隔离的测试环境
   - 不要在生产系统测试
   - 遵循安全编码规范

2. **最佳实践**
   - 始终验证用户输入
   - 使用参数化查询
   - 实施适当的访问控制
   - 保持依赖包更新

3. **安全意识**
   - 定期安全培训
   - 关注安全公告
   - 及时修复漏洞

## 注意事项

1. 本项目仅用于教育目的
2. 不要在未授权的系统上测试
3. 遵守相关法律法规
4. 负责任地使用这些知识

## 贡献指南

欢迎通过以下方式参与项目：
1. 提交新的漏洞演示
2. 改进现有代码
3. 补充文档说明
4. 报告问题或建议

## 参考资源

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Web安全测试指南](https://owasp.org/www-project-web-security-testing-guide/)
- [CWE/SANS Top 25](https://www.sans.org/top25-software-errors/)
- [Mozilla Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)

## 版权声明

本项目采用MIT许可证。

## 免责声明

本项目中的所有内容仅用于教育和研究目的。作者不对任何人使用本项目中的内容造成的任何直接或间接损失负责。在使用本项目进行测试时，请确保您有合法的授权。
