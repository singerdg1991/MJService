package service

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/constants"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/customer/models"
	"github.com/hoitek/Maja-Service/internal/customer/ports"
	"github.com/hoitek/Maja-Service/storage"
)

type CustomerService struct {
	PostgresRepository ports.CustomerRepositoryPostgresDB
	MongoDBRepository  ports.CustomerRepositoryMongoDB
	MinIOStorage       *storage.MinIO
}

func NewCustomerService(pDB ports.CustomerRepositoryPostgresDB, mDB ports.CustomerRepositoryMongoDB, m *storage.MinIO) CustomerService {
	go minio.SetupMinIOStorage(constants.CUSTOMER_BUCKET_NAME, m)
	return CustomerService{
		PostgresRepository: pDB,
		MongoDBRepository:  mDB,
		MinIOStorage:       m,
	}
}

func (s *CustomerService) FindCustomerServicesForSpecificShift(cyclePickupShiftID int64, date time.Time, shiftName string, shiftMorningStartHour int64, shiftMorningEndHour int64, shiftEveningStartHour int64, shiftEveningEndHour int64, shiftNightStartHour int64, shiftNightEndHour int64) ([]*domain.CustomerServices, error) {
	return s.PostgresRepository.FindCustomerServicesForSpecificShift(cyclePickupShiftID, date, shiftName, shiftMorningStartHour, shiftMorningEndHour, shiftEveningStartHour, shiftEveningEndHour, shiftNightStartHour, shiftNightEndHour)
}

func (s *CustomerService) Query(q *models.CustomersQueryRequestParams) (*restypes.QueryResponse, error) {
	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 1, 1, q.Limit)

	customers, err := s.PostgresRepository.Query(q)
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
		Items:      customers,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) CreatePersonalInfo(payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	customer, err := s.PostgresRepository.CreatePersonalInfo(payload)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) UpdateUserInformation(customerID int64, payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	return s.PostgresRepository.UpdateUserInformation(customerID, payload)
}

func (s *CustomerService) UpdateAdditionalInfo(payload *models.CustomersUpdateAdditionalInfoRequestBody) (*domain.Customer, error) {
	return s.PostgresRepository.UpdateAdditionalInfo(payload)
}

func (s *CustomerService) FindByID(id int) (*domain.Customer, error) {
	r, err := s.Query(&models.CustomersQueryRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("customer not found")
	}
	customers := r.Items.([]*domain.Customer)
	if len(customers) == 0 {
		return nil, errors.New("customer not found")
	}
	return customers[0], nil
}

func (s *CustomerService) CreateCreditDetails(payload *models.CustomersCreateCreditDetailsRequestBody) (*domain.CustomerCreditDetail, error) {
	return s.PostgresRepository.CreateCreditDetails(payload)
}

