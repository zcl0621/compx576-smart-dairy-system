package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zcl0621/compx576-smart-dairy-system/config"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
)

var ErrBadToken = errors.New("token looks wrong")
var ErrBadJWTSecret = errors.New("jwt secret is empty")

const tokenTTL = 24 * time.Hour
const authUserIDKey = "auth_user_id"

type jwtClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(tokenTTL)
	claims := jwtClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := jwtSecret()
	if err != nil {
		return "", time.Time{}, err
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func ParseToken(tokenString string) (*jwtClaims, error) {
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrBadToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrBadToken
	}
	if !token.Valid {
		return nil, ErrBadToken
	}
	if claims.UserID == "" {
		return nil, ErrBadToken
	}

	return claims, nil
}

func GetAuthUserID(c *gin.Context) (string, error) {
	if value, ok := c.Get(authUserIDKey); ok {
		userID, _ := value.(string)
		if userID != "" {
			return userID, nil
		}
	}

	tokenString, err := parseBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		return "", err
	}

	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	c.Set(authUserIDKey, claims.UserID)
	return claims.UserID, nil
}

func NeedAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := GetAuthUserID(c); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}

func parseBearerToken(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", ErrBadToken
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return "", ErrBadToken
	}
	if parts[1] == "" {
		return "", ErrBadToken
	}

	return parts[1], nil
}

func jwtSecret() (string, error) {
	secret := config.Get().JWT.Secret
	if secret != "" {
		return secret, nil
	}

	projectlog.L().Error("jwt secret is empty")
	return "", ErrBadJWTSecret
}

const cowTokenTTL = 7 * 24 * time.Hour
const authCowIDKey = "auth_cow_id"

type cowJWTClaims struct {
	CowID string `json:"cow_id"`
	jwt.RegisteredClaims
}

func GenerateCowToken(cowID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(cowTokenTTL)
	claims := cowJWTClaims{
		CowID: cowID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   cowID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := jwtSecret()
	if err != nil {
		return "", time.Time{}, err
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func ParseCowToken(tokenString string) (*cowJWTClaims, error) {
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}

	claims := &cowJWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrBadToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrBadToken
	}
	if !token.Valid {
		return nil, ErrBadToken
	}
	if claims.CowID == "" {
		return nil, ErrBadToken
	}

	return claims, nil
}

func NeedCowAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := parseBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		claims, err := ParseCowToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Set(authCowIDKey, claims.CowID)
		c.Next()
	}
}

func GetAuthCowID(c *gin.Context) (string, error) {
	if value, ok := c.Get(authCowIDKey); ok {
		cowID, _ := value.(string)
		if cowID != "" {
			return cowID, nil
		}
	}
	return "", ErrBadToken
}
