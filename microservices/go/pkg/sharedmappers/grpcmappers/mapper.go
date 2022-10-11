package grpcmappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
)

func UserReadDtoToUserReply(dto *userdtos.UserReadDto, reply *userpb.UserReply) {
	reply.Id = dto.Id
	reply.Exists = dto.Exists
	reply.UserName = dto.UserName
	reply.DisplayName = dto.DisplayName
	reply.CreatedAt = dto.CreatedAt
	reply.UpdatedAt = dto.UpdatedAt
}

func UserReplyToUserReadDto(reply *userpb.UserReply, dto *userdtos.UserReadDto) {
	dto.Id = reply.Id
	dto.Exists = reply.Exists
	dto.UserName = reply.UserName
	dto.DisplayName = reply.DisplayName
	dto.CreatedAt = reply.CreatedAt
	dto.UpdatedAt = reply.UpdatedAt
}

func UserKeySessionDtoToUserKeySession(source *keydtos.UserKeySessionDto, dest *userkeypb.UserKeySession) {
	dest.ProxyKid = source.ProxyKid
	dest.Token = source.Token
}

func UserKeySessionToUserKeySessionDto(source *userkeypb.UserKeySession, dest *keydtos.UserKeySessionDto) {
	dest.ProxyKid = source.ProxyKid
	dest.Token = source.Token
}

func UserKeyDtoToUserKey(source *keydtos.UserKeyDto, dest *userkeypb.UserKey) {
	dest.KeyBase64 = source.KeyBase64
}

func UserKeyToUserKeyDto(source *userkeypb.UserKey, dest *keydtos.UserKeyDto) {
	dest.KeyBase64 = source.KeyBase64
}
