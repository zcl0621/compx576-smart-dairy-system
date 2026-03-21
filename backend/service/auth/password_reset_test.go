package auth

import (
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	redisdb "github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	authdto "github.com/zcl0621/compx576-smart-dairy-system/dto/auth"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestRequestPasswordReset_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "reset-user", "reset@example.com", "oldpass")

		var sentEmail string
		var sentCode string
		originalSendResetCode := sendResetCode
		sendResetCode = func(toEmail, code string) error {
			sentEmail = toEmail
			sentCode = code
			return nil
		}
		defer func() {
			sendResetCode = originalSendResetCode
		}()

		resp, err := RequestPasswordReset(&authdto.PasswordResetRequestEmailRequest{Email: "reset@example.com"})

		require.NoError(t, err)
		require.Equal(t, "reset code sent", resp.Message)
		require.Equal(t, "reset@example.com", sentEmail)
		require.Len(t, sentCode, 6)

		storedEmail, err := redisdb.Get(buildPasswordResetCodeKey(sentCode))
		require.NoError(t, err)
		assert.Equal(t, "reset@example.com", storedEmail)

		latestCode, err := redisdb.Get(buildPasswordResetEmailCodeKey("reset@example.com"))
		require.NoError(t, err)
		assert.Equal(t, sentCode, latestCode)
	})
}

func TestRequestPasswordReset_ReplacesOldCode(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "reset-user", "reset@example.com", "oldpass")

		codes := make([]string, 0, 2)
		originalSendResetCode := sendResetCode
		sendResetCode = func(_ string, code string) error {
			codes = append(codes, code)
			return nil
		}
		defer func() {
			sendResetCode = originalSendResetCode
		}()

		_, err := RequestPasswordReset(&authdto.PasswordResetRequestEmailRequest{Email: "reset@example.com"})
		require.NoError(t, err)
		_, err = RequestPasswordReset(&authdto.PasswordResetRequestEmailRequest{Email: "reset@example.com"})
		require.NoError(t, err)
		require.Len(t, codes, 2)

		_, err = redisdb.Get(buildPasswordResetCodeKey(codes[0]))
		assert.True(t, errors.Is(err, redis.Nil))

		latestCode, err := redisdb.Get(buildPasswordResetEmailCodeKey("reset@example.com"))
		require.NoError(t, err)
		assert.Equal(t, codes[1], latestCode)
	})
}

func TestVerifyPasswordResetCode_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "reset-user", "reset@example.com", "oldpass")

		var sentCode string
		originalSendResetCode := sendResetCode
		sendResetCode = func(_ string, code string) error {
			sentCode = code
			return nil
		}
		defer func() {
			sendResetCode = originalSendResetCode
		}()

		_, err := RequestPasswordReset(&authdto.PasswordResetRequestEmailRequest{Email: "reset@example.com"})
		require.NoError(t, err)

		resp, err := VerifyPasswordResetCode(&authdto.PasswordResetVerifyRequest{Code: sentCode})

		require.NoError(t, err)
		require.NotEmpty(t, resp.ResetToken)

		storedEmail, err := redisdb.Get(buildPasswordResetTokenKey(resp.ResetToken))
		require.NoError(t, err)
		assert.Equal(t, "reset@example.com", storedEmail)

		_, err = redisdb.Get(buildPasswordResetCodeKey(sentCode))
		assert.True(t, errors.Is(err, redis.Nil))
	})
}

func TestConfirmPasswordReset_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "reset-user", "reset@example.com", "oldpass")

		resetToken := "reset-token-1"
		err := redisdb.Set(buildPasswordResetTokenKey(resetToken), "reset@example.com", passwordResetTokenTTL)
		require.NoError(t, err)

		resp, err := ConfirmPasswordReset(&authdto.PasswordResetConfirmRequest{
			ResetToken:  resetToken,
			NewPassword: "newpass123",
		})

		require.NoError(t, err)
		assert.Equal(t, "password changed", resp.Message)

		_, err = Login(&authdto.LoginRequest{Email: "reset@example.com", Password: "oldpass"})
		assert.ErrorIs(t, err, ErrBadLogin)

		loginResp, err := Login(&authdto.LoginRequest{Email: "reset@example.com", Password: "newpass123"})
		require.NoError(t, err)
		assert.NotEmpty(t, loginResp.Token)

		_, err = redisdb.Get(buildPasswordResetTokenKey(resetToken))
		assert.True(t, errors.Is(err, redis.Nil))
	})
}
