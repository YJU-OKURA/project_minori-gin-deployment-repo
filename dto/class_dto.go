package dto

// CreateClassRequest クラス作成リクエストDTO
type CreateClassRequest struct {
	Name        string  `form:"name"`                   // クラス名
	Limitation  *int    `form:"limitation"`             // 参加制限人数
	Description *string `form:"description"`            // クラス説明
	UID         uint    `form:"uid" binding:"required"` // ユーザID
}
