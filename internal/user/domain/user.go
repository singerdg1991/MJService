package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/security"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
 * @apiDefine: UserLanguageSkillRes
 */
type UserLanguageSkillRes struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:English"`
}

/*
 * @apiDefine: UserRolePermission
 */
type UserRolePermission struct {
	ID    uint   `json:"id" openapi:"example:1"`
	Name  string `json:"name" openapi:"example:Dashboard"`
	Title string `json:"title" openapi:"example:Can See Dashboard"`
}

/*
 * @apiDefine: UserRole
 */
type UserRole struct {
	ID          uint                 `json:"id" openapi:"example:1"`
	Name        string               `json:"name" openapi:"example:John;required"`
	Permissions []UserRolePermission `json:"permissions" openapi:"$ref:UserRolePermission"`
}

/*
 * @apiDefine: UserRoleID
 */
type UserRoleID struct {
	ID uint `json:"id" openapi:"example:1"`
}

/*
 * @apiDefine: User
 */
type User struct {
	ID                      uint                            `json:"id" openapi:"example:1"`
	MongoID                 primitive.ObjectID              `bson:"_id,omitempty" json:"-" openapi:"example:5f7b5f5b9b9b9b9b9b9b9b9b"`
	RoleIDs                 []UserRoleID                    `json:"roleIds" openapi:"$ref:UserRoleID;type:array;" bson:"roleIds"`
	StaffID                 *uint                           `json:"staffId" openapi:"-"`
	CustomerID              *uint                           `json:"customerId" openapi:"-"`
	FirstName               string                          `json:"firstName" openapi:"example:John;required"`
	LastName                string                          `json:"lastName" openapi:"example:Doe;ignored"`
	Username                string                          `json:"username" openapi:"example:johndoe"`
	Password                string                          `json:"password" openapi:"ignore:true"`
	Email                   string                          `json:"email" openapi:"example:saeed@gmail.com"`
	Phone                   string                          `json:"phone" openapi:"example:09123456789"`
	Telephone               *string                         `json:"telephone" openapi:"example:02112345678"`
	LanguageSkills          []UserLanguageSkillRes          `json:"languageSkills" openapi:"ignored"`
	RegistrationNumber      *string                         `json:"registrationNumber" openapi:"example:1234567890"`
	WorkPhoneNumber         *string                         `json:"workPhoneNumber" openapi:"example:02112345678"`
	Gender                  *string                         `json:"gender" openapi:"example:male"`
	AccountNumber           *string                         `json:"accountNumber" openapi:"example:02112345678"`
	NationalCode            string                          `json:"nationalCode" openapi:"example:1234567890"`
	BirthDate               time.Time                       `json:"birthDate" openapi:"example:1990-01-01T00:00:00Z"`
	AvatarUrl               string                          `json:"avatarUrl" openapi:"example:https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
	ForcedChangePassword    bool                            `json:"forcedChangePassword" openapi:"example:true"`
	Roles                   []UserRole                      `json:"roles" openapi:"$ref:UserRole;type:array;" bson:"roles"`
	PrivacyPolicyAcceptedAt *time.Time                      `json:"privacy_policy_accepted_at" openapi:"example:2021-01-01T00:00:00Z"`
	UserType                string                          `json:"userType" openapi:"example:staff"`
	VehicleTypes            interface{}                     `json:"vehicleTypes" openapi:"example:[\"car\",\"bicycle\",\"public_transportation\"];type:array;required;"`
	VehicleLicenseTypes     interface{}                     `json:"vehicleLicenseTypes" openapi:"example:[\"automatic\",\"manual\"];type:array;required;"`
	Limitations             []sharedmodels.SharedLimitation `json:"limitations" openapi:"$ref:SharedLimitation;type:array;required;"`
	SuspendedAt             *time.Time                      `json:"suspended_at" openapi:"example:2021-01-01T00:00:00Z"`
	CreatedAt               time.Time                       `json:"created_at" openapi:"example:2021-01-01T00:00:00Z" bson:"created_at"`
	UpdatedAt               time.Time                       `json:"updated_at" openapi:"example:2021-01-01T00:00:00Z" bson:"updated_at"`
	DeletedAt               *time.Time                      `json:"deleted_at" openapi:"example:2021-01-01T00:00:00Z" bson:"deleted_at"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *User) ValidatePassword(password string) error {
	isValid := security.ValidatePassword(password, u.Password)
	if isValid {
		return nil
	}
	return errors.New("user not found")
}

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.Roles {
		for _, rolePermission := range role.Permissions {
			if rolePermission.Name == permission {
				return true
			}
		}
	}
	return false
}
