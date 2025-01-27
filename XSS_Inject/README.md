# XSS（跨站脚本）攻击演示

这个项目演示了三种主要的XSS（跨站脚本）攻击类型及其防护方法。

## XSS攻击类型

### 1. 反射型XSS
- 描述：用户输入直接返回到浏览器并执行
- 位置：搜索框
- 测试payload：
```html
<script>alert('Reflected XSS!');</script>
```

### 2. 存储型XSS
- 描述：恶意代码存储在服务器上，影响所有访问页面的用户
- 位置：评论系统
- 测试payload：
```html
<script>alert('Stored XSS!');</script>
```

### 3. DOM型XSS
- 描述：通过JavaScript动态修改DOM时引入的漏洞
- 位置：问候消息和URL片段
- 测试payload：
```html
<img src=x onerror="alert('DOM XSS!');">
```

## 运行方法

1. 安装依赖：
```bash
go mod tidy
```

2. 运行服务器：
```bash
go run main.go
```

3. 访问演示页面：
```
http://localhost:8080
```

## 漏洞演示步骤

### 反射型XSS测试
1. 在搜索框中输入XSS payload
2. 提交表单
3. 观察返回页面中的脚本执行情况

### 存储型XSS测试
1. 在评论框中输入XSS payload
2. 提交评论
3. 刷新页面，观察存储的脚本执行情况

### DOM型XSS测试
1. 在姓名输入框中输入XSS payload
2. 点击"Show Greeting"按钮
3. 观察动态生成的内容中的脚本执行情况

或者：
1. 在URL中添加锚点：`#<img src=x onerror="alert('DOM XSS!');">`
2. 访问修改后的URL
3. 观察页面加载时脚本的执行情况

## 安全实现示例

项目同时提供了安全的实现方式，展示如何防止XSS攻击：
1. 使用HTML转义
2. 使用安全的模板系统
3. 输入验证和过滤

## 防护建议

1. 始终对用户输入进行HTML转义
2. 使用安全的模板引擎
3. 实施内容安全策略（CSP）
4. 设置适当的Cookie安全标志
5. 对输入进行验证和过滤
6. 使用现代框架的XSS防护功能

## 注意事项

1. 本项目仅用于教育目的
2. 不要在生产环境中使用不安全的示例代码
3. 在实际应用中应始终使用安全的实现方式
4. 建议在隔离的环境中进行测试

## 参考资源

- [OWASP XSS Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)
- [MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)
- [Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
