package controller

import (
	"jade-mes/app/infrastructure/log"
	"jade-mes/app/user_permission/services"

	"github.com/gin-gonic/gin"
)

// RoleController ...
type RoleController struct{}

func (rc *RoleController) CreateRole(c *gin.Context) {
	var (
		err   error
		ret   interface{}
		param struct {
			Id   int64  `json:"id" binding:"required"`
			Name string `json:"name" binding:"required"`
			Desc string `json:"desc" binding:"required"`
		}
	)
	defer func() {
		log.Debug("CreateRole", log.Reflect("param", param), log.Reflect("ret", ret), log.Err(err))
		response(c, ret, err)
	}()

	err = services.CreateRole(c.Request.Context(), param.Id, param.Name, param.Desc)
	return
}

func (rc *RoleController) FindRoleByID(c *gin.Context) {
	var (
		err   error
		ret   interface{}
		param struct {
			Id   int64  `json:"id" binding:"required"`
		}
	)
	defer func() {
		log.Debug("FindRoleByID", log.Reflect("param", param), log.Reflect("ret", ret), log.Err(err))
		response(c, ret, err)
	}()

	ret, err = services.FindRoleByID(c.Request.Context(), param.Id)
	return
}

func (rc *RoleController) FindRoleByName(c *gin.Context) {
	var (
		err   error
		ret   interface{}
		param struct {
			Name   string  `json:"name" binding:"required"`
		}
	)
	defer func() {
		log.Debug("FindRoleByName", log.Reflect("param", param), log.Reflect("ret", ret), log.Err(err))
		response(c, ret, err)
	}()

	ret, err = services.FindRoleByName(c.Request.Context(), param.Name)
	return
}
