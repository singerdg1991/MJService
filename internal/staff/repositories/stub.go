package repositories

import (
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
	"github.com/hoitek/Maja-Service/internal/staff/models"
)

type StaffRepositoryStub struct {
	Staffs []*domain.Staff
}

type staffTestCondition struct {
	HasError bool
}

var UserTestCondition *staffTestCondition = &staffTestCondition{}

func NewStaffRepositoryStub() *StaffRepositoryStub {
	return &StaffRepositoryStub{
		Staffs: []*domain.Staff{
			{
				ID: 1,
			},
		},
	}
}

func (r *StaffRepositoryStub) Query(dataModel *models.StaffsQueryRequestParams) ([]*domain.Staff, error) {
	var staffs []*domain.Staff
	for _, v := range r.Staffs {
		if v.ID == uint(dataModel.ID) {
			staffs = append(staffs, v)
			break
		}
	}
	return staffs, nil
}

func (r *StaffRepositoryStub) Count(dataModel *models.StaffsQueryRequestParams) (int64, error) {
	var staffs []*domain.Staff
	for _, v := range r.Staffs {
		if v.ID == uint(dataModel.ID) {
			staffs = append(staffs, v)
			break
		}
	}
	return int64(len(staffs)), nil
}

func (r *StaffRepositoryStub) Migrate() {
	// do stuff
}

func (r *StaffRepositoryStub) Seed(staffs []*domain.Staff) error {
	// do stuff
	return nil
}

func (r *StaffRepositoryStub) Delete(payload *models.StaffsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateOrUpdateContract(payload *models.StaffsCreateOrUpdateContractRequestBody) (*domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) Drop() error {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateLicenses(staff *domain.Staff, payload *models.StaffsCreateLicensesRequestBody) (*domain.Staff, *domain.StaffLicensesRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) HasLicense(staffID uint, licenseID uint) (bool, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateEmptyStaffForUser(userId int) (*domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) GenerateStaffOrganizationNumber() (string, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateAbsences(staff *domain.Staff, payload *models.StaffsCreateAbsencesRequestBody) (*domain.StaffAbsencesQueryRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) DeleteLicenses(payload *models.StaffsDeleteLicensesRequestBody) ([]int64, *domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryLicenses(dataModel *models.StaffsQueryLicensesRequestParams) ([]*domain.StaffLicensesRes, int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) FindStaffLicenseByID(id int) (*domain.StaffLicensesRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) HasLicenseExcept(staffID uint, licenseID uint, staffLicenseID uint) (bool, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateLicense(staffLicense *domain.StaffLicensesRes, payload *models.StaffsUpdateLicenseRequestBody) (*domain.Staff, *domain.StaffLicensesRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryAbsences(dataModel *models.StaffsQueryAbsencesRequestParams) ([]*domain.StaffAbsencesQueryRes, int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) FindStaffAbsenceByID(id int) (*domain.StaffAbsenceRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateAbsence(staffAbsence *domain.StaffAbsenceRes, payload *models.StaffsUpdateAbsenceRequestBody) (*domain.Staff, *domain.StaffAbsenceRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) DeleteAbsences(payload *models.StaffsDeleteAbsencesRequestBody) ([]int64, *domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryProfile(query *models.StaffsQueryProfileRequestParams) (*domain.StaffProfile, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateStaff(userId int64, payload *models.StaffsCreatePersonalInfoRequestBody) (*domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateStaff(userId int64, payload *models.StaffsUpdatePersonalInfoRequestBody) (*domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateAbsenceAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffAbsencesQueryRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateLicenseAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffLicensesRes, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateOtherAttachments(staff *domain.Staff, payload *models.StaffsCreateOtherAttachmentsRequestBody) (*domain.StaffOtherAttachment, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateStaffOtherAttachments(attachments []*types.UploadMetadata, id int64) (*domain.StaffOtherAttachment, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryOtherAttachments(dataModel *models.StaffsQueryOtherAttachmentsRequestParams) ([]*domain.StaffOtherAttachment, int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateStaffOtherAttachment(staffOtherAttachment *domain.StaffOtherAttachment, payload *models.StaffsUpdateOtherAttachmentRequestBody) (*domain.StaffOtherAttachment, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) DeleteStaffOtherAttachments(payload *models.StaffsDeleteOtherAttachmentsRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateLibraries(staff *domain.Staff, payload *models.StaffsCreateLibrariesRequestBody) (*domain.StaffLibrary, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateStaffLibraries(libraries []*types.UploadMetadata, id int64) (*domain.StaffLibrary, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryLibraries(dataModel *models.StaffsQueryLibrariesRequestParams) ([]*domain.StaffLibrary, int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateStaffLibrary(staffLibrary *domain.StaffLibrary, payload *models.StaffsUpdateLibraryRequestBody) (*domain.StaffLibrary, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) DeleteStaffLibraries(payload *models.StaffsDeleteLibrariesRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateStaffAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Staff, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryChats(queries *models.StaffsQueryChatsRequestParams) ([]*domain.StaffChat, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CountChats(queries *models.StaffsQueryChatsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) QueryChatMessages(queries *models.StaffsQueryChatMessagesRequestParams) ([]*domain.StaffChatMessage, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CountChatMessages(queries *models.StaffsQueryChatMessagesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) CreateChatMessage(payload *models.StaffsCreateChatMessageRequestBody) (*domain.StaffChatMessage, error) {
	panic("implement me")
}

func (r *StaffRepositoryStub) UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffChatMessage, error) {
	panic("implement me")
}
