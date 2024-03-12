package entity

type ResponseWrapper struct {
	Data interface{} `json:"data"`
}

type ErrorWrapper struct {
	Error error `json:"error"`
}
