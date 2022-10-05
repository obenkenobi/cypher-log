package keydtos

type PasscodeCreateDto struct {
	Passcode string `json:"Passcode" binding:"exists,alphanumunicode,min=4,max=20"`
}

type PasscodeDto struct {
	Passcode string `json:"Passcode" binding:"required"`
}
