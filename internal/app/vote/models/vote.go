package models

type Vote struct {
	Nickname string `json:"nickname"`
	ThreadID int    `json:"thread_id"`
	Voice    int    `json:"voice"`
}
