package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/runtimebag"
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

func (c *Reservation) ConfirmReservation(ctx *gin.Context) {
	confirmationToken := ctx.Query("token")
	if confirmationToken == "" {
		frontendUrl := runtimebag.GetEnvString("FRONTEND_URL", "")
		failedRedirect := fmt.Sprintf("%s/reservation-confirm?status=invalid", frontendUrl)
		ctx.Redirect(http.StatusTemporaryRedirect, failedRedirect)
		return
	}

	var request applications.UpdateReservationRequest
	request.ReservationID = ctx.Param("id")
	request.Status = models.ReservationConfirmed
	request.ConfirmationToken = confirmationToken

	frontendUrl := runtimebag.GetEnvString("FRONTEND_URL", "")

	if _, appErr := c.reservationApplication.UpdateReservation(&request); appErr != nil {
		failedRedirect := fmt.Sprintf("%s/reservation-confirm?status=failed", frontendUrl)
		ctx.Redirect(http.StatusTemporaryRedirect, failedRedirect)
		return
	}

	successRedirect := fmt.Sprintf("%s/reservation-confirm?status=success", frontendUrl)
	ctx.Redirect(http.StatusTemporaryRedirect, successRedirect)
}

func (c *Reservation) DeclineReservation(ctx *gin.Context) {
	confirmationToken := ctx.Query("token")
	if confirmationToken == "" {
		frontendUrl := runtimebag.GetEnvString("FRONTEND_URL", "")
		failedRedirect := fmt.Sprintf("%s/reservation-decline?status=invalid", frontendUrl)
		ctx.Redirect(http.StatusTemporaryRedirect, failedRedirect)
		return
	}

	var request applications.UpdateReservationRequest
	request.ReservationID = ctx.Param("id")
	request.Status = models.ReservationCancelled
	request.ConfirmationToken = confirmationToken

	frontendUrl := runtimebag.GetEnvString("FRONTEND_URL", "")

	if _, appErr := c.reservationApplication.UpdateReservation(&request); appErr != nil {
		failedRedirect := fmt.Sprintf("%s/reservation-decline?status=failed", frontendUrl)
		ctx.Redirect(http.StatusTemporaryRedirect, failedRedirect)
		return
	}

	successRedirect := fmt.Sprintf("%s/reservation-decline?status=success", frontendUrl)
	ctx.Redirect(http.StatusTemporaryRedirect, successRedirect)
}
