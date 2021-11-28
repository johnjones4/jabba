module main

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/jackc/pgx/v4 v4.11.0
	github.com/johnjones4/Jabba/core v0.0.0-00010101000000-000000000000
	github.com/swaggest/rest v0.2.16
	github.com/swaggest/usecase v1.1.0
)

replace github.com/johnjones4/Jabba/core => ../core
