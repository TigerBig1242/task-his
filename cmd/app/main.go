package main

import (
	"log"

	"github.com/bigthamm/task-his/internal/config"
	"github.com/bigthamm/task-his/internal/database/postgres"
	model "github.com/bigthamm/task-his/internal/model"
	"github.com/gin-gonic/gin"
)

func main() {

	databaseConfig := config.LoadConfig()

	db, err := postgres.ConnectDatabase(databaseConfig)

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.Hospital{}, &model.Patient{}, &model.Staff{})
	if err != nil {
		log.Fatalf("Database Migration failed: %v", err)
	}
	router := gin.Default()
	// router.SetTrustedProxies([]string{"192.168.1.222"})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		// fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})
	router.Run()
}
