package api

import (
	"blog/config"
	"blog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有文章
func GetPosts(c *gin.Context) {
	posts := []models.Post{}
	if err := config.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "posts are not found"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// 获取单个文章
func GetPost(c *gin.Context) {
	postId := c.Param("id")
	post := models.Post{}
	if err := config.DB.Preload("User").Where("id =?", postId).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post is not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// 创建文章
func CreatePost(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authorized"})
		return
	}
	post := models.Post{}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.UserID = userId.(uint)

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// 更新文章
func UpdatePost(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// 获取请求体中的文章信息
	reqPost := models.Post{}
	if err := c.ShouldBindJSON(&reqPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取数据库中的文章信息
	storagePost := models.Post{}
	if err := config.DB.Where("id =?", reqPost.ID).First(&storagePost).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 只有文章作者才能更新
	if userId != storagePost.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have no permission to update this post"})
		return
	}

	if err := config.DB.Debug().Updates(&reqPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, reqPost)
}

// 删除文章
func DeletePost(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}
	// 获取请求参数中的文章ID
	postId := c.Param("id")
	if _, err := strconv.Atoi(postId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalide postID"})
		return
	}
	// 获取数据库中的文章信息
	post := models.Post{}
	if err := config.DB.Where("id =?", postId).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not exist"})
		return
	}
	// 校验：只有作者才能删除文章
	if userId.(uint) != post.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have no permission to delete this post"})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failded to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfuly"})
}
