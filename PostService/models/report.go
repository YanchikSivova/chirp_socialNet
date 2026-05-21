package models

import (
	"github.com/google/uuid"
	"time"
)

type Report struct {
	ReportID  uuid.UUID `json:"report_id"`
	PostID    uuid.UUID `json:"post_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
