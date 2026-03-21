package auth

import (
	"errors"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	authdto "github.com/zcl0621/compx576-smart-dairy-system/dto/auth"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
	"go.uber.org/zap"
)

var ErrBadLogin = errors.New("email or password is wrong")

func Login(r *authdto.LoginRequest) (*authdto.LoginResponse, error) {
	var user model.User
	err := pg.DB.Model(&model.User{}).
		Where("email = ?", r.Email).
		First(&user).Error
	if err != nil {
		return nil, ErrBadLogin
	}

	if err := util.CheckPassword(user.Password, r.Password); err != nil {
		return nil, ErrBadLogin
	}

	return buildLoginResponse(&user)
}

func RefreshToken(userID string) (*authdto.LoginResponse, error) {
	var user model.User
	err := pg.DB.Model(&model.User{}).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return buildLoginResponse(&user)
}

func buildLoginResponse(user *model.User) (*authdto.LoginResponse, error) {
	tokenString, expiresAt, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		projectlog.L().Error("sign jwt failed",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		return nil, err
	}

	response := authdto.LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User: authdto.LoginUser{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		},
	}

	return &response, nil
}
