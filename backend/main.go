package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	initDB()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", health)

	r.OPTIONS("/api/quizzes", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.GET("/api/quizzes", listQuizzes)
	r.POST("/api/quizzes", createQuiz)

	r.OPTIONS("/api/quizzes/:id", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.GET("/api/quizzes/:id", getQuiz)

	// Auth routes
	r.OPTIONS("/api/register", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.POST("/api/register", registerUser)
	r.OPTIONS("/api/login", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.POST("/api/login", loginUser)
	r.OPTIONS("/api/me", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.GET("/api/me", getMe)

	// Submission/results routes
	r.OPTIONS("/api/quizzes/:id/submit", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.POST("/api/quizzes/:id/submit", submitQuiz)
	r.OPTIONS("/api/quizzes/:id/results/:submissionId", func(c *gin.Context) { enableCorsHeaders(c); c.Status(200) })
	r.GET("/api/quizzes/:id/results/:submissionId", getResults)

	log.Println("starting backend on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
