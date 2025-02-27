package booking

import (
	"context"
	"fmt"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	GetByIDAndPassengerID(ctx context.Context, id int64, passengerID int64) (*domain.Booking, error)
	GetListByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.Booking, error)
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

func (s *Service) Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	// check if flight has available seats
	flight, err := s.flightRepo.GetByID(ctx, booking.FlightID)
	if err != nil {
		return nil, err
	}

	// check if flight is canceled
	if flight.Status == domain.FlightStatusCanceled {
		return nil, domain.ErrFlightCanceled
	}

	if flight.AvailableSeats < booking.Seats {
		return nil, domain.ErrFlightNoAvailableSeats
	}

	// create booking
	booking, err = s.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, booking.ID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

func (s *Service) GetByIDAndPassengerID(ctx context.Context, bookingID int64, passengerID int64) (*domain.Booking, error) {
	return s.bookingRepo.GetByIDAndPassengerID(ctx, bookingID, passengerID)
}

func (s *Service) GetListByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.Booking, error) {
	return s.bookingRepo.GetListByPassengerID(ctx, passengerID, page, perPage)
}

func (s *Service) Updates(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	bookingToUpdate, err := s.bookingRepo.GetByID(ctx, booking.ID)
	if err != nil {
		return nil, err
	}

	if booking.FlightID != 0 {
		return nil, fmt.Errorf("flight id is not allowed to be updated")
	}

	if booking.Status == domain.BookingStatusCanceled {
		bookingToUpdate.Status = booking.Status
		return s.bookingRepo.UpdateStatus(ctx, booking)
	}

	if booking.Seats != 0 {
		seatsChange := booking.Seats - bookingToUpdate.Seats
		bookingToUpdate.Seats = booking.Seats
		return s.bookingRepo.UpdateSeats(ctx, bookingToUpdate, seatsChange)
	}

	return bookingToUpdate, nil
}
