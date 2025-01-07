package endpoints

import (
	"encoding/json"
	"errors"
	"learn/db/redis"
	"learn/middlewares"
	"learn/models"
	"learn/schemas"
	"learn/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
	router.POST("/recend-verification-code", RecendVerificationCode)
	router.GET("/recovery-codes", middlewares.AuthMiddleware, MakeRecoveryCodes)
}

// @Schemes	http
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Param		login	body		schemas.AuthRequestSchema	true	"Login request"
// @Success	200		{object}	schemas.AuthResponseSchema
// @Router		/auth/login [post]
func Login(ctx *gin.Context) {
	loginData, err := utils.GetRequestBody[schemas.AuthRequestSchema](ctx)
	if err != nil {
		return
	}
	var user models.User
	tx := ctx.MustGet("tx").(*gorm.DB)
	result := tx.Model(&models.User{}).Where("email = ? and is_active", loginData.Email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
		return
	}
	if !utils.CheckPassword(loginData.Password, user.Password) {
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
	newUser, err := utils.GetRequestBody[schemas.UserCreate](ctx)
	if err != nil {
		return
	}
	tx := ctx.MustGet("tx").(*gorm.DB)
	newUser.Password = utils.MakePasswordHash(newUser.Password)
	result := tx.Create(&newUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
		return
	}

	verification_key, _ := utils.SendVerificationCode(newUser.Email)

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

// @Tags		Auth
// @Param		RecendVerification	body		schemas.RecendVerificationSchema	true	"Recend Verification Code"
// @Success	200	{object}	schemas.AuthVerificationKeySchema
// @Router		/auth/recend-verification-code [post]
func RecendVerificationCode(ctx *gin.Context) {
	data, err := utils.GetRequestBody[schemas.RecendVerificationSchema](ctx)
	if err != nil {
		return
	}

	var user models.User

	tx := ctx.MustGet("tx").(*gorm.DB)
	tx.Model(&models.User{}).Where("email = ?", data.Email).First(&user)

	if user.IsActive {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Email already verified"})
		return
	}

	verification_key, _ := utils.SendVerificationCode(data.Email)
	ctx.JSON(http.StatusOK, schemas.AuthVerificationKeySchema{VerificationKey: verification_key})
}

// @Tags		Auth
// @Success	200	{object}	schemas.AuthVerificationKeySchema
// @Router		/auth/recovery-codes [get]
// @Security BearerAuth
func MakeRecoveryCodes(ctx *gin.Context) {
	user := ctx.MustGet("user").(*schemas.UserRetrieveSchema)
	tx := ctx.MustGet("tx").(*gorm.DB)

	codes := utils.GenerateRecoveryCodes()
	codesJson, _ := json.Marshal(codes)

	recoveryCodes := models.UserRecoveryCodes{
		UserId: user.Id,
		Codes:  codesJson,
	}
	if err := tx.Save(&recoveryCodes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update recovery codes"})
		return
	}
	ctx.JSON(http.StatusOK, codes)
}
