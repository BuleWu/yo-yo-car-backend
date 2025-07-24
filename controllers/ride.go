package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

func NewRideController(
	rideApplication *applications.Ride,
) *Ride {
	return &Ride{
		rideApplication: rideApplication,
	}
}

type Ride struct {
	Controller
	rideApplication *applications.Ride
}

func (c Ride) GetRides(ctx *gin.Context) {
	data, appErr := c.rideApplication.GetRides()

	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, data, http.StatusOK)
}

func (c Ride) GetRideById(ctx *gin.Context) {
	id := ctx.Param("id")

	ride, appErr := c.rideApplication.GetRideById(id)

	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, ride, http.StatusOK)
}

func (c Ride) CreateRide(ctx *gin.Context) {
	var request applications.CreateRideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusInternalServerError)
		return
	}

	ride, appErr := c.rideApplication.CreateRide(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, ride, http.StatusCreated)
}

func (c Ride) UpdateRide(ctx *gin.Context) {
	var request applications.UpdateRideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusInternalServerError)
		return
	}

	request.RideID = ctx.Param("id")

	ride, appErr := c.rideApplication.UpdateRide(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, ride, http.StatusOK)
}

func (c Ride) DeleteRide(ctx *gin.Context) {
	if appErr := c.rideApplication.DeleteRide(ctx.Param("id")); appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, nil, http.StatusNoContent)
}
