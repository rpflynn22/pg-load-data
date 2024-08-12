pg-up:
	docker run -d --rm \
		--name pg-load-data \
		-p 5432:5432 \
		-e POSTGRES_HOST_AUTH_METHOD=trust \
		postgres:15.5
	sleep 1
	
pg-down:
	docker kill pg-load-data

apply-schema:
	psql -h localhost -d postgres -U postgres -f ./cases/$(CASE)/schema.sql

insert:
	go run ./cases/$(CASE)

psql:
	psql -h localhost -d postgres -U postgres
