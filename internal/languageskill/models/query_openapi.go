package models

import "github.com/hoitek/Maja-Service/internal/languageskill/domain"

/*
 * @apiDefine: LanguageSkillsQueryResponseDataItem
 */
type LanguageSkillsQueryResponseDataItem struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:John;required"`
	Description *string `json:"description" openapi:"example:John;required"`
}

/*
 * @apiDefine: LanguageSkillsQueryResponseData
 */
type LanguageSkillsQueryResponseData struct {
	Limit      int                                   `json:"limit" openapi:"example:10"`
	Offset     int                                   `json:"offset" openapi:"example:0"`
	Page       int                                   `json:"page" openapi:"example:1"`
	TotalRows  int                                   `json:"totalRows" openapi:"example:1"`
	TotalPages int                                   `json:"totalPages" openapi:"example:1"`
	Items      []LanguageSkillsQueryResponseDataItem `json:"items" openapi:"$ref:LanguageSkillsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: LanguageSkillsQueryResponse
 */
type LanguageSkillsQueryResponse struct {
	StatusCode int                             `json:"statusCode" openapi:"example:200;"`
	Data       LanguageSkillsQueryResponseData `json:"data" openapi:"$ref:LanguageSkillsQueryResponseData;type:object;"`
}

/*
 * @apiDefine: LanguageSkillsQueryNotFoundResponse
 */
type LanguageSkillsQueryNotFoundResponse struct {
	LanguageSkills []domain.LanguageSkill `json:"languageskills" openapi:"$ref:LanguageSkill;type:array"`
}
