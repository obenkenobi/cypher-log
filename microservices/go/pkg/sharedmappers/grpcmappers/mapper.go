package grpcmappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
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
	dto.Id = reply.GetId()
	dto.Exists = reply.GetExists()
	dto.UserName = reply.GetUserName()
	dto.DisplayName = reply.GetDisplayName()
	dto.CreatedAt = reply.GetCreatedAt()
	dto.UpdatedAt = reply.GetUpdatedAt()
}

func UserKeySessionDtoToUserKeySession(source *commondtos.UKeySessionDto, dest *userkeypb.UserKeySession) {
	dest.ProxyKid = source.ProxyKid
	dest.Token = source.Token
	dest.UserId = source.UserId
	dest.KeyVersion = source.KeyVersion
	dest.StartTime = source.StartTime
	dest.DurationMilli = source.DurationMilli
}

func UserKeySessionToUserKeySessionDto(source *userkeypb.UserKeySession, dest *commondtos.UKeySessionDto) {
	dest.ProxyKid = source.GetProxyKid()
	dest.Token = source.GetToken()
	dest.UserId = source.GetUserId()
	dest.KeyVersion = source.GetKeyVersion()
	dest.StartTime = source.GetStartTime()
	dest.DurationMilli = source.GetDurationMilli()
}

func UserKeyDtoToUserKey(source *keydtos.UserKeyDto, dest *userkeypb.UserKey) {
	dest.KeyBase64 = source.KeyBase64
	dest.KeyVersion = source.KeyVersion
}

func UserKeyToUserKeyDto(source *userkeypb.UserKey, dest *keydtos.UserKeyDto) {
	dest.KeyBase64 = source.GetKeyBase64()
	dest.KeyVersion = source.GetKeyVersion()
}
