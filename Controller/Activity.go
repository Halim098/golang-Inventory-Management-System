package Controller

import (
	"ims/Model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActivityByID(c *gin.Context) {
	id := c.Param("id")

	activity, err := Model.GetActivityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activity": activity})
}

func GetActivities(c *gin.Context) {
	changeType := c.Query("changetype")
	userName := c.Query("username")
	time := c.Query("time")
	productName := c.Query("product")
	query := Model.ActivityQuery(changeType, userName, time, productName)

	activities, err := Model.FilterActivity(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}
