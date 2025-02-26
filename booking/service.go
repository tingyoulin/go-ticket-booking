package booking

import (
	"context"
	"fmt"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	GetByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.BookingListResponse, error)
	UpdateStatus(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	UpdateSeats(ctx context.Context, booking *domain.Booking, seatsChange int32) (*domain.Booking, error)
}

type FlightRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Flight, error)
}

type Service struct {
	bookingRepo BookingRepository
	flightRepo  FlightRepository
}

func NewService(b BookingRepository, f FlightRepository) *Service {
	return &Service{
		bookingRepo: b,
		flightRepo:  f,
	}
}

func (s *Service) Create(ctx context.Context, booking *domain.BookingCreateRequest) (*domain.Booking, error) {
	// check if flight has available seats
	flight, err := s.flightRepo.GetByID(ctx, booking.FlightID)
	if err != nil {
		return nil, err
	}

	// check if flight is canceled
	if flight.Status == domain.FlightStatusCanceled {
		return nil, fmt.Errorf("flight %d is canceled", booking.FlightID)
	}

	if flight.AvailableSeats < booking.Seats {
		return nil, fmt.Errorf("flight %d has only %d available seats", booking.FlightID, flight.AvailableSeats)
	}

	// create booking
	return s.bookingRepo.Create(ctx, &domain.Booking{
		PassengerID: booking.PassengerID,
		FlightID:    booking.FlightID,
		Flight:      *flight,
		Seats:       booking.Seats,
		Status:      domain.BookingStatusConfirmed,
	})
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

func (s *Service) GetByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.BookingListResponse, error) {
	return s.bookingRepo.GetByPassengerID(ctx, passengerID, page, perPage)
}

func (s *Service) Updates(ctx context.Context, bookingReq *domain.BookingPatchRequest) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingReq.BookingID)
	if err != nil {
		return nil, err
	}

	if bookingReq.FlightID != 0 {
		return nil, fmt.Errorf("flight id is not allowed to be updated")
	}

	if bookingReq.Status == domain.BookingStatusCanceled {
		booking.Status = bookingReq.Status
		return s.bookingRepo.UpdateStatus(ctx, booking)
	}

	if bookingReq.Seats != 0 {
		seatsChange := bookingReq.Seats - booking.Seats
		booking.Seats = bookingReq.Seats
		return s.bookingRepo.UpdateSeats(ctx, booking, seatsChange)
	}

	return booking, nil
}
