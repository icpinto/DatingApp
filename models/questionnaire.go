package models

type Question struct {
	ID           int
	QuestionText string
	QuestionType string
	Options      []string
}

type Answer struct {
	QuestionID  int    `json:"question_id"`
	AnswerText  string `json:"answer_text"`
	AnswerValue int    `json:"answer_value"`
}

type SubmitAnswersRequest struct {
	Answers []Answer
}
