package models

type Question struct {
	ID           int
	QuestionText string
	QuestionType string
	Options      []string
}

type Answer struct {
	QuestionID  int
	AnswerText  string
	AnswerValue int
}

type SubmitAnswersRequest struct {
	Answers []Answer
}
