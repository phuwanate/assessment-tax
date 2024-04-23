package csv

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/phuwanate/assessment-tax/utils"
)

type TaxEntry struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

func CalculateTaxFromCSV(c echo.Context) error {
	// Retrieve uploaded file
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "failed to get file"})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to open file"})
	}
	defer src.Close()

	// Parse CSV file
	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to parse CSV"})
	}

	// Prepare tax calculations
	var taxes []TaxEntry
	var hasError bool // Flag to indicate if an error occurred

	// Iterate over records starting from the second row (skipping the header)
	for i, record := range records {
		if i == 0 {
			// Skip the header row
			continue
		}

		if len(record) < 3 {
			// Log an error for incomplete record
			log.Printf("Incomplete record: %v", record)
			hasError = true
			break // Stop processing the CSV file upon encountering an error
		}

		// Convert CSV record values to float64
		totalIncome, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			// Log an error for invalid totalIncome
			log.Printf("Invalid totalIncome: %v, error: %v", record[0], err)
			hasError = true
			break // Stop processing the CSV file upon encountering an error
		}
		wht, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			// Log an error for invalid wht
			log.Printf("Invalid wht: %v, error: %v", record[1], err)
			hasError = true
			break // Stop processing the CSV file upon encountering an error
		}
		donation, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			// Log an error for invalid donation
			log.Printf("Invalid donation: %v, error: %v", record[2], err)
			hasError = true
			break // Stop processing the CSV file upon encountering an error
		}

		tax := utils.CalculateTaxAmountForCSV(totalIncome, wht, donation)

		// Determine if it's a tax or a refund
		if tax < 0 {
			tax = 0
		} 
		// Append tax entry to result
		taxes = append(taxes, TaxEntry{
			TotalIncome: totalIncome,
			Tax:         tax,
		})
	}

	if hasError {
		// Return a response indicating that an error occurred
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "CSV processing stopped due to errors"})
	}

	// Construct response
	response := map[string][]TaxEntry{"taxes": taxes}

	return c.JSON(http.StatusOK, response)
}