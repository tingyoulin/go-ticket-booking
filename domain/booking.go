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

func (b Booking) ToBookingResponse() BookingResponse {
	return BookingResponse{
		ID: b.ID,
		Flight: FlightSummaryResponse{
			ID:            b.Flight.ID,
			Departure:     b.Flight.Departure,
			Destination:   b.Flight.Destination,
			DepartureTime: b.Flight.DepartureTime,
			Status:        b.Flight.Status,
		},
		Passenger: PassengerResponse{
			ID:        b.Passenger.ID,
			Name:      b.Passenger.Name,
			Email:     b.Passenger.Email,
			CreatedAt: b.Passenger.CreatedAt,
			UpdatedAt: b.Passenger.UpdatedAt,
		},
		Seats:     b.Seats,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

type BookingCreateRequest struct {
	FlightID int64 `json:"flight_id" validate:"required"`
	Seats    int32 `json:"seats" validate:"required"`
}

type BookingPatchRequest struct {
	FlightID int64         `json:"flight_id"`
	Seats    int32         `json:"seats"`
	Status   BookingStatus `json:"status"`
}

type BookingResponse struct {
	ID        int64                 `json:"id"`
	Flight    FlightSummaryResponse `json:"flight"`
	Passenger PassengerResponse     `json:"passenger"`
	Seats     int32                 `json:"seats"`
	Status    BookingStatus         `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type BookingListResponse struct {
	ID        int64                 `json:"id"`
	Flight    FlightSummaryResponse `json:"flight"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
