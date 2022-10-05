package commondtos

type SuccessDto struct {
	success bool
}

func NewSuccessTrue() SuccessDto {
	return SuccessDto{success: true}
}
