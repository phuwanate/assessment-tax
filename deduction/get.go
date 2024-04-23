package deduction

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetPersonalAllowance(db *sql.DB) (float64, error) {
	var personalDeduction float64
	err := db.QueryRow("SELECT personalDeduction FROM allowance WHERE id = $1", 1).Scan(&personalDeduction)
	if err != nil {
		return 0, fmt.Errorf("failed to query personal deduction: %v", err)
	}
	return personalDeduction, nil
}

func GetKReceipt(db *sql.DB) (float64, error) {
	var kReceipt float64
	err := db.QueryRow("SELECT kReceipt FROM allowance WHERE id = $1", 1).Scan(&kReceipt)
	if err != nil {
		return 0, fmt.Errorf("failed to query kReceipt maximum: %v", err)
	}
	return kReceipt, nil
}