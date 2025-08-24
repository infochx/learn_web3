package api

import (
	"blog/config"
	"blog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建评论
func CreateComment(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// 获取请求体中的评论信息
	comment := models.Comment{}
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 设置评论的用户ID
	comment.UserID = userID.(uint)
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// 获取文章的所有评论信息
func GetCommentsByPostId(c *gin.Context) {
	postID := c.Param("postId")
	// 校验文章ID有效性
	if _, err := strconv.Atoi(postID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid postID"})
		return
	}
	// 根据文章ID查询所有评论信息
	comments := []models.Comment{}
	if err := config.DB.Where("post_id =?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}
