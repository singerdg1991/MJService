package sharedmodels

type AuthenticatedUser struct {
	ID        uint   `json:"id" openapi:"example:1"`
	FirstName string `json:"firstName" openapi:"example:firstName;required"`
	LastName  string `json:"lastName" openapi:"example:lastName;required"`
	Email     string `json:"email" openapi:"example:email;required"`
	AvatarUrl string `json:"avatarUrl" openapi:"example:avatarUrl;required"`
}
