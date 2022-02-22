package controllers

import (
	"userservice/helper"
	"userservice/models"
)

type UserProfileResponse struct {
	Uuid            string `json:"uuid"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PenName         string `json:"pen_name"`
	Description     string `json:"description"`
	IsAuthor        bool   `json:"is_author"`
	IsPremiumAuthor bool   `json:"is_premium_author"`
	FollowersCount  int64  `json:"followers_count"`
	FolloweeCount   int64  `json:"followee_count"`
}

func buildUserProfileResponse(userDetails helper.UserDetails, userProfile *models.UserProfile) UserProfileResponse {
	res := UserProfileResponse{
		Uuid:            userDetails.Uuid,
		Name:            userDetails.Name,
		Email:           userDetails.Email,
		PenName:         userProfile.PenName,
		Description:     userProfile.Description,
		IsAuthor:        userProfile.IsAuthor,
		IsPremiumAuthor: userProfile.IsPremiumAuthor,
		FollowersCount:  userProfile.FollowersCount,
		FolloweeCount:   userProfile.FolloweeCount,
	}

	return res
}
