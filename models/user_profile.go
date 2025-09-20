package models

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}

type Profile struct {
	ID                   int      `json:"id"`
	UserID               int      `json:"user_id"` // Foreign key to users table
	Bio                  string   `json:"bio"`
	Gender               string   `json:"gender"`
	DateOfBirth          string   `json:"date_of_birth"`
	LocationLegacy       string   `json:"location"`
	Interests            []string `json:"interests"` // Array of interests
	CivilStatus          string   `json:"civil_status"`
	Religion             string   `json:"religion"`
	ReligionDetail       string   `json:"religion_detail"`
	Caste                string   `json:"caste"`
	HeightCM             int      `json:"height_cm"`
	WeightKG             int      `json:"weight_kg"`
	DietaryPreference    string   `json:"dietary_preference"`
	Smoking              string   `json:"smoking"`
	Alcohol              string   `json:"alcohol"`
	Languages            []string `json:"languages"`
	PhoneNumber          string   `json:"phone_number"`
	ContactVerified      bool     `json:"contact_verified"`
	IdentityVerified     bool     `json:"identity_verified"`
	CountryCode          string   `json:"country_code"`
	Province             string   `json:"province"`
	District             string   `json:"district"`
	City                 string   `json:"city"`
	PostalCode           string   `json:"postal_code"`
	HighestEducation     string   `json:"highest_education"`
	FieldOfStudy         string   `json:"field_of_study"`
	Institution          string   `json:"institution"`
	EmploymentStatus     string   `json:"employment_status"`
	Occupation           string   `json:"occupation"`
	FatherOccupation     string   `json:"father_occupation"`
	MotherOccupation     string   `json:"mother_occupation"`
	SiblingsCount        int      `json:"siblings_count"`
	Siblings             string   `json:"siblings"`
	HoroscopeAvailable   bool     `json:"horoscope_available"`
	BirthTime            string   `json:"birth_time"`
	BirthPlace           string   `json:"birth_place"`
	SinhalaRaasi         string   `json:"sinhala_raasi"`
	Nakshatra            string   `json:"nakshatra"`
	Horoscope            string   `json:"horoscope"`
	ProfileImageURL      string   `json:"profile_image_url"`
	ProfileImageThumbURL string   `json:"profile_image_thumb_url"`
	Verified             bool     `json:"verified"`
	ModerationStatus     string   `json:"moderation_status"`
	LastActiveAt         string   `json:"last_active_at"`
	Metadata             string   `json:"metadata"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}

type UserProfile struct {
	Profile
	Username string `json:"username"`
}

// ProfileVerificationStatus captures persisted verification state for a profile.
type ProfileVerificationStatus struct {
	PhoneNumber      string
	ContactVerified  bool
	IdentityVerified bool
	Verified         bool
}

// ProfileEnums represents available enum values for profile fields.
type ProfileEnums struct {
	CivilStatus       []string `json:"civil_status"`
	DietaryPreference []string `json:"dietary_preference"`
	HabitFrequency    []string `json:"habit_frequency"`
	EducationLevel    []string `json:"education_level"`
	EmploymentStatus  []string `json:"employment_status"`
}

// ProfileFilters represents optional filters when querying profiles.
type ProfileFilters struct {
	Gender             string
	Age                *int
	CivilStatus        string
	Religion           string
	DietaryPreference  string
	Smoking            string
	CountryCode        string
	HighestEducation   string
	EmploymentStatus   string
	HoroscopeAvailable *bool
}
