package main

import (
	"github.com/TRu-S3/backend/internal/application/usecase"
	"github.com/TRu-S3/backend/internal/infrastructure/database"
	"github.com/TRu-S3/backend/internal/infrastructure/web"
	"github.com/gin-gonic/gin"
)

func main() {
	// 依存性の注入
	userRepo := database.NewUserRepositoryImpl()
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := web.NewUserHandler(userUseCase)

	// Ginルーターの設定
	r := gin.Default()

	// ルートの設定
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// サーバーの起動
	r.Run(":8080")
}
