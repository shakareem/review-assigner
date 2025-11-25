package api

import (
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
	AuthorId        string `json:"author_id"`
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
}

type PullRequestMerge struct {
	PullRequestId string `json:"pull_request_id"`
}

type PullRequestReassign struct {
	OldUserId     string `json:"old_user_id"`
	PullRequestId string `json:"pull_request_id"`
}

func (h *Handler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	// TODO
}
