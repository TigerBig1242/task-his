package main

import (
	"log"

	"github.com/bigthamm/task-his/internal/config"
	"github.com/bigthamm/task-his/internal/database/postgres"
	"github.com/bigthamm/task-his/internal/delivery/handler"
	router "github.com/bigthamm/task-his/internal/delivery/router"
	model "github.com/bigthamm/task-his/internal/model"
	"github.com/bigthamm/task-his/internal/repository"
	"github.com/bigthamm/task-his/internal/service"
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
	app := gin.Default()
	// router.SetTrustedProxies([]string{"192.168.1.222"})

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		// fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})

	// Wiring api create staff
	createStaffRepo := repository.NewStaffRepository(db)
	createStaffService := service.NewStaffService(createStaffRepo)
	staffHandler := handler.NewStaffHandler(createStaffService)
	router.SetStaffRoute(app, staffHandler)

	app.Run()
}
