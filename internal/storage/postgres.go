package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

const (
	PRStatusMERGED PRStatus = "MERGED"
	PRStatusOPEN   PRStatus = "OPEN"
)

type PRStatus string

type PullRequest struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            PRStatus `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt"`
	MergedAt          *string  `json:"mergedAt"`
}

type Storage interface {
	AddTeam(teamName string, users []User) error
	GetTeam(teamName string) ([]User, error)
	SetUserIsActive(userID string, isActive bool) (*User, error)
	CreatePullRequest(prID, prName, authorID string) (*PullRequest, error)
	MergePullRequest(prID string) (*PullRequest, error)
	ReassignPullRequest(prID, oldUserID string) (*PullRequest, error)
}

var (
	ErrAlreadyExists = errors.New("object already exists")
	ErrDoesNotExist  = errors.New("object does not exist")
	ErrPRMerged      = errors.New("PR already merged")
)

func NewPostgresStorage() (*PostgresStorage, error) {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port= %d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) AddTeam(teamName string, users []User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var exists bool
	if err = tx.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1);`,
		teamName).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	if len(users) > 0 {
		ids := make([]string, 0, len(users))
		for _, u := range users {
			ids = append(ids, u.ID)
		}

		if err = tx.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ANY($1));`,
			pq.Array(ids)).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return ErrAlreadyExists
		}
	}

	if _, err = tx.Exec(
		`INSERT INTO teams (team_name) VALUES ($1);`,
		teamName); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrAlreadyExists
		}
		return err
	}

	if len(users) > 0 {
		args := make([]any, 0, len(users)*4)
		placeholders := ""
		for i, u := range users {
			if i > 0 {
				placeholders += ","
			}
			base := i * 4
			placeholders += fmt.Sprintf("($%d,$%d,$%d,$%d)", base+1, base+2, base+3, base+4)
			args = append(args, u.ID, u.Name, teamName, u.IsActive)
		}

		query := `INSERT INTO users (user_id, user_name, team_name, is_active) VALUES ` + placeholders
		if _, err = tx.Exec(query, args...); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				return ErrAlreadyExists
			}
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) GetTeam(teamName string) ([]User, error) {
	users := []User{}

	var exists bool
	if err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1);`,
		teamName).Scan(&exists); err != nil {
		return users, err
	}
	if !exists {
		return users, ErrDoesNotExist
	}

	rows, err := s.db.Query(
		`SELECT user_id, user_name, is_active FROM users WHERE team_name = $1;`,
		teamName)
	if err != nil {
		return users, err
	}

	var id, name string
	var isActive bool
	for rows.Next() {
		err := rows.Scan(&id, &name, &isActive)
		if err != nil {
			return users, err
		}

		users = append(users, User{
			ID:       id,
			Name:     name,
			TeamName: teamName,
			IsActive: isActive,
		})
	}

	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (s *PostgresStorage) SetUserIsActive(userID string, isActive bool) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(
		`UPDATE users SET is_active = $1 WHERE user_id = $2 
		RETURNING user_id, user_name, team_name, is_active;`,
		isActive, userID).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDoesNotExist
		}
		return nil, err
	}
	return user, nil
}

func (s *PostgresStorage) CreatePullRequest(prID, prName, authorID string) (*PullRequest, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var authorTeam string
	if err := tx.QueryRow(
		`SELECT team_name FROM users WHERE user_id = $1;`,
		authorID).Scan(&authorTeam); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDoesNotExist
		}
		return nil, err
	}

	createdAt := time.Now()

	if _, err = tx.Exec(
		`INSERT INTO pull_requests (pr_id, pr_name, author_id, pr_status, created_at) 
		VALUES ($1, $2, $3, $4, $5);`,
		prID, prName, authorID, PRStatusOPEN, createdAt); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}

	reviewersIDs, err := s.assignReviewers(prID, authorTeam, authorID, tx)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &PullRequest{
		ID:                prID,
		Name:              prName,
		Status:            PRStatusOPEN,
		CreatedAt:         createdAt.Format(time.RFC3339),
		AssignedReviewers: reviewersIDs,
	}, nil
}

func (s *PostgresStorage) assignReviewers(prID, teamName, prAuthor string, tx *sql.Tx) ([]string, error) {
	rows, err := tx.Query(
		`SELECT user_id FROM users 
		WHERE team_name = $1 AND is_active = TRUE AND user_id != $2
		LIMIT 2;`,
		teamName, prAuthor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewersIDs := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewersIDs = append(reviewersIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(reviewersIDs) > 0 {
		args := make([]any, 0, len(reviewersIDs)*2)
		placeholders := ""
		for i, id := range reviewersIDs {
			if i > 0 {
				placeholders += ","
			}
			base := i * 2
			placeholders += fmt.Sprintf("($%d,$%d)", base+1, base+2)
			args = append(args, prID, id)
		}

		query := `INSERT INTO pr_reviewers (pr_id, user_id) VALUES ` + placeholders
		if _, err := tx.Exec(query, args...); err != nil {
			return nil, err
		}
	}

	return reviewersIDs, nil
}

func (s *PostgresStorage) MergePullRequest(prID string) (*PullRequest, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	mergedAt := time.Now()
	pr := &PullRequest{}
	var createdAt time.Time

	if err := tx.QueryRow(
		`UPDATE pull_requests SET pr_status = $1, merged_at = $2 WHERE pr_id = $3
		 RETURNING pr_id, pr_name, author_id, pr_status, created_at;`,
		PRStatusMERGED, mergedAt, prID,
	).Scan(&pr.ID, &pr.Name, &pr.AuthorId, &pr.Status, &createdAt); err != nil {
		return nil, err
	}

	rows, err := tx.Query(
		`SELECT user_id FROM pr_reviewers WHERE pr_id = $1;`,
		prID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewersIDs := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewersIDs = append(reviewersIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewersIDs
	mergedAtString := mergedAt.Format(time.RFC3339)
	pr.MergedAt = &mergedAtString
	pr.CreatedAt = createdAt.Format(time.RFC3339)

	return pr, nil
}

func (s *PostgresStorage) ReassignPullRequest(prID, oldUserID string) (*PullRequest, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	pr := &PullRequest{}

	var createdAt, mergedAt *time.Time
	if err := tx.QueryRow(
		`SELECT pr_id, pr_name, author_id, pr_status, created_at, merged_at
		 FROM pull_requests WHERE pr_id = $1;`,
		prID,
	).Scan(&pr.ID, &pr.Name, &pr.AuthorId, &pr.Status, &createdAt, &mergedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDoesNotExist
		}
		return nil, err
	}

	if pr.Status == PRStatusMERGED {
		return nil, ErrPRMerged
	}

	var exists bool
	if err = tx.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2);`,
		prID, oldUserID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDoesNotExist
	}

	var authorTeam string
	if err := tx.QueryRow(
		`SELECT team_name FROM users WHERE user_id = $1;`,
		pr.AuthorId).Scan(&authorTeam); err != nil {
		return nil, err
	}

	var newReviewerID string
	if err := tx.QueryRow(
		`SELECT user_id FROM users 
		 WHERE team_name = $1 AND is_active = TRUE 
		 AND user_id != $2 AND user_id != $3
		 LIMIT 1;`,
		authorTeam, pr.AuthorId, oldUserID).Scan(&newReviewerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDoesNotExist // No available reviewer
		}
		return nil, err
	}

	if _, err := tx.Exec(
		`UPDATE pr_reviewers SET user_id = $1 WHERE pr_id = $2 AND user_id = $3;`,
		newReviewerID, prID, oldUserID); err != nil {
		return nil, err
	}

	rows, err := tx.Query(
		`SELECT user_id FROM pr_reviewers WHERE pr_id = $1;`,
		prID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewersIDs := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewersIDs = append(reviewersIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewersIDs
	if createdAt != nil {
		pr.CreatedAt = createdAt.Format(time.RFC3339)
	}
	if mergedAt != nil {
		mergedAtStr := mergedAt.Format(time.RFC3339)
		pr.MergedAt = &mergedAtStr
	}

	return pr, nil
}
