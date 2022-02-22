package models

import (
	"time"
)

type UserFollower struct {
	Id         int64
	UserId     int64
	FollowerId int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUserFollower() *UserFollower {
	return &UserFollower{}
}
