package domain

import (
	"time"
)

type Booking struct {
	ID          int64         `json:"id"`
	Passenger   Passenger     `gorm:"foreignKey:PassengerID"`
	PassengerID int64         `json:"passenger_id"`
	Flight      Flight        `gorm:"foreignKey:FlightID"`
	FlightID    int64         `json:"flight_id"`
	Seats       int32         `json:"seats"`
	Status      BookingStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCanceled  BookingStatus = "canceled"
)

type BookingCreateRequest struct {
	PassengerID int64         `json:"passenger_id" validate:"required"`
	FlightID    int64         `json:"flight_id" validate:"required"`
	Seats       int32         `json:"seats" validate:"required"`
	Status      BookingStatus `json:"status"`
}

type BookingPatchRequest struct {
	BookingID int64         `json:"booking_id" validate:"required"`
	FlightID  int64         `json:"flight_id"`
	Seats     int32         `json:"seats"`
	Status    BookingStatus `json:"status"`
}

type BookingListResponse struct {
	ID        int64                 `json:"id"`
	Flight    FlightSummaryResponse `json:"flight" gorm:"foreignKey:flight_id"`
	FlightID  int64                 `json:"flight_id"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
