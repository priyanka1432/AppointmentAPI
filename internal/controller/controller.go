package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	apperrors "AppointmentAPI/internal/errors"
	"AppointmentAPI/internal/models"
	"AppointmentAPI/internal/service"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *service.Service
}

func NewController(s *service.Service) *Controller {
	return &Controller{service: s}
}

type availabilityDTO struct {
	CoachID   uint   `json:"coach_id" binding:"required"`
	Day       string `json:"day" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}
type bookingDTO struct {
	UserID   uint   `json:"user_id" binding:"required"`
	CoachID  uint   `json:"coach_id" binding:"required"`
	Datetime string `json:"datetime" binding:"required"`
}


func (controller *Controller) AddAvailability(c *gin.Context) {
	var details availabilityDTO

	if err := c.ShouldBindJSON(&details); err != nil {
		c.Error(apperrors.NewBadRequest(err.Error()))
		return
	}

	a := &models.Availability{
		CoachID:   details.CoachID,
		Day:       details.Day,
		StartTime: details.StartTime,
		EndTime:   details.EndTime,
	}
	if appErr := controller.service.AddAvailability(a); appErr != nil {
		c.Error(appErr)
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (ctl *Controller) GetAvailableSlots(c *gin.Context) {
	coachStr := c.Query("coach_id")
	dateStr := c.Query("date")

	fmt.Println("DATE RECEIVED:", dateStr)
	if coachStr == "" || dateStr == "" {
		c.Error(apperrors.NewBadRequest("coach_id and date required"))
		return
	}
	coachID, err := strconv.Atoi(coachStr)
	if err != nil {
		c.Error(apperrors.NewBadRequest("invalid coach_id"))
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.Error(apperrors.NewBadRequest("invalid date, use YYYY-MM-DD"))
		return
	}
	slots, appErr := ctl.service.GetAvailableSlots(uint(coachID), date)
	if appErr != nil {
		c.Error(appErr)
		return
	}
	out := make([]string, 0, len(slots))
	for _, t := range slots {
		out = append(out, t.UTC().Format(time.RFC3339))
	}
	c.JSON(http.StatusOK, gin.H{"slots": out})
}

func (ctl *Controller) BookSlot(c *gin.Context) {
	var details bookingDTO
	if err := c.ShouldBindJSON(&details); err != nil {
		c.Error(apperrors.NewBadRequest(err.Error()))
		return
	}
	dt, err := time.Parse(time.RFC3339, details.Datetime)
	if err != nil {
		c.Error(apperrors.NewBadRequest("invalid datetime; use RFC3339"))
		return
	}
	booking, appErr := ctl.service.BookSlot(details.UserID, details.CoachID, dt)
	if appErr != nil {
		c.Error(appErr)
		return
	}
	c.JSON(http.StatusCreated, booking)
}

func (ctl *Controller) GetUserBookings(c *gin.Context) {
	uidStr := c.Query("user_id")
	if uidStr == "" {
		c.Error(apperrors.NewBadRequest("user_id required"))
		return
	}
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.Error(apperrors.NewBadRequest("invalid user_id"))
		return
	}
	bookings, appErr := ctl.service.GetUserBookings(uint(uid))
	if appErr != nil {
		c.Error(appErr)
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (ctl *Controller) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	userIdStr := c.Param("userid")
	if idStr == "" {
		c.Error(apperrors.NewBadRequest("booking id required"))
		return
	}
	if userIdStr == "" {
		c.Error(apperrors.NewBadRequest("userId id required"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(apperrors.NewBadRequest("invalid booking id"))
		return
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.Error(apperrors.NewBadRequest("invalid booking id"))
		return
	}

	if appErr := ctl.service.CancelBooking(uint(id), uint(userID)); appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cancelled", "id": id})
}
