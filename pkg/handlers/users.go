package handlers

type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type GetUsersGetReviewParams struct {
	UserId string `form:"user_id" json:"user_id"`
}

type UsersSetIsActive struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
}
