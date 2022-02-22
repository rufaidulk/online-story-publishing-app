package models

import (
	"time"
)

type UserProfile struct {
	Id              int64
	UserId          int64
	Description     string
	IsAuthor        bool
	IsPremiumAuthor bool
	FollowersCount  int64
	FolloweeCount   int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewUserProfile() *UserProfile {
	return &UserProfile{}
}
