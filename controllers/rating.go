package controllers

import "zavrsni/yo-yo-car/applications"

func NewRatingController(
	ratingApplication *applications.Rating
) *Rating {
	return &Rating{
		ratingApplication: ratingApplication,
	}
}

type Rating struct {
	Controller
	applications.Rating
}
