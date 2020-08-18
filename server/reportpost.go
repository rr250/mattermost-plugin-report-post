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
	ChannelName      string    `json:"channel_name"`
	ReportedUserName string    `json:"reported_username"`
	ReportedEmail    string    `json:"reported_email"`
	ReportedText     string    `json:"reported_text"`
	ReportedTextID   string    `json:"reported_text_id"`
	Reason           string
}

type PostDetails struct {
	PostID        string `json:"post_id"`
	CurrentUserID string `json:"current_user_id"`
}
