package endpoints

import (
	"learn/db"
	"learn/middlewares"
	"learn/schemas"
	"learn/utils"

	"github.com/gin-gonic/gin"

	"learn/models"
)

func AddUsers(eng *gin.Engine) {
	router := eng.Group("/users", middlewares.TransactionMiddleware, middlewares.AuthMiddleware)
	router.GET("/", GetUsers)
	router.POST("/", CreateUser)
	router.DELETE("/", DeleteUser)
}

func GetUsers(c *gin.Context) {
	c.Get("userId")
	var users []models.User // Create a slice to hold the users
	result := db.Database.Select("*").Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()}) // Handle error
		return
	}

	c.JSON(200, users)
}

func CreateUser(c *gin.Context) {
	var data schemas.UserCreate
	c.BindJSON(&data)
	password_hash := utils.MakePasswordHash(data.Password)
	user := &models.User{Email: data.Email, Password: password_hash}
	db.Database.Create(user)
}

func DeleteUser(c *gin.Context) {
	userId := c.GetInt("userId")
	var user models.User
	db.Database.Model(&models.User{}).Where("id = ?", userId).First(user)
	db.Database.Delete(&models.User{}, &models.User{Email: user.Email})
}
