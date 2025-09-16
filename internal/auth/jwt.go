package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userID string) (string, error) {
	secret := []byte(viper.GetString("jwt.secret"))
	dur, err := time.ParseDuration(viper.GetString("jwt.access_ttl"))
	if err != nil {
		dur = 15 * time.Minute
	}
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := []byte(viper.GetString("jwt.secret"))
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		return nil, err
	}
	return claims, nil
}
