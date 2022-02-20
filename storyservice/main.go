package main

import (
	"flag"
	"log"
	"net/http"
	v1 "storyservice/controllers/v1"
	"storyservice/helper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	ModeServer = "server"
	ModeSetup  = "setup"
)

func main() {
	modeFlag := flag.String("mode", ModeServer, "mode of execution")
	flag.Parse()
	log.Println(*modeFlag)
	if *modeFlag == ModeSetup {
		InitSetup()
		return
	}

	e := echo.New()

	// Route level middleware
	userIdentity := getUserIdentityMiddleware()

	configRoutes(e, userIdentity)

	e.Logger.Fatal(e.Start(":" + helper.GetEnv("APP_PORT")))
}

func configRoutes(e *echo.Echo, userIdentity echo.MiddlewareFunc) {
	e.POST("/category", v1.CreateCategory)
	e.POST("/story", v1.CreateStory, userIdentity)
	e.PUT("/story/:id/promotional-info", v1.UpdateStoryPromotionalInfo, userIdentity)
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
