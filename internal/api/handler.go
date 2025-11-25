package api

import (
	"encoding/json"
	"net/http"

	"github.com/shakareem/review-assigner/internal/storage"
)

type Storage interface {
	AddTeam(teamName string, users []storage.User) error
	GetTeam(teamName string) ([]storage.User, error)
	SetUserIsActive(userID string, isActive bool) (*storage.User, error)
	GetUserReviewPRs(userID string) ([]storage.PullRequest, error)
	CreatePullRequest(prID, prName, authorID string) (*storage.PullRequest, error)
	MergePullRequest(prID string) (*storage.PullRequest, error)
	ReassignPullRequest(prID, oldUserID string) (*storage.PullRequest, error)
}

type Handler struct {
	Storage Storage
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type ErrorCode string

const (
	NOCANDIDATE ErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorCode = "NOT_FOUND"
	PREXISTS    ErrorCode = "PR_EXISTS"
	PRMERGED    ErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorCode = "TEAM_EXISTS"
	USEREXISTS  ErrorCode = "USER_EXISTS"
	BADREQUEST  ErrorCode = "BAD_REQUEST"
)

func NewHandler(s Storage) *Handler {
	return &Handler{
		Storage: s,
	}
}

func internalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: Error{Code: NOTFOUND, Message: err.Error()},
	})
}

func invalidBody(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: Error{Code: BADREQUEST, Message: "invalid request body"},
	})
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: Error{Code: NOTFOUND, Message: "resource not found"},
	})
}
