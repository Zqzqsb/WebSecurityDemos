# CSRF（跨站请求伪造）攻击演示

这个项目演示了CSRF（Cross-Site Request Forgery）攻击的原理、实现和防护方法。

## 项目概述

本项目模拟了一个简单的资金转账系统，包含以下功能：
- 不安全的转账接口（易受CSRF攻击）
- 安全的转账接口（使用CSRF Token保护）
- CSRF攻击模拟
- 完整的防护示例

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

## 功能说明

### 1. 不安全的转账接口
- 路径：`/transfer/unsafe`
- 特点：没有任何CSRF保护
- 漏洞：可以被跨站请求伪造攻击利用

### 2. 安全的转账接口
- 路径：`/transfer/safe`
- 特点：使用CSRF Token保护
- 安全措施：每个请求都需要验证CSRF Token

### 3. 攻击演示
项目包含了以下攻击场景：
- 基本CSRF攻击
- 自动提交表单攻击
- 隐藏表单攻击

## 攻击演示步骤

1. 基本CSRF攻击：
   - 访问主页面
   - 点击"Simulate CSRF Attack"按钮
   - 观察未经授权的转账是否成功

2. 跨站点攻击：
   - 复制攻击代码到另一个域名下
   - 访问含有攻击代码的页面
   - 观察是否触发自动转账

## 防护措施

1. CSRF Token
```go
// 生成Token
token := generateCSRFToken()
csrfTokens.Store(userID, token)

// 验证Token
if token != expectedToken {
    return errors.New("invalid CSRF token")
}
```

2. 验证请求来源
```go
// 检查Referer头
if !isValidReferer(c.Request.Referer()) {
    return errors.New("invalid referer")
}
```

3. SameSite Cookie
```go
// 设置SameSite属性
c.SetCookie("session", token, 3600, "/", "", true, true)
```

4. 自定义请求头
```javascript
// 添加自定义头
fetch('/api/transfer', {
    headers: {
        'X-CSRF-Token': token
    }
})
```

## 最佳实践

1. 始终使用CSRF Token保护POST请求
2. 实施适当的会话管理
3. 使用SameSite Cookie属性
4. 验证请求来源
5. 避免在GET请求中修改状态

## 注意事项

1. 本项目仅用于教育目的
2. 不要在生产环境中使用不安全的示例代码
3. 在实际应用中应始终使用安全的实现方式
4. 建议在隔离的环境中进行测试

## 参考资源

- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)
- [CSRF Attacks and Defense](https://www.neuralegion.com/blog/csrf-attack/)
