package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChapterForm struct {
	Title            string
	Body             string
	IsStoryCompleted bool `json:"is_story_completed"`
}

type ChapterResponse struct {
	Id        primitive.ObjectID `json:"_id"`
	UserUuid  string             `json:"user_uuid"`
	StoryId   primitive.ObjectID `json:"story_id"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	Rating    int8               `json:"rating"`
	ReadCount int64              `json:"read_count"`
}

func CreateChapter(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	form := new(ChapterForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	storyId := ctx.Param("id")
	story := collections.NewStory()
	if err := story.LoadById(storyId); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "invalid story"))
	}
	if statusCode, err := validateChapterForm(story, form, userUuid); err != nil {
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}
	chapter := collections.NewChapter()
	chapter.UserUuid = userUuid
	chapter.Title = form.Title
	chapter.Body = form.Body
	chapter.StoryId = story.Id
	if err := chapter.CreateDocument(); err != nil {
		return err
	}
	if form.IsStoryCompleted {
		story.IsCompleted = true
	}
	story.AddChapter(chapter.Id, chapter.Title)
	if err := story.Update(); err != nil {
		return err
	}

	res := buildChapterResponse(chapter)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "chapter created", res))
}

func validateChapterForm(story *collections.Story, form *ChapterForm, userUuid string) (int, error) {
	if story.UserUuid != userUuid {
		return http.StatusForbidden, errors.New("forbidden")
	}

	if story.IsCompleted {
		return http.StatusUnprocessableEntity, errors.New("story already completed")
	}

	if form.Body == "" {
		return http.StatusUnprocessableEntity, errors.New("body is required")
	}

	return 0, nil
}

func buildChapterResponse(chapter *collections.Chapter) ChapterResponse {
	res := ChapterResponse{
		Id:        chapter.Id,
		UserUuid:  chapter.UserUuid,
		StoryId:   chapter.StoryId,
		Title:     chapter.Title,
		Body:      chapter.Body,
		Rating:    chapter.Rating,
		ReadCount: chapter.ReadCount,
	}

	return res
}
