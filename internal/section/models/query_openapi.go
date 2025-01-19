package models

/*
 * @apiDefine: SectionsQueryResponseDataItemParentParent
 */
type SectionsQueryResponseDataItemParentParent struct {
}

/*
 * @apiDefine: SectionsQueryResponseDataItemParent
 */
type SectionsQueryResponseDataItemParent struct {
	ID          int     `json:"id" openapi:"example:1"`
	Name        string  `json:"name" openapi:"example:test"`
	Color       *string `json:"color" openapi:"example:#000000;nullable"`
	Description *string `json:"description" openapi:"example:description;nullable"`
}

/*
 * @apiDefine: SectionsQueryResponseDataItemChildren
 */
type SectionsQueryResponseDataItemChildren struct {
	ID          int                                 `json:"id" openapi:"example:1"`
	Name        string                              `json:"name" openapi:"example:test"`
	Parent      SectionsQueryResponseDataItemParent `json:"parent" openapi:"$ref:SectionsQueryResponseDataItemParent;type:object;nullable"`
	Color       *string                             `json:"color" openapi:"example:#000000;nullable"`
	Description *string                             `json:"description" openapi:"example:description;nullable"`
}

/*
 * @apiDefine: SectionsQueryResponseDataItem
 */
type SectionsQueryResponseDataItem struct {
	ID          int                                     `json:"id" openapi:"example:1"`
	Name        string                                  `json:"name" openapi:"example:test"`
	Parent      SectionsQueryResponseDataItemParent     `json:"parent" openapi:"$ref:SectionsQueryResponseDataItemParent;type:object;nullable"`
	Children    []SectionsQueryResponseDataItemChildren `json:"children" openapi:"$ref:SectionsQueryResponseDataItemChildren;type:array;nullable"`
	Color       *string                                 `json:"color" openapi:"example:#000000;nullable"`
	Description *string                                 `json:"description" openapi:"example:description;nullable"`
}

/*
 * @apiDefine: SectionsQueryResponseData
 */
type SectionsQueryResponseData struct {
	Limit      int                             `json:"limit" openapi:"example:10"`
	Offset     int                             `json:"offset" openapi:"example:0"`
	Page       int                             `json:"page" openapi:"example:1"`
	TotalRows  int                             `json:"totalRows" openapi:"example:1"`
	TotalPages int                             `json:"totalPages" openapi:"example:1"`
	Items      []SectionsQueryResponseDataItem `json:"items" openapi:"$ref:SectionsQueryResponseDataItem;type:array"`
}

/*
 * @apiDefine: SectionsQueryResponse
 */
type SectionsQueryResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200"`
	Data       SectionsQueryResponseData `json:"data" openapi:"$ref:SectionsQueryResponseData"`
}

/*
 * @apiDefine: SectionsQueryNotFoundResponse
 */
type SectionsQueryNotFoundResponse struct {
	StatusCode int    `json:"statusCode" openapi:"example:404"`
	Message    string `json:"message" openapi:"example:Not Found"`
}
