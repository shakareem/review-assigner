package handler

type GetUsersGetReviewParams struct {
	UserId string `form:"user_id" json:"user_id"`
}

type UsersSetIsActive struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
}
