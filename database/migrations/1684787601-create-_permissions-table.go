package migrations

import (
	"fmt"
	"log"

	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/permissions"
)

type PermissionMigration1684787601 struct {
}

func NewPermissionMigration1684787601() *PermissionMigration1684787601 {
	return &PermissionMigration1684787601{}
}

func generateValues(permissions map[string]string) string {
	values := ""
	var i int
	for p, t := range permissions {
		if i > 0 {
			values += ",\n"
		}
		values += fmt.Sprintf("(%d, '%s', '%s', '2020-01-01', '2020-01-01', '2020-01-01')", i+1, p, t)
		i++
	}
	return values
}

func (m *PermissionMigration1684787601) MigrateUp() error {
	perms := map[string]string{}
	perms[permissions.DASHBOARD] = "dashboard"
	perms[permissions.COMMUNICATION_ORGANIZATION_CHAT] = "Organization Chat"
	perms[permissions.COMMUNICATION_CUSTOMER_CHAT] = "Customer Chat"
	perms[permissions.CYCLES_INCOMING_CYCLE] = "Incoming Cycle"
	perms[permissions.CYCLES_CURRENT_CYCLE] = "Current Cycle"
	perms[permissions.CYCLES_VIEW_ALL_CYCLES] = "View All Cycles"
	perms[permissions.CYCLES_CREATE_NEW_CYCLE] = "Create New Cycle"
	perms[permissions.KEIKKALA_LIST] = "Keikkala List"
	perms[permissions.KEIKKALA_ADD_KEIKKALA_SHIFT] = "Keikkala Add Keikkala Shift"
	perms[permissions.TRANSPORTATION_DETAILS] = "Transportation Details"
	perms[permissions.REPORTS_SYSTEM_REPORTS] = "System Reports"
	perms[permissions.STAFFS_VIEW_ALL_STAFFS] = "View All Staffs"
	perms[permissions.STAFFS_CREATE_NEW_STAFF] = "Create New Staff"
	perms[permissions.CUSTOMERS_VIEW_ALL_CUSTOMERS] = "View All Customers"
	perms[permissions.CUSTOMERS_CREATE_NEW_CUSTOMER] = "Create New Customer"
	perms[permissions.STAFF_CLUB_DETAILS] = "Staff Club Details"
	perms[permissions.SECTIONS_CREATE_NEW_SECTION] = "Create New Section"
	perms[permissions.SECTIONS_VIEW_ALL_SECTIONS] = "View All Sections"
	perms[permissions.TEAMS_CREATE_NEW_TEAM] = "Create New Team"
	perms[permissions.TEAMS_VIEW_ALL_TEAMS] = "View All Teams"
	perms[permissions.SERVICE_GRADES_CREATE_NEW_SERVICE_GRADE] = "Create New Service Grade"
	perms[permissions.SERVICE_GRADES_VIEW_ALL_SERVICE_GRADES] = "View All Service Grades"
	perms[permissions.STAFF_PERMISSIONS_CREATE_NEW_STAFF_PERMISSION] = "Create New Staff Permission"
	perms[permissions.STAFF_PERMISSIONS_VIEW_ALL_STAFF_PERMISSIONS] = "View All Staff Permissions"
	perms[permissions.ROLE_PERMISSIONS_CREATE_NEW_ROLE_PERMISSION] = "Create New Role Permission"
	perms[permissions.ROLE_PERMISSIONS_VIEW_ALL_ROLE_PERMISSIONS] = "View All Role Permissions"
	perms[permissions.LANGUAGE_SKILLS_CREATE_NEW_LANGUAGE_SKILL] = "Create New Language Skill"
	perms[permissions.LANGUAGE_SKILLS_VIEW_ALL_LANGUAGE_SKILLS] = "View All Language Skills"
	perms[permissions.CONTRACT_TYPES_CREATE_NEW_CONTRACT_TYPE] = "Create New Contract Type"
	perms[permissions.CONTRACT_TYPES_VIEW_ALL_CONTRACT_TYPES] = "View All Contract Types"
	perms[permissions.EQUIPMENTS_CREATE_NEW_EQUIPMENT] = "Create New Equipment"
	perms[permissions.EQUIPMENTS_VIEW_ALL_EQUIPMENTS] = "View All Equipments"
	perms[permissions.LIMITATIONS_CREATE_NEW_LIMITATION] = "Create New Limitation"
	perms[permissions.LIMITATIONS_VIEW_ALL_LIMITATIONS] = "View All Limitations"
	perms[permissions.DIAGNOSES_CREATE_NEW_DIAGNOSE] = "Create New Diagnose"
	perms[permissions.DIAGNOSES_VIEW_ALL_DIAGNOSES] = "View All Diagnoses"
	perms[permissions.LICENSES_CREATE_NEW_LICENSE] = "Create New License"
	perms[permissions.LICENSES_VIEW_ALL_LICENSES] = "View All Licenses"
	perms[permissions.ARCHIVES_CREATE_NEW_ARCHIVE] = "Create New Archive"
	perms[permissions.ARCHIVES_VIEW_ALL_ARCHIVES] = "View All Archives"
	perms[permissions.MEDICINES_CREATE_NEW_MEDICINE] = "Create New Medicine"
	perms[permissions.MEDICINES_VIEW_ALL_MEDICINES] = "View All Medicines"
	perms[permissions.SERVICES_CREATE_NEW_SERVICE] = "Create New Service"
	perms[permissions.SERVICES_VIEW_ALL_SERVICES] = "View All Services"
	perms[permissions.SERVICE_TYPES_CREATE_NEW_SERVICE_TYPE] = "Create New Service Type"
	perms[permissions.SERVICE_TYPES_VIEW_ALL_SERVICE_TYPES] = "View All Service Types"
	perms[permissions.SERVICE_OPTIONS_CREATE_NEW_SERVICE_OPTION] = "Create New Service Option"
	perms[permissions.SERVICE_OPTIONS_VIEW_ALL_SERVICE_OPTIONS] = "View All Service Options"
	perms[permissions.NOTIFICATIONS] = "Notifications"
	perms[permissions.REQUESTS] = "Requests"
	perms[permissions.EMAIL] = "Email"
	perms[permissions.QUIZZES_ALL_QUIZZES] = "All Quizzes"
	perms[permissions.TICKETS_CREATE_NEW_TICKET] = "Create New Ticket"
	perms[permissions.TICKETS_VIEW_ALL_TICKETS] = "View All Tickets"
	perms[permissions.SETTING_GENERAL] = "General"
	perms[permissions.SETTING_CYCLE] = "Cycle"
	perms[permissions.SETTING_NOTIFICATION] = "Notification"
	perms[permissions.SETTING_SECURITY] = "Security"
	perms[permissions.SETTING_KEIKKALA] = "Keikkala"
	perms[permissions.SETTING_STAFF_CLUB] = "Staff Club"

	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS _permissions (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		ALTER TABLE _permissions ALTER COLUMN id SET DEFAULT nextval('_permissions_id_seq'::regclass);
		ALTER TABLE _permissions ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE _permissions ALTER COLUMN updated_at SET DEFAULT now();
		INSERT INTO _permissions
			(id, name, title, created_at, updated_at, deleted_at)
		VALUES
			%s
		ON CONFLICT(id) DO NOTHING;
		SELECT setval('_permissions_id_seq', (SELECT MAX(id) FROM _permissions));
	`, generateValues(perms))
	database.PostgresDB.Exec(query)
	return nil
}

func (m *PermissionMigration1684787601) MigrateDown() error {
	_, err := database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'PermissionMigration1684787601';`)
	if err != nil {
		log.Printf("-------- Error deleting migration PermissionMigration1684787601: %v\n", err)
	}
	_, err = database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS _permissions;
    `)
	if err != nil {
		log.Printf("-------- Error deleting table _permissions: %v\n", err)
	}
	return nil
}
