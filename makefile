run:
	@echo "Building and running test instance"
	docker build --tag "kb-auth:local" .
	docker-compose -f docker-compose-local.yml up

run-prod:
	docker-compose -f docker-compose-prod.yml up