package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
)

func setupRequestRouter(db *sql.DB, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	frService := services.NewFriendRequestService(db)
	r.Use(middlewares.ServiceMiddleware(middlewares.Services{FriendRequestService: frService}))
	if withUser {
		r.Use(func(c *gin.Context) {
			c.Set("username", "john")
			c.Next()
		})
	}
	r.POST("/sendRequest", controllers.SendFriendRequest)
	r.POST("/acceptRequest", controllers.AcceptFriendRequest)
	r.POST("/rejectRequest", controllers.RejectFriendRequest)
	r.GET("/requests", controllers.GetPendingRequests)
	r.GET("/sentRequests", controllers.GetSentRequests)
	r.GET("/checkReqStatus/:reciver_id", controllers.CheckReqStatus)
	return r
}

func TestSendFriendRequestSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("jane"))
	mock.ExpectQuery("SELECT status FROM friend_requests WHERE sender_id = \\$1 AND receiver_id = \\$2").
		WithArgs(1, 2).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO friend_requests \\(sender_id, sender_username, receiver_id, receiver_username, status, created_at, updated_at\\)").
		WithArgs(1, "john", 2, "jane", "pending", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := setupRequestRouter(db, true)

	body, _ := json.Marshal(models.FriendRequest{ReceiverID: 2})
	req := httptest.NewRequest(http.MethodPost, "/sendRequest", bytes.NewReader(body))
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

func TestSendFriendRequestDuplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("jane"))
	mock.ExpectQuery("SELECT status FROM friend_requests WHERE sender_id = \\$1 AND receiver_id = \\$2").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("pending"))

	router := setupRequestRouter(db, true)

	body, _ := json.Marshal(models.FriendRequest{ReceiverID: 2})
	req := httptest.NewRequest(http.MethodPost, "/sendRequest", bytes.NewReader(body))
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

func TestAcceptFriendRequestSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE friend_requests SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs("accepted", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT sender_id, receiver_id FROM friend_requests WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"sender_id", "receiver_id"}).AddRow(1, 2))
	mock.ExpectExec("INSERT INTO conversation_outbox \\(event_id, user1_id, user2_id, processed, created_at\\) VALUES \\(\\$1, \\$2, \\$3, false, \\$4\\)").
		WithArgs(sqlmock.AnyArg(), 1, 2, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := setupRequestRouter(db, false)

	body, _ := json.Marshal(models.AcceptRequest{RequestID: 1})
	req := httptest.NewRequest(http.MethodPost, "/acceptRequest", bytes.NewReader(body))
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

func TestRejectFriendRequestSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE friend_requests SET status = \$1, updated_at = \$2 WHERE id = \$3`).
		WithArgs("rejected", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupRequestRouter(db, false)

	body, _ := json.Marshal(models.RejectRequest{RequestID: 1})
	req := httptest.NewRequest(http.MethodPost, "/rejectRequest", bytes.NewReader(body))
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

func TestGetPendingRequestsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT id, sender_id, sender_username, receiver_id, receiver_username, status, created_at FROM friend_requests WHERE receiver_id = \\$1 AND status = 'pending'").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id", "sender_username", "receiver_id", "receiver_username", "status", "created_at"}).
			AddRow(10, 2, "", 1, "", "pending", time.Now()))
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("alice"))
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("john"))

	router := setupRequestRouter(db, true)

	req := httptest.NewRequest(http.MethodGet, "/requests", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestGetSentRequestsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT id, sender_id, sender_username, receiver_id, receiver_username, status, created_at, updated_at FROM friend_requests WHERE sender_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id", "sender_username", "receiver_id", "receiver_username", "status", "created_at", "updated_at"}).
			AddRow(11, 1, "", 2, "", "accepted", now, now))
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("john"))
	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("alice"))

	router := setupRequestRouter(db, true)

	req := httptest.NewRequest(http.MethodGet, "/sentRequests", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}

func TestCheckReqStatusSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM users WHERE username=\\$1").
		WithArgs("john").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM friend_requests WHERE sender_id = \$1 AND receiver_id = \$2`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	router := setupRequestRouter(db, true)

	req := httptest.NewRequest(http.MethodGet, "/checkReqStatus/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
	}
}
