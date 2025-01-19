package domain

/*
 * @apiDefine: CustomerLifecycle
 */
type CustomerLifecycle struct {
	Stage      string `json:"stage" openapi:"example:active"`  // active, inactive, dead
	Count      int    `json:"count" openapi:"example: 1"`      // number of customers in this stage
	StageOrder int    `json:"stageOrder" openapi:"example: 1"` // order for display (1: active, 2: inactive, 3: dead)
}
