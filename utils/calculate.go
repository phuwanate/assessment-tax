package utils

import (
	"log"
	"net/http"

	"github.com/phuwanate/assessment-tax/db"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type RefundResponse struct {
	Refund float64 `json:"refund"`
}

type TaxLevel struct {
	Tax float64 `json:"tax"`
}

func CalculateTax(c echo.Context) error {
	db := database.DB
	req := new(TaxRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	stmt, err := db.Prepare("SELECT personalDeduction FROM allowance where id=$1")
	if err != nil {
		log.Fatal("can't prepare query one row statment", err)
	}

	rowId := 1
	row := stmt.QueryRow(rowId)
	var personalDeduction float64

	err = row.Scan(&personalDeduction)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}

	// Calculate taxable income
	var taxAmount float64
	var taxLevels []TaxLevel
	var level1Tax, level2Tax, level3Tax, level4Tax float64

	if req.Allowances[0].Amount > 100000 {
		req.Allowances[0].Amount = 100000
	}
	
	//Before tax levels
	taxableIncome := req.TotalIncome - personalDeduction - req.Allowances[0].Amount

	// Define tax levels
	taxLevels = append(taxLevels, TaxLevel{Tax: 0.0})
	if taxableIncome > 150000 {
		level1Tax = (min(taxableIncome, 500000) - 150000) * 0.10
	}
	taxLevels = append(taxLevels, TaxLevel{Tax: level1Tax})
	if taxableIncome > 500000 {
		level2Tax = (min(taxableIncome, 1000000) - 500000) * 0.15
	}
	taxLevels = append(taxLevels, TaxLevel{Tax: level2Tax})
	if taxableIncome > 1000000 {
		level3Tax = (min(taxableIncome, 2000000) - 1000000) * 0.20
	}
	taxLevels = append(taxLevels, TaxLevel{Tax: level3Tax})
	if taxableIncome > 2000000 {
		level4Tax = (taxableIncome - 2000000) * 0.35
	}
	taxLevels = append(taxLevels, TaxLevel{Tax: level4Tax})

	// Calculate total tax
	for _, level := range taxLevels {
		taxAmount += level.Tax
	}

	taxAmount -= req.WHT
	if taxAmount < 0 {
		return c.JSON(http.StatusOK, RefundResponse{
			Refund: -taxAmount,
		})
	} else {
		return c.JSON(http.StatusOK, TaxResponse{
			Tax: taxAmount,
		})
	}
}