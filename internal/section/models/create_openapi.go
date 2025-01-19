package models

/*
 * @apiDefine: SectionsCreateResponseDataParent
 */
type SectionsCreateResponseDataParent struct {
	ID          int     `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:test"`
	Color       *string `json:"color" openapi:"example:#000000;nullable"`
	Description *string `json:"description" openapi:"example:description;nullable"`
}

/*
 * @apiDefine: SectionsCreateResponseDataItem
 */
type SectionsCreateResponseDataItem struct {
	ID          int                              `json:"id" openapi:"example:1"`
	Name        string                           `json:"name" openapi:"example:test"`
	Parent      SectionsCreateResponseDataParent `json:"parent" openapi:"$ref:SectionsCreateResponseDataParent;type:object;nullable"`
	Color       *string                          `json:"color" openapi:"example:#000000;nullable"`
	Description *string                          `json:"description" openapi:"example:description;nullable"`
}

/*
 * @apiDefine: SectionsCreateResponseData
 */
type SectionsCreateResponseData struct {
	ID          int                              `json:"id" openapi:"example:1"`
	Name        string                           `json:"name" openapi:"example:test"`
	Parent      SectionsCreateResponseDataParent `json:"parent" openapi:"$ref:SectionsCreateResponseDataParent;type:object;nullable"`
	Children    []SectionsCreateResponseDataItem `json:"children" openapi:"$ref:SectionsCreateResponseDataItem;type:array;nullable"`
	Color       *string                          `json:"color" openapi:"example:#000000;nullable"`
	Description *string                          `json:"description" openapi:"example:description;nullable"`
}

/*
 * @apiDefine: SectionsCreateResponse
 */
type SectionsCreateResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200"`
	Data       SectionsCreateResponseData `json:"data" openapi:"$ref:SectionsCreateResponseData"`
}
