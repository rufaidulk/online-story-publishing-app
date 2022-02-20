package v1

import (
	"errors"
	"mime/multipart"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryForm struct {
	Title        string
	Body         string
	LanguageCode string `json:"language_code"`
	Categories   []string
	IsSeries     bool `json:"is_series"`
	IsPremium    bool `json:"is_premium"`
}

type StoryPromotionalInfoForm struct {
	PromotionalTitle string                `form:"promotional_title"` //optional
	PromotionalImage *multipart.FileHeader `form:"promotional_image"` //optional
}

type StoryResponse struct {
	Id               primitive.ObjectID          `json:"_id"`
	UserUuid         string                      `json:"user_uuid"`
	Slug             string                      `json:"slug"`
	Title            string                      `json:"title"`
	PromotionalTitle string                      `json:"promotional_title"`
	PromotionalImage string                      `json:"promotional_image"`
	LanguageCode     string                      `json:"language_code"`
	Categories       []primitive.ObjectID        `json:"categories"`
	Chapters         map[int]ChapterInfoResponse `json:"chapters"`
	IsPremium        bool                        `json:"is_premium"`
	IsCompleted      bool                        `json:"is_completed"`
	Rating           int8                        `json:"rating"`
	AvgReadCount     int64                       `json:"avg_read_count"`
}

type ChapterInfoResponse struct {
	ChapterId    primitive.ObjectID `json:"chapter_id"`
	ChapterTitle string             `json:"chapter_title"`
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
	storyResponse := buildStoryResponse(story)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "story created", storyResponse))
}

func UpdateStoryPromotionalInfo(ctx echo.Context) error {
	form := new(StoryPromotionalInfoForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	if form.PromotionalTitle == "" && form.PromotionalImage == nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "invalid request"))
	}
	storyId := ctx.Param("id")
	story := collections.NewStory()
	if err := story.LoadById(storyId); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, "invalid story"))
	}
	// Source
	file, err := ctx.FormFile("promotional_image")
	if err != nil {
		return err
	}
	allowedImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedImageTypes[file.Header.Get("content-type")] {
		return ctx.JSON(http.StatusUnsupportedMediaType,
			helper.NewErrorResponse(http.StatusUnsupportedMediaType, "image type not supported"))
	}
	fileName, err := helper.FileUpload(file)
	if err != nil {
		return err
	}
	oldFile := story.PromotionalImage
	story.PromotionalTitle = form.PromotionalTitle
	story.PromotionalImage = fileName
	story.Update()
	if oldFile != "" {
		defer helper.FileDelete(oldFile)
	}

	storyResponse := buildStoryResponse(story)
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusOK, "story promotional info updated", storyResponse))
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

func buildStoryResponse(story *collections.Story) *StoryResponse {
	infoRes := make(map[int]ChapterInfoResponse)
	for k, v := range story.Chapters {
		info := ChapterInfoResponse{
			ChapterId:    v.ChapterId,
			ChapterTitle: v.ChapterTitle,
		}
		infoRes[k] = info
	}

	res := StoryResponse{
		Id:               story.Id,
		UserUuid:         story.UserUuid,
		Slug:             story.Slug,
		Title:            story.Title,
		PromotionalTitle: story.PromotionalTitle,
		PromotionalImage: story.PromotionalImage,
		LanguageCode:     story.LanguageCode,
		Categories:       story.Categories,
		Chapters:         infoRes,
		IsPremium:        story.IsPremium,
		IsCompleted:      story.IsCompleted,
		Rating:           story.Rating,
		AvgReadCount:     story.AvgReadCount,
	}

	return &res
}
