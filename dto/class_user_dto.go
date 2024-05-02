package dto

type UserClassInfoDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Limitation  int    `json:"limitation"`
	Description string `json:"description"`
	Image       string `json:"image"`
	IsFavorite  bool   `json:"is_favorite"`
	Role        string `json:"role"`
}

type ClassMemberDTO struct {
	Uid      uint   `json:"uid"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Image    string `json:"image"`
}
