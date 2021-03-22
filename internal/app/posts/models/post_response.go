package models

type PostResponse struct {
	Post *Post `json:"post"`
}

func test() {
	p := &PostResponse{}
	p.MarshalJSON()
}