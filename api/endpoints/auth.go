package endpoints

import (
	"errors"
	"fmt"
	"learn/db/redis"
	"learn/middlewares"
	"learn/models"
	"learn/schemas"
	"learn/utils"
	"learn/utils/smtp"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ext_redis "github.com/redis/go-redis/v9"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/swaggo/swag"        // swag
	"gorm.io/gorm"
)

func AddAuth(eng *gin.Engine) {
	router := eng.Group("/api/auth", middlewares.TransactionMiddleware)
	router.POST("/login", Login)
	router.POST("/register", Register)
	router.POST("/verify", Verify)
}

// @Schemes	http
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Param		login	body		schemas.AuthRequestSchema	true	"Login request"
// @Success	200		{object}	schemas.AuthResponseSchema
// @Router		/auth/login [post]
func Login(ctx *gin.Context) {
	var LoginData schemas.AuthRequestSchema
	var user models.User
	if err := ctx.BindJSON(&LoginData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}
	tx := ctx.MustGet("tx").(*gorm.DB)
	result := tx.Model(&models.User{}).Where("email = ? and is_active", LoginData.Email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
		return
	}
	if !utils.CheckPassword(LoginData.Password, user.Password) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Invalid credentials"})
		return
	}
	token := utils.MakeJwt(user)
	ctx.JSON(http.StatusOK, gin.H{"access": token})
}

// @Schemes	http
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Param		login	body		schemas.UserCreate	true	"Login request"
// @Success	200		{object}	schemas.AuthVerificationKeySchema
// @Router		/auth/register [post]
func Register(ctx *gin.Context) {
	var newUser models.User
	if err := ctx.BindJSON(&newUser); err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}
	tx := ctx.MustGet("tx").(*gorm.DB)
	newUser.Password = utils.MakePasswordHash(newUser.Password)
	result := tx.Create(&newUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
		return
	}
	verification_key := uuid.New().String()
	verification_code := fmt.Sprint(utils.MakeVerificationCode())
	go redis.Client.SetJson(verification_key, map[string]string{"email": newUser.Email, "code": verification_code})
	go smtp.SendMail([]string{newUser.Email}, verification_code, "asd")
	ctx.JSON(http.StatusOK, schemas.AuthVerificationKeySchema{VerificationKey: verification_key})
}

// @Schemes	http
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Param		verification	body	schemas.AuthVerificationSchema true	"Verification data"
// @Success	200		{object}	schemas.AuthVerificationKeySchema
// @Router		/auth/verify [post]
func Verify(ctx *gin.Context) {
	data, err := utils.GetRequestBody[schemas.AuthVerificationSchema](ctx)
	if err != nil {
		return
	}
	redis_data, err := redis.Client.GetJson(data.VerificationKey)
	if errors.Is(err, ext_redis.Nil) {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if data.Code != redis_data["code"] {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"detail": "Invalid code"})
		return
	}

	tx := ctx.MustGet("tx").(*gorm.DB)
	tx.Model(&models.User{}).Where("email = ?", redis_data["email"]).Update("is_active", true)
}
