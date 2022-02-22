package controllers

import (
	"net/http"
	"userservice/helper"

	"github.com/labstack/echo/v4"
)

func AuthorizeUser(ctx echo.Context) error {
	userDetails := ctx.Get("userDetails").(helper.UserDetails)

	return ctx.JSON(http.StatusOK,
		helper.NewSuccessResponse(http.StatusOK, "authorized", userDetails))
}
