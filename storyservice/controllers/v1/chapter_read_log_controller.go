package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
)

func CreateChapterReadLog(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	form := new(ChapterRatingForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	story, chapter, statusCode, err := validateChapterReadLogRequest(ctx, userUuid)
	if err != nil {
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}
	chapterReadLog := collections.NewChapterReadLog()
	chapterReadLog.UserUuid = userUuid
	chapterReadLog.StoryId = story.Id
	chapterReadLog.ChapterId = chapter.Id
	if isFirstTimeRead, err := chapterReadLog.UpsertDocument(); err != nil {
		return err
	} else {
		if isFirstTimeRead {
			if err := chapter.IncrementReadCount(); err != nil {
				return err
			}
			if err := story.UpdateAvgReadCount(); err != nil {
				return err
			}
		}
	}

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "chapter read log created", ""))
}

func validateChapterReadLogRequest(ctx echo.Context, userUuid string) (story collections.Story, chapter collections.Chapter, statusCode int, err error) {
	storyId := ctx.Param("id")
	if err := story.LoadById(storyId); err != nil {
		return story, chapter, http.StatusNotFound, errors.New("requested story not found")
	}

	chapterId := ctx.Param("chapterId")
	if err := chapter.LoadById(chapterId); err != nil {
		return story, chapter, http.StatusNotFound, errors.New("requested chapter not found")
	}
	if story.UserUuid == userUuid {
		return story, chapter, http.StatusUnprocessableEntity, errors.New("this read cannot be logged")
	}

	if chapter.StoryId != story.Id {
		return story, chapter, http.StatusForbidden, errors.New("invalid story and chapter")
	}

	return
}
