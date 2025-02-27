package domain

import (
	"time"
)

type Flight struct {
	ID             int64        `json:"id"`
	Departure      string       `json:"departure" validate:"required"`
	Destination    string       `json:"destination" validate:"required"`
	DepartureTime  time.Time    `json:"departure_time" gorm:"type:datetime" validate:"required"`
	AvailableSeats int32        `json:"available_seats" validate:"required"`
	Price          float32      `json:"price" validate:"required"`
	Status         FlightStatus `json:"status" validate:"required"`
	LockVersion    int32        `json:"lock_version" validate:"required"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type FlightStatus string

const (
	FlightStatusOnTime   FlightStatus = "on_time"
	FlightStatusDelayed  FlightStatus = "delayed"
	FlightStatusCanceled FlightStatus = "canceled"
)

type FlightSummaryResponse struct {
	ID            int64        `json:"id"`
	Departure     string       `json:"departure"`
	Destination   string       `json:"destination"`
	DepartureTime time.Time    `json:"departure_time"`
	Status        FlightStatus `json:"status"`
}

type FlightResponse struct {
	ID            int64        `json:"id"`
	Departure     string       `json:"departure"`
	Destination   string       `json:"destination"`
	DepartureTime time.Time    `json:"departure_time"`
	Price         float32      `json:"price"`
	Status        FlightStatus `json:"status"`
}
