package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
)

type ChapterRatingForm struct {
	Rating int8
}

func RateChapter(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	form := new(ChapterRatingForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	story, chapter, statusCode, err := validateRatingForm(ctx, form, userUuid)
	if err != nil {
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}
	userRating := collections.NewChapterUserRating()
	userRating.UserUuid = userUuid
	userRating.StoryId = story.Id
	userRating.ChapterId = chapter.Id
	userRating.Rating = form.Rating
	if err := userRating.UpsertDocument(); err != nil {
		return err
	}
	//todo:: set as a background job
	if err := chapter.UpdateRating(); err != nil {
		return err
	}
	if err := story.UpdateRating(); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "chapter rated successfully", ""))
}

func validateRatingForm(ctx echo.Context, form *ChapterRatingForm, userUuid string) (story collections.Story, chapter collections.Chapter, statusCode int, err error) {
	storyId := ctx.Param("id")
	if err := story.LoadById(storyId); err != nil {
		return story, chapter, http.StatusNotFound, errors.New("requested story not found")
	}

	chapterId := ctx.Param("chapterId")
	if err := chapter.LoadById(chapterId); err != nil {
		return story, chapter, http.StatusNotFound, errors.New("requested chapter not found")
	}
	if story.UserUuid == userUuid {
		return story, chapter, http.StatusForbidden, errors.New("user cannot rate his own story")
	}

	if chapter.UserUuid == userUuid {
		return story, chapter, http.StatusForbidden, errors.New("user cannot rate user own story")
	}

	if chapter.StoryId != story.Id {
		return story, chapter, http.StatusForbidden, errors.New("invalid story and chapter")
	}

	if form.Rating < 1 || form.Rating > 5 {
		return story, chapter, http.StatusUnprocessableEntity, errors.New("rating must be between 1 and 5")
	}

	return
}
