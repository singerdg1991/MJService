package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	cPorts "github.com/hoitek/Maja-Service/internal/customer/ports"
	"github.com/hoitek/Maja-Service/internal/report/config"
	"github.com/hoitek/Maja-Service/internal/report/models"
	rPorts "github.com/hoitek/Maja-Service/internal/report/ports"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type ReportHandler struct {
	UserService     uPorts.UserService
	CustomerService cPorts.CustomerService
	ReportService   rPorts.ReportService
	S3Service       s3Ports.S3Service
}

func NewReportHandler(r *mux.Router,
	nService cPorts.CustomerService,
	rService rPorts.ReportService,
	s3Service s3Ports.S3Service,
	uService uPorts.UserService,
) (ReportHandler, error) {
	reportHandler := ReportHandler{
		UserService:     uService,
		CustomerService: nService,
		ReportService:   rService,
		S3Service:       s3Service,
	}
	if r == nil {
		return ReportHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ReportConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ReportConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(uService, []string{}))

	// Add Handlers
	rAuth.Handle("/reports/arrangements/table", reportHandler.QueryArrangementsTable()).Methods(http.MethodGet)
	rAuth.Handle("/reports/arrangements/chart/shift-scheduling", reportHandler.QueryArrangementsChartShiftScheduling()).Methods(http.MethodGet)
	rAuth.Handle("/reports/arrangements/chart/shift-duration-analysis", reportHandler.QueryArrangementsChartShiftDurationAnalysis()).Methods(http.MethodGet)
	rAuth.Handle("/reports/shifts/table", reportHandler.QueryShiftsTable()).Methods(http.MethodGet)
	rAuth.Handle("/reports/visits/table", reportHandler.QueryVisitsTable()).Methods(http.MethodGet)
	rAuth.Handle("/reports/customers/table", reportHandler.QueryCustomersTable()).Methods(http.MethodGet)
	rAuth.Handle("/reports/arrangements/chart/shift-distribution", reportHandler.QueryShiftDistributionByCustomer()).Methods(http.MethodGet)
	rAuth.Handle("/reports/shifts/chart/shifts-by-staff-customer", reportHandler.QueryShiftsByStaffAndCustomer()).Methods(http.MethodGet)
	rAuth.Handle("/reports/shifts/chart/shifts-per-day", reportHandler.QueryShiftCountPerDay()).Methods(http.MethodGet)
	rAuth.Handle("/reports/shifts/chart/shift-by-day-and-hour", reportHandler.QueryShiftHeatmap()).Methods(http.MethodGet)
	rAuth.Handle("/reports/customers/chart/lifecycle", reportHandler.QueryCustomerLifecycle()).Methods(http.MethodGet)
	rAuth.Handle("/reports/customers/chart/location", reportHandler.QueryCustomerLocation()).Methods(http.MethodGet)
	rAuth.Handle("/reports/customers/chart/age-group", reportHandler.QueryCustomerAgeGroup()).Methods(http.MethodGet)

	return reportHandler, nil
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/arrangements/table
* @apiResponseRef: ReportsQueryArrangementsTableResponse
* @apiSummary: Query reports arrangements table
* @apiParametersRef: ReportsQueryArrangementsTableRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryArrangementsTableNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryArrangementsTableNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryArrangementsTable() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryArrangementsTableRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		data, err := h.ReportService.QueryArrangementsTable(queries)
		if err != nil {
			log.Printf("Error querying arrangements table: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/shifts/table
* @apiResponseRef: ReportsQueryShiftsTableResponse
* @apiSummary: Query reports shifts table
* @apiParametersRef: ReportsQueryShiftsTableRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryShiftsTableNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryShiftsTableNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryShiftsTable() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryShiftsTableRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		data, err := h.ReportService.QueryShiftsTable(queries)
		if err != nil {
			log.Printf("Error querying shifts table: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/visits/table
* @apiResponseRef: ReportsQueryVisitsTableResponse
* @apiSummary: Query reports visits table
* @apiParametersRef: ReportsQueryVisitsTableRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryVisitsTableNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryVisitsTableNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryVisitsTable() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryVisitsTableRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		data, err := h.ReportService.QueryVisitsTable(queries)
		if err != nil {
			log.Printf("Error querying visits table: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/customers/table
* @apiResponseRef: ReportsQueryCustomersTableResponse
* @apiSummary: Query reports customers table
* @apiParametersRef: ReportsQueryCustomersTableRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryCustomersTableNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryCustomersTableNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryCustomersTable() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryCustomersTableRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		data, err := h.ReportService.QueryCustomersTable(queries)
		if err != nil {
			log.Printf("Error querying customers table: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/arrangements/chart/shift-scheduling
* @apiResponseRef: ReportsQueryArrangementsChartShiftSchedulingResponse
* @apiSummary: Query reports arrangements chart shift scheduling
* @apiParametersRef: ReportsQueryArrangementsTableRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryArrangementsChartShiftSchedulingNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryArrangementsChartShiftSchedulingNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryArrangementsChartShiftScheduling() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryArrangementsTableRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		data, err := h.ReportService.GetShiftSchedulingChart(queries)
		if err != nil {
			log.Printf("Error querying arrangements chart: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/arrangements/chart/shift-duration-analysis
* @apiResponseRef: ReportsQueryShiftDurationAnalysisResponse
* @apiSummary: Query shift duration analysis chart
* @apiParametersRef: ReportsQueryShiftDurationAnalysisChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryShiftDurationAnalysisNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryShiftDurationAnalysisNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryArrangementsChartShiftDurationAnalysis() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryShiftDurationAnalysisChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		data, err := h.ReportService.GetShiftDurationAnalysis(queries)
		if err != nil {
			log.Printf("Error querying shift duration analysis: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/arrangements/chart/shift-distribution
* @apiResponseRef: ReportsQueryShiftDistributionByCustomerChartDistributionResponse
* @apiSummary: Query shift distribution by customer pie chart
* @apiParametersRef: ReportsQueryShiftDistributionByCustomersChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryShiftDistributionByCustomerChartDistributionNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ReportsQueryShiftDistributionByCustomerChartDistributionNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryShiftDistributionByCustomer() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryShiftDistributionByCustomersChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Get period from query params (default to daily if not specified)
		period := r.URL.Query().Get("period")
		if period != "daily" && period != "weekly" {
			period = "daily"
		}

		// Calculate data
		data, err := h.ReportService.GetShiftDistributionByCustomer(queries, period)
		if err != nil {
			log.Printf("Error querying shift distribution: %v\n", err.Error())
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again")
		}

		// Return response
		return response.Success(data)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/shifts/chart/shifts-by-staff-customer
* @apiResponseRef: ReportsQueryShiftsByStaffAndCustomerChartResponse
* @apiSummary: Query shifts by staff and customer stacked bar chart
* @apiParametersRef: ReportsQueryShiftsByStaffAndCustomerChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryShiftsByStaffAndCustomerChartNotFoundResponse
* @api404ResponseDescription: No data found
* @api500ResponseRef: ReportsQueryShiftsByStaffAndCustomerChartNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryShiftsByStaffAndCustomer() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		result, err := h.ReportService.GetShiftsByStaffAndCustomer(queries)
		if err != nil {
			log.Printf("Error querying shifts by staff and customer: %v\n", err.Error())
			return response.ErrorInternalServerError(err, "Something went wrong, please try again")
		}

		return response.Success(result)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/shifts/chart/shifts-per-day
* @apiResponseRef: ReportsQueryShiftCountPerDayChartResponse
* @apiSummary: Query shifts count per day of week
* @apiParametersRef: ReportsQueryShiftCountPerDayChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryShiftCountPerDayChartNotFoundResponse
* @api404ResponseDescription: No data found
* @api500ResponseRef: ReportsQueryShiftCountPerDayChartNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryShiftCountPerDay() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryShiftCountPerDayChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		result, err := h.ReportService.GetShiftCountPerDay(queries)
		if err != nil {
			log.Printf("Error querying shifts per day: %v\n", err.Error())
			return response.ErrorInternalServerError(err, "Something went wrong, please try again")
		}

		return response.Success(result)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/shifts/chart/shift-by-day-and-hour
* @apiResponseRef: ReportsQueryShiftHeatmapChartResponse
* @apiSummary: Query shifts heatmap by day and hour
* @apiParametersRef: ReportsQueryShiftHeatmapChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryShiftHeatmapChartNotFoundResponse
* @api404ResponseDescription: No data found
* @api500ResponseRef: ReportsQueryShiftHeatmapChartNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryShiftHeatmap() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryShiftHeatmapChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		result, err := h.ReportService.GetShiftHeatmap(queries)
		if err != nil {
			log.Printf("Error querying shift heatmap: %v\n", err.Error())
			return response.ErrorInternalServerError(err, "Something went wrong, please try again")
		}

		return response.Success(result)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/customers/chart/lifecycle
* @apiResponseRef: ReportsQueryCustomerLifecycleChartResponse
* @apiSummary: Query customer lifecycle funnel data
* @apiParametersRef: ReportsQueryCustomerLifecycleChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryCustomerLifecycleChartNotFoundResponse
* @api404ResponseDescription: No data found
* @api500ResponseRef: ReportsQueryCustomerLifecycleChartNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryCustomerLifecycle() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryCustomerLifecycleChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		result, err := h.ReportService.GetCustomerLifecycle(queries)
		if err != nil {
			log.Printf("Error querying customer lifecycle: %v\n", err.Error())
			return response.ErrorInternalServerError(err, "Something went wrong, please try again")
		}

		return response.Success(result)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/customers/chart/location
* @apiResponseRef: ReportsQueryCustomerLocationChartResponse
* @apiSummary: Query customer count by location
* @apiParametersRef: ReportsQueryCustomerLocationChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryCustomerLocationChartNotFoundResponse
* @api404ResponseDescription: No data found
* @api500ResponseRef: ReportsQueryCustomerLocationChartNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryCustomerLocation() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryCustomerLocationChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		result, err := h.ReportService.GetCustomerLocation(queries)
		if err != nil {
			log.Printf("Error querying customer locations: %v\n", err.Error())
			return response.ErrorInternalServerError(err, "Something went wrong, please try again")
		}

		return response.Success(result)
	})
}

/*
* @apiTag: report
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /reports/customers/chart/age-group
* @apiResponseRef: ReportsQueryCustomerAgeGroupChartResponse
* @apiSummary: Query customer count by age group
* @apiParametersRef: ReportsQueryCustomerAgeGroupChartRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ReportsQueryCustomerAgeGroupChartNotFoundResponse
* @api404ResponseDescription: No data found
* @api500ResponseRef: ReportsQueryCustomerAgeGroupChartNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ReportHandler) QueryCustomerAgeGroup() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ReportsQueryCustomerAgeGroupChartRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Calculate data
		result, err := h.ReportService.GetCustomerAgeGroup(queries)
		if err != nil {
			log.Printf("Error querying customer age groups: %v\n", err.Error())
			return response.ErrorInternalServerError(err, "Something went wrong, please try again")
		}

		return response.Success(result)
	})
}
