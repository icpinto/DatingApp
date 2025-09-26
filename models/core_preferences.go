package models

// CorePreferences represents the matching preferences for a user.
type CorePreferences struct {
	UserID             int    `json:"user_id"`
	MinAge             int    `json:"minAge"`
	MaxAge             int    `json:"maxAge"`
	Gender             string `json:"gender"`
	DrinkingHabit      string `json:"drinkingHabit"`
	EducationLevel     string `json:"educationLevel"`
	SmokingHabit       string `json:"smokingHabit"`
	CountryOfResidence string `json:"countryOfResidence"`
	OccupationStatus   string `json:"occupationStatus"`
	CivilStatus        string `json:"civilStatus"`
	Religion           string `json:"religion"`
	MinHeight          int    `json:"minHeight"`
	MaxHeight          int    `json:"maxHeight"`
	FoodPreference     string `json:"foodPreference"`
}
