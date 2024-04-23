package deduction

import (
	"log"
	"github.com/phuwanate/assessment-tax/db"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
}

func InitDeduction() {
	db := database.DB
	var err error
	createTb := `CREATE TABLE IF NOT EXISTS allowance ( id SERIAL PRIMARY KEY, personalDeduction INT);`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	db.QueryRow("INSERT INTO  allowance (personalDeduction) values ($1)  RETURNING id", 60000)
}