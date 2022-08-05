package gtools

// GrpcAction flags that indicate what kind of operation is done in a GRPC
// service. This is not used in the grpc standard and instead is used to augment
// functionality on top of GRPC like error handling.
type GrpcAction int

const (
	ReadAction GrpcAction = iota
	CreateAction
	UpdateAction
	DeleteAction
	OtherAction
)
