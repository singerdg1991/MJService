package models

/*
 * @apiDefine: PushesResponseData
 */
type PushesResponseData struct {
	ID         uint    `json:"id" openapi:"example:1"`
	UserID     uint    `json:"userId" openapi:"example:1"`
	Endpoint   string  `json:"endpoint" openapi:"example:endpoint;required"`
	KeysAuth   string  `json:"keysAuth" openapi:"example:keysAuth;required"`
	KeysP256dh string  `json:"keysP256dh" openapi:"example:keysP256dh;required"`
	CreatedAt  string  `json:"created_at" openapi:"example:2021-01-01T00:00:00Z"`
	UpdatedAt  string  `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z"`
	DeletedAt  *string `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z"`
}

/*
 * @apiDefine: PushesCreateResponse
 */
type PushesCreateResponse struct {
	StatusCode int                `json:"statusCode" openapi:"example:200"`
	Data       PushesResponseData `json:"data" openapi:"$ref:PushesResponseData"`
}
