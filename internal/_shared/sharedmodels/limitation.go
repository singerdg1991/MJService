package sharedmodels

/*
 * @apiDefine: SharedLimitation
 */
type SharedLimitation struct {
	ID          uint   `json:"id" openapi:"example:1"`
	Name        string `json:"name" openapi:"example:limitation"`
	Description string `json:"description" openapi:"example:limitation"`
}
