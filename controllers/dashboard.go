package controllers

import (
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
)

type DashboardRes struct{
	TotalSales float64 `json:"total_sales"`
	TotalExpenses float64 `json:"total_expenses"`
	TotalSupplies float64 `json:"total_supplies"`
	TotalSalaryAdvance float64 `json:"total_salary_advance"`
	TotalDippings float64 `json:"total_dippings"`
	TotalSuppliesPaid float64 `json:"total_supplies_paid"`
	TotalDebt float64 `json:"total_debt"`
	TotalEmployees float64 `json:"total_balance"`
}

func Dashboard(c *fiber.Ctx) error {
	var totalSales float64
	var totalExpenses float64
	var totalSupplies float64
	var totalSalaryAdvance float64
	var totalDippings float64
	var totalSuppliesPaid float64
	var totalDebt float64
	var totalEmployees int64

	db.Table("sales").Select("COALESCE(SUM(amount), 0)").Scan(&totalSales)
	db.Table("expenses").Select("COALESCE(SUM(amount), 0)").Scan(&totalExpenses)
	db.Table("supplies").Select("COALESCE(SUM(amount), 0)").Scan(&totalSupplies)
	db.Table("salary_advances").Select("COALESCE(SUM(amount), 0)").Scan(&totalSalaryAdvance)
	db.Table("dippings").Select("COALESCE(SUM(amount), 0)").Scan(&totalDippings)
	db.Table("supplies_payments").Select("COALESCE(SUM(amount), 0)").Scan(&totalSuppliesPaid)
	db.Table("debts").Select("COALESCE(SUM(amount), 0)").Scan(&totalDebt)
	db.Table("employees").Count(&totalEmployees)

	dashboard := DashboardRes{
		TotalSales:         totalSales,
		TotalExpenses:      totalExpenses,
		TotalSupplies:      totalSupplies,
		TotalSalaryAdvance: totalSalaryAdvance,
		TotalDippings:      totalDippings,
		TotalSuppliesPaid:  totalSuppliesPaid,
		TotalDebt:          totalDebt,
		TotalEmployees:     float64(totalEmployees),
	}

	return utils.SuccessResponse(c,"Dashboard", dashboard)
}
