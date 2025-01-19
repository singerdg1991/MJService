package models

import "github.com/hoitek/Maja-Service/internal/medicine/domain"

/*
 * @apiDefine: MedicinesQueryResponseData
 */
type MedicinesQueryResponseData struct {
	Limit      int                     `json:"limit" openapi:"example:10"`
	Offset     int                     `json:"offset" openapi:"example:0"`
	Page       int                     `json:"page" openapi:"example:1"`
	TotalRows  int                     `json:"totalRows" openapi:"example:1"`
	TotalPages int                     `json:"totalPages" openapi:"example:1"`
	Items      []MedicinesResponseData `json:"items" openapi:"$ref:MedicinesResponseData;type:array"`
}

/*
 * @apiDefine: MedicinesQueryResponse
 */
type MedicinesQueryResponse struct {
	StatusCode int                        `json:"statusCode" openapi:"example:200;"`
	Data       MedicinesQueryResponseData `json:"data" openapi:"$ref:MedicinesQueryResponseData;type:object;"`
}

/*
 * @apiDefine: MedicinesQueryNotFoundResponse
 */
type MedicinesQueryNotFoundResponse struct {
	Medicines []domain.Medicine `json:"medicines" openapi:"$ref:Medicine;type:array"`
}
