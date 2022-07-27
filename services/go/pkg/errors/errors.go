package errors

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return v.Message
}

type RuleError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (v RuleError) Error() string {
	return v.Message
}

type BadRequestError struct {
	RuleErrors       []RuleError       `json:"ruleErrors"`
	ValidationErrors []ValidationError `json:"validationErrors"`
}

func (e BadRequestError) Error() string {
	return "Bad request error"
}

func (e BadRequestError) ContainsRuleErrorCode(code string) bool {
	for _, ruleError := range e.RuleErrors {
		if ruleError.Code == code {
			return true
		}
	}
	return false
}

func NewBadReqErrorFromRuleError(ruleError RuleError) BadRequestError {
	return NewBadReqErrorFromRuleErrors([]RuleError{ruleError}...)
}

func NewBadReqErrorFromRuleErrors(ruleError ...RuleError) BadRequestError {
	return BadRequestError{
		RuleErrors:       ruleError,
		ValidationErrors: []ValidationError{},
	}
}

func NewBadReqErrorFromValidationErrors(validationErrors []ValidationError) BadRequestError {
	return BadRequestError{
		RuleErrors:       []RuleError{},
		ValidationErrors: validationErrors,
	}
}
