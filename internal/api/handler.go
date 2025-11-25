package api

import "github.com/shakareem/review-assigner/internal/storage"

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

type ErrorResponse struct {
	Error struct {
		Code    ErrorCode `json:"code"`
		Message string    `json:"message"`
	} `json:"error"`
}

type ErrorCode string

const (
	NOCANDIDATE ErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorCode = "NOT_FOUND"
	PREXISTS    ErrorCode = "PR_EXISTS"
	PRMERGED    ErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorCode = "TEAM_EXISTS"
)

func NewHandler(s Storage) *Handler {
	return &Handler{
		Storage: s,
	}
}
