# JWT 认证使用指南

本文档介绍如何在 Plik 中使用 JWT 认证进行文件上传。

## 概述

Plik 现在支持使用 JWT (JSON Web Token) 进行认证，允许客户端直接使用 JWT 令牌上传文件，无需先通过 Web 界面登录。

## 认证优先级

Plik 的认证中间件按以下优先级处理认证：

1. **JWT 认证** - `Authorization: Bearer <JWT_TOKEN>`
2. **令牌认证** - `X-PlikToken: <TOKEN>`
3. **会话 Cookie 认证** - `Cookie: plik-session=<JWT_SESSION>`

## 生成 JWT 令牌

### 使用 Plik 的签名密钥

JWT 令牌必须使用与 Plik 服务器相同的签名密钥进行签名。签名密钥存储在数据库的 `settings` 表中，键为 `authentication_signature_key`。

### JWT 载荷格式

```json
{
  "uid": "user-id-here",
  "exp": 1234567890,
  "iat": 1234567890
}
```

**必需字段：**
- `uid`: 用户 ID（字符串）
- `exp`: 过期时间（Unix 时间戳）
- `iat`: 签发时间（Unix 时间戳）

### 示例代码

#### Go 语言

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/dgrijalva/jwt-go"
)

func generateJWTToken(userID, signatureKey string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS512)
    claims := token.Claims.(jwt.MapClaims)
    
    claims["uid"] = userID
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 24小时过期
    claims["iat"] = time.Now().Unix()
    
    return token.SignedString([]byte(signatureKey))
}

func main() {
    userID := "your-user-id"
    signatureKey := "your-signature-key" // 从 Plik 数据库获取
    
    token, err := generateJWTToken(userID, signatureKey)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("JWT Token: %s\n", token)
}
```

#### Python 语言

```python
import jwt
import time

def generate_jwt_token(user_id, signature_key):
    payload = {
        'uid': user_id,
        'exp': int(time.time()) + 86400,  # 24小时过期
        'iat': int(time.time())
    }
    
    token = jwt.encode(payload, signature_key, algorithm='HS512')
    return token

# 使用示例
user_id = "your-user-id"
signature_key = "your-signature-key"  # 从 Plik 数据库获取

token = generate_jwt_token(user_id, signature_key)
print(f"JWT Token: {token}")
```

#### JavaScript/Node.js

```javascript
const jwt = require('jsonwebtoken');

function generateJWTToken(userId, signatureKey) {
    const payload = {
        uid: userId,
        exp: Math.floor(Date.now() / 1000) + 86400, // 24小时过期
        iat: Math.floor(Date.now() / 1000)
    };
    
    return jwt.sign(payload, signatureKey, { algorithm: 'HS512' });
}

// 使用示例
const userId = "your-user-id";
const signatureKey = "your-signature-key"; // 从 Plik 数据库获取

const token = generateJWTToken(userId, signatureKey);
console.log(`JWT Token: ${token}`);
```

## 使用 JWT 认证上传文件

### 1. 直接上传（快速模式）

```bash
curl --form 'file=@test.txt' \
     --header 'Authorization: Bearer YOUR_JWT_TOKEN' \
     http://127.0.0.1:8080/
```

### 2. 创建上传后添加文件

```bash
# 1. 创建上传
curl -X POST http://127.0.0.1:8080/upload \
     --header 'Authorization: Bearer YOUR_JWT_TOKEN' \
     --header 'Content-Type: application/json' \
     -d '{"ttl": 3600}'

# 2. 添加文件到上传
curl --form 'file=@test.txt' \
     --header 'Authorization: Bearer YOUR_JWT_TOKEN' \
     http://127.0.0.1:8080/file/UPLOAD_ID
```

### 3. 使用 plik 客户端

```bash
# 设置 JWT 令牌环境变量
export PLIK_JWT_TOKEN="YOUR_JWT_TOKEN"

# 上传文件
plik /path/to/file
```

## 配置要求

### 1. 启用认证

在 `server/plikd.cfg` 中：

```toml
FeatureAuthentication = "enabled"  # 或 "forced"
```

### 2. 获取签名密钥

**重要**：Plik 的签名密钥存储在数据库的 `settings` 表中，而不是通过环境变量配置。可以通过以下方式获取：

#### 方法 1：通过数据库查询

```sql
SELECT value FROM settings WHERE key = 'authentication_signature_key';
```

#### 方法 2：通过 Plik 管理命令

```bash
# 查看当前配置（需要先启动服务器）
./plikd config
```

#### 方法 3：通过 API 端点

```bash
# 启动服务器后通过 API 获取
curl http://127.0.0.1:8080/config | jq '.authenticationSignatureKey'
```

**注意**：如果数据库中不存在签名密钥，Plik 会在首次启动时自动生成一个随机密钥并存储到数据库中。

## 安全注意事项

1. **密钥安全**：
   - 签名密钥必须保密
   - 定期轮换签名密钥
   - 使用强随机密钥

2. **令牌安全**：
   - 设置合理的过期时间
   - 使用 HTTPS 传输
   - 不要在日志中记录令牌

3. **用户验证**：
   - 确保用户 ID 存在且有效
   - 验证用户权限

## 错误处理

### 常见错误

| HTTP 状态码 | 错误信息 | 原因 |
|-------------|----------|------|
| 403 | "invalid JWT token" | JWT 格式错误或签名无效 |
| 403 | "invalid JWT claims" | JWT 载荷格式错误 |
| 403 | "missing user ID in JWT" | JWT 中缺少 uid 字段 |
| 403 | "user not found" | 用户不存在 |
| 500 | "unable to get user" | 数据库错误 |

### 调试步骤

1. **验证 JWT 格式**：
   ```bash
   # 使用 jwt.io 在线工具验证 JWT 格式
   echo "YOUR_JWT_TOKEN" | base64 -d
   ```

2. **检查签名**：
   ```bash
   # 使用相同的签名密钥验证
   jwt verify YOUR_JWT_TOKEN YOUR_SIGNATURE_KEY
   ```

3. **检查用户**：
   ```bash
   # 验证用户是否存在
   curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
        http://127.0.0.1:8080/me
   ```

## 示例：完整的文件上传流程

```bash
#!/bin/bash

# 配置
PLIK_SERVER="http://127.0.0.1:8080"
USER_ID="your-user-id"
SIGNATURE_KEY="your-signature-key"
FILE_PATH="/path/to/file.txt"

# 1. 生成 JWT 令牌
JWT_TOKEN=$(python3 -c "
import jwt
import time

payload = {
    'uid': '$USER_ID',
    'exp': int(time.time()) + 86400,
    'iat': int(time.time())
}

token = jwt.encode(payload, '$SIGNATURE_KEY', algorithm='HS512')
print(token)
")

echo "Generated JWT Token: $JWT_TOKEN"

# 2. 上传文件
echo "Uploading file..."
RESPONSE=$(curl -s --form "file=@$FILE_PATH" \
                --header "Authorization: Bearer $JWT_TOKEN" \
                "$PLIK_SERVER/")

echo "Upload response: $RESPONSE"

# 3. 验证上传
if [[ $RESPONSE == http* ]]; then
    echo "Upload successful! File URL: $RESPONSE"
else
    echo "Upload failed: $RESPONSE"
fi
```

## 与现有认证方式的兼容性

JWT 认证与现有的认证方式完全兼容：

- **会话 Cookie**：Web 界面登录后仍使用 Cookie
- **X-PlikToken**：API 调用仍可使用令牌
- **HTTP Basic**：密码保护的上传仍使用 Basic 认证

JWT 认证只是添加了一种新的认证方式，不会影响现有功能。
