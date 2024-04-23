package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phuwanate/assessment-tax/csv"
	"github.com/phuwanate/assessment-tax/db"
	"github.com/phuwanate/assessment-tax/deduction"
	"github.com/phuwanate/assessment-tax/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func calculateTax(c echo.Context) error {
	req := new(utils.TaxRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	// Query personal allowance from database
	personalAllowance, err := deduction.GetPersonalAllowance(database.DB)
	if err != nil {
		return err
	}
	maxKReceipt, err := deduction.GetKReceipt(database.DB)
	if err != nil {
		return err
	}

	//Specific limits based on allowance type
	for i := range req.Allowances {
		switch req.Allowances[i].AllowanceType {
		case "k-receipt":
			// Limit k-receipt allowance to maximum of 50,000
			if req.Allowances[i].Amount >  maxKReceipt{
				req.Allowances[i].Amount = maxKReceipt
			}
		case "donation":
			// Limit donation allowance to maximum of 100,000
			if req.Allowances[i].Amount > 100000 {
				req.Allowances[i].Amount = 100000
			}
		default:
			
		}
	}

	// Calculate tax levels and total tax
	taxAmount, taxLevels := utils.CalculateTaxAmount(req.TotalIncome, req.WHT, personalAllowance, req.Allowances)
	if taxAmount < 0 {
		taxAmount = 0
	}

	// Prepare response
	res := utils.TaxResponse{
		Tax:      taxAmount,
		TaxLevel: taxLevels,
	}

	return c.JSON(http.StatusOK, res)
}

func AuthValidator(username, password string, c echo.Context) (bool, error) {
	expectedUsername := os.Getenv("ADMIN_USERNAME")
	expectedPassword := os.Getenv("ADMIN_PASSWORD")

	// Validate the provided username and password against environment variables
	if username == expectedUsername && password == expectedPassword {
		return true, nil
	}

	return false, nil
}

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	deduction.InitDeduction()
	csv.InitCSVTable()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// User Routes
	e.POST("/tax/calculations", calculateTax)
	e.POST("/tax/calculations/upload-csv", csv.CalculateTaxFromCSV)

	//Admin Routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(AuthValidator))
	adminGroup.POST("/deductions/personal", deduction.UpdatePersonalDeduction)
	adminGroup.POST("/deductions/k-receipt", deduction.UpdateMaximumKReceipt)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server on %s", addr)
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("\nBye Bye...")
}
