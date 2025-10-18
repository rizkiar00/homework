package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rizkiar00/homework/internal/controller"
	"github.com/rizkiar00/homework/internal/repository"
	"github.com/rizkiar00/homework/internal/usecase"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// wire hexagonal components: repository -> usecase -> controller
	repo := repository.NewInMemoryRepo()
	uc := usecase.NewEchoUsecase(repo)
	ctrl := controller.NewEchoController(uc)

	// routes
	r.GET("/echo/:value", ctrl.GetEcho)
	r.POST("/echo", ctrl.PostEcho)

	r.Run(":8080")
}
