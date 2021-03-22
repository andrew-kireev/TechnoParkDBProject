package models

import (
	modelsForum "TechnoParkDBProject/internal/app/forum/models"
	modelsThread "TechnoParkDBProject/internal/app/thread/models"
	modelsUser "TechnoParkDBProject/internal/app/user/models"
)

type PostResponse struct {
	Post   *Post                `json:"post"`
	User   *modelsUser.User     `json:"author,omitempty"`
	Forum  *modelsForum.Forum   `json:"forum,omitempty"`
	Thread *modelsThread.Thread `json:"thread,omitempty"`
}
