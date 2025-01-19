package repositories

import (
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/customer/models"
)

type CustomerRepositoryStub struct {
	Customers []*domain.Customer
}

type customerTestCondition struct {
	HasError bool
}

var UserTestCondition *customerTestCondition = &customerTestCondition{}

func NewCustomerRepositoryStub() *CustomerRepositoryStub {
	return &CustomerRepositoryStub{
		Customers: []*domain.Customer{
			{
				ID: 1,
			},
		},
	}
}

func (r *CustomerRepositoryStub) FindCustomerServicesForSpecificShift(cyclePickupShiftID int64, date time.Time, shiftName string, shiftMorningStartHour int64, shiftMorningEndHour int64, shiftEveningStartHour int64, shiftEveningEndHour int64, shiftNightStartHour int64, shiftNightEndHour int64) ([]*domain.CustomerServices, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) Query(dataModel *models.CustomersQueryRequestParams) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	for _, v := range r.Customers {
		if v.ID == int64(dataModel.ID) {
			customers = append(customers, v)
			break
		}
	}
	return customers, nil
}

func (r *CustomerRepositoryStub) Count(dataModel *models.CustomersQueryRequestParams) (int64, error) {
	var customers []*domain.Customer
	for _, v := range r.Customers {
		if v.ID == int64(dataModel.ID) {
			customers = append(customers, v)
			break
		}
	}
	return int64(len(customers)), nil
}

func (r *CustomerRepositoryStub) CreatePersonalInfo(payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateUserInformation(customerID int64, payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateAdditionalInfo(payload *models.CustomersUpdateAdditionalInfoRequestBody) (*domain.Customer, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CreateCreditDetails(payload *models.CustomersCreateCreditDetailsRequestBody) (*domain.CustomerCreditDetail, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryCreditDetails(queries *models.CustomersQueryCreditDetailsRequestParams) ([]*domain.CustomerCreditDetail, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CountCreditDetails(dataModel *models.CustomersQueryCreditDetailsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) DeleteCustomerCreditDetails(payload *models.CustomersDeleteCreditDetailsRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateCustomerCreditDetails(payload *models.CustomersUpdateCreditDetailsRequestBody) (*domain.CustomerCreditDetail, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CreateAbsences(customer *domain.Customer, payload *models.CustomersCreateAbsencesRequestBody) (*domain.CustomerAbsence, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryAbsences(q *models.CustomersQueryAbsencesRequestParams) ([]*domain.CustomerAbsence, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CountAbsences(q *models.CustomersQueryAbsencesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateAbsence(customerAbsence *domain.CustomerAbsence, payload *models.CustomersUpdateAbsenceRequestBody) (*domain.CustomerAbsence, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) DeleteAbsences(payload *models.CustomersDeleteAbsencesRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CreateServices(customer *domain.Customer, payload *models.CustomersCreateServicesRequestBody) (*domain.CustomerServices, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryServices(queries *models.CustomersQueryServicesRequestParams) ([]*domain.CustomerServices, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CountServices(queries *models.CustomersQueryServicesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateService(customerService *domain.CustomerServices, payload *models.CustomersCreateServicesRequestBody) (*domain.CustomerServices, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) DeleteServices(payload *models.CustomersDeleteServicesRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryMedicines(queries *models.CustomersQueryMedicinesRequestParams) ([]*domain.CustomerMedicine, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CountMedicines(queries *models.CustomersQueryMedicinesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CreateMedicines(customer *domain.Customer, payload *models.CustomersCreateMedicinesRequestBody) (*domain.CustomerMedicine, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateMedicine(customerMedicine *domain.CustomerMedicine, payload *models.CustomersUpdateMedicinesRequestBody) (*domain.CustomerMedicine, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) DeleteMedicines(payload *models.CustomersDeleteMedicinesRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdatePersonalInfo(customerId int64, payload *models.CustomersCreatePersonalInfoRequestBody) (*domain.Customer, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CreateOtherAttachments(customer *domain.Customer, payload *models.CustomersCreateOtherAttachmentsRequestBody) (*domain.CustomerOtherAttachment, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateCustomerOtherAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CustomerOtherAttachment, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryOtherAttachments(dataModel *models.CustomersQueryOtherAttachmentsRequestParams) ([]*domain.CustomerOtherAttachment, int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateCustomerOtherAttachment(customerOtherAttachment *domain.CustomerOtherAttachment, payload *models.CustomersUpdateOtherAttachmentRequestBody) (*domain.CustomerOtherAttachment, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) DeleteCustomerOtherAttachments(payload *models.CustomersDeleteOtherAttachmentsRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CreateRelatives(customer *domain.Customer, payload *models.CustomersCreateRelativesRequestBody) (*domain.CustomerRelative, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryRelatives(queries *models.CustomersQueryRelativesRequestParams) ([]*domain.CustomerRelative, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CountRelatives(queries *models.CustomersQueryRelativesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateRelative(customerRelative *domain.CustomerRelative, payload *models.CustomersCreateRelativesRequestBody) (*domain.CustomerRelative, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) DeleteRelatives(payload *models.CustomersDeleteRelativesRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryContractualMobilityRestrictionLogs(queries *models.CustomersQueryContractualMobilityRestrictionLogsRequestParams) ([]*domain.CustomerContractualMobilityRestrictionLog, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) CountContractualMobilityRestrictionLogs(queries *models.CustomersQueryContractualMobilityRestrictionLogsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateAbsenceAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CustomerAbsence, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) UpdateMedicineAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CustomerMedicine, error) {
	panic("implement me")
}

func (r *CustomerRepositoryStub) QueryStatusLogs(queries *models.CustomersQueryStatusLogsRequestParams) ([]*domain.CustomerStatusLog, error) {
	panic("implement me")
}
func (r *CustomerRepositoryStub) CountStatusLogs(queries *models.CustomersQueryStatusLogsRequestParams) (int64, error) {
	panic("implement me")
}
