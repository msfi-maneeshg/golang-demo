package post

import "golang-demo/api/user"

type PostInfo struct {
	ID       int    `json:"id,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	Content  string `json:"content,omitempty"`
	Image    string `json:"image,omitempty"`
	Datetime string `json:"datetime,omitempty"`
	user.UserInformation
}
