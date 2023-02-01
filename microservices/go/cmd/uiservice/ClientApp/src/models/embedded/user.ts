interface BaseUserCommon {
  userName: string,
  displayName: string
}

interface BaseUserAuthId {
  authId: string,
}

interface BaseUserPublicDto extends BaseUserCommon, BaseCRUDObject {}