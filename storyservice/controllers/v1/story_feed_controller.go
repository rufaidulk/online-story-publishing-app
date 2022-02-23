package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
)

type StoryFeedInterestedCategoriesForm struct {
	Categories []string
}

func FetchStoryFeed(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	storyFeed := collections.NewStoryFeed()
	storyFeed.LoadByUser(userUuid)
	res := make(map[string]interface{})
	if !storyFeed.ReadStatus.StoryId.IsZero() && !storyFeed.ReadStatus.ChapterId.IsZero() {
		currentStory := collections.NewStory()
		if err := currentStory.LoadByObjectId(storyFeed.ReadStatus.StoryId); err != nil {
			return err
		}
		currentChapter := collections.NewChapter()
		if err := currentChapter.LoadByObjectId(storyFeed.ReadStatus.ChapterId); err != nil {
			return err
		}
		res["current_story"] = buildStoryResponse(currentStory)
		res["current_chapter"] = buildChapterResponse(currentChapter)
	}
	data, err := storyFeed.FetchRecommendedStories()
	if err != nil {
		return err
	}
	res["recommendations"] = data
	return ctx.JSON(http.StatusCreated, helper.NewSuccessResponse(http.StatusOK, "story feed", res))
}

func CreateInterestedCategoriesInStoryFeed(ctx echo.Context) error {
	userUuid := ctx.Get("userUuid").(string)
	form := new(StoryFeedInterestedCategoriesForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	if err := validateStoryFeedInterestedCategoriesForm(form); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	StoryFeed := collections.NewStoryFeed()
	StoryFeed.UserUuid = userUuid
	StoryFeed.SetCategories(form.Categories)
	if err := StoryFeed.UpsertDocument(); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, helper.NewSuccessResponse(http.StatusCreated, "interested categories saved", ""))
}

func validateStoryFeedInterestedCategoriesForm(form *StoryFeedInterestedCategoriesForm) error {
	if len(form.Categories) == 0 {
		return errors.New("at least one category is required")
	}

	category := collections.NewCategory()
	if err := category.CheckExistsAllByIds(form.Categories); err != nil {
		return errors.New("invalid categories")
	}

	return nil
}
