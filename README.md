# Web安全漏洞演示项目

这是一个用于学习和研究Web安全漏洞的综合性演示项目。通过实际的代码示例和详细的说明，帮助开发者理解常见的Web安全漏洞以及相应的防护措施。

## 项目结构

```
WebSecurity/
├── SQL_Inject/            # SQL注入漏洞演示
│   ├── main.go           # SQL注入示例代码
│   ├── go.mod           # Go模块依赖
│   └── README.md        # SQL注入项目说明
├── XSS/                  # 跨站脚本攻击演示（计划中）
├── CSRF/                 # 跨站请求伪造演示（计划中）
└── README.md            # 项目总体说明
```

## 已实现的漏洞演示

### 1. SQL注入（SQL_Inject）
- 多种SQL注入技术的演示
- 包含安全和不安全的实现对比
- 详细的攻击原理说明
- [查看详情](SQL_Inject/README.md)

## 计划实现的漏洞演示

### 2. 跨站脚本攻击（XSS）
- 反射型XSS
- 存储型XSS
- DOM型XSS
- XSS防护措施

### 3. 跨站请求伪造（CSRF）
- CSRF攻击原理
- Token验证
- 双重Cookie验证
- SameSite Cookie

### 4. 其他计划
- 文件上传漏洞
- 命令注入
- 目录遍历
- 会话劫持
- 点击劫持
- 服务器端请求伪造（SSRF）

