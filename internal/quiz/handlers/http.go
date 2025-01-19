package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	permPorts "github.com/hoitek/Maja-Service/internal/permission/ports"
	"github.com/hoitek/Maja-Service/internal/quiz/config"
	"github.com/hoitek/Maja-Service/internal/quiz/models"
	rPorts "github.com/hoitek/Maja-Service/internal/quiz/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type QuizHandler struct {
	QuizService       rPorts.QuizService
	PermissionService permPorts.PermissionService
	UserService       uPorts.UserService
}

func NewQuizHandler(r *mux.Router, s rPorts.QuizService, p permPorts.PermissionService, u uPorts.UserService) (QuizHandler, error) {
	quizHandler := QuizHandler{
		QuizService:       s,
		PermissionService: p,
		UserService:       u,
	}
	if r == nil {
		return QuizHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.QuizConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.QuizConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(u, []string{}))

	rAuth.Handle("/quizzes", quizHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/quizzes", quizHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/quizzes", quizHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/quizzes/answer", quizHandler.CreateAnswer()).Methods(http.MethodPost)
	rAuth.Handle("/quizzes/answers", quizHandler.QueryAnswers()).Methods(http.MethodGet)
	rAuth.Handle("/quizzes/end/{quizId}", quizHandler.UpdateEnd()).Methods(http.MethodPut)
	rAuth.Handle("/quizzes/question", quizHandler.CreateQuestion()).Methods(http.MethodPost)
	rAuth.Handle("/quizzes/question/{questionId}", quizHandler.UpdateQuestion()).Methods(http.MethodPut)
	rAuth.Handle("/quizzes/questions", quizHandler.QueryQuestions()).Methods(http.MethodGet)
	rAuth.Handle("/quizzes/question", quizHandler.QueryQuestions()).Methods(http.MethodPut)
	rAuth.Handle("/quizzes/start", quizHandler.CreateStart()).Methods(http.MethodPost)
	rAuth.Handle("/quizzes/participants", quizHandler.QueryParticipants()).Methods(http.MethodGet)
	rAuth.Handle("/quizzes/status/{id}", quizHandler.UpdateStatus()).Methods(http.MethodPut)
	rAuth.Handle("/quizzes/{id}", quizHandler.Update()).Methods(http.MethodPut)

	return quizHandler, nil
}

/*
* @apiTag: quiz
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: quizzes
* @apiResponseRef: QuizzesQueryResponse
* @apiSummary: Query quizzes
* @apiParametersRef: QuizzesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: QuizzesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: QuizzesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.QuizzesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}

		// Get quizzes
		quizzes, err := h.QuizService.Query(queries, &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(quizzes)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: QuizzesCreateRequestBody
 * @apiResponseRef: QuizzesCreateResponse
 * @apiSummary: Create quiz
 * @apiDescription: Create quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.QuizzesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByTitle(payload.Title, &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if foundQuiz != nil {
			return response.ErrorNotFound(nil, "Quiz already exists")
		}

		// Convert interface slice to slice of int64
		if payload.ParticipantUserIDs != nil {
			ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.ParticipantUserIDs)
			if err != nil {
				return response.ErrorBadRequest(nil, "ParticipantUserIDs are invalid")
			}
			payload.ParticipantUserIDsAsInt64 = ids
		}

		// Check Permissions Exists
		//permissions, err := h.PermissionService.GetPermissionsByIds(payload.PermissionsInt64)
		//if err != nil {
		//	log.Printf("Error when get permissions by ids: %v\n", err)
		//	return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		//}
		//if len(permissions) != len(payload.PermissionsInt64) {
		//	return response.ErrorBadRequest(nil, "Permissions is invalid")
		//}

		// Create quiz
		insertedQuiz, err := h.QuizService.Create(payload)
		if err != nil {
			log.Printf("Error when create quiz: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(insertedQuiz)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: QuizzesUpdateRequestParams
 * @apiRequestRef: QuizzesCreateRequestBody
 * @apiResponseRef: QuizzesCreateResponse
 * @apiSummary: Update quiz
 * @apiDescription: Update quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.QuizzesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}

		// Validate body
		payload := &models.QuizzesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Convert interface slice to slice of int64
		if payload.ParticipantUserIDs != nil {
			ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.ParticipantUserIDs)
			if err != nil {
				return response.ErrorBadRequest(nil, "ParticipantUserIDs are invalid")
			}
			payload.ParticipantUserIDsAsInt64 = ids
		}

		// Check Permissions Exists
		//permissions, err := h.PermissionService.GetPermissionsByIds(payload.PermissionsInt64)
		//if err != nil {
		//	log.Printf("Error when get permissions by ids: %v\n", err)
		//	return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		//}
		//if len(permissions) != len(payload.PermissionsInt64) {
		//	return response.ErrorBadRequest(nil, "Permissions is invalid")
		//}

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByID(params.ID, &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if foundQuiz == nil {
			return response.ErrorNotFound(nil, "Quiz not found")
		}

		// Validate quiz title
		foundQuizByTitle := h.QuizService.GetQuizByTitle(payload.Title, &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if foundQuizByTitle != nil && foundQuizByTitle.ID != foundQuiz.ID {
			return response.ErrorBadRequest(nil, "Quiz with this title is already exists")
		}

		// Update quiz
		data, err := h.QuizService.Update(payload, int64(params.ID))
		if err != nil {
			log.Printf("Error when update quiz: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/status/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: QuizzesUpdateStatusRequestParams
 * @apiRequestRef: QuizzesUpdateStatusRequestBody
 * @apiResponseRef: QuizzesCreateResponse
 * @apiSummary: Update quiz
 * @apiDescription: Update quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) UpdateStatus() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.QuizzesUpdateStatusRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate body
		payload := &models.QuizzesUpdateStatusRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByID(params.ID, &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if foundQuiz == nil {
			return response.ErrorNotFound(nil, "Quiz not found")
		}

		// Update quiz status
		updatedQuiz, err := h.QuizService.UpdateStatus(payload, int64(params.ID))
		if err != nil {
			log.Printf("Error when update quiz: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(updatedQuiz)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: QuizzesDeleteRequestBody
 * @apiResponseRef: QuizzesCreateResponse
 * @apiSummary: Delete quiz
 * @apiDescription: Delete quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.QuizzesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Convert interface slice to slice of int64
		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "Quiz IDs are invalid")
		}
		payload.IDsInt64 = ids

		// Delete quiz
		data, err := h.QuizService.Delete(payload)
		if err != nil {
			log.Printf("Error when delete quiz: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/question
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: QuizzesCreateQuestionRequestBody
 * @apiResponseRef: QuizzesCreateQuestionResponse
 * @apiSummary: Create quiz question
 * @apiDescription: Create quiz question
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) CreateQuestion() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.QuizzesCreateQuestionRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByID(int(payload.QuizID), &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if foundQuiz == nil {
			return response.ErrorNotFound(nil, "Quiz not found")
		}

		// Create quiz
		insertedQuizQuestion, err := h.QuizService.CreateQuestion(payload)
		if err != nil {
			log.Printf("Error when create quiz question: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(insertedQuizQuestion)
	})
}

/*
* @apiTag: quiz
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /quizzes/questions
* @apiResponseRef: QuizzesQueryQuestionsResponse
* @apiSummary: Query quizzes questions
* @apiParametersRef: QuizzesQueryQuestionsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: QuizzesQueryQuestionsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: QuizzesQueryQuestionsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) QueryQuestions() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.QuizzesQueryQuestionsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		quizQuestions, err := h.QuizService.QueryQuestions(queries)
		if err != nil {
			log.Printf("Error when query quiz questions: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(quizQuestions)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/question/{questionId}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: QuizzesUpdateQuestionRequestParams
 * @apiRequestRef: QuizzesCreateQuestionRequestBody
 * @apiResponseRef: QuizzesCreateResponse
 * @apiSummary: Update quiz
 * @apiDescription: Update quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) UpdateQuestion() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.QuizzesUpdateQuestionRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate body
		payload := &models.QuizzesCreateQuestionRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check quiz exists
		foundQuizQuestion := h.QuizService.GetQuizQuestionByID(int64(params.QuestionID))
		if foundQuizQuestion == nil {
			return response.ErrorBadRequest(nil, "Quiz question not found")
		}
		if int64(foundQuizQuestion.QuizID) != payload.QuizID {
			return response.ErrorBadRequest(nil, "Quiz question not found")
		}

		// Update quiz question
		updatedQuizQuestion, err := h.QuizService.UpdateQuizQuestion(payload, int64(params.QuestionID))
		if err != nil {
			log.Printf("Error when update quiz question: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(updatedQuizQuestion)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/start
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: QuizzesCreateStartRequestBody
 * @apiResponseRef: QuizzesCreateStartResponse
 * @apiSummary: Start quiz
 * @apiDescription: Start quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) CreateStart() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.QuizzesCreateStartRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.UserID = int64(authenticatedUser.ID)

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByID(int(payload.QuizID), &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		})
		if foundQuiz == nil {
			return response.ErrorNotFound(nil, "Quiz not found")
		}

		// Check quiz password if is lock
		if *foundQuiz.IsLock {
			if !foundQuiz.ValidatePassword(payload.Password) {
				return response.ErrorBadRequest(nil, "Quiz password is invalid")
			}
		}

		// Find quiz participant
		quizParticipant := h.QuizService.GetQuizParticipantByQuizIDAndUserID(payload.QuizID, payload.UserID)
		if quizParticipant != nil {
			return response.ErrorBadRequest(nil, "You already participated in this quiz")
		}

		// Check user is in quiz available participants
		canStartQuiz := false
		for _, availableParticipant := range foundQuiz.AvailableParticipants {
			if int64(availableParticipant.UserID) == payload.UserID {
				canStartQuiz = true
				break
			}
		}
		if !canStartQuiz {
			return response.ErrorBadRequest(nil, "You are not allowed to participate in this quiz")
		}

		// Check quiz questions
		quizQuestions, err := h.QuizService.QueryQuestions(&models.QuizzesQueryQuestionsRequestParams{
			QuizID: payload.QuizID,
		})
		if err != nil {
			log.Printf("Error when query quiz questions: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}
		if quizQuestions.TotalRows == 0 {
			return response.ErrorBadRequest(nil, "For this quiz there are no questions")
		}

		// Create quiz
		startedQuizParticipant, err := h.QuizService.StartQuiz(payload)
		if err != nil {
			log.Printf("Error when create quiz question: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(startedQuizParticipant)
	})
}

/*
* @apiTag: quiz
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /quizzes/participants
* @apiResponseRef: QuizzesQueryParticipantsResponse
* @apiSummary: Query quizzes participants
* @apiParametersRef: QuizzesQueryParticipantsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: QuizzesQueryParticipantsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: QuizzesQueryParticipantsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) QueryParticipants() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.QuizzesQueryParticipantsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		quizParticipants, err := h.QuizService.QueryParticipants(queries)
		if err != nil {
			log.Printf("Error when query quiz participants: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(quizParticipants)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/answer
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: QuizzesCreateAnswerRequestBody
 * @apiResponseRef: QuizzesCreateAnswerResponse
 * @apiSummary: Answer quiz question
 * @apiDescription: Answer quiz question
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) CreateAnswer() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.QuizzesCreateAnswerRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.User = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByID(int(payload.QuizID), payload.User)
		if foundQuiz == nil {
			return response.ErrorNotFound(nil, "Quiz not found")
		}

		// Check quiz is started
		foundQuizParticipant := h.QuizService.GetQuizParticipantByQuizIDAndUserID(payload.QuizID, int64(payload.User.ID))
		if foundQuizParticipant == nil {
			return response.ErrorBadRequest(nil, "You are not allowed to answer this quiz, please start this quiz first")
		}

		// Check quiz is finished
		if foundQuizParticipant.EndedAt != nil {
			return response.ErrorBadRequest(nil, "This quiz is already finished, you are not allowed to answer this quiz")
		}

		// Check question is exists
		foundQuizQuestion := h.QuizService.GetQuizQuestionByID(payload.QuestionID)
		if foundQuizQuestion == nil {
			return response.ErrorNotFound(nil, "Quiz question not found")
		}

		// Check question is belongs to this quiz
		if int64(foundQuizQuestion.QuizID) != payload.QuizID {
			return response.ErrorBadRequest(nil, "This question is not belongs to this quiz")
		}

		// Check quiz question option is exists
		foundQuizQuestionOption := h.QuizService.GetQuizQuestionOptionByID(payload.QuizQuestionOptionID)
		if foundQuizQuestionOption == nil {
			return response.ErrorNotFound(nil, "Quiz question option not found")
		}

		// Check quiz question option is belongs to this quiz question
		if int64(foundQuizQuestionOption.QuizQuestionId) != payload.QuestionID {
			return response.ErrorBadRequest(nil, "This question option is not belongs to this question, please try again")
		}

		// Create quiz
		insertedQuizQuestion, err := h.QuizService.CreateAnswer(payload)
		if err != nil {
			log.Printf("Error when create quiz question: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(insertedQuizQuestion)
	})
}

/*
* @apiTag: quiz
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /quizzes/answers
* @apiResponseRef: QuizzesQueryQuestionAnswersResponse
* @apiSummary: Query quizzes answers
* @apiParametersRef: QuizzesQueryQuestionAnswersRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: QuizzesQueryQuestionAnswersNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: QuizzesQueryQuestionAnswersNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) QueryAnswers() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.QuizzesQueryQuestionAnswersRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}

		// Get quizzes
		quizzes, err := h.QuizService.QueryAnswers(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(quizzes)
	})
}

/*
 * @apiTag: quiz
 * @apiPath: /quizzes/end/{quizId}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: QuizzesUpdateEndRequestParams
 * @apiResponseRef: QuizzesUpdateEndResponse
 * @apiSummary: End quiz
 * @apiDescription: End quiz
 * @apiSecurity: apiKeySecurity
 */
func (h *QuizHandler) UpdateEnd() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.QuizzesUpdateEndRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate body
		payload := &models.QuizzesUpdateEndRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get authenticated user
		authenticatedUser := h.UserService.GetUserFromContext(r.Context())
		if authenticatedUser == nil {
			return response.ErrorBadRequest(nil, "You are not authorized to perform this action")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        authenticatedUser.ID,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
			Email:     authenticatedUser.Email,
			AvatarUrl: authenticatedUser.AvatarUrl,
		}

		// Check quiz exists
		foundQuiz := h.QuizService.GetQuizByID(params.QuizID, payload.AuthenticatedUser)
		if foundQuiz == nil {
			return response.ErrorBadRequest(nil, "Quiz not found")
		}

		// Check quiz is started
		foundQuizParticipant := h.QuizService.GetQuizParticipantByQuizIDAndUserID(int64(foundQuiz.ID), int64(payload.AuthenticatedUser.ID))
		if foundQuizParticipant == nil {
			return response.ErrorBadRequest(nil, "You are not allowed to end this quiz, please start this quiz first")
		}

		// Check quiz is ended
		if foundQuizParticipant.EndedAt != nil {
			return response.ErrorBadRequest(nil, "This quiz is already ended, you are not allowed to end this quiz")
		}

		// Check all questions is answered
		quizAnswers := h.QuizService.GetQuizAnswersByQuizIDAndUserID(int64(foundQuiz.ID), int64(payload.AuthenticatedUser.ID))
		quizQuestions := h.QuizService.GetQuizQuestionsByQuizID(int64(foundQuiz.ID))
		log.Println("quizQuestions", quizQuestions)
		if len(quizAnswers) != len(quizQuestions) {
			return response.ErrorBadRequest(nil, "You are not allowed to end this quiz, please answer all questions first")
		}

		// Update quiz participant
		updatedQuizParticipant, err := h.QuizService.UpdateQuizEnd(payload, int64(foundQuiz.ID))
		if err != nil {
			log.Printf("Error when update quiz end: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(updatedQuizParticipant)
	})
}
