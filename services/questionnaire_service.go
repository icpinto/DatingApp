package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

// QuestionnaireService handles questionnaire-related operations.
type QuestionnaireService struct {
	db   *sql.DB
	repo *repositories.QuestionRepository
}

// NewQuestionnaireService creates a new QuestionnaireService.
func NewQuestionnaireService(db *sql.DB) *QuestionnaireService {
	return &QuestionnaireService{db: db, repo: repositories.NewQuestionRepository(db)}
}

// GetQuestionnaire retrieves all questions from the repository.
func (s *QuestionnaireService) GetQuestionnaire() ([]models.Question, error) {
	questions, err := s.repo.GetAll()
	if err != nil {
		log.Printf("GetQuestionnaire service error: %v", err)
		return nil, err
	}
	return questions, nil
}

// SubmitAnswer stores or updates a user's answer.
func (s *QuestionnaireService) SubmitAnswer(username string, answer models.Answer) error {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("SubmitAnswer user lookup error for %s: %v", username, err)
		return err
	}
	if err := s.repo.UpsertAnswer(userID, answer); err != nil {
		log.Printf("SubmitAnswer repository error for user %d question %d: %v", userID, answer.QuestionID, err)
		return err
	}
	return nil
}

// GetUserAnswers retrieves all answers for a given user.
func (s *QuestionnaireService) GetUserAnswers(username string) ([]models.Answer, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetUserAnswers user lookup error for %s: %v", username, err)
		return nil, err
	}
	answers, err := s.repo.GetAnswersByUserID(userID)
	if err != nil {
		log.Printf("GetUserAnswers repository error for user %d: %v", userID, err)
		return nil, err
	}
	return answers, nil
}
