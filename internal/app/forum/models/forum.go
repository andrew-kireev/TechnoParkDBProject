package models

//easyjson:json
type Forum struct {
	Tittle       string `json:"title"`
	UserNickname string `json:"user"`
	Slug         string `json:"slug"`
	Posts        int    `json:"posts"`
	Threads      int    `json:"threads"`
}
