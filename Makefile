#!make
include ~/.env

run:
	go mod tidy
	go run main.go

test:
	go test

run_nats:
	docker run -p 4222:4222 -p 8222:8222 nats

run_beanstalkd:
	docker run -d -p 11300:11300 schickling/beanstalkd

run_postgres:
	docker run --name gsb-postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=pass -p 5432:5432 -d postgres
	PGPASSWORD=pass psql -h "0.0.0.0" -U admin -c '\l'

install_sonar_server:
	docker run -d --name sonarqube -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p 9000:9000 sonarqube:latest

sonar_local:
	go test -coverprofile=coverage.out gsb
	docker run \
    	--rm \
		--net=host \
	    -e SONAR_HOST_URL="${SONARQUBE_SERVER}" \
    	-e SONAR_TOKEN="${SONAR_GSB_TOKEN}" \
    	-v ".:/usr/src" \
    	sonarsource/sonar-scanner-cli
