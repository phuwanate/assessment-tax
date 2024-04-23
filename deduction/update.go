package deduction

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/phuwanate/assessment-tax/db"
)


func UpdatePersonalDeduction(c echo.Context) error {
	var res float64
	db := database.DB

	// Create a struct to represent the request payload
	type RequestBody struct {
		Amount float64 `json:"amount"`
	}

	// Instantiate a new RequestBody struct
	reqBody := new(RequestBody)

	// Bind the request body to the RequestBody struct
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if (reqBody.Amount > 100000) {
		res = 100000
	} else if (reqBody.Amount < 10000) {
		res = 10000
	}else {
		res = reqBody.Amount
	}
	// Execute SQL statement to update personalDeduction in the database
	stmt, err := db.Prepare(`UPDATE allowance SET personalDeduction = $1`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(res); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	// Prepare the desired response format
	response := map[string]float64{"personalDeduction": res}

	// Return the response with HTTP status OK (200)
	return c.JSON(http.StatusOK, response)
}

func UpdateMaximumKReceipt(c echo.Context) error { 

	var res float64
	db := database.DB

	// Create a struct to represent the request payload
	type RequestBody struct {
		Amount float64 `json:"amount"`
	}

	// Instantiate a new RequestBody struct
	reqBody := new(RequestBody)

	// Bind the request body to the RequestBody struct
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	// Extract the amount from the RequestBody struct
	if (reqBody.Amount > 100000) {
		res = 100000
	} else if (reqBody.Amount < 0) {
		res = 1
	} else {
		res = reqBody.Amount
	}
	// Execute SQL statement to update personalDeduction in the database
	stmt, err := db.Prepare(`UPDATE allowance SET kReceipt = $1`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(res); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	// Prepare the desired response format
	response := map[string]float64{"kReceipt": res}

	// Return the response with HTTP status OK (200)
	return c.JSON(http.StatusOK, response)
}