package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

type FlightRepository struct {
	DB *gorm.DB
}

func NewFlightRepository(db *gorm.DB) *FlightRepository {
	return &FlightRepository{DB: db}
}

func (r *FlightRepository) Fetch(ctx context.Context, departure, destination string, departureTime time.Time, seatsToBook, page, perPage int) ([]domain.Flight, error) {
	var flights []domain.Flight
	err := r.DB.WithContext(ctx).
		Table("flights").
		Where("departure = ? AND destination = ? AND departure_time >= ? AND available_seats >= ?",
			departure, destination, departureTime, seatsToBook).
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&flights).Error
	if err != nil {
		return nil, err
	}
	return flights, nil
}

func (r *FlightRepository) GetByID(ctx context.Context, id int64) (*domain.Flight, error) {
	var flight domain.Flight
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&flight).Error; err != nil {
		return nil, err
	}
	return &flight, nil
}

func (r *FlightRepository) IncrementAvailableSeats(ctx context.Context, flight *domain.Flight, seats int32) (*domain.Flight, error) {
	flight.AvailableSeats += seats

	if err := r.DB.WithContext(ctx).
		Model(&domain.Flight{}).
		Where("id = ? AND lock_version = ?", flight.ID, flight.LockVersion).
		Updates(map[string]interface{}{
			"available_seats": flight.AvailableSeats,
			"lock_version":    flight.LockVersion + 1,
		}).Error; err != nil {
		return nil, err
	}
	return flight, nil
}

func (r *FlightRepository) DecrementAvailableSeats(ctx context.Context, flight *domain.Flight, seats int32) (*domain.Flight, error) {
	if err := r.DB.WithContext(ctx).
		Model(&domain.Flight{}).
		Where("id = ? AND available_seats >= ? AND lock_version = ?", flight.ID, seats, flight.LockVersion).
		Updates(map[string]interface{}{
			"available_seats": gorm.Expr("available_seats - ?", seats),
			"lock_version":    gorm.Expr("lock_version + 1"),
		}).Error; err != nil {
		return nil, err
	}
	flight.AvailableSeats -= seats

	return flight, nil
}
