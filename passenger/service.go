package passenger

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

const (
	tokenExpired = 24 * time.Hour
)

type PassengerRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.Passenger, error)
	Create(ctx context.Context, passenger *domain.Passenger) (*domain.Passenger, error)
}

type TokenRedisRepository interface {
	Set(ctx context.Context, token string, expiration time.Duration) error
}

type Service struct {
	passengerRepo  PassengerRepository
	tokenRedisRepo TokenRedisRepository
}

func NewService(passengerRepo PassengerRepository, tokenRedisRepo TokenRedisRepository) *Service {
	return &Service{
		passengerRepo:  passengerRepo,
		tokenRedisRepo: tokenRedisRepo,
	}
}

func (s *Service) Register(ctx context.Context, passenger *domain.Passenger) (*domain.Passenger, error) {
	// check if passenger already exists
	if existingPassenger, err := s.passengerRepo.GetByEmail(ctx, passenger.Email); err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if existingPassenger != nil {
		return nil, domain.ErrConflict
	}

	// using bcrypt to hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passenger.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passenger.Password = string(hashedPassword)

	// create passenger
	return s.passengerRepo.Create(ctx, passenger)
}

func (s *Service) Login(ctx context.Context, passenger *domain.Passenger) (*domain.Token, error) {
	// check if passenger exists
	existingPassenger, err := s.passengerRepo.GetByEmail(ctx, passenger.Email)
	if err != nil || existingPassenger == nil {
		return nil, domain.ErrNotFound
	}

	// check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(existingPassenger.Password), []byte(passenger.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passenger_id": existingPassenger.ID,
		"email":        existingPassenger.Email,
		"exp":          time.Now().Add(tokenExpired).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &domain.Token{AccessToken: tokenString}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	// store token to redis as blacklist
	if err := s.tokenRedisRepo.Set(ctx, token, tokenExpired); err != nil {
		return err
	}

	return nil
}
