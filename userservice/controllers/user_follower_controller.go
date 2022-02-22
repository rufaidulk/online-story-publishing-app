package controllers

import (
	"errors"
	"log"
	"net/http"
	"userservice/adapters"
	"userservice/helper"
	"userservice/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Person A follows Person B and C follows A. Then A has a follower C and B is followee of A.
func CreateFollower(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	user, ok := ctx.Get("userData").(models.UserData)
	if !ok {
		return errors.New("type assertion failed")
	}
	follower, statusCode, err := validateAddFollowerRequest(user, ctx.Param("uuid"), db)
	if err != nil {
		log.Println(err)
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}

	userFollower := models.NewUserFollower()
	userFollower.UserId = user.Id
	userFollower.FollowerId = follower.Id

	if err := db.Create(&userFollower).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, helper.NewSuccessResponse(http.StatusCreated, "added follower", ""))
}

func validateAddFollowerRequest(user models.UserData, followerUuid string, db *gorm.DB) (follower models.UserData, statusCode int, err error) {
	uuid := models.NewUuid(followerUuid)
	err = db.Where("uuid = ?", uuid).Take(&follower).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return follower, http.StatusUnprocessableEntity, errors.New("requested user not found.")
	} else if err != nil {
		return follower, http.StatusInternalServerError, err
	}

	userFollower := models.NewUserFollower()
	err = db.Where("user_id = ? AND follower_id = ?", user.Id, follower.Id).Take(&userFollower).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return follower, 0, nil
	} else if err != nil {
		return follower, http.StatusInternalServerError, err
	} else if userFollower.Id != 0 {
		return follower, http.StatusUnprocessableEntity, errors.New("already followed.")
	}

	return
}
