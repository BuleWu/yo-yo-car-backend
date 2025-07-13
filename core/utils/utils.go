package utils

type HttpError struct {
	Message string `json:"message"`
}

func NewHttpError(message string) *HttpError {
	return &HttpError{Message: message}
}
