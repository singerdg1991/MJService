package migrations

import (
	"database/sql"

	"github.com/hoitek/Maja-Service/cmd/migrator"
	"github.com/hoitek/Maja-Service/cmd/migrator/drivers"
)

// Apply do migration functionalities
func Apply(pDB *sql.DB) (*drivers.Postgres, error) {
	m, err := migrator.New(drivers.NewPostgresDriver(pDB))
	if err != nil {
		return nil, err
	}

	// Migrate down all migrations
	// NewPushesMigration1687787691().MigrateDown()
	// NewStaffChatMessagesMigration1687787690().MigrateDown()
	// NewStaffChatsMigration1686787689().MigrateDown()
	// NewCyclePickupShiftCustomersMigration1688787688().MigrateDown()
	// NewCycleChatMessagesMigration1687787687().MigrateDown()
	// NewCycleChatsMigration1686787686().MigrateDown()
	// NewCyclePickupShiftTodosMigration1685787685().MigrateDown()
	// NewCycleShiftCustomerHomeKeysMigration1684787684().MigrateDown()
	// NewIncomingCycleShiftsMigration1684787682().MigrateDown()
	// NewCyclePickupShiftsMigration1684787681().MigrateDown()
	// NewCycleShiftsMigration1684787683().MigrateDown()
	// NewCycleNextStaffTypesMigration1684787680().MigrateDown()

	// NewCustomerStatusLogsMigration1684787679().MigrateDown()
	// NewCustomersDiagnosesMigration1684787678().MigrateDown()
	// NewCustomersRelativesMigration1684787677().MigrateDown()
	// NewCustomerContractualMobilityRestrictionLogsMigration1684787676().MigrateDown()
	// NewCustomerRelativesMigration1684787675().MigrateDown()
	// NewUserLibrariesMigration1684787674().MigrateDown()
	// NewUserOtherAttachmentsMigration1684787673().MigrateDown()
	// NewCustomersMedicinesMigration1684787672().MigrateDown()
	// NewStaffLicensesMigration1684787671().MigrateDown()
	// NewStaffLimitationsMigration1684787669().MigrateDown()
	// NewEvaluationsMigration1684787669().MigrateDown()
	// NewKeikkalaShiftsMigration1684787668().MigrateDown()
	// NewEmailsMigration1684787667().MigrateDown()
	// NewQuizQuestionAnswersMigration1684787666().MigrateDown()
	// NewStaffClubHolidaysMigration1684787665().MigrateDown()
	// NewStaffClubWarningsMigration1684787664().MigrateDown()
	// NewStaffClubAttentionsMigration1684787663().MigrateDown()
	// NewStaffClubGracesMigration1684787662().MigrateDown()
	// NewQuizParticipantsMigration1684787661().MigrateDown()
	// NewQuizQuestionOptionsMigration1684787660().MigrateDown()
	// NewQuizQuestionsMigration1684787659().MigrateDown()
	// NewQuizzesMigration1684787658().MigrateDown()
	// NewNotificationsMigration1684787657().MigrateDown()
	// NewTodosMigration1684787656().MigrateDown()
	// NewTicketsMigration1684787655().MigrateDown()
	// NewCustomersServicesMigration1684787654().MigrateDown()
	// NewCustomersAbsencesMigration1684787653().MigrateDown()
	// NewCustomersCreditDetailsMigration1684787652().MigrateDown()
	// NewUsersAddRelationToCustomerIdAndStaffIdMigration1684787651().MigrateDown()
	// NewPrescriptionsMigration1684787628().MigrateDown()
	// NewAddressesMigration1684787626().MigrateDown()
	// NewCustomersMigration1684787650().MigrateDown()

	// NewCycleStaffTypesMigration1684787649().MigrateDown()
	// NewCyclesAddColumnsMigration1684787645().MigrateDown()

	// NewArchivesMigration1684787644().MigrateDown()
	// NewPunishmentsMigration1684787643().MigrateDown()
	// NewRewardsMigration1684787642().MigrateDown()
	// NewEquipmentsMigration1684787641().MigrateDown()
	// NewServiceGradesMigration1684787640().MigrateDown()
	// NewStaffsMigration1684787624().MigrateDown()
	// NewStaffAbsencesMigration1684787624().MigrateDown()
	// NewContractTypesMigration1684787639().MigrateDown()
	// NewLicensesMigration1684787638().MigrateDown()
	// NewLimitationsMigration1684787637().MigrateDown()
	// NewPermissionsAddDescriptionColumnMigration1684787636().MigrateDown()
	// NewLanguageSkillsAddDescriptionColumnMigration1684787635().MigrateDown()
	// NewTrashesMigration1684787634().MigrateDown()
	// NewOTPsMigration1684787633().MigrateDown()

	// NewCyclesMigration1684787632().MigrateDown()

	// NewServiceOptionsMigration1684787648().MigrateDown()
	// NewServiceTypesMigration1684787647().MigrateDown()
	// NewServicesMigration1684787646().MigrateDown()
	// NewDiagnosesMigration1684787629().MigrateDown()
	// NewMedicinesMigration1684787627().MigrateDown()
	// NewVehiclesMigration1684787623().MigrateDown()
	// NewSeedUsers1684787622().MigrateDown()
	// NewUsersRolesMigration1684787622().MigrateDown()
	// NewUsersMigration1684787622().MigrateDown()
	// NewPaymentTypesMigration1684787621().MigrateDown()
	// NewAbilitiesMigration1684787620().MigrateDown()
	// NewVehicleTypesMigration1684787619().MigrateDown()
	// NewCompaniesMigration1684787618().MigrateDown()
	// NewSectionsMigration1684787617().MigrateDown()
	// NewShiftTypesMigration1684787615().MigrateDown()
	// NewPermissionsMigration1684787613().MigrateDown()
	// NewStaffTypesMigration1684787612().MigrateDown()
	// NewLanguageSkillsMigration1684787611().MigrateDown()
	// NewGearTypesMigration1684787610().MigrateDown()
	// NewCityMigration1684787609().MigrateDown()
	// NewRolesPermissionsMigration1684787602().MigrateDown()
	// NewPermissionMigration1684787601().MigrateDown()
	// NewRoleMigration1684787600().MigrateDown()

	// Add all migrations
	m.AddMigrations(
		NewRoleMigration1684787600(),
		NewPermissionMigration1684787601(),
		NewRolesPermissionsMigration1684787602(),
		NewCityMigration1684787609(),
		NewGearTypesMigration1684787610(),
		NewLanguageSkillsMigration1684787611(),
		NewStaffTypesMigration1684787612(),
		NewPermissionsMigration1684787613(),
		NewShiftTypesMigration1684787615(),
		NewSectionsMigration1684787617(),
		NewCompaniesMigration1684787618(),
		NewVehicleTypesMigration1684787619(),
		NewAbilitiesMigration1684787620(),
		NewPaymentTypesMigration1684787621(),
		NewUsersMigration1684787622(),
		NewUsersRolesMigration1684787622(),
		NewSeedUsers1684787622(),
		NewVehiclesMigration1684787623(),
		NewMedicinesMigration1684787627(),
		NewDiagnosesMigration1684787629(),
		NewServicesMigration1684787646(),
		NewServiceTypesMigration1684787647(),
		NewServiceOptionsMigration1684787648(),
		NewCyclesMigration1684787632(),
		NewOTPsMigration1684787633(),
		NewTrashesMigration1684787634(),
		NewLanguageSkillsAddDescriptionColumnMigration1684787635(),
		NewPermissionsAddDescriptionColumnMigration1684787636(),
		NewLimitationsMigration1684787637(),
		NewLicensesMigration1684787638(),
		NewContractTypesMigration1684787639(),
		NewStaffsMigration1684787624(),
		NewStaffAbsencesMigration1684787624(),
		NewServiceGradesMigration1684787640(),
		NewEquipmentsMigration1684787641(),
		NewRewardsMigration1684787642(),
		NewPunishmentsMigration1684787643(),
		NewArchivesMigration1684787644(),
		NewCyclesAddColumnsMigration1684787645(),
		NewCycleStaffTypesMigration1684787649(),
		NewCustomersMigration1684787650(),
		NewPrescriptionsMigration1684787628(),
		NewAddressesMigration1684787626(),
		NewUsersAddRelationToCustomerIdAndStaffIdMigration1684787651(),
		NewCustomersCreditDetailsMigration1684787652(),
		NewCustomersAbsencesMigration1684787653(),
		NewCustomersServicesMigration1684787654(),
		NewTicketsMigration1684787655(),
		NewTodosMigration1684787656(),
		NewNotificationsMigration1684787657(),
		NewQuizzesMigration1684787658(),
		NewQuizQuestionsMigration1684787659(),
		NewQuizQuestionOptionsMigration1684787660(),
		NewQuizParticipantsMigration1684787661(),
		NewStaffClubGracesMigration1684787662(),
		NewStaffClubAttentionsMigration1684787663(),
		NewStaffClubWarningsMigration1684787664(),
		NewStaffClubHolidaysMigration1684787665(),
		NewQuizQuestionAnswersMigration1684787666(),
		NewEmailsMigration1684787667(),
		NewKeikkalaShiftsMigration1684787668(),
		NewEvaluationsMigration1684787669(),
		NewStaffLimitationsMigration1684787669(),
		NewStaffLicensesMigration1684787671(),
		NewCustomersMedicinesMigration1684787672(),
		NewUserOtherAttachmentsMigration1684787673(),
		NewUserLibrariesMigration1684787674(),
		NewCustomerRelativesMigration1684787675(),
		NewCustomerContractualMobilityRestrictionLogsMigration1684787676(),
		NewCustomersRelativesMigration1684787677(),
		NewCustomersDiagnosesMigration1684787678(),
		NewCustomerStatusLogsMigration1684787679(),
		NewCycleNextStaffTypesMigration1684787680(),
		NewCycleShiftsMigration1684787683(),
		NewCyclePickupShiftsMigration1684787681(),
		NewIncomingCycleShiftsMigration1684787682(),
		NewCycleShiftCustomerHomeKeysMigration1684787684(),
		NewCyclePickupShiftTodosMigration1685787685(),
		NewCycleChatsMigration1686787686(),
		NewCycleChatMessagesMigration1687787687(),
		NewCyclePickupShiftCustomersMigration1688787688(),
		NewStaffChatsMigration1686787689(),
		NewStaffChatMessagesMigration1687787690(),
		NewPushesMigration1687787691(),
	)
	if err := m.MigrateUp(false); err != nil {
		return nil, err
	}

	return m, nil
}
