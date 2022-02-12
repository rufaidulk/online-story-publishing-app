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

type LoginForm struct {
	Email    string
	Password string
}

func Login(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	form := new(LoginForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	user, err := validateLoginRequest(form, db)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	token, err := helper.CreateJwt(helper.NewUserDetails(user.UuidText, user.Name, user.Email))
	if err != nil {
		return err
	}

	res := make(map[string]string)
	res["token"] = token
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusOK, "SUCCESS", res))
}

func validateLoginRequest(form *LoginForm, db *gorm.DB) (user models.UserData, err error) {
	if (LoginForm{}) == *form {
		return user, errors.New("All fields required.")
	}

	if ok, _ := regexp.MatchString("^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$", form.Email); !ok {
		return user, errors.New("Invalid email address.")
	}

	if err := db.Where("email = ?", form.Email).Take(&user).Error; err != nil {
		return user, errors.New("User not found.")
	} else if !helper.ValidatePassword(form.Password, user.Password) {
		return user, errors.New("Wrong password.")
	}

	return user, nil
}
