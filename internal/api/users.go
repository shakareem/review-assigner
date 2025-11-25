package api

import (
	"encoding/json"
	"net/http"

	"github.com/shakareem/review-assigner/internal/storage"
)

type GetUsersGetReviewParams struct {
	UserId string `form:"user_id" json:"user_id"`
}

type UsersSetIsActive struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
}

type UserResponse struct {
	User storage.User `json:"user"`
}

type UserReviewPRsResponse struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: Error{Code: BADREQUEST, Message: "user_id parameter required"},
		})
		return
	}

	prs, err := h.Storage.GetUserReviewPRs(userID)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			notFound(w)
		default:
			internalError(w, err)
		}
		return
	}

	shortPRs := make([]PullRequestShort, len(prs))
	for i, pr := range prs {
		shortPRs[i] = PullRequestShort{
			PullRequestId:   pr.ID,
			PullRequestName: pr.Name,
			AuthorId:        pr.AuthorId,
			Status:          pr.Status,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserReviewPRsResponse{
		UserId:       userID,
		PullRequests: shortPRs,
	})
}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UsersSetIsActive
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		invalidBody(w)
		return
	}

	user, err := h.Storage.SetUserIsActive(req.UserId, req.IsActive)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			notFound(w)
		default:
			internalError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		User: storage.User{
			ID:       user.ID,
			Name:     user.Name,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	})
}
