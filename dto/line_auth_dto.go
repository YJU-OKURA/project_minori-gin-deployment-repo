package dto

type LineUserInput struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	StatusMessage string `json:"statusMessage"`
	PictureURL    string `json:"pictureUrl"`
}
