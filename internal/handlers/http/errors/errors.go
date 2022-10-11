package errors

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ValidationError struct {
	FailedField string `json:"failedField"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

type ValidationErrorResponse struct {
	ErrorResponse

	Errors []ValidationError `json:"errors"`
}