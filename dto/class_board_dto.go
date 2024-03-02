package dto

import "mime/multipart"

// ClassBoardCreateDTO - グループ掲示板を作成するためのDTO
type ClassBoardCreateDTO struct {
	Title       string                `json:"title" form:"title"  binding:"required" example:"Sample Title"`
	Content     string                `json:"content" form:"content"  binding:"required" example:"Sample Content"`
	Image       *multipart.FileHeader `form:"image"`
	ImageURL    string
	IsAnnounced bool `json:"is_announced" form:"is_announced" default:"false"`
	CID         uint `json:"cid" form:"cid"  binding:"required"`
	UID         uint `json:"uid" form:"uid"  binding:"required"`
}

// ClassBoardUpdateDTO - グループ掲示板を更新するためのDTO
type ClassBoardUpdateDTO struct {
	ID          uint   `json:"id" form:"id"  binding:"required"`
	Title       string `json:"title" form:"title"`
	Content     string `json:"content" form:"content"`
	Image       string `json:"image" form:"image"`
	IsAnnounced bool   `json:"is_announced" form:"is_announced"`
}
