# Prepare these environments before starting the program
```
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/ktaxes?sslmode=disable"
export PORT=8080
export ADMIN_USERNAME=adminTax
export ADMIN_PASSWORD=admin!
```

# How to Start the Program
```
$> docker compose down
$> sudo docker compose up
```
```
prepare all environments
```
```
$> go run main.go
```

# How to Upload CSV File
## Upload with Thunder Client
```
1. Go to Body
2. Select Form
3. Use "taxFile" as the field name and select "taxes.csv" from this repository as the CSV file.
```
### Upload with Command line
```
$> curl -X POST -F "taxFile=@/path/to/taxes.csv" http://localhost:8080/tax/calculations/upload-csv
```
