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
