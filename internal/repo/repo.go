package repo

import (
	"AppointmentAPI/internal/models"
	"time"

	"gorm.io/gorm"
)

type Repo interface {
	AddAvailability(a *models.Availability) error
	GetAvailabilities(coachID uint) ([]models.Availability, error)
	BookSlot(b *models.Booking) error
	GetBookingsByCoachDay(coachID uint, dayStart, dayEnd time.Time) ([]models.Booking, error)
	GetBookingsByUser(userID uint) ([]models.Booking, error)
	DeleteBooking(bookingID uint) error
	GetBookingByID(bookingID uint) (*models.Booking, error)
}

type gormRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &gormRepo{db: db}
}
