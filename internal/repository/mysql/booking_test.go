package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/tingyoulin/go-ticket-booking/domain"
	"github.com/tingyoulin/go-ticket-booking/internal/repository"
	"github.com/tingyoulin/go-ticket-booking/internal/repository/mysql"
)

func TestCreateBooking(t *testing.T) {
	gormDB, mock := repository.SetupTestDB(t)

	booking := &domain.Booking{
		PassengerID: 1,
		FlightID:    1,
		Seats:       1,
		Status:      domain.BookingStatusConfirmed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `bookings`").
		WithArgs(booking.PassengerID, booking.FlightID, booking.Seats, booking.Status, booking.CreatedAt, booking.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := mysql.NewBookingRepository(gormDB)
	booking, err := repo.Create(context.TODO(), booking)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), booking.ID)
}

func TestCreateBookingWithInvalidData(t *testing.T) {
	gormDB, mock := repository.SetupTestDB(t)

	booking := &domain.Booking{
		PassengerID: 0, // Invalid PassengerID
		FlightID:    1,
		Seats:       1,
		Status:      domain.BookingStatusConfirmed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := mysql.NewBookingRepository(gormDB)
	booking, err := repo.Create(context.TODO(), booking)
	assert.Error(t, err)
	assert.Nil(t, booking)

	mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM ` + "`bookings`" + `$`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	var count int
	err = gormDB.Raw("SELECT COUNT(*) FROM `bookings`").Scan(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "database should not have changed")
}

func TestCreateBookingWithDuplicate(t *testing.T) {
	gormDB, mock := repository.SetupTestDB(t)

	booking := &domain.Booking{
		PassengerID: 1,
		FlightID:    1,
		Seats:       1,
		Status:      domain.BookingStatusConfirmed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `bookings`").
		WithArgs(booking.PassengerID, booking.FlightID, booking.Seats, booking.Status, booking.CreatedAt, booking.UpdatedAt).
		WillReturnError(fmt.Errorf("duplicate entry"))
	mock.ExpectRollback()

	repo := mysql.NewBookingRepository(gormDB)
	booking, err := repo.Create(context.TODO(), booking)
	assert.Error(t, err)
	assert.Nil(t, booking)
}

func TestGetByID(t *testing.T) {
	gormDB, mock := repository.SetupTestDB(t)

	repo := mysql.NewBookingRepository(gormDB)
	mockBooking := domain.Booking{ID: 1, PassengerID: 123, FlightID: 456}

	mock.ExpectQuery("^SELECT \\* FROM `bookings` WHERE `bookings`.`id` = \\? ORDER BY `bookings`.`id` LIMIT \\?$").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "passenger_id", "flight_id"}).
			AddRow(mockBooking.ID, mockBooking.PassengerID, mockBooking.FlightID))

	booking, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, mockBooking.ID, booking.ID)
	assert.Equal(t, mockBooking.PassengerID, booking.PassengerID)
	assert.Equal(t, mockBooking.FlightID, booking.FlightID)
}

func TestGetListByPassengerID(t *testing.T) {
	gormDB, mock := repository.SetupTestDB(t)
	repo := mysql.NewBookingRepository(gormDB)
	mockBookings := []domain.Booking{
		{ID: 1, PassengerID: 123, FlightID: 456},
		{ID: 2, PassengerID: 123, FlightID: 789},
	}

	mock.ExpectQuery(`^SELECT \* FROM `+"`bookings`"+` WHERE passenger_id = \? LIMIT \?$`).
		WithArgs(123, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "passenger_id", "flight_id"}).
			AddRow(mockBookings[0].ID, mockBookings[0].PassengerID, mockBookings[0].FlightID).
			AddRow(mockBookings[1].ID, mockBookings[1].PassengerID, mockBookings[1].FlightID))

	bookings, err := repo.GetListByPassengerID(context.Background(), 123, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, mockBookings[0].ID, bookings[0].ID)
	assert.Equal(t, mockBookings[1].ID, bookings[1].ID)
}

func TestUpdateStatus(t *testing.T) {
	gormDB, mock := repository.SetupTestDB(t)
	repo := mysql.NewBookingRepository(gormDB)

	booking := &domain.Booking{ID: 1, PassengerID: 123, FlightID: 456}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `bookings`").
		WithArgs(booking.ID, booking.PassengerID, booking.FlightID, sqlmock.AnyArg(), booking.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	booking, err := repo.UpdateStatus(context.Background(), booking)
	assert.NoError(t, err)
	assert.Equal(t, domain.BookingStatusCanceled, booking.Status)
}
