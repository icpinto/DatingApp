package services

import (
	"database/sql"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func GetQuestionnaire(db *sql.DB) ([]models.Question, error) {
	return repositories.GetQuestions(db)
}

func SubmitAnswer(db *sql.DB, username string, answer models.Answer) error {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return err
	}
	return repositories.UpsertAnswer(db, userID, answer)
}

func GetUserAnswers(db *sql.DB, username string) ([]models.Answer, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return nil, err
	}
	return repositories.GetAnswersByUserID(db, userID)
}
