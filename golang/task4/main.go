package main

import (
	"blog/config"
	"blog/models"
	"blog/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	config.InitDB()

	// 自动迁移模型
	config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	// 设置路由
	router := gin.Default()

	// 强制设置 Gin 使用 UTF-8 编码
	router.Use(func(c *gin.Context) {
		// 确保响应以UTF-8编码
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	routes.SetupRoutes(router)

	// 启动服务器
	router.Run(":8080")
}
