package controllers

import (
	"net/http"
	"userservice/adapters"
	"userservice/helper"
	"userservice/models"

	"github.com/labstack/echo/v4"
)

type UserProfileUpdateForm struct {
	Name        string
	Password    string
	PenName     string
	Description string
}

func ViewUserProfile(ctx echo.Context) error {
	userDetails := ctx.Get("userDetails").(helper.UserDetails)
	db := adapters.GetDbHandle(ctx)
	sql := "SELECT uuid_text as uuid, name, email, pen_name, description, is_author, is_premium_author, followers_count, followee_count FROM users "
	sql += "INNER JOIN user_profiles ON user_profiles.user_id = users.id WHERE email = ? AND uuid = ?"
	profileRes := UserProfileResponse{}
	uuid := models.NewUuid(userDetails.Uuid)
	err := db.Raw(sql, userDetails.Email, uuid).Scan(&profileRes).Error
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK,
		helper.NewSuccessResponse(http.StatusOK, "profile details", profileRes))
}

func UpdateUserProfile(ctx echo.Context) error {
	db := adapters.GetDbHandle(ctx)
	form := new(UserProfileUpdateForm)
	if err := ctx.Bind(form); err != nil {
		return err
	}
	userDetails := ctx.Get("userDetails").(helper.UserDetails)
	user := models.NewUserData()
	uuid := models.NewUuid(userDetails.Uuid)
	if err := db.Where("uuid = ?", uuid).Take(&user).Error; err != nil {
		return err
	}
	userProfile := models.NewUserProfile()
	if err := db.Where("user_id = ?", user.Id).Take(&userProfile).Error; err != nil {
		return err
	}
	data := map[string]interface{}{
		"name": form.Name,
	}
	tx := db.Begin()
	user.Name = form.Name
	if form.Password != "" {
		passwd, err := helper.HashPassword(form.Password)
		if err != nil {
			return err
		}
		data["password"] = passwd
	}
	if err := tx.Model(&user).Updates(data).Error; err != nil {
		return err
	}
	userProfile.PenName = form.PenName
	userProfile.Description = form.Description
	if err := tx.Save(&userProfile).Error; err != nil {
		return err
	}

	tx.Commit()
	profileRes := buildUserProfileResponse(userDetails, userProfile)
	return ctx.JSON(http.StatusOK,
		helper.NewSuccessResponse(http.StatusOK, "profile details updated", profileRes))
}
