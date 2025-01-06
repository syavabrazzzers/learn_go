package schemas

type AuthRequestSchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponseSchema struct {
	Access string `json:"access"`
}

type RegisterSchema struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthVerificationKeySchema struct {
	VerificationKey string `json:"verification_key"`
}

type AuthVerificationSchema struct {
	VerificationKey string `json:"verification_key"`
	Code            string `json:"code"`
}

type RecendVerificationSchema struct {
	Email string `json:"email" binding:"required"`
}
