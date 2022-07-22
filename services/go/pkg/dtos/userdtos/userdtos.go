package userdtos

type UserDto struct {
	Id          string `json:"id"`
	Exists      bool   `json:"exists"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
}

type UserSaveDto struct {
	UserName    string `json:"userName" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
}
