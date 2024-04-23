package deduction

import (
	_ "database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/phuwanate/assessment-tax/db"

)

type Err struct {
	Message string `json:"message"`
}

func InitDeduction() {
	var err error
	db := database.DB
	createTb := `CREATE TABLE IF NOT EXISTS allowance ( id SERIAL PRIMARY KEY, personalDeduction FLOAT, kReceipt FLOAT);`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	db.QueryRow("INSERT INTO allowance (personalDeduction, kReceipt) values ($1, $2) RETURNING id", 60000, 50000)
}