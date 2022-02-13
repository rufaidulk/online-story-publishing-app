package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"subscriptionservice/adapters"
	"subscriptionservice/helper"
	"subscriptionservice/models"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UpgradeForm struct {
	PeriodType string
}

type SubscriptionResponse struct {
	Id         int64      `json:"id"`
	UserUuid   string     `json:"user_uuid"`
	IsPremium  bool       `json:"is_premium"`
	ExpiryDate *time.Time `json:"expiry_date"`
}

func CurrentSubscription(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	userUuid := ctx.Get("userUuid").(string)
	uuidData := models.NewUuidData(userUuid)
	subscription := models.NewSubscriptionData()
	err := db.Where("status = ? AND user_uuid = ?", models.SubscriptionActive, uuidData).Take(subscription).Error
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "No subscription"))
	}
	isPremium := true
	if subscription.IsPremium == models.FreeSubscription {
		isPremium = false
	}

	res := SubscriptionResponse{
		Id:         subscription.Id,
		UserUuid:   subscription.UserUuidText,
		IsPremium:  isPremium,
		ExpiryDate: subscription.ExpiryDate,
	}

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusOK, "SUCCESS", res))
}

func CreateBaseSubscription(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	userUuid := ctx.Get("userUuid").(string)
	_, err := validateBaseSubscriptionRequest(userUuid, db)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	subscriptionPlan := models.NewSubscriptionPlan()
	subscriptionPlan.GetBasicSubscriptionPlan(ctx)

	tx := db.Begin()
	paymentTxnId, err := handlePayment(subscriptionPlan.Amount, tx)
	if err != nil {
		return err
	}
	subscription := models.NewSubscription()
	subscription.SetUuid(userUuid)
	subscription.SubscriptionPlanId = subscriptionPlan.Id
	subscription.PaymentTransactionId = paymentTxnId
	subscription.IsPremium = models.FreeSubscription
	subscription.SetSubscriptionActive()

	if err := tx.Create(&subscription).Error; err != nil {
		return err
	}

	tx.Commit()

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "SUCCESS", ""))
}

func UpgradeSubscriptionToPremium(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	form := new(UpgradeForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	userUuid := ctx.Get("userUuid").(string)
	baseSubscription, err := validateBaseSubscriptionRequest(userUuid, db)
	if err == nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "user doesn't have a base subscription"))
	}

	subscriptionPlan := models.NewSubscriptionPlan()
	if err := subscriptionPlan.GetPremiumSubscriptionPlan(form.PeriodType, ctx); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	tx := db.Begin()
	if err := db.Model(&baseSubscription).Update("status", models.SubscriptionCompleted).Error; err != nil {
		return err
	}

	paymentTxnId, err := handlePayment(subscriptionPlan.Amount, tx)
	if err != nil {
		return err
	}

	subscription := models.NewSubscription()
	subscription.SetUuid(userUuid)
	subscription.SubscriptionPlanId = subscriptionPlan.Id
	subscription.PaymentTransactionId = paymentTxnId
	subscription.IsPremium = models.PremiumSubscription
	expiryDate := time.Now().AddDate(0, 0, int(subscriptionPlan.PeriodInDays))
	subscription.ExpiryDate = &expiryDate
	subscription.SetSubscriptionActive()
	if err := tx.Create(&subscription).Error; err != nil {
		return err
	}

	tx.Commit()

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "SUCCESS", ""))
}

func validateBaseSubscriptionRequest(userUuid string, db *gorm.DB) (models.Subscription, error) {
	uuidData := models.NewUuidData(userUuid)
	subscription := models.NewSubscription()
	err := db.Where("status = ? AND user_uuid = ? and is_premium = ?", models.SubscriptionActive, uuidData, models.FreeSubscription).Take(subscription).Error
	if err == nil {
		return *subscription, errors.New("requested user has an existing subscription")
	}

	return *subscription, nil
}

func handlePayment(amt float64, tx *gorm.DB) (int64, error) {
	sequence := models.NewSequence()
	if err := tx.Where("type = ?", models.SubscriptionTxn).Take(&sequence).Error; err != nil {
		return 0, err
	}
	tx.Model(&sequence).Update("seq_no", sequence.SeqNo+1)

	refNo := fmt.Sprintf("TX%06d", sequence.SeqNo)
	if amt > 0 {
		fmt.Println("Calling Payment Gateway...")
	}

	paymentTxn := models.PaymentTransaction{
		RefNo:  refNo,
		Amount: amt,
		Type:   "CREDIT",
	}
	paymentTxn.SetPaymentCompleted()
	if err := tx.Create(&paymentTxn).Error; err != nil {
		return 0, err
	}

	return paymentTxn.Id, nil
}
