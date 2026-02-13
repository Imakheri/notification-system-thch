
If the database is just created and want to execute a Model AutoMigrate
go run cmd/db/main.go -migrate

If the migration is already done. One table needs to be seeded, please execute
go run cmd/db/main.go -seed

If you want to do those two steps at once, please execute
go run cmd/db/main.go -all