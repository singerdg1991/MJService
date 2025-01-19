package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"

	"github.com/hoitek/Maja-Service/internal/archive/domain"
	"github.com/hoitek/Maja-Service/internal/archive/models"
)

type ArchiveRepositoryStub struct {
	Archives []*domain.Archive
}

type archiveTestCondition struct {
	HasError bool
}

var UserTestCondition *archiveTestCondition = &archiveTestCondition{}

func NewArchiveRepositoryStub() *ArchiveRepositoryStub {
	return &ArchiveRepositoryStub{
		Archives: []*domain.Archive{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *ArchiveRepositoryStub) Query(dataModel *models.ArchivesQueryRequestParams) ([]*domain.Archive, error) {
	var archives []*domain.Archive
	for _, v := range r.Archives {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			archives = append(archives, v)
			break
		}
	}
	return archives, nil
}

func (r *ArchiveRepositoryStub) Count(dataModel *models.ArchivesQueryRequestParams) (int64, error) {
	var archives []*domain.Archive
	for _, v := range r.Archives {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			archives = append(archives, v)
			break
		}
	}
	return int64(len(archives)), nil
}

func (r *ArchiveRepositoryStub) Migrate() {
	// do stuff
}

func (r *ArchiveRepositoryStub) Seed() {
	// do stuff
}

func (r *ArchiveRepositoryStub) Create(payload *models.ArchivesCreateRequestBody) (*domain.Archive, error) {
	panic("implement me")
}

func (r *ArchiveRepositoryStub) Delete(payload *models.ArchivesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *ArchiveRepositoryStub) Update(payload *models.ArchivesCreateRequestBody, id int64) (*domain.Archive, error) {
	panic("implement me")
}

func (r *ArchiveRepositoryStub) GetArchivesByIds(ids []int64) ([]*domain.Archive, error) {
	panic("implement me")
}

func (r *ArchiveRepositoryStub) UpdateAttachments(attachments []*types.UploadMetadata, id int64) (*domain.Archive, error) {
	panic("implement me")
}
