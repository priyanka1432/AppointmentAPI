package service

import (
	"errors"
	"strings"
	"time"

	apperrors "AppointmentAPI/internal/errors"
	"AppointmentAPI/internal/models"
	"AppointmentAPI/internal/repo"
	"AppointmentAPI/internal/utils"

	"github.com/go-sql-driver/mysql"
)

type Service struct {
	repo repo.Repo
}

func NewService(r repo.Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) AddAvailability(a *models.Availability) *apperrors.AppError {
	if a == nil {
		return apperrors.NewBadRequest("invalid payload")
	}
	if a.CoachID == 0 {
		return apperrors.NewBadRequest("coach_id required")
	}
	dayStr := strings.TrimSpace(a.Day)
	if dayStr == "" {
		return apperrors.NewBadRequest("day is required")
	}
	if _, err := utils.DayStringToInt(dayStr); err != nil {
		return apperrors.NewBadRequest("invalid day; use Monday..Sunday")
	}

	a.Day = strings.Title(strings.ToLower(dayStr))
	a.StartTime = strings.TrimSpace(a.StartTime)
	a.EndTime = strings.TrimSpace(a.EndTime)
	starthour, startminute, err := utils.ParseTimeOfDay(a.StartTime)
	if err != nil {
		return apperrors.NewBadRequest("start_time must be HH:MM")
	}
	endhour, endminute, err := utils.ParseTimeOfDay(a.EndTime)
	if err != nil {
		return apperrors.NewBadRequest("end_time must be HH:MM")
	}
	if starthour*60+startminute >= endhour*60+endminute {
		return apperrors.NewBadRequest("start_time must be before end_time")
	}
	if endhour*60+endminute-starthour*60-startminute < 30 {
		return apperrors.NewBadRequest("availability window must be >= 30 minutes")
	}

	if err := s.repo.AddAvailability(a); err != nil {
		return apperrors.NewInternal("failed to save availability")
	}
	return nil
}

func (s *Service) GetAvailableSlots(coachID uint, date time.Time) ([]time.Time, *apperrors.AppError) {

	wins, err := s.repo.GetAvailabilities(coachID)
	if err != nil {
		return nil, apperrors.NewInternal("failed to get availability")
	}

	var result []time.Time
	for _, w := range wins {

		applies, err := utils.AvailabilityAppliesToDate(w.Day, date)
		if err != nil {

			continue
		}
		if !applies {
			continue
		}
		starthour, startminute, err := utils.ParseTimeOfDay(w.StartTime)
		if err != nil {
			continue
		}
		endhour, endminute, err := utils.ParseTimeOfDay(w.EndTime)
		if err != nil {
			continue
		}

		start := utils.Combine(date, starthour, startminute)
		end := utils.Combine(date, endhour, endminute)

		cur := start
		for cur.Add(30 * time.Minute).Before(end.Add(time.Second)) {
			result = append(result, cur)
			cur = cur.Add(30 * time.Minute)
		}
	}

	return result, nil
}

func (s *Service) BookSlot(userID, coachID uint, t time.Time) (*models.Booking, *apperrors.AppError) {
	if userID == 0 || coachID == 0 {
		return nil, apperrors.NewBadRequest("user_id and coach_id required")
	}
	slot := t.UTC()
	if !utils.Is30(slot) {
		return nil, apperrors.NewBadRequest("datetime must be :00 or :30 minutes")
	}
	ok, err := s.isSlotInAvailability(coachID, slot)
	if err != nil {
		return nil, apperrors.NewInternal("failed to validate availability")
	}
	if !ok {
		return nil, apperrors.NewBadRequest("slot not within coach availability")
	}
	b := &models.Booking{UserID: userID, CoachID: coachID, StartTime: slot}
	if err := s.repo.BookSlot(b); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, apperrors.NewConflict("slot already booked")
		}
		return nil, apperrors.NewInternal("failed to create booking")
	}
	return b, nil
}

func (s *Service) isSlotInAvailability(coachID uint, slot time.Time) (bool, error) {
	slot = slot.UTC()

	wins, err := s.repo.GetAvailabilities(coachID)
	if err != nil {
		return false, err
	}

	for _, w := range wins {

		applies, err := utils.AvailabilityAppliesToDate(w.Day, slot)
		if err != nil {

			continue
		}
		if !applies {
			continue
		}

		sh, sm, err := utils.ParseTimeOfDay(w.StartTime)
		if err != nil {
			continue
		}
		eh, em, err := utils.ParseTimeOfDay(w.EndTime)
		if err != nil {
			continue
		}

		start := utils.Combine(slot, sh, sm).UTC()
		end := utils.Combine(slot, eh, em).UTC()
		slotEnd := slot.Add(30 * time.Minute)

		if (slot.Equal(start) || slot.After(start)) && (slotEnd.Equal(end) || slotEnd.Before(end)) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) GetUserBookings(userID uint) ([]models.Booking, *apperrors.AppError) {
	if userID == 0 {
		return nil, apperrors.NewBadRequest("user_id required")
	}
	bookings, err := s.repo.GetBookingsByUser(userID)
	if err != nil {
		return nil, apperrors.NewInternal("failed to fetch bookings")
	}
	return bookings, nil
}

func (s *Service) CancelBooking(bookingID, userID uint) *apperrors.AppError {
	if bookingID == 0 {
		return apperrors.NewBadRequest("booking_id required")
	}
	b, err := s.repo.GetBookingByID(bookingID)
	if err != nil {
		return apperrors.NewInternal("failed to fetch booking")
	}
	if b.UserID != userID {
		return apperrors.NewBadRequest("not allowed to cancel this booking")
	}
	if err := s.repo.DeleteBooking(bookingID); err != nil {
		return apperrors.NewInternal("failed to delete booking")
	}
	return nil
}
