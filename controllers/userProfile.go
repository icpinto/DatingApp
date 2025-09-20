package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func filterEmptyStrings(values []string) []string {
	var result []string
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			result = append(result, v)
		}
	}
	return result
}

// CreateProfile godoc
// @Summary      Create or update the authenticated user's profile
// @Description  Updates the profile information for the authenticated user. Supports multipart form data with optional profile image upload.
// @Tags         Profiles
// @Accept       multipart/form-data
// @Produce      json
// @Param        bio                   formData  string false "Biography"
// @Param        gender                formData  string false "Gender"
// @Param        date_of_birth         formData  string false "Date of birth (YYYY-MM-DD)"
// @Param        location              formData  string false "Location"
// @Param        interests             formData  string false "Interests (can be repeated)"
// @Param        civil_status          formData  string false "Civil status"
// @Param        religion              formData  string false "Religion"
// @Param        dietary_preference    formData  string false "Dietary preference"
// @Param        smoking               formData  string false "Smoking habit"
// @Param        alcohol               formData  string false "Alcohol habit"
// @Param        languages             formData  string false "Languages (can be repeated)"
// @Param        highest_education     formData  string false "Highest education"
// @Param        employment_status     formData  string false "Employment status"
// @Param        occupation            formData  string false "Occupation"
// @Param        siblings_count        formData  int    false "Number of siblings"
// @Param        horoscope_available   formData  bool   false "Whether a horoscope is available"
// @Param        profile_image         formData  file   false "Profile image"
// @Success      200                   {object}  utils.MessageResponse
// @Failure      400                   {object}  utils.ErrorResponse
// @Failure      401                   {object}  utils.ErrorResponse
// @Failure      500                   {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/profile [post]
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
	profile.Interests = filterEmptyStrings(ctx.PostFormArray("interests"))
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
	profile.Languages = filterEmptyStrings(ctx.PostFormArray("languages"))
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

// GetProfile godoc
// @Summary      Retrieve the authenticated user's profile
// @Tags         Profiles
// @Produce      json
// @Success      200  {object}  models.UserProfile
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/profile [get]
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

// GetProfiles godoc
// @Summary      List user profiles with optional filters
// @Tags         Profiles
// @Produce      json
// @Param        gender               query     string false "Filter by gender"
// @Param        civil_status         query     string false "Filter by civil status"
// @Param        religion             query     string false "Filter by religion"
// @Param        dietary_preference   query     string false "Filter by dietary preference"
// @Param        smoking              query     string false "Filter by smoking habit"
// @Param        country_code         query     string false "Filter by country code"
// @Param        highest_education    query     string false "Filter by education"
// @Param        employment_status    query     string false "Filter by employment status"
// @Param        age                  query     int    false "Filter by age"
// @Param        horoscope_available  query     bool   false "Filter by horoscope availability"
// @Success      200                  {array}   models.UserProfile
// @Failure      400                  {object}  utils.ErrorResponse
// @Failure      500                  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/profiles [get]
func GetProfiles(ctx *gin.Context) {
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	filters := models.ProfileFilters{
		Gender:            ctx.Query("gender"),
		CivilStatus:       ctx.Query("civil_status"),
		Religion:          ctx.Query("religion"),
		DietaryPreference: ctx.Query("dietary_preference"),
		Smoking:           ctx.Query("smoking"),
		CountryCode:       ctx.Query("country_code"),
		HighestEducation:  ctx.Query("highest_education"),
		EmploymentStatus:  ctx.Query("employment_status"),
	}

	if ageStr := ctx.Query("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, err, "GetProfiles invalid age filter", "Invalid age filter")
			return
		}
		filters.Age = &age
	}

	if horoscopeStr := ctx.Query("horoscope_available"); horoscopeStr != "" {
		horoscope, err := strconv.ParseBool(horoscopeStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, err, "GetProfiles invalid horoscope filter", "Invalid horoscope filter")
			return
		}
		filters.HoroscopeAvailable = &horoscope
	}

	profiles, err := profileService.GetProfiles(filters)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetProfiles service error", "Failed to retrieve profiles")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, profiles)
}

// GetUserProfile godoc
// @Summary      Retrieve a user profile by ID
// @Tags         Profiles
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Success      200      {object}  models.UserProfile
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/profile/{user_id} [get]
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
// GetProfileEnums godoc
// @Summary      Retrieve supported enum values for profile fields
// @Tags         Profiles
// @Produce      json
// @Success      200  {object}  models.ProfileEnums
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /profile/enums [get]
func GetProfileEnums(ctx *gin.Context) {
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	enums, err := profileService.GetProfileEnums()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetProfileEnums service error", "Failed to retrieve enums")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, enums)
}
