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
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"detail": "Unathorized"})
		return
	}
	verifiedToken, err := utils.VerifyJwt(bearer[7:])
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	user, err := verifiedToken.Claims.GetSubject()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
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
