package responses

type Response struct {
	Message string `json:"message"`
}

func test() {
	r := &Response{}
	r.MarshalJSON()
}
