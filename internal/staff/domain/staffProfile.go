package domain

import "encoding/json"

/*
 * @apiDefine: StaffProfileSection
 */
type StaffProfileSection struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

/*
 * @apiDefine: StaffProfileTeam
 */
type StaffProfileTeam struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

/*
 * @apiDefine: StaffProfileRole
 */
type StaffProfileRole struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

/*
 * @apiDefine: StaffProfileAbility
 */
type StaffProfileAbility struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffProfilePermission
 */
type StaffProfilePermission struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffProfileLimitation
 */
type StaffProfileLimitation struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffProfileLanguageSkill
 */
type StaffProfileLanguageSkill struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffProfileVehicleType
 */
type StaffProfileVehicleType struct {
	ID   uint   `json:"id" openapi:"example:1"`
	Name string `json:"name" openapi:"example:team"`
}

/*
 * @apiDefine: StaffProfile
 */
type StaffProfile struct {
	StaffID        uint                        `json:"staffId" openapi:"example:1"`
	Sections       []StaffProfileSection       `json:"sections" openapi:"$ref:StaffProfileSection;type:array"`
	Teams          []StaffProfileTeam          `json:"teams" openapi:"$ref:StaffProfileTeam;type:array"`
	Roles          []StaffProfileRole          `json:"roles" openapi:"$ref:StaffProfileRole;type:array"`
	Abilities      []StaffProfileAbility       `json:"abilities" openapi:"$ref:StaffProfileAbility;type:array"`
	Permissions    []StaffProfilePermission    `json:"permissions" openapi:"$ref:StaffProfilePermission;type:array"`
	Limitations    []StaffProfileLimitation    `json:"limitations" openapi:"$ref:StaffProfileLimitation;type:array"`
	LanguageSkills []StaffProfileLanguageSkill `json:"languageSkills" openapi:"$ref:StaffProfileLanguageSkill;type:array"`
	VehicleTypes   []StaffProfileVehicleType   `json:"vehicleTypes" openapi:"$ref:StaffProfileVehicleType;type:array"`
}

func NewStaffProfile() *StaffProfile {
	return &StaffProfile{}
}

func (u *StaffProfile) ToJson() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *StaffProfile) ToMap() (map[string]interface{}, error) {
	jsonString, err := u.ToJson()
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
