package errordtos

type ValidationErrDto struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppErrorDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (v AppErrorDto) Error() string {
	return v.Message
}

type ErrorResponseDto struct {
	IsInternalError  bool               `json:"isSystemError"`
	AppErrors        []AppErrorDto      `json:"brViolations"`
	ValidationErrors []ValidationErrDto `json:"validationErrors"`
}

func (e ErrorResponseDto) ContainsAppErrorCode(code string) bool {
	for _, appError := range e.AppErrors {
		if appError.Code == code {
			return true
		}
	}
	return false
}

func NewInternalErrResponse() *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  true,
		AppErrors:        []AppErrorDto{},
		ValidationErrors: []ValidationErrDto{},
	}
}

func NewSingleAppErrorResponse(appErrorDto AppErrorDto) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  false,
		AppErrors:        []AppErrorDto{appErrorDto},
		ValidationErrors: []ValidationErrDto{},
	}
}

func NewAppErrorsResponse(appErrors []AppErrorDto) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  false,
		AppErrors:        appErrors,
		ValidationErrors: []ValidationErrDto{},
	}
}

func NewValidationErrorResponse(validationErrors []ValidationErrDto) *ErrorResponseDto {
	return &ErrorResponseDto{
		IsInternalError:  false,
		AppErrors:        []AppErrorDto{},
		ValidationErrors: validationErrors,
	}
}
