package middlewares

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"userservice/adapters"
	"userservice/helper"
	"userservice/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var UserJwtMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		db := adapters.GetDbHandle(ctx)
		log.Println("Verifying JWT...")
		authorization := ctx.Request().Header.Get("authorization")
		if userDetails, err := validateUserJwt(ctx, authorization, db); err != nil {
			return ctx.JSON(http.StatusUnauthorized,
				helper.NewErrorResponse(http.StatusUnauthorized, err.Error()))
		} else {
			ctx.Set("userDetails", userDetails)
		}
		log.Println("JWT verified")
		return next(ctx)
	}
}

func validateUserJwt(ctx echo.Context, authorization string, db *gorm.DB) (userDetails helper.UserDetails, err error) {
	if len(authorization) == 0 {
		return userDetails, errors.New("Unauthorized")
	}

	bearerToken := strings.Split(authorization, " ")
	if len(bearerToken) != 2 {
		return userDetails, errors.New("Supports only bearer token.")
	}

	userDetails, err = helper.DecodeJwt(bearerToken[1])
	if err != nil {
		return
	}

	user := models.NewUserData()
	uuid := models.NewUuid(userDetails.Uuid)
	dbErr := db.Where("email = ? AND uuid = ?", userDetails.Email, uuid).Take(&user).Error
	if dbErr != nil {
		return userDetails, errors.New("Corrupted JWT.")
	}

	ctx.Set("userData", *user)
	return
}
