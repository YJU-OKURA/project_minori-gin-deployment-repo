package dto

import (
	"time"
)

// ClassScheduleDTO クラススケジュールDTO
type ClassScheduleDTO struct {
	Title     string    `json:"title" binding:"required"`
	StartedAt time.Time `json:"started_at" binding:"required"`
	EndedAt   time.Time `json:"ended_at" binding:"required"`
	CID       uint      `json:"cid" binding:"required"`
	IsLive    bool      `json:"is_live"`
}

// UpdateClassScheduleDTO クラススケジュール更新DTO
type UpdateClassScheduleDTO struct {
	Title     *string    `json:"title"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	IsLive    *bool      `json:"is_live"`
}
