package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListPremiumStories(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	var limit int64 = 15
	data, _ := collections.ListPremiumStories(userUuid, limit)

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "premium stories list", data))
}

func ListTrendingStories(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	var limit int64 = 15
	data, _ := collections.ListTrendingStories(userUuid, limit)

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "trending stories list", data))
}

func ListMostRatedStories(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	var limit int64 = 15
	data, _ := collections.ListTrendingStories(userUuid, limit)

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "top stories list", data))
}

func ListStoriesByCategory(ctx echo.Context) error {
	_, category, statusCode, err := validateListStoriesByCategoryRequest(ctx)
	if err != nil {
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}
	var limit int64 = 2
	var skip int64
	page := ctx.Param("page")
	if page != "" {
		p, _ := strconv.ParseInt(page, 10, 64)
		skip = limit*p - 1
	}

	data, _ := collections.ListStoriesByCategory(category.Id, skip, limit)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "story list by category", data))
}

func validateListStoriesByCategoryRequest(ctx echo.Context) (story collections.Story, category collections.Category, statusCode int, err error) {
	storyId := ctx.Param("id")
	if err := story.LoadById(storyId); err != nil {
		return story, category, http.StatusNotFound, errors.New("requested story not found")
	}

	categoryId, _ := primitive.ObjectIDFromHex(ctx.Param("categoryId"))
	if err := category.LoadById(categoryId); err != nil {
		return story, category, http.StatusNotFound, errors.New("requested category not found")
	}
	categoryNotFound := true
	for _, v := range story.Categories {
		if v == categoryId {
			categoryNotFound = false
		}
	}
	if categoryNotFound {
		return story, category, http.StatusUnprocessableEntity, errors.New("invalid story and category")
	}

	return
}
