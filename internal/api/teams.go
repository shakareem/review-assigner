package api

import (
	"encoding/json"
	"net/http"

	"github.com/shakareem/review-assigner/internal/storage"
)

type Team struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type AddTeamResponse struct {
	Team Team `json:"team"`
}

type GetTeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		invalidBody(w)
		return
	}

	users := make([]storage.User, len(req.Members))
	for i, member := range req.Members {
		users[i] = storage.User{
			ID:       member.UserId,
			Name:     member.Username,
			TeamName: req.TeamName,
			IsActive: member.IsActive,
		}
	}

	if err := h.Storage.AddTeam(req.TeamName, users); err != nil {
		switch err {
		case storage.ErrTeamAlreadyExists:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: Error{Code: TEAMEXISTS, Message: req.TeamName + " already exists"},
			})
		case storage.ErrUserAlreadyExists:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: Error{Code: USEREXISTS, Message: "user from team " + req.TeamName + " already exists"},
			})
		default:
			internalError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AddTeamResponse{Team: req})
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: Error{Code: NOTFOUND, Message: "team_name parameter required"},
		})
		return
	}

	users, err := h.Storage.GetTeam(teamName)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			notFound(w)
		default:
			internalError(w, err)
		}
		return
	}

	members := make([]TeamMember, len(users))
	for i, user := range users {
		members[i] = TeamMember{
			UserId:   user.ID,
			Username: user.Name,
			IsActive: user.IsActive,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Team{
		TeamName: teamName,
		Members:  members,
	})
}
