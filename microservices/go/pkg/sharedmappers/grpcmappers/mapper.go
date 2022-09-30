package grpcmappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
)

func MapUserReadDtoToUserReply(dto *userdtos.UserReadDto, reply *userpb.UserReply) {
	reply.Id = dto.Id
	reply.Exists = dto.Exists
	reply.UserName = dto.UserName
	reply.DisplayName = dto.DisplayName
	reply.CreatedAt = dto.CreatedAt
	reply.UpdatedAt = dto.UpdatedAt
}

func MapUserReplyToUserReadDto(reply *userpb.UserReply, dto *userdtos.UserReadDto) {
	dto.Id = reply.Id
	dto.Exists = reply.Exists
	dto.UserName = reply.UserName
	dto.DisplayName = reply.DisplayName
	dto.CreatedAt = reply.CreatedAt
	dto.UpdatedAt = reply.UpdatedAt
}
