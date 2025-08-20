package model

type Trip struct {
	ID          int64  `json:"id"`
	LeaderID    int64  `json:"leader_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
