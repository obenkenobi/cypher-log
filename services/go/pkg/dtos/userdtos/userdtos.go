package userdtos

type UserDto struct {
	Id          string `json:"id"`
	UserAdded   bool   `json:"userAdded"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
}

type UserSaveDto struct {
	UserName    string `json:"userName" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
}
