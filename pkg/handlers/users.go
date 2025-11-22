package handlers

type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type TeamNameQuery = string

type UserIdQuery = string

type GetUsersGetReviewParams struct {
	UserId UserIdQuery `form:"user_id" json:"user_id"`
}

type PostUsersSetIsActiveJSONBody struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
}
