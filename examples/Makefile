#!make
include ~/.env

run_postgres:
	go mod tidy
	GSB_SINGLE_LOOP=Y GSB_APR_Db=postgres://admin:pass@0.0.0.0/testdb go run main_postgres.go

run_inmem:
	go mod tidy
	GSB_SINGLE_LOOP=Y go run main_inmem.go

run_retry:
	go mod tidy
	GSB_SINGLE_LOOP=Y go run main_retry.go

run_agent_reply:
	go mod tidy
	GSB_SINGLE_LOOP=Y go run main_agent_reply.go

run_error:
	go mod tidy
	GSB_SINGLE_LOOP=Y go run main_error.go

run_registry:
	go mod tidy
	GSB_SINGLE_LOOP=Y go run main_registry.go

setup_postgres:
	PGPASSWORD=pass psql -h "0.0.0.0" -U admin -c "CREATE DATABASE testdb;"
	cat ./sql/setup_pg.sql | PGPASSWORD=pass psql -h "0.0.0.0" -U admin -d testdb
	PGPASSWORD=pass psql -h "0.0.0.0" -U admin -d testdb -c "SELECT * FROM test_tbl"
