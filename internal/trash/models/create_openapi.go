package models

import "github.com/hoitek/Maja-Service/internal/trash/domain"

/*
 * @apiDefine: TrashesCreateResponse
 */
type TrashesCreateResponse struct {
	Trash domain.Trash `json:"trash" openapi:"$ref:Trash;type:object;"`
}
