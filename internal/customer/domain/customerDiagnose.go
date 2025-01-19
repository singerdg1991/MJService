package domain

/*
 * @apiDefine: CustomerDiagnoseDiagnose
 */
type CustomerDiagnoseDiagnose struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Title string `json:"title" openapi:"example:title;required"`
}

/*
 * @apiDefine: CustomerDiagnose
 */
type CustomerDiagnose struct {
	ID         uint                      `json:"id" openapi:"example:1"`
	CustomerID uint                      `json:"customerId" openapi:"example:1"`
	DiagnoseID uint                      `json:"diagnoseId" openapi:"example:1"`
	Diagnose   *CustomerDiagnoseDiagnose `json:"diagnose" openapi:"$ref:CustomerDiagnoseDiagnose"`
}
