package openapi

/*
 * @apiEnum: GenderEnum
 */
type GenderEnum struct {
	Select string `json:"select"`
	Male   string `json:"male"`
	Female string `json:"female"`
}

/*
 * @apiEnum: FilterOperatorsEnum
 */
type FilterOperatorsEnum struct {
	Equals            string `json:"equals"`
	Contains          string `json:"contains"`
	StartsWith        string `json:"startsWith"`
	EndsWith          string `json:"endsWith"`
	IsEmpty           string `json:"isEmpty"`
	IsNotEmpty        string `json:"isNotEmpty"`
	IsAnyOf           string `json:"isAnyOf"`
	NumberEquals      string `json:"="`
	NumberNotEquals   string `json:"!="`
	GreaterThan       string `json:">"`
	GreaterThanEquals string `json:">="`
	LessThan          string `json:"<"`
	LessThanEquals    string `json:"<="`
	DateIs            string `json:"is"`
	DateIsNot         string `json:"isNot"`
}

/*
 * @apiEnum: UserTypeEnum
 */
type UserTypeEnum struct {
	Customer string `json:"customer"`
	Staff    string `json:"staff"`
}

/*
 * @apiEnum: ExperienceAmountUnitEnum
 */
type ExperienceAmountUnitEnum struct {
	Hour  string `json:"hour"`
	Day   string `json:"day"`
	Week  string `json:"week"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

/*
 * @apiEnum: PaymentMethodEnum
 */
type PaymentMethodEnum struct {
	Seteli       string `json:"seteli"`
	Own          string `json:"own"`
	SeteliAndOwn string `json:"seteliAndOwn"`
}

/*
 * @apiEnum: CustomerStatusEnum
 */
type CustomerStatusEnum struct {
	Active                     string `json:"active"`
	Death                      string `json:"death"`
	FormerCustomerOrDischarged string `json:"formerCustomerOrDischarged"`
}

/*
 * @apiEnum: CustomerServiceTypeEnum
 */
type CustomerServiceTypeEnum struct {
	Online string `json:"online"`
	OnSite string `json:"onsite"`
}

/*
 * @apiEnum: CustomerServiceRepeatEnum
 */
type CustomerServiceRepeatEnum struct {
	Daily          string `json:"daily"`
	Weekly         string `json:"weekly"`
	Monthly        string `json:"monthly"`
	EveryMonday    string `json:"everyMonday"`
	EveryTuesday   string `json:"everyTuesday"`
	EveryWednesday string `json:"everyWednesday"`
	EveryThursday  string `json:"everyThursday"`
	EveryFriday    string `json:"everyFriday"`
	EverySaturday  string `json:"everySaturday"`
	EverySunday    string `json:"everySunday"`
}

/*
 * @apiEnum: CycleStatusEnum
 */
type CycleStatusEnum struct {
	Active  string `json:"active"`
	Expired string `json:"expired"`
	Frozen  string `json:"frozen"`
}

/*
 * @apiEnum: MedicineAvailabilityEnum
 */
type MedicineAvailabilityEnum struct {
	Hospital           string `json:"hospital"`
	Pharmacy           string `json:"pharmacy"`
	ByPrescriptionOnly string `json:"byPrescriptionOnly"`
	OnlineStore        string `json:"onlineStore"`
}

/*
 * @apiEnum: TodoStatusEnum
 */
type TodoStatusEnum struct {
	Active string `json:"active"`
	Done   string `json:"done"`
}

/*
 * @apiEnum: SectionTypeEnum
 */
type SectionTypeEnum struct {
	All      string `json:"all"`
	Parent   string `json:"parent"`
	Children string `json:"children"`
}

/*
 * @apiEnum: KeikkalaPeymentTypeEnum
 */
type KeikkalaPeymentTypeEnum struct {
	PaySoon string `json:"paySoon"`
	Bonus   string `json:"bonus"`
	Nothing string `json:"nothing"`
}

/*
 * @apiEnum: KeikkalaShiftNameEnum
 */
type KeikkalaShiftNameEnum struct {
	Morning string `json:"morning"`
	Evening string `json:"evening"`
	Night   string `json:"night"`
}

/*
 * @apiEnum: NotificationTypeEnum
 */
type NotificationTypeEnum struct {
	All          string `json:"all"`
	Notification string `json:"notification"`
	Request      string `json:"request"`
}

/*
 * @apiEnum: VehicleTypeEnum
 */
type VehicleTypeEnum struct {
	Car                  string `json:"car"`
	Bicycle              string `json:"bicycle"`
	PublicTransportation string `json:"public_transportation"`
}

/*
 * @apiEnum: ShiftVehicleTypeEnum
 */
type ShiftVehicleTypeEnum struct {
	Own                  string `json:"own"`
	Company              string `json:"company"`
	PublicTransportation string `json:"public-transportation"`
}

/*
 * @apiEnum: ShiftStartLocationTypeEnum
 */
type ShiftStartLocationTypeEnum struct {
	Office        string `json:"office"`
	OtherLocation string `json:"other-location"`
}

/*
 * @apiEnum: ShiftStatusEnum
 */
type ShiftStatusEnum struct {
	NotStarted string `json:"not-started"`
	Started    string `json:"started"`
	Ended      string `json:"ended"`
}

/*
 * @apiEnum: VisitStatusEnum
 */
type VisitStatusEnum struct {
	NotStarted string `json:"not-started"`
	Started    string `json:"started"`
	Ended      string `json:"ended"`
	Cancelled  string `json:"cancelled"`
	Delayed    string `json:"delayed"`
	Paused     string `json:"paused"`
	Resumed    string `json:"resumed"`
}

/*
 * @apiEnum: RolesCoreEnum
 */
type RolesCoreEnum struct {
	Owner          string `json:"Owner"`
	StateManager   string `json:"State-Manager"`
	CityManager    string `json:"City-Manager"`
	SectionManager string `json:"Section-Manager"`
	TeamManager    string `json:"Team-Manager"`
}

/*
 * @apiEnum: PermissionsEnum
 */
type PermissionsEnum struct {
	Dashboard                                string `json:"dashboard"`
	CommunicationOrganizationChat            string `json:"communication-organization-chat"`
	CommunicationCustomerChat                string `json:"communication-customer-chat"`
	CyclesIncomingCycle                      string `json:"cycles-incoming-cycle"`
	CyclesCurrentCycle                       string `json:"cycles-current-cycle"`
	CyclesViewAllCycles                      string `json:"cycles-view-all-cycles"`
	CyclesCreateNewCycle                     string `json:"cycles-create-new-cycle"`
	KeikkalaList                             string `json:"keikkala-list"`
	KeikkalaAddKeikkalaShift                 string `json:"keikkala-add-keikkala-shift"`
	TransportationDetails                    string `json:"transportation-details"`
	ReportsSystemReports                     string `json:"reports-system-reports"`
	StaffsViewAllStaffs                      string `json:"staffs-view-all-staffs"`
	StaffsCreateNewStaff                     string `json:"staffs-create-new-staff"`
	CustomersViewAllCustomers                string `json:"customers-view-all-customers"`
	CustomersCreateNewCustomer               string `json:"customers-create-new-customer"`
	StaffclubDetails                         string `json:"staffclub-details"`
	SectionsCreateNewSection                 string `json:"sections-create-new-section"`
	SectionsViewAllSections                  string `json:"sections-view-all-sections"`
	TeamsCreateNewTeam                       string `json:"teams-create-new-team"`
	TeamsViewAllTeams                        string `json:"teams-view-all-teams"`
	ServicegradesCreateNewServiceGrade       string `json:"servicegrades-create-new-service-grade"`
	ServicegradesViewAllServiceGrades        string `json:"servicegrades-view-all-service-grades"`
	StaffpermissionsCreateNewStaffPermission string `json:"staffpermissions-create-new-staff-permission"`
	StaffpermissionsViewAllStaffPermissions  string `json:"staffpermissions-view-all-staff-permissions"`
	RolepermissionsCreateNewRolePermission   string `json:"rolepermissions-create-new-role-permission"`
	RolepermissionsViewAllRolePermissions    string `json:"rolepermissions-view-all-role-permissions"`
	LanguageskillsCreateNewLanguageSkill     string `json:"languageskills-create-new-language-skill"`
	LanguageskillsViewAllLanguageSkills      string `json:"languageskills-view-all-language-skills"`
	ContracttypesCreateNewContractType       string `json:"contracttypes-create-new-contract-type"`
	ContracttypesViewAllContractTypes        string `json:"contracttypes-view-all-contract-types"`
	EquipmentsCreateNewEquipment             string `json:"equipments-create-new-equipment"`
	EquipmentsViewAllEquipments              string `json:"equipments-view-all-equipments"`
	LimitationsCreateNewLimitation           string `json:"limitations-create-new-limitation"`
	LimitationsViewAllLimitations            string `json:"limitations-view-all-limitations"`
	DiagnosesCreateNewDiagnose               string `json:"diagnoses-create-new-diagnose"`
	DiagnosesViewAllDiagnoses                string `json:"diagnoses-view-all-diagnoses"`
	LicensesCreateNewLicense                 string `json:"licenses-create-new-license"`
	LicensesViewAllLicenses                  string `json:"licenses-view-all-licenses"`
	ArchivesCreateNewArchive                 string `json:"archives-create-new-archive"`
	ArchivesViewAllArchives                  string `json:"archives-view-all-archives"`
	MedicinesCreateNewMedicine               string `json:"medicines-create-new-medicine"`
	MedicinesViewAllMedicines                string `json:"medicines-view-all-medicines"`
	ServicesCreateNewService                 string `json:"services-create-new-service"`
	ServicesViewAllServices                  string `json:"services-view-all-services"`
	ServicetypesCreateNewServiceType         string `json:"servicetypes-create-new-service-type"`
	ServicetypesViewAllServiceTypes          string `json:"servicetypes-view-all-service-types"`
	ServiceoptionsCreateNewServiceOption     string `json:"serviceoptions-create-new-service-option"`
	ServiceoptionsViewAllServiceOptions      string `json:"serviceoptions-view-all-service-options"`
	Notifications                            string `json:"notifications"`
	Requests                                 string `json:"requests"`
	Email                                    string `json:"email"`
	QuizzesAllQuizzes                        string `json:"quizzes-all-quizzes"`
	TicketsCreateNewTicket                   string `json:"tickets-create-new-ticket"`
	TicketsViewAllTickets                    string `json:"tickets-view-all-tickets"`
	SettingGeneral                           string `json:"setting-general"`
	SettingCycle                             string `json:"setting-cycle"`
	SettingNotification                      string `json:"setting-notification"`
	SettingSecurity                          string `json:"setting-security"`
	SettingKeikkala                          string `json:"setting-keikkala"`
	SettingStaffclub                         string `json:"setting-staffclub"`
}
