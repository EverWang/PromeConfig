package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 配置
var (
	SecretKey     string
	TokenDuration time.Duration
)

// Claims JWT声明
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID string) (string, error) {
	// 创建声明
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// InitJWT 初始化JWT配置
func InitJWT(secretKey string, duration time.Duration) {
	SecretKey = secretKey
	TokenDuration = duration
}