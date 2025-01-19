package models

/*
 * @apiDefine: RolesCreateResponseDataPermission
 */
type RolesCreateResponseDataPermission struct {
	ID    int64  `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:John;required"`
	Title string `json:"title" openapi:"example:John;required"`
}

/*
 * @apiDefine: RolesCreateResponseData
 */
type RolesCreateResponseData struct {
	ID          int64                               `json:"id" openapi:"example:1"`
	Name        string                              `json:"name" openapi:"example:John;required"`
	Type        string                              `json:"type" openapi:"example:core;required"`
	Permissions []RolesCreateResponseDataPermission `json:"permissions" openapi:"$ref:RolesCreateResponseDataPermission;type:array;required"`
	CreatedAt   string                              `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt   string                              `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt   *string                             `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: RolesCreateResponse
 */
type RolesCreateResponse struct {
	StatusCode int                       `json:"statusCode" openapi:"example:200"`
	Data       []RolesCreateResponseData `json:"data" openapi:"$ref:RolesCreateResponseData;type:array;required"`
}
