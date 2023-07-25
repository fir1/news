
# For development it will get up to date proto files current branches proto files
.PHONY: dev
dev:
	go mod tidy
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -d "./" -g "http/server.go"

.PHONY: swagger
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init  -d "./" -g "http/server.go"  --outputTypes "go,json" --overridesFile docs/.swaggo