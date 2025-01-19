package models

/*
 * @apiDefine: QuizzesUpdateEndResponse
 */
type QuizzesUpdateEndResponse struct {
	StatusCode int                                        `json:"statusCode" openapi:"example:200"`
	Data       []QuizzesQueryParticipantsResponseDataItem `json:"data" openapi:"$ref:QuizzesQueryParticipantsResponseDataItem;type:array;required"`
}
