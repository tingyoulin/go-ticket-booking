package flight

import (
	"context"
	"time"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

type FlightRepository interface {
	Fetch(ctx context.Context, departure, destination string, departureTime time.Time, seatsToBook, page, perPage int) ([]domain.Flight, error)
}

type Service struct {
	flightRepo FlightRepository
}

func NewService(f FlightRepository) *Service {
	return &Service{
		flightRepo: f,
	}
}

func (s *Service) Fetch(ctx context.Context, departure, destination string, departureTime time.Time, seatsToBook, page, perPage int) ([]domain.Flight, error) {
	flights, err := s.flightRepo.Fetch(ctx, departure, destination, departureTime, seatsToBook, page, perPage)
	if err != nil {
		return nil, err
	}
	return flights, nil
}
