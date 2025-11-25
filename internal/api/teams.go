package api

import "net/http"

type Team struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type GetTeamGetParams struct {
	TeamName string `form:"team_name" json:"team_name"`
}

func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request) {
	// TODO
}
