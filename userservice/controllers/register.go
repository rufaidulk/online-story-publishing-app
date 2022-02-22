package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"userservice/adapters"
	"userservice/helper"
	"userservice/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RegisterForm struct {
	Name     string
	Email    string
	Password string
}

func Registration(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	form := new(RegisterForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	if err := validateRegistrationForm(form, db); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	passwd, err := helper.HashPassword(form.Password)
	if err != nil {
		return err
	}

	user := models.NewUser()
	user.Name = form.Name
	user.Email = form.Email
	user.Password = passwd
	user.Status = 10
	user.SetUuid()
	user.SetUserActive()
	tx := db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		return err
	}
	userProfile := models.NewUserProfile()
	userProfile.UserId = user.Id
	if err := tx.Create(&userProfile).Error; err != nil {
		return err
	}

	tx.Commit()
	userDetails := helper.UserDetails{
		Uuid:  user.Uuid.UuidStr(),
		Name:  user.Name,
		Email: user.Email,
	}

	return ctx.JSON(http.StatusCreated,
		helper.NewSuccessResponse(http.StatusCreated, "SUCCESS", userDetails))
}

func validateRegistrationForm(form *RegisterForm, db *gorm.DB) error {
	if (RegisterForm{}) == *form {
		return errors.New("All fields required.")
	}

	if ok, _ := regexp.MatchString("^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$", form.Email); !ok {
		return errors.New("Invalid email address.")
	}

	var count int64
	if db.Model(models.NewUser()).Where("email = ?", form.Email).Count(&count); count != 0 {
		return errors.New("Email already taken.")
	}

	return nil
}
