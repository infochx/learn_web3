package api

import (
	"blog/config"
	"blog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	userId := c.Param("id")
	// 验证用户
	if _, err := strconv.Atoi(userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user := models.User{}
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 不返回密码信息
	user.Password = ""
	c.JSON(http.StatusOK, user)
}
