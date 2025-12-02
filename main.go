package main

import (
	"AppointmentAPI/internal/controller"
	"AppointmentAPI/internal/db"
	"AppointmentAPI/internal/middleware"
	"AppointmentAPI/internal/repo"
	"AppointmentAPI/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	database := db.InitDb()
	repository := repo.NewRepo(database)
	service := service.NewService(repository)
	controller := controller.NewController(service)
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api")
	{
		api.POST("/coaches/availability", controller.AddAvailability)
		api.GET("/users/slots", controller.GetAvailableSlots)
		api.POST("/users/bookings", controller.BookSlot)
		api.GET("/users/bookings", controller.GetUserBookings)
		api.DELETE("/users/bookings/:userid/:id", controller.CancelBooking)
	}

	router.Run(":8080")
}
