package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userdto "github.com/zcl0621/compx576-smart-dairy-system/dto/user"
	userservice "github.com/zcl0621/compx576-smart-dairy-system/service/user"
)

// UserList godoc
// @Summary list users
// @Description get user list
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param page query int false "page num" default(1)
// @Param page_size query int false "page size" default(20)
// @Param name query string false "user name"
// @Success 200 {object} userdto.ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/list [get]
func (h *Handler) UserList(c *gin.Context) {
	var query userdto.ListQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := userservice.UserListService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UserInfo godoc
// @Summary get user info
// @Description get one user by id
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id query string true "User ID"
// @Success 200 {object} userdto.InfoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/info [get]
func (h *Handler) UserInfo(c *gin.Context) {
	var query userdto.InfoQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := userservice.UserInfoService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UserCreate godoc
// @Summary create user
// @Description create one user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userdto.CreateRequest true "create user request"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/create [post]
func (h *Handler) UserCreate(c *gin.Context) {
	var request userdto.CreateRequest
	if !bindJSON(c, &request) {
		return
	}

	if err := userservice.UserCreateService(&request); err != nil {
		writeError(c, err)
		return
	}
	writeOK(c)
}

// UserUpdate godoc
// @Summary update user
// @Description update one user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userdto.UpdateRequest true "update user request"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/update [post]
func (h *Handler) UserUpdate(c *gin.Context) {
	var request userdto.UpdateRequest
	if !bindJSON(c, &request) {
		return
	}

	if err := userservice.UserUpdateService(&request); err != nil {
		writeError(c, err)
		return
	}
	writeOK(c)
}

// UserUpdatePassword godoc
// @Summary update user password
// @Description update password for one user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userdto.UpdatePasswordRequest true "update password request"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/update_password [post]
func (h *Handler) UserUpdatePassword(c *gin.Context) {
	var request userdto.UpdatePasswordRequest
	if !bindJSON(c, &request) {
		return
	}

	if err := userservice.UserUpdatePasswordService(&request); err != nil {
		writeError(c, err)
		return
	}
	writeOK(c)
}

// UserDelete godoc
// @Summary delete user
// @Description delete one user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userdto.DeleteRequest true "delete user request"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/delete [post]
func (h *Handler) UserDelete(c *gin.Context) {
	var request userdto.DeleteRequest
	if !bindJSON(c, &request) {
		return
	}

	if err := userservice.UserDeleteService(&request); err != nil {
		writeError(c, err)
		return
	}
	writeOK(c)
}
