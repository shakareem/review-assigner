package api

import (
	"encoding/json"
	"net/http"

	"github.com/shakareem/review-assigner/internal/storage"
)

type PullRequestShort struct {
	AuthorId        string           `json:"author_id"`
	PullRequestId   string           `json:"pull_request_id"`
	PullRequestName string           `json:"pull_request_name"`
	Status          storage.PRStatus `json:"status"`
}

type PullRequestCreate struct {
	AuthorId string `json:"author_id"`
	PRID     string `json:"pull_request_id"`
	PRName   string `json:"pull_request_name"`
}

type PullRequestMerge struct {
	PRID string `json:"pull_request_id"`
}

type PullRequestReassign struct {
	OldUserId string `json:"old_user_id"`
	PRID      string `json:"pull_request_id"`
}

type CreatePRResponse struct {
	PR storage.PullRequest `json:"pr"`
}

type MergePRResponse struct {
	PR storage.PullRequest `json:"pr"`
}

type ReassignPRResponse struct {
	PR         storage.PullRequest `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

func (h *Handler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req PullRequestCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		invalidBody(w)
		return
	}

	pr, err := h.Storage.CreatePullRequest(req.PRID, req.PRName, req.AuthorId)
	if err != nil {

		switch err {
		case storage.ErrPRAlreadyExists:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: Error{Code: PREXISTS, Message: req.PRID + " already exists"},
			})
		case storage.ErrNotFound:
			notFound(w)
		default:
			internalError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatePRResponse{
		PR: storage.PullRequest{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorId:          pr.AuthorId,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
			CreatedAt:         pr.CreatedAt,
			MergedAt:          pr.MergedAt,
		},
	})
}

func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req PullRequestMerge
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		invalidBody(w)
		return
	}

	pr, err := h.Storage.MergePullRequest(req.PRID)
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
	json.NewEncoder(w).Encode(MergePRResponse{
		PR: storage.PullRequest{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorId:          pr.AuthorId,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
			CreatedAt:         pr.CreatedAt,
			MergedAt:          pr.MergedAt,
		},
	})
}

func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var req PullRequestReassign
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		invalidBody(w)
		return
	}

	pr, err := h.Storage.ReassignPullRequest(req.PRID, req.OldUserId)
	if err != nil {

		switch err {
		case storage.ErrPRMerged:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: Error{Code: PRMERGED, Message: "cannot reassign on merged PR"},
			})
		case storage.ErrNotAssigned:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: Error{Code: NOTASSIGNED, Message: "reviewer is not assigned to this PR"},
			})
		case storage.ErrNoCandidate:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: Error{Code: NOCANDIDATE, Message: "no active replacement candidate in team"},
			})
		case storage.ErrNotFound:
			notFound(w)
		default:
			internalError(w, err)
		}
		return
	}

	newReviewerID := getNewReviewerID(pr.AssignedReviewers, req.OldUserId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ReassignPRResponse{
		PR: storage.PullRequest{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorId:          pr.AuthorId,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
			CreatedAt:         pr.CreatedAt,
			MergedAt:          pr.MergedAt,
		},
		ReplacedBy: newReviewerID,
	})
}

func getNewReviewerID(reviewers []string, oldUserID string) string {
	for _, id := range reviewers {
		if id != oldUserID {
			return id
		}
	}
	return ""
}
