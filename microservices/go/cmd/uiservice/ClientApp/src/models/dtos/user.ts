interface UserIdentityDto extends UserReadDto, BaseUserAuthId {
  authorities: string[]
}

interface UserReadDto extends BaseUserPublicDto, ExistsDto {
}

interface UserSaveDto extends BaseUserCommon {}