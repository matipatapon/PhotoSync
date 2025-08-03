package jwt

import (
	"crypto/rand"
	"errors"
	"log"
	"os"
	"photosync/src/helper"
	"strconv"

	ejwt "github.com/golang-jwt/jwt/v5"
)

var logger *log.Logger = log.New(os.Stdout, "[JwtManager]: ", log.LstdFlags)

type IJwtManager interface {
	Create(data JwtPayload) (string, error)
	Decode(tokenString string) (JwtPayload, error)
}

type JwtManager struct {
	secretKey []byte
	th        helper.ITimeHelper
}

func NewJwtManager(th helper.ITimeHelper) JwtManager {
	secretKey := make([]byte, 64)
	rand.Read(secretKey)
	jm := JwtManager{secretKey: secretKey, th: th}
	logger.Print("Created JwtManager")
	return jm
}

func (jm *JwtManager) Create(data JwtPayload) (string, error) {
	claims := ejwt.NewWithClaims(
		ejwt.SigningMethodHS256,
		ejwt.MapClaims{
			"user_id":         strconv.FormatInt(data.UserId, 10),
			"username":        data.Username,
			"expiration_time": strconv.FormatInt(data.ExpirationTime, 10),
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

	jp.ExpirationTime, _ = strconv.ParseInt(claimsMap["expiration_time"].(string), 10, 64)
	currentTime := jm.th.TimeNow()
	if jp.ExpirationTime < currentTime {
		logger.Printf("Token expired, expiration time '%d', current time '%d'", jp.ExpirationTime, currentTime)
		return jp, errors.New("token expired")
	}

	jp.Username = claimsMap["username"].(string)
	jp.UserId, _ = strconv.ParseInt(claimsMap["user_id"].(string), 10, 64)

	return jp, nil
}
