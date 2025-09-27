package controllers_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
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

func setupProfileRouter(db *sql.DB, matchService *services.MatchService, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	profileService := services.NewProfileService(db)
	r.Use(middlewares.ServiceMiddleware(middlewares.Services{ProfileService: profileService, MatchService: matchService}))
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
		"languages", "phone_number", "contact_verified", "identity_verified", "country_code", "province", "district", "city", "postal_code", "highest_education", "field_of_study",
		"institution", "employment_status", "occupation", "father_occupation", "mother_occupation", "siblings_count", "siblings",
		"horoscope_available", "birth_time", "birth_place", "sinhala_raasi", "nakshatra", "horoscope",
		"profile_image_url", "profile_image_thumb_url", "verified", "moderation_status", "last_active_at", "metadata",
		"created_at", "updated_at",
	}
	now := time.Now()
	row := make([]driver.Value, len(columns))
	for i, column := range columns {
		switch column {
		case "id", "user_id":
			row[i] = 1
		case "username":
			row[i] = "john"
		case "interests", "languages":
			row[i] = pq.StringArray{}
		case "height_cm", "weight_kg", "siblings_count":
			row[i] = 0
		case "horoscope_available", "contact_verified", "identity_verified", "verified":
			row[i] = false
		case "created_at", "updated_at":
			row[i] = now
		default:
			row[i] = ""
		}
	}
	return sqlmock.NewRows(columns).AddRow(row...)
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
	mock.ExpectQuery("SELECT COALESCE\\(phone_number, ''\\), COALESCE\\(contact_verified, false\\), COALESCE\\(identity_verified, false\\), COALESCE\\(verified, false\\) FROM profiles WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"phone_number", "contact_verified", "identity_verified", "verified"}))

	args := make([]driver.Value, 45)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}
	mock.ExpectExec("INSERT INTO profiles").
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT p.id, p.user_id").
		WithArgs(1).
		WillReturnRows(mockProfileRows())

	payloadCh := make(chan map[string]any, 1)
	errCh := make(chan error, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errCh <- fmt.Errorf("failed to read request body: %w", err)
			http.Error(w, "bad request", http.StatusInternalServerError)
			return
		}
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			errCh <- fmt.Errorf("failed to unmarshal match payload: %w", err)
			http.Error(w, "bad request", http.StatusInternalServerError)
			return
		}
		payloadCh <- payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"userId":1}`))
	}))
	defer server.Close()

	router := setupProfileRouter(db, services.NewMatchService(server.URL), true)

	form := url.Values{}
	form.Set("bio", "hello")
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}

	select {
	case err := <-errCh:
		t.Fatalf("match service request error: %v", err)
	default:
	}

	select {
	case payload := <-payloadCh:
		if _, found := payload["id"]; found {
			t.Fatalf("expected match payload to omit id field, got: %v", payload)
		}
		if userID, found := payload["user_id"]; !found || userID != float64(1) {
			t.Fatalf("expected user_id 1 in match payload, got: %v", payload["user_id"])
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for match service payload")
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"userId":1}`))
	}))
	defer server.Close()

	router := setupProfileRouter(db, services.NewMatchService(server.URL), true)
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"userId":1}`))
	}))
	defer server.Close()

	router := setupProfileRouter(db, services.NewMatchService(server.URL), false)
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

func TestGetProfilesWithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT p.id, p.user_id").
		WithArgs("female", 30, true).
		WillReturnRows(mockProfileRows())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"userId":1}`))
	}))
	defer server.Close()

	router := setupProfileRouter(db, services.NewMatchService(server.URL), false)
	req := httptest.NewRequest(http.MethodGet, "/profiles?gender=female&age=30&horoscope_available=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestGetProfilesInvalidAgeFilter(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"userId":1}`))
	}))
	defer server.Close()

	router := setupProfileRouter(db, services.NewMatchService(server.URL), false)
	req := httptest.NewRequest(http.MethodGet, "/profiles?age=notanumber", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 got %d: %s", w.Code, w.Body.String())
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"userId":1}`))
	}))
	defer server.Close()

	router := setupProfileRouter(db, services.NewMatchService(server.URL), false)
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
