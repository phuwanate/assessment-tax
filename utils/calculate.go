package utils

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/phuwanate/assessment-tax/deduction"
	"github.com/phuwanate/assessment-tax/db"
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

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxResponse struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}

func CalculateTaxAmount(totalIncome, wht, personalAllowance float64, allowances []Allowance) (float64, []TaxLevel) {
	var taxAmount float64
	var taxLevels []TaxLevel
	var level1Tax, level2Tax, level3Tax, level4Tax float64

	// Calculate total amount of allowances
	var totalAllowanceAmount float64
	for _, a := range allowances {
		totalAllowanceAmount += a.Amount
	}

	taxableIncome := (totalIncome - personalAllowance) - totalAllowanceAmount

	// Define tax levels
	taxLevels = append(taxLevels, TaxLevel{Level: "0-150,000", Tax: 0.0})
	if taxableIncome > 150000 {
		level1Tax = (min(taxableIncome, 500000) - 150000) * 0.10
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "150,001-500,000", Tax: level1Tax})
	if taxableIncome > 500000 {
		level2Tax = (min(taxableIncome, 1000000) - 500000) * 0.15
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "500,001-1,000,000", Tax: level2Tax})
	if taxableIncome > 1000000 {
		level3Tax = (min(taxableIncome, 2000000) - 1000000) * 0.20
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "1,000,001-2,000,000", Tax: level3Tax})
	if taxableIncome > 2000000 {
		level4Tax = (taxableIncome - 2000000) * 0.35
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: level4Tax})

	// Calculate total tax
	for _, level := range taxLevels {
		taxAmount += level.Tax
	}

	// Subtract WHT
	taxAmount -= wht

	return taxAmount, taxLevels
}



func CalculateTaxAmountForCSV(totalIncome, wht, donation float64) float64 {
	var taxAmount float64
	var taxLevels []TaxLevel
	var level1Tax, level2Tax, level3Tax, level4Tax float64

	personalAllowance, err := deduction.GetPersonalAllowance(database.DB)
	if err != nil {
		log.Fatal("can't get personal deduction", err)
	}

	if donation > 100000 {
		donation = 100000
	}
	taxableIncome := (totalIncome - personalAllowance) - donation

	// Define tax levels
	taxLevels = append(taxLevels, TaxLevel{Level: "0-150,000", Tax: 0.0})
	if taxableIncome > 150000 {
		level1Tax = (min(taxableIncome, 500000) - 150000) * 0.10
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "150,001-500,000", Tax: level1Tax})
	if taxableIncome > 500000 {
		level2Tax = (min(taxableIncome, 1000000) - 500000) * 0.15
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "500,001-1,000,000", Tax: level2Tax})
	if taxableIncome > 1000000 {
		level3Tax = (min(taxableIncome, 2000000) - 1000000) * 0.20
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "1,000,001-2,000,000", Tax: level3Tax})
	if taxableIncome > 2000000 {
		level4Tax = (taxableIncome - 2000000) * 0.35
	}
	taxLevels = append(taxLevels, TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: level4Tax})

	// Calculate total tax
	for _, level := range taxLevels {
		taxAmount += level.Tax
	}

	// Subtract WHT
	taxAmount -= wht

	return taxAmount
}