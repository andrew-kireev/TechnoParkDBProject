package models

type Request struct {
	Title   string      `json:"title"`
	Author  string      `json:"author"`
	Message string      `json:"message"`
	Created interface{} `json:"created"`
}