package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authdto "github.com/zcl0621/compx576-smart-dairy-system/dto/auth"
	userdto "github.com/zcl0621/compx576-smart-dairy-system/dto/user"
	"github.com/zcl0621/compx576-smart-dairy-system/service/auth"
	"github.com/zcl0621/compx576-smart-dairy-system/service/user"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestUserCreate_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		err := user.UserCreateService(&userdto.CreateRequest{
			Username: "dave",
			Email:    "dave@example.com",
			Password: "pass",
		})

		require.NoError(t, err)

		resp, err := auth.Login(&authdto.LoginRequest{
			Email:    "dave@example.com",
			Password: "pass",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})
}

func TestUserList_Pagination(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "u1", "u1@example.com", "pass")
		testhelper.SeedUser(t, tx, "u2", "u2@example.com", "pass")
		testhelper.SeedUser(t, tx, "u3", "u3@example.com", "pass")

		resp, err := user.UserListService(&userdto.ListQuery{})

		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.List, 3)
	})
}

func TestUserList_SearchByName(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "searchme", "searchme@example.com", "pass")
		testhelper.SeedUser(t, tx, "other", "other@example.com", "pass")

		resp, err := user.UserListService(&userdto.ListQuery{Name: "searchme"})

		require.NoError(t, err)
		require.Equal(t, int64(1), resp.Total)
		assert.Equal(t, "searchme", resp.List[0].Username)
	})
}

func TestUserInfo_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		u := testhelper.SeedUser(t, tx, "eve", "eve@example.com", "pass")

		resp, err := user.UserInfoService(&userdto.InfoQuery{ID: u.ID})

		require.NoError(t, err)
		assert.Equal(t, u.ID, resp.ID)
		assert.Equal(t, "eve", resp.Username)
		assert.Equal(t, "eve@example.com", resp.Email)
	})
}

func TestUserInfo_NotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		_, err := user.UserInfoService(&userdto.InfoQuery{ID: "nonexistent"})

		assert.Error(t, err)
	})
}

func TestUserUpdate_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		u := testhelper.SeedUser(t, tx, "frank", "frank@example.com", "pass")

		err := user.UserUpdateService(&userdto.UpdateRequest{
			ID:       u.ID,
			Username: "frank-updated",
			Email:    "frank2@example.com",
		})
		require.NoError(t, err)

		info, err := user.UserInfoService(&userdto.InfoQuery{ID: u.ID})
		require.NoError(t, err)
		assert.Equal(t, "frank-updated", info.Username)
		assert.Equal(t, "frank2@example.com", info.Email)
	})
}

func TestUserUpdatePassword_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		u := testhelper.SeedUser(t, tx, "grace", "grace@example.com", "oldpass")

		err := user.UserUpdatePasswordService(&userdto.UpdatePasswordRequest{
			ID:       u.ID,
			Password: "newpass",
		})
		require.NoError(t, err)

		resp, err := auth.Login(&authdto.LoginRequest{
			Email:    "grace@example.com",
			Password: "newpass",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})
}

func TestUserDelete_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		u := testhelper.SeedUser(t, tx, "heidi", "heidi@example.com", "pass")

		err := user.UserDeleteService(&userdto.DeleteRequest{ID: u.ID})
		require.NoError(t, err)

		_, err = user.UserInfoService(&userdto.InfoQuery{ID: u.ID})
		assert.Error(t, err)
	})
}

func TestUserList_SearchByEmail(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedUser(t, tx, "ivan", "findbyemail@example.com", "pass")
		testhelper.SeedUser(t, tx, "judy", "other@example.com", "pass")

		resp, err := user.UserListService(&userdto.ListQuery{Name: "findbyemail"})

		require.NoError(t, err)
		require.Equal(t, int64(1), resp.Total)
		assert.Equal(t, "ivan", resp.List[0].Username)
	})
}

func TestUserUpdate_NotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		err := user.UserUpdateService(&userdto.UpdateRequest{
			ID:       "nonexistent-user-id",
			Username: "ghost",
			Email:    "ghost@example.com",
		})

		assert.Error(t, err)
	})
}

func TestUserUpdatePassword_NotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		err := user.UserUpdatePasswordService(&userdto.UpdatePasswordRequest{
			ID:       "nonexistent-user-id",
			Password: "newpass",
		})

		assert.Error(t, err)
	})
}
