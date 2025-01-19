package service

import (
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/report/constants"
	"github.com/hoitek/Maja-Service/internal/report/models"
	"github.com/hoitek/Maja-Service/internal/report/ports"
	"github.com/hoitek/Maja-Service/internal/report/types"
	"github.com/hoitek/Maja-Service/storage"

	"math"
)

type ReportService struct {
	PostgresRepository ports.ReportRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewReportService(pDB ports.ReportRepositoryPostgresDB, m *storage.MinIO) ReportService {
	go minio.SetupMinIOStorage(constants.REPORT_BUCKET_NAME, m)
	return ReportService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ReportService) QueryArrangementsTable(q *models.ReportsQueryArrangementsTableRequestParams) (*restypes.QueryResponse, error) {
	items, err := s.PostgresRepository.QueryArrangementsTable(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountArrangementsTable(q)
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 10, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
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

func (s *ReportService) QueryShiftsTable(q *models.ReportsQueryShiftsTableRequestParams) (*restypes.QueryResponse, error) {
	items, err := s.PostgresRepository.QueryShiftsTable(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountShiftsTable(q)
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 10, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
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

func (s *ReportService) QueryVisitsTable(q *models.ReportsQueryVisitsTableRequestParams) (*restypes.QueryResponse, error) {
	items, err := s.PostgresRepository.QueryVisitsTable(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountVisitsTable(q)
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 10, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
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

func (s *ReportService) QueryCustomersTable(q *models.ReportsQueryCustomersTableRequestParams) (*restypes.QueryResponse, error) {
	items, err := s.PostgresRepository.QueryCustomersTable(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountCustomersTable(q)
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 10, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
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

func (s *ReportService) GetShiftSchedulingChart(q *models.ReportsQueryArrangementsTableRequestParams) (*restypes.QueryResponse, error) {
	items, err := s.PostgresRepository.GetShiftSchedulingChart(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.CountArrangementsTable(q)
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

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *ReportService) GetShiftDurationAnalysis(q *models.ReportsQueryShiftDurationAnalysisChartRequestParams) (*types.QueryResponseChart, error) {
	items, err := s.PostgresRepository.GetShiftDurationAnalysis(q)
	if err != nil {
		return nil, err
	}
	return &types.QueryResponseChart{
		Items: items,
	}, nil
}

func (s *ReportService) GetShiftDistributionByCustomer(q *models.ReportsQueryShiftDistributionByCustomersChartRequestParams, period string) (*types.QueryResponseChart, error) {
	items, err := s.PostgresRepository.GetShiftDistributionByCustomer(q, period)
	if err != nil {
		return nil, err
	}
	return &types.QueryResponseChart{
		Items: items,
	}, nil
}

func (s *ReportService) GetShiftsByStaffAndCustomer(q *models.ReportsQueryShiftsByStaffAndCustomerChartRequestParams) (*types.QueryResponseChart, error) {
	items, err := s.PostgresRepository.GetShiftsByStaffAndCustomer(q)
	if err != nil {
		return nil, err
	}

	return &types.QueryResponseChart{
		Items: items,
	}, nil
}

func (s *ReportService) GetShiftCountPerDay(q *models.ReportsQueryShiftCountPerDayChartRequestParams) (*types.QueryResponseChart, error) {
	items, err := s.PostgresRepository.GetShiftCountPerDay(q)
	if err != nil {
		return nil, err
	}

	return &types.QueryResponseChart{
		Items: items,
	}, nil
}

func (s *ReportService) GetShiftHeatmap(q *models.ReportsQueryShiftHeatmapChartRequestParams) (*types.QueryResponseChart, error) {
	items, err := s.PostgresRepository.GetShiftHeatmap(q)
	if err != nil {
		return nil, err
	}

	return &types.QueryResponseChart{
		Items: items,
	}, nil
}

func (s *ReportService) GetCustomerLifecycle(q *models.ReportsQueryCustomerLifecycleChartRequestParams) (*types.QueryResponseChart, error) {
	data, err := s.PostgresRepository.GetCustomerLifecycle(q)
	if err != nil {
		return nil, err
	}

	return &types.QueryResponseChart{
		Items: data,
	}, nil
}

func (s *ReportService) GetCustomerLocation(q *models.ReportsQueryCustomerLocationChartRequestParams) (*types.QueryResponseChart, error) {
	data, err := s.PostgresRepository.GetCustomerLocation(q)
	if err != nil {
		return nil, err
	}

	return &types.QueryResponseChart{
		Items: data,
	}, nil
}

func (s *ReportService) GetCustomerAgeGroup(q *models.ReportsQueryCustomerAgeGroupChartRequestParams) (*types.QueryResponseChart, error) {
	data, err := s.PostgresRepository.GetCustomerAgeGroup(q)
	if err != nil {
		return nil, err
	}

	return &types.QueryResponseChart{
		Items: data,
	}, nil
}
