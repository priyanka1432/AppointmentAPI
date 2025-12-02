package errors

type AppError struct {
	Code    string
	Message string
	HTTP    int
}

func (e *AppError) Error() string { return e.Message }

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: "bad_request", Message: msg, HTTP: 400}
}
func NewConflict(msg string) *AppError { return &AppError{Code: "conflict", Message: msg, HTTP: 409} }
func NewInternal(msg string) *AppError {
	return &AppError{Code: "internal_error", Message: msg, HTTP: 500}
}
