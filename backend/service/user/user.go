package user

import (
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/user"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
)

func UserListService(r *user.ListQuery) (*user.ListResponse, error) {
	db := pg.DB.Model(&model.User{})
	db = db.Order("id desc")
	if r.Name != "" {
		keyword := "%" + r.Name + "%"
		db = db.Where("username LIKE ? OR email LIKE ?", keyword, keyword)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []user.ListItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	response := user.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}

	return &response, nil
}

func UserInfoService(r *user.InfoQuery) (*user.InfoResponse, error) {
	db := pg.DB.Model(&model.User{})
	db = db.Where("id = ?", r.ID)

	var response user.InfoResponse
	if err := db.First(&response).Error; err != nil {
		return nil, err
	}

	return &response, nil
}

func UserCreateService(r *user.CreateRequest) error {
	hashedPassword, err := util.HashPassword(r.Password)
	if err != nil {
		return err
	}

	item := model.User{
		Username: r.Username,
		Password: hashedPassword,
		Email:    r.Email,
	}

	return pg.DB.Create(&item).Error
}

func UserUpdateService(r *user.UpdateRequest) error {
	db := pg.DB.Model(&model.User{})
	db = db.Where("id = ?", r.ID)

	var item model.User
	if err := db.First(&item).Error; err != nil {
		return err
	}

	item.Username = r.Username
	item.Email = r.Email

	return pg.DB.Save(&item).Error
}

func UserUpdatePasswordService(r *user.UpdatePasswordRequest) error {
	db := pg.DB.Model(&model.User{})
	db = db.Where("id = ?", r.ID)

	var item model.User
	if err := db.First(&item).Error; err != nil {
		return err
	}

	hashedPassword, err := util.HashPassword(r.Password)
	if err != nil {
		return err
	}

	item.Password = hashedPassword

	return pg.DB.Save(&item).Error
}

func UserDeleteService(r *user.DeleteRequest) error {
	db := pg.DB.Model(&model.User{})
	db = db.Where("id = ?", r.ID)
	return db.Delete(&model.User{}).Error
}
