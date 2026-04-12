up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f --tail=100
	
test:
	go test -v ./internal/application/... -coverprofile=coverage.out
