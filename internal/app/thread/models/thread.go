package models

type Thread struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int    `json:"votes"`
	Slug    string `json:"slug,omitempty"`
	Created string `json:"created"`
}
