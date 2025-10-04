package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func setupRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	userService := services.NewUserService(db)
	r.Use(middlewares.ServiceMiddleware(middlewares.Services{UserService: userService}))
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	return r
}

func TestRegisterSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", "john@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := setupRouter(db)

	body, _ := json.Marshal(models.User{Username: "john", Email: "john@example.com", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestRegisterDuplicateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", "john@example.com", sqlmock.AnyArg()).
		WillReturnError(repositories.ErrDuplicateUser)

	router := setupRouter(db)

	body, _ := json.Marshal(models.User{Username: "john", Email: "john@example.com", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	hashed, err := utils.HashPassword("pass")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}

	mock.ExpectQuery("SELECT id, password FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(1, hashed))

	router := setupRouter(db)

	body := []byte(`{"username":"john","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("token")) {
		t.Fatalf("expected token in response: %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestLoginUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, password FROM users WHERE username=\\$1").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	router := setupRouter(db)

	body := []byte(`{"username":"missing","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	hashed, err := utils.HashPassword("pass")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}

	mock.ExpectQuery("SELECT id, password FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(1, hashed))

	router := setupRouter(db)

	body := []byte(`{"username":"john","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestReactivateAllowsInactiveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1 AND is_active = true").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectBegin()
	mock.ExpectExec("\\s*UPDATE users\\s+SET is_active = true, deactivated_at = NULL\\s+WHERE id = \\$1 AND is_active = false").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("\\s*INSERT INTO user_lifecycle_outbox").
		WithArgs(sqlmock.AnyArg(), 1, "reactivated", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	userService := services.NewUserService(db)
	router.Use(middlewares.ServiceMiddleware(middlewares.Services{UserService: userService}))
	router.POST("/user/reactivate", middlewares.Authenticate, controllers.ReactivateCurrentUser)

	token, err := utils.GenerateToken(1)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	body := []byte(`{"reason":"let me back"}`)
	req := httptest.NewRequest(http.MethodPost, "/user/reactivate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}
