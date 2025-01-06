package schemas

type AuthRequestSchema struct {
	EMAIL    string `json:"email" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

type AuthResponseSchema struct {
	Access string `json:"access"`
}

type RegisterSchema struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
