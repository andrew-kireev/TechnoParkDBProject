package models

//easyjson:json
type User struct {
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}


func test() {
	us := &User{}
	us.MarshalJSON()
}