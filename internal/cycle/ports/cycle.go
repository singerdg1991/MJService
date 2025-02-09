package ports

import (
	"time"

	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/models"
)

type CycleService interface {
	Query(dataModel *models.CyclesQueryRequestParams) (*restypes.QueryResponse, error)
	QueryStaffTypes(dataModel *models.CyclesQueryStaffTypesRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.CyclesCreateRequestBody) (*domain.Cycle, error)
	Delete(payload *models.CyclesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.CyclesCreateRequestBody, id int64) (*domain.Cycle, error)
	UpdateStaffType(payload *models.CyclesUpdateStaffTypeRequestBody, id int64, isUnplanned bool) (*domain.Cycle, error)
	UpdateStaffTypes(payload *models.CyclesUpdateStaffTypesRequestBody, id int64) (*domain.Cycle, error)
	UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(payload *models.CyclesUpdateStaffTypeRequestBody, migratedCycleID int64, currentCycleID int64) (*domain.Cycle, error)
	FindByID(id int64) (*domain.Cycle, error)
	GetCurrent() (*domain.Cycle, error)
	GetLast() (*domain.Cycle, error)
	Duplicate(payload *models.CyclesDuplicateRequestBody) (*domain.Cycle, error)
	FindAllStaffTypesByCycleID(cycleID int64) ([]*domain.CycleStaffType, error)
	QueryNextStaffTypes(dataModel *models.CyclesQueryNextStaffTypesRequestParams) (*restypes.QueryResponse, error)
	UpdateNextStaffType(payload *models.CyclesUpdateNextStaffTypeRequestBody, QueryIncomingCyclePickupShiftsid int64) ([]*domain.CycleNextStaffType, error)
	UpdateNextStaffTypes(payload *models.CyclesUpdateNextStaffTypesRequestBody, id int64) ([]*domain.CycleNextStaffType, error)
	PickupShift(payload *models.CyclesCreatePickupShiftRequestBody) (*restypes.QueryResponse, []*domain.CyclePickupShift, error)
	QueryPickupShifts(dataModel *models.CyclesQueryPickupShiftsRequestParams) (*restypes.QueryResponse, error)
	FindPickupShiftByID(id int64) (*domain.CyclePickupShift, error)
	CountPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) (int64, error)
	FindPickedUpShiftByID(id int64) (*domain.CyclePickupShift, error)
	FindPickedUpShiftForStaff(cycleId int64, staffId int64, staffTypeId int64) (*domain.CyclePickupShift, error)
	PickupShiftIncomingCycle(payload *models.CyclesCreateIncomingCyclePickupShiftRequestBody) (*restypes.QueryResponse, error)
	QueryIncomingCyclePickupShifts(dataModel *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (*restypes.QueryResponse, error)
	CountIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (int64, error)
	FindIncomingCyclePickedUpShiftByID(id int64) (*domain.CycleIncomingCyclePickupShift, error)
	FindIncomingCyclePickedUpShiftForStaff(cycleId int64, staffId int64, staffTypeId int64) (*domain.CycleIncomingCyclePickupShift, error)
	FindAllNextStaffTypesByCycleID(cycleID int64) ([]*domain.CycleNextStaffType, error)
	FindVisitsForStaffInSpecificShift(cycleID int64, staffID int64, datetime *time.Time, shiftName string) ([]*domain.CyclePickupShift, error)
	AssignShiftsToStaff(payload *models.CyclesCreateShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CyclePickupShift, error)
	SwapShifts(payload *models.CyclesCreateShiftSwapRequestBody) ([]*domain.CyclePickupShift, []*domain.CyclePickupShift, error)
	CreateShiftIfNotExist(cycleId int64, visitIds []int64, shiftName string, datetime *time.Time, status string) (*int64, error)
	ShiftStart(payload *models.CyclesCreateShiftStartRequestBody) (*domain.CycleShift, error)
	ShiftEnd(payload *models.CyclesCreateShiftEndRequestBody) (*domain.CycleShift, error)
	QueryCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) (*restypes.QueryResponse, error)
	QueryShiftCustomerHomeKeys(dataModel *models.CyclesQueryShiftCustomerHomeKeysRequestParams) (*restypes.QueryResponse, error)
	FindShiftByID(id int64) (*domain.CycleShift, error)
	CreateCycleShiftCustomerHomeKey(payload *models.CyclesShiftCustomerHomeKeyRequestBody) (*domain.CycleShiftCustomerHomeKey, error)
	VisitStart(payload *models.CyclesCreateVisitStartRequestBody) (*domain.CyclePickupShift, error)
	VisitEnd(payload *models.CyclesCreateVisitEndRequestBody) (*domain.CyclePickupShift, error)
	VisitCancel(payload *models.CyclesCreateVisitCancelRequestBody) (*domain.CyclePickupShift, error)
	VisitDelay(payload *models.CyclesCreateVisitDelayRequestBody) (*domain.CyclePickupShift, error)
	VisitPause(payload *models.CyclesCreateVisitPauseRequestBody) (*domain.CyclePickupShift, error)
	VisitResume(payload *models.CyclesCreateVisitResumeRequestBody) (*domain.CyclePickupShift, error)
	VisitReactive(payload *models.CyclesCreateVisitReactiveRequestBody) (*domain.CyclePickupShift, error)
	AssignVisitToStaff(payload *models.CyclesCreateVisitAssignRequestBody) (*domain.CyclePickupShift, error)
	SwapVisits(payload *models.CyclesCreateVisitSwapRequestBody) (*domain.CyclePickupShift, *domain.CyclePickupShift, error)
	CreateUnplannedVisit(payload *models.CyclesCreateVisitUnplannedRequestBody) (*domain.CyclePickupShift, error)
	CreateVisitTodo(payload *models.CyclesCreateVisitTodoRequestBody) (*domain.CyclePickupShiftTodo, error)
	QueryVisitTodos(dataModel *models.CyclesQueryVisitsTodosRequestParams) (*restypes.QueryResponse, error)
	FindVisitTodoByID(id int64) (*domain.CyclePickupShiftTodo, error)
	UpdateVisitTodoStatus(payload *models.CyclesUpdateVisitTodoStatusRequestBody, id int64) (*domain.CyclePickupShiftTodo, error)
	UpdateVisitTodoAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CyclePickupShiftTodo, error)
	AssignIncomingShiftsToStaff(payload *models.CyclesCreateIncomingShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CycleIncomingCyclePickupShift, error)
	SwapIncomingShifts(payload *models.CyclesCreateIncomingShiftSwapRequestBody) ([]*domain.CycleIncomingCyclePickupShift, []*domain.CycleIncomingCyclePickupShift, error)
	QueryChats(dataModel *models.CyclesQueryChatsRequestParams) (*restypes.QueryResponse, error)
	QueryChatMessages(dataModel *models.CyclesQueryChatMessagesRequestParams) (*restypes.QueryResponse, error)
	CreateChatMessage(payload *models.CyclesCreateChatMessageRequestBody) (*domain.CycleChatMessage, error)
	UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CycleChatMessage, error)
}

