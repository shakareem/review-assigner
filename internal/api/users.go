package api

import "net/http"

type GetUsersGetReviewParams struct {
	UserId string `form:"user_id" json:"user_id"`
}

type UsersSetIsActive struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	// TODO
}
