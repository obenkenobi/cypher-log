package errordtos

type ValidationErrDto struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BrViolationDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponseDto struct {
	IsInternalError  bool               `json:"isSystemError"`
	BrViolations     []BrViolationDto   `json:"brViolations"`
	ValidationErrors []ValidationErrDto `json:"validationErrors"`
	Message          string             `json:"message"`
}

func NewInternalErrorResponse() *ErrorResponseDto {
	return NewErrorResponse("Internal Error", true)
}

func NewErrorResponse(msg string, isInternalError bool) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  isInternalError,
		BrViolations:     []BrViolationDto{},
		ValidationErrors: []ValidationErrDto{},
		Message:          msg,
	}
}

func NewSingleBRViolationErrorResponse(code string, msg string) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  false,
		BrViolations:     []BrViolationDto{{Code: code, Message: msg}},
		ValidationErrors: []ValidationErrDto{},
		Message:          "",
	}
}

func NewBRViolationErrorResponse(brViolations []BrViolationDto) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  false,
		BrViolations:     brViolations,
		ValidationErrors: []ValidationErrDto{},
		Message:          "",
	}
}

func NewValidationErrorResponse(validationErrors []ValidationErrDto) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  false,
		BrViolations:     []BrViolationDto{},
		ValidationErrors: validationErrors,
		Message:          "",
	}
}
