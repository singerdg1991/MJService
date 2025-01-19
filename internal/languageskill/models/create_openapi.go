package models

/*
 * @apiDefine: LanguageSkillsResponseData
 */
type LanguageSkillsResponseData struct {
	ID          uint    `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:saeed"`
	Description *string `json:"description" openapi:"example:test"`
}

/*
 * @apiDefine: LanguageSkillsCreateResponse
 */
type LanguageSkillsCreateResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200"`
	Data       LanguageSkillsResponseData `json:"data" openapi:"$ref:LanguageSkillsResponseData"`
}
