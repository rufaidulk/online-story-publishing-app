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

func CreateBaseSubscription(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	userUuid := ctx.Get("userUuid").(string)
	err := validateBaseSubscriptionRequest(userUuid, db)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	subscriptionPlan := models.NewSubscriptionPlan()
	subscriptionPlan.GetBasicSubscriptionPlan(ctx)
	fmt.Println(subscriptionPlan)

	tx := db.Begin()
	paymentTxnId, err := handlePayment(subscriptionPlan.Amount, tx)
	if err != nil {
		return err
	}
	subscription := models.NewSubscription()
	subscription.SetUuid(userUuid)
	subscription.SubscriptionPlanId = subscription.Id
	subscription.PaymentTransactionId = paymentTxnId
	days := time.Duration(subscriptionPlan.PeriodInDays)
	subscription.ExpiryDate = time.Now().Add(time.Hour * days)
	subscription.SetSubscriptionActive()

	if err := tx.Create(&subscription).Error; err != nil {
		return err
	}

	tx.Commit()

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "SUCCESS", ""))
}

func validateBaseSubscriptionRequest(userUuid string, db *gorm.DB) error {
	uuidData := models.NewUuidData(userUuid)
	subscription := models.NewSubscription()
	err := db.Where("status = ? AND user_uuid = ?", models.ActiveSubscription, uuidData).Take(&subscription).Error
	if err == nil {
		return errors.New("Requested user has an existing subscription.")
	}

	return nil
}

func handlePayment(amt float64, tx *gorm.DB) (int64, error) {
	sequence := models.NewSequence()
	if err := tx.Where("type = ?", models.SubscriptionTxn).Take(&sequence).Error; err != nil {
		return 0, err
	}
	tx.Model(&sequence).Update("seq_no", sequence.SeqNo+1)

	refNo := fmt.Sprintf("TX%06d", sequence.SeqNo)
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
