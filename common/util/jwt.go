package util

import (
	"github.com/unknwon/com"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

// Claims JWT负载
type Claims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

// GenerateToken 生成token
func GenerateToken(id int64) (string, error) {
	claims := Claims{
		com.ToStr(id),
		jwt.StandardClaims{
			Issuer:   "gin-gin-gorm-demo",
			IssuedAt: time.Now().Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
