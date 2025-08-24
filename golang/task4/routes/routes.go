package routes

import (
	"blog/api"
	"blog/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 用户认证
	auth := router.Group("/auth")
	{
		auth.POST("/register", api.Register)
		auth.POST("/login", api.Login)
	}

	// 用户路由（需要认证）
	user := router.Group("/users")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/:id", api.GetUser)
	}

	// 文章路由（需要认证）
	post := router.Group("/post")
	post.Use(middleware.AuthMiddleware())
	{
		// 获取所有文章
		post.GET("/", api.GetPosts)
		// 获取单个文章
		post.GET("/:id", api.GetPost)
		// 创建文章
		post.POST("/", api.CreatePost)
		// 更新文章
		post.PUT("/", api.UpdatePost)
		// 删除文章
		post.DELETE("/:id", api.DeletePost)
	}

	// 评论路由
	comment := router.Group("/comments")
	{
		// 创建评论
		comment.POST("/", middleware.AuthMiddleware(), api.CreateComment)
		// 根据文章ID查询所有评论信息
		comment.GET("/post/:postId", api.GetCommentsByPostId)
	}

}
