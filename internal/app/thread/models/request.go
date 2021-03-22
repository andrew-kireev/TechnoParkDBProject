package models

type Request struct {
	Title   string      `json:"title"`
	Author  string      `json:"author"`
	Message string      `json:"message"`
	Created interface{} `json:"created"`
}

type UpdateRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
