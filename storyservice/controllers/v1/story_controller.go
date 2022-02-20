package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
)

type StoryForm struct {
	Title        string
	Body         string
	LanguageCode string `json:"language_code"`
	Categories   []string
	IsSeries     bool `json:"is_series"`
	IsPremium    bool `json:"is_premium"`
}

func CreateStory(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	form := new(StoryForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	if err := validateStoryForm(form, userUuid); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	story := collections.NewStory()
	story.UserUuid = userUuid
	story.Title = form.Title
	story.LanguageCode = form.LanguageCode
	story.SetCategories(form.Categories)
	if !form.IsSeries {
		story.IsCompleted = true
	}
	//todo:: check the eligibility to write premium stories
	story.IsPremium = form.IsPremium
	if err := story.CreateDocument(); err != nil {
		return err
	}

	chapter := collections.NewChapter()
	chapter.UserUuid = userUuid
	chapter.StoryId = story.Id
	chapter.Title = form.Title
	chapter.Body = form.Body
	if err := chapter.CreateDocument(); err != nil {
		return err
	}
	if err := story.AddChapter(chapter.Id, chapter.Title); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "story created", story))
}

func validateStoryForm(form *StoryForm, userUuid string) error {
	if form.Title == "" {
		return errors.New("title is required")
	}
	if form.Body == "" {
		return errors.New("body is required")
	}
	if form.LanguageCode == "" {
		return errors.New("language code is required")
	}
	if len(form.Categories) == 0 {
		return errors.New("at least one category is required")
	}

	category := collections.NewCategory()
	if err := category.CheckExistsAllByIds(form.Categories); err != nil {
		return errors.New("invalid categories")
	}

	return nil
}
