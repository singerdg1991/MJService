package repositories

import (
	"fmt"
	"time"

	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/models"
)

type CycleRepositoryStub struct {
	Cycles []*domain.Cycle
}

type cycleTestCondition struct {
	HasError bool
}

var UserTestCondition *cycleTestCondition = &cycleTestCondition{}

func NewCycleRepositoryStub() *CycleRepositoryStub {
	return &CycleRepositoryStub{
		Cycles: []*domain.Cycle{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *CycleRepositoryStub) Query(dataModel *models.CyclesQueryRequestParams) ([]*domain.Cycle, error) {
	var cycles []*domain.Cycle
	for _, v := range r.Cycles {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			cycles = append(cycles, v)
			break
		}
	}
	return cycles, nil
}

func (r *CycleRepositoryStub) Count(dataModel *models.CyclesQueryRequestParams) (int64, error) {
	var cycles []*domain.Cycle
	for _, v := range r.Cycles {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			cycles = append(cycles, v)
			break
		}
	}
	return int64(len(cycles)), nil
}

func (r *CycleRepositoryStub) Create(payload *models.CyclesCreateRequestBody) (*domain.Cycle, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) Delete(payload *models.CyclesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) Update(payload *models.CyclesCreateRequestBody, id int64) (*domain.Cycle, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateStaffType(payload *models.CyclesUpdateStaffTypeRequestBody, id int64, isUnplanned bool) (*domain.Cycle, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateStaffTypes(payload *models.CyclesUpdateStaffTypesRequestBody, id int64) (*domain.Cycle, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryStaffTypes(queries *models.CyclesQueryStaffTypesRequestParams) ([]*domain.CycleStaffType, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountStaffTypes(queries *models.CyclesQueryStaffTypesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) Duplicate(payload *models.CyclesDuplicateRequestBody) (*domain.Cycle, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) FindAllStaffTypesByCycleID(cycleID int64) ([]*domain.CycleStaffType, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateNextStaffType(payload *models.CyclesUpdateNextStaffTypeRequestBody, id int64) ([]*domain.CycleNextStaffType, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateNextStaffTypes(payload *models.CyclesUpdateNextStaffTypesRequestBody, id int64) ([]*domain.CycleNextStaffType, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryNextStaffTypes(queries *models.CyclesQueryNextStaffTypesRequestParams) ([]*domain.CycleNextStaffType, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountNextStaffTypes(queries *models.CyclesQueryNextStaffTypesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) PickupShift(payload *models.CyclesCreatePickupShiftRequestBody) ([]*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) ([]*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) PickupShiftIncomingCycle(payload *models.CyclesCreateIncomingCyclePickupShiftRequestBody) ([]*domain.CycleIncomingCyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) ([]*domain.CycleIncomingCyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) FindAllNextStaffTypesByCycleID(cycleID int64) ([]*domain.CycleNextStaffType, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) FindVisitsForStaffInSpecificShift(cycleID int64, staffID int64, datetime *time.Time, shiftName string) ([]*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) AssignShiftsToStaff(payload *models.CyclesCreateShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) SwapShifts(payload *models.CyclesCreateShiftSwapRequestBody) ([]*domain.CyclePickupShift, []*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CreateShiftIfNotExist(cycleId int64, visitIds []int64, shiftName string, datetime *time.Time, status string) (*int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) ShiftStart(payload *models.CyclesCreateShiftStartRequestBody) (*domain.CycleShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) ShiftEnd(payload *models.CyclesCreateShiftEndRequestBody) (*domain.CycleShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) ([]*domain.CycleShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryShiftCustomerHomeKeys(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) ([]*domain.CycleShiftCustomerHomeKey, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountShiftCustomerHomeKeys(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CreateCycleShiftCustomerHomeKey(payload *models.CyclesShiftCustomerHomeKeyRequestBody) (*domain.CycleShiftCustomerHomeKey, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitStart(payload *models.CyclesCreateVisitStartRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitEnd(payload *models.CyclesCreateVisitEndRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitCancel(payload *models.CyclesCreateVisitCancelRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitDelay(payload *models.CyclesCreateVisitDelayRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitPause(payload *models.CyclesCreateVisitPauseRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitResume(payload *models.CyclesCreateVisitResumeRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) VisitReactive(payload *models.CyclesCreateVisitReactiveRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) AssignVisitToStaff(payload *models.CyclesCreateVisitAssignRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) SwapVisits(payload *models.CyclesCreateVisitSwapRequestBody) (*domain.CyclePickupShift, *domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CreateUnplannedVisit(payload *models.CyclesCreateVisitUnplannedRequestBody) (*domain.CyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CreateVisitTodo(payload *models.CyclesCreateVisitTodoRequestBody) (*domain.CyclePickupShiftTodo, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryVisitTodos(queries *models.CyclesQueryVisitsTodosRequestParams) ([]*domain.CyclePickupShiftTodo, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountVisitTodos(queries *models.CyclesQueryVisitsTodosRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateVisitTodoStatus(payload *models.CyclesUpdateVisitTodoStatusRequestBody, id int64) (*domain.CyclePickupShiftTodo, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateVisitTodoAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CyclePickupShiftTodo, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) AssignIncomingShiftsToStaff(payload *models.CyclesCreateIncomingShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CycleIncomingCyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) SwapIncomingShifts(payload *models.CyclesCreateIncomingShiftSwapRequestBody) ([]*domain.CycleIncomingCyclePickupShift, []*domain.CycleIncomingCyclePickupShift, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryChats(queries *models.CyclesQueryChatsRequestParams) ([]*domain.CycleChat, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountChats(queries *models.CyclesQueryChatsRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) QueryChatMessages(queries *models.CyclesQueryChatMessagesRequestParams) ([]*domain.CycleChatMessage, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CountChatMessages(queries *models.CyclesQueryChatMessagesRequestParams) (int64, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) CreateChatMessage(payload *models.CyclesCreateChatMessageRequestBody) (*domain.CycleChatMessage, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CycleChatMessage, error) {
	panic("implement me")
}

func (r *CycleRepositoryStub) UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(payload *models.CyclesUpdateStaffTypeRequestBody, migratedCycleID int64, currentCycleID int64) (*domain.Cycle, error) {
	panic("implement me")
}
