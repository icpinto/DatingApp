package middlewares_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func setupAuthRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	userService := services.NewUserService(db)
	router.Use(middlewares.ServiceMiddleware(middlewares.Services{UserService: userService}))
	return router
}

func TestAuthenticateAllowsInactiveUserForRequestListingRoutes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1 AND is_active = true").
		WithArgs(42).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("inactive_jane"))

	router := setupAuthRouter(db)
	router.GET("/user/requests", middlewares.Authenticate, func(c *gin.Context) {
		username := c.GetString("username")
		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	token, err := utils.GenerateToken(42)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/user/requests", nil)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "inactive_jane") {
		t.Fatalf("expected response to include username, got: %s", w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestAuthenticateBlocksInactiveUserForDisallowedReadRoutes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1 AND is_active = true").
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	router := setupAuthRouter(db)
	handlerCalled := false
	router.GET("/user/profile", middlewares.Authenticate, func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	token, err := utils.GenerateToken(99)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 got %d: %s", w.Code, w.Body.String())
	}
	if handlerCalled {
		t.Fatalf("handler should not be called for disallowed inactive read access")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestAuthenticateBlocksInactiveUserForWriteRoutes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1 AND is_active = true").
		WithArgs(7).
		WillReturnError(sql.ErrNoRows)

	router := setupAuthRouter(db)
	handlerCalled := false
	router.POST("/user/sendRequest", middlewares.Authenticate, func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	token, err := utils.GenerateToken(7)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user/sendRequest", nil)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 got %d: %s", w.Code, w.Body.String())
	}
	if handlerCalled {
		t.Fatalf("handler should not be called for inactive write access")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}
