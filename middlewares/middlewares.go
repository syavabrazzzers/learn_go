package middlewares

import (
	"fmt"
	"learn/db"
	"learn/models"
	"learn/schemas"
	"learn/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionMiddleware(ctx *gin.Context) {
	tx := db.Database.Begin()
	ctx.Set("tx", tx)
	ctx.Next()
	if ctx.Errors.Last() != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func AuthMiddleware(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"detail": "Unathorized"})
	}
	verifiedToken, err := utils.VerifyJwt(token)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusUnauthorized, err)
	}
	user, err := verifiedToken.Claims.GetSubject()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
	}
	var userSchema schemas.UserRetrieveSchema
	tx, ok := c.MustGet("tx").(*gorm.DB)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tx.Model(&models.User{}).Where("id = ?", user).First(&userSchema)
	c.Set("user", &userSchema)
	c.Next()
}
