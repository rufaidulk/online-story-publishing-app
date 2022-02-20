package v1

import (
	"errors"
	"net/http"
	"storyservice/collections"
	"storyservice/helper"

	"github.com/labstack/echo/v4"
)

type CategoryForm struct {
	Name string
}

func CreateCategory(ctx echo.Context) error {
	form := new(CategoryForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}

	if err := validateCategoryForm(form); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity,
			helper.NewErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	category := collections.NewCategory()
	err := category.CreateDocument(form.Name)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	res := make(map[string]string)
	res["id"] = category.Id.Hex()
	res["name"] = category.Name
	return ctx.JSON(http.StatusOK, helper.NewSuccessResponse(http.StatusCreated, "category created", res))
}

func validateCategoryForm(form *CategoryForm) error {
	if form.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
