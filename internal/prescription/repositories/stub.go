package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/prescription/domain"
	"github.com/hoitek/Maja-Service/internal/prescription/models"
)

type PrescriptionRepositoryStub struct {
	Prescriptions []*domain.Prescription
}

type prescriptionTestCondition struct {
	HasError bool
}

var UserTestCondition *prescriptionTestCondition = &prescriptionTestCondition{}

func NewPrescriptionRepositoryStub() *PrescriptionRepositoryStub {
	return &PrescriptionRepositoryStub{
		Prescriptions: []*domain.Prescription{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *PrescriptionRepositoryStub) Query(dataModel *models.PrescriptionsQueryRequestParams) ([]*domain.Prescription, error) {
	var prescriptions []*domain.Prescription
	for _, v := range r.Prescriptions {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			prescriptions = append(prescriptions, v)
			break
		}
	}
	return prescriptions, nil
}

func (r *PrescriptionRepositoryStub) Count(dataModel *models.PrescriptionsQueryRequestParams) (int64, error) {
	var prescriptions []*domain.Prescription
	for _, v := range r.Prescriptions {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			prescriptions = append(prescriptions, v)
			break
		}
	}
	return int64(len(prescriptions)), nil
}

func (r *PrescriptionRepositoryStub) Migrate() {
	// do stuff
}

func (r *PrescriptionRepositoryStub) Seed() {
	// do stuff
}

func (r *PrescriptionRepositoryStub) Create(payload *models.PrescriptionsCreateRequestBody) (*domain.Prescription, error) {
	panic("implement me")
}

func (r *PrescriptionRepositoryStub) Delete(payload *models.PrescriptionsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *PrescriptionRepositoryStub) Update(payload *models.PrescriptionsUpdateRequestBody, id int) (*domain.Prescription, error) {
	panic("implement me")
}

func (r *PrescriptionRepositoryStub) UpdatePrescriptionAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Prescription, error) {
	panic("implement me")
}
