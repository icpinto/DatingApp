package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

// QuestionRepository manages questions and user answers.
type QuestionRepository struct {
	db *sql.DB
}

// NewQuestionRepository creates a new QuestionRepository.
func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// GetAll retrieves all questions.
func (r *QuestionRepository) GetAll() ([]models.Question, error) {
	rows, err := r.db.Query(`SELECT id, question_text, question_type, options FROM questions`)
	if err != nil {
		log.Printf("QuestionRepository.GetAll query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType, pq.Array(&q.Options)); err != nil {
			log.Printf("QuestionRepository.GetAll scan error: %v", err)
			return nil, err
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		log.Printf("QuestionRepository.GetAll rows error: %v", err)
		return nil, err
	}
	return questions, nil
}

// UpsertAnswer creates or updates a user's answer.
func (r *QuestionRepository) UpsertAnswer(userID int, answer models.Answer) error {
	_, err := r.db.Exec(`
                INSERT INTO user_answers (user_id, question_id, answer_text, answer_value, created_at)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (user_id, question_id)
                DO UPDATE SET
                        answer_text = EXCLUDED.answer_text,
                        answer_value = EXCLUDED.answer_value,
                        created_at = EXCLUDED.created_at`,
		userID, answer.QuestionID, answer.AnswerText, answer.AnswerValue, time.Now())
	if err != nil {
		log.Printf("QuestionRepository.UpsertAnswer exec error for user %d question %d: %v", userID, answer.QuestionID, err)
	}
	return err
}

// GetAnswersByUserID retrieves all answers for a user.
func (r *QuestionRepository) GetAnswersByUserID(userID int) ([]models.Answer, error) {
	rows, err := r.db.Query(`
        SELECT question_id, answer_text, answer_value
        FROM user_answers WHERE user_id = $1`, userID)
	if err != nil {
		log.Printf("QuestionRepository.GetAnswersByUserID query error for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var a models.Answer
		if err := rows.Scan(&a.QuestionID, &a.AnswerText, &a.AnswerValue); err != nil {
			log.Printf("QuestionRepository.GetAnswersByUserID scan error: %v", err)
			return nil, err
		}
		answers = append(answers, a)
	}
	if err := rows.Err(); err != nil {
		log.Printf("QuestionRepository.GetAnswersByUserID rows error: %v", err)
		return nil, err
	}
	return answers, nil
}

// DeleteAnswer removes an answer for a user and question.
func (r *QuestionRepository) DeleteAnswer(userID, questionID int) error {
	if _, err := r.db.Exec(`DELETE FROM user_answers WHERE user_id = $1 AND question_id = $2`, userID, questionID); err != nil {
		log.Printf("QuestionRepository.DeleteAnswer exec error for user %d question %d: %v", userID, questionID, err)
		return err
	}
	return nil
}
