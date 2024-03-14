package dto

// CreateClassRequest クラス作成リクエストDTO
type CreateClassRequest struct {
	Name        string  `form:"name"`        // 클래스 이름
	Limitation  *int    `form:"limitation"`  // 참여 인원 제한
	Description *string `form:"description"` // 클래스 설명
}
