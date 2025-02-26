package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/tingyoulin/go-ticket-booking/domain"
	"github.com/tingyoulin/go-ticket-booking/internal/utils"
)

type FlightService interface {
	Fetch(ctx context.Context, departure, destination string, departureTime time.Time, seatsToBook, page, perPage int) ([]domain.Flight, error)
}

type FlightHandler struct {
	Service FlightService
}

func NewFlightHandler(e *echo.Echo, svc FlightService) {
	handler := &FlightHandler{
		Service: svc,
	}
	e.GET("/api/flights", handler.Fetch)
}

// GET /api/flights
func (handler *FlightHandler) Fetch(c echo.Context) error {
	departure := c.QueryParam("departure")
	destination := c.QueryParam("destination")
	if departure == "" || destination == "" {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: "departure and destination are required"})
	}
	departureTime, err := time.Parse(time.RFC3339, c.QueryParam("departure_time"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: fmt.Sprintf("invalid departure time: %s", err.Error())})
	}
	seats, err := strconv.Atoi(c.QueryParam("seats"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: fmt.Sprintf("invalid seats: %s", err.Error())})
	}
	page, perPage := utils.ParsePage(c)

	ctx := c.Request().Context()

	flights, err := handler.Service.Fetch(ctx, departure, destination, departureTime, seats, page, perPage)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, flights)
}
