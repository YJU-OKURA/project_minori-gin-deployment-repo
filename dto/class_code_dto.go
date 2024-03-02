package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ClassCodeCheckRequest はグループコードの存在を確認するリクエストです。
type ClassCodeCheckRequest struct {
	Code string `json:"code" binding:"required"`
}

// ClassCodeRequest Bind はリクエストのバインドを行います。
type ClassCodeRequest struct {
	Code   string `json:"code" binding:"required"`
	Secret string `json:"secret,omitempty"`
	UID    uint   `json:"uid" binding:"required"`
}

// Bind はリクエストのバインドを行います。
func (r *ClassCodeRequest) Bind(c *gin.Context) error {
	return c.ShouldBindWith(r, binding.JSON)
}
