package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/core/utils"
)

func NewRatingController(
	ratingApplication *applications.Rating,
) *Rating {
	return &Rating{
		ratingApplication: ratingApplication,
	}
}

type Rating struct {
	Controller
	ratingApplication *applications.Rating
}

func (c *Rating) GetRatings(ctx *gin.Context) {
	ratings, appErr := c.ratingApplication.GetRatings()
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, ratings, http.StatusOK)
}

func (c *Rating) GetRatingById(ctx *gin.Context) {
	rating, appErr := c.ratingApplication.GetRatingById(ctx.Param("id"))
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, rating, http.StatusOK)
}

func (c *Rating) GetRatingsByUserId(ctx *gin.Context) {
	ratings, appErr := c.ratingApplication.GetRatingsByUserId(ctx.Param("userId"))
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, ratings, http.StatusOK)
}

func (c *Rating) CreateRating(ctx *gin.Context) {
	var request applications.CreateRatingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}

	rating, appErr := c.ratingApplication.CreateRating(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, rating, http.StatusCreated)
}

func (c *Rating) UpdateRating(ctx *gin.Context) {
	var request applications.UpdateRatingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.returnJSON(ctx, utils.NewHttpError(err.Error()), http.StatusBadRequest)
		return
	}
	request.RatingID = ctx.Param("id")

	rating, appErr := c.ratingApplication.UpdateRating(&request)
	if appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}

	c.returnJSON(ctx, rating, http.StatusOK)
}

func (c *Rating) DeleteRating(ctx *gin.Context) {
	if appErr := c.ratingApplication.DeleteRating(ctx.Param("id")); appErr != nil {
		c.returnJSON(ctx, utils.NewHttpError(appErr.GetMessage()), appErr.GetCode())
		return
	}
	c.returnJSON(ctx, nil, http.StatusNoContent)
}
