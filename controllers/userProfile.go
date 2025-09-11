package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func CreateProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "CreateProfile unauthorized", "Unauthorized")
		return
	}

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	var profile models.Profile
	profile.Bio = ctx.PostForm("bio")
	profile.Gender = ctx.PostForm("gender")
	profile.DateOfBirth = ctx.PostForm("date_of_birth")
	profile.LocationLegacy = ctx.PostForm("location")
	profile.Interests = ctx.PostFormArray("interests")
	profile.CivilStatus = ctx.PostForm("civil_status")
	profile.Religion = ctx.PostForm("religion")
	profile.ReligionDetail = ctx.PostForm("religion_detail")
	profile.Caste = ctx.PostForm("caste")
	if v, err := strconv.Atoi(ctx.DefaultPostForm("height_cm", "0")); err == nil {
		profile.HeightCM = v
	}
	if v, err := strconv.Atoi(ctx.DefaultPostForm("weight_kg", "0")); err == nil {
		profile.WeightKG = v
	}
	profile.DietaryPreference = ctx.PostForm("dietary_preference")
	profile.Smoking = ctx.PostForm("smoking")
	profile.Alcohol = ctx.PostForm("alcohol")
	profile.Languages = ctx.PostFormArray("languages")
	profile.CountryCode = ctx.DefaultPostForm("country_code", "LK")
	profile.Province = ctx.PostForm("province")
	profile.District = ctx.PostForm("district")
	profile.City = ctx.PostForm("city")
	profile.PostalCode = ctx.PostForm("postal_code")
	profile.HighestEducation = ctx.PostForm("highest_education")
	profile.FieldOfStudy = ctx.PostForm("field_of_study")
	profile.Institution = ctx.PostForm("institution")
	profile.EmploymentStatus = ctx.PostForm("employment_status")
	profile.Occupation = ctx.PostForm("occupation")
	profile.FatherOccupation = ctx.PostForm("father_occupation")
	profile.MotherOccupation = ctx.PostForm("mother_occupation")
	if v, err := strconv.Atoi(ctx.DefaultPostForm("siblings_count", "0")); err == nil {
		profile.SiblingsCount = v
	}
	profile.Siblings = ctx.PostForm("siblings")
	if v, err := strconv.ParseBool(ctx.DefaultPostForm("horoscope_available", "false")); err == nil {
		profile.HoroscopeAvailable = v
	}
	profile.BirthTime = ctx.PostForm("birth_time")
	profile.BirthPlace = ctx.PostForm("birth_place")
	profile.SinhalaRaasi = ctx.PostForm("sinhala_raasi")
	profile.Nakshatra = ctx.PostForm("nakshatra")
	profile.Horoscope = ctx.PostForm("horoscope")

	file, err := ctx.FormFile("profile_image")
	if err == nil {
		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err, "CreateProfile mkdir error", "Failed to save image")
			return
		}
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		dst := filepath.Join("uploads", filename)
		if err := ctx.SaveUploadedFile(file, dst); err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err, "CreateProfile save error", "Failed to save image")
			return
		}
		url := fmt.Sprintf("http://%s/uploads/%s", ctx.Request.Host, filename)
		profile.ProfileImageURL = url
		profile.ProfileImageThumbURL = url
	}

	if err := profileService.CreateOrUpdateProfile(username.(string), profile); err != nil {
		logMsg := fmt.Sprintf("CreateProfile service error for %s", username.(string))
		status := http.StatusInternalServerError
		clientMsg := "Failed to update profile"
		if errors.Is(err, services.ErrInvalidEnum) {
			status = http.StatusBadRequest
			clientMsg = "Invalid profile data"
		}
		utils.RespondError(ctx, status, err, logMsg, clientMsg)
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetProfile unauthorized", "Unauthorized")
		return
	}

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	profile, err := profileService.GetProfile(username.(string))
	if err != nil {
		logMsg := fmt.Sprintf("GetProfile service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve profile")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, profile)
}

func GetProfiles(ctx *gin.Context) {
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	profiles, err := profileService.GetProfiles()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetProfiles service error", "Failed to retrieve profiles")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, profiles)
}

func GetUserProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "GetUserProfile invalid user id", "Invalid user id")
		return
	}

	profile, err := profileService.GetProfileByUserID(userID)
	if err != nil {
		logMsg := fmt.Sprintf("GetUserProfile service error for user %d", userID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve profile")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, profile)
}

// GetProfileEnums returns enum values for profile-related fields.
func GetProfileEnums(ctx *gin.Context) {
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	enums, err := profileService.GetProfileEnums()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetProfileEnums service error", "Failed to retrieve enums")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, enums)
}
