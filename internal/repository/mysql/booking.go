package mysql

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

type BookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	tx := r.DB.WithContext(ctx).Begin()
	err := tx.Model(&domain.Flight{}).
		Where("id = ? AND available_seats >= ? AND lock_version = ?", booking.FlightID, booking.Seats, booking.Flight.LockVersion).
		Updates(map[string]interface{}{
			"available_seats": gorm.Expr("available_seats - ?", booking.Seats),
			"lock_version":    gorm.Expr("lock_version + 1"),
		}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Omit(clause.Associations).Create(booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return booking, tx.Commit().Error
}

func (r *BookingRepository) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.DB.WithContext(ctx).
		Joins("Passenger").
		Joins("Flight").
		First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) GetByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.BookingListResponse, error) {
	var bookings []domain.BookingListResponse
	if err := r.DB.WithContext(ctx).
		Model(&domain.Booking{}).
		Joins("Flight").
		Where("bookings.passenger_id = ?", passengerID).
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if err := tx.Model(booking).Update("status", booking.Status).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&domain.Flight{}).
		Where("id = ? AND lock_version = ?", booking.FlightID, booking.Flight.LockVersion).
		Updates(map[string]interface{}{
			"available_seats": gorm.Expr("available_seats + ?", booking.Seats),
			"lock_version":    gorm.Expr("lock_version + 1"),
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return booking, tx.Commit().Error
}

func (r *BookingRepository) UpdateSeats(ctx context.Context, booking *domain.Booking, seatsChange int32) (*domain.Booking, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if err := tx.Model(booking).Update("seats", booking.Seats).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	booking.Flight.AvailableSeats -= seatsChange

	if err := tx.Model(&domain.Flight{}).
		Where("id = ? AND available_seats >= ? AND lock_version = ?", booking.FlightID, seatsChange, booking.Flight.LockVersion).
		Updates(map[string]interface{}{
			"available_seats": booking.Flight.AvailableSeats,
			"lock_version":    gorm.Expr("lock_version + 1"),
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return booking, tx.Commit().Error
}
