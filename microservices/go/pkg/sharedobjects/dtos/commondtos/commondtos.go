package commondtos

type SuccessDto struct {
	Success bool `json:"success"`
}

func NewSuccessTrue() SuccessDto {
	return SuccessDto{Success: true}
}

type ExistsDto struct {
	Exists bool `json:"exists"`
}

func NewExistsDto(exists bool) ExistsDto {
	return ExistsDto{Exists: exists}
}
