package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authdto "github.com/zcl0621/compx576-smart-dairy-system/dto/auth"
	"github.com/zcl0621/compx576-smart-dairy-system/service/auth"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestLogin_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		u := testhelper.SeedUser(t, tx, "alice", "alice@example.com", "pass123")

		resp, err := auth.Login(&authdto.LoginRequest{
			Email:    u.Email,
			Password: "pass123",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, u.ID, resp.User.ID)
		assert.Equal(t, u.Email, resp.User.Email)
	})
}

func TestLogin_WrongPassword(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "bob", "bob@example.com", "correct")

		_, err := auth.Login(&authdto.LoginRequest{
			Email:    "bob@example.com",
			Password: "wrong",
		})

		assert.ErrorIs(t, err, auth.ErrBadLogin)
	})
}

func TestLogin_UserNotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		_, err := auth.Login(&authdto.LoginRequest{
			Email:    "nobody@example.com",
			Password: "pass",
		})

		assert.ErrorIs(t, err, auth.ErrBadLogin)
	})
}

func TestRefreshToken_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		u := testhelper.SeedUser(t, tx, "carol", "carol@example.com", "pass")

		resp, err := auth.RefreshToken(u.ID)

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, u.ID, resp.User.ID)
	})
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		_, err := auth.RefreshToken("nonexistent-id")

		assert.Error(t, err)
	})
}
