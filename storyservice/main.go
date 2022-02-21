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
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		log.Println("Path: ", ctx.Path())
		log.Println("Query Params: ", ctx.QueryParams())
		log.Println("Path Param Names: ", ctx.ParamNames())
		log.Println("Path Param Values: ", ctx.ParamValues())
		log.Println(err)

		// Call the default handler to return the HTTP response
		e.DefaultHTTPErrorHandler(err, ctx)
	}
	// Route level middleware
	userIdentity := getUserIdentityMiddleware()

	configRoutes(e, userIdentity)

	e.Logger.Fatal(e.Start(":" + helper.GetEnv("APP_PORT")))
}

func configRoutes(e *echo.Echo, userIdentity echo.MiddlewareFunc) {
	e.POST("/category", v1.CreateCategory)
	e.POST("/story", v1.CreateStory, userIdentity)
	e.PUT("/story/:id/promotional-info", v1.UpdateStoryPromotionalInfo, userIdentity)
	e.GET("/story/:id", v1.ViewStory, userIdentity)
	e.PUT("/story/:id", v1.UpdateStory, userIdentity)
	e.POST("/story/:id/chapter", v1.CreateChapter, userIdentity)
	e.PUT("/story/:id/chapter/:chapterId", v1.UpdateChapter, userIdentity)
	e.GET("/story/:id/chapter/:chapterId", v1.ViewChapter, userIdentity)
	e.DELETE("/story/:id/chapter/:chapterId", v1.DeleteChapter, userIdentity)
	e.POST("/story/:id/chapter/:chapterId/read-log", v1.CreateChapterReadLog, userIdentity)
	e.POST("/story/:id/chapter/:chapterId/rating", v1.RateChapter, userIdentity)
	e.GET("/story/:id/category/:categoryId/:page", v1.ListStoriesByCategory, userIdentity)
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
