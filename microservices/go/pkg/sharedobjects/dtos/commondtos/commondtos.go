package commondtos

type SuccessDto struct {
	success bool
}

func NewSuccessTrue() SuccessDto {
	return SuccessDto{success: true}
}

type ExistsDto struct {
	exists bool
}

func NewExistsDto(exists bool) ExistsDto {
	return ExistsDto{exists: exists}
}