func (s *CustomerService) QueryCreditDetails(q *models.CustomersQueryCreditDetailsRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer credit details", q)
	creditDetails, err := s.PostgresRepository.QueryCreditDetails(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountCreditDetails(&models.CustomersQueryCreditDetailsRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      creditDetails,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) DeleteCustomerCreditDetails(payload *models.CustomersDeleteCreditDetailsRequestBody) ([]int64, error) {
	return s.PostgresRepository.DeleteCustomerCreditDetails(payload)
}

func (s *CustomerService) FindCreditDetailsByID(id int) (*domain.CustomerCreditDetail, error) {
	r, err := s.QueryCreditDetails(&models.CustomersQueryCreditDetailsRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("credit details not found")
	}
	creditDetails := r.Items.([]*domain.CustomerCreditDetail)
	if len(creditDetails) == 0 {
		return nil, errors.New("credit details not found")
	}
	return creditDetails[0], nil
}

func (s *CustomerService) FindCreditDetailsByIDAndCustomerID(id int, customerId int) (*domain.CustomerCreditDetail, error) {
	r, err := s.QueryCreditDetails(&models.CustomersQueryCreditDetailsRequestParams{
		ID:         id,
		CustomerID: customerId,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("credit details not found")
	}
	creditDetails := r.Items.([]*domain.CustomerCreditDetail)
	if len(creditDetails) == 0 {
		return nil, errors.New("credit details not found")
	}
	return creditDetails[0], nil
}

func (s *CustomerService) UpdateCustomerCreditDetails(payload *models.CustomersUpdateCreditDetailsRequestBody) (*domain.CustomerCreditDetail, error) {
	return s.PostgresRepository.UpdateCustomerCreditDetails(payload)
}

func (s *CustomerService) CreateAbsences(customer *domain.Customer, payload *models.CustomersCreateAbsencesRequestBody) (*domain.CustomerAbsence, error) {
	return s.PostgresRepository.CreateAbsences(customer, payload)
}

func (s *CustomerService) QueryAbsences(q *models.CustomersQueryAbsencesRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer absences", q)
	absences, err := s.PostgresRepository.QueryAbsences(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountAbsences(&models.CustomersQueryAbsencesRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
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

func (s *CustomerService) FindCustomerAbsenceByID(id int) (*domain.CustomerAbsence, error) {
	r, err := s.QueryAbsences(&models.CustomersQueryAbsencesRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("absence not found")
	}
	absences, ok := r.Items.([]*domain.CustomerAbsence)
	if !ok {
		return nil, errors.New("absence not found")
	}
	if len(absences) == 0 {
		return nil, errors.New("absence not found")
	}
	return absences[0], nil
}

func (s *CustomerService) UpdateAbsence(customerAbsence *domain.CustomerAbsence, payload *models.CustomersUpdateAbsenceRequestBody) (*domain.CustomerAbsence, error) {
	return s.PostgresRepository.UpdateAbsence(customerAbsence, payload)
}

func (s *CustomerService) DeleteAbsences(payload *models.CustomersDeleteAbsencesRequestBody) ([]int64, error) {
	return s.PostgresRepository.DeleteAbsences(payload)
}

func (s *CustomerService) QueryServices(queries *models.CustomersQueryServicesRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer services", queries)
	services, err := s.PostgresRepository.QueryServices(queries)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountServices(&models.CustomersQueryServicesRequestParams{
		ID:      queries.ID,
		Page:    queries.Page,
		Limit:   0,
		Filters: queries.Filters,
	})
	if err != nil {
		return nil, err
	}

	queries.Page = exp.TerIf(queries.Page < 1, 1, queries.Page)
	queries.Limit = exp.TerIf(queries.Limit < 10, 1, queries.Limit)

	page := queries.Page
	limit := queries.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      services,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) CreateServices(customer *domain.Customer, payload *models.CustomersCreateServicesRequestBody) (*domain.CustomerServices, error) {
	return s.PostgresRepository.CreateServices(customer, payload)
}

func (s *CustomerService) FindCustomerServiceByID(id int) (*domain.CustomerServices, error) {
	r, err := s.QueryServices(&models.CustomersQueryServicesRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("service not found")
	}
	services, ok := r.Items.([]*domain.CustomerServices)
	if !ok {
		return nil, errors.New("service not found")
	}
	if len(services) == 0 {
		return nil, errors.New("service not found")
	}
	return services[0], nil
}

func (s *CustomerService) FindCustomerServiceByIDAndCustomerID(id int, customerId int) (*domain.CustomerServices, error) {
	r, err := s.QueryServices(&models.CustomersQueryServicesRequestParams{
		ID:         id,
		CustomerID: customerId,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("service not found")
	}
	services, ok := r.Items.([]*domain.CustomerServices)
	if !ok {
		return nil, errors.New("service not found")
	}
	if len(services) == 0 {
		return nil, errors.New("service not found")
	}
	return services[0], nil
}

func (s *CustomerService) UpdateService(customerService *domain.CustomerServices, payload *models.CustomersCreateServicesRequestBody) (*domain.CustomerServices, error) {
	return s.PostgresRepository.UpdateService(customerService, payload)
}

func (s *CustomerService) DeleteServices(payload *models.CustomersDeleteServicesRequestBody) ([]int64, error) {
	return s.PostgresRepository.DeleteServices(payload)
}

func (s *CustomerService) QueryMedicines(queries *models.CustomersQueryMedicinesRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer medicines", queries)
	medicines, err := s.PostgresRepository.QueryMedicines(queries)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountMedicines(&models.CustomersQueryMedicinesRequestParams{
		ID:      queries.ID,
		Page:    queries.Page,
		Limit:   0,
		Filters: queries.Filters,
	})
	if err != nil {
		return nil, err
	}

	queries.Page = exp.TerIf(queries.Page < 1, 1, queries.Page)
	queries.Limit = exp.TerIf(queries.Limit < 10, 1, queries.Limit)

	page := queries.Page
	limit := queries.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      medicines,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) CreateMedicines(customer *domain.Customer, payload *models.CustomersCreateMedicinesRequestBody) (*domain.CustomerMedicine, error) {
	return s.PostgresRepository.CreateMedicines(customer, payload)
}

func (s *CustomerService) FindCustomerMedicineByIDAndCustomerID(id int, customerId int) (*domain.CustomerMedicine, error) {
	r, err := s.QueryMedicines(&models.CustomersQueryMedicinesRequestParams{
		ID:         id,
		CustomerID: customerId,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("service not found")
	}
	medicines, ok := r.Items.([]*domain.CustomerMedicine)
	if !ok {
		return nil, errors.New("service not found")
	}
	if len(medicines) == 0 {
		return nil, errors.New("service not found")
	}
	return medicines[0], nil
}

func (s *CustomerService) UpdateMedicine(customerMedicine *domain.CustomerMedicine, payload *models.CustomersUpdateMedicinesRequestBody) (*domain.CustomerMedicine, error) {
	return s.PostgresRepository.UpdateMedicine(customerMedicine, payload)
}

func (s *CustomerService) DeleteMedicines(payload *models.CustomersDeleteMedicinesRequestBody) ([]int64, error) {
	return s.PostgresRepository.DeleteMedicines(payload)
}

func (s *CustomerService) UpdatePersonalInfo(customerId int64, payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	return s.PostgresRepository.UpdatePersonalInfo(customerId, payload)
}

func (s *CustomerService) CreateOtherAttachments(customer *domain.Customer, payload *models.CustomersCreateOtherAttachmentsRequestBody) (*domain.CustomerOtherAttachment, error) {
	return s.PostgresRepository.CreateOtherAttachments(customer, payload)
}

func (s *CustomerService) UpdateCustomerOtherAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CustomerOtherAttachment, error) {
	return s.PostgresRepository.UpdateCustomerOtherAttachments(attachments, id)
}

func (s *CustomerService) QueryOtherAttachments(dataModel *models.CustomersQueryOtherAttachmentsRequestParams) (*restypes.QueryResponse, error) {
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

func (s *CustomerService) FindCustomerOtherAttachmentByID(id int) (*domain.CustomerOtherAttachment, error) {
	r, err := s.QueryOtherAttachments(&models.CustomersQueryOtherAttachmentsRequestParams{
		ID:    id,
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("customer other attachment not found")
	}
	customerOtherAttachments, ok := r.Items.([]*domain.CustomerOtherAttachment)
	if !ok {
		return nil, errors.New("customer other attachment not found")
	}
	return customerOtherAttachments[0], nil
}

func (s *CustomerService) UpdateCustomerOtherAttachment(customerOtherAttachment *domain.CustomerOtherAttachment, payload *models.CustomersUpdateOtherAttachmentRequestBody) (*domain.CustomerOtherAttachment, error) {
	return s.PostgresRepository.UpdateCustomerOtherAttachment(customerOtherAttachment, payload)
}

func (s *CustomerService) DeleteCustomerOtherAttachments(payload *models.CustomersDeleteOtherAttachmentsRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.DeleteCustomerOtherAttachments(payload)
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

func (s *CustomerService) CreateRelatives(customer *domain.Customer, payload *models.CustomersCreateRelativesRequestBody) (*domain.CustomerRelative, error) {
	return s.PostgresRepository.CreateRelatives(customer, payload)
}

func (s *CustomerService) QueryRelatives(q *models.CustomersQueryRelativesRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer relatives", q)
	relatives, err := s.PostgresRepository.QueryRelatives(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountRelatives(&models.CustomersQueryRelativesRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      relatives,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) FindCustomerRelativeByID(id int) (*domain.CustomerRelative, error) {
	r, err := s.QueryRelatives(&models.CustomersQueryRelativesRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("relative not found")
	}
	relatives, ok := r.Items.([]*domain.CustomerRelative)
	if !ok {
		return nil, errors.New("relative not found")
	}
	if len(relatives) == 0 {
		return nil, errors.New("relative not found")
	}
	return relatives[0], nil
}

func (s *CustomerService) UpdateRelative(customerRelative *domain.CustomerRelative, payload *models.CustomersCreateRelativesRequestBody) (*domain.CustomerRelative, error) {
	return s.PostgresRepository.UpdateRelative(customerRelative, payload)
}

func (s *CustomerService) DeleteRelatives(payload *models.CustomersDeleteRelativesRequestBody) ([]int64, error) {
	return s.PostgresRepository.DeleteRelatives(payload)
}

func (s *CustomerService) QueryContractualMobilityRestrictionLogs(q *models.CustomersQueryContractualMobilityRestrictionLogsRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer contractualMobilityRestrictionLogs", q)
	contractualMobilityRestrictionLogs, err := s.PostgresRepository.QueryContractualMobilityRestrictionLogs(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountContractualMobilityRestrictionLogs(&models.CustomersQueryContractualMobilityRestrictionLogsRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      contractualMobilityRestrictionLogs,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) FindCustomerContractualMobilityRestrictionLogByID(id int) (*domain.CustomerContractualMobilityRestrictionLog, error) {
	r, err := s.QueryContractualMobilityRestrictionLogs(&models.CustomersQueryContractualMobilityRestrictionLogsRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("contractualMobilityRestrictionLog not found")
	}
	contractualMobilityRestrictionLogs, ok := r.Items.([]*domain.CustomerContractualMobilityRestrictionLog)
	if !ok {
		return nil, errors.New("contractualMobilityRestrictionLog not found")
	}
	if len(contractualMobilityRestrictionLogs) == 0 {
		return nil, errors.New("contractualMobilityRestrictionLog not found")
	}
	return contractualMobilityRestrictionLogs[0], nil
}

func (s *CustomerService) UpdateAbsenceAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CustomerAbsence, error) {
	return s.PostgresRepository.UpdateAbsenceAttachments(previousAttachments, attachments, id)
}

func (s *CustomerService) UpdateMedicineAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CustomerMedicine, error) {
	return s.PostgresRepository.UpdateMedicineAttachments(previousAttachments, attachments, id)
}

func (s *CustomerService) QueryStatusLogs(q *models.CustomersQueryStatusLogsRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying customer statusLogs", q)
	statusLogs, err := s.PostgresRepository.QueryStatusLogs(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountStatusLogs(&models.CustomersQueryStatusLogsRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      statusLogs,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *CustomerService) FindCustomerStatusLogByID(id int) (*domain.CustomerStatusLog, error) {
	r, err := s.QueryStatusLogs(&models.CustomersQueryStatusLogsRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("statusLog not found")
	}
	statusLogs, ok := r.Items.([]*domain.CustomerStatusLog)
	if !ok {
		return nil, errors.New("statusLog not found")
	}
	if len(statusLogs) == 0 {
		return nil, errors.New("statusLog not found")
	}
	return statusLogs[0], nil
}
