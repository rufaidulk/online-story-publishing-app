package models

import (
	"time"
	"userservice/helper"
)

const UserActive int = 10

type User struct {
	Id        int64
	Uuid      Uuid
	Name      string
	Email     string
	Password  string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserData struct {
	Id        int64
	UuidText  string
	Name      string
	Email     string
	Password  string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserData) TableName() string {
	return "users"
}

func NewUserData() *UserData {
	return &UserData{}
}

func NewUser() *User {
	return &User{}
}

func (u *User) SetUuid() {
	u.Uuid = Uuid{uuidStr: helper.GenerateUuid()}
}

func (u *User) SetUserActive() {
	u.Status = UserActive
}
