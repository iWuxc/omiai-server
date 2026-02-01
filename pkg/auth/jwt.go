package auth

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	// 尝试加载 RSA 密钥
	prvKeyPath := "configs/certs/private.pem"
	pubKeyPath := "configs/certs/public.pem"

	prvKeyBytes, err := os.ReadFile(prvKeyPath)
	if err == nil {
		signKey, _ = jwt.ParseRSAPrivateKeyFromPEM(prvKeyBytes)
	}

	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err == nil {
		verifyKey, _ = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	}
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint64, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	if signKey != nil {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		return token.SignedString(signKey)
	}

	// Fallback to HS256 if RSA keys are missing
	secret := []byte("omiai-server-secret-key-2026")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if signKey != nil {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return verifyKey, nil
		}
		return []byte("omiai-server-secret-key-2026"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
