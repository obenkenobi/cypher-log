package userdtos

type UserIdentityDto struct {
	AuthId      string   `json:"authId"`
	Authorities []string `json:"authorities"`
	User        UserDto  `json:"user"`
}

type UserDto struct {
	Id          string `json:"id"`
	Exists      bool   `json:"exists"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type UserSaveDto struct {
	UserName    string `json:"userName" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
}
