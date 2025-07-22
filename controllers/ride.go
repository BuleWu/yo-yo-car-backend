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
