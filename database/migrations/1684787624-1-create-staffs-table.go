package migrations

import (
	"github.com/hoitek/Maja-Service/database"
	"log"
)

type StaffsMigration1684787624 struct {
}

func NewStaffsMigration1684787624() *StaffsMigration1684787624 {
	return &StaffsMigration1684787624{}
}

func (m *StaffsMigration1684787624) MigrateUp() error {
	_, err := database.PostgresDB.Exec(`
        CREATE TABLE IF NOT EXISTS staffs (
			id SERIAL PRIMARY KEY,
			userId INT NOT NULL,
		    paymentTypeId INT DEFAULT NULL,
			jobTitle VARCHAR(255) DEFAULT NULL,
		    certificateCode VARCHAR(255) DEFAULT NULL,
			experienceAmount INT DEFAULT NULL,
			experienceAmountUnit VARCHAR(255) DEFAULT NULL,
            isSubcontractor BOOLEAN DEFAULT FALSE,
		    companyRegistrationNumber VARCHAR(255) DEFAULT NULL,
		    organizationNumber VARCHAR(255) DEFAULT NULL,
		    percentLengthInContract INT DEFAULT NULL,
		    hourLengthInContract INT DEFAULT NULL,
		    salary INT DEFAULT NULL,
            vehicleTypes JSONB DEFAULT NULL,
            vehicleLicenseTypes JSONB DEFAULT NULL,
		    trial_time TIMESTAMP DEFAULT NULL,
			attachments JSONB DEFAULT NULL,
		    joined_at TIMESTAMP DEFAULT NULL,
		    contract_started_at TIMESTAMP DEFAULT NULL,
		    contract_expires_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
		    CONSTRAINT fk_user_id FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_payment_type_id FOREIGN KEY (paymentTypeId) REFERENCES paymentTypes(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS staffSections (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			sectionId INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_section_id FOREIGN KEY (sectionId) REFERENCES sections(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS staffShiftTypes (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			shiftTypeId INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_shift_type_id FOREIGN KEY (shiftTypeId) REFERENCES shiftTypes(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS staffContractTypes (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			contractTypeId INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_contract_type_id FOREIGN KEY (contractTypeId) REFERENCES contractTypes(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS staffStaffTypes (
			id SERIAL PRIMARY KEY,
			staffId INT NOT NULL,
			staffTypeId INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL,
			CONSTRAINT fk_staff_id FOREIGN KEY (staffId) REFERENCES staffs(id) ON DELETE CASCADE,
			CONSTRAINT fk_staff_type_id FOREIGN KEY (staffTypeId) REFERENCES staffTypes(id) ON DELETE CASCADE
		);
		ALTER TABLE staffs ALTER COLUMN id SET DEFAULT nextval('staffs_id_seq'::regclass);
		ALTER TABLE staffs ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffs ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE staffs ALTER COLUMN deleted_at SET DEFAULT NULL;
		ALTER TABLE staffSections ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffSections ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE staffShiftTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffShiftTypes ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE staffContractTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffContractTypes ALTER COLUMN updated_at SET DEFAULT now();
		ALTER TABLE staffStaffTypes ALTER COLUMN created_at SET DEFAULT now();
		ALTER TABLE staffStaffTypes ALTER COLUMN updated_at SET DEFAULT now();
    `)
	if err != nil {
		log.Println("****************************************", err)
	}

	return nil
}

func (m *StaffsMigration1684787624) MigrateDown() error {
	database.PostgresDB.Exec(`DELETE FROM _migrations WHERE name = 'StaffsMigration1684787624'`)
	database.PostgresDB.Exec(`
		DELETE FROM userLanguageSkills WHERE userId IN (SELECT userId FROM staffs);
	`)
	database.PostgresDB.Exec(`
		DELETE FROM userLanguageSkills WHERE userId > 10;
	`)
	database.PostgresDB.Exec(`
		DELETE FROM users WHERE id IN (SELECT userId FROM staffs);
	`)
	database.PostgresDB.Exec(`
		DELETE FROM users WHERE id > 10 AND userType = 'staff';
	`)
	_, err := database.PostgresDB.Exec(`
        DROP TABLE IF EXISTS staffSections CASCADE;
		DROP TABLE IF EXISTS staffAbsences CASCADE;
		DROP TABLE IF EXISTS staffShiftTypes CASCADE;
		DROP TABLE IF EXISTS staffContractTypes CASCADE;
		DROP TABLE IF EXISTS staffStaffTypes CASCADE;
		DROP TABLE IF EXISTS staffs CASCADE;
    `)
	if err != nil {
		log.Println("****************************************", err)
	}
	return nil
}
