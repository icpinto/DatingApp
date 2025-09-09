package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func GetQuestionnaire(db *sql.DB) ([]models.Question, error) {
	questions, err := repositories.GetQuestions(db)
	if err != nil {
		log.Printf("GetQuestionnaire service error: %v", err)
		return nil, err
	}
	return questions, nil
}

func SubmitAnswer(db *sql.DB, username string, answer models.Answer) error {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("SubmitAnswer user lookup error for %s: %v", username, err)
		return err
	}
	if err := repositories.UpsertAnswer(db, userID, answer); err != nil {
		log.Printf("SubmitAnswer repository error for user %d question %d: %v", userID, answer.QuestionID, err)
		return err
	}
	return nil
}

func GetUserAnswers(db *sql.DB, username string) ([]models.Answer, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("GetUserAnswers user lookup error for %s: %v", username, err)
		return nil, err
	}
	answers, err := repositories.GetAnswersByUserID(db, userID)
	if err != nil {
		log.Printf("GetUserAnswers repository error for user %d: %v", userID, err)
		return nil, err
	}
	return answers, nil
}
