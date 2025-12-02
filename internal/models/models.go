package models

import "time"

type User struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Coach struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Availability struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	CoachID   uint   `gorm:"index" json:"coach_id"`
	Day       string `gorm:"column:day" json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	CoachID   uint      `gorm:"index;uniqueIndex:ux_coach_start" json:"coach_id"`
	StartTime time.Time `gorm:"uniqueIndex:ux_coach_start" json:"start_time"`
	CreatedAt time.Time `json:"created_at"`
}
