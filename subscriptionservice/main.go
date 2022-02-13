package main

import (
	"log"
	"net/http"
	"subscriptionservice/controllers"
	"subscriptionservice/helper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Route level middleware
	userIdentity := getUserIdentityMiddleware()

	configRoutes(e, userIdentity)

	e.Logger.Fatal(e.Start(":" + helper.GetEnv("APP_PORT")))
}

func configRoutes(e *echo.Echo, userIdentity echo.MiddlewareFunc) {
	e.GET("/subscriptions", controllers.CurrentSubscription, userIdentity)
	e.POST("/subscriptions", controllers.CreateBaseSubscription, userIdentity)
	e.POST("/subscriptions/upgrade", controllers.UpgradeSubscriptionToPremium, userIdentity)
}

func getUserIdentityMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			log.Println("Verifying UUID...")
			userUuid := ctx.Request().Header.Get("useruuid")
			if !IsValidUUID(userUuid) {
				return ctx.JSON(http.StatusUnauthorized,
					helper.NewErrorResponse(http.StatusUnauthorized, "UUID is not valid"))
			} else {
				ctx.Set("userUuid", userUuid)
			}
			log.Println("UUID verified")
			return next(ctx)
		}
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)

	return err == nil
}
