package Controller

import (
	"ims/Helper"
	"ims/Model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user Model.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = user.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	var userLogin Model.UserLogin
	err := c.BindJSON(&userLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Model.FindUserByUsername(userLogin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = user.ValidatePassword(userLogin.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid login credentials",
		})
		return
	}

	token, err := Helper.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"massage": "Login success", "token": token})
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "IMS_Cookie",
		Value:  "",
		MaxAge: -1,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logout success."})
}

func GetUser(c *gin.Context) {
	id := c.MustGet("user_id")
	users, err := Model.GetUserByID(id.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
