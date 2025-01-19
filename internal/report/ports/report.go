package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/report/domain"
	"github.com/hoitek/Maja-Service/internal/report/models"
	"github.com/hoitek/Maja-Service/internal/report/types"
)

type ReportService interface {
	QueryArrangementsTable(q *models.ReportsQueryArrangementsTableRequestParams) (*restypes.QueryResponse, error)
	QueryShiftsTable(q *models.ReportsQueryShiftsTableRequestParams) (*restypes.QueryResponse, error)
	QueryVisitsTable(q *models.ReportsQueryVisitsTableRequestParams) (*restypes.QueryResponse, error)
	QueryCustomersTable(q *models.ReportsQueryCustomersTableRequestParams) (*restypes.QueryResponse, error)
	GetShiftSchedulingChart(q *models.ReportsQueryArrangementsTableRequestParams) (*restypes.QueryResponse, error)
	GetShiftDurationAnalysis(q *models.ReportsQueryShiftDurationAnalysisChartRequestParams) (*types.QueryResponseChart, error)
	GetShiftDistributionByCustomer(q *models.ReportsQueryShiftDistributionByCustomersChartRequestParams, period string) (*types.QueryResponseChart, error)
	GetShiftsByStaffAndCustomer(q *models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams) (*types.QueryResponseChart, error)
	GetShiftCountPerDay(q *models.ReportsQueryShiftCountPerDayChartRequestParams) (*types.QueryResponseChart, error)
	GetShiftHeatmap(q *models.ReportsQueryShiftHeatmapChartRequestParams) (*types.QueryResponseChart, error)
	GetCustomerLifecycle(q *models.ReportsQueryCustomerLifecycleChartRequestParams) (*types.QueryResponseChart, error)
	GetCustomerLocation(q *models.ReportsQueryCustomerLocationChartRequestParams) (*types.QueryResponseChart, error)
	GetCustomerAgeGroup(q *models.ReportsQueryCustomerAgeGroupChartRequestParams) (*types.QueryResponseChart, error)
}

type ReportRepositoryPostgresDB interface {
	QueryArrangementsTable(q *models.ReportsQueryArrangementsTableRequestParams) ([]*domain.ReportArrangementTable, error)
	QueryShiftsTable(q *models.ReportsQueryShiftsTableRequestParams) ([]*domain.ReportShiftTable, error)
	QueryVisitsTable(q *models.ReportsQueryVisitsTableRequestParams) ([]*domain.ReportVisitTable, error)
	CountArrangementsTable(queries *models.ReportsQueryArrangementsTableRequestParams) (int64, error)
	CountShiftsTable(queries *models.ReportsQueryShiftsTableRequestParams) (int64, error)
	CountVisitsTable(queries *models.ReportsQueryVisitsTableRequestParams) (int64, error)
	QueryCustomersTable(q *models.ReportsQueryCustomersTableRequestParams) ([]*domain.ReportCustomerTable, error)
	CountCustomersTable(queries *models.ReportsQueryCustomersTableRequestParams) (int64, error)
	GetShiftSchedulingChart(q *models.ReportsQueryArrangementsTableRequestParams) ([]domain.ShiftSchedulingChartData, error)
	GetShiftDurationAnalysis(q *models.ReportsQueryShiftDurationAnalysisChartRequestParams) ([]domain.ShiftDurationAnalysis, error)
	GetShiftDistributionByCustomer(q *models.ReportsQueryShiftDistributionByCustomersChartRequestParams, period string) ([]domain.ShiftDistributionByCustomer, error)
	GetShiftsByStaffAndCustomer(q *models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams) ([]domain.ShiftsByStaffAndCustomer, error)
	GetShiftCountPerDay(q *models.ReportsQueryShiftCountPerDayChartRequestParams) ([]domain.ShiftCountPerDay, error)
	GetShiftHeatmap(q *models.ReportsQueryShiftHeatmapChartRequestParams) ([]domain.ShiftHeatmap, error)
	GetCustomerLifecycle(q *models.ReportsQueryCustomerLifecycleChartRequestParams) ([]domain.CustomerLifecycle, error)
	GetCustomerLocation(q *models.ReportsQueryCustomerLocationChartRequestParams) ([]domain.CustomerLocation, error)
	GetCustomerAgeGroup(q *models.ReportsQueryCustomerAgeGroupChartRequestParams) ([]domain.CustomerAgeGroup, error)
}
