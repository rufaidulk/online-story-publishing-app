package models

import (
	"time"
)

const SubscriptionActive int = 10
const SubscriptionCompleted int = 11

const FreeSubscription int8 = 10
const PremiumSubscription int8 = 11

type Subscription struct {
	Id                   int64
	UserUuid             UuidData
	SubscriptionPlanId   int64
	PaymentTransactionId int64
	IsPremium            int8
	ExpiryDate           *time.Time
	Status               int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type SubscriptionData struct {
	Id                   int64
	UserUuidText         string
	SubscriptionPlanId   int64
	PaymentTransactionId int64
	IsPremium            int8
	ExpiryDate           *time.Time
	Status               int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (u *SubscriptionData) TableName() string {
	return "subscriptions"
}

func NewSubscriptionData() *SubscriptionData {
	return &SubscriptionData{}
}

func NewSubscription() *Subscription {
	return &Subscription{}
}

func (s *Subscription) SetUuid(uuid string) {
	s.UserUuid = UuidData{uuidStr: uuid}
}

func (s *Subscription) SetSubscriptionActive() {
	s.Status = SubscriptionActive
}

func (s *Subscription) SetSubscriptionCompleted() {
	s.Status = SubscriptionCompleted
}
