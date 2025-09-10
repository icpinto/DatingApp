package models

import "time"

type FriendRequest struct {
	RequestId        int       `json:"id"`
	SenderID         int       `json:"sender_id"`
	SenderUsername   string    `json:"sender_username"`
	ReceiverID       int       `json:"receiver_id"`
	ReceiverUsername string    `json:"receiver_username"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AcceptRequest struct {
	RequestID int `json:"id"`
}

type RejectRequest struct {
	RequestID int `json:"id"`
}
