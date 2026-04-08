package models

import "time"

type TicketStatus string

const (
	StatusTodo       TicketStatus = "todo"
	StatusInProgress TicketStatus = "in-progress"
	StatusDone       TicketStatus = "done"
)

type Ticket struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Status      TicketStatus `json:"status"`
	Folder      string       `json:"folder"`
	Branch      string       `json:"branch"`
	ProjectRoot string       `json:"project_root"`
	CreatedAt   time.Time    `json:"created_at"`
}

type WorkSession struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"` // seconds
}