type CycleRepositoryPostgresDB interface {
	Query(dataModel *models.CyclesQueryRequestParams) ([]*domain.Cycle, error)
	QueryStaffTypes(queries *models.CyclesQueryStaffTypesRequestParams) ([]*domain.CycleStaffType, error)
	Count(dataModel *models.CyclesQueryRequestParams) (int64, error)
	CountStaffTypes(queries *models.CyclesQueryStaffTypesRequestParams) (int64, error)
	Create(payload *models.CyclesCreateRequestBody) (*domain.Cycle, error)
	Delete(payload *models.CyclesDeleteRequestBody) ([]int64, error)
	Update(payload *models.CyclesCreateRequestBody, id int64) (*domain.Cycle, error)
	UpdateStaffType(payload *models.CyclesUpdateStaffTypeRequestBody, id int64, isUnplanned bool) (*domain.Cycle, error)
	UpdateStaffTypes(payload *models.CyclesUpdateStaffTypesRequestBody, id int64) (*domain.Cycle, error)
	UpdateStaffTypeAndPickupShiftsMigratedFromLastIncomingCycle(payload *models.CyclesUpdateStaffTypeRequestBody, migratedCycleID int64, currentCycleID int64) (*domain.Cycle, error)
	Duplicate(payload *models.CyclesDuplicateRequestBody) (*domain.Cycle, error)
	FindAllStaffTypesByCycleID(cycleID int64) ([]*domain.CycleStaffType, error)
	QueryNextStaffTypes(queries *models.CyclesQueryNextStaffTypesRequestParams) ([]*domain.CycleNextStaffType, error)
	UpdateNextStaffType(payload *models.CyclesUpdateNextStaffTypeRequestBody, id int64) ([]*domain.CycleNextStaffType, error)
	UpdateNextStaffTypes(payload *models.CyclesUpdateNextStaffTypesRequestBody, id int64) ([]*domain.CycleNextStaffType, error)
	CountNextStaffTypes(queries *models.CyclesQueryNextStaffTypesRequestParams) (int64, error)
	PickupShift(payload *models.CyclesCreatePickupShiftRequestBody) ([]*domain.CyclePickupShift, error)
	QueryPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) ([]*domain.CyclePickupShift, error)
	CountPickupShifts(queries *models.CyclesQueryPickupShiftsRequestParams) (int64, error)
	PickupShiftIncomingCycle(payload *models.CyclesCreateIncomingCyclePickupShiftRequestBody) ([]*domain.CycleIncomingCyclePickupShift, error)
	QueryIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) ([]*domain.CycleIncomingCyclePickupShift, error)
	CountIncomingCyclePickupShifts(queries *models.CyclesQueryIncomingCyclePickupShiftsRequestParams) (int64, error)
	FindAllNextStaffTypesByCycleID(cycleID int64) ([]*domain.CycleNextStaffType, error)
	FindVisitsForStaffInSpecificShift(cycleID int64, staffID int64, datetime *time.Time, shiftName string) ([]*domain.CyclePickupShift, error)
	AssignShiftsToStaff(payload *models.CyclesCreateShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CyclePickupShift, error)
	SwapShifts(payload *models.CyclesCreateShiftSwapRequestBody) ([]*domain.CyclePickupShift, []*domain.CyclePickupShift, error)
	CreateShiftIfNotExist(cycleId int64, visitIds []int64, shiftName string, datetime *time.Time, status string) (*int64, error)
	ShiftStart(payload *models.CyclesCreateShiftStartRequestBody) (*domain.CycleShift, error)
	ShiftEnd(payload *models.CyclesCreateShiftEndRequestBody) (*domain.CycleShift, error)
	CountCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) (int64, error)
	QueryCycleShifts(queries *models.CyclesQueryCycleShiftsRequestParams) ([]*domain.CycleShift, error)
	QueryShiftCustomerHomeKeys(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) ([]*domain.CycleShiftCustomerHomeKey, error)
	CountShiftCustomerHomeKeys(queries *models.CyclesQueryShiftCustomerHomeKeysRequestParams) (int64, error)
	CreateCycleShiftCustomerHomeKey(payload *models.CyclesShiftCustomerHomeKeyRequestBody) (*domain.CycleShiftCustomerHomeKey, error)
	VisitStart(payload *models.CyclesCreateVisitStartRequestBody) (*domain.CyclePickupShift, error)
	VisitEnd(payload *models.CyclesCreateVisitEndRequestBody) (*domain.CyclePickupShift, error)
	VisitCancel(payload *models.CyclesCreateVisitCancelRequestBody) (*domain.CyclePickupShift, error)
	VisitDelay(payload *models.CyclesCreateVisitDelayRequestBody) (*domain.CyclePickupShift, error)
	VisitPause(payload *models.CyclesCreateVisitPauseRequestBody) (*domain.CyclePickupShift, error)
	VisitResume(payload *models.CyclesCreateVisitResumeRequestBody) (*domain.CyclePickupShift, error)
	VisitReactive(payload *models.CyclesCreateVisitReactiveRequestBody) (*domain.CyclePickupShift, error)
	AssignVisitToStaff(payload *models.CyclesCreateVisitAssignRequestBody) (*domain.CyclePickupShift, error)
	SwapVisits(payload *models.CyclesCreateVisitSwapRequestBody) (*domain.CyclePickupShift, *domain.CyclePickupShift, error)
	CreateUnplannedVisit(payload *models.CyclesCreateVisitUnplannedRequestBody) (*domain.CyclePickupShift, error)
	CreateVisitTodo(payload *models.CyclesCreateVisitTodoRequestBody) (*domain.CyclePickupShiftTodo, error)
	QueryVisitTodos(queries *models.CyclesQueryVisitsTodosRequestParams) ([]*domain.CyclePickupShiftTodo, error)
	CountVisitTodos(queries *models.CyclesQueryVisitsTodosRequestParams) (int64, error)
	UpdateVisitTodoStatus(payload *models.CyclesUpdateVisitTodoStatusRequestBody, id int64) (*domain.CyclePickupShiftTodo, error)
	UpdateVisitTodoAttachments(attachments []*types.UploadMetadata, id int64) (*domain.CyclePickupShiftTodo, error)
	AssignIncomingShiftsToStaff(payload *models.CyclesCreateIncomingShiftAssignToMeRequestBody, targetStaffID int64) ([]*domain.CycleIncomingCyclePickupShift, error)
	SwapIncomingShifts(payload *models.CyclesCreateIncomingShiftSwapRequestBody) ([]*domain.CycleIncomingCyclePickupShift, []*domain.CycleIncomingCyclePickupShift, error)
	QueryChats(queries *models.CyclesQueryChatsRequestParams) ([]*domain.CycleChat, error)
	CountChats(queries *models.CyclesQueryChatsRequestParams) (int64, error)
	QueryChatMessages(queries *models.CyclesQueryChatMessagesRequestParams) ([]*domain.CycleChatMessage, error)
	CountChatMessages(queries *models.CyclesQueryChatMessagesRequestParams) (int64, error)
	CreateChatMessage(payload *models.CyclesCreateChatMessageRequestBody) (*domain.CycleChatMessage, error)
	UpdateChatMessageAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.CycleChatMessage, error)
}

type CycleRepositoryMongoDB interface {
	Query(queries *models.CyclesQueryRequestParams) ([]*domain.Cycle, error)
	Count(queries *models.CyclesQueryRequestParams) (int64, error)
	Create(postgresID int, payload interface{}) (interface{}, error)
	UpdateByPostgresID(postgresID int, payload interface{}) (interface{}, error)
	Update(payload interface{}, id int) error
	Delete(ids []uint) error
}
