package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTClaims 自定义 JWT 载荷结构
type JWTClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

// generateJWTToken 生成 JWT token
func generateJWTToken(userID, signatureKey string, expirationHours int) (string, error) {
	// 创建 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, JWTClaims{
		UID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expirationHours) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "plik-jwt-generator",
		},
	})

	// 签名 token
	tokenString, err := token.SignedString([]byte(signatureKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// validateJWTToken 验证 JWT token
func validateJWTToken(tokenString, signatureKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signatureKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("JWT Token Generator for Plik")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  jwt_gen -user <user_id> -key <signature_key> [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -user string")
	fmt.Println("        User ID (required)")
	fmt.Println("  -key string")
	fmt.Println("        Signature key (required)")
	fmt.Println("  -expire int")
	fmt.Println("        Token expiration in hours (default: 24)")
	fmt.Println("  -validate string")
	fmt.Println("        Validate existing token")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Generate token for user 'user123' with 24-hour expiration")
	fmt.Println("  jwt_gen -user user123 -key your-secret-key")
	fmt.Println("")
	fmt.Println("  # Generate token with 7-day expiration")
	fmt.Println("  jwt_gen -user user123 -key your-secret-key -expire 168")
	fmt.Println("")
	fmt.Println("  # Validate existing token")
	fmt.Println("  jwt_gen -validate 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...' -key your-secret-key")
	fmt.Println("")
	fmt.Println("  # Get signature key from Plik database")
	fmt.Println("  sqlite3 plik.db \"SELECT value FROM settings WHERE key = 'authentication_signature_key';\"")
	fmt.Println("")
	fmt.Println("  # Get signature key from Plik API")
	fmt.Println("  curl -s http://127.0.0.1:8080/config | jq -r '.authenticationSignatureKey'")
}

func main() {
	var (
		userID        = flag.String("user", "", "User ID (required)")
		signatureKey  = flag.String("key", "", "Signature key (required)")
		expireHours   = flag.Int("expire", 24, "Token expiration in hours")
		validateToken = flag.String("validate", "", "Validate existing token")
		help          = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	// 验证必需参数
	if *signatureKey == "" {
		fmt.Fprintf(os.Stderr, "Error: signature key is required\n")
		fmt.Fprintf(os.Stderr, "Use -help for usage information\n")
		os.Exit(1)
	}

	// 如果提供了验证参数，则验证 token
	if *validateToken != "" {
		claims, err := validateJWTToken(*validateToken, *signatureKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Token validation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Token is valid!")
		fmt.Printf("User ID: %s\n", claims.UID)
		fmt.Printf("Issued At: %s\n", time.Unix(claims.IssuedAt, 0).Format(time.RFC3339))
		fmt.Printf("Expires At: %s\n", time.Unix(claims.ExpiresAt, 0).Format(time.RFC3339))
		fmt.Printf("Issuer: %s\n", claims.Issuer)

		// 检查是否即将过期
		timeUntilExpiry := time.Until(time.Unix(claims.ExpiresAt, 0))
		if timeUntilExpiry < time.Hour {
			fmt.Printf("⚠️  Warning: Token expires in %v\n", timeUntilExpiry.Round(time.Minute))
		}

		os.Exit(0)
	}

	// 验证用户 ID
	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: user ID is required\n")
		fmt.Fprintf(os.Stderr, "Use -help for usage information\n")
		os.Exit(1)
	}

	// 生成 JWT token
	token, err := generateJWTToken(*userID, *signatureKey, *expireHours)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	// 输出结果
	fmt.Println("🔑 Generated JWT Token:")
	fmt.Println("")
	fmt.Println(token)
	fmt.Println("")
	fmt.Println("📋 Token Information:")
	fmt.Printf("User ID: %s\n", *userID)
	fmt.Printf("Expires in: %d hours\n", *expireHours)
	fmt.Printf("Expires at: %s\n", time.Now().Add(time.Duration(*expireHours)*time.Hour).Format(time.RFC3339))
	fmt.Println("")
	fmt.Println("🚀 Usage Examples:")
	fmt.Println("")
	fmt.Println("# Upload file with JWT authentication:")
	fmt.Printf("curl --form 'file=@test.txt' \\\n")
	fmt.Printf("     --header 'Authorization: Bearer %s' \\\n", token)
	fmt.Printf("     http://127.0.0.1:8080/\n")
	fmt.Println("")
	fmt.Println("# Create upload with JWT authentication:")
	fmt.Printf("curl -X POST http://127.0.0.1:8080/upload \\\n")
	fmt.Printf("     --header 'Authorization: Bearer %s' \\\n", token)
	fmt.Printf("     --header 'Content-Type: application/json' \\\n")
	fmt.Printf("     -d '{\"ttl\": 3600}'\n")
	fmt.Println("")
	fmt.Println("# Validate this token:")
	fmt.Printf("jwt_gen -validate '%s' -key '%s'\n", token, *signatureKey)
}
