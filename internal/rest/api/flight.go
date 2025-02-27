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

	// Parse departure time
	departureTime := time.Now()
	if c.QueryParam("departure_time") != "" {
		var err error
		departureTime, err = time.Parse(time.RFC3339, c.QueryParam("departure_time"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: fmt.Sprintf("invalid departure time: %s", err.Error())})
		}
	}

	// Parse seats
	seats := 1
	if c.QueryParam("seats") != "" {
		var err error
		seats, err = strconv.Atoi(c.QueryParam("seats"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError{Message: fmt.Sprintf("invalid seats: %s", err.Error())})
		}
	}

	// Parse page and perPage
	page, perPage := utils.ParsePage(c)

	ctx := c.Request().Context()

	flights, err := handler.Service.Fetch(ctx, departure, destination, departureTime, seats, page, perPage)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseError{Message: err.Error()})
	}

	flightResponses := make([]domain.FlightResponse, len(flights))
	for i, flight := range flights {
		flightResponses[i] = domain.FlightResponse{
			ID:            flight.ID,
			Departure:     flight.Departure,
			Destination:   flight.Destination,
			DepartureTime: flight.DepartureTime,
			Price:         flight.Price,
			Status:        flight.Status,
		}
	}

	return c.JSON(http.StatusOK, flightResponses)
}
