package handler

import (
	"github.com/eydeveloper/highload-social/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("api")
	{
		auth := api.Group("auth")
		{
			auth.POST("login", h.login)
			auth.POST("register", h.register)
		}

		user := api.Group("user")
		{
			user.GET(":id", h.getUserById)
			user.GET("search", h.searchUsers)
		}

		api.PUT("follow/:id", h.authenticationMiddleware(), h.follow)
		api.PUT("unfollow/:id", h.authenticationMiddleware(), h.unfollow)
	}

	return router
}

func (h *Handler) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := tokenParts[1]
		userId, err := h.services.Authorization.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token"})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
