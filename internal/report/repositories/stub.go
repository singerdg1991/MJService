package repositories

import (
	"github.com/hoitek/Maja-Service/internal/report/domain"
	"github.com/hoitek/Maja-Service/internal/report/models"
)

type ReportRepositoryStub struct {
	Reports []*domain.ReportArrangementTable
}

type customerTestCondition struct {
	HasError bool
}

var UserTestCondition *customerTestCondition = &customerTestCondition{}

func NewReportRepositoryStub() *ReportRepositoryStub {
	return &ReportRepositoryStub{
		Reports: []*domain.ReportArrangementTable{},
	}
}

func (r *ReportRepositoryStub) QueryArrangementsTable(q *models.ReportsQueryArrangementsTableRequestParams) ([]*domain.ReportArrangementTable, error) {
	return r.Reports, nil
}

func (r *ReportRepositoryStub) CountArrangementsTable(queries *models.ReportsQueryArrangementsTableRequestParams) (int64, error) {
	return int64(len(r.Reports)), nil
}

func (r *ReportRepositoryStub) QueryShiftsTable(q *models.ReportsQueryShiftsTableRequestParams) ([]*domain.ReportShiftTable, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) CountShiftsTable(queries *models.ReportsQueryShiftsTableRequestParams) (int64, error) {
	return 0, nil
}

func (r *ReportRepositoryStub) QueryVisitsTable(q *models.ReportsQueryVisitsTableRequestParams) ([]*domain.ReportVisitTable, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) CountVisitsTable(queries *models.ReportsQueryVisitsTableRequestParams) (int64, error) {
	return 0, nil
}

func (r *ReportRepositoryStub) QueryCustomersTable(q *models.ReportsQueryCustomersTableRequestParams) ([]*domain.ReportCustomerTable, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) CountCustomersTable(queries *models.ReportsQueryCustomersTableRequestParams) (int64, error) {
	return 0, nil
}

func (r *ReportRepositoryStub) GetShiftSchedulingChart(q *models.ReportsQueryArrangementsTableRequestParams) ([]domain.ShiftSchedulingChartData, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetShiftDurationAnalysis(q *models.ReportsQueryShiftDurationAnalysisChartRequestParams) ([]domain.ShiftDurationAnalysis, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetShiftDistributionByCustomer(q *models.ReportsQueryShiftDistributionByCustomersChartRequestParams, period string) ([]domain.ShiftDistributionByCustomer, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetShiftsByStaffAndCustomer(q *models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams) ([]domain.ShiftsByStaffAndCustomer, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetShiftCountPerDay(q *models.ReportsQueryShiftCountPerDayChartRequestParams) ([]domain.ShiftCountPerDay, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetShiftHeatmap(q *models.ReportsQueryShiftHeatmapChartRequestParams) ([]domain.ShiftHeatmap, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetCustomerLifecycle(q *models.ReportsQueryCustomerLifecycleChartRequestParams) ([]domain.CustomerLifecycle, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetCustomerLocation(q *models.ReportsQueryCustomerLocationChartRequestParams) ([]domain.CustomerLocation, error) {
	return nil, nil
}

func (r *ReportRepositoryStub) GetCustomerAgeGroup(q *models.ReportsQueryCustomerAgeGroupChartRequestParams) ([]domain.CustomerAgeGroup, error) {
	return nil, nil
}
