package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/tingyoulin/go-ticket-booking/domain"
	"github.com/tingyoulin/go-ticket-booking/internal/utils"
)

type PassengerService interface {
	Register(ctx context.Context, passenger *domain.Passenger) (*domain.Passenger, error)
	Login(ctx context.Context, passenger *domain.Passenger) (*domain.Token, error)
	Logout(ctx context.Context, token string) error
}

type PassengerHandler struct {
	Service PassengerService
}

func NewPassengerHandler(e *echo.Echo, svc PassengerService, authMiddleware echo.MiddlewareFunc) {
	handler := &PassengerHandler{
		Service: svc,
	}
	e.POST("/api/passengers/register", handler.Register)
	e.POST("/api/passengers/login", handler.Login)
	e.POST("/api/passengers/logout", handler.Logout, authMiddleware)
}

// POST /api/passengers/register
func (h *PassengerHandler) Register(c echo.Context) error {
	var passenger domain.PassengerRegisterRequest
	if err := c.Bind(&passenger); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.ResponseError{Message: err.Error()})
	}

	if ok, err := utils.IsRequestValid(passenger); !ok {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: err.Error()})
	}

	passengerResponse, err := h.Service.Register(c.Request().Context(), &domain.Passenger{
		Name:     passenger.Name,
		Email:    passenger.Email,
		Password: passenger.Password,
	})
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, &domain.PassengerResponse{
		ID:        passengerResponse.ID,
		Name:      passengerResponse.Name,
		Email:     passengerResponse.Email,
		CreatedAt: passengerResponse.CreatedAt,
		UpdatedAt: passengerResponse.UpdatedAt,
	})
}

// POST /api/passengers/login
func (h *PassengerHandler) Login(c echo.Context) error {
	var loginRequest domain.PassengerLoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.ResponseError{Message: err.Error()})
	}

	if ok, err := utils.IsRequestValid(loginRequest); !ok {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: err.Error()})
	}

	// login and return JWT token
	token, err := h.Service.Login(c.Request().Context(), &domain.Passenger{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, token)
}

// POST /api/passengers/logout
func (h *PassengerHandler) Logout(c echo.Context) error {
	// get passenger id from context
	passengerID := c.Get("passenger_id")
	if passengerID == nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	// logout
	if err := h.Service.Logout(c.Request().Context(), c.Request().Header.Get("Authorization")); err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "logout success")
}
