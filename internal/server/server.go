package server

import (
	"net/http"

	"github.com/shakareem/review-assigner/internal/handler"
	"github.com/shakareem/review-assigner/internal/storage"
)

type Handler interface {
	PostPullRequestCreate(w http.ResponseWriter, r *http.Request)
	PostPullRequestMerge(w http.ResponseWriter, r *http.Request)
	PostPullRequestReassign(w http.ResponseWriter, r *http.Request)
	PostTeamAdd(w http.ResponseWriter, r *http.Request)
	GetTeamGet(w http.ResponseWriter, r *http.Request)
	GetUsersGetReview(w http.ResponseWriter, r *http.Request)
	PostUsersSetIsActive(w http.ResponseWriter, r *http.Request)
}

type Storage interface {
	AddTeam(handler.Team) (handler.Team, error)
	GetTeam(teamName string) (handler.Team, error)
	SetUserIsActive(userID string, isActive bool) (storage.User, error)
	CreatePullRequest(prID, prName, authorID string) (handler.PullRequest, error)
	MergePullRequest(prID string) (handler.PullRequest, error)
	ReassignPullRequest(prID, oldUserID string) (handler.PullRequest, error)
}

type Server struct {
	Server *http.Server
}

const PORT = ":8080"

func NewServer(h Handler) *Server {
	m := http.NewServeMux()

	m.HandleFunc("POST /pullRequest/create", h.PostPullRequestCreate)
	m.HandleFunc("POST /pullRequest/merge", h.PostPullRequestMerge)
	m.HandleFunc("POST /pullRequest/reassign", h.PostPullRequestReassign)
	m.HandleFunc("POST /team/add", h.PostTeamAdd)
	m.HandleFunc("GET /team/get", h.GetTeamGet)
	m.HandleFunc("GET /users/getReview", h.GetUsersGetReview)
	m.HandleFunc("POST /users/setIsActive", h.PostUsersSetIsActive)

	return &Server{
		&http.Server{
			Addr:    PORT,
			Handler: m,
		},
	}
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
