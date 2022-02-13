package models

import (
	"errors"
	"subscriptionservice/adapters"
	"time"

	"github.com/labstack/echo/v4"
)

//possible values for Type
const Monthly int8 = 10
const Yearly int8 = 11

//possible values for IsPremium
const PremiumPlan int8 = 10
const BasicPlan int8 = 11

const ActivePlan int8 = 10
const InactivePlan int8 = 11

type SubscriptionPlan struct {
	Id           int64
	PeriodType   int8
	Name         string
	PlanType     int8
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

	db.Where("period_type = ? AND plan_type = ? AND status = ?", Monthly, BasicPlan, ActivePlan).First(s)
}

func (s *SubscriptionPlan) GetPremiumSubscriptionPlan(periodType string, ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	switch periodType {
	case "monthly":
		db.Where("period_type = ? AND plan_type = ? AND status = ?", Monthly, PremiumPlan, ActivePlan).First(s)
	case "yearly":
		db.Where("period_type = ? AND plan_type = ? AND status = ?", Yearly, PremiumPlan, ActivePlan).First(s)
	default:
		return errors.New("invalid period type")
	}

	return nil
}
