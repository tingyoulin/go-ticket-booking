package mysql

import (
	"context"

	"github.com/tingyoulin/go-ticket-booking/domain"
	"gorm.io/gorm"
)

type PassengerRepository struct {
	DB *gorm.DB
}

func NewPassengerRepository(db *gorm.DB) *PassengerRepository {
	return &PassengerRepository{DB: db}
}

func (r *PassengerRepository) GetByEmail(ctx context.Context, email string) (*domain.Passenger, error) {
	var passenger domain.Passenger
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&passenger).Error; err != nil {
		return nil, err
	}
	return &passenger, nil
}

func (r *PassengerRepository) Create(ctx context.Context, passenger *domain.Passenger) (*domain.Passenger, error) {
	if err := r.DB.WithContext(ctx).Create(passenger).Error; err != nil {
		return nil, err
	}
	return passenger, nil
}
