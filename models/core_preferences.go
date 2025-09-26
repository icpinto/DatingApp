package models

// CorePreferences represents the matching preferences for a user.
type CorePreferences struct {
	UserID             int    `json:"user_id"`
	MinAge             int    `json:"min_age"`
	MaxAge             int    `json:"max_age"`
	Gender             string `json:"gender"`
	DrinkingHabit      string `json:"drinking_habit"`
	EducationLevel     string `json:"education_level"`
	SmokingHabit       string `json:"smoking_habit"`
	CountryOfResidence string `json:"country_of_residence"`
	OccupationStatus   string `json:"occupation_status"`
	CivilStatus        string `json:"civil_status"`
	Religion           string `json:"religion"`
	MinHeight          int    `json:"min_height"`
	MaxHeight          int    `json:"max_height"`
	FoodPreference     string `json:"food_preference"`
}
