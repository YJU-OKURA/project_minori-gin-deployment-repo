package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GroupCodeCheckRequest はグループコードの存在を確認するリクエストです。
type GroupCodeCheckRequest struct {
	Code string `json:"code" binding:"required"`
}

// GroupCodeRequest Bind はリクエストのバインドを行います。
type GroupCodeRequest struct {
	Code   string `json:"code" binding:"required"`
	Secret string `json:"secret,omitempty"`
}

// Bind はリクエストのバインドを行います。
func (r *GroupCodeRequest) Bind(c *gin.Context) error {
	return c.ShouldBindWith(r, binding.JSON)
}
