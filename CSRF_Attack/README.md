# CSRF（跨站请求伪造）攻击演示

这个项目演示了CSRF（Cross-Site Request Forgery）攻击的原理、实现和防护方法。通过一个模拟的银行转账系统，展示了CSRF攻击的危险性以及相应的防护措施。

## 功能特点

1. **转账系统模拟**
   - 用户余额管理
   - 转账交易记录
   - 实时余额更新
   - 完整的事务处理

2. **不安全的转账接口**
   - 路径：`/transfer/unsafe`
   - 无CSRF保护
   - 易受攻击的表单提交
   - 完整的错误处理

3. **安全的转账接口**
   - 路径：`/transfer/safe`
   - CSRF Token保护
   - 安全的表单提交
   - 事务完整性保证

4. **安全特性**
   - CSRF Token生成和验证
   - 数据库事务处理
   - 输入验证和过滤
   - 错误处理和日志记录

## 系统要求

- Go 1.16+
- SQLite3
- 现代Web浏览器

## 安装步骤

1. 克隆仓库：
```bash
git clone [repository-url]
cd CSRF_Attack
```

2. 安装依赖：
```bash
go mod tidy
```

3. 运行应用：
```bash
go run main.go
```

4. 访问演示页面：
```
http://localhost:8080
```

## 功能演示

### 1. 基本转账操作

1. 访问主页面
2. 在转账表单中输入：
   - 接收方用户名
   - 转账金额
3. 点击"Transfer"按钮
4. 观察转账结果和余额变化

### 2. CSRF攻击演示

1. 不安全的转账：
   - 使用不带CSRF保护的表单
   - 观察转账是否成功
   - 查看交易记录

2. 安全的转账：
   - 使用带CSRF Token的表单
   - 尝试绕过Token验证
   - 观察安全机制的效果

### 3. 攻击模拟

点击"Simulate CSRF Attack"按钮，将会：
- 创建隐藏的恶意表单
- 自动提交转账请求
- 展示攻击结果

## 安全特性说明

1. **CSRF Token保护**
   - 每个会话生成唯一Token
   - 表单提交需要验证Token
   - Token通过安全的方式传输

2. **交易安全**
   - 完整的事务处理
   - 余额检查
   - 用户验证
   - 金额验证

3. **错误处理**
   - 详细的错误信息
   - 事务回滚
   - 用户友好的提示

## API说明

### 不安全的转账接口
```
POST /transfer/unsafe
参数：
- to: 接收方用户名
- amount: 转账金额
```

### 安全的转账接口
```
POST /transfer/safe
头部：
- X-CSRF-Token: CSRF令牌
参数：
- to: 接收方用户名
- amount: 转账金额
```

## 最佳实践

1. **CSRF防护**
   - 使用CSRF Token
   - 验证请求来源
   - 使用安全的Cookie设置

2. **数据安全**
   - 使用事务处理
   - 验证所有输入
   - 适当的错误处理

3. **用户体验**
   - 清晰的错误提示
   - 实时的余额更新
   - 交易历史记录

## 注意事项

1. 本项目仅用于教育目的
2. 不要在生产环境中使用不安全的示例代码
3. 在实际应用中应始终使用安全的实现方式
4. 建议在隔离的环境中进行测试

## 故障排除

1. 数据库错误
   - 检查SQLite数据库文件权限
   - 确保数据库连接正确

2. CSRF Token错误
   - 检查Token是否正确传递
   - 确认Token未过期

3. 转账失败
   - 检查余额是否充足
   - 验证用户名是否正确
   - 确认金额格式正确

## 参考资源

- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)
- [Go Web Examples](https://gowebexamples.com/)
