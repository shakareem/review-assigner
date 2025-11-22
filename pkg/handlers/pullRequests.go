package handlers

import "time"

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequestStatus string

type PullRequest struct {
	AssignedReviewers []string          `json:"assigned_reviewers"`
	AuthorId          string            `json:"author_id"`
	CreatedAt         *time.Time        `json:"createdAt"`
	MergedAt          *time.Time        `json:"mergedAt"`
	PullRequestId     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	Status            PullRequestStatus `json:"status"`
}

type PullRequestShort struct {
	AuthorId        string            `json:"author_id"`
	PullRequestId   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	Status          PullRequestStatus `json:"status"`
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
