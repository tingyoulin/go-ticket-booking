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
	Create(ctx context.Context, booking *domain.BookingCreateRequest) (*domain.Booking, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	GetByPassengerID(ctx context.Context, passengerID int64, page, perPage int) ([]domain.BookingListResponse, error)
	Updates(ctx context.Context, bookingReq *domain.BookingPatchRequest) (*domain.Booking, error)
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
	e.GET("/api/bookings/list", handler.GetByPassengerID, authMiddleware)
	e.PATCH("/api/bookings", handler.Updates, authMiddleware)
}

// POST /api/bookings
func (h *BookingHandler) Create(c echo.Context) error {
	var booking domain.BookingCreateRequest
	if err := c.Bind(&booking); err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}
	if ok, err := utils.IsRequestValid(&booking); !ok {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	bookResponse, err := h.Service.Create(ctx, &booking)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, bookResponse)
}

// GET /api/bookings/:id
func (h *BookingHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: "invalid booking ID"})
	}
	booking, err := h.Service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, booking)
}

// GET /api/bookings/list
func (h *BookingHandler) GetByPassengerID(c echo.Context) error {
	passengerID, err := strconv.ParseInt(c.QueryParam("passenger_id"), 10, 64)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: "invalid passenger ID"})
	}
	page, perPage := utils.ParsePage(c)
	bookings, err := h.Service.GetByPassengerID(c.Request().Context(), passengerID, page, perPage)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, bookings)
}

// PATCH /api/bookings
func (h *BookingHandler) Updates(c echo.Context) error {
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

	booking, err := h.Service.Updates(c.Request().Context(), &requestBody)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, booking)
}
