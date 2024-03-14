package utils

type CustomError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *CustomError) Error() string {
	return e.Message
}
