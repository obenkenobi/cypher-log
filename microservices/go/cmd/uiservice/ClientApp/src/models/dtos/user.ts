interface UserIdentityDto extends UserReadDto, BaseUserAuthId {
  authorities: string[]
}

interface UserReadDto extends BaseUserPublicDto {
  exists: boolean
}

interface UserSaveDto extends BaseUserCommon {}