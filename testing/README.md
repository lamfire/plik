# JWT Token Generator for Plik

这是一个用于生成和验证 Plik JWT 认证令牌的命令行工具。

## 功能特性

- ✅ 生成 JWT 认证令牌
- ✅ 验证现有令牌
- ✅ 自定义过期时间
- ✅ 详细的令牌信息显示
- ✅ 使用示例和命令生成

## 编译

```bash
# 进入 testing 目录
cd testing

# 编译工具
make build

# 或者直接使用 go build
go build -o jwt_gen jwt_gen.go
```

## 使用方法

### 基本用法

```bash
# 生成 JWT 令牌
./jwt_gen -user <user_id> -key <signature_key> [options]
```

### 参数说明

| 参数 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `-user` | string | ✅ | - | 用户 ID |
| `-key` | string | ✅ | - | 签名密钥 |
| `-expire` | int | ❌ | 24 | 过期时间（小时） |
| `-validate` | string | ❌ | - | 验证现有令牌 |
| `-help` | bool | ❌ | false | 显示帮助信息 |

### 使用示例

#### 1. 生成基本令牌

```bash
# 生成 24 小时过期的令牌
./jwt_gen -user user123 -key your-secret-key
```

#### 2. 生成长期令牌

```bash
# 生成 7 天过期的令牌
./jwt_gen -user user123 -key your-secret-key -expire 168
```

#### 3. 验证令牌

```bash
# 验证现有令牌
./jwt_gen -validate "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..." -key your-secret-key
```

#### 4. 查看帮助

```bash
./jwt_gen -help
```

## 获取签名密钥

### 方法 1：从 Plik 数据库获取

```bash
# SQLite 数据库
sqlite3 plik.db "SELECT value FROM settings WHERE key = 'authentication_signature_key';"

# PostgreSQL 数据库
psql -d plik -c "SELECT value FROM settings WHERE key = 'authentication_signature_key';"

# MySQL 数据库
mysql -u username -p plik -e "SELECT value FROM settings WHERE key = 'authentication_signature_key';"
```

### 方法 2：从 Plik API 获取

```bash
# 通过 API 获取配置
curl -s http://127.0.0.1:8080/config | jq -r '.authenticationSignatureKey'
```

### 方法 3：从 Plik 管理命令获取

```bash
# 启动 Plik 服务器后查看配置
./plikd config | grep -i signature
```

## 使用生成的令牌

### 上传文件

```bash
# 使用 JWT 令牌上传文件
curl --form 'file=@test.txt' \
     --header 'Authorization: Bearer YOUR_JWT_TOKEN' \
     http://127.0.0.1:8080/
```

### 创建上传

```bash
# 使用 JWT 令牌创建上传
curl -X POST http://127.0.0.1:8080/upload \
     --header 'Authorization: Bearer YOUR_JWT_TOKEN' \
     --header 'Content-Type: application/json' \
     -d '{"ttl": 3600}'
```

### 添加文件到上传

```bash
# 使用 JWT 令牌添加文件到现有上传
curl --form 'file=@test.txt' \
     --header 'Authorization: Bearer YOUR_JWT_TOKEN' \
     http://127.0.0.1:8080/file/UPLOAD_ID
```

## 令牌格式

生成的 JWT 令牌包含以下信息：

```json
{
  "uid": "user_id",
  "exp": 1234567890,
  "iat": 1234567890,
  "iss": "plik-jwt-generator"
}
```

- `uid`: 用户 ID
- `exp`: 过期时间（Unix 时间戳）
- `iat`: 签发时间（Unix 时间戳）
- `iss`: 签发者

## 安全注意事项

1. **保护签名密钥**：签名密钥必须保密，不要泄露给未授权人员
2. **合理设置过期时间**：不要设置过长的过期时间
3. **使用 HTTPS**：在生产环境中使用 HTTPS 传输令牌
4. **定期轮换密钥**：定期更换签名密钥以提高安全性

## 故障排除

### 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| "signature key is required" | 未提供签名密钥 | 使用 `-key` 参数提供密钥 |
| "user ID is required" | 未提供用户 ID | 使用 `-user` 参数提供用户 ID |
| "Token validation failed" | 令牌无效或过期 | 检查令牌格式和签名密钥 |
| "unexpected signing method" | 签名算法不匹配 | 确保使用正确的签名密钥 |

### 调试步骤

1. **验证签名密钥**：
   ```bash
   # 检查密钥是否正确
   ./jwt_gen -validate "YOUR_TOKEN" -key "YOUR_KEY"
   ```

2. **检查令牌格式**：
   ```bash
   # 使用 jwt.io 在线工具检查令牌格式
   echo "YOUR_TOKEN" | base64 -d
   ```

3. **验证用户存在**：
   ```bash
   # 检查用户是否存在于 Plik 中
   curl -H "Authorization: Bearer YOUR_TOKEN" http://127.0.0.1:8080/me
   ```

## 开发

### 项目结构

```
testing/
├── jwt_gen.go      # 主程序文件
├── Makefile        # 构建配置
└── README.md       # 说明文档
```

### 依赖

- Go 1.16+
- github.com/dgrijalva/jwt-go

### 构建

```bash
# 安装依赖
make deps

# 构建工具
make build

# 运行测试
make test

# 清理文件
make clean
```

## 许可证

本工具遵循与 Plik 项目相同的许可证。