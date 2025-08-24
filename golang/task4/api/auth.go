package api

import (
	"blog/config"
	"blog/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 定义 JWT 签名密钥，用于 token 的签名和验证
var JwtKey = []byte("your_secret_key")

/*
定义自定义的 JWT claims 结构体，包含：

	UserID：用户 ID
	Username：用户名
	jwt.StandardClaims：标准 JWT claims，包含过期时间等信息
*/
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 用户注册
func Register(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加密密码
	if hashPassWord, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash passwrod"})
		return
	} else {
		user.Password = string(hashPassWord)
	}

	// 创建用户
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 用户登录
func Login(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	storedUser := models.User{}
	if err := config.DB.Where("username =?", user.Username).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成JWT
	// expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		UserID:         storedUser.ID,
		Username:       storedUser.Username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(24 * time.Hour).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
