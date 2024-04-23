package csv

import (
	_ "database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/phuwanate/assessment-tax/db"
)

func InitCSVTable() {
	var err error
	db := database.DB
	createTb := `
	CREATE TABLE IF NOT EXISTS taxes_csv ( 
		id SERIAL PRIMARY KEY, 
		totalIncome INT, 
		wht INT, 
		donation INT );`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}