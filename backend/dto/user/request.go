package user

import "github.com/zcl0621/compx576-smart-dairy-system/dto/common"

type ListQuery struct {
	common.PageQuery
	Name string `form:"name"`
}

type InfoQuery struct {
	ID string `form:"id" binding:"required"`
}

type CreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateRequest struct {
	ID       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdatePasswordRequest struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DeleteRequest struct {
	ID string `json:"id" binding:"required"`
}
