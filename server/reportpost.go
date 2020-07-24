package main

import (
	"time"
)

type ReportPost struct {
	ID               string
	ReportedBy       string    `json:"reported_by"`
	ReportedByID     string    `json:"reported_by_id"`
	CreatedAt        time.Time `json:"created_at"`
	ReportedName     string    `json:"reported_name"`
	ReportedID       string    `json:"reported_id"`
	ChannelID        string    `json:"channel_id"`
	ReportedUserName string    `json:"reported_username"`
	ReportedEmail    string    `json:"reported_email"`
	ReportedText     string    `json:"reported_text"`
	ReportedTextID   string    `json:"reported_text_id"`
}
