package sharedmodels

type CustomersCreatePersonalInfo struct {
	CustomerID           int64  `json:"customerId"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Gender               string `json:"gender"`
	NationalCode         string `json:"nationalCode"`
	DateOfBirth          string `json:"dateOfBirth"`
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	ForcedChangePassword bool   `json:"forcedChangePassword"`
}
