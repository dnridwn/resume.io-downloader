package main

import (
	"log"
	"net/http"
	resumeio "resume-downloader/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"error": false, "message": "Welcome to the resume downloader service"})
	})

	r.GET("/generate", func(ctx *gin.Context) {
		renderingToken := ctx.Query("rendering_token")
		if renderingToken == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Rendering token is required"})
			return
		}

		pdfBytes, err := resumeio.GeneratePDF(renderingToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
			return
		}

		ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
