package main

import (
	"log"
	"userservice/controllers"
	"userservice/helper"
	"userservice/middlewares"

	"github.com/labstack/echo/v4"
)

func main() {
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
	userJwtAuth := middlewares.UserJwtMiddleware

	configRoutes(e, userJwtAuth)

	e.Logger.Fatal(e.Start(":" + helper.GetEnv("APP_PORT")))
}

func configRoutes(e *echo.Echo, jwtAuth echo.MiddlewareFunc) {
	e.POST("/register", controllers.Registration)
	e.POST("/authenticate", controllers.Login)
	e.POST("/authorize", controllers.AuthorizeUser, jwtAuth)
	e.GET("/user-profile", controllers.ViewUserProfile, jwtAuth)
	e.PUT("/user-profile", controllers.UpdateUserProfile, jwtAuth)
	e.POST("/user/:uuid/follow", controllers.CreateFollower, jwtAuth)
	e.DELETE("/user/:uuid/follow", controllers.DeleteFollower, jwtAuth)
}
