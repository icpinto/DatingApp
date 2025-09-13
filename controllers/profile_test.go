package controllers_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/services"
	"github.com/lib/pq"
)

func setupProfileRouter(db *sql.DB, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	profileService := services.NewProfileService(db)
	r.Use(middlewares.ServiceMiddleware(middlewares.Services{ProfileService: profileService}))
	if withUser {
		r.Use(func(c *gin.Context) {
			c.Set("username", "john")
			c.Next()
		})
	}
	r.POST("/profile", controllers.CreateProfile)
	r.GET("/profile", controllers.GetProfile)
	r.GET("/profiles", controllers.GetProfiles)
	r.GET("/profile/:user_id", controllers.GetUserProfile)
	return r
}

func mockProfileRows() *sqlmock.Rows {
	columns := []string{
		"id", "user_id", "username", "bio", "gender", "date_of_birth", "location_legacy", "interests", "civil_status",
		"religion", "religion_detail", "caste", "height_cm", "weight_kg", "dietary_preference", "smoking", "alcohol",
		"languages", "country_code", "province", "district", "city", "postal_code", "highest_education", "field_of_study",
		"institution", "employment_status", "occupation", "father_occupation", "mother_occupation", "siblings_count", "siblings",
		"horoscope_available", "birth_time", "birth_place", "sinhala_raasi", "nakshatra", "horoscope",
		"profile_image_url", "profile_image_thumb_url", "verified", "moderation_status", "last_active_at", "metadata",
		"created_at", "updated_at",
	}
	now := time.Now()
	return sqlmock.NewRows(columns).AddRow(
		1, 1, "john", "", "", "", "", pq.StringArray{}, "", "", "", "", 0, 0, "", "", "", pq.StringArray{}, "", "", "", "", "", "", "", "", "", "", "", "", 0, "",
		false, "", "", "", "", "", "", "", false, "", "", "", now, now,
	)
}

func TestCreateProfileSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	args := make([]driver.Value, 42)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}
	mock.ExpectExec("INSERT INTO profiles").
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := setupProfileRouter(db, true)

	form := url.Values{}
	form.Set("bio", "hello")
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestGetProfileSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT p.id, p.user_id").
		WithArgs(1).
		WillReturnRows(mockProfileRows())

	router := setupProfileRouter(db, true)
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestGetProfilesSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT p.id, p.user_id").
		WillReturnRows(mockProfileRows())

	router := setupProfileRouter(db, false)
	req := httptest.NewRequest(http.MethodGet, "/profiles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestGetUserProfileSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT p.id, p.user_id").
		WithArgs(2).
		WillReturnRows(mockProfileRows())

	router := setupProfileRouter(db, false)
	req := httptest.NewRequest(http.MethodGet, "/profile/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}
