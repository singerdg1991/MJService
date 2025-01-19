package models

import (
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"net/http"
)

/*
 * @apiDefine: QuizzesCreateAnswerRequestBody
 */
type QuizzesCreateAnswerRequestBody struct {
	QuizID               int64                           `json:"quizId" openapi:"example:1;required;"`
	QuestionID           int64                           `json:"questionId" openapi:"example:1;required;"`
	QuizQuestionOptionID int64                           `json:"quizQuestionOptionId" openapi:"example:1;required;"`
	User                 *sharedmodels.AuthenticatedUser `json:"-" openapi:"ignored"`
}

func (data *QuizzesCreateAnswerRequestBody) ValidateBody(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"quizId":               govalidity.New("quizId").Int().Min(1).Required(),
		"questionId":           govalidity.New("questionId").Int().Min(1).Required(),
		"quizQuestionOptionId": govalidity.New("quizQuestionOptionId").Int().Min(1).Required(),
	}

	errs := govalidity.ValidateBody(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}
