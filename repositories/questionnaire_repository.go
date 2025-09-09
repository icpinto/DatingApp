package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

func GetQuestions(db *sql.DB) ([]models.Question, error) {
	rows, err := db.Query(`SELECT id, question_text, question_type, options FROM questions`)
	if err != nil {
		log.Printf("GetQuestions query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType, pq.Array(&q.Options)); err != nil {
			log.Printf("GetQuestions scan error: %v", err)
			return nil, err
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetQuestions rows error: %v", err)
		return nil, err
	}
	return questions, nil
}

func UpsertAnswer(db *sql.DB, userID int, answer models.Answer) error {
	_, err := db.Exec(`
                INSERT INTO user_answers (user_id, question_id, answer_text, answer_value, created_at)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (user_id, question_id)
                DO UPDATE SET
                        answer_text = EXCLUDED.answer_text,
                        answer_value = EXCLUDED.answer_value,
                        created_at = EXCLUDED.created_at`,
		userID, answer.QuestionID, answer.AnswerText, answer.AnswerValue, time.Now())
	if err != nil {
		log.Printf("UpsertAnswer exec error for user %d question %d: %v", userID, answer.QuestionID, err)
	}
	return err
}

func GetAnswersByUserID(db *sql.DB, userID int) ([]models.Answer, error) {
	rows, err := db.Query(`
        SELECT question_id, answer_text, answer_value
        FROM user_answers WHERE user_id = $1`, userID)
	if err != nil {
		log.Printf("GetAnswersByUserID query error for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var a models.Answer
		if err := rows.Scan(&a.QuestionID, &a.AnswerText, &a.AnswerValue); err != nil {
			log.Printf("GetAnswersByUserID scan error: %v", err)
			return nil, err
		}
		answers = append(answers, a)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetAnswersByUserID rows error: %v", err)
		return nil, err
	}
	return answers, nil
}
