package dto

type UserClassInfoDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Limitation  int    `json:"limitation"`
	Description string `json:"description"`
	Image       string `json:"image"`
	IsFavorite  bool   `json:"is_favorite"`
}
