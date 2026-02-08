package repositories

import (
	"aplikasi-kasir/models"
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	report := &models.DailyReport{}

	// Get total revenue for today
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalRevenue)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get total transactions for today
	err = repo.db.QueryRow(`
		SELECT COALESCE(COUNT(*), 0) 
		FROM transactions 
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalTransaksi)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get best-selling product for today
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity)
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return report, nil
}

func (repo *ReportRepository) GetReportByDateRange(startDate, endDate time.Time) (*models.DateRangeReport, error) {
	report := &models.DateRangeReport{}

	// Get total revenue for date range
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&report.TotalRevenue)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get total transactions for date range
	err = repo.db.QueryRow(`
		SELECT COALESCE(COUNT(*), 0) 
		FROM transactions 
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&report.TotalTransaksi)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get best-selling product for date range
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity)
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.id, p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`, startDate, endDate).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return report, nil
}
