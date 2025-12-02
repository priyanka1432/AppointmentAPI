package repo

import (
	"AppointmentAPI/internal/models"
	"time"
)

func (r *gormRepo) AddAvailability(a *models.Availability) error {
	return r.db.Create(a).Error
}

func (r *gormRepo) GetAvailabilities(coachID uint) ([]models.Availability, error) {
	var arr []models.Availability
	err := r.db.Where("coach_id = ?", coachID).Find(&arr).Error
	return arr, err
}

func (r *gormRepo) BookSlot(b *models.Booking) error {
	return r.db.Create(b).Error
}

func (r *gormRepo) GetBookingsByCoachDay(coachID uint, dayStart, dayEnd time.Time) ([]models.Booking, error) {
	var arr []models.Booking
	err := r.db.Where("coach_id = ? AND start_time >= ? AND start_time < ?", coachID, dayStart, dayEnd).Find(&arr).Error
	return arr, err
}

func (r *gormRepo) GetBookingsByUser(userID uint) ([]models.Booking, error) {
	var arr []models.Booking
	err := r.db.Where("user_id = ? AND start_time >= ?", userID, time.Now().UTC()).Order("start_time asc").Find(&arr).Error
	return arr, err
}

func (r *gormRepo) DeleteBooking(bookingID uint) error {
	return r.db.Delete(&models.Booking{}, bookingID).Error
}
func (r *gormRepo) GetBookingByID(bookingID uint) (*models.Booking, error) {
	var b models.Booking
	if err := r.db.First(&b, bookingID).Error; err != nil {
		return nil, err
	}
	return &b, nil
}
