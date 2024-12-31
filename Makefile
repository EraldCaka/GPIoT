run:
	go run gpio-service/cmd/main.go
analog:
	go run gpio-service/integration-test/analog/integration_analog.go
digital:
	go run gpio-service/integration-test/digital/integration_digital.go
up:
	docker compose up
