package handlers

import (
	"errors"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	csPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	"github.com/hoitek/Maja-Service/internal/cycle/config"
	"github.com/hoitek/Maja-Service/internal/cycle/constants"
	"github.com/hoitek/Maja-Service/internal/cycle/domain"
	"github.com/hoitek/Maja-Service/internal/cycle/models"
	"github.com/hoitek/Maja-Service/internal/cycle/ports"
	rolePorts "github.com/hoitek/Maja-Service/internal/role/ports"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	ssPorts "github.com/hoitek/Maja-Service/internal/service/ports"
	stypePorts "github.com/hoitek/Maja-Service/internal/servicetype/ports"
	stPorts "github.com/hoitek/Maja-Service/internal/staff/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"

	"log"
	"net/http"
)

type CycleHandler struct {
	CycleService       ports.CycleService
	RoleService        rolePorts.RoleService
	UserService        uPorts.UserService
	StaffService       stPorts.StaffService
	CustomerService    csPorts.CustomerService
	ServiceService     ssPorts.ServiceService
	ServiceTypeService stypePorts.ServiceTypeService
	S3Service          s3Ports.S3Service
}

func NewCycleHandler(
	r *mux.Router,
	cs ports.CycleService,
	rs rolePorts.RoleService,
	us uPorts.UserService,
	ss stPorts.StaffService,
	css csPorts.CustomerService,
	ssv ssPorts.ServiceService,
	stsv stypePorts.ServiceTypeService,
	s3s s3Ports.S3Service,
) (CycleHandler, error) {
	cycleHandler := CycleHandler{
		CycleService:       cs,
		RoleService:        rs,
		UserService:        us,
		StaffService:       ss,
		CustomerService:    css,
		ServiceService:     ssv,
		ServiceTypeService: stsv,
		S3Service:          s3s,
	}
	if r == nil {
		return CycleHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.CycleConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.CycleConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	// Authenticated routes
	rAuth.Handle("/cycles", cycleHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/cycles", cycleHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/cycles", cycleHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/cycles/current", cycleHandler.QueryCurrent()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/last", cycleHandler.QueryLast()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/duplicate", cycleHandler.Duplicate()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/stafftypes", cycleHandler.QueryStaffTypes()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/stafftype/{cycleid}", cycleHandler.UpdateStaffType()).Methods(http.MethodPut)
	rAuth.Handle("/cycles/stafftypes/{cycleid}", cycleHandler.UpdateStaffTypes()).Methods(http.MethodPut)
	rAuth.Handle("/cycles/{id}", cycleHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/cycles/csv/download", cycleHandler.Download()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/visits/{cycleid}", cycleHandler.QueryCycleVisits()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/nextstafftypes", cycleHandler.QueryNextStaffTypes()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/nextstafftype/{currentcycleid}", cycleHandler.UpdateNextStaffType()).Methods(http.MethodPut)
	rAuth.Handle("/cycles/nextstafftypes/{currentcycleid}", cycleHandler.UpdateNextStaffTypes()).Methods(http.MethodPut)
	rAuth.Handle("/cycles/arrangement/pickup", cycleHandler.CreateArrangementPickupShift()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/arrangement/shifts", cycleHandler.QueryArrangementShifts()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/incoming/arrangement/pickup", cycleHandler.CreateIncomingArrangementPickupShift()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/incoming/arrangement/shifts", cycleHandler.QueryIncomingArrangementShifts()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/incoming/shift/assign-to-me", cycleHandler.CreateIncomingShiftAssignToMe()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/incoming/shift/swap", cycleHandler.CreateIncomingShiftSwap()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shifts/{shiftid}/start", cycleHandler.CreateShiftsStart()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shifts/{shiftid}/end", cycleHandler.CreateShiftsEnd()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shift/assign-to-me", cycleHandler.CreateShiftAssignToMe()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shift/swap", cycleHandler.CreateShiftSwap()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shifts", cycleHandler.QueryCycleShifts()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/shifts/customer-home-key", cycleHandler.CreateCycleShiftsCustomerHomeKey()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shifts/customer-home-keys", cycleHandler.QueryCycleShiftsCustomerHomeKeys()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/visits/start", cycleHandler.CreateCycleVisitStart()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visits/end", cycleHandler.CreateCycleVisitEnd()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visits/cancel", cycleHandler.CreateCycleVisitCancel()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visits/delay", cycleHandler.CreateCycleVisitDelay()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visits/pause", cycleHandler.CreateCycleVisitPause()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visits/resume", cycleHandler.CreateCycleVisitResume()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visits/reactive", cycleHandler.CreateCycleVisitReactive()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shifts/customer-home-key/release", cycleHandler.CreateCycleShiftsCustomerHomeKeyRelease()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/shifts/customer-home-key/not-release", cycleHandler.CreateCycleShiftsCustomerHomeKeyNotRelease()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visit/assign", cycleHandler.CreateVisitAssign()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visit/swap", cycleHandler.CreateVisitSwap()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visit/unplanned", cycleHandler.CreateVisitUnplanned()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visit/todo", cycleHandler.CreateVisitTodo()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visit/sos", cycleHandler.CreateVisitSos()).Methods(http.MethodPost)
	rAuth.Handle("/cycles/visit/todos", cycleHandler.QueryVisitTodos()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/visit/todo/status/{visittodoid}", cycleHandler.UpdateVisitTodoStatus()).Methods(http.MethodPut)
	rAuth.Handle("/cycles/chats", cycleHandler.QueryChats()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/chats/messages", cycleHandler.QueryChatMessages()).Methods(http.MethodGet)
	rAuth.Handle("/cycles/chats/messages", cycleHandler.CreateChatMessage()).Methods(http.MethodPost)

	return cycleHandler, nil
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: cycles
* @apiResponseRef: CyclesQueryResponse
* @apiSummary: Query cycles
* @apiParametersRef: CyclesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		cycles, err := h.CycleService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(cycles)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/visits/{cycleid}
* @apiResponseRef: CyclesQueryCycleVisitsResponse
* @apiSummary: Query cycle visits
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryCycleVisitsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryCycleVisitsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryCycleVisits() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		return response.Success(true)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/current
* @apiResponseRef: CyclesQueryCurrentResponse
* @apiSummary: Query current cycle
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryCurrentNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryCurrentNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryCurrent() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		cycle, err := h.CycleService.GetCurrent()
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Check cycle status
		if cycle.Status != constants.STATUS_ACTIVE {
			return response.ErrorNotFound(nil, "There is not active cycle yet")
		}

		return response.Success(cycle)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/last
* @apiResponseRef: CyclesQueryLastResponse
* @apiSummary: Query current cycle
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryLastNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryLastNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryLast() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		cycle, err := h.CycleService.GetLast()
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		return response.Success(cycle)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/stafftypes
* @apiResponseRef: CyclesQueryStaffTypesResponse
* @apiSummary: Query staff types
* @apiParametersRef: CyclesQueryStaffTypesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryStaffTypesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryStaffTypesNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryStaffTypes() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryStaffTypesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate roleIds
		for _, roleId := range queries.RoleIDsInt64 {
			role := h.RoleService.GetRoleByID(int(roleId))
			if role == nil {
				return response.ErrorBadRequest(nil, "Role not found")
			}
		}

		// Query staff types
		staffTypes, err := h.CycleService.QueryStaffTypes(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffTypes)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: cycles/csv/download
* @apiResponseRef: CyclesQueryResponse
* @apiSummary: Query cycles
* @apiParametersRef: CyclesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		cycles, err := h.CycleService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(cycles)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/arrangement/pickup
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreatePickupShiftRequestBody
 * @apiResponseRef: CyclesCreatePickupShiftResponse
 * @apiSummary: Pickup shifts for staff in arrangement
 * @apiDescription: Pickup shifts for staff in arrangement
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateArrangementPickupShift() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		//user := h.UserService.GetUserFromContext(r.Context())

		// Generate unique correlation ID
		//correlationID := utils.GenerateCorrelationID()
		//
		//// Send the request to the queue with correlation ID
		//message := "test"
		//err := queues.DefaultCycleArrangementQueue.Publish(constants.ARRANGEMENT_QUEUE_WISH_NAME, message, correlationID)
		//if err != nil {
		//	return response.ErrorInternalServerError("Failed to send request to queue")
		//}
		//
		//// Wait for response with the same correlation ID
		//responseMsg, err := queues.DefaultCycleArrangementQueue.WaitForResponseFromQueue(constants.ARRANGEMENT_QUEUE_WISH_RESPONSE_NAME, correlationID)
		//if err != nil {
		//	log.Println(err)
		//	return response.ErrorInternalServerError("Failed to get response from queue")
		//}

		// Validate request body
		payload := &models.CyclesCreatePickupShiftRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff by id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil || staff == nil {
			return response.ErrorBadRequest(nil, "Staff not found")
		}
		payload.Staff = &domain.CyclePickupShiftStaff{
			ID:        staff.ID,
			UserID:    staff.UserID,
			FirstName: staff.User.FirstName,
			LastName:  staff.User.LastName,
		}

		// Validate cycleStaffTypes
		cycleStaffTypes, err := h.CycleService.FindAllStaffTypesByCycleID(int64(payload.CycleID))
		if err != nil || cycleStaffTypes == nil {
			if err != nil {
				log.Printf("Error finding cycle staff types: %v\n", err.Error())
			}
			return response.ErrorBadRequest(nil, "cycleStaffTypes not found")
		}
		var foundStaffType = false
		var shiftName string
		for _, cycleStaffTypeID := range payload.CycleStaffTypeIDsInt64 {
			if foundStaffType {
				break
			}
			for _, cycleStaffType := range cycleStaffTypes {
				if cycleStaffType.ID == uint(cycleStaffTypeID) {
					foundStaffType = true
					shiftName = cycleStaffType.ShiftName
					break
				}
			}
		}
		if !foundStaffType {
			return response.ErrorBadRequest(nil, "cycleStaffTypes not found")
		}
		for _, cycleStaffType := range cycleStaffTypes {
			for _, cycleStaffTypeID := range payload.CycleStaffTypeIDsInt64 {
				if cycleStaffType.ID == uint(cycleStaffTypeID) {
					if cycleStaffType.RemindStaffCount <= 0 {
						return response.ErrorBadRequest(nil, fmt.Sprintf("cycleStaffType with id %d has no remind staff count", cycleStaffType.ID))
					}
				}
			}
		}
		payload.ShiftName = shiftName

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Check pickedUpShift already exists
		//for _, cycleStaffTypeID := range payload.CycleStaffTypeIDsInt64 {
		//	pickedUpShift, _ := h.CycleService.FindPickedUpShiftForStaff(int64(payload.CycleID), int64(payload.StaffID), cycleStaffTypeID)
		//	if pickedUpShift != nil {
		//		return response.ErrorBadRequest(nil, fmt.Sprintf("cycleStaffType with id \"%d\" has already picked up shift for staff with id \"%d\"\n", cycleStaffTypeID, payload.StaffID))
		//	}
		//}

		// Pickup shift
		res, _, err := h.CycleService.PickupShift(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(res)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: cycles/arrangement/shifts
* @apiResponseRef: CyclesQueryPickupShiftsResponse
* @apiSummary: Query pickedup shifts
* @apiParametersRef: CyclesQueryPickupShiftsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryPickupShiftsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryPickupShiftsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryArrangementShifts() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryPickupShiftsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		pickupShifts, err := h.CycleService.QueryPickupShifts(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(pickupShifts)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/incoming/arrangement/pickup
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateIncomingCyclePickupShiftRequestBody
 * @apiResponseRef: CyclesCreateIncomingCyclePickupShiftResponse
 * @apiSummary: Pickup shifts for staff in incoming cycle arrangement
 * @apiDescription: Pickup shifts for staff in incoming cycle arrangement
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateIncomingArrangementPickupShift() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CyclesCreateIncomingCyclePickupShiftRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff by id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil || staff == nil {
			return response.ErrorBadRequest(nil, "Staff not found")
		}
		payload.Staff = &domain.CycleIncomingCyclePickupShiftStaff{
			ID:        staff.ID,
			UserID:    staff.UserID,
			FirstName: staff.User.FirstName,
			LastName:  staff.User.LastName,
		}

		// Validate cycleNextStaffTypes
		cycleNextStaffTypes, err := h.CycleService.FindAllNextStaffTypesByCycleID(int64(payload.CycleID))
		if err != nil || cycleNextStaffTypes == nil {
			if err != nil {
				log.Printf("Error finding cycle next staff types: %v\n", err.Error())
			}
			return response.ErrorBadRequest(nil, "cycleNextStaffTypes not found")
		}
		var foundNextStaffType = false
		for _, cycleNextStaffTypeID := range payload.CycleNextStaffTypeIDsInt64 {
			if foundNextStaffType {
				break
			}
			for _, cycleNextStaffType := range cycleNextStaffTypes {
				if cycleNextStaffType.ID == uint(cycleNextStaffTypeID) {
					foundNextStaffType = true
					break
				}
			}
		}
		if !foundNextStaffType {
			return response.ErrorBadRequest(nil, "cycleNextStaffTypes not found")
		}
		for _, cycleNextStaffType := range cycleNextStaffTypes {
			for _, cycleStaffTypeID := range payload.CycleNextStaffTypeIDsInt64 {
				if cycleNextStaffType.ID == uint(cycleStaffTypeID) {
					if cycleNextStaffType.RemindStaffCount <= 0 {
						return response.ErrorBadRequest(nil, fmt.Sprintf("cycleNextStaffType with id %d has no remind staff count", cycleNextStaffType.ID))
					}
				}
			}
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Check pickedUpShift already exists
		for _, cycleNextStaffTypeID := range payload.CycleNextStaffTypeIDsInt64 {
			pickedUpShift, _ := h.CycleService.FindIncomingCyclePickedUpShiftForStaff(int64(payload.CycleID), int64(payload.StaffID), cycleNextStaffTypeID)
			if pickedUpShift != nil {
				return response.ErrorBadRequest(nil, fmt.Sprintf("cycleNextStaffType with id \"%d\" has already picked up shift for staff with id \"%d\"\n", cycleNextStaffTypeID, payload.StaffID))
			}
		}

		// Pickup shift
		data, err := h.CycleService.PickupShiftIncomingCycle(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/incoming/arrangement/shifts
* @apiResponseRef: CyclesQueryIncomingCyclePickupShiftsResponse
* @apiSummary: Query staff types
* @apiParametersRef: CyclesQueryIncomingCyclePickupShiftsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryPickupShiftsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryIncomingCyclePickupShiftsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryIncomingArrangementShifts() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryIncomingCyclePickupShiftsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		pickupShifts, err := h.CycleService.QueryIncomingCyclePickupShifts(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(pickupShifts)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateRequestBody
 * @apiResponseRef: CyclesCreateResponse
 * @apiSummary: Create cycle
 * @apiDescription: Create cycle
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CyclesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get last cycle
		// TODO: Comment out this code when test is done
		//lastCycle, err := h.CycleService.GetLast()
		//if err != nil {
		//	return response.ErrorInternalServerError(nil, err.Error())
		//}
		//if lastCycle.Status != constants.STATUS_EXPIRED {
		//	return response.ErrorBadRequest(nil, "Last cycle is not expired")
		//}

		// Create cycle
		data, err := h.CycleService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/duplicate
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesDuplicateRequestBody
 * @apiResponseRef: CyclesCreateResponse
 * @apiSummary: Duplicate cycle
 * @apiDescription: Duplicate cycle
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) Duplicate() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CyclesDuplicateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(payload.CycleID)
		if err != nil {
			log.Printf("CycleHandler.Duplicate: %s", err.Error())
			return response.ErrorInternalServerError(nil, "cycle not found")
		}
		if cycle == nil {
			return response.ErrorNotFound(nil, "cycle not found")
		}
		payload.Cycle = cycle

		// Find staff types
		staffTypes, err := h.CycleService.FindAllStaffTypesByCycleID(payload.CycleID)
		if err != nil {
			log.Printf("CycleHandler.Duplicate: %s", err.Error())
			return response.ErrorInternalServerError(nil, "something went wrong, please try again later")
		}
		payload.CycleStaffTypes = staffTypes

		// Duplicate cycle
		data, err := h.CycleService.Duplicate(payload)
		if err != nil {
			log.Printf("CycleHandler.Duplicate: %s", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: CyclesUpdateRequestParams
 * @apiRequestRef: CyclesCreateRequestBody
 * @apiResponseRef: CyclesCreateResponse
 * @apiSummary: Update cycle
 * @apiDescription: Update cycle
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CyclesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CyclesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(params.ID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Update cycle
		data, err := h.CycleService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/stafftype/{cycleid}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: CyclesUpdateStaffTypeRequestParams
 * @apiRequestRef: CyclesUpdateStaffTypeRequestBody
 * @apiResponseRef: CyclesCreateResponse
 * @apiSummary: Update cycle staff type
 * @apiDescription: Update cycle staff type
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) UpdateStaffType() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CyclesUpdateStaffTypeRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CyclesUpdateStaffTypeRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find role by id
		role := h.RoleService.GetRoleByID(payload.RoleID)
		if role == nil {
			return response.ErrorBadRequest(nil, "Role not found")
		}
		payload.Role = &domain.CycleStaffTypeRole{
			ID:   role.ID,
			Name: role.Name,
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(params.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Update cycle
		data, err := h.CycleService.UpdateStaffType(payload, int64(params.CycleID), false)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/stafftypes/{cycleid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CyclesUpdateStaffTypesRequestParams
 * @apiRequestRef: CyclesUpdateStaffTypesRequestBody
 * @apiResponseRef: CyclesCreateResponse
 * @apiSummary: Update cycle staff types
 * @apiDescription: Update cycle staff types
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) UpdateStaffTypes() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CyclesUpdateStaffTypesRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CyclesUpdateStaffTypesRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(params.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Find role by id and set to staff type for each staff type
		for _, staffType := range payload.StaffTypes {
			// Find role by id
			role := h.RoleService.GetRoleByID(staffType.RoleID)
			if role == nil {
				return response.ErrorBadRequest(nil, "Role not found")
			}
			staffType.Role = &domain.CycleStaffTypeRole{
				ID:   role.ID,
				Name: role.Name,
			}
		}

		// Update cycle
		data, err := h.CycleService.UpdateStaffTypes(payload, int64(params.CycleID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: CyclesDeleteRequestBody
 * @apiResponseRef: CyclesCreateResponse
 * @apiSummary: Delete cycle
 * @apiDescription: Delete cycle
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.CyclesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.CycleService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/nextstafftypes
* @apiResponseRef: CyclesQueryNextStaffTypesResponse
* @apiSummary: Query staff types
* @apiParametersRef: CyclesQueryNextStaffTypesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryNextStaffTypesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryNextStaffTypesNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryNextStaffTypes() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryNextStaffTypesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate roleIds
		for _, roleId := range queries.RoleIDsInt64 {
			role := h.RoleService.GetRoleByID(int(roleId))
			if role == nil {
				return response.ErrorBadRequest(nil, "Role not found")
			}
		}

		// Query staff types
		staffTypes, err := h.CycleService.QueryNextStaffTypes(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffTypes)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/nextstafftype/{currentcycleid}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: CyclesUpdateNextStaffTypeRequestParams
 * @apiRequestRef: CyclesUpdateNextStaffTypeRequestBody
 * @apiResponseRef: CyclesUpdateNextStaffTypeResponse
 * @apiSummary: Update cycle next staff type
 * @apiDescription: Update cycle next staff type
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) UpdateNextStaffType() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CyclesUpdateNextStaffTypeRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CyclesUpdateNextStaffTypeRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find role by id
		role := h.RoleService.GetRoleByID(payload.RoleID)
		if role == nil {
			return response.ErrorBadRequest(nil, "Role not found")
		}
		payload.Role = &domain.CycleNextStaffTypeRole{
			ID:   role.ID,
			Name: role.Name,
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(params.CurrentCycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Update cycle
		data, err := h.CycleService.UpdateNextStaffType(payload, int64(params.CurrentCycleID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/nextstafftypes/{currentcycleid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CyclesUpdateNextStaffTypesRequestParams
 * @apiRequestRef: CyclesUpdateNextStaffTypesRequestBody
 * @apiResponseRef: CyclesUpdateNextStaffTypesResponse
 * @apiSummary: Update cycle next staff types
 * @apiDescription: Update cycle next staff types
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) UpdateNextStaffTypes() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CyclesUpdateNextStaffTypesRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CyclesUpdateNextStaffTypesRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(params.CurrentCycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Find role by id and set to staff type for each staff type
		for _, staffType := range payload.StaffTypes {
			// Find role by id
			role := h.RoleService.GetRoleByID(staffType.RoleID)
			if role == nil {
				return response.ErrorBadRequest(nil, "Role not found")
			}
			staffType.Role = &domain.CycleNextStaffTypeRole{
				ID:   role.ID,
				Name: role.Name,
			}
		}

		// Update cycle
		data, err := h.CycleService.UpdateNextStaffTypes(payload, int64(params.CurrentCycleID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shifts/{shiftid}/start
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateShiftStartRequestBody
 * @apiResponseRef: CyclesCreateShiftStartResponse
 * @apiSummary: Start shift
 * @apiDescription: Start shift
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateShiftsStart() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CyclesCreateShiftStartRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Pickup shift
		data, err := h.CycleService.ShiftStart(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shifts/{shiftid}/end
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateShiftEndRequestBody
 * @apiResponseRef: CyclesCreateShiftEndResponse
 * @apiSummary: End shift
 * @apiDescription: End shift
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateShiftsEnd() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request body
		payload := &models.CyclesCreateShiftEndRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Pickup shift
		data, err := h.CycleService.ShiftEnd(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shift/assign-to-me
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateShiftAssignToMeRequestBody
 * @apiResponseRef: CyclesCreateShiftAssignToMeResponse
 * @apiSummary: Assign shift to me
 * @apiDescription: Assign shift to me
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateShiftAssignToMe() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateShiftAssignToMeRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff by id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil || staff == nil {
			return response.ErrorBadRequest(nil, "Staff not found")
		}
		payload.Staff = &domain.CyclePickupShiftStaff{
			ID:        staff.ID,
			UserID:    staff.UserID,
			FirstName: staff.User.FirstName,
			LastName:  staff.User.LastName,
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Assign to me
		targetStaffID := user.StaffID
		if targetStaffID == nil {
			return response.ErrorBadRequest(nil, "You are authenticated as a test user without staff. Please create a staff and authenticate with it and try again.")
		}
		data, err := h.CycleService.AssignShiftsToStaff(payload, int64(*targetStaffID))
		if err != nil {
			log.Printf("CycleHandler.CreateShiftAssignToMe: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shift/swap
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateShiftSwapRequestBody
 * @apiResponseRef: CyclesCreateShiftSwapResponse
 * @apiSummary: Swap shift to another staff
 * @apiDescription: Swap shift to another staff
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateShiftSwap() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateShiftSwapRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find source staff by id
		sourceStaff, err := h.StaffService.FindByID(payload.SourceStaffID)
		if err != nil || sourceStaff == nil {
			return response.ErrorBadRequest(nil, "Source staff not found")
		}
		payload.SourceStaff = &domain.CyclePickupShiftStaff{
			ID:        sourceStaff.ID,
			UserID:    sourceStaff.UserID,
			FirstName: sourceStaff.User.FirstName,
			LastName:  sourceStaff.User.LastName,
		}

		// Find source staff by id
		targetStaff, err := h.StaffService.FindByID(payload.TargetStaffID)
		if err != nil || targetStaff == nil {
			return response.ErrorBadRequest(nil, "Target staff not found")
		}
		payload.TargetStaff = &domain.CyclePickupShiftStaff{
			ID:        targetStaff.ID,
			UserID:    targetStaff.UserID,
			FirstName: targetStaff.User.FirstName,
			LastName:  targetStaff.User.LastName,
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Swap shifts
		sourceShifts, targetShifts, err := h.CycleService.SwapShifts(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateShiftSwap: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(models.CyclesCreateShiftSwapResponseData{
			SourceShifts: sourceShifts,
			TargetShifts: targetShifts,
		})
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/shifts
* @apiResponseRef: CyclesQueryCycleShiftsResponse
* @apiSummary: Query cycle shifts
* @apiParametersRef: CyclesQueryCycleShiftsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryCycleShiftsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryCycleShiftsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryCycleShifts() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryCycleShiftsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		cycles, err := h.CycleService.QueryCycleShifts(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(cycles)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shifts/customer-home-key/release
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesShiftCustomerHomeKeyReleaseRequestBody
 * @apiResponseRef: CyclesCreateShiftCustomerHomeKeyReleaseResponse
 * @apiSummary: Release cycle shift customer home key
 * @apiDescription: Release cycle shift customer home key
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleShiftsCustomerHomeKeyRelease() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesShiftCustomerHomeKeyReleaseRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find shift by id
		shift, err := h.CycleService.FindShiftByID(int64(payload.ShiftID))
		if err != nil || shift == nil {
			return response.ErrorBadRequest(nil, "Shift not found")
		}

		// Release cycle shift customer home key
		data, err := h.CycleService.CreateCycleShiftCustomerHomeKey(&models.CyclesShiftCustomerHomeKeyRequestBody{
			ShiftID:           payload.ShiftID,
			KeyNo:             payload.KeyNo,
			Status:            payload.Status,
			Reason:            payload.Reason,
			AuthenticatedUser: payload.AuthenticatedUser,
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shifts/customer-home-key/not-release
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesShiftCustomerHomeKeyNotReleaseRequestBody
 * @apiResponseRef: CyclesCreateShiftCustomerHomeKeyNotReleaseResponse
 * @apiSummary: Make Not Release cycle shift customer home key
 * @apiDescription: Make Not Release cycle shift customer home key
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleShiftsCustomerHomeKeyNotRelease() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesShiftCustomerHomeKeyNotReleaseRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find shift by id
		shift, err := h.CycleService.FindShiftByID(int64(payload.ShiftID))
		if err != nil || shift == nil {
			return response.ErrorBadRequest(nil, "Shift not found")
		}

		// Release cycle shift customer home key
		data, err := h.CycleService.CreateCycleShiftCustomerHomeKey(&models.CyclesShiftCustomerHomeKeyRequestBody{
			ShiftID:           payload.ShiftID,
			KeyNo:             payload.KeyNo,
			Status:            payload.Status,
			Reason:            payload.Reason,
			AuthenticatedUser: payload.AuthenticatedUser,
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/shifts/customer-home-key
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesShiftCustomerHomeKeyRequestBody
 * @apiResponseRef: CyclesCreateShiftCustomerHomeKeyResponse
 * @apiSummary: Create cycle shift customer home key
 * @apiDescription: Create cycle shift customer home key
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleShiftsCustomerHomeKey() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesShiftCustomerHomeKeyRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find shift by id
		shift, err := h.CycleService.FindShiftByID(int64(payload.ShiftID))
		if err != nil || shift == nil {
			return response.ErrorBadRequest(nil, "Shift not found")
		}

		// Create cycle shift customer home key
		data, err := h.CycleService.CreateCycleShiftCustomerHomeKey(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/shifts/customer-home-keys
* @apiResponseRef: CyclesQueryShiftCustomerHomeKeysResponse
* @apiSummary: Query cycle shifts customer home keys
* @apiParametersRef: CyclesQueryShiftCustomerHomeKeysRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryShiftCustomerHomeKeysNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryShiftCustomerHomeKeysNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryCycleShiftsCustomerHomeKeys() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryShiftCustomerHomeKeysRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query cycle shifts customer home keys
		customerHomeKeys, err := h.CycleService.QueryShiftCustomerHomeKeys(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Return response
		return response.Success(customerHomeKeys)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/start
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitStartRequestBody
 * @apiResponseRef: CyclesCreateVisitStartResponse
 * @apiSummary: Start visit
 * @apiDescription: Start Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitStart() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitStartRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// Start visit
		data, err := h.CycleService.VisitStart(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/end
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitEndRequestBody
 * @apiResponseRef: CyclesCreateVisitEndResponse
 * @apiSummary: End visit
 * @apiDescription: End Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitEnd() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitEndRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// End visit
		data, err := h.CycleService.VisitEnd(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/cancel
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitCancelRequestBody
 * @apiResponseRef: CyclesCreateVisitCancelResponse
 * @apiSummary: Cancel visit
 * @apiDescription: Cancel Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitCancel() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitCancelRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// Cancel visit
		data, err := h.CycleService.VisitCancel(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/delay
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitDelayRequestBody
 * @apiResponseRef: CyclesCreateVisitDelayResponse
 * @apiSummary: Delay visit
 * @apiDescription: Delay Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitDelay() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitDelayRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// Delay visit
		data, err := h.CycleService.VisitDelay(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/pause
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitPauseRequestBody
 * @apiResponseRef: CyclesCreateVisitPauseResponse
 * @apiSummary: Pause visit
 * @apiDescription: Pause Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitPause() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitPauseRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// Pause visit
		data, err := h.CycleService.VisitPause(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/resume
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitResumeRequestBody
 * @apiResponseRef: CyclesCreateVisitResumeResponse
 * @apiSummary: Resume visit
 * @apiDescription: Resume Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitResume() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitResumeRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// Pause Resume
		data, err := h.CycleService.VisitResume(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visits/reactive
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitReactiveRequestBody
 * @apiResponseRef: CyclesCreateVisitReactiveResponse
 * @apiSummary: Reactive visit
 * @apiDescription: Reactive Visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateCycleVisitReactive() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitReactiveRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Set authenticated user
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "Visit not found")
		}

		// Cancel visit
		data, err := h.CycleService.VisitReactive(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visit/assign
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitAssignRequestBody
 * @apiResponseRef: CyclesCreateVisitAssignResponse
 * @apiSummary: Assign visit to staff
 * @apiDescription: Assign visit to staff
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateVisitAssign() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitAssignRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff by id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil || staff == nil {
			return response.ErrorBadRequest(nil, "Staff not found")
		}
		payload.Staff = &domain.CyclePickupShiftStaff{
			ID:        staff.ID,
			UserID:    staff.UserID,
			FirstName: staff.User.FirstName,
			LastName:  staff.User.LastName,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickupShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "cyclePickupShift not found")
		}

		data, err := h.CycleService.AssignVisitToStaff(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateShiftAssignToMe: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visit/swap
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitSwapRequestBody
 * @apiResponseRef: CyclesCreateVisitSwapResponse
 * @apiSummary: Swap visit to another staff
 * @apiDescription: Swap visit to another staff
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateVisitSwap() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitSwapRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find sourceCyclePickupShift by id
		sourceCyclePickupShift, err := h.CycleService.FindPickupShiftByID(int64(payload.SourceCyclePickupShiftID))
		if err != nil || sourceCyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "sourceCyclePickupShift not found")
		}

		// Find targetCyclePickupShift by id
		targetCyclePickupShift, err := h.CycleService.FindPickupShiftByID(int64(payload.TargetCyclePickupShiftID))
		if err != nil || targetCyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "targetCyclePickupShift not found")
		}

		// Swap visits
		sourceVisit, targetVisit, err := h.CycleService.SwapVisits(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateShiftSwap: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(map[string]interface{}{
			"sourceVisit": sourceVisit,
			"targetVisit": targetVisit,
		})
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visit/unplanned
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitUnplannedRequestBody
 * @apiResponseRef: CyclesCreateVisitAssignResponse
 * @apiSummary: Create unplanned visit
 * @apiDescription: Create unplanned visit
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateVisitUnplanned() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitUnplannedRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Find customer by id
		customer, err := h.CustomerService.FindByID(payload.CustomerID)
		if err != nil || customer == nil {
			return response.ErrorBadRequest(nil, "Customer not found")
		}
		payload.Customer = customer

		// Find service by id
		service, err := h.ServiceService.FindByID(int64(payload.ServiceID))
		if err != nil || service == nil {
			return response.ErrorBadRequest(nil, "Service not found")
		}

		// Find service type by id
		serviceType, err := h.ServiceTypeService.FindByID(int64(payload.ServiceTypeID))
		if err != nil || serviceType == nil {
			return response.ErrorBadRequest(nil, "Service type not found")
		}

		// Find staff by id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil || staff == nil {
			return response.ErrorBadRequest(nil, "Staff not found")
		}
		var roles []*domain.CyclePickupShiftStaffRole
		for _, role := range staff.User.Roles {
			roles = append(roles, &domain.CyclePickupShiftStaffRole{
				ID:   role.ID,
				Name: role.Name,
			})
		}
		payload.Staff = &domain.CyclePickupShiftStaff{
			ID:        staff.ID,
			UserID:    staff.UserID,
			FirstName: staff.User.FirstName,
			LastName:  staff.User.LastName,
			Roles:     roles,
		}

		data, err := h.CycleService.CreateUnplannedVisit(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateShiftAssignToMe: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visit/todo
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitTodoRequestBody
 * @apiResponseRef: CyclesCreateVisitTodoResponse
 * @apiSummary: Create visit todos
 * @apiDescription: Create visit todos
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateVisitTodo() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitTodoRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Find cyclePickupShift by id
		cyclePickupShift, err := h.CycleService.FindPickedUpShiftByID(int64(payload.CyclePickupShiftID))
		if err != nil || cyclePickupShift == nil {
			return response.ErrorBadRequest(nil, "cyclePickupShift not found")
		}

		// Create visit todos
		data, err := h.CycleService.CreateVisitTodo(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateVisitTodo: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/visit/todos
* @apiResponseRef: CyclesQueryVisitsTodosResponse
* @apiSummary: Query visit todos
* @apiParametersRef: CyclesQueryVisitsTodosRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryVisitsTodosNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryStaffTypesNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryVisitTodos() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryVisitsTodosRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query visit todos
		todos, err := h.CycleService.QueryVisitTodos(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(todos)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visit/todo/status/{visittodoid}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: CyclesUpdateVisitTodoStatusRequestParams
 * @apiRequestRef: CyclesUpdateVisitTodoStatusRequestBody
 * @apiResponseRef: CyclesUpdateVisitTodoStatusResponse
 * @apiSummary: Update visit todo status
 * @apiDescription: Update visit todo status
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) UpdateVisitTodoStatus() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate request params
		p := mux.Vars(r)
		params := &models.CyclesUpdateVisitTodoStatusRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate request body
		payload := &models.CyclesUpdateVisitTodoStatusRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find visitTodo by id
		visitTodo, err := h.CycleService.FindVisitTodoByID(int64(params.VisitTodoID))
		if err != nil || visitTodo == nil {
			return response.ErrorBadRequest(nil, "Todo not found")
		}

		// Update visit todo status
		todo, err := h.CycleService.UpdateVisitTodoStatus(payload, int64(params.VisitTodoID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CYCLE_BUCKET_NAME, payload.Attachments, int64(todo.ID))
		if len(uploadedFilesMetadata) > 0 {
			updatedTodo, err := h.CycleService.UpdateVisitTodoAttachments(uploadedFilesMetadata, int64(todo.ID))
			if err != nil {
				return response.ErrorInternalServerError(nil, err.Error())
			}
			todo = updatedTodo
		}

		return response.Success(todo)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/visit/sos
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateVisitSosRequestBody
 * @apiResponseRef: CyclesCreateVisitSosResponse
 * @apiSummary: Create visit sos
 * @apiDescription: Create visit sos
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateVisitSos() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateVisitSosRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// TODO: Get timeValueInSeconds from settings
		var timeValueInSeconds = 120

		return response.Success(map[string]interface{}{
			"timeValueInSeconds": timeValueInSeconds,
		})
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/incoming/shift/assign-to-me
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateIncomingShiftAssignToMeRequestBody
 * @apiResponseRef: CyclesCreateIncomingShiftAssignToMeResponse
 * @apiSummary: Assign shift to me in incoming cycle
 * @apiDescription: Assign shift to me in incoming cycle
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateIncomingShiftAssignToMe() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateIncomingShiftAssignToMeRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find staff by id
		staff, err := h.StaffService.FindByID(payload.StaffID)
		if err != nil || staff == nil {
			return response.ErrorBadRequest(nil, "Staff not found")
		}
		payload.Staff = &domain.CyclePickupShiftStaff{
			ID:        staff.ID,
			UserID:    staff.UserID,
			FirstName: staff.User.FirstName,
			LastName:  staff.User.LastName,
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Assign to me
		targetStaffID := user.StaffID
		if targetStaffID == nil {
			return response.ErrorBadRequest(nil, "You are authenticated as a test user without staff. Please create a staff and authenticate with it and try again.")
		}
		data, err := h.CycleService.AssignIncomingShiftsToStaff(payload, int64(*targetStaffID))
		if err != nil {
			log.Printf("CycleHandler.CreateIncomingShiftAssignToMe: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/incoming/shift/swap
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateIncomingShiftSwapRequestBody
 * @apiResponseRef: CyclesCreateIncomingShiftSwapResponse
 * @apiSummary: Swap shift to another staff in incoming cycle
 * @apiDescription: Swap shift to another staff in incoming cycle
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateIncomingShiftSwap() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateIncomingShiftSwapRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find source staff by id
		sourceStaff, err := h.StaffService.FindByID(payload.SourceStaffID)
		if err != nil || sourceStaff == nil {
			return response.ErrorBadRequest(nil, "Source staff not found")
		}
		payload.SourceStaff = &domain.CyclePickupShiftStaff{
			ID:        sourceStaff.ID,
			UserID:    sourceStaff.UserID,
			FirstName: sourceStaff.User.FirstName,
			LastName:  sourceStaff.User.LastName,
		}

		// Find source staff by id
		targetStaff, err := h.StaffService.FindByID(payload.TargetStaffID)
		if err != nil || targetStaff == nil {
			return response.ErrorBadRequest(nil, "Target staff not found")
		}
		payload.TargetStaff = &domain.CyclePickupShiftStaff{
			ID:        targetStaff.ID,
			UserID:    targetStaff.UserID,
			FirstName: targetStaff.User.FirstName,
			LastName:  targetStaff.User.LastName,
		}

		// Find cycle by id
		cycle, err := h.CycleService.FindByID(int64(payload.CycleID))
		if err != nil || cycle == nil {
			return response.ErrorBadRequest(nil, "Cycle not found")
		}

		// Swap shifts
		sourceShifts, targetShifts, err := h.CycleService.SwapIncomingShifts(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateShiftSwap: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(models.CyclesCreateIncomingShiftSwapResponseData{
			SourceShifts: sourceShifts,
			TargetShifts: targetShifts,
		})
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/chats
* @apiResponseRef: CyclesQueryChatsResponse
* @apiSummary: Query cycle chats
* @apiParametersRef: CyclesQueryChatsRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryChatsNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryChatsNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryChats() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryChatsRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query chats
		chats, err := h.CycleService.QueryChats(queries)
		if err != nil {
			log.Printf("CycleHandler.QueryChats: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(chats)
	})
}

/*
* @apiTag: cycle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /cycles/chats/messages
* @apiResponseRef: CyclesQueryChatMessagesResponse
* @apiSummary: Query cycle chat messages
* @apiParametersRef: CyclesQueryChatMessagesRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: CyclesQueryChatMessagesNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: CyclesQueryChatMessagesNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) QueryChatMessages() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.CyclesQueryChatMessagesRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query chat messages
		chatMessages, err := h.CycleService.QueryChatMessages(queries)
		if err != nil {
			log.Printf("CycleHandler.QueryChatMessages: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(chatMessages)
	})
}

/*
 * @apiTag: cycle
 * @apiPath: /cycles/chats/messages
 * @apiMethod: POST
 * @apiStatusCode: 200
 * @apiRequestRef: CyclesCreateChatMessageRequestBody
 * @apiResponseRef: CyclesCreateChatMessageResponse
 * @apiSummary: Create chat message
 * @apiDescription: Create chat message
 * @apiSecurity: apiKeySecurity
 */
func (h *CycleHandler) CreateChatMessage() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Get user from context
		user := h.UserService.GetUserFromContext(r.Context())
		if user == nil {
			return response.ErrorUnAuthorized(nil, "Unauthorized")
		}

		// Validate request body
		payload := &models.CyclesCreateChatMessageRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}
		payload.AuthenticatedUser = &sharedmodels.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		}

		// Create chat message
		chatMessage, err := h.CycleService.CreateChatMessage(payload)
		if err != nil {
			log.Printf("CycleHandler.CreateChatMessage: %s\n", err.Error())
			return response.ErrorInternalServerError(nil, err.Error())
		}

		// Upload attachments to S3
		uploadedFilesMetadata, _ := h.S3Service.UploadFiles(constants.CYCLE_BUCKET_NAME, payload.Attachments, int64(chatMessage.ID))
		if len(uploadedFilesMetadata) > 0 {
			updated, err := h.CycleService.UpdateChatMessageAttachments(nil, uploadedFilesMetadata, int64(chatMessage.ID))
			if err != nil {
				log.Printf("CycleHandler.CreateChatMessage: %s\n", err.Error())
				return response.ErrorInternalServerError(nil, err.Error())
			}
			chatMessage = updated
		}

		// Return created chat message
		return response.Success(chatMessage)
	})
}
