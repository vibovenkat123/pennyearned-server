module main/expenses

go 1.19

require (
	github.com/jmoiron/sqlx v1.3.5 // direct
	github.com/lib/pq v1.10.7 // direct
)

require (
	github.com/aws/aws-lambda-go v1.38.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.14.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/rs/cors v1.8.3
	go.uber.org/zap v1.24.0
)

require (
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)
