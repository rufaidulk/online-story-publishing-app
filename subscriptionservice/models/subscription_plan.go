package models

import (
	"subscriptionservice/adapters"
	"time"

	"github.com/labstack/echo/v4"
)

//possible values for Type
const Monthly int8 = 10
const Yearly int8 = 11

//possible values for IsPremium
const PremiumSubscription int8 = 10
const BasicSubscription int8 = 11

const ActiveSubscription int8 = 10
const InactiveSubscription int8 = 11

type SubscriptionPlan struct {
	Id           int64
	Type         int8
	Name         string
	IsPremium    int8
	Amount       float64
	PeriodInDays int8
	Status       int8
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewSubscriptionPlan() *SubscriptionPlan {
	return &SubscriptionPlan{}
}

func (s *SubscriptionPlan) GetBasicSubscriptionPlan(ctx echo.Context) {
	db := adapters.GetDbHandle(ctx)
	db.Where("type = ? AND is_premium = ? AND status = ?", Monthly, BasicSubscription, ActiveSubscription).First(s)

}
