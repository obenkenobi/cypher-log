package gtools

import (
	"github.com/akrennmair/slice"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"google.golang.org/grpc"
)

type DialOptionSingleCreator func() single.Single[grpc.DialOption]

func CreateSingleWithDialOptions(
	dialOptionSingles []single.Single[grpc.DialOption]) single.Single[[]grpc.DialOption] {
	return slice.ReduceWithInitialValue(
		dialOptionSingles,
		single.Just([]grpc.DialOption{}),
		func(
			dialOptionsSrc single.Single[[]grpc.DialOption],
			dialOptSrc single.Single[grpc.DialOption],
		) single.Single[[]grpc.DialOption] {
			return single.Map(single.Zip2(dialOptionsSrc, dialOptSrc),
				func(zipped tuple.T2[[]grpc.DialOption, grpc.DialOption]) []grpc.DialOption {
					return append(zipped.V1, zipped.V2)
				},
			)
		},
	)
}

func CreateConnectionWithOptions(addr string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, options...)
}
