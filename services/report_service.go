package services

import (
	"aplikasi-kasir/models"
	"aplikasi-kasir/repositories"
	"time"
)

type ReportService struct {
	reportRepo *repositories.ReportRepository
}

func NewReportService(reportRepo *repositories.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func (service *ReportService) GetDailyReport() (*models.DailyReport, error) {
	return service.reportRepo.GetDailyReport()
}

func (service *ReportService) GetReportByDateRange(startDate, endDate time.Time) (*models.DateRangeReport, error) {
	return service.reportRepo.GetReportByDateRange(startDate, endDate)
}
