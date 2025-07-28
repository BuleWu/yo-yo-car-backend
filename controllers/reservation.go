package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

func NewReservationController(reservationApplication *applications.Reservation) *Reservation {
	return &Reservation{
		reservationApplication: reservationApplication,
	}
}

type Reservation struct {
	Controller
	reservationApplication *applications.Reservation
}

// CreateReservation POST /reservations
func (c *Reservation) CreateReservation(ctx *gin.Context) {
	var request applications.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	reservation, appErr := c.reservationApplication.CreateReservation(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, reservation, http.StatusCreated)
}

// GetReservationById GET /reservations/:id
func (c *Reservation) GetReservationById(ctx *gin.Context) {
	id := ctx.Param("id")

	reservation, appErr := c.reservationApplication.GetReservationById(id)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, reservation, http.StatusOK)
}

// GetAllReservations GET /reservations
func (c *Reservation) GetAllReservations(ctx *gin.Context) {
	reservations, appErr := c.reservationApplication.GetAllReservations()
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, reservations, http.StatusOK)
}

// UpdateReservation PUT /reservations/:id
func (c *Reservation) UpdateReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	var request applications.UpdateReservationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	request.ReservationID = id

	reservation, appErr := c.reservationApplication.UpdateReservation(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, reservation, http.StatusOK)
}

// DeleteReservation DELETE /reservations/:id
func (c *Reservation) DeleteReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	appErr := c.reservationApplication.DeleteReservation(id)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, nil, http.StatusNoContent)
}
