package auth

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/xid"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	redisdb "github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	authdto "github.com/zcl0621/compx576-smart-dairy-system/dto/auth"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	mailservice "github.com/zcl0621/compx576-smart-dairy-system/service/mail"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
)

const (
	passwordResetCodeTTL     = 10 * time.Minute
	passwordResetTokenTTL    = 15 * time.Minute
	passwordResetMaxAttempts = 5
)

var (
	ErrResetCodeNotFound  = errors.New("invalid or expired code")
	ErrResetTokenNotFound = errors.New("reset token is invalid or expired")
	sendResetCode         = mailservice.SendResetCode
)

func RequestPasswordReset(r *authdto.PasswordResetRequestEmailRequest) (*authdto.MessageResponse, error) {
	var user model.User
	err := pg.DB.Model(&model.User{}).
		Where("email = ?", r.Email).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	code, err := generateResetCode()
	if err != nil {
		return nil, err
	}

	emailCodeKey := buildPasswordResetEmailCodeKey(r.Email)
	oldCode, err := redisdb.Get(emailCodeKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if oldCode != "" {
		if err := redisdb.Del(
			buildPasswordResetCodeKey(oldCode),
			buildPasswordResetAttemptsKey(oldCode),
		); err != nil {
			return nil, err
		}
	}

	if err := redisdb.Set(buildPasswordResetCodeKey(code), r.Email, passwordResetCodeTTL); err != nil {
		return nil, err
	}
	if err := redisdb.Set(emailCodeKey, code, passwordResetCodeTTL); err != nil {
		return nil, err
	}
	if err := redisdb.Del(buildPasswordResetAttemptsKey(code)); err != nil {
		return nil, err
	}

	if err := sendResetCode(r.Email, code); err != nil {
		_ = redisdb.Del(buildPasswordResetCodeKey(code), emailCodeKey)
		return nil, err
	}

	return &authdto.MessageResponse{Message: "reset code sent"}, nil
}

func VerifyPasswordResetCode(r *authdto.PasswordResetVerifyRequest) (*authdto.PasswordResetVerifyResponse, error) {
	codeKey := buildPasswordResetCodeKey(r.Code)
	email, err := redisdb.Get(codeKey)
	if errors.Is(err, redis.Nil) {
		if increaseErr := increasePasswordResetAttempts(r.Code); increaseErr != nil {
			return nil, increaseErr
		}
		return nil, ErrResetCodeNotFound
	}
	if err != nil {
		return nil, err
	}

	resetToken := xid.New().String()
	if err := redisdb.Set(buildPasswordResetTokenKey(resetToken), email, passwordResetTokenTTL); err != nil {
		return nil, err
	}

	if err := redisdb.Del(
		codeKey,
		buildPasswordResetAttemptsKey(r.Code),
		buildPasswordResetEmailCodeKey(email),
	); err != nil {
		return nil, err
	}

	return &authdto.PasswordResetVerifyResponse{ResetToken: resetToken}, nil
}

func ConfirmPasswordReset(r *authdto.PasswordResetConfirmRequest) (*authdto.MessageResponse, error) {
	tokenKey := buildPasswordResetTokenKey(r.ResetToken)
	email, err := redisdb.Get(tokenKey)
	if errors.Is(err, redis.Nil) {
		return nil, ErrResetTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(r.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := pg.DB.Model(&model.User{}).
		Where("email = ?", email).
		Update("password", hashedPassword).Error; err != nil {
		return nil, err
	}

	if err := redisdb.Del(tokenKey); err != nil {
		return nil, err
	}

	return &authdto.MessageResponse{Message: "password changed"}, nil
}

func increasePasswordResetAttempts(code string) error {
	attemptsKey := buildPasswordResetAttemptsKey(code)
	attempts, err := redisdb.Incr(attemptsKey)
	if err != nil {
		return err
	}
	if attempts == 1 {
		if err := redisdb.Expire(attemptsKey, passwordResetCodeTTL); err != nil {
			return err
		}
	}
	if attempts < passwordResetMaxAttempts {
		return nil
	}

	email, err := redisdb.Get(buildPasswordResetCodeKey(code))
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	keys := []string{buildPasswordResetCodeKey(code), attemptsKey}
	if email != "" {
		keys = append(keys, buildPasswordResetEmailCodeKey(email))
	}

	return redisdb.Del(keys...)
}

func generateResetCode() (string, error) {
	for i := 0; i < 10; i++ {
		buf := make([]byte, 4)
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}

		code := fmt.Sprintf("%06d", binary.BigEndian.Uint32(buf)%1000000)

		exists, err := redisdb.Exists(buildPasswordResetCodeKey(code))
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", errors.New("could not generate reset code")
}

func buildPasswordResetCodeKey(code string) string {
	return "reset:code:" + code
}

func buildPasswordResetAttemptsKey(code string) string {
	return "reset:attempts:" + code
}

func buildPasswordResetEmailCodeKey(email string) string {
	return "reset:email-code:" + email
}

func buildPasswordResetTokenKey(token string) string {
	return "reset:token:" + token
}
