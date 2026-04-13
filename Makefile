up:
	docker compose up --build
	
test:
	go test -v ./internal/application/... -coverprofile=coverage.out
