package handler

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
