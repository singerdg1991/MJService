package models

import (
	"github.com/hoitek/Maja-Service/internal/user/domain"
)

/*
 * @apiDefine: OAuth2UserInfoResponse
 */
type OAuth2UserInfoResponse struct {
	User domain.User `json:"user" openapi:"$ref:User;type:object"`
}

/*
 * @apiDefine: OAuth2UserInfoNotFoundResponse
 */
type OAuth2UserInfoNotFoundResponse struct {
	User domain.User `json:"user" openapi:"$ref:User;type:object"`
}
