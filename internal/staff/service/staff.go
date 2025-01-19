package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/staff/constants"
	"github.com/hoitek/Maja-Service/internal/staff/domain"
	"github.com/hoitek/Maja-Service/internal/staff/models"
	"github.com/hoitek/Maja-Service/internal/staff/ports"
	"github.com/hoitek/Maja-Service/storage"
)

type StaffService struct {
	PostgresRepository ports.StaffRepositoryPostgresDB
	MongoDBRepository  ports.StaffRepositoryMongoDB
	MinIOStorage       *storage.MinIO
}

func NewStaffService(pDB ports.StaffRepositoryPostgresDB, mDB ports.StaffRepositoryMongoDB, m *storage.MinIO) StaffService {
	go minio.SetupMinIOStorage(constants.STAFF_BUCKET_NAME, m)
	return StaffService{
		PostgresRepository: pDB,
		MongoDBRepository:  mDB,
		MinIOStorage:       m,
	}
}

func (s *StaffService) Query(q *models.StaffsQueryRequestParams) (*restypes.QueryResponse, error) {
	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 1, 10, q.Limit)

	staffs, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(q)
	if err != nil {
		return nil, err
	}

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      staffs,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *StaffService) QueryLicenses(q *models.StaffsQueryLicensesRequestParams) (*restypes.QueryResponse, error) {
	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 1, 10, q.Limit)

	licenses, count, err := s.PostgresRepository.QueryLicenses(q)
	if err != nil {
		return nil, err
	}

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      licenses,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *StaffService) QueryAbsences(dataModel *models.StaffsQueryAbsencesRequestParams) (*restypes.QueryResponse, error) {
	dataModel.Page = exp.TerIf(dataModel.Page < 1, 1, dataModel.Page)
	dataModel.Limit = exp.TerIf(dataModel.Limit < 1, 10, dataModel.Limit)

	absences, count, err := s.PostgresRepository.QueryAbsences(dataModel)
	if err != nil {
		return nil, err
	}

	page := dataModel.Page
	limit := dataModel.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      absences,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *StaffService) FindByUserID(userId int) (*domain.Staff, error) {
	staffs, err := s.PostgresRepository.Query(&models.StaffsQueryRequestParams{
		Filters: models.StaffFilterType{
			UserId: filters.FilterValue[int]{
				Op:    operators.EQUALS,
				Value: userId,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(staffs) < 1 {
		return nil, errors.New("user not found")
	}
	return staffs[0], nil
}

func (s *StaffService) FindByID(id int) (*domain.Staff, error) {
	r, err := s.Query(&models.StaffsQueryRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("staff not found")
	}
	staffs := r.Items.([]*domain.Staff)
	return staffs[0], nil
}

func (s *StaffService) Delete(payload *models.StaffsDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *StaffService) DeleteLicenses(payload *models.StaffsDeleteLicensesRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, _, err := s.PostgresRepository.DeleteLicenses(payload)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *StaffService) CreateOrUpdateContract(payload *models.StaffsCreateOrUpdateContractRequestBody) (*domain.Staff, error) {
	return s.PostgresRepository.CreateOrUpdateContract(payload)
}

func (s *StaffService) ExportToCsvAndSave(staffs []*domain.Staff, dest string) (string, error) {
	// Check if there is data to export
	if len(staffs) < 1 {
		return "", errors.New("no data to export")
	}

	var firstRow = map[string]interface{}{}
	if len(staffs) > 0 {
		row, err := staffs[0].ToMap()
		if err != nil {
			return "", err
		}
		firstRow = row
	}
	if len(firstRow) < 1 {
		return "", errors.New("no data to export")
	}

	// Create the file
	file, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get the header row dynamically
	dataMapKeys := make([]string, 0)
	var headerKeys []string
	utils.JoinMapKeysWithDot(firstRow, "", &headerKeys)
	for _, fieldName := range headerKeys {
		if !utils.Contains(dataMapKeys, fieldName) {
			dataMapKeys = append(dataMapKeys, fieldName)
		}
	}

	// Write the data rows dynamically
	for _, staff := range staffs {
		for _, key := range dataMapKeys {
			// Get the row
			staffRow, err := staff.ToMap()
			if err != nil {
				return "", err
			}

			// Get the value
			value := utils.GetMapValueFromDotJoinedKeys(staffRow, key)

			// Create a row
			row := make([]string, 0)
			row = append(row, key)
			row = append(row, fmt.Sprintf("%v", value))

			// Write the header row
			writer.Write(row)
		}
	}

	// Return the file path
	return fmt.Sprintf("/%s", dest), nil
}

func (s *StaffService) CreateLicenses(staff *domain.Staff, payload *models.StaffsCreateLicensesRequestBody) (*domain.StaffLicensesRes, error) {
	_, staffLicense, err := s.PostgresRepository.CreateLicenses(staff, payload)
	return staffLicense, err
}

func (s *StaffService) HasLicense(staffId uint, licenseId uint) (bool, error) {
	return s.PostgresRepository.HasLicense(staffId, licenseId)
}

func (s *StaffService) CreateEmptyStaffForUser(userId int) (*domain.Staff, error) {
	return s.PostgresRepository.CreateEmptyStaffForUser(userId)
}

func (s *StaffService) GenerateStaffOrganizationNumber() (string, error) {
	return s.PostgresRepository.GenerateStaffOrganizationNumber()
}

func (s *StaffService) CreateAbsences(staff *domain.Staff, payload *models.StaffsCreateAbsencesRequestBody) (*domain.StaffAbsencesQueryRes, error) {
	return s.PostgresRepository.CreateAbsences(staff, payload)
}

func (s *StaffService) FindStaffLicenseByID(id int) (*domain.StaffLicensesRes, error) {
	return s.PostgresRepository.FindStaffLicenseByID(id)
}

func (s *StaffService) HasLicenseExcept(staffID uint, licenseID uint, staffLicenseID uint) (bool, error) {
	return s.PostgresRepository.HasLicenseExcept(staffID, licenseID, staffLicenseID)
}

func (s *StaffService) UpdateLicense(staffLicense *domain.StaffLicensesRes, payload *models.StaffsUpdateLicenseRequestBody) (*domain.StaffLicensesRes, error) {
	_, staffLicense, err := s.PostgresRepository.UpdateLicense(staffLicense, payload)
	return staffLicense, err
}

func (s *StaffService) FindStaffAbsenceByID(id int) (*domain.StaffAbsenceRes, error) {
	return s.PostgresRepository.FindStaffAbsenceByID(id)
}

func (s *StaffService) UpdateAbsence(staffAbsence *domain.StaffAbsenceRes, payload *models.StaffsUpdateAbsenceRequestBody) (*domain.StaffAbsenceRes, error) {
	_, staffAbsence, err := s.PostgresRepository.UpdateAbsence(staffAbsence, payload)
	return staffAbsence, err
}

func (s *StaffService) DeleteAbsences(payload *models.StaffsDeleteAbsencesRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, _, err := s.PostgresRepository.DeleteAbsences(payload)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *StaffService) QueryProfile(query *models.StaffsQueryProfileRequestParams) (*domain.StaffProfile, error) {
	profile, err := s.PostgresRepository.QueryProfile(query)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *StaffService) CreateStaff(userId int64, payload *models.StaffsCreatePersonalInfoRequestBody) (*domain.Staff, error) {
	return s.PostgresRepository.CreateStaff(userId, payload)
}

func (s *StaffService) UpdateStaff(userId int64, payload *models.StaffsUpdatePersonalInfoRequestBody) (*domain.Staff, error) {
	return s.PostgresRepository.UpdateStaff(userId, payload)
}

func (s *StaffService) GenerateRegistrationNumber() (string, error) {
	// Generate random 6 unique digits with rand.Intn
	src := rand.NewSource(time.Now().UnixNano())

	// Create a new rand.Rand using the source
	r := rand.New(src)

	randomNum := r.Intn(999999-100000) + 100000

	// T- for staff, A- for subcontractor, K- for keikkala
	return fmt.Sprintf("S-%d", randomNum), nil
}

func (s *StaffService) UpdateAbsenceAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffAbsencesQueryRes, error) {
	return s.PostgresRepository.UpdateAbsenceAttachments(previousAttachments, attachments, id)
}

func (s *StaffService) UpdateLicenseAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffLicensesRes, error) {
	return s.PostgresRepository.UpdateLicenseAttachments(previousAttachments, attachments, id)
}

func (s *StaffService) CreateOtherAttachments(staff *domain.Staff, payload *models.StaffsCreateOtherAttachmentsRequestBody) (*domain.StaffOtherAttachment, error) {
	return s.PostgresRepository.CreateOtherAttachments(staff, payload)
}

func (s *StaffService) UpdateStaffOtherAttachments(attachments []*types.UploadMetadata, id int64) (*domain.StaffOtherAttachment, error) {
	return s.PostgresRepository.UpdateStaffOtherAttachments(attachments, id)
}

func (s *StaffService) QueryOtherAttachments(dataModel *models.StaffsQueryOtherAttachmentsRequestParams) (*restypes.QueryResponse, error) {
	dataModel.Page = exp.TerIf(dataModel.Page < 1, 1, dataModel.Page)
	dataModel.Limit = exp.TerIf(dataModel.Limit < 1, 10, dataModel.Limit)

	otherAttachments, count, err := s.PostgresRepository.QueryOtherAttachments(dataModel)
	if err != nil {
		return nil, err
	}

	page := dataModel.Page
	limit := dataModel.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      otherAttachments,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *StaffService) FindStaffOtherAttachmentByID(id int) (*domain.StaffOtherAttachment, error) {
	r, err := s.QueryOtherAttachments(&models.StaffsQueryOtherAttachmentsRequestParams{
		ID:    id,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("staff other attachment not found")
	}
	staffOtherAttachments, ok := r.Items.([]*domain.StaffOtherAttachment)
	if !ok {
		return nil, errors.New("staff other attachment not found")
	}
	return staffOtherAttachments[0], nil
}

func (s *StaffService) UpdateStaffOtherAttachment(staffOtherAttachment *domain.StaffOtherAttachment, payload *models.StaffsUpdateOtherAttachmentRequestBody) (*domain.StaffOtherAttachment, error) {
	return s.PostgresRepository.UpdateStaffOtherAttachment(staffOtherAttachment, payload)
}

func (s *StaffService) DeleteStaffOtherAttachments(payload *models.StaffsDeleteOtherAttachmentsRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.DeleteStaffOtherAttachments(payload)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *StaffService) CreateLibraries(staff *domain.Staff, payload *models.StaffsCreateLibrariesRequestBody) (*domain.StaffLibrary, error) {
	return s.PostgresRepository.CreateLibraries(staff, payload)
}

func (s *StaffService) UpdateStaffLibraries(attachments []*types.UploadMetadata, id int64) (*domain.StaffLibrary, error) {
	return s.PostgresRepository.UpdateStaffLibraries(attachments, id)
}

func (s *StaffService) QueryLibraries(dataModel *models.StaffsQueryLibrariesRequestParams) (*restypes.QueryResponse, error) {
	dataModel.Page = exp.TerIf(dataModel.Page < 1, 1, dataModel.Page)
	dataModel.Limit = exp.TerIf(dataModel.Limit < 1, 10, dataModel.Limit)

	libraries, count, err := s.PostgresRepository.QueryLibraries(dataModel)
	if err != nil {
		return nil, err
	}

	page := dataModel.Page
	limit := dataModel.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      libraries,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *StaffService) FindStaffLibraryByID(id int) (*domain.StaffLibrary, error) {
	r, err := s.QueryLibraries(&models.StaffsQueryLibrariesRequestParams{
		ID:    id,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("staff library not found")
	}
	staffLibraries, ok := r.Items.([]*domain.StaffLibrary)
	if !ok {
		return nil, errors.New("staff library not found")
	}
	return staffLibraries[0], nil
}

func (s *StaffService) UpdateStaffLibrary(staffLibrary *domain.StaffLibrary, payload *models.StaffsUpdateLibraryRequestBody) (*domain.StaffLibrary, error) {
	return s.PostgresRepository.UpdateStaffLibrary(staffLibrary, payload)
}

func (s *StaffService) DeleteStaffLibraries(payload *models.StaffsDeleteLibrariesRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.DeleteStaffLibraries(payload)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *StaffService) UpdateStaffAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Staff, error) {
	return s.PostgresRepository.UpdateStaffAttachments(previousAttachments, attachments, id)
}

// QueryChats queries the chats based on the provided data model.
//
// Parameters:
// - dataModel: a pointer to models.StaffsQueryChatsRequestParams containing the query parameters.
// Returns:
// - *restypes.QueryResponse: a pointer to the query response.
// - error: an error if the query process encounters any issues.
func (s *StaffService) QueryChats(dataModel *models.StaffsQueryChatsRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying Chats", dataModel)
	chats, err := s.PostgresRepository.QueryChats(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountChats(dataModel)
	if err != nil {
		return nil, err
	}

	dataModel.Page = exp.TerIf(dataModel.Page < 1, 1, dataModel.Page)
	dataModel.Limit = exp.TerIf(dataModel.Limit < 10, 1, dataModel.Limit)

	page := dataModel.Page
	limit := dataModel.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))
	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      chats,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// QueryChatMessages queries the chat messages based on the provided data model.
//
// Parameters:
// - dataModel: a pointer to models.StaffsQueryChatMessagesRequestParams containing the query parameters.
// Returns:
// - *restypes.QueryResponse: a pointer to the query response.
// - error: an error if the query process encounters any issues.
func (s *StaffService) QueryChatMessages(dataModel *models.StaffsQueryChatMessagesRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying Chat Messages", dataModel)
	messages, err := s.PostgresRepository.QueryChatMessages(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountChatMessages(dataModel)
	if err != nil {
		return nil, err
	}

	dataModel.Page = exp.TerIf(dataModel.Page < 1, 1, dataModel.Page)
	dataModel.Limit = exp.TerIf(dataModel.Limit < 10, 1, dataModel.Limit)

	page := dataModel.Page
	limit := dataModel.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))
	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      messages,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// CreateChatMessage creates a new chat message in the StaffService.
//
// It takes a payload of type *models.StaffsCreateChatMessageRequestBody, which contains the staff chat ID,
// sender user ID, recipient user ID, message, and attachments.
// It returns a pointer to a domain.StaffChatMessage and an error.
func (s *StaffService) CreateChatMessage(payload *models.StaffsCreateChatMessageRequestBody) (*domain.StaffChatMessage, error) {
	return s.PostgresRepository.CreateChatMessage(payload)
}

// UpdateChatMessageAttachments updates the attachments of a chat message.
//
// Parameters:
// - previousAttachments: a slice of types.UploadMetadata representing the previous attachments.
// - attachments: a slice of pointers to types.UploadMetadata representing the new attachments.
// - id: an int64 representing the ID of the chat message.
// Returns:
// - *domain.StaffChatMessage: a pointer to the updated chat message.
// - error: an error if the update process encounters any issues.
func (s *StaffService) UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.StaffChatMessage, error) {
	return s.PostgresRepository.UpdateChatMessageAttachments(previousAttachments, attachments, id)
}
