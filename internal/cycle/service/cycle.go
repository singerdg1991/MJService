package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/models"
	"github.com/hoitek/Maja-Service/internal/cycle/ports"
	"github.com/hoitek/Maja-Service/storage"
)

type CycleService struct {
	PostgresRepository ports.CycleRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

// NewCycleService creates a new instance of CycleService.
//
// pDB is the Postgres database repository for cycles, and m is the MinIO storage.
// Returns a new CycleService instance.
func NewCycleService(pDB ports.CycleRepositoryPostgresDB, m *storage.MinIO) CycleService {
	go minio.SetupMinIOStorage(constants.CYCLE_BUCKET_NAME, m)
	return CycleService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

// Query queries cycles based on the provided query parameters.
//
// q is the query request parameters, including page, limit, and other filters.
// Returns a QueryResponse containing the queried cycles and pagination information, or an error if the query fails.
func (s *CycleService) Query(q *models.CyclesQueryRequestParams) (*restypes.QueryResponse, error) {
	cycles, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}
	for _, cycle := range cycles {
		cycle.SetDefaultStatus()
	}

	count, err := s.PostgresRepository.Count(q)
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

	// Transform the response to the format that the frontend expects
	var items []*domain.Cycle
	for _, item := range cycles {
		items = append(items, &domain.Cycle{
			ID:                    item.ID,
			SectionID:             item.SectionID,
			Name:                  item.Name,
			StartDate:             item.StartDate,
			EndDate:               item.EndDate,
			PeriodLength:          item.PeriodLength,
			ShiftMorningStartTime: item.ShiftMorningStartTime,
			ShiftMorningEndTime:   item.ShiftMorningEndTime,
			ShiftEveningStartTime: item.ShiftEveningStartTime,
			ShiftEveningEndTime:   item.ShiftEveningEndTime,
			ShiftNightStartTime:   item.ShiftNightStartTime,
			ShiftNightEndTime:     item.ShiftNightEndTime,
			FreezePeriodDate:      item.FreezePeriodDate,
			WishDays:              item.WishDays,
			StaffTypes:            item.StaffTypes,
			NextStaffTypes:        item.NextStaffTypes,
			Status:                item.Status,
		})
	}

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// FindAllStaffTypesByCycleID retrieves all staff types associated with the specified cycle ID.
//
// cycleID is the ID of the cycle for which to retrieve staff types.
// []*domain.CycleStaffType, error
func (s *CycleService) FindAllStaffTypesByCycleID(cycleID int64) ([]*domain.CycleStaffType, error) {
	return s.PostgresRepository.FindAllStaffTypesByCycleID(cycleID)
}

// QueryStaffTypes retrieves staff types based on the provided query parameters.
//
// q is a pointer to the models.CyclesQueryStaffTypesRequestParams struct, which contains the query parameters.
// *restypes.QueryResponse, error
func (s *CycleService) QueryStaffTypes(q *models.CyclesQueryStaffTypesRequestParams) (*restypes.QueryResponse, error) {
	staffTypes, err := s.PostgresRepository.QueryStaffTypes(q)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountStaffTypes(q)
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

	// Transform the response to the format that the frontend expects
	var items []*domain.CycleStaffType
	for _, item := range staffTypes {
		items = append(items, &domain.CycleStaffType{
			ID:               item.ID,
			CycleID:          item.CycleID,
			RoleID:           item.RoleID,
			Role:             item.Role,
			DateTime:         item.DateTime,
			ShiftName:        item.ShiftName,
			NeededStaffCount: item.NeededStaffCount,
			StartHour:        item.StartHour,
			EndHour:          item.EndHour,
			UsedStaffCount:   item.UsedStaffCount,
			RemindStaffCount: item.RemindStaffCount,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
			DeletedAt:        item.DeletedAt,
		})
	}

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new cycle.
//
// It takes a CyclesCreateRequestBody object as a parameter.
// Returns a Cycle object and an error.
func (s *CycleService) Create(payload *models.CyclesCreateRequestBody) (*domain.Cycle, error) {
	cycle, err := s.PostgresRepository.Create(payload)
	if err != nil {
		return nil, err
	}

	// Set default status
	cycle.SetDefaultStatus()

	return cycle, nil
}

// Delete deletes a cycle by its ID.
//
// It takes a CyclesDeleteRequestBody object as a parameter.
// Returns a DeleteResponse object and an error.
func (s *CycleService) Delete(payload *models.CyclesDeleteRequestBody) (*restypes.DeleteResponse, error) {
	ids, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}
	var resultIDs []uint
	for _, id := range ids {
		resultIDs = append(resultIDs, uint(id))
	}

	return &restypes.DeleteResponse{
		IDs: resultIDs,
	}, nil
}

// Update updates a cycle.
//
// It takes a CyclesCreateRequestBody object and an ID as parameters.
// Returns a Cycle object and an error.
func (s *CycleService) Update(payload *models.CyclesCreateRequestBody, id int64) (*domain.Cycle, error) {
	// Update the cycle
	cycle, err := s.PostgresRepository.Update(payload, id)
	if err != nil {
		return nil, err
	}

	// Set default status
	cycle.SetDefaultStatus()

	return cycle, nil
}

// UpdateStaffType updates the staff type for a cycle in the CycleService.
//
// It takes a payload of type *models.CyclesUpdateStaffTypeRequestBody,
// an ID of type int64, and a boolean isUnplanned as parameters.
// It returns a Cycle object and an error.
func (s *CycleService) UpdateStaffType(payload *models.CyclesUpdateStaffTypeRequestBody, id int64, isUnplanned bool) (*domain.Cycle, error) {
	cycle, err := s.PostgresRepository.UpdateStaffType(payload, id, isUnplanned)
	if err != nil {
		return nil, err
	}
	cycle.SetDefaultStatus()
	return cycle, nil
}

// UpdateStaffTypes updates multiple staff types for a cycle in the CycleService.
//
// It takes a payload of type *models.CyclesUpdateStaffTypesRequestBody and an ID of type int64 as parameters.
// Returns a Cycle object and an error.
func (s *CycleService) UpdateStaffTypes(payload *models.CyclesUpdateStaffTypesRequestBody, id int64) (*domain.Cycle, error) {
	cycle, err := s.PostgresRepository.UpdateStaffTypes(payload, id)
	if err != nil {
		return nil, err
	}
	cycle.SetDefaultStatus()
	return cycle, nil
}

// UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle updates staff type for a cycle and also assigns incoming shifts to a staff member.
//
// It takes a payload of type *models.CyclesUpdateStaffTypesRequestBody and an ID of type int64 as parameters.
// Returns a Cycle object and an error.
func (s *CycleService) UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(payload *models.CyclesUpdateStaffTypeRequestBody, migratedCycleID int64, currentCycleID int64) (*domain.Cycle, error) {
	cycle, err := s.PostgresRepository.UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(payload, migratedCycleID, currentCycleID)
	if err != nil {
		return nil, err
	}
	cycle.SetDefaultStatus()
	return cycle, nil
}

// FindByID finds a cycle by its ID.
//
// It takes an id of type int64 as a parameter.
// Returns a *domain.Cycle object and an error.
func (s *CycleService) FindByID(id int64) (*domain.Cycle, error) {
	r, err := s.Query(&models.CyclesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("cycle not found")
	}
	Cycles := r.Items.([]*domain.Cycle)
	if len(Cycles) == 0 {
		return nil, errors.New("cycle not found")
	}
	return Cycles[0], nil
}

// GetLast retrieves the most recent cycle.
//
// It takes no parameters.
// Returns a *domain.Cycle object and an error.
func (s *CycleService) GetLast() (*domain.Cycle, error) {
	r, err := s.Query(&models.CyclesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Sorts: models.CycleSortType{
			ID: models.CycleSortValue{
				Op: "desc",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("there is no cycle yet")
	}
	Cycles := r.Items.([]*domain.Cycle)
	if len(Cycles) == 0 {
		return nil, errors.New("there is no cycle yet")
	}
	return Cycles[0], nil
}

// GetCurrent retrieves the most recent cycle.
//
// It takes no parameters.
// Returns a *domain.Cycle object and an error.
func (s *CycleService) GetCurrent() (*domain.Cycle, error) {
	r, err := s.Query(&models.CyclesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Sorts: models.CycleSortType{
			ID: models.CycleSortValue{
				Op: "desc",
			},
		},
		Filters: models.CycleFilterType{
			StartDate: filters.FilterValue[string]{
				Op:    "<=",
				Value: fmt.Sprintf("%sT23:59:59Z", time.Now().Format("2006-01-02")),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("there is no cycle yet")
	}
	Cycles := r.Items.([]*domain.Cycle)
	if len(Cycles) == 0 {
		return nil, errors.New("there is no cycle yet")
	}
	return Cycles[0], nil
}

// Duplicate duplicates a cycle.
//
// It takes a models.CyclesDuplicateRequestBody object as a parameter.
// Returns a *domain.Cycle object and an error.
func (s *CycleService) Duplicate(payload *models.CyclesDuplicateRequestBody) (*domain.Cycle, error) {
	cycle, err := s.PostgresRepository.Duplicate(payload)
	if err != nil {
		return nil, err
	}

	// Set default status
	cycle.SetDefaultStatus()
	return cycle, nil
}

// QueryNextStaffTypes retrieves next staff types based on the provided query parameters.
//
// It takes a pointer to the models.CyclesQueryNextStaffTypesRequestParams struct as a parameter, which contains the query parameters.
// Returns a *restypes.QueryResponse object and an error.
func (s *CycleService) QueryNextStaffTypes(dataModel *models.CyclesQueryNextStaffTypesRequestParams) (*restypes.QueryResponse, error) {
	staffTypes, err := s.PostgresRepository.QueryNextStaffTypes(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountNextStaffTypes(dataModel)
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

	// Transform the response to the format that the frontend expects
	var items []*domain.CycleNextStaffType
	for _, item := range staffTypes {
		items = append(items, &domain.CycleNextStaffType{
			ID:               item.ID,
			CurrentCycleID:   item.CurrentCycleID,
			RoleID:           item.RoleID,
			Role:             item.Role,
			DateTime:         item.DateTime,
			ShiftName:        item.ShiftName,
			NeededStaffCount: item.NeededStaffCount,
			StartHour:        item.StartHour,
			EndHour:          item.EndHour,
			UsedStaffCount:   item.UsedStaffCount,
			RemindStaffCount: item.RemindStaffCount,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
			DeletedAt:        item.DeletedAt,
		})
	}

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// UpdateNextStaffType updates the next staff type for a given cycle.
//
// payload - The request body containing the updated staff type details.
// id - The ID of the cycle to update.
//
// []*domain.CycleNextStaffType, error
func (s *CycleService) UpdateNextStaffType(payload *models.CyclesUpdateNextStaffTypeRequestBody, id int64) ([]*domain.CycleNextStaffType, error) {
	return s.PostgresRepository.UpdateNextStaffType(payload, id)
}

// UpdateNextStaffTypes updates the next staff types for a cycle in the CycleService.
//
// It takes a payload of type *models.CyclesUpdateNextStaffTypesRequestBody and an ID of type int64 as parameters.
// It returns a slice of pointers to *domain.CycleNextStaffType and an error.
func (s *CycleService) UpdateNextStaffTypes(payload *models.CyclesUpdateNextStaffTypesRequestBody, id int64) ([]*domain.CycleNextStaffType, error) {
	return s.PostgresRepository.UpdateNextStaffTypes(payload, id)
}

// PickupShift creates a new pickup shift for a given cycle.
//
// payload - The request body containing the pickup shift details.
//
// *restypes.QueryResponse, []*domain.CyclePickupShift, error
func (s *CycleService) PickupShift(payload *models.CyclesCreatePickupShiftRequestBody) (*restypes.QueryResponse, []*domain.CyclePickupShift, error) {
	pickupShifts, err := s.PostgresRepository.PickupShift(payload)
	if err != nil {
		return nil, nil, err
	}
	dataModel := &models.CyclesQueryPickupShiftsRequestParams{
		CycleID: payload.CycleID,
		StaffID: payload.StaffID,
		Filters: models.CycleQueryPickupShiftsFilterType{
			DateTime: filters.FilterValue[string]{
				Op:    "equals",
				Value: payload.DateTimeAsDate.Format("2006-01-02T15:04:05"),
			},
		},
	}
	count, err := s.PostgresRepository.CountPickupShifts(dataModel)
	if err != nil {
		return nil, nil, err
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
		Items:      pickupShifts,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, pickupShifts, nil
}

// QueryPickupShifts queries the pickup shifts based on the provided data model.
//
// dataModel is a pointer to models.CyclesQueryPickupShiftsRequestParams.
// Returns a pointer to restypes.QueryResponse and an error.
func (s *CycleService) QueryPickupShifts(dataModel *models.CyclesQueryPickupShiftsRequestParams) (*restypes.QueryResponse, error) {
	pickupShifts, err := s.PostgresRepository.QueryPickupShifts(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountPickupShifts(dataModel)
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
		Items:      pickupShifts,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// CountPickupShifts returns the count of pickup shifts based on the provided query parameters.
//
// It takes a pointer to models.CyclesQueryPickupShiftsRequestParams as an argument.
// Returns the count as an int64 and an error.
func (s *CycleService) CountPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) (int64, error) {
	return s.PostgresRepository.CountPickupShifts(queries)
}

// FindPickedUpShiftByID retrieves a picked up shift by its ID.
//
// id - The ID of the picked up shift to retrieve.
// Returns a pointer to domain.CyclePickupShift and an error.
func (s *CycleService) FindPickedUpShiftByID(id int64) (*domain.CyclePickupShift, error) {
	// Create the query model
	dataModel := &models.CyclesQueryPickupShiftsRequestParams{
		ID: int(id),
	}

	// Query the pickedUp shift
	pickedUpShifts, err := s.PostgresRepository.QueryPickupShifts(dataModel)
	if err != nil {
		return nil, err
	}

	// Check if the pickedUp shift is found
	if len(pickedUpShifts) == 0 {
		return nil, errors.New("pickedUp shift not found")
	}

	return pickedUpShifts[0], nil
}

// FindPickedUpShiftForStaff finds a picked up shift for a staff member based on the provided cycle ID, staff ID, and staff type ID.
//
// cycleId - The ID of the cycle to search for the picked up shift.
// staffId - The ID of the staff member to search for the picked up shift.
// staffTypeId - The ID of the staff type to search for the picked up shift.
// Returns a pointer to domain.CyclePickupShift and an error.
func (s *CycleService) FindPickedUpShiftForStaff(cycleId int64, staffId int64, staffTypeId int64) (*domain.CyclePickupShift, error) {
	// Create the query model
	dataModel := &models.CyclesQueryPickupShiftsRequestParams{
		CycleID:                int(cycleId),
		StaffID:                int(staffId),
		CycleStaffTypeIDs:      fmt.Sprintf("[%d]", int(staffTypeId)),
		CycleStaffTypeIDsInt64: []int64{staffTypeId},
	}

	// Query the pickedUp shift
	pickedUpShifts, err := s.PostgresRepository.QueryPickupShifts(dataModel)
	if err != nil {
		return nil, err
	}

	// Check if the pickedup shift is found
	if len(pickedUpShifts) == 0 {
		return nil, errors.New("pickedUp shift not found")
	}

	return pickedUpShifts[0], nil
}

// PickupShiftIncomingCycle retrieves the incoming cycle pickup shifts for a given payload.
//
// payload - The request body containing the incoming cycle pickup shift details.
// Returns a pointer to restypes.QueryResponse and an error.
func (s *CycleService) PickupShiftIncomingCycle(payload *models.CyclesCreateIncomingCyclePickupShiftRequestBody) (*restypes.QueryResponse, error) {
	pickupShifts, err := s.PostgresRepository.PickupShiftIncomingCycle(payload)
	if err != nil {
		return nil, err
	}
	dataModel := &models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
		CycleID: payload.CycleID,
		StaffID: payload.StaffID,
		Filters: models.CycleQueryIncomingCyclePickupShiftsFilterType{
			DateTime: filters.FilterValue[string]{
				Op:    "equals",
				Value: payload.DateTimeAsDate.Format("2006-01-02T15:04:05"),
			},
		},
	}
	count, err := s.PostgresRepository.CountIncomingCyclePickupShifts(dataModel)
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
		Items:      pickupShifts,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// QueryIncomingCyclePickupShifts queries the incoming cycle pickup shifts based on the provided data model.
//
// dataModel is a pointer to models.CyclesQueryIncomingCyclePickupShiftsRequestParams.
// Returns a pointer to restypes.QueryResponse and an error.
func (s *CycleService) QueryIncomingCyclePickupShifts(dataModel *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (*restypes.QueryResponse, error) {
	pickupShifts, err := s.PostgresRepository.QueryIncomingCyclePickupShifts(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountIncomingCyclePickupShifts(dataModel)
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
		Items:      pickupShifts,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// CountIncomingCyclePickupShifts returns the count of incoming cycle pickup shifts based on the provided query parameters.
//
// queries is a pointer to models.CyclesQueryIncomingCyclePickupShiftsRequestParams.
// Returns the count as an int64 and an error.
func (s *CycleService) CountIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (int64, error) {
	return s.PostgresRepository.CountIncomingCyclePickupShifts(queries)
}

// FindIncomingCyclePickedUpShiftByID retrieves an incoming cycle picked up shift by its ID.
//
// id - The ID of the incoming cycle picked up shift to retrieve.
// Returns a pointer to domain.CycleIncomingCyclePickupShift and an error.
func (s *CycleService) FindIncomingCyclePickedUpShiftByID(id int64) (*domain.CycleIncomingCyclePickupShift, error) {
	// Create the query model
	dataModel := &models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
		ID: int(id),
	}

	// Query the pickedUp shift
	pickedUpShifts, err := s.PostgresRepository.QueryIncomingCyclePickupShifts(dataModel)
	if err != nil {
		return nil, err
	}

	// Check if the pickedUp shift is found
	if len(pickedUpShifts) == 0 {
		return nil, errors.New("pickedUp shift not found")
	}

	return pickedUpShifts[0], nil
}

// FindIncomingCyclePickedUpShiftForStaff retrieves an incoming cycle picked up shift for a specific staff member.
//
// cycleId - The ID of the cycle.
// staffId - The ID of the staff member.
// nextStaffTypeId - The ID of the next staff type.
// Returns a pointer to domain.CycleIncomingCyclePickupShift and an error.
func (s *CycleService) FindIncomingCyclePickedUpShiftForStaff(cycleId int64, staffId int64, nextStaffTypeId int64) (*domain.CycleIncomingCyclePickupShift, error) {
	// Create the query model
	dataModel := &models.CyclesQueryIncomingCyclePickupShiftsRequestParams{
		CycleID:                    int(cycleId),
		StaffID:                    int(staffId),
		CycleNextStaffTypeIDs:      fmt.Sprintf("[%d]", int(nextStaffTypeId)),
		CycleNextStaffTypeIDsInt64: []int64{nextStaffTypeId},
	}

	// Query the pickedUp shift
	pickedUpShifts, err := s.PostgresRepository.QueryIncomingCyclePickupShifts(dataModel)
	if err != nil {
		return nil, err
	}

	// Check if the pickedup shift is found
	if len(pickedUpShifts) == 0 {
		return nil, errors.New("pickedUp shift not found")
	}

	return pickedUpShifts[0], nil
}

// FindAllNextStaffTypesByCycleID retrieves all next staff types for a given cycle ID.
//
// cycleID is the ID of the cycle for which to retrieve next staff types.
// []*domain.CycleNextStaffType, error
func (s *CycleService) FindAllNextStaffTypesByCycleID(cycleID int64) ([]*domain.CycleNextStaffType, error) {
	return s.PostgresRepository.FindAllNextStaffTypesByCycleID(cycleID)
}

// FindVisitsForStaffInSpecificShift retrieves visits for staff in a specific shift.
//
// cycleID - The ID of the cycle.
// staffID - The ID of the staff member.
// datetime - The date and time of the shift.
// shiftName - The name of the shift.
// Returns a slice of pointers to domain.CyclePickupShift and an error.
func (s *CycleService) FindVisitsForStaffInSpecificShift(cycleID int64, staffID int64, datetime *time.Time, shiftName string) ([]*domain.CyclePickupShift, error) {
	return s.PostgresRepository.FindVisitsForStaffInSpecificShift(cycleID, staffID, datetime, shiftName)
}

// AssignShiftsToStaff assigns a shift to the staff member with the specified targetStaffID.
//
// payload - The CyclesCreateShiftAssignToMeRequestBody containing the shift details.
// targetStaffID - The ID of the staff member to whom the shift is being assigned.
// []*domain.CyclePickupShift, error
func (s *CycleService) AssignShiftsToStaff(payload *models.CyclesCreateShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CyclePickupShift, error) {
	return s.PostgresRepository.AssignShiftsToStaff(payload, targetStaffID)
}

// SwapShifts swaps two shifts.
//
// payload - The CyclesCreateShiftSwapRequestBody containing the shift swap details.
// []*domain.CyclePickupShift, []*domain.CyclePickupShift, error
func (s *CycleService) SwapShifts(payload *models.CyclesCreateShiftSwapRequestBody) ([]*domain.CyclePickupShift, []*domain.CyclePickupShift, error) {
	return s.PostgresRepository.SwapShifts(payload)
}

// CreateShiftIfNotExist creates a shift if it does not exist.
//
// cycleId - The ID of the cycle.
// visitIds - A list of visit IDs.
// shiftName - The name of the shift.
// datetime - The date and time of the shift.
// status - The status of the shift.
// Returns a pointer to the ID of the created shift and an error.
func (s *CycleService) CreateShiftIfNotExist(cycleId int64, visitIds []int64, shiftName string, datetime *time.Time, status string) (*int64, error) {
	return s.PostgresRepository.CreateShiftIfNotExist(cycleId, visitIds, shiftName, datetime, status)
}

// ShiftStart starts a new shift.
//
// payload - The CyclesCreateShiftStartRequestBody containing the shift start details.
// Returns a pointer to domain.CycleShift and an error.
func (s *CycleService) ShiftStart(payload *models.CyclesCreateShiftStartRequestBody) (*domain.CycleShift, error) {
	return s.PostgresRepository.ShiftStart(payload)
}

// ShiftEnd ends a shift.
//
// payload - The CyclesCreateShiftEndRequestBody containing the shift end details.
// Returns a pointer to domain.CycleShift and an error.
func (s *CycleService) ShiftEnd(payload *models.CyclesCreateShiftEndRequestBody) (*domain.CycleShift, error) {
	return s.PostgresRepository.ShiftEnd(payload)
}

// QueryCycleShifts queries the cycle shifts based on the provided data model.
//
// dataModel is a pointer to models.CyclesQueryCycleShiftsRequestParams.
// Returns a pointer to restypes.QueryResponse and an error.
func (s *CycleService) QueryCycleShifts(dataModel *models.CyclesQueryCycleShiftsRequestParams) (*restypes.QueryResponse, error) {
	pickupShifts, err := s.PostgresRepository.QueryCycleShifts(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountCycleShifts(dataModel)
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
		Items:      pickupShifts,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// QueryShiftCustomerHomeKeys queries the shift customer home keys based on the provided data model.
//
// dataModel is a pointer to models.CyclesQueryShiftCustomerHomeKeysRequestParams.
// Returns a pointer to restypes.QueryResponse and an error.
func (s *CycleService) QueryShiftCustomerHomeKeys(dataModel *models.CyclesQueryShiftCustomerHomeKeysRequestParams) (*restypes.QueryResponse, error) {
	pickupShifts, err := s.PostgresRepository.QueryShiftCustomerHomeKeys(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountShiftCustomerHomeKeys(dataModel)
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
		Items:      pickupShifts,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// FindShiftByID retrieves a CycleShift by its ID from the CycleService.
//
// Parameters:
// - id: an int64 representing the ID of the CycleShift to retrieve.
//
// Returns:
// - *domain.CycleShift: a pointer to the retrieved CycleShift, or nil if not found.
// - error: an error if the retrieval process encounters any issues.
func (s *CycleService) FindShiftByID(id int64) (*domain.CycleShift, error) {
	r, err := s.QueryCycleShifts(&models.CyclesQueryCycleShiftsRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("cycle shift not found")
	}
	cycleShifts, ok := r.Items.([]*domain.CycleShift)
	if !ok {
		return nil, errors.New("cycle shift not found")
	}
	if len(cycleShifts) == 0 {
		return nil, errors.New("cycle shift not found")
	}
	return cycleShifts[0], nil
}

// FindPickupShiftByID retrieves a CyclePickupShift by its ID from the CycleService.
//
// Parameters:
// - id: an int64 representing the ID of the CyclePickupShift to retrieve.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the retrieved CyclePickupShift, or nil if not found.
// - error: an error if the retrieval process encounters any issues.
func (s *CycleService) FindPickupShiftByID(id int64) (*domain.CyclePickupShift, error) {
	r, err := s.QueryPickupShifts(&models.CyclesQueryPickupShiftsRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("pickup shift not found")
	}
	pickupShifts, ok := r.Items.([]*domain.CyclePickupShift)
	if !ok {
		return nil, errors.New("pickup shift not found")
	}
	if len(pickupShifts) == 0 {
		return nil, errors.New("pickup shift not found")
	}
	return pickupShifts[0], nil
}

// CreateCycleShiftCustomerHomeKey creates a new cycle shift customer home key.
//
// Parameters:
// - payload: a pointer to models.CyclesShiftCustomerHomeKeyRequestBody containing the request body.
//
// Returns:
// - *domain.CycleShiftCustomerHomeKey: a pointer to the created cycle shift customer home key.
// - error: an error if the creation process encounters any issues.
func (s *CycleService) CreateCycleShiftCustomerHomeKey(payload *models.CyclesShiftCustomerHomeKeyRequestBody) (*domain.CycleShiftCustomerHomeKey, error) {
	return s.PostgresRepository.CreateCycleShiftCustomerHomeKey(payload)
}

// VisitStart initiates a visit for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitStartRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit start process encounters any issues.
func (s *CycleService) VisitStart(payload *models.CyclesCreateVisitStartRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitStart(payload)
}

// VisitEnd marks the end of a visit.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitEndRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the update process encounters any issues.
func (s *CycleService) VisitEnd(payload *models.CyclesCreateVisitEndRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitEnd(payload)
}

// VisitCancel cancels a visit for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitCancelRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit cancel process encounters any issues.
func (s *CycleService) VisitCancel(payload *models.CyclesCreateVisitCancelRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitCancel(payload)
}

// VisitDelay initiates a delay for a visit.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitDelayRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit delay process encounters any issues.
func (s *CycleService) VisitDelay(payload *models.CyclesCreateVisitDelayRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitDelay(payload)
}

// VisitPause pauses a visit for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitPauseRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit pause process encounters any issues.
func (s *CycleService) VisitPause(payload *models.CyclesCreateVisitPauseRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitPause(payload)
}

// VisitResume resumes a visit for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitResumeRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit resume process encounters any issues.
func (s *CycleService) VisitResume(payload *models.CyclesCreateVisitResumeRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitResume(payload)
}

// VisitReactive initiates a reactive visit for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitReactiveRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit reactive process encounters any issues.
func (s *CycleService) VisitReactive(payload *models.CyclesCreateVisitReactiveRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.VisitReactive(payload)
}

// AssignVisitToStaff assigns a visit to a staff member.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitAssignRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the updated cycle pickup shift.
// - error: an error if the visit assignment process encounters any issues.
func (s *CycleService) AssignVisitToStaff(payload *models.CyclesCreateVisitAssignRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.AssignVisitToStaff(payload)
}

// SwapVisits swaps two visits for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitSwapRequestBody containing the visit swap details.
// Returns:
// - *domain.CyclePickupShift: a pointer to the first updated cycle pickup shift.
// - *domain.CyclePickupShift: a pointer to the second updated cycle pickup shift.
// - error: an error if the visit swap process encounters any issues.
func (s *CycleService) SwapVisits(payload *models.CyclesCreateVisitSwapRequestBody) (*domain.CyclePickupShift, *domain.CyclePickupShift, error) {
	return s.PostgresRepository.SwapVisits(payload)
}

// CreateUnplannedVisit creates an unplanned visit for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitUnplannedRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShift: a pointer to the created cycle pickup shift.
// - error: an error if the visit creation process encounters any issues.
func (s *CycleService) CreateUnplannedVisit(payload *models.CyclesCreateVisitUnplannedRequestBody) (*domain.CyclePickupShift, error) {
	return s.PostgresRepository.CreateUnplannedVisit(payload)
}

// CreateVisitTodo creates a visit todo for a cycle pickup shift.
//
// Parameters:
// - payload: a pointer to models.CyclesCreateVisitTodoRequestBody containing the request body.
//
// Returns:
// - *domain.CyclePickupShiftTodo: a pointer to the created cycle pickup shift todo.
// - error: an error if the visit todo creation process encounters any issues.
func (s *CycleService) CreateVisitTodo(payload *models.CyclesCreateVisitTodoRequestBody) (*domain.CyclePickupShiftTodo, error) {
	return s.PostgresRepository.CreateVisitTodo(payload)
}

// QueryVisitTodos queries the visit todos based on the provided data model.
//
// Parameters:
// - dataModel: a pointer to models.CyclesQueryVisitsTodosRequestParams containing the query parameters.
// Returns:
// - *restypes.QueryResponse: a pointer to the query response.
// - error: an error if the query process encounters any issues.
func (s *CycleService) QueryVisitTodos(dataModel *models.CyclesQueryVisitsTodosRequestParams) (*restypes.QueryResponse, error) {
	visitTodos, err := s.PostgresRepository.QueryVisitTodos(dataModel)
	if err != nil {
		return nil, err
	}
	count, err := s.PostgresRepository.CountVisitTodos(dataModel)
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
		Items:      visitTodos,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

// FindVisitTodoByID retrieves a visit todo by its ID.
//
// Parameters:
// - id: the ID of the visit todo to be retrieved.
// Returns:
// - *domain.CyclePickupShiftTodo: the retrieved visit todo.
// - error: an error if the retrieval process encounters any issues.
func (s *CycleService) FindVisitTodoByID(id int64) (*domain.CyclePickupShiftTodo, error) {
	r, err := s.QueryVisitTodos(&models.CyclesQueryVisitsTodosRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("visit todo not found")
	}
	visitTodos, ok := r.Items.([]*domain.CyclePickupShiftTodo)
	if !ok {
		return nil, errors.New("visit todo not found")
	}
	if len(visitTodos) == 0 {
		return nil, errors.New("visit todo not found")
	}
	return visitTodos[0], nil
}

// UpdateVisitTodoStatus updates the status of a visit todo.
//
// Parameters:
// - payload: the request body containing the updated status.
// - id: the ID of the visit todo to be updated.
// Returns:
// - *domain.CyclePickupShiftTodo: the updated visit todo.
// - error: an error if the update process encounters any issues.
func (s *CycleService) UpdateVisitTodoStatus(payload *models.CyclesUpdateVisitTodoStatusRequestBody, id int64) (*domain.CyclePickupShiftTodo, error) {
	return s.PostgresRepository.UpdateVisitTodoStatus(payload, id)
}

// UpdateVisitTodoAttachments updates the attachments of a visit todo.
//
// Parameters:
// - attachments: the list of attachments to be updated.
// - id: the ID of the visit todo to be updated.
// Returns:
// - *domain.CyclePickupShiftTodo: the updated visit todo.
// - error: an error if the update process encounters any issues.
func (s *CycleService) UpdateVisitTodoAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CyclePickupShiftTodo, error) {
	return s.PostgresRepository.UpdateVisitTodoAttachments(attachments, id)
}

// AssignIncomingShiftsToStaff assigns incoming shifts to staff members.
//
// Parameters:
// - payload: the request body containing the shift assignment details.
// - targetStaffID: the ID of the staff member to whom the shifts are being assigned.
// Returns:
// - []*domain.CycleIncomingCyclePickupShift: the assigned shifts.
// - error: an error if the assignment process encounters any issues.
func (s *CycleService) AssignIncomingShiftsToStaff(payload *models.CyclesCreateIncomingShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CycleIncomingCyclePickupShift, error) {
	return s.PostgresRepository.AssignIncomingShiftsToStaff(payload, targetStaffID)
}

// SwapIncomingShifts swaps two incoming shifts.
//
// Parameters:
// - payload: the request body containing the shift swap details.
// Returns:
// - []*domain.CycleIncomingCyclePickupShift: the swapped shifts.
// - []*domain.CycleIncomingCyclePickupShift: the swapped shifts.
// - error: an error if the swap process encounters any issues.
func (s *CycleService) SwapIncomingShifts(payload *models.CyclesCreateIncomingShiftSwapRequestBody) ([]*domain.CycleIncomingCyclePickupShift, []*domain.CycleIncomingCyclePickupShift, error) {
	return s.PostgresRepository.SwapIncomingShifts(payload)
}

// QueryChats queries the chats based on the provided data model.
//
// Parameters:
// - dataModel: a pointer to models.CyclesQueryChatsRequestParams containing the query parameters.
// Returns:
// - *restypes.QueryResponse: a pointer to the query response.
// - error: an error if the query process encounters any issues.
func (s *CycleService) QueryChats(dataModel *models.CyclesQueryChatsRequestParams) (*restypes.QueryResponse, error) {
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
// - dataModel: a pointer to models.CyclesQueryChatMessagesRequestParams containing the query parameters.
// Returns:
// - *restypes.QueryResponse: a pointer to the query response.
// - error: an error if the query process encounters any issues.
func (s *CycleService) QueryChatMessages(dataModel *models.CyclesQueryChatMessagesRequestParams) (*restypes.QueryResponse, error) {
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

// CreateChatMessage creates a new chat message in the CycleService.
//
// It takes a payload of type *models.CyclesCreateChatMessageRequestBody, which contains the cycle chat ID,
// sender user ID, recipient user ID, message, and attachments.
// It returns a pointer to a domain.CycleChatMessage and an error.
func (s *CycleService) CreateChatMessage(payload *models.CyclesCreateChatMessageRequestBody) (*domain.CycleChatMessage, error) {
	return s.PostgresRepository.CreateChatMessage(payload)
}

// UpdateChatMessageAttachments updates the attachments of a chat message.
//
// Parameters:
// - previousAttachments: a slice of types.UploadMetadata representing the previous attachments.
// - attachments: a slice of pointers to types.UploadMetadata representing the new attachments.
// - id: an int64 representing the ID of the chat message.
// Returns:
// - *domain.CycleChatMessage: a pointer to the updated chat message.
// - error: an error if the update process encounters any issues.
func (s *CycleService) UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CycleChatMessage, error) {
	return s.PostgresRepository.UpdateChatMessageAttachments(previousAttachments, attachments, id)
}
