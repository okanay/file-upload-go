package httperrors

type HttpError struct {
	Message string
	Code    int
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(message string, code int) *HttpError {
	return &HttpError{Message: message, Code: code}
}
