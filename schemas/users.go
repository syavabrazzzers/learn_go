package schemas

type UserCreate struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UserRetrieveSchema struct {
	Id         int    `json:"id" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Created_at string `json:"created_at"`
	Username   string `json:"username"`
}
