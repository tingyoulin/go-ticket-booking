package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/tingyoulin/go-ticket-booking/domain"
	"github.com/tingyoulin/go-ticket-booking/internal/utils"
)

type BookingService interface {
	Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetByIDAndPassengerID(ctx context.Context, id int64, passengerID int64) (*domain.Booking, error)
	GetListByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.Booking, error)
	Updates(ctx context.Context, bookingReq *domain.Booking) (*domain.Booking, error)
}

type BookingHandler struct {
	Service BookingService
}

func NewBookingHandler(e *echo.Echo, svc BookingService, authMiddleware echo.MiddlewareFunc) {
	handler := &BookingHandler{
		Service: svc,
	}
	e.POST("/api/bookings", handler.Create, authMiddleware)
	e.GET("/api/bookings/:id", handler.GetByID, authMiddleware)
	e.GET("/api/bookings/list", handler.GetListByPassengerID, authMiddleware)
	e.PATCH("/api/bookings", handler.Updates, authMiddleware)
}

// POST /api/bookings
func (h *BookingHandler) Create(c echo.Context) error {
	var booking domain.BookingCreateRequest
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.ResponseError{Message: err.Error()})
	}
	if ok, err := utils.IsRequestValid(&booking); !ok {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: err.Error()})
	}

	bookResponse, err := h.Service.Create(c.Request().Context(), &domain.Booking{
		PassengerID: c.Get("passenger_id").(int64),
		FlightID:    booking.FlightID,
		Seats:       booking.Seats,
		Status:      domain.BookingStatusConfirmed,
	})
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, bookResponse.ToBookingResponse())
}

// GET /api/bookings/:id
func (h *BookingHandler) GetByID(c echo.Context) error {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: "invalid booking ID"})
	}
	booking, err := h.Service.GetByIDAndPassengerID(c.Request().Context(), bookingID, c.Get("passenger_id").(int64))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, booking.ToBookingResponse())
}

// GET /api/bookings/list
func (h *BookingHandler) GetListByPassengerID(c echo.Context) error {
	passengerID := c.Get("passenger_id").(int64)
	page, perPage := utils.ParsePage(c)
	bookings, err := h.Service.GetListByPassengerID(c.Request().Context(), passengerID, page, perPage)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}
	bookingResponses := make([]domain.BookingListResponse, len(bookings))
	for i, booking := range bookings {
		bookingResponses[i] = domain.BookingListResponse{
			ID: booking.ID,
			Flight: domain.FlightSummaryResponse{
				ID:            booking.Flight.ID,
				Departure:     booking.Flight.Departure,
				Destination:   booking.Flight.Destination,
				DepartureTime: booking.Flight.DepartureTime,
				Status:        booking.Flight.Status,
			},
			CreatedAt: booking.CreatedAt,
			UpdatedAt: booking.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, bookingResponses)
}

// PATCH /api/bookings/:id
func (h *BookingHandler) Updates(c echo.Context) error {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: "invalid booking ID"})
	}

	var requestBody domain.BookingPatchRequest

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if ok, err := utils.IsRequestValid(&requestBody); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// only allow to update status and seats
	if requestBody.FlightID != 0 {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: "flight id is not allowed to be updated"})
	}

	if requestBody.Status != "" && requestBody.Status != domain.BookingStatusCanceled {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: "invalid status"})
	}

	booking, err := h.Service.Updates(c.Request().Context(), &domain.Booking{
		ID:     bookingID,
		Seats:  requestBody.Seats,
		Status: requestBody.Status,
	})
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, booking.ToBookingResponse())
}
