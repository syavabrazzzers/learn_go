package endpoints

import (
	"errors"
	"fmt"
	"learn/middlewares"
	"learn/models"
	"learn/schemas"
	"learn/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/swaggo/swag"        // swag
	"gorm.io/gorm"
)

func AddAuth(eng *gin.Engine) {
	router := eng.Group("/api/auth", middlewares.TransactionMiddleware)
	router.POST("/login", Login)
	router.POST("/register", Register)
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
	result := tx.Model(&models.User{}).Where("email = ?", LoginData.EMAIL).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
		return
	}
	if !utils.CheckPassword(LoginData.PASSWORD, user.Password) {
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
// @Success	200		{object}	schemas.AuthResponseSchema
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
}
