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

type ChapterUpdateForm struct {
	Title string
	Body  string
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
	if err := story.UpdateDocument(); err != nil {
		return err
	}

	res := buildChapterResponse(chapter)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "chapter created", res))
}

func UpdateChapter(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	form := new(ChapterUpdateForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	storyId := ctx.Param("id")
	story := collections.NewStory()
	if err := story.LoadById(storyId); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "invalid story"))
	}
	chapterId := ctx.Param("chapterId")
	chapter := collections.NewChapter()
	if err := chapter.LoadById(chapterId); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "invalid chapter"))
	}
	if statusCode, err := validateChapterUpdateForm(story, chapter, form, userUuid); err != nil {
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}
	chapter.Title = form.Title
	chapter.Body = form.Body
	if err := chapter.Update(); err != nil {
		return err
	}
	story.EditChapter(chapter.Id, chapter.Title)
	if err := story.UpdateDocument(); err != nil {
		return err
	}

	res := buildChapterResponse(chapter)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusOK, "chapter updated", res))
}

func ViewChapter(ctx echo.Context) error {
	storyId := ctx.Param("id")
	story := collections.NewStory()
	if err := story.LoadById(storyId); err != nil {
		return ctx.JSON(http.StatusNotFound,
			helper.NewErrorResponse(http.StatusNotFound, "requested story not found"))
	}
	chapterId := ctx.Param("chapterId")
	chapter := collections.NewChapter()
	if err := chapter.LoadById(chapterId); err != nil {
		return ctx.JSON(http.StatusNotFound,
			helper.NewErrorResponse(http.StatusNotFound, "requested chapter not found"))
	}
	if statusCode, err := validateChapterViewRequest(story, chapter); err != nil {
		return ctx.JSON(statusCode, helper.NewErrorResponse(statusCode, err.Error()))
	}

	res := buildChapterResponse(chapter)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusOK, "chapter details", res))
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

func validateChapterUpdateForm(story *collections.Story, chapter *collections.Chapter, form *ChapterUpdateForm, userUuid string) (int, error) {
	if story.UserUuid != userUuid {
		return http.StatusForbidden, errors.New("forbidden")
	}

	if chapter.UserUuid != userUuid {
		return http.StatusForbidden, errors.New("forbidden")
	}

	if chapter.StoryId != story.Id {
		return http.StatusForbidden, errors.New("invalid story and chapter")
	}

	if form.Body == "" {
		return http.StatusUnprocessableEntity, errors.New("body is required")
	}

	return 0, nil
}

func validateChapterViewRequest(story *collections.Story, chapter *collections.Chapter) (int, error) {
	if chapter.StoryId != story.Id {
		return http.StatusForbidden, errors.New("invalid story and chapter")
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
