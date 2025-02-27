package utils

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"github.com/tingyoulin/go-ticket-booking/domain"
)

const (
	defaultPage    = 1
	defaultPerPage = 10
)

func ParsePage(c echo.Context) (int, int) {
	page := defaultPage
	if value, err := strconv.Atoi(c.QueryParam("page")); value > 0 && err == nil {
		page = value
	}
	perPage := defaultPerPage
	if value, err := strconv.Atoi(c.QueryParam("per_page")); value > 0 && err == nil {
		perPage = value
	}
	return page, perPage
}

func IsRequestValid(model interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(model)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrFlightCanceled, domain.ErrFlightNoAvailableSeats:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}
