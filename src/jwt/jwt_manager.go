package jwt

import (
	"crypto/rand"
	"errors"
	"log"
	"os"
	"time"

	ejwt "github.com/golang-jwt/jwt/v5"
)

var logger *log.Logger = log.New(os.Stdout, "[JwtManager]: ", log.LstdFlags)

type IJwtManager interface {
	Create(data JwtPayload) (string, error)
	Decode(tokenString string) (JwtPayload, error)
}

type JwtManager struct {
	secretKey []byte
}

func NewJwtManager() JwtManager {
	jm := JwtManager{}
	jm.secretKey = make([]byte, 64)
	rand.Read(jm.secretKey)
	logger.Print("Created JwtManager")
	return jm
}

func (jm *JwtManager) Create(data JwtPayload) (string, error) {
	claims := ejwt.NewWithClaims(
		ejwt.SigningMethodHS256,
		ejwt.MapClaims{
			"username":        data.Username,
			"expiration_time": data.ExpirationTime,
		})
	tokenString, err := claims.SignedString(jm.secretKey)
	if err != nil {
		logger.Printf("Failed to sign token: '%s'", err.Error())
		return "", err
	}
	logger.Print("Created token")
	return tokenString, err
}

func (jm *JwtManager) Decode(tokenString string) (JwtPayload, error) {
	token, err := ejwt.Parse(tokenString, func(token *ejwt.Token) (interface{}, error) {
		return jm.secretKey, nil
	})
	if err != nil {
		logger.Printf("Failed to parse token: '%s'", err.Error())
		return JwtPayload{}, err
	}

	claimsMap := token.Claims.(ejwt.MapClaims)
	jp := JwtPayload{}

	jp.ExpirationTime = int64(claimsMap["expiration_time"].(float64))
	if jp.ExpirationTime < time.Now().Unix() {
		logger.Print("Token expired")
		return jp, errors.New("token expired")
	}

	jp.Username = claimsMap["username"].(string)
	return jp, nil
}
